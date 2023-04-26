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
	From string `json:"from"`
	Text string `json:"text"`
	To   string `json:"to"`
}

func (s *Server) hndlrTranslate(w http.ResponseWriter, r *http.Request) {
	prms := &prmsTranslate{}
	if err := s.getPrms(r, prms); err != nil {
		s.sendError(w, r, http.StatusBadRequest)
		return
	}

	result := map[string]interface{}{}
	blob := map[string]interface{}{}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd := exec.Command(
			"/usr/bin/python3",
			s.pythonScriptsPath+"/google_translator.py",
			"-t", prms.Text,
			"-s", prms.From,
			"-d", prms.To,
		)
		out, err := cmd.Output()
		if err != nil {
			log.Printf("python exec err: %v", err)
			return
		}

		result["google"] = string(out)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd := exec.Command(
			"/usr/bin/python3",
			s.pythonScriptsPath+"/reverso.py",
			"-t", prms.Text,
			"-s", prms.From,
			"-d", prms.To,
			"-c", "e",
		)
		out, err := cmd.Output()
		if err != nil {
			log.Printf("python exec err: %v", err)
			return
		}

		result["example"] = s.findExample(prms.Text, prms.To) + string(out)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd := exec.Command(
			"/usr/bin/python3",
			s.pythonScriptsPath+"/reverso.py",
			"-t", prms.Text,
			"-s", prms.From,
			"-d", prms.To,
			"-c", "t",
		)
		out, err := cmd.Output()
		if err != nil {
			log.Printf("python exec err: %v", err)
			return
		}

		result["translation"] = string(out)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd := exec.Command(
			"/usr/bin/python3",
			s.pythonScriptsPath+"/blob.py",
			"-t", prms.Text,
			"-s", prms.From,
		)
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		result["glossary"] = s.findGloss(prms.Text, prms.To)
	}()

	wg.Wait()

	s.send(w, r, result)
}
