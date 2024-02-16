package main

import (
	"time"

	"github.com/fatih/color"
)

type Barbershop struct {
	ShopCapacity    int
	HaircutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *Barbershop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false // barber is not sleeping

		color.Yellow(barber + " is waiting for clients")

		// check if is any clients
		for {
			// if no clients, barber goes to sleep
			if len(shop.ClientsChan) == 0 { // if no clients, in the channel
				color.Yellow(barber + " goes to sleep")
				isSleeping = true
			}

			// we have at least one client, il deservim
			// shopOpen is true is channel is open and false if channel is closed
			client, shopOpen := <-shop.ClientsChan

			if shopOpen {
				if isSleeping {
					color.Yellow(barber + " is waking up for " + client)
					isSleeping = false
				}

				// cut hair
				shop.cutHair(barber, client)
			} else {
				// shop is closed, barber goes home, goroutine closes
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *Barbershop) cutHair(barber string, client string) {
	color.Green(client + " is getting a haircut from " + barber)
	time.Sleep(shop.HaircutDuration)
	color.Green(client + " is done")
}

func (shop *Barbershop) sendBarberHome(barber string) {
	color.Green(barber + " is going home")
	shop.BarbersDoneChan <- true
}

func (shop *Barbershop) closeShopForDay() {
	color.Cyan("closing shop for the day")
	close(shop.ClientsChan)
	shop.Open = false

	for a := 0; a < shop.NumberOfBarbers; a++ {
		// blocks untill each barber have sent the done value
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)

	color.Green("The barbsershop is closed for the day, and everyone has gone home.")
	color.Green("--------------------------- 		---------------------------")
}

func (shop *Barbershop) addClient(client string) {
	// print out a message
	color.Green(client + " has arrived")

	if shop.Open {
		select {
		case shop.ClientsChan <- client: // if we have space in client chan(10 clients max), we add a client
			color.Blue(client + " takes a set in the waiting room")
		default: // no space remains in client chan
			color.Red(client + " is leaving, no seats available in waiting room")
		}
	} else {
		color.Red(client + " has arrived, but the shop is closed")
	}
}
