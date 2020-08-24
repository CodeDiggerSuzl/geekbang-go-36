package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var mailbox uint8
	var lock sync.RWMutex // 读写锁
	sendCond := sync.NewCond(&lock)
	recvCond := sync.NewCond(lock.RLocker())
	sign := make(chan struct{}, 3) // 用于传递演示
	max := 5
	// sender
	go func(max int) {
		defer func() {
			log.Println("sender defer ......")
			sign <- struct{}{}
		}()
		for i := 0; i < max; i++ {
			time.Sleep(time.Microsecond * 500)
			lock.Lock()
			for mailbox == 1 {
				sendCond.Wait()
			}
			log.Printf("sender [%d]: the mail box is empty:", i)
			mailbox = 1
			log.Printf("sender: [%d] the letter has been send", i)
			lock.Unlock()
			recvCond.Signal()
		}
	}(max)
	// recv
	go func(max int) {
		defer func() {
			sign <- struct{}{}
			log.Println("receiver defer ......")
		}()
		for j := 0; j < max; j++ {
			time.Sleep(time.Millisecond * 500)
			lock.RLock()
			for mailbox == 0 {
				recvCond.Wait()
			}
			log.Printf("receiver [%d]: the mail box is full.", j)
			mailbox = 0
			log.Printf("receiver [%d]: the letter has been received.", j)
			lock.RUnlock()
			sendCond.Signal()
		}
	}(max)

	<-sign
	<-sign
}
