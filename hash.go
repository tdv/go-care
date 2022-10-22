// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash/crc32"
)

type Hash interface {
	Calc(key string) (string, error)
}

const (
	hashPiecesDelimiter = "_"
)

type defaultHash struct {
	crc32Table *crc32.Table
}

func (s *defaultHash) Calc(key string) (string, error) {
	keyLen := len(key)

	if keyLen < 1 {
		return "", errors.New("You can't calc the hash from an empty string.")
	}

	data := []byte(key)
	crc32Hash := crc32.Checksum(data, s.crc32Table)
	sha256Hash := sha256.Sum256(data)
	hash := fmt.Sprintf("%08x%s%x",
		crc32Hash,
		hashPiecesDelimiter,
		sha256Hash,
	)
	return hash, nil
}

func newDefaultHash() Hash {
	return &defaultHash{
		crc32Table: crc32.MakeTable(crc32.Castagnoli),
	}
}
