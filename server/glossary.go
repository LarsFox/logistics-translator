package server

import (
	"os"
	"strings"
)

func newGlossary(path string) ([]glossaryTerm, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	glossary := []glossaryTerm{}
	for _, line := range strings.Split(string(b), "\n") {
		separated := strings.Split(line, ";")
		if len(separated) != 5 {
			continue
		}

		glossary = append(glossary, glossaryTerm{
			lang:       separated[0],
			term:       strings.ToLower(separated[1]),
			gloss:      separated[2],
			example:    separated[3],
			definition: separated[4],
		})
	}

	return glossary, nil
}

func (s *Server) findExample(text, lang string) string {
	text = strings.ToLower(text)
	for _, g := range s.glossary {
		if g.lang == lang && g.term == text {
			return g.example + "\n"
		}
	}

	return ""
}

func (s *Server) findGloss(text, lang string) string {
	text = strings.ToLower(text)
	for _, g := range s.glossary {
		if g.lang == lang && g.term == text {
			return g.gloss
		}
	}

	return "—"
}
