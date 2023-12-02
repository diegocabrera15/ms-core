package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sirupsen/logrus"
	"time"
)

// CloudWatchLogsHook es un hook logrus que envía logs a CloudWatch Logs.
type CloudWatchLogsHook struct {
	logGroupName  string
	logStreamName string
	cwl           *cloudwatchlogs.CloudWatchLogs
}

// NewCloudWatchLogsHook crea un nuevo hook para enviar logs a CloudWatch Logs.
func NewCloudWatchLogsHook(cwl *cloudwatchlogs.CloudWatchLogs, logGroupName, logStreamName string) (*CloudWatchLogsHook, error) {
	return &CloudWatchLogsHook{
		logGroupName:  logGroupName,
		logStreamName: logStreamName,
		cwl:           cwl,
	}, nil
}

// Levels devuelve los niveles de log para los cuales este hook debe estar activo.
func (hook *CloudWatchLogsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire se llama cuando se emite un log y envía el log a CloudWatch Logs.
func (hook *CloudWatchLogsHook) Fire(entry *logrus.Entry) error {
	message, err := entry.String()
	if err != nil {
		return err
	}

	params := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(hook.logGroupName),
		LogStreamName: aws.String(hook.logStreamName),
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(aws.TimeUnixMilli(time.Now())),
			},
		},
	}

	if _, err := hook.cwl.PutLogEvents(params); err != nil {
		return fmt.Errorf("error al enviar log a CloudWatch Logs: %v", err)
	}

	return nil
}

// SetupCloudWatchLogs configura el hook para enviar logs a CloudWatch Logs.
func SetupCloudWatchLogs(logger *logrus.Logger) {
	logGroupName := "/aws/lambda/my-microservicio-dev-hello" // Reemplaza con el nombre correcto
	logStreamName := "my-log-stream"                         // Puedes personalizar según tus necesidades

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Reemplaza con tu región AWS
	})
	if err != nil {
		logger.Fatalf("Error creando la sesión de AWS: %v", err)
	}

	cwl := cloudwatchlogs.New(sess)

	hook, err := NewCloudWatchLogsHook(cwl, logGroupName, logStreamName)
	if err != nil {
		logger.Fatalf("Error configurando el hook de CloudWatch Logs: %v", err)
	}

	logger.AddHook(hook)
}
