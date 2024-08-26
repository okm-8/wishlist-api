package wishlist

import (
	internalHttp "api/internal/service/integration/http"
	"api/internal/system/http/middleware"
	"api/internal/system/http/server"
	"net/http"
)

func Register(_server *server.Server) {
	_server.HandleFunc("/wishlist", middleware.AuthorizedOnly(
		internalHttp.Method(internalHttp.MethodMap{
			http.MethodGet:  list,
			http.MethodPost: create,
		}),
	))
	_server.HandleFunc("/wishlist/{wishlistId}", middleware.AuthorizedOnly(
		internalHttp.Method(internalHttp.MethodMap{
			http.MethodPatch: update,
		}),
	))
	_server.HandleFunc("/wishlist/{wishlistId}/wish", middleware.AuthorizedOnly(
		internalHttp.Method(internalHttp.MethodMap{
			http.MethodPost: addWish,
		}),
	))
	_server.HandleFunc("/wishlist/{wishlistId}/wish/{wishId}", middleware.AuthorizedOnly(
		internalHttp.Method(internalHttp.MethodMap{
			http.MethodPatch: updateWish,
		}),
	))
	_server.HandleFunc("/wishlist/promises", middleware.AuthorizedOnly(
		internalHttp.Method(internalHttp.MethodMap{
			http.MethodGet: listPromises,
		}),
	))
	_server.HandleFunc("/wisher", internalHttp.Method(internalHttp.MethodMap{
		http.MethodGet: wishers,
	}))
	_server.HandleFunc("/wisher/{wisherId}/wishlist", internalHttp.Method(internalHttp.MethodMap{
		http.MethodGet: wisherWishlists,
	}))
	_server.HandleFunc("/wisher/{wisherId}/wishlist/{wishlistId}/wish/{wishId}/promise", middleware.AuthorizedOnly(
		internalHttp.Method(internalHttp.MethodMap{
			http.MethodPost:   promiseWish,
			http.MethodDelete: renageWish,
		}),
	))
}
