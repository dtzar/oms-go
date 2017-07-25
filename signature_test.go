
package main

import (
	"fmt"
	"testing"
	"crypto/sha256"
	"encoding/base64"
)

const (
	fullString = "POST\n1024\nContent-Type\nx-ms-date:Mon, 04 Apr 2016 08:00:00 GMT\n/api/logs"
	postString = "POST"
	lenString = "1024"
	typeString = "Content-Type"
	dateString = "x-ms-date:Mon, 04 Apr 2016 08:00:00 GMT"
	resourceString = "/api/logs"
)

func TestValidateLogOpt(t *testing.T) {
	sum256 := sha256.Sum256([]byte(fullString));
	fullbase64 := base64.StdEncoding.EncodeToString( sum256[:]);

	testbase64 := Signature(postString,lenString,typeString,dateString,resourceString);
	fmt.Print(fullbase64,"---",testbase64)

	if testbase64 != fullbase64 {
		t.Fatal("signature did not calculate the same");
	}
}
