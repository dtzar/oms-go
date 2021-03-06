// OMS Data Collector package adapted from 
// https://github.com/Azure/oms-log-analytics-firehose-nozzle
package oms_data_collector

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OmsLogClient interface {
	PostData(*[]byte, string) error
}

// OmsLogClient posts messages to OMS
type omslogclient struct {
	customerID      string
	sharedKey       string
	url             string
	httpPostTimeout time.Duration
	client			*http.Client
}

const (
	method      = "POST"
	contentType = "application/json"
	resource    = "/api/logs"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 30
}

// New instance of the OmsLogClient
func NewOmsLogClient(customerID string, sharedKey string, postTimeout time.Duration ) OmsLogClient {
	return &omslogclient{
		customerID:      customerID,
		sharedKey:       sharedKey,
		url:             "https://" + customerID + ".ods.opinsights.azure.com" + resource + "?api-version=2016-04-01",
		httpPostTimeout: postTimeout,
		client: 		 &http.Client{ Timeout: postTimeout	},
	}
}

// PostData posts message to OMS
func (c *omslogclient) PostData(msg *[]byte, logType string) error {
	// Headers
	contentLength := len(*msg)
	rfc1123date := time.Now().UTC().Format(time.RFC1123)
	//TODO: rfc1123 date should have UTC offset
	rfc1123date = strings.Replace(rfc1123date, "UTC", "GMT", 1)
	//Signature
	signature, err := c.buildSignature(rfc1123date, contentLength, method, contentType, resource)
	if err != nil {
		return err
	}
	// Create request
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(*msg))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", signature)
	req.Header.Set("Log-Type", logType)
	//TODO: headers should be case insentitive
	//req.Header.Set("x-ms-date", rfc1123date)
	req.Header["x-ms-date"] = []string{rfc1123date}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return fmt.Errorf("Post Error. HTTP response code:%d message:%s", resp.StatusCode, resp.Status)
	}
	return nil
}

func (c *omslogclient) buildSignature(date string, contentLength int, method string, contentType string, resource string) (string, error) {
	xHeaders := "x-ms-date:" + date
	stringToHash := method + "\n" + strconv.Itoa(contentLength) + "\n" + contentType + "\n" + xHeaders + "\n" + resource
	bytesToHash := []byte(stringToHash)
	keyBytes, err := base64.StdEncoding.DecodeString(c.sharedKey)
	if err != nil {
		return "", err
	}
	hasher := hmac.New(sha256.New, keyBytes)
	hasher.Write(bytesToHash)
	encodedHash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	authorization := fmt.Sprintf("SharedKey %s:%s", c.customerID, encodedHash)
	return authorization, err
}