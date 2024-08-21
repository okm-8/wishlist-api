package iam

import (
	"api/internal/model/log"
	"api/internal/model/user"
	internalHttp "api/internal/service/integration/http"
	userStore "api/internal/service/integration/pgx/user"
	"api/internal/system/validation"
	"errors"
	"net/http"
)

type iamResponse struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

func me(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	user := ctx.User()

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, iamResponse{
		Id:      user.Id().String(),
		Email:   user.Email(),
		Name:    user.Name(),
		IsAdmin: user.IsAdmin(),
	}, nil)
}

type updateRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

var validateUpdateRequest = validation.Struct(
	validation.StructField(
		"name",
		func(request *updateRequest) any {
			return request.Name
		},
		validation.Optional[string](
			validation.String(
				validation.NotEmpty(),
				validation.MaxLength(255),
			),
		),
	),
	validation.StructField(
		"email",
		func(request *updateRequest) any {
			return request.Email
		},
		validation.Optional[string](
			validation.String(
				validation.NotEmpty(),
				validation.Email(),
				validation.MaxLength(255),
			),
		),
	),
	validation.StructField(
		"password",
		func(request *updateRequest) any {
			return request.Password
		},
		validation.Optional[string](
			validation.String(
				validation.NotEmpty(),
			),
		),
	),
)

var ErrMeUpdateFailed = errors.New("me update failed")

func updateMe(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	_user := ctx.User()

	_updateRequest := new(updateRequest)

	err := internalHttp.ReadJson(request, _updateRequest)
	if err != nil {
		ctx.Log(log.Debug, "failed to read request", log.NewLabel("error", err.Error()))

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", nil, nil)

		return
	}

	errs := validateUpdateRequest(_updateRequest)
	if len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	var updatedUser *user.User
	err = userStore.Update(ctx.UserStoreContext(), _user.Id(), func(_user *user.User) (*user.User, error) {
		if _user == nil {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "user not found", nil, nil)

			return nil, ErrMeUpdateFailed
		}

		updatedUser = _user

		if _updateRequest.Name != nil {
			updatedUser.ChangeName(*_updateRequest.Name)
		}

		if _updateRequest.Email != nil {
			updatedUser.ChangeEmail(*_updateRequest.Email)
		}

		if _updateRequest.Password != nil {
			ctx.IamContext().SetUserPassword(updatedUser, *_updateRequest.Password)
		}

		return updatedUser, nil
	})

	if err != nil {
		if errors.Is(err, userStore.ErrEmailExists) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusConflict, "email already exists", nil, nil)
		} else if !errors.Is(err, ErrMeUpdateFailed) {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		}

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, iamResponse{
		Id:      updatedUser.Id().String(),
		Email:   updatedUser.Email(),
		Name:    updatedUser.Name(),
		IsAdmin: updatedUser.IsAdmin(),
	}, nil)
}
