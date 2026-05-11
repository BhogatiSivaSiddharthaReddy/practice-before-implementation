package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println(".....Day 4 excercises.....")
	ex1()
	ex2()
	ex3()
	ex4()
	ex5()
	ex6and7()
	ex8()
}

func task(i int, name string) {
	fmt.Println(i, name)
}

func ex1() {
	for i := 0; i < 5; i++ {
		go task(i, "Sid")
	}

	time.Sleep(1 * time.Second)
}

func ex2() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Work Completed"
	}()

	m := <-ch

	fmt.Println(m)
}

func ex3() {
	ch := make(chan string, 2)

	func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "Sid"
	}()

	m := <-ch
	fmt.Println(m)
}

func workers(i int, j int, ch chan string) {
	fmt.Printf("worker %d started \n", i)
	time.Sleep(time.Duration(j) * time.Second)
	ch <- fmt.Sprintf("Worker %d strated after %d delay", i, j)
}

func ex4() {
	ch := make(chan string)
	go workers(1, 2, ch)
	go workers(2, 1, ch)
	go workers(3, 1, ch)

	for i := 0; i < 3; i++ {
		r := <-ch
		fmt.Println(r)
	}

}

func ex5() {
	ch := make(chan string)

	go func() {
		time.Sleep(5 * time.Second)
		ch <- "A"
	}()

	select {
	case r := <-ch:
		fmt.Println(r)

	case <-time.After(1 * time.Second):
		fmt.Println("Time out")
	}
}

var (
	counter int
	mu      sync.Mutex
)

func increment() {
	mu.Lock()
	counter++
	mu.Unlock()
}

func ex6and7() {
	for i := 0; i < 100; i++ {
		go increment()
	}

	time.Sleep(3 * time.Second) // still you will get data race condition, this will be solved by waigroup functionality.

	fmt.Println(counter)
}

var (
	value int
	rwmut sync.RWMutex
)

func adder() {
	rwmut.Lock()
	value++
	rwmut.Unlock()
}

func reader() {
	rwmut.RLock()
	fmt.Println(value)
	rwmut.RUnlock()
}

func ex8() {
	for i := 0; i < 5; i++ {
		go reader()
	}

	for j := 0; j < 3; j++ {
		go adder()
	}

	time.Sleep(2 * time.Second)
	fmt.Println(value)
}
