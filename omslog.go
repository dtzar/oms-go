package omslog

import (
	"fmt"
	"log"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
	"github.com/Azure/oms-log-analytics-firehose-nozzle/client"
)

const (
	name = "omslog"
)

type omsLogger struct {
	//client *client.Client
}

func init() {
	if err := logger.RegisterLogDriver(name, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

func ValidateLogOpt(cfg map[string]string) error {
	return nil
}

func New(info logger.Info) (logger.Logger, error) {
	l := &omsLogger{} //client := client.NewOmsClient()
	return l, nil
}

func (l *omsLogger) Log(message *logger.Message) error {
	//l.client.PostData()
	return nil
}

func (l *omsLogger) Name() string {
	return name
}

func (l *omsLogger) Close() error {
	return nil
}
