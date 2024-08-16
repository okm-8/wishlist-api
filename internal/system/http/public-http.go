package http

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	"api/internal/system/http/controller/healthcheck"
	"api/internal/system/http/controller/iam"
)

func PublicHttp(ctx *systemContext.Context) error {
	return serveHttp(
		ctx.WithLabels(log.NewLabel("system", "http"), log.NewLabel("server", "public")),
		ctx.Config().Http().PublicServer(),
		healthcheck.Register,
		iam.RegisterPublic,
	)
}
