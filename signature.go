
package main

import (
	"fmt"
	"crypto/sha256"
	"encoding/base64"
)

func Signature(verb,contentLength,contentType,date,resource string) string {
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",verb,contentLength,contentType,date,resource);
	sum256 := sha256.Sum256([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(sum256[:])
}

