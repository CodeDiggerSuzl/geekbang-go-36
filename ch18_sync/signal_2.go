package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var mailbox uint8
	var lock sync.Mutex
	sendCond := sync.NewCond(&lock)
	recvCond := sync.NewCond(&lock)

	send := func(id, idx int) {
		lock.Lock()
		for mailbox == 1 {
			sendCond.Wait()
		}
		log.Printf("send [%d-%d]: the mailbox is empty", id, idx)
		mailbox = 1
		log.Printf("sender [%d-%d]: the letter has been send ", id, idx)
		lock.Unlock()
		recvCond.Broadcast()
	}

	recv := func(id, idx int) {
		lock.Lock()
		for mailbox == 0 {
			recvCond.Wait()
		}
		log.Printf("recv [%d-%d]: the mailbox is full", id, idx)
		mailbox = 0
		log.Printf("recv [%d-%d]: the letter has ben received", id, idx)
		lock.Unlock()
		sendCond.Broadcast()
	}

	sign := make(chan struct{}, 3)
	max := 6
	go func(id, max int) {
		defer func() {
			sign <- struct{}{}
		}()
		for i := 1; i <= max; i++ {
			time.Sleep(time.Millisecond * 500)
			send(id, i)
		}
	}(0, max)

	go func(id, max int) {
		defer func() {
			sign <- struct{}{}
		}()
		for j := 0; j < max; j++ {
			time.Sleep(time.Millisecond * 200)
			recv(id, j)
		}
	}(1, max/2)
	go func(id, max int) {
		defer func() {
			sign <- struct{}{}
		}()
		for k := 0; k < max; k++ {
			time.Sleep(time.Millisecond * 200)
			recv(id, k)
		}
	}(2, max/2)
	<-sign
	<-sign
	<-sign
}
