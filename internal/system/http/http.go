package http

import (
	"api/internal/model/config"
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	systemHttpServer "api/internal/system/http/server"
	"errors"
	"net/http"
)

type Register func(server *systemHttpServer.Server)

func serveHttp(
	ctx *systemContext.Context,
	_config *config.HttpServerConfig,
	controls ...Register,
) error {
	ctx = ctx.WithLabels(log.NewLabel("address", _config.Address()))

	_systemHttpServer := systemHttpServer.New(ctx)

	for _, control := range controls {
		control(_systemHttpServer)
	}

	server := &http.Server{
		Addr:    _config.Address(),
		Handler: _systemHttpServer,
	}

	serverErrChan := make(chan error, 1)

	go func() {
		ctx.Log(log.Info, "http server is running")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ctx.Log(
				log.Error,
				"http server is broken",
				log.NewLabel("error", err.Error()),
			)

			serverErrChan <- err
		}

		ctx.Log(log.Info, "http server is stopped")
		serverErrChan <- nil
	}()

	select {
	case <-ctx.Done():
		ctx.Log(log.Info, "http server is shutting down")
	case serverErr := <-serverErrChan:
		return serverErr
	}

	return server.Shutdown(ctx)
}
