package wishlist

import (
	"api/internal/model/log"
	"api/internal/model/wishlist"
	internalHttp "api/internal/service/integration/http"
	wishlistStore "api/internal/service/integration/pgx/wishlist"
	"api/internal/system/validation"
	"errors"
	"net/http"
)

type addWishRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var validateAddWishRequest = validation.Struct(
	validation.StructField(
		"name",
		func(request *addWishRequest) any {
			return request.Name
		},
		validation.String(
			validation.NotEmpty(),
			validation.MaxLength(255),
		),
	),
)

func addWish(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	wishlistId, err := ctx.WishlistId()

	if err != nil {
		ctx.Log(log.Debug, "failed to parse id", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wishlist not found", nil, nil)

		return
	}

	_addWishRequest := new(addWishRequest)

	if err = internalHttp.ReadJson(request, _addWishRequest); err != nil {
		ctx.Log(log.Debug, "failed to read request", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", nil, nil)

		return
	}

	if errs := validateAddWishRequest(_addWishRequest); len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	var updatedWishlist *wishlist.Wishlist
	err = wishlistStore.Update(ctx.WishlistStoreContext(), wishlistId, func(_wishlist *wishlist.Wishlist) (*wishlist.Wishlist, error) {
		if _wishlist == nil {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wishlist not found", nil, nil)

			return nil, ErrWishlistUpdateFailed
		}

		if _wishlist.Wisher().Id() != ctx.WisherFromUser().Id() {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wishlist not found", nil, nil)

			return nil, ErrWishlistUpdateFailed
		}

		updatedWishlist = _wishlist

		_ = updatedWishlist.AddWish(_addWishRequest.Name, _addWishRequest.Description)

		return updatedWishlist, nil
	})

	if err != nil {
		if !errors.Is(err, ErrWishlistUpdateFailed) {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		}

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWishlist(updatedWishlist), nil)
}
