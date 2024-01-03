package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
)



func main(){	
	_, err := flags.Parse(&tag)
	if err != nil {
		return
	}
	CheckBin(tag.Postcli)
	app := New()
	app.Connect(tag.Address)
	time.Sleep(2 * time.Second)
	go app.Task_Start()

	go func ()  {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
		for{
			s := <-c
			switch s {
			case syscall.SIGINT:
				app.Close()
			case syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT:
				app.Close()
			}
		}
	}()
	<- app.close
}