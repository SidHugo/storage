package utils

import (
	"fmt"
	"log"
	"os"
)

func SetUpLogger(loggerName string) *os.File {
	f, err := os.OpenFile(loggerName+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
	}

	log.SetOutput(f)

	return f
}
