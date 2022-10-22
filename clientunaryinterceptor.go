package care

import (
	"context"
	"errors"
	"google.golang.org/grpc"
)

type unaryClientInterceptor struct {
	interceptor *interceptor
}

func (this *unaryClientInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		resp, err := this.interceptor.execute(
			ctx,
			method,
			req,
			func(c context.Context, r interface{}) (interface{}, error) {
				e := invoker(c, method, r, reply, cc, opts...)
				return reply, e
			},
		)

		if resp != nil {
			reply = resp
		}

		return err
	}
}

func NewClientUnaryInterceptor(opts *Options) (grpc.UnaryClientInterceptor, error) {
	if opts == nil {
		return nil, errors.New("The options must not be provided as a nil-pointer.")
	}

	interceptor := unaryClientInterceptor{
		interceptor: newInterceptor(opts),
	}

	return interceptor.Unary(), nil
}
