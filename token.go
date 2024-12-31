// Package token generates randomized base62 encoded tokens based on a single integer.
// It's ideal for shorturl services or for semi-secured randomized api primary keys.
//
// # How it Works
//
// `Token` is an alias for `uint64`.
// Its `Token.Encode()` method interface returns a `Base62` encoded string based off
// of the number. It's implementation of the `encoding.TextMarshaler` and
// `encoding.TextUnmarshaler` interfaces encodes and decodes the `Token` when it's
// being marshalled or unmarshalled as json or xml.
//
// Basically, the outside world will always address the token as its string equivalent
// and internally we can always be used as an `uint64` for fast, indexed, unique,
// lookups in various databases.
//
// NOTE: Remember to always check for collisions when adding randomized tokens to a
// database.
package token

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
)

// Token is an alias of an uint64 that is marshalled into a base62 encoded token.
type Token uint64

const (
	// Base62 is a string respresentation of every possible base62 character.
	Base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// MaxTokenLength is the largest possible character length of a token.
	MaxTokenLength = 10
	// MinTokenLength is the smallest possible character length of a token.
	MinTokenLength = 1
	// DefaultTokenLength is the default size of a token.
	DefaultTokenLength = 9

	base62Len = Token(len(Base62))
)

// Possible errors.
var (
	ErrTokenTooSmall    = errors.New("the base62 token is smaller than MinTokenLength")
	ErrTokenTooBig      = errors.New("the base62 token is larger than MaxTokenLength")
	ErrInvalidCharacter = errors.New("there was a non base62 character in the token")
)

// New returns a [Base62] encoded [Token] of *up to* [DefaultTokenLength].
// If you pass in a `tokenLength` between [MinTokenLength] and [MaxTokenLength]
// this will return a [Token] of *up to* that length. If you pass in an out of
// range [tokenLength] it will be adjusted. It never emits 0 as a Token.
func New(opts ...int) (t Token) {
	maX := maxHashInt(DefaultTokenLength)

	if len(opts) > 0 {
		x := opts[0]
		maX = maxHashInt(min(max(x, MinTokenLength), MaxTokenLength))
	}

	t = Token(rand.Int63n(int64(maX & math.MaxInt64))) //nolint:gosec // not applicable
	if t == 0 {
		t++
	}

	return
}

// Decode returns a token from a 1-12 character base62 encoded string.
func Decode(token string) (t Token, err error) {
	err = (&t).UnmarshalText([]byte(token))
	return
}

// Encode encodes the token into a base62 string.
func (t Token) Encode() string {
	bs, _ := t.MarshalText() //nolint:errcheck // it never errors out
	return string(bs)
}

func (t *Token) UnmarshalText(data []byte) (err error) {
	number := Token(0)
	idx := 0.0
	chars := []byte(Base62)

	charsLength := float64(len(chars))
	tokenLength := float64(len(data))

	if tokenLength > MaxTokenLength {
		return ErrTokenTooBig
	} else if tokenLength < MinTokenLength {
		return ErrTokenTooSmall
	}

	for _, c := range data {
		power := tokenLength - (idx + 1)

		index := bytes.IndexByte(chars, c)
		if index < 0 {
			return ErrInvalidCharacter
		}

		number += Token(index) * Token(math.Pow(charsLength, power))
		idx++
	}

	*t = number

	return
}

func (t Token) MarshalText() (chars []byte, err error) {
	if t == 0 {
		return
	}

	for t > 0 {
		result := t / base62Len
		remainder := t % base62Len
		chars = append(chars, Base62[remainder])
		t = result
	}

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return
}

// maxHashInt returns the largest possible int that will yield a base62
// encoded token of the specified length.
func maxHashInt(length int) uint64 {
	return uint64(max(0, min(math.MaxUint64, math.Pow(float64(base62Len), float64(length)))))
}
