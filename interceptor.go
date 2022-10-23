package care

import (
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc/metadata"
	"reflect"
	"sort"
	"strings"
)

type interceptor struct {
	opts  *Options
	types typeStorage
}

type interceptorFunc func(context.Context, interface{}) (interface{}, error)

func (this *interceptor) processMeta(ctx context.Context, builder *strings.Builder) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	data := make([]string, 0)

	for k, v := range md {
		if !this.opts.MetaFilter.Allowed(k, v) {
			continue
		}

		sort.Strings(v)
		data = append(data, k+strings.Join(v, ";"))
	}

	sort.Strings(data)
	builder.WriteString(strings.Join(data, ";"))

}

func (this *interceptor) makeKey(
	ctx context.Context,
	method string,
	req interface{},
) (string, error) {
	typ := reflect.TypeOf(req)
	if typ == nil {
		return "", errors.New("An empty request.")
	}

	buf, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	key := strings.Builder{}

	// Add method name
	key.WriteString(method)
	// A bit salt has been added in order to get more unique key.
	key.WriteString(typ.String())
	// Add meta
	this.processMeta(ctx, &key)
	// Add serialized request
	key.Write(buf)

	hash, err := this.opts.Hash.Calc(key.String())
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (this *interceptor) restoreResponse(typ reflect.Type, buf []byte) (resp interface{}, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("Failed to construct new value.")
		}
	}()

	val := reflect.New(typ).Interface()
	err = json.Unmarshal(buf, &val)
	resp = val

	return val, err
}

func (this *interceptor) cacheResponse(key string, val interface{}) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return this.opts.Cache.Put(key, buf)
}

func (this *interceptor) execute(
	ctx context.Context,
	method string,
	req interface{},
	handler interceptorFunc,
) (
	interface{},
	error,
) {
	if this.opts.Switch.IsTurnedOn() && this.opts.Methods.Exists(method) {
		key, err := this.makeKey(ctx, method, req)
		if err == nil {
			typ, hasType := this.types.Get(key)
			if hasType {
				buf, err := this.opts.Cache.Get(key)
				if err == nil && len(buf) > 0 {
					if resp, err := this.restoreResponse(typ, buf); err == nil {
						return resp, nil
					}
				}
			}

			resp, err := handler(ctx, req)

			if !hasType {
				this.types.Put(key, resp)
			}

			go func() {
				this.cacheResponse(key, resp)
			}()

			return resp, err
		}
	}

	return handler(ctx, req)
}

func newInterceptor(opts *Options) *interceptor {
	return &interceptor{
		opts:  opts,
		types: newTypeStorage(),
	}
}
