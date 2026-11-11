package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"strings"

	"example.com/open-monastery/internal/app/abstractions"
	"golang.org/x/crypto/argon2"
)

const (
	defaultMemoryKiB   uint32 = 19 * 1024
	defaultIterations  uint32 = 2
	defaultParallelism uint8  = 1
	defaultSaltLength  uint32 = 16
	defaultHashLength  uint32 = 32
	defaultPrefix      string = "$argon2id$v=19$m=19456,t=2,p=1$"
)

var _ abstractions.PasswordHasher = (*ArgonPasswordHasher)(nil)

type ArgonPasswordHasher struct {
	MemoryKiB   uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	HashLength  uint32
	Prefix      string
}

func NewArgonPasswordHasher() *ArgonPasswordHasher {
	return &ArgonPasswordHasher{
		MemoryKiB:   defaultMemoryKiB,
		Iterations:  defaultIterations,
		Parallelism: defaultParallelism,
		SaltLength:  defaultSaltLength,
		HashLength:  defaultHashLength,
		Prefix:      defaultPrefix,
	}
}

func (a *ArgonPasswordHasher) derive(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		a.Iterations,
		a.MemoryKiB,
		a.Parallelism,
		a.HashLength,
	)
}

func (a *ArgonPasswordHasher) HashPassword(password string) string {
	salt := make([]byte, a.SaltLength)
	rand.Read(salt)

	hash := a.derive(password, salt)
	encoding := base64.RawStdEncoding
	return a.Prefix + encoding.EncodeToString(salt) + "$" + encoding.EncodeToString(hash)
}

func (a *ArgonPasswordHasher) VerifyPassword(password string, stored string) bool {
	payload, ok := strings.CutPrefix(stored, a.Prefix)
	if !ok {
		return false
	}

	saltText, hashText, ok := strings.Cut(payload, "$")
	if !ok {
		return false
	}

	encoding := base64.RawStdEncoding
	salt, err := encoding.DecodeString(saltText)
	if err != nil || len(salt) != int(a.SaltLength) {
		return false
	}
	expected, err := encoding.DecodeString(hashText)
	if err != nil || len(expected) != int(a.HashLength) {
		return false
	}

	actual := a.derive(password, salt)
	return subtle.ConstantTimeCompare(actual, expected) == 1
}
