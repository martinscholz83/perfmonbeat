package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/maddin2016/perfmonbeat/beater"
)

func main() {
	err := beat.Run("perfmonbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
