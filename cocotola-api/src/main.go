package main

import (
	"fmt"
	"log/slog"

	"github.com/kujilabo/cocotola-1.21/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola-1.21/proto"
)

func main() {
	logger := slog.Default()
	logger.Info(fmt.Sprintf("%+v", proto.HelloRequest{}))

	logger.Info("")
	logger.Info(domain.Lang2EN.String())
	logger.Info("Hello")
	logger.Info("Hello")
	service.A()
}
