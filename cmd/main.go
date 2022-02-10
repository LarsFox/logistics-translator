package main

import (
	"github.com/LarsFox/logistics-translator/server"
)

func main() {
	s := server.New()
	s.Listen(":9090")
}
