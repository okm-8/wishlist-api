package wishlist

import (
	"api/internal/model/log"
	"api/internal/model/wishlist"
	internalHttp "api/internal/service/integration/http"
	wishlistStore "api/internal/service/integration/pgx/wishlist"
	"errors"
	"net/http"
)

func wishers(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)

	var currentWisher *wishlist.Wisher

	_, err := ctx.User()
	if err == nil {
		currentWisher = ctx.WisherFromUser()
	}

	wishers, err := wishlistStore.GetWishers(ctx.WishlistStoreContext())
	if err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	filteredWishers := make([]*wishlist.Wisher, 0, len(wishers))

	for _, wisher := range wishers {
		if currentWisher == nil || wisher.Id() != currentWisher.Id() {
			filteredWishers = append(filteredWishers, wisher)
		}
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWishers(filteredWishers), nil)
}

func wisherWishlists(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)

	wisherId, err := ctx.WisherId()
	if err != nil {
		ctx.Log(log.Debug, "failed to parse wisher id", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wisher not found", nil, nil)

		return
	}

	wishlists, err := wishlistStore.GetByWisherIdActive(ctx.WishlistStoreContext(), wisherId)
	if err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWishLists(wishlists), nil)
}

func promiseWish(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)

	wisherId, err := ctx.WisherId()
	if err != nil {
		ctx.Log(log.Debug, "failed to parse wisher id", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wisher not found", nil, nil)
	}

	wishlistId, err := ctx.WishlistId()
	if err != nil {
		ctx.Log(log.Debug, "failed to parse wishlist id", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wishlist not found", nil, nil)
	}

	wishId, err := ctx.WishId()
	if err != nil {
		ctx.Log(log.Debug, "failed to parse wish id", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wish not found", nil, nil)

		return
	}

	assignee := ctx.AssigneeFromUser()
	var promisedWish *wishlist.Wish

	err = wishlistStore.UpdateWish(ctx.WishlistStoreContext(), wishId, func(_wish *wishlist.Wish) (*wishlist.Wish, error) {
		if _wish == nil || _wish.Wishlist().Wisher().Id() != wisherId || _wish.Wishlist().Id() != wishlistId {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wish not found", nil, nil)

			return nil, ErrWishUpdateFailed
		}

		promisedWish = _wish

		err = promisedWish.Promise(assignee)
		if err != nil {
			return nil, err
		}

		return promisedWish, nil
	})

	if err != nil {
		if !errors.Is(err, ErrWishUpdateFailed) {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		} else if errors.Is(err, wishlist.ErrWishAlreadyPromised) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusConflict, "wish already promised", nil, nil)
		} else if errors.Is(err, wishlist.ErrWishAlreadyFulfilled) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusConflict, "wish already fulfilled", nil, nil)
		}

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWish(promisedWish), nil)
}

func renageWish(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)

	wisherId, err := ctx.WisherId()
	if err != nil {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wisher not found", nil, nil)
	}

	wishlistId, err := ctx.WishlistId()
	if err != nil {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wishlist not found", nil, nil)
	}

	wishId, err := ctx.WishId()
	if err != nil {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wish not found", nil, nil)

		return
	}

	assignee := ctx.AssigneeFromUser()
	var renagedWish *wishlist.Wish

	err = wishlistStore.UpdateWish(ctx.WishlistStoreContext(), wishId, func(_wish *wishlist.Wish) (*wishlist.Wish, error) {
		if _wish == nil || _wish.Wishlist().Wisher().Id() != wisherId || _wish.Wishlist().Id() != wishlistId {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "wish not found", nil, nil)

			return nil, ErrWishUpdateFailed
		}

		renagedWish = _wish

		err = renagedWish.Renege(assignee)
		if err != nil {
			return nil, err
		}

		return renagedWish, nil
	})

	if err != nil {
		if !errors.Is(err, ErrWishUpdateFailed) {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		} else if errors.Is(err, wishlist.ErrorWishNotPromised) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusConflict, "wish not promised", nil, nil)
		} else if errors.Is(err, wishlist.ErrWishAlreadyFulfilled) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusConflict, "wish already fulfilled", nil, nil)
		}

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWish(renagedWish), nil)
}
