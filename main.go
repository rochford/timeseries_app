package main

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	"github.com/rochford/timeseries"
)

const (
	tagKeyStockName = iota
	tagKeyStockPrice
)
const (
	tagValueStockAlpha = iota
	tagValueStockBravo
)

func observationProducer(stockName timeseries.TagValue, ch chan timeseries.Observation) {
	rand.Seed(time.Now().UnixNano())
	for {
		tagValues := make([]timeseries.Tag, 0, 5)
		tag := timeseries.Tag{Key: tagKeyStockName, Value: stockName}
		tagValues = append(tagValues, tag)

		value := rand.Intn(100)
		observation := timeseries.TagValue(value)
		tagValues = append(tagValues, timeseries.Tag{Key: tagKeyStockPrice, Value: observation})
		obs, _ := timeseries.NewObservation(tagValues)

		ch <- obs
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
		if value < 2 {
			close(ch)
			break
		}
	}
}

func main() {
	stockAlphaObservationChannel := make(chan (timeseries.Observation))
	stockBravoObservationChannel := make(chan (timeseries.Observation))

	ts := timeseries.NewTimeSeries("stock.quote", time.Minute)

	go observationProducer(tagValueStockAlpha, stockAlphaObservationChannel)
	go observationProducer(tagValueStockBravo, stockBravoObservationChannel)

	exit := false
	for exit == false {
		select {
		case obs, ok := <-stockAlphaObservationChannel:
			if !ok {
				exit = true
				break
			}
			processObservationEvent(ts, obs)
		case obs, ok := <-stockBravoObservationChannel:
			if !ok {
				exit = true
				break
			}
			processObservationEvent(ts, obs)
		}
	}

	var b bytes.Buffer
	ts.Flush(&b)
	log.Println(b)

}

func processObservationEvent(ts *timeseries.TimeSeries, obs timeseries.Observation) {
	timestamp := time.Now()

	ts.AddPoint(obs, timestamp)
	log.Println("observation:", obs.String())

}
