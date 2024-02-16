package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// variables
var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second
var concomitentClients = 3

func main() {
	// seed our random number generator
	rand.Seed(time.Now().UnixNano())

	// wlc msg
	color.Yellow("Welcome to the barbershop")
	color.Yellow("---------------------------")

	// create channels
	// max 10 people in this channel
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := Barbershop{
		ShopCapacity:    seatingCapacity,
		HaircutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("Barbershop is open for business")

	// add barbers
	//  1 goroutine for each barber
	shop.addBarber("Frank")
	shop.addBarber("George")
	shop.addBarber("Susan")
	shop.addBarber("Kelly")

	// start the barbershop as a goroutine

	shopClosing := make(chan bool)
	closed := make(chan bool)
	// ensure goroutine stays open $timeOpen
	go func() {
		// blocks until timeOpen passes
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients
	i := 0

	go func() {
		for {
			// rand avg arrivalRate
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Duration(randomMilliseconds) * time.Millisecond): // when client arrives
				for client := 0; client < concomitentClients; client++ {
					shop.addClient(fmt.Sprintf("Client #%d", i))
					i++
				}
			}
		}
	}()

	// block for the barbershop to close
	// blocks untill we receive smth on close chanell
	// the only thing we can receive here is from line 57 "closed <- true"
	<-closed
}
