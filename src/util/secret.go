package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

func NewHmac(s string) string {
	h := hmac.New(sha256.New, []byte(s))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
