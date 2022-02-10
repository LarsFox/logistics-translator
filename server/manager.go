package server

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	router  *mux.Router
	tmplMap map[string]*template.Template
}

// route is a single path for a mux handler.
type route struct {
	Method   string
	Name     string
	Path     string
	Handler  http.HandlerFunc
	Wrappers []wrapper
}

// New ...
func New() *Server {
	s := &Server{
		router:  mux.NewRouter().StrictSlash(true),
		tmplMap: readTemplates("html/"),
	}

	// Notes.
	s.addHandlers([]route{
		{
			Method:   "POST",
			Path:     "/translate/",
			Name:     "Notes",
			Handler:  s.hndlrTranslate,
			Wrappers: []wrapper{s.wrapContentTypeJSON, s.wrapReadTemplates},
		},
	})

	var indexHandler http.Handler = http.HandlerFunc(s.hndlrIndex)
	for _, wrapper := range []wrapper{s.wrapDuration, s.wrapEasterEggHeader, s.wrapReadTemplates, s.wrapRecover} {
		indexHandler = wrapper(indexHandler)
	}

	s.router.PathPrefix("/").Handler(indexHandler)

	return s
}

// Listen servers the addr.
func (s *Server) Listen(addr string) error {
	log.Println("Translator started on addr", addr)
	return http.ListenAndServe(addr, s.router)
}

// addHandlers adds routes to mux.
func (s *Server) addHandlers(routes []route) {
	essentialWrappers := []wrapper{s.wrapEasterEggHeader, s.wrapDuration, s.wrapRecover}
	for _, r := range routes {
		var wrapper http.Handler = r.Handler
		for _, w := range r.Wrappers {
			wrapper = w(wrapper)
		}
		for _, w := range essentialWrappers {
			wrapper = w(wrapper)
		}
		s.router.Methods(r.Method).Path(r.Path).Name(r.Name).Handler(wrapper)
	}
}

// send ...
func (s *Server) send(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"ok":     true,
		"result": data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if strings.Contains(err.Error(), "write: broken pipe") {
			return
		}
		notify(err)
	}

	log.Println(r.RequestURI, resp)
}

// sendError ...
func (s *Server) sendError(w http.ResponseWriter, _ *http.Request, code int) {
	w.WriteHeader(code)
}

// getPrms unmarshals request body and validates the result.
// Body size is limited to 1 MB.
func (s *Server) getPrms(r *http.Request, prms interface{}) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	if err := json.Unmarshal(body, prms); err != nil {
		log.Println(err, map[string]interface{}{
			"uri":  r.RequestURI,
			"body": string(body),
		})
		return err
	}

	return nil
}
