// Code generated by protoc-gen-go-rerpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-rerpc v0.0.1
// - protoc             v3.17.3
// source: internal/ping/v1test/ping.proto

package pingpb

import (
	context "context"
	errors "errors"
	rerpc "github.com/rerpc/rerpc"
	proto "google.golang.org/protobuf/proto"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the
// rerpc package are compatible. If you get a compiler error that this constant
// isn't defined, this code was generated with a version of rerpc newer than the
// one compiled into your binary. You can fix the problem by either regenerating
// this code with an older version of rerpc or updating the rerpc version
// compiled into your binary.
const _ = rerpc.SupportsCodeGenV0 // requires reRPC v0.0.1 or later

// PingServiceClientReRPC is a client for the internal.ping.v1test.PingService
// service.
type PingServiceClientReRPC interface {
	Ping(ctx context.Context, req *PingRequest, opts ...rerpc.CallOption) (*PingResponse, error)
	Fail(ctx context.Context, req *FailRequest, opts ...rerpc.CallOption) (*FailResponse, error)
}

type pingServiceClientReRPC struct {
	doer    rerpc.Doer
	baseURL string
	options []rerpc.CallOption
}

// NewPingServiceClientReRPC constructs a client for the
// internal.ping.v1test.PingService service. Call options passed here apply to
// all calls made with this client.
//
// The URL supplied here should be the base URL for the gRPC server (e.g.,
// https://api.acme.com or https://acme.com/grpc).
func NewPingServiceClientReRPC(baseURL string, doer rerpc.Doer, opts ...rerpc.CallOption) PingServiceClientReRPC {
	return &pingServiceClientReRPC{
		baseURL: strings.TrimRight(baseURL, "/"),
		doer:    doer,
		options: opts,
	}
}

func (c *pingServiceClientReRPC) mergeOptions(opts []rerpc.CallOption) []rerpc.CallOption {
	merged := make([]rerpc.CallOption, 0, len(c.options)+len(opts))
	for _, o := range c.options {
		merged = append(merged, o)
	}
	for _, o := range opts {
		merged = append(merged, o)
	}
	return merged
}

// Ping calls internal.ping.v1test.PingService.Ping. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPC) Ping(ctx context.Context, req *PingRequest, opts ...rerpc.CallOption) (*PingResponse, error) {
	merged := c.mergeOptions(opts)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeUnary,
		c.baseURL,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Ping",                 // protobuf method
		merged...,
	)
	wrapped := rerpc.Func(func(ctx context.Context, msg proto.Message) (proto.Message, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res PingResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic := rerpc.ConfiguredCallInterceptor(merged...); ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*PingResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.ping.v1test.PingResponse, got %v", res.ProtoReflect().Descriptor().FullName())
	}
	return typed, nil
}

// Fail calls internal.ping.v1test.PingService.Fail. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPC) Fail(ctx context.Context, req *FailRequest, opts ...rerpc.CallOption) (*FailResponse, error) {
	merged := c.mergeOptions(opts)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeUnary,
		c.baseURL,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Fail",                 // protobuf method
		merged...,
	)
	wrapped := rerpc.Func(func(ctx context.Context, msg proto.Message) (proto.Message, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res FailResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic := rerpc.ConfiguredCallInterceptor(merged...); ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*FailResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.ping.v1test.FailResponse, got %v", res.ProtoReflect().Descriptor().FullName())
	}
	return typed, nil
}

// PingServiceReRPC is a server for the internal.ping.v1test.PingService
// service. To make sure that adding methods to this protobuf service doesn't
// break all implementations of this interface, all implementations must embed
// UnimplementedPingServiceReRPC.
//
// By default, recent versions of grpc-go have a similar forward compatibility
// requirement. See https://github.com/grpc/grpc-go/issues/3794 for a longer
// discussion.
type PingServiceReRPC interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Fail(context.Context, *FailRequest) (*FailResponse, error)
	Sum(context.Context, *PingServiceReRPC_SumServer) error
	CountUp(context.Context, *CountUpRequest, *PingServiceReRPC_CountUpServer) error
	CumSum(context.Context, *PingServiceReRPC_CumSumServer) error
	mustEmbedUnimplementedPingServiceReRPC()
}

