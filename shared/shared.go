package shared

import (
	"crypto/sha256"
	"fmt"
)

func PanicOnErr(err error) {
	if err == nil {
		return
	}

	panic(err)
}

func Unwrap[T any](result T, err error) T {
	PanicOnErr(err)
	return result
}

func Sha256(text string) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(text)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
