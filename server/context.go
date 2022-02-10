package server

import (
	"net/http"
)

type contextKey int

const (
	ctxKeyBodyParams contextKey = iota
	ctxKeyTester
)

func getFromCtxBodyParams(r *http.Request) interface{} {
	return r.Context().Value(ctxKeyBodyParams)
}
