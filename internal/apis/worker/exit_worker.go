package api

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *WorkerServer) ExitWorker(ctx echo.Context) error {
	defer func() {
		time.Sleep(1 * time.Second)
		s.server.Shutdown(context.Background())
	}()

	s.worker.Cleanup()

	return nil
}
