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
		if len(separated) != 4 {
			continue
		}

		glossary = append(glossary, glossaryTerm{
			term:       separated[0],
			gloss:      separated[1],
			example:    separated[2],
			definition: separated[3],
		})
		glossary = append(glossary, glossaryTerm{
			gloss:      separated[0],
			term:       separated[1],
			example:    separated[2],
			definition: separated[3],
		})
	}

	return glossary, nil
}

func (s *Server) findGloss(text string) string {
	for _, g := range s.glossary {
		if g.term == text {
			return g.gloss
		}
	}

	return "—"
}
