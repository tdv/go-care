// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package rediscache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/tdv/go-care"
	"time"
)

type redisCache struct {
	client  *redis.Client
	ctx     context.Context
	timeout time.Duration
}

func (s *redisCache) Put(key string, val []byte, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	return s.client.Set(ctx, key, val, ttl).Err()
}

func (s *redisCache) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	return s.client.Get(ctx, key).Bytes()
}

func (s *redisCache) init(
	ctx context.Context, timeout time.Duration,
	host string, port int, db int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	err := func() error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return client.Ping(ctx).Err()
	}()

	if err != nil {
		return err
	}

	if client == nil {
		return errors.New("Failed to create Redis client.")
	}
	s.client = client
	s.ctx = ctx
	s.timeout = timeout
	return nil
}

func New(ctx context.Context, opTimeout time.Duration,
	host string, port int, db int) (care.Cache, error) {
	cache := redisCache{}
	if err := cache.init(ctx, opTimeout, host, port, db); err != nil {
		return nil, err
	}
	return &cache, nil
}
