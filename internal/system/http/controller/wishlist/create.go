package wishlist

import (
	"api/internal/model/log"
	"api/internal/model/wishlist"
	internalHttp "api/internal/service/integration/http"
	wishlistStore "api/internal/service/integration/pgx/wishlist"
	"api/internal/system/validation"
	"net/http"
)

type createRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var validateCreateRequest = validation.Struct(
	validation.StructField(
		"name",
		func(request *createRequest) any {
			return request.Name
		},
		validation.String(
			validation.NotEmpty(),
			validation.MaxLength(255),
		),
	),
)

func create(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	_createRequest := new(createRequest)

	err := internalHttp.ReadJson(request, _createRequest)
	if err != nil {
		ctx.Log(log.Debug, "failed to read request", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", nil, nil)

		return
	}

	errs := validateCreateRequest(_createRequest)

	if len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	_wishlist := wishlist.New(
		ctx.WisherFromUser(),
		_createRequest.Name,
		_createRequest.Description,
	)

	err = wishlistStore.Store(ctx.WishlistStoreContext(), _wishlist)

	if err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, serializeWishlist(_wishlist), nil)
}
