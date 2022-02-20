package server

import (
	"os"
	"regexp"
	"strings"
)

func newGlossary(path string) (map[string]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	glossary := map[string]string{}
	for _, line := range strings.Split(string(b), "\n") {
		separated := strings.Split(line, ";")
		if len(separated) != 2 {
			continue
		}

		glossary[separated[0]] = separated[1]
	}
	return glossary, nil
}

var puncts = regexp.MustCompile("[a-zа-я-]+")

func (s *Server) findGlossaryEntries(text string) []string {
	var entries []string
	duplicates := map[string]bool{}
	for _, word := range strings.Split(strings.ToLower(text), " ") {
		trimmed := strings.Trim(word, " ")
		if duplicates[trimmed] {
			continue
		}

		found := puncts.FindAllString(trimmed, -1)
		for _, match := range found {
			value, ok := s.glossary[match]
			if !ok {
				continue
			}

			if duplicates[match] {
				continue
			}

			entries = append(entries, value)
			duplicates[match] = true
		}
	}

	return entries
}
