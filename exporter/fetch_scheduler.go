package main

import (
	"log"
	"time"
)

type FetchScheduler interface {
	Run(f func())
}

type PeriodicFetchScheduler struct {
	ticker *time.Ticker
	logger *log.Logger
}

func NewPeriodicFetchScheduler(
	interval time.Duration,
	logger *log.Logger,
) *PeriodicFetchScheduler {

	return &PeriodicFetchScheduler{
		time.NewTicker(interval),
		logger,
	}
}

func (fs *PeriodicFetchScheduler) Run(f func()) {
	fs.logger.Println("periodic-fetch-scheduler.run")

	for {
		select {
		case <-fs.ticker.C:
			fs.logger.Println("periodic-fetch-scheduler.tick")
			f()
		}
	}
}
