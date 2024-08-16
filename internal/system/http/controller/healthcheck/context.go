package healthcheck

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	httpContext "api/internal/system/http/context"
	"net/http"
)

type Context struct {
	request    *http.Request
	requestCtx *httpContext.RequestContext
	systemCtx  *systemContext.Context
}

func NewContext(request *http.Request) *Context {
	requestCtx := request.Context().(*httpContext.RequestContext)
	systemCtx := requestCtx.SystemContext().WithLabels(log.NewLabel("controller", "healthcheck"))

	return &Context{
		request:    request,
		requestCtx: request.Context().(*httpContext.RequestContext),
		systemCtx:  systemCtx,
	}
}

func (ctx *Context) Log(level log.Level, message string, labels ...*log.Label) {
	ctx.systemCtx.Log(level, message, labels...)
}
