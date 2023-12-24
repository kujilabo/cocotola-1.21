package main

import (
	"log/slog"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/usecase"
)

func main() {
	var _ = new(usecase.Authentication)
	logger := slog.Default()
	logger.Info("AUTH 2")
}
