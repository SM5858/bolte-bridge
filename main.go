package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"bolte-bridge/internal/config"
	"bolte-bridge/internal/store"
)

// run performs one invocation of the bridge.
func run() error {
	ctx := context.Background()

	cfg, err := config.Load(os.Args[1:], config.DefaultSections...)
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	if err := store.Init(ctx, cfg.Store); err != nil {
		return fmt.Errorf("initialise store: %w", err)
	}
	defer func() { _ = store.Client().Close(ctx) }()

	log.Printf("store initialised at %s", cfg.Store.SQLite.Path)
	// TODO: implement the rest of the relay here.
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("bolte-bridge: %v", err)
	}
}
