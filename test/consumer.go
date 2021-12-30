package test

import "log"

type Consumer struct {
}

func (c *Consumer) GetConcurrency() int {
	return 1
}

func (c *Consumer) Run(q chan string) {
	for item := range q {
		log.Println("Consumer:", item)
	}
}
