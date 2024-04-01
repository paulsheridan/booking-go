package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/paulsheridan/booking-go/application"
)

func main() {
	app := application.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
