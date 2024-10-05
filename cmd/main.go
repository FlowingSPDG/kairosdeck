package main

import (
	"context"
	"os"

	"github.com/FlowingSPDG/kairosdeck/handlers"
	"github.com/FlowingSPDG/streamdeck"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	h := handlers.SetupHandlers(client)
	return h.Run(ctx)
}
