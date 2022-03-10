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
	if err := s.tmplMap["index"].Execute(w, &tmplIndex{}); err != nil {
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

	result := map[string]interface{}{}
	blob := map[string]interface{}{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		cmd := exec.Command("/usr/bin/python3", "python/google_translator.py", "-t", prms.Text)
		out, err := cmd.Output()
		if err != nil {
			log.Printf("python exec err: %v", err)
			return
		}

		result["google"] = string(out)
	}()

	go func() {
		defer wg.Done()
		result["glossary"] = s.findGlossaryEntries(prms.Text)
	}()

	go func() {
		defer wg.Done()
		cmd := exec.Command("/usr/bin/python3", "python/blob.py", "-t", prms.Text)
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
