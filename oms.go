package main

import (
	client "oms-go/oms_data_collector"
	"encoding/json"
	"fmt"
	"time"
)

type omsMessage struct {
	SourceSystem   string `json:"sourceSystem,omitempty"`
	ContainerID    string `json:"containerId"`
	ContainerName  string `json:"containerName"`
	TimeGenerated  int64  `json:"timeGenerated"`
	LogEntry       string `json:"logEntry"`
}

// Required parameters
var (
	// Update customerId to your Operations Management Suite workspace ID
	omscustomerID    = "xxxxxxxx-xxx-xxx-xxx-xxxxxxxxxxxx"
	// For sharedKey, use either the primary or the secondary Connected Sources client authentication key
	omssharedKey     = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	// HTTP timeout for posting events to OMS Log Analytics
	omsPostTimeout   = 5 * time.Second
	// Interval to post an OMS batch
	//omsBatchTime         = kingpin.Flag("oms-batch-time", "Interval to post an OMS batch").Default("5s").OverrideDefaultFromEnvar("OMS_BATCH_TIME").Duration()
	// Max number of messages per OMS batch
	//omsMaxMsgNumPerBatch = kingpin.Flag("oms-max-msg-num-per-batch", "").Default("1000").OverrideDefaultFromEnvar("OMS_MAX_MSG_NUM_PER_BATCH").Int()
)

func main() {
	omsclient := client.NewOmsLogClient(omscustomerID, omssharedKey, omsPostTimeout)

	// An example JSON data message to post
	msg := &omsMessage{
			SourceSystem:   "MySystemName",
			ContainerID:    "1234567890",
			ContainerName:  "mycontainer",
			TimeGenerated:  int64(time.Millisecond),
			LogEntry:       "super important log event",
		}

	buffer, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	postErr := omsclient.PostData(&buffer, "Container Log")
	if postErr != nil {
		fmt.Println("error:", postErr)
	}
}