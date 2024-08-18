package iam

import (
	internalHttp "api/internal/service/integration/http"
	"api/internal/system/http/server"
	"net/http"
)

func RegisterCommon(_server *server.Server) {
	_server.HandleFunc("/iam/login", internalHttp.Method(internalHttp.MethodMap{
		http.MethodPost: login,
	}))
	_server.HandleFunc("/iam/logout", internalHttp.Method(internalHttp.MethodMap{
		http.MethodGet:  logout,
		http.MethodPost: logout,
	}))
}

func RegisterPublic(_server *server.Server) {
	RegisterCommon(_server)

	_server.HandleFunc("/iam/signup", internalHttp.Method(internalHttp.MethodMap{
		http.MethodPost: signUp,
	}))
}

func RegisterPrivate(_server *server.Server) {
	RegisterCommon(_server)

	_server.HandleFunc("/iam/signup", internalHttp.Method(internalHttp.MethodMap{
		http.MethodPost: signupNoToken,
	}))
}
