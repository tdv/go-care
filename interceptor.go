// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc/metadata"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"
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
	key := strings.Builder{}

	// Add method name
	key.WriteString(method)
	// Add meta
	this.processMeta(ctx, &key)
	// Add serialized request
	if err := robustHashingData(req, &key); err != nil {
		return "", err
	}

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

func (this *interceptor) cacheResponse(key string, val interface{}, ttl time.Duration) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return this.opts.Cache.Put(key, buf, ttl)
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
	if this.opts.Switch.IsTurnedOn() {
		if cacheable, ttl := this.opts.Methods.Cacheable(method); cacheable {
			key, err := this.makeKey(ctx, method, req)
			if err != nil {
				log.Printf("Failed to make the key. Error: %v\n", err)
			} else {
				typ, hasType := this.types.Get(key)
				if hasType {
					buf, err := this.opts.Cache.Get(key)
					if err == nil && len(buf) > 0 {
						if resp, err := this.restoreResponse(typ, buf); err == nil {
							return resp, nil
						} else {
							log.Printf("Failed to restore a response from the cache. Error: %v\n", err)
						}
					}
				}

				resp, err := handler(ctx, req)

				if !hasType {
					err = this.types.Put(key, resp)
					if err != nil {
						log.Printf("Failed to memorize the Type. Error: %v\n", err)
					}
				}

				err = this.cacheResponse(key, resp, ttl)
				if err != nil {
					log.Printf("Failed to cache the response. Error: %v\n", err)
				}

				return resp, nil
			}
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
