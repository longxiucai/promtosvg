package main

import (
	httpserver "github.com/longxiucai/promtosvg/pkg/http"
)

func main() {
	err := httpserver.RunWeb()
	if err != nil {
		panic(err)
	}
}
