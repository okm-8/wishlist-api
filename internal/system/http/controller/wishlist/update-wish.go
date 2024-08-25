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

type UpdateWishRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Hidden      *bool   `json:"hidden"`
	Fulfilled   *bool   `json:"fulfilled"`
}

var validateUpdateWishRequest = validation.Struct(
	validation.StructField(
		"name",
		func(request *UpdateWishRequest) any {
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

var ErrWishUpdateFailed = errors.New("failed to update wish")

func updateWish(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	wishId, err := ctx.WishId()

	if err != nil {
		ctx.Log(log.Debug, "failed to parse id", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wish not found", nil, nil)

		return
	}

	_updateWishRequest := new(UpdateWishRequest)

	err = internalHttp.ReadJson(request, _updateWishRequest)
	if err != nil {
		ctx.Log(log.Debug, "failed to read request", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", nil, nil)

		return
	}

	errs := validateUpdateWishRequest(_updateWishRequest)
	if len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	var updatedWish *wishlist.Wish
	wisher := ctx.WisherFromUser()
	err = wishlistStore.UpdateWish(ctx.WishlistStoreContext(), wishId, func(_wish *wishlist.Wish) (*wishlist.Wish, error) {
		if _wish == nil {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wish not found", nil, nil)

			return nil, ErrWishUpdateFailed
		}

		if _wish.Wishlist().Wisher().Id() != wisher.Id() {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wish not found", nil, nil)

			return nil, ErrWishUpdateFailed
		}

		updatedWish = _wish

		if _updateWishRequest.Name != nil {
			err = updatedWish.Rename(*_updateWishRequest.Name)

			if err != nil {
				return nil, err
			}
		}

		if _updateWishRequest.Description != nil {
			err = updatedWish.UpdateDescription(*_updateWishRequest.Description)

			if err != nil {
				return nil, err
			}
		}

		if _updateWishRequest.Hidden != nil {
			if *_updateWishRequest.Hidden {
				updatedWish.Hide()
			} else {
				updatedWish.Show()
			}
		}

		if _updateWishRequest.Fulfilled != nil {
			wisherAsAssignee := wishlist.RestoreAssignee(
				wishlist.AssigneeId(wisher.Id()),
				wisher.Name(),
				wisher.Email(),
			)

			if *_updateWishRequest.Fulfilled {
				if !updatedWish.Promised() {
					updatedWish.Promise(wisherAsAssignee)
				}

				if err = updatedWish.Fulfill(); err != nil {
					return nil, err
				}
			} else {
				updatedWish.Rollback()

				assignee := updatedWish.Assignee()
				if assignee != nil && assignee.Id() == wisherAsAssignee.Id() {
					err = updatedWish.Renege(wisherAsAssignee)
				}
			}
		}

		return updatedWish, nil
	})

	if err != nil {
		if errors.Is(err, wishlist.ErrWishAlreadyFulfilled) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusForbidden, "wish already fulfilled", nil, nil)
		} else if errors.Is(err, wishlist.ErrWishAlreadyPromised) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusForbidden, "wish already promised", nil, nil)
		} else if !errors.Is(err, ErrWishUpdateFailed) {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		}

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWish(updatedWish), nil)
}
