package test

import (
	"github.com/beckjiang/nunchakus"
	"log"
	"time"
)

type Producer struct {

}

func (p *Producer) Run(progress nunchakus.Progress) chan string {
	data := make(chan string)
	go func() {
		t := time.NewTicker(time.Second)

		i := progress.Load()
		log.Println("last i:", i)
		for {
			select {
			case <-t.C:
				// do something
				data <- "hello"
				i ++
				log.Println(i)
				progress.Save(i)
			}
		}
	}()
	return data
}