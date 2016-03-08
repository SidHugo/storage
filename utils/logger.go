package utils

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
)

func SetUpLogger(loggerName string) *logging.Logger {
	var logger = logging.MustGetLogger(loggerName)

	// TODO: change log files location
	f, err := os.OpenFile(loggerName+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file: %v", err)
	}

	backend := logging.NewLogBackend(f, "", 0)

	var format = logging.MustStringFormatter(
		`%{color}%{time} %{shortfunc} ▶ %{level} %{color:reset} %{message}`,
	)

	backendFormatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(backendFormatter)

	return logger
}
