package ratecounter

import (
	"fmt"
	"log"
	"time"

	"testing"
)

func TestMilisecond(t *testing.T) { // count miliseconds in a second
	counter := NewRateCounter(time.Second)
	counter.Start()
	
	label := "miliseconds"

	stop_tick := make(chan bool)
	go func(stop chan bool) {
		ticker := time.NewTicker(time.Millisecond)
		for {
			select {
			case <- ticker.C:
				counter.Inc(label)
			case <-stop:
				return
			}
		}
	}(stop_tick)

	stop_print := make(chan bool)
	go func(stop chan bool) {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <- ticker.C:
				fmt.Printf("%s: %d \n", label, counter.Count(label))
			case <-stop:
				return
			}
		}
	}(stop_print)

	<-time.After(10*time.Second)
	stop_tick <- true
	stop_print <- true
	counter.Stop()
	log.Println("Exiting.")
}

func TestNanosecond(t *testing.T) { // count nanoseconds in a second
	counter := NewRateCounter(time.Second)
	counter.Start()
	
	label := "nanoseconds"

	stop_tick := make(chan bool)
	go func(stop chan bool) {
		ticker := time.NewTicker(time.Nanosecond)
		for {
			select {
			case <- ticker.C:
				counter.Inc(label)
			case <-stop:
				return
			}
		}
	}(stop_tick)

	stop_print := make(chan bool)
	go func(stop chan bool) {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <- ticker.C:
				fmt.Printf("%s: %d \n", label, counter.Count(label))
			case <-stop:
				return
			}
		}
	}(stop_print)

	<-time.After(10*time.Second)
	stop_tick <- true
	stop_print <- true
	counter.Stop()
	log.Println("Exiting.")
}