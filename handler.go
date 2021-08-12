package rerpc

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	// Always advertise that reRPC accepts gzip compression.
	acceptEncodingValue    = strings.Join([]string{CompressionGzip, CompressionIdentity}, ",")
	acceptPostValueDefault = strings.Join(
		[]string{TypeDefaultGRPC, TypeProtoGRPC, TypeJSON},
		",",
	)
	acceptPostValueWithoutJSON = strings.Join(
		[]string{TypeDefaultGRPC, TypeProtoGRPC},
		",",
	)
)

type handlerCfg struct {
	DisableGzipResponse bool
	DisableTwirp        bool
	MaxRequestBytes     int64
	Registrar           *Registrar
	Interceptor         Interceptor
	Package             string
	Service             string
	Method              string
}

// A HandlerOption configures a Handler.
//
// In addition to any options grouped in the documentation below, remember that
// Registrars and Options are also valid HandlerOptions.
type HandlerOption interface {
	applyToHandler(*handlerCfg)
}

type serveTwirpOption struct {
	Disable bool
}

func (o *serveTwirpOption) applyToHandler(cfg *handlerCfg) {
	cfg.DisableTwirp = o.Disable
}

// ServeTwirp enables or disables support for Twirp's JSON and protobuf
// formats. Disable Twirp if you only want your handlers to speak the gRPC
// protocol.
//
// By default, handlers support Twirp.
func ServeTwirp(enable bool) HandlerOption {
	return &serveTwirpOption{!enable}
}

// A Handler is the server-side implementation of a single RPC defined by a
// protocol buffer service. It's the interface between the reRPC library and
// the code generated by the reRPC protoc plugin; most users won't ever need to
// deal with it directly.
//
// To see an example of how Handler is used in the generated code, see the
// internal/pingpb/v0 package.
type Handler struct {
	stype          StreamType
	config         handlerCfg
	implementation func(context.Context, StreamFunc)
}

