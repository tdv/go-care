package care

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"reflect"
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

		realComputation := false

		resp, err := this.interceptor.execute(
			ctx,
			method,
			req,
			func(c context.Context, r interface{}) (interface{}, error) {
				e := invoker(c, method, r, reply, cc, opts...)
				realComputation = true
				return reply, e
			},
		)

		if !realComputation && resp != nil {
			replyVal := reflect.ValueOf(reply).Elem()
			respVal := reflect.ValueOf(resp).Elem()
			tmp := replyVal.Interface()

			replyVal.Set(respVal)
			respVal.Set(reflect.ValueOf(tmp))
		}

		return err
	}
}

func NewClientUnaryInterceptor(opts *Options) (grpc.DialOption, error) {
	if opts == nil {
		return nil, errors.New("The options must not be provided as a nil-pointer.")
	}

	interceptor := unaryClientInterceptor{
		interceptor: newInterceptor(opts),
	}

	return grpc.WithUnaryInterceptor(interceptor.Unary()), nil
}
