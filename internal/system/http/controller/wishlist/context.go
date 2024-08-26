package wishlist

import (
	"api/internal/model/log"
	"api/internal/model/user"
	"api/internal/model/wishlist"
	wishlistStore "api/internal/service/integration/pgx/wishlist"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service/integration"
	httpContext "api/internal/system/http/context"
	"net/http"
)

type Context struct {
	request          *http.Request
	requestCtx       *httpContext.RequestContext
	systemCtx        *systemContext.Context
	wishlistStoreCtx wishlistStore.Context
}

func NewContext(request *http.Request) *Context {
	requestCtx := request.Context().(*httpContext.RequestContext).WithLabels(log.NewLabel("controller", "wishlist"))
	systemCtx := requestCtx.SystemContext()

	return &Context{
		request:          request,
		requestCtx:       requestCtx,
		systemCtx:        systemCtx,
		wishlistStoreCtx: integration.NewWishlistStoreContext(systemCtx),
	}
}

func (ctx *Context) WishlistStoreContext() wishlistStore.Context {
	return ctx.wishlistStoreCtx
}

func (ctx *Context) WishlistId() (wishlist.Id, error) {
	return wishlist.ParseId(ctx.request.PathValue("wishlistId"))
}

func (ctx *Context) WishId() (wishlist.WishId, error) {
	return wishlist.ParseWishId(ctx.request.PathValue("wishId"))
}

func (ctx *Context) WisherId() (wishlist.WisherId, error) {
	return wishlist.ParseWisherId(ctx.request.PathValue("wisherId"))
}

func (ctx *Context) User() *user.User {
	_user, err := ctx.requestCtx.User()

	if err != nil {
		panic(err) // This should never happen
	}

	return _user
}

func (ctx *Context) WisherFromUser() *wishlist.Wisher {
	_user := ctx.User()

	return wishlist.RestoreWisher(
		wishlist.WisherId(_user.Id()),
		_user.Name(),
		_user.Email(),
	)
}

func (ctx *Context) AssigneeFromUser() *wishlist.Assignee {
	_user := ctx.User()

	return wishlist.RestoreAssignee(
		wishlist.AssigneeId(_user.Id()),
		_user.Name(),
		_user.Email(),
	)
}

func (ctx *Context) Log(level log.Level, message string, labels ...*log.Label) {
	ctx.requestCtx.Log(level, message, labels...)
}
