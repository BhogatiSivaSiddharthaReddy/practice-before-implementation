package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("---Day 3 practice -----")
	calling_function()
	multiple_goroutines()
	channel_example()
	signal_dmy()
	owner()
	example_channel()
	context_example()
}

func blocking_function() {
	fmt.Println("Blocking function")
	a := 2 * time.Second
	time.Sleep(a)
	fmt.Println("waited for 2 seconds")
}

func calling_function() {
	fmt.Println("calling function")
	go blocking_function()
	time.Sleep(3 * time.Second)
	fmt.Println("closing the calling function")
}

func called(s string) {
	fmt.Println(s)
	for i := 1; i <= 3; i++ {
		fmt.Println(i, s)
		time.Sleep(500 * time.Microsecond)
	}
}

func multiple_goroutines() {
	go called("A")
	go called("B")
	time.Sleep(3 * time.Second)
}

func channel_example() {
	ch := make(chan int)

	go func() {
		fmt.Println("Inputting value 5 into channel ch")
		ch <- 5
	}()

	fmt.Println("Outputting value 5 from channel to variable i")
	i := <-ch

	fmt.Println(i)

}

func signal_dmy() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT)

	a := syscall.SIGINT
	fmt.Println(a, int(a))

	fmt.Println("Waiting for SIGINT")

	si := <-ch

	sig := si.(syscall.Signal)

	fmt.Println(sig)
	fmt.Println(int(sig))

	fmt.Println("Received:", si)

}

func worker() {
	for {
		fmt.Println("working...")
		time.Sleep(1 * time.Second)
	}
}

func owner() {
	go worker()
	time.Sleep(4 * time.Second)
}

func recieving_channel(ch <-chan struct{}) {

	fmt.Println("waiting for data from channel")

	<-ch

	fmt.Println("I got executed second")

}

func sending_channel(ch chan<- struct{}) {

	fmt.Println("Sending data to channel")

	close(ch)

	fmt.Println("I got executed first")

}

func example_channel() {

	ch := make(chan struct{})

	go sending_channel(ch)
	go recieving_channel(ch)

	time.Sleep(1 * time.Second)
}

func context_example() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	<-ctx.Done()

	fmt.Println("Timeout reached!")
}
