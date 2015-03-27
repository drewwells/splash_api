package main

import (
	"flag"
	"log"

	"github.com/drewwells/splash_api"
	//"github.com/mitchellh/cli"
)

func main() {

	p := splash_api.Params{
		Endpoint: splash_api.RANDOM,
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
)

func init() {
	flag.BoolVar(&Help, "help", false, "Show help")
	flag.BoolVar(&Help, "h", false, "Show help")
	flag.BoolVar(&Check, "check", false, "Check for but do not download new iamges")
}
