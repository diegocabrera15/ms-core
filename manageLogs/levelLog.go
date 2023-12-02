package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func handler() (string, error) {
	logrus.Info("This is an info message.")
	logrus.Warn("This is a warning message.")
	logrus.Error("This is an error message.")

	err := someFunction()
	if err != nil {
		logrus.WithError(err).Error("An error occurred in someFunction.")
		return "", err
	}

	return "Hello from Lambda in Go with log handling!", nil
}

func someFunction() error {
	// Simula un error
	return fmt.Errorf("Simulated error in someFunction")
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{})

	cloudwatchlogs.SetupCloudWatchLogs(logger)

	logger.SetOutput(os.Stdout)

	return logger
}

func main() {
	logger := setupLogger()

	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		// Ejecución local
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			message, err := handler()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, message)
		})

		port := 3000
		fmt.Printf("Local server listening on :%d...\n", port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		// Ejecución en AWS Lambda
		lambda.Start(handler)
	}
}
