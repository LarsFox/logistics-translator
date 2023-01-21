package main

import (
	"os"

	"github.com/LarsFox/logistics-translator/server"
)

func main() {
	s := server.New(
		os.Getenv("LOGISTICS_TRANSLATOR_GLOSSARY_PATH"),
		os.Getenv("LOGISTICS_TRANSLATOR_HTML_PATH"),
		os.Getenv("LOGISTICS_TRANSLATOR_PYTHON_SCRIPTS_PATH"),
	)
	s.Listen(":9090")
}
