package main

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	log *log.Logger
	ch  chan string
}

func NewLogger() *Logger {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)
	return &Logger{log: logger, ch: make(chan string)}
}

func (l *Logger) Log(message string) {
	l.ch <- message
}

func (l *Logger) Run() {
	for {
		select {
		case message := <-l.ch:
			l.log.Printf(message)
		case <-time.After(5 * time.Minute):
			l.log.Printf("Запись логов в файл...")
		}
	}
}
