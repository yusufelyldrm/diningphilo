package main

import (
	"fmt"
	"sync"
)

const numPhilosophers = 5
const numMealsPerPhilosopher = 3

type Chopstick struct{ sync.Mutex }

type Philosopher struct {
	id                            int
	leftChopstick, rightChopstick *Chopstick
	mealsEaten                    int
	host                          *Host
}

type Host struct {
	permission chan bool
}

func (p Philosopher) eat() {
	for p.mealsEaten < numMealsPerPhilosopher {
		// Ask permission from the host
		<-p.host.permission

		// Try to pick up the chopsticks
		p.leftChopstick.Lock()
		p.rightChopstick.Lock()

		// Start eating
		fmt.Printf("starting to eat %d\n", p.id)
		p.mealsEaten++

		// Finish eating
		fmt.Printf("finishing eating %d\n", p.id)

		// Release the chopsticks
		p.rightChopstick.Unlock()
		p.leftChopstick.Unlock()

		// Release permission to the host
		p.host.permission <- true
	}
}

func main() {
	// Create the chopsticks and the host
	chopsticks := make([]*Chopstick, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		chopsticks[i] = new(Chopstick)
	}
	host := Host{permission: make(chan bool, 2)}

	// Create the philosophers
	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:             i + 1,
			leftChopstick:  chopsticks[i],
			rightChopstick: chopsticks[(i+1)%numPhilosophers],
			host:           &host,
		}
	}

	// Start the philosophers
	var wg sync.WaitGroup
	wg.Add(numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		go func(p *Philosopher) {
			defer wg.Done()
			p.eat()
		}(philosophers[i])
	}

	// Allow the philosophers to start eating
	for i := 0; i < 2; i++ {
		host.permission <- true
	}

	// Wait for the philosophers to finish
	wg.Wait()
}

