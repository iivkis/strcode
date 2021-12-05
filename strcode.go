package strcode

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	nanosec int64 = 1_000_000_000
)

var (
	ErrExpired            = errors.New("the expiration date has expired")
	ErrIncorrectParametrs = errors.New("incorrect parametrs")
	ErrIncorrectHash      = errors.New("the code has been edited")
	ErrExpiresInIsZero    = errors.New("expires in cannot be zero")
)

type Strcode struct {
	secret    rune
	seperator string
	expiresIn int64
}

func NewStrcode(secret string, separator string, expiresIn time.Duration) (*Strcode, error) {
	if expiresIn == 0 {
		return &Strcode{}, ErrExpiresInIsZero
	}

	var sum rune = 1
	for _, char := range secret {
		sum += char
	}

	return &Strcode{
		secret:    sum,
		seperator: separator,
		expiresIn: int64(expiresIn) / nanosec,
	}, nil
}

func (s *Strcode) hash(str string, expiresAt int64) (hash int64) {
	var sum rune = 1
	for _, char := range str {
		sum += char
	}

	hash = int64(sum * s.secret * rune(expiresAt))
	if hash < 0 {
		hash *= -1
	}

	return hash
}

func (s *Strcode) Encode(str string) string {
	expiresAt := time.Now().Unix() + s.expiresIn
	return str + s.seperator + strconv.FormatInt(s.hash(str, expiresAt), 10) + s.seperator + strconv.FormatInt(expiresAt, 10)
}

func (s *Strcode) Decode(code string) (string, error) {
	codeSplit := strings.Split(code, s.seperator)
	if len(codeSplit) != 3 {
		return "", ErrIncorrectParametrs
	}

	expiresAt, err := strconv.ParseInt(codeSplit[2], 10, 0)
	if err != nil {
		return "", err
	}

	if time.Now().Unix() >= expiresAt {
		return "", ErrExpired
	}

	hash, err := strconv.ParseInt(codeSplit[1], 10, 0)
	if err != nil {
		return "", err
	}

	if hash != s.hash(codeSplit[0], expiresAt) {
		return "", ErrIncorrectHash
	}
	return codeSplit[0], nil
}
