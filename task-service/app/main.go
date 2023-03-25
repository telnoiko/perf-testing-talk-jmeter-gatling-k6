package main

import (
	"context"
	"os"
	"os/signal"
	"task-service/app/rest"
	"time"
)

func main() {

	rest := rest.New()
	go func() {
		rest.Start(":1323")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rest.Stop(ctx)
}
