package main

import (
	"app-websocket/internal/components"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	components, err := components.InitComponents()
	if err != nil {
		logrus.Fatal(err)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		return components.HTTPServer.Run(ctx)
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case s := <-sigQuit:
			logrus.Info("Captured signal", s.String())
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	err = eg.Wait()

}
