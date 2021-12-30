package nunchakus

import (
	"strings"
	"time"
)

// Producer -> queue -> Consumer
type Consumer interface {
	GetConcurrency () int
	Run	(q chan string)
}
type Producer interface {
	Run(progress Progress) chan string
}

// 保存进度，续跑使用
type Progress interface {
	Save(i int)
	Load() int
}

type Queue interface {
	Push(item string) bool
	Pop() string
	PopN(n int) []string
}

type Nunchakus struct {
	taskName string
	consumer Consumer
	producer Producer
	
	progress Progress
	queue Queue
}

type Setter func (nc *Nunchakus)

func WithTaskName(name string) Setter {
	return func(nc *Nunchakus) {
		name = strings.Replace(name, " ", "_", -1)
		nc.taskName = name
	}
}

func WithQueue(q Queue) Setter {
	return func(nc *Nunchakus) {
		nc.queue = q
	}
}
func WithProducer(p Producer) Setter {
	return func(nc *Nunchakus) {
		nc.producer = p
	}
}

func WithConsumer(c Consumer) Setter {
	return func(nc *Nunchakus) {
		nc.consumer = c
	}
}

func NewNunchakus(opts ... Setter ) *Nunchakus {
	nc := &Nunchakus{}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(nc)
		}
	}
	if nc.taskName == "" {
		nc.taskName = "Nunchakus"
	}

	if nc.progress == nil {
		nc.progress = NewFileProcess("")
	}

	if nc.queue == nil {
		nc.queue = NewRedisQueue(nil, nc.taskName)
	}

	return nc
}

func (nc *Nunchakus) StartProducer() {
	if nc.producer == nil {
		panic("producer is required")
	}

	data := nc.producer.Run(nc.progress)
	for item := range data {
		nc.queue.Push(item)
	}
}

func (nc *Nunchakus) StartConsumer() {
	if nc.consumer == nil {
		panic("consumer is required")
	}

	n := nc.consumer.GetConcurrency()
	dataChan := make(chan string, n)
	println("Concurrency:", n)
	for i := 0; i < n; i++ {
		go nc.consumer.Run(dataChan)
	}

	for {
		message := nc.queue.Pop()
		if message == "" {
			time.Sleep(time.Second)
			continue
		}

		dataChan <- message
	}
}