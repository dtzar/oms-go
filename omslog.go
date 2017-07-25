package omslog

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
	"oms-go/oms_data_collector"
)

const (
	name = "omslog"

	// Options
	workspaceKey = "workspaceId"
	sharedKey = "sharedKey"
	timeout = 30 * time.Second

	// Errors
	errOptRequired = "must specify a value for log opt '%s'"
)

type omsLogger struct {
	containerID string
	containerName string
	imageID string
	imageName string
	client oms_data_collector.OmsLogClient
}

type omsMessage struct {
	ContainerID string `json:"containerId"`
	ContainerName string `json:"containerName"`
	ImageID string `json:"imageId"`
	ImageName string `json:"imageName"`
	Timestamp int64 `json:"timestamp"`
	Source string `json:"source"`
	Line string `json:"line"`
}

func init() {
	if err := logger.RegisterLogDriver(name, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// ValidateLogOpt looks for workspaceKey and sharedKey
func ValidateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
			case workspaceKey:
			case sharedKey:
			default:
				return fmt.Errorf("unknown log opt '%s' for %s log driver", key, name)
		}
	}

	if cfg[workspaceKey] == "" {
		return fmt.Errorf(errOptRequired, workspaceKey)
	}

	if cfg[sharedKey] == "" {
		return fmt.Errorf(errOptRequired, sharedKey)
	}

	return nil
}

// New creates an omslog using configuration options passed in via the context.
func New(info logger.Info) (logger.Logger, error) {
	workspaceID := info.Config[workspaceKey]
	sharedKey := info.Config[sharedKey]

	l := &omsLogger{
		containerID: info.ContainerID,
		containerName: info.ContainerName,
		imageID: info.ContainerImageID,
		imageName: info.ContainerImageName,
		client: oms_data_collector.NewOmsLogClient(workspaceID, sharedKey, timeout),
	}
	
	return l, nil
}

func (l *omsLogger) Log(message *logger.Message) error {
	msg := &omsMessage {
		ContainerID: l.containerID,
		ContainerName: l.containerName,
		ImageID: l.imageID,
		ImageName: l.imageName,
		Timestamp: message.Timestamp.UnixNano() / int64(time.Millisecond),
		Source: message.Source,
		Line: string(message.Line),
	}

	buffer, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := l.client.PostData(&buffer, "Line"); err != nil {
		return err
	}

	return nil
}

func (l *omsLogger) Name() string {
	return name
}

func (l *omsLogger) Close() error {
	return nil
}
