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

type updateWishlistRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Hidden      *bool   `json:"hidden,omitempty"`
}

var validateUpdateWishlistRequest = validation.Struct(
	validation.StructField(
		"name",
		func(request *updateWishlistRequest) any {
			return request.Name
		},
		validation.Optional[string](
			validation.String(
				validation.NotEmpty(),
				validation.MaxLength(255),
			),
		),
	),
)

var ErrWishlistUpdateFailed = errors.New("wishlist update failed")

func update(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	wishlistId, err := ctx.WishlistId()

	if err != nil {
		ctx.Log(log.Debug, "failed to parse id", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wishlist not found", nil, nil)

		return
	}

	_updateWishlistRequest := new(updateWishlistRequest)

	if err = internalHttp.ReadJson(request, _updateWishlistRequest); err != nil {
		ctx.Log(log.Debug, "failed to read request", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", nil, nil)

		return
	}

	if errs := validateUpdateWishlistRequest(_updateWishlistRequest); len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	var updatedWishlist *wishlist.Wishlist

	err = wishlistStore.Update(ctx.WishlistStoreContext(), wishlistId, func(_wishlist *wishlist.Wishlist) (*wishlist.Wishlist, error) {
		if _wishlist == nil || _wishlist.Wisher().Id() != ctx.WisherFromUser().Id() {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wishlist not found", nil, nil)

			return nil, ErrWishlistUpdateFailed
		}

		updatedWishlist = _wishlist

		if _updateWishlistRequest.Name != nil {
			updatedWishlist.Rename(*_updateWishlistRequest.Name)
		}

		if _updateWishlistRequest.Description != nil {
			updatedWishlist.UpdateDescription(*_updateWishlistRequest.Description)
		}

		if _updateWishlistRequest.Hidden != nil {
			if *_updateWishlistRequest.Hidden {
				updatedWishlist.Hide()
			} else {
				updatedWishlist.Show()
			}
		}

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