// NewPingServiceHandlerReRPC wraps the service implementation in an HTTP
// handler. It returns the handler and the path on which to mount it.
func NewPingServiceHandlerReRPC(svc PingServiceReRPC, opts ...rerpc.HandlerOption) (string, *http.ServeMux) {
	mux := http.NewServeMux()
	ic := rerpc.ConfiguredHandlerInterceptor(opts...)

	pingFunc := rerpc.Func(func(ctx context.Context, req proto.Message) (proto.Message, error) {
		typed, ok := req.(*PingRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.ping.v1test.PingService.Ping with a %v",
				req.ProtoReflect().Descriptor().FullName(),
			)
		}
		return svc.Ping(ctx, typed)
	})
	if ic != nil {
		pingFunc = ic.Wrap(pingFunc)
	}
	ping := rerpc.NewHandler(
		rerpc.StreamTypeUnary,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Ping",                 // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			stream := sf(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req PingRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := pingFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
		},
		opts...,
	)
	mux.Handle(ping.Path(), ping)

	failFunc := rerpc.Func(func(ctx context.Context, req proto.Message) (proto.Message, error) {
		typed, ok := req.(*FailRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.ping.v1test.PingService.Fail with a %v",
				req.ProtoReflect().Descriptor().FullName(),
			)
		}
		return svc.Fail(ctx, typed)
	})
	if ic != nil {
		failFunc = ic.Wrap(failFunc)
	}
	fail := rerpc.NewHandler(
		rerpc.StreamTypeUnary,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Fail",                 // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			stream := sf(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req FailRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := failFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
		},
		opts...,
	)
	mux.Handle(fail.Path(), fail)

	sum := rerpc.NewHandler(
		rerpc.StreamTypeClient,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Sum",                  // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewPingServiceReRPC_SumServer(stream)
			err := svc.Sum(stream.Context(), typed)
			_ = stream.CloseReceive()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	mux.Handle(sum.Path(), sum)

	countUp := rerpc.NewHandler(
		rerpc.StreamTypeServer,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"CountUp",              // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewPingServiceReRPC_CountUpServer(stream)
			var req CountUpRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseReceive()
				_ = stream.CloseSend(err)
				return
			}
			if err := stream.CloseReceive(); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			err := svc.CountUp(stream.Context(), &req, typed)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	mux.Handle(countUp.Path(), countUp)

	cumSum := rerpc.NewHandler(
		rerpc.StreamTypeBidirectional,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"CumSum",               // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewPingServiceReRPC_CumSumServer(stream)
			err := svc.CumSum(stream.Context(), typed)
			_ = stream.CloseReceive()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	mux.Handle(cumSum.Path(), cumSum)

	// Respond to unknown protobuf methods with gRPC and Twirp's 404 equivalents.
	mux.Handle("/", rerpc.NewBadRouteHandler(opts...))

	return cumSum.ServicePath(), mux
}

var _ PingServiceReRPC = (*UnimplementedPingServiceReRPC)(nil) // verify interface implementation

// UnimplementedPingServiceReRPC returns CodeUnimplemented from all methods. To
// maintain forward compatibility, all implementations of PingServiceReRPC must
// embed UnimplementedPingServiceReRPC.
type UnimplementedPingServiceReRPC struct{}

func (UnimplementedPingServiceReRPC) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.Ping isn't implemented")
}

func (UnimplementedPingServiceReRPC) Fail(context.Context, *FailRequest) (*FailResponse, error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.Fail isn't implemented")
}

func (UnimplementedPingServiceReRPC) Sum(context.Context, *PingServiceReRPC_SumServer) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.Sum isn't implemented")
}

func (UnimplementedPingServiceReRPC) CountUp(context.Context, *CountUpRequest, *PingServiceReRPC_CountUpServer) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.CountUp isn't implemented")
}

func (UnimplementedPingServiceReRPC) CumSum(context.Context, *PingServiceReRPC_CumSumServer) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.CumSum isn't implemented")
}

func (UnimplementedPingServiceReRPC) mustEmbedUnimplementedPingServiceReRPC() {}

// PingServiceReRPC_SumServer is the server-side stream for the
// internal.ping.v1test.PingService.Sum procedure.
type PingServiceReRPC_SumServer struct {
	stream rerpc.Stream
}

func NewPingServiceReRPC_SumServer(stream rerpc.Stream) *PingServiceReRPC_SumServer {
	return &PingServiceReRPC_SumServer{stream}
}

func (s *PingServiceReRPC_SumServer) Receive() (*SumRequest, error) {
	var req SumRequest
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *PingServiceReRPC_SumServer) SendAndClose(msg *SumResponse) error {
	if err := s.stream.CloseReceive(); err != nil {
		return err
	}
	return s.stream.Send(msg)
}

// PingServiceReRPC_CountUpServer is the server-side stream for the
// internal.ping.v1test.PingService.CountUp procedure.
type PingServiceReRPC_CountUpServer struct {
	stream rerpc.Stream
}

func NewPingServiceReRPC_CountUpServer(stream rerpc.Stream) *PingServiceReRPC_CountUpServer {
	return &PingServiceReRPC_CountUpServer{stream}
}

func (s *PingServiceReRPC_CountUpServer) Send(msg *CountUpResponse) error {
	return s.stream.Send(msg)
}

// PingServiceReRPC_CumSumServer is the server-side stream for the
// internal.ping.v1test.PingService.CumSum procedure.
type PingServiceReRPC_CumSumServer struct {
	stream rerpc.Stream
}

func NewPingServiceReRPC_CumSumServer(stream rerpc.Stream) *PingServiceReRPC_CumSumServer {
	return &PingServiceReRPC_CumSumServer{stream}
}

func (s *PingServiceReRPC_CumSumServer) Receive() (*CumSumRequest, error) {
	var req CumSumRequest
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *PingServiceReRPC_CumSumServer) Send(msg *CumSumResponse) error {
	return s.stream.Send(msg)
}
