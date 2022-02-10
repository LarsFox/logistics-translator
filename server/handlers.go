package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

type tmplMapData struct {
	CanvasStyle   string   // дополнительный класс холста.
	CustomCSS     []string // дополнительные ЦСС-файлы.
	CSSLanguages  []string // ЦСС-файлы языков.
	Description   string   // Тест на знание ...
	International bool     // есть ли версия на английском.
	Name          string   // sakha.
	Language      string   // язык шаблона.
	Title         string   // Как хорошо ....
	Subjects      string   // regions, countries.
	URL           string   // все, что после домена.
}

func (s *Server) hndlrIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.tmplMap["index"].Execute(w, &tmplMapData{}); err != nil {
		notify(err)
		s.sendError(w, r, http.StatusInternalServerError)
	}
}

type prmsTranslate struct {
	Text string `json:"text"`
}

type respTranslate struct {
	Text string `json:"text"`
}

// todo: queues
func (s *Server) hndlrTranslate(w http.ResponseWriter, r *http.Request) {
	prms := &prmsTranslate{}
	if err := s.getPrms(r, prms); err != nil {
		s.sendError(w, r, http.StatusBadRequest)
		return
	}

	cmd := exec.Command("/usr/local/bin/python3", "python/blob.py", "-t", prms.Text)
	// cmd.Stderr = log.Default().Writer()
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		s.sendError(w, r, http.StatusInternalServerError)
		return
	}

	blob := map[string]interface{}{}
	if err := json.Unmarshal(out, &blob); err != nil {
		log.Println(err)
		s.sendError(w, r, http.StatusInternalServerError)
		return
	}

	s.send(w, r, blob)
}
