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
	maxKeySuze          = 1024 * 1024
	hashPiecesDelimiter = "_"
)

type defaultHash struct {
	crc32Table *crc32.Table
}

func (this *defaultHash) Calc(key string) (string, error) {
	keyLen := len(key)

	if keyLen < 1 {
		return "", errors.New("You can't calc the hash from an empty string.")
	}

	if keyLen > maxKeySuze {
		return "", errors.New(fmt.Sprintf(
			"The key contains %d characters, but the one must be less than %d characters.",
			keyLen,
			maxKeySuze,
		))
	}

	data := []byte(key)
	crc32Hash := crc32.Checksum(data, this.crc32Table)
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
