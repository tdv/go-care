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
	timeout time.Duration
}

func (s *redisCache) Put(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.client.Set(ctx, key, val, ttl).Err()
}

func (s *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.client.Get(ctx, key).Bytes()
}

func (s *redisCache) init(
	host string, port int, db int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	if client == nil {
		return errors.New("Failed to create Redis client.")
	}
	s.client = client
	return nil
}

func New(host string, port int, db int) (care.Cache, error) {
	cache := redisCache{}
	if err := cache.init(host, port, db); err != nil {
		return nil, err
	}
	return &cache, nil
}
