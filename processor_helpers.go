// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"net/http"

	"github.com/aporeto-inc/elemental"
)

// CheckAuthentication checks if the current context has been authenticated if there is any authenticator registered.
//
// If it is not authenticated it stops the normal processing execution flow, and will write the Unauthorized response to the given writer.
// If not Authenticator is set, then it will always return true.
//
// This is mostly used by autogenerated code, and you should not need to use it manually.
func CheckAuthentication(ctx *Context, w http.ResponseWriter) bool {

	server := DefaultBahamut()
	authenticator, err := server.Authenticator()
	if err != nil {
		return true
	}

	ok, err := authenticator.IsAuthenticated(ctx)

	if err != nil {
		WriteHTTPError(w, http.StatusInternalServerError, elemental.NewError("Internal Server Error", err.Error(), "http", http.StatusInternalServerError))
		return false
	}

	if !ok {
		WriteHTTPError(w, http.StatusUnauthorized, elemental.NewError("Unauthorized", "You are not authorized to access this resource.", "http", http.StatusUnauthorized))
		return false
	}

	return true
}

// CheckAuthorization checks if the current context has been authorized if there is any authorizer registered.
//
// If it is not authorized it stops the normal processing execution flow, and will write the Unauthorized response to the given writer.
// If not Authorizer is set, then it will always return true.
//
// This is mostly used by autogenerated code, and you should not need to use it manually.
func CheckAuthorization(ctx *Context, w http.ResponseWriter) bool {

	server := DefaultBahamut()
	authorizer, err := server.Authorizer()
	if err != nil {
		return true
	}

	ok, err := authorizer.IsAuthorized(ctx)

	if err != nil {
		WriteHTTPError(w, http.StatusInternalServerError, elemental.NewError("Internal Server Error", err.Error(), "http", http.StatusInternalServerError))
		return false
	}

	if !ok {
		WriteHTTPError(w, http.StatusForbidden, elemental.NewError("Forbidden", "You are not allowed to access this resource.", "http", http.StatusForbidden))
		return false
	}

	return true
}
