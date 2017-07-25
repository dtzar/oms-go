
package oms_data_collector

import (
	"testing"
	"time"
)


//date string, contentLength int, method string, contentType string, resource string
func TestNewOmeClient(t *testing.T) {
	client := NewOmsLogClient("customerID string", "sharedKey string", time.Second * 30)

	if client == nil {
		t.Fatal("Did not Create a new Client")
	}
}
