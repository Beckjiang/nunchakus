package main

import (
	"github.com/beckjiang/nunchakus"
	"github.com/beckjiang/nunchakus/test"
)

func main() {
	nc := nunchakus.NewNunchakus(nunchakus.WithProducer(&test.Producer{}), nunchakus.WithConsumer(&test.Consumer{}), nunchakus.WithTaskName("happy"))

	go nc.StartConsumer()
	nc.StartProducer()
}
