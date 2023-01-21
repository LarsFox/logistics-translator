package server

import (
	"fmt"
	"os"
	"regexp"
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
		if len(separated) != 2 {
			continue
		}

		glossary = append(glossary, glossaryTerm{
			gloss: separated[0],
			term:  separated[1],
		})
		glossary = append(glossary, glossaryTerm{
			gloss: separated[1],
			term:  separated[0],
		})
	}
	return glossary, nil
}

func (s *Server) tagGlossaryEntries(text string) string {
	lowered := strings.ToLower(text)
	for i, term := range s.glossary {
		lowered = strings.ReplaceAll(lowered, term.term, fmt.Sprintf("<term-%d>%s</term-%d>", i, term.term, i))
	}

	return lowered
}

func (s *Server) replaceGlossaryEntries(translated string) (string, error) {
	lowered := strings.ToLower(translated)

	for i, term := range s.glossary {
		re, err := regexp.Compile(fmt.Sprintf("<term-%d>.*?</term-%d>", i, i))
		if err != nil {
			return "", err
		}

		replaced := fmt.Sprintf("<span class='match'>%s</span>", term.gloss)
		for _, match := range re.FindAllString(lowered, -1) {
			lowered = strings.ReplaceAll(lowered, match, replaced)
		}
	}

	return lowered, nil
}
