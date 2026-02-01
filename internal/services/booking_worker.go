package services

import "log"

type BookingWorker struct {
	Queue chan string
}

func NewBookingWorker() *BookingWorker {
	return &BookingWorker{
		Queue: make(chan string, 10),
	}
}

func (w *BookingWorker) Start() {
	go func() {
		for msg := range w.Queue {
			log.Println("Booking processed:", msg)
		}
	}()
}
