package http

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	"api/internal/system/http/controller/healthcheck"
	"api/internal/system/http/controller/iam"
	"api/internal/system/http/controller/root"
	"api/internal/system/http/controller/wishlist"
)

func PublicHttp(ctx *systemContext.Context) error {
	return serveHttp(
		ctx.WithLabels(log.NewLabel("system", "http"), log.NewLabel("server", "public")),
		ctx.Config().Http().PublicServer(),
		root.Register,
		healthcheck.Register,
		iam.RegisterPublic,
		wishlist.Register,
	)
}
