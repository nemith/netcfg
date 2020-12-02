package main

import (
	"log"
	"netcfg/interp"
)

func main() {
	if err := interp.ExecFile("main.star"); err != nil {
		log.Fatal(err)
	}
}
