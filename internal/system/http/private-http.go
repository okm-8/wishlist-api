package http

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	"api/internal/system/http/controller/healthcheck"
	"api/internal/system/http/controller/iam"
)

func PrivateHttp(ctx *systemContext.Context) error {
	return serveHttp(
		ctx.WithLabels(log.NewLabel("system", "http"), log.NewLabel("server", "private")),
		ctx.Config().Http().PrivateServer(),
		healthcheck.Register,
		iam.RegisterPrivate,
	)
}