// NewHandler constructs a Handler. The supplied package, service, and method
// names must be protobuf identifiers. For example, a handler for the URL
// "/acme.foo.v1.FooService/Bar" would have package "acme.foo.v1", service
// "FooService", and method "Bar".
//
// Remember that NewHandler is usually called from generated code - most users
// won't need to deal with protobuf identifiers directly.
func NewHandler(
	stype StreamType,
	pkg, service, method string,
	implementation func(context.Context, StreamFunc),
	opts ...HandlerOption,
) *Handler {
	cfg := handlerCfg{
		Package: pkg,
		Service: service,
		Method:  method,
	}
	for _, opt := range opts {
		opt.applyToHandler(&cfg)
	}
	if reg := cfg.Registrar; reg != nil {
		reg.register(cfg.Package, cfg.Service)
	}
	return &Handler{
		stype:          stype,
		config:         cfg,
		implementation: implementation,
	}
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We don't need to defer functions  to close the request body or read to
	// EOF: the stream we construct later on already does that, and we only
	// return early when dealing with misbehaving clients. In those cases, it's
	// okay if we can't re-use the connection.

	isBidi := (h.stype & StreamTypeBidirectional) == StreamTypeBidirectional
	if isBidi && r.ProtoMajor < 2 {
		w.WriteHeader(http.StatusHTTPVersionNotSupported)
		io.WriteString(w, "bidirectional streaming requires HTTP/2")
		return
	}
	if r.Method != http.MethodPost {
		// grpc-go returns a 500 here, but interoperability with non-gRPC HTTP
		// clients is better if we return a 405.
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctype := r.Header.Get("Content-Type")
	if (ctype == TypeJSON || ctype == TypeProtoTwirp) && h.config.DisableTwirp {
		w.Header().Set("Accept-Post", acceptPostValueWithoutJSON)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	if ctype != TypeDefaultGRPC && ctype != TypeProtoGRPC && ctype != TypeProtoTwirp && ctype != TypeJSON {
		// grpc-go returns 500, but the spec recommends 415.
		// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests
		w.Header().Set("Accept-Post", acceptPostValueDefault)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	spec := &Specification{
		Type:                h.stype,
		Package:             h.config.Package,
		Service:             h.config.Service,
		Method:              h.config.Method,
		Path:                r.URL.Path,
		ContentType:         ctype,
		RequestCompression:  CompressionIdentity,
		ResponseCompression: CompressionIdentity,
		ReadMaxBytes:        h.config.MaxRequestBytes,
	}

	// We need to parse metadata before entering the interceptor stack, but we'd
	// like to report errors to the client in a format they understand (if
	// possible). We'll collect any such errors here and use them to
	// short-circuit early later on.
	//
	// NB, future refactorings will need to take care to avoid typed nils here.
	var failed *Error

	timeout, err := parseTimeout(r.Header.Get("Grpc-Timeout"))
	if err != nil && err != errNoTimeout {
		// Errors here indicate that the client sent an invalid timeout header, so
		// the error text is safe to send back.
		failed = wrap(CodeInvalidArgument, err)
	} else if err == nil {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		r = r.WithContext(ctx)
	} // else err == errNoTimeout, nothing to do

	if spec.ContentType == TypeJSON || spec.ContentType == TypeProtoTwirp {
		if r.Header.Get("Content-Encoding") == "gzip" {
			spec.RequestCompression = CompressionGzip
		}
		// TODO: Actually parse Accept-Encoding instead of this hackery.
		if !h.config.DisableGzipResponse && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			spec.ResponseCompression = CompressionGzip
		}
	} else {
		spec.RequestCompression = CompressionIdentity
		if me := r.Header.Get("Grpc-Encoding"); me != "" {
			switch me {
			case CompressionIdentity:
				spec.RequestCompression = CompressionIdentity
			case CompressionGzip:
				spec.RequestCompression = CompressionGzip
			default:
				// Per https://github.com/grpc/grpc/blob/master/doc/compression.md, we
				// should return CodeUnimplemented and specify acceptable compression(s)
				// (in addition to setting the Grpc-Accept-Encoding header).
				if failed == nil {
					failed = errorf(
						CodeUnimplemented,
						"unknown compression %q: accepted grpc-encoding values are %v",
						me, acceptEncodingValue,
					)
				}
			}
		}
		// Follow https://github.com/grpc/grpc/blob/master/doc/compression.md.
		// (The grpc-go implementation doesn't read the "grpc-accept-encoding" header
		// and doesn't support compression method asymmetry.)
		spec.ResponseCompression = spec.RequestCompression
		if h.config.DisableGzipResponse {
			spec.ResponseCompression = CompressionIdentity
		} else if mae := r.Header.Get("Grpc-Accept-Encoding"); mae != "" {
			for _, enc := range strings.FieldsFunc(mae, splitOnCommasAndSpaces) {
				switch enc {
				case CompressionGzip:
					spec.ResponseCompression = CompressionGzip
					// prefer gzip, so no continue
				case CompressionIdentity:
					spec.ResponseCompression = CompressionIdentity
					continue
				default:
					continue
				}
				break
			}
		}
	}

	// We should write any remaining headers here, since: (a) the implementation
	// may write to the body, thereby sending the headers, and (b) interceptors
	// should be able to see this data.
	w.Header().Set("Content-Type", spec.ContentType)
	if spec.ContentType != TypeJSON && spec.ContentType != TypeProtoTwirp {
		w.Header().Set("Grpc-Accept-Encoding", acceptEncodingValue)
		w.Header().Set("Grpc-Encoding", spec.ResponseCompression)
		// Every gRPC response will have these trailers.
		w.Header().Add("Trailer", "Grpc-Status")
		w.Header().Add("Trailer", "Grpc-Message")
		w.Header().Add("Trailer", "Grpc-Status-Details-Bin")
	}

	// Unlike gRPC, Twirp manages compression using the standard HTTP mechanisms.
	// Since they apply to the whole stream, it's easiest to handle it here.
	var requestBody io.Reader = r.Body
	if spec.ContentType == TypeJSON || spec.ContentType == TypeProtoTwirp {
		if spec.RequestCompression == CompressionGzip {
			gr, err := gzip.NewReader(requestBody)
			if err != nil && failed == nil {
				failed = errorf(CodeInvalidArgument, "can't read gzipped body: %w", err)
			} else if err == nil {
				defer gr.Close()
				requestBody = gr
			}
		}
		// Checking Content-Encoding ensures that some other user-supplied
		// middleware isn't already compressing the response.
		if spec.ResponseCompression == CompressionGzip && w.Header().Get("Content-Encoding") == "" {
			w.Header().Set("Content-Encoding", "gzip")
			gw := getGzipWriter(w)
			defer putGzipWriter(gw)
			w = &gzipResponseWriter{ResponseWriter: w, gw: gw}
		}
	}

	ctx := NewHandlerContext(r.Context(), *spec, r.Header, w.Header())
	sf := StreamFunc(func(ctx context.Context) Stream {
		return newServerStream(
			ctx,
			w,
			&readCloser{Reader: requestBody, Closer: r.Body},
			spec.ContentType,
			h.config.MaxRequestBytes,
			spec.ResponseCompression == CompressionGzip,
		)
	})
	if failed != nil {
		stream := sf(ctx)
		_ = stream.CloseReceive()
		_ = stream.CloseSend(failed)
		return
	}
	h.implementation(ctx, sf)
}

// Path returns the URL pattern to use when registering this handler. It's used
// by the generated code.
func (h *Handler) Path() string {
	if h.config.Package == "" && h.config.Service == "" && h.config.Method == "" {
		return "/"
	}
	return fmt.Sprintf("/%s.%s/%s", h.config.Package, h.config.Service, h.config.Method)
}

// ServicePath returns the URL pattern for the protobuf service. It's used by
// the generated code.
func (h *Handler) ServicePath() string {
	if h.config.Package == "" && h.config.Service == "" && h.config.Method == "" {
		return "/"
	}
	return fmt.Sprintf("/%s.%s/", h.config.Package, h.config.Service)
}

func splitOnCommasAndSpaces(c rune) bool {
	return c == ',' || c == ' '
}

type readCloser struct {
	io.Reader
	io.Closer
}
