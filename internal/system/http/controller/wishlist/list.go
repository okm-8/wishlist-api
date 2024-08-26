package wishlist

import (
	"api/internal/model/wishlist"
	internalHttp "api/internal/service/integration/http"
	wishlistStore "api/internal/service/integration/pgx/wishlist"
	"net/http"
)

type ListItemResponse struct {
}

func list(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	user := ctx.MustUser()

	wishlists, err := wishlistStore.GetByWisherId(ctx.WishlistStoreContext(), wishlist.WisherId(user.Id()))

	if err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWishLists(wishlists), nil)
}
