package splash_api

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func init() {
	setupDir = func() (string, error) {
		dir, err := ioutil.TempDir(os.TempDir(), "splash_dir")
		log.Println("created", dir)
		return dir, err
	}
}

func TestDir(t *testing.T) {
	setupDir()
}

func TestGet_latest(t *testing.T) {
	p := Params{
		Endpoint: LATEST,
	}
	Get(p)
}

func TestGet_rndm(t *testing.T) {
	p := Params{
		Endpoint: RANDOM,
	}
	err := Get(p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetch_shortpath(t *testing.T) {
	p := Params{
		Fetch: false,
	}
	i := Image{ID: 1880,
		URL:   "https://splashbase.s3.amazonaws.com/unsplash/regular/unsplash_524010c76b52a_1.JPG%3Ffit%3Dcrop%26fm%3Djpg%26h%3D650%26q%3D75%26w%3D950",
		Large: "", Site: "unsplash", Path: ""}

	err := i.Fetch(p)
	if err != nil {
		log.Fatal(err)
	}
}

func TestFetch_nopath(t *testing.T) {
	p := Params{
		Fetch: false,
	}
	i := Image{
		ID: 2018, URL: "", Large: "",
		Site: "unsplash",
	}

	err := i.Fetch(p)

	if err != ErrResolvePath {
		log.Fatal(err)
	}
}
