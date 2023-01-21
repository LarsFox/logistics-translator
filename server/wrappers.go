package server

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type wrapper func(http.Handler) http.Handler

/*
// wrapBodyParams ...
func (s *Server) wrapBodyParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var handlerName string
		if route := mux.CurrentRoute(r); route != nil {
			handlerName = route.GetName()
			if handlerName == "" {
				notify(errUnnamedMethod, bugsnag.MetaData{"ctx": {"uri": r.URL}})
				s.sendError(w, r, http.StatusInternalServerError)
				return
			}
		}

		prms, err := getPrmsScaffold(handlerName)
		if err != nil {
			s.sendError(w, r, http.StatusInternalServerError)
			return
		}

		if ok := m.getPrms(r, prms); !ok {
			s.sendError(w, r, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyBodyParams, prms)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/

// wrapContentTypeJSON wraps json content type validation.
func (s *Server) wrapContentTypeJSON(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if strings.ToLower(ct) != "application/json; charset=utf-8" {
			s.sendError(w, r, http.StatusBadRequest)
			return
		}
		inner.ServeHTTP(w, r)
	})
}

// wrapEasterEggHeader adds easter egg headers.
func (s *Server) wrapEasterEggHeader(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("hey", "what are you trying to find here")
		w.Header().Set("Leon-Motovskikh", "is the best")
		w.Header().Set("x-files", "Scully approves")
		inner.ServeHTTP(w, r)
	})
}

// wrapDuration logs the response duration.
func (s *Server) wrapDuration(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()

		log.Println(r.Method, r.RequestURI, duration)
	})
}

// wrapReadTemplates re-reads the templates, so they could be refreshed easily.
// TODO: disable for prod version.
func (s *Server) wrapReadTemplates(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmplMap, err := readTemplates(s.htmlPath)
		if err != nil {
			log.Println(err)
			inner.ServeHTTP(w, r)
			return
		}

		s.tmplMap = tmplMap
		inner.ServeHTTP(w, r)
	})
}

// wrapRecover recovers panics, should one occur.
func (s *Server) wrapRecover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			rec := recover()
			if rec != nil {
				switch t := rec.(type) {
				case error:
					err = t
				default:
					err = errUnknownError
				}

				log.Println(err, map[string]interface{}{
					"stacktrace": string(debug.Stack()),
					"uri":        r.RequestURI,
				})
			}
		}()
		h.ServeHTTP(w, r)
	})
}

/*
// wrapScheme ...
func (s *Server) wrapScheme(inner http.Handler) http.Handler {
	if m.local {
		return inner
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := r.Header.Get("X-Forwarded-Proto")
		if scheme != "https" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			return
		}

		if strings.HasPrefix(r.Host, "www.") {
			http.Redirect(w, r, "https://"+r.Host[4:]+r.RequestURI, http.StatusMovedPermanently)
			return
		}

		inner.ServeHTTP(w, r)
	})
}
*/
