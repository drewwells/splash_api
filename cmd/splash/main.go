package main

import (
	"flag"
	"log"

	"github.com/drewwells/splash_api"
	//"github.com/mitchellh/cli"
)

func main() {
	flag.Parse()
	ep := splash_api.LATEST
	if Rndm {
		ep = splash_api.RANDOM
	}
	p := splash_api.Params{
		Endpoint: ep,
		Fetch:    !Check,
	}
	err := splash_api.Get(p)
	if err != nil {
		log.Fatal("Failed to retrieve latest images:\n    ", err)
	}

}

var (
	Help  bool
	Check bool
	Rndm  bool
)

func init() {
	flag.BoolVar(&Help, "help", false, "Show help")
	flag.BoolVar(&Help, "h", false, "Show help")
	flag.BoolVar(&Check, "check", false, "Check for but do not download new iamges")
	flag.BoolVar(&Rndm, "random", false, "Pull a random image")
}
