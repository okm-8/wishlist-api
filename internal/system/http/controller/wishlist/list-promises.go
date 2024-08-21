package wishlist

import (
	"api/internal/model/wishlist"
	internalHttp "api/internal/service/integration/http"
	wishlistStore "api/internal/service/integration/pgx/wishlist"
	"net/http"
)

func listPromises(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	user := ctx.User()

	wishlists, err := wishlistStore.GetByAssigneeId(ctx.WishlistStoreContext(), wishlist.AssigneeId(user.Id()))

	if err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWishLists(wishlists), nil)
}
