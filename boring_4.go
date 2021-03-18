package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str  string
	wait chan bool
}

func boring(msg string) <-chan Message {
	c := make(chan Message)
	waitForIt := make(chan bool)

	go func() {
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s: %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c
}

func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
	// go func() {
	//     for {
	//         c <- <-input1
	//     }
	// }()
	// go func() {
	//     for {
	//         c <- <-input2
	//     }
	// }()
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}

	}()
	return c
}

func main() {
	c := fanIn(boring("joe"), boring("ann"))

	for i := 0; i < 5; i++ {
		//fmt.Printf("You say %q\n", <-c)
		//fmt.Println(<-c)
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)

		msg1.wait <- true
		msg2.wait <- true

		fmt.Println("--------------")
	}

	fmt.Println("Your boring, im leaving")
}
