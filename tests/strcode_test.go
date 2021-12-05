package strcode

import (
	"errors"
	"testing"
	"time"

	"github.com/iivkis/strcode"
)

const exampleSecret = "secretKey123"
const exampleStr = "strcode123@gmail.com"

func TestStrcodeInit(t *testing.T) {
	t.Run("simple init", func(t *testing.T) {
		_, err := strcode.NewStrCode(exampleSecret, ":", time.Hour)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("err: expires in is zero", func(t *testing.T) {
		_, err := strcode.NewStrCode(exampleSecret, ":", 0)
		if !errors.Is(err, strcode.ErrExpiresInIsZero) {
			t.Error(err)
		}
	})
}

func TestStrcodeEncode(t *testing.T) {
	t.Run("simple encode", func(t *testing.T) {
		sc, _ := strcode.NewStrCode(exampleSecret, ":", time.Second)
		encode := sc.Encode(exampleStr)
		t.Log(encode)
	})

	t.Run("empty string encode", func(t *testing.T) {
		sc, _ := strcode.NewStrCode(exampleSecret, ":", time.Second)
		encode := sc.Encode("")
		t.Log(encode)
	})
}

func TestStrcodeDecode(t *testing.T) {
	sc, _ := strcode.NewStrCode(exampleSecret, ":", time.Second)
	encode := sc.Encode(exampleStr)

	t.Run("simple decode", func(t *testing.T) {
		decode, err := sc.Decode(encode)
		t.Log(decode)

		if err != nil {
			t.Error(err)
		}
	})

	t.Run("edit code", func(t *testing.T) {
		edit := string([]byte(encode)[1:])
		decode, err := sc.Decode(edit)
		t.Log(decode)

		if !errors.Is(err, strcode.ErrIncorrectHash) {
			t.Error(err)
		}
	})

	t.Run("expired", func(t *testing.T) {
		time.Sleep(time.Second)
		decode, err := sc.Decode(encode)
		t.Log(decode)

		if !errors.Is(err, strcode.ErrExpired) {
			t.Error(err)
		}
	})
}

func BenchmarkStrcodeEncode(b *testing.B) {
	sc, _ := strcode.NewStrCode(exampleSecret, ":", time.Hour)
	for i := 0; i < b.N; i++ {
		sc.Encode(exampleStr)
	}
}

func BenchmarkStrcodeDecode(b *testing.B) {
	sc, _ := strcode.NewStrCode(exampleSecret, ":", time.Hour)
	encode := sc.Encode(exampleStr)
	for i := 0; i < b.N; i++ {
		sc.Decode(encode)
	}
}
