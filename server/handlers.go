package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"sync"
)

type tmplIndex struct{}

func (s *Server) hndlrIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.tmplMap["/index"].Execute(w, &tmplIndex{}); err != nil {
		notify(err)
		s.sendError(w, r, http.StatusInternalServerError)
	}
}

type prmsTranslate struct {
	Text string `json:"text"`
}

func (s *Server) hndlrTranslate(w http.ResponseWriter, r *http.Request) {
	prms := &prmsTranslate{}
	if err := s.getPrms(r, prms); err != nil {
		s.sendError(w, r, http.StatusBadRequest)
		return
	}

	result := map[string]interface{}{}
	blob := map[string]interface{}{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		out, err := s.googleTranslate(prms.Text)
		if err != nil {
			log.Printf("python exec err: %v", err)
			return
		}

		result["google"] = out
	}()

	go func() {
		defer wg.Done()
		tagged := s.tagGlossaryEntries(prms.Text)
		translated, err := s.googleTranslate(tagged)
		if err != nil {
			log.Printf("python exec err: %v", err)
			return
		}

		replaced, err := s.replaceGlossaryEntries(string(translated))
		if err != nil {
			log.Printf("replace tag err: %v", err)
			return
		}

		result["glossary"] = replaced
	}()

	go func() {
		defer wg.Done()
		cmd := exec.Command("/usr/bin/python3", s.pythonScriptsPath+"/blob.py", "-t", prms.Text)
		out, err := cmd.Output()
		if err != nil {
			log.Printf("python exec err: %v", err)
			return
		}

		if err := json.Unmarshal(out, &blob); err != nil {
			log.Println(err)
			return
		}

		result["blob"] = blob
	}()

	wg.Wait()

	s.send(w, r, result)
}

func (s *Server) googleTranslate(text string) (string, error) {
	cmd := exec.Command("/usr/bin/python3", s.pythonScriptsPath+"/google_translator.py", "-t", text)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
