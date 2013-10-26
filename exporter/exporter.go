package main

import (
	"github.com/cppforlife/checkman_jenkins_proxy/storer"
	"log"
)

type Exporter interface {
	Run()
}

type PeriodicExporter struct {
	fetchScheduler FetchScheduler
	fetcher        Fetcher
	storer         storer.Storer
	storeKey       string
	logger         *log.Logger
}

func NewPeriodicExporter(
	fetchScheduler FetchScheduler,
	fetcher Fetcher,
	storer storer.Storer,
	storeKey string,
	logger *log.Logger,
) *PeriodicExporter {

	return &PeriodicExporter{
		fetchScheduler,
		fetcher,
		storer,
		storeKey,
		logger,
	}
}

func (pe *PeriodicExporter) Run() {
	eachTick := func() {
		content, err := pe.fetcher.Fetch()
		if err != nil {
			pe.logger.Printf("exporter.fetch.fail key=%s err=%v\n", pe.storeKey, err)
			return
		}

		defer content.Close()

		err = pe.storer.Put(pe.storeKey, content)
		if err != nil {
			pe.logger.Printf("exporter.store.fail key=%s err=%v\n", pe.storeKey, err)
			return
		}

		pe.logger.Printf("exporter.success key=%s\n", pe.storeKey)
	}

	pe.logger.Printf("exporter.run key=%s\n", pe.storeKey)
	pe.fetchScheduler.Run(eachTick)
}
