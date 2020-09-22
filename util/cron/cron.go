package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func CronSetup()  {
	fmt.Println("cron starting...")

	c := cron.New()

	c.AddFunc("0 15 17 * * *", func() {
		fmt.Println("cron doing 1...")
		//model.TestCreateUser()
		for i := 200; i <= 400 ; i++ {
			//go controller.RequestColumn(i)
			//time.Sleep(1 * time.Second)
		}
	})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <- t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
