package splash_api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

const LATEST = "http://www.splashbase.co/api/v1/images/latest?images_only=true"
const RANDOM = "http://www.splashbase.co/api/v1/images/random?images_only=true"

var (
	ErrResolvePath = errors.New("Failed to resolve path")
	ErrDir         = errors.New("Slash API directory could not be created")
	ErrNoPath      = errors.New("No resolvable path to image")
)

var setupDir func() (string, error)

func init() {

	setupDir = func() (string, error) {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Join(usr.HomeDir, ".splash_api")
		info, err := os.Stat(dir)
		if err == nil {
			return dir, err
		}
		_ = info

		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return dir, err
			}
		}
		return dir, err
	}
}

//{"id":390,"url":"https://splashbase.s3.amazonaws.com/unsplash/regular/tumblr_n6rzkfxeOR1st5lhmo1_1280.jpg","large_url":"https://splashbase.s3.amazonaws.com/unsplash/large/1rSok7e","source_id":244,"copyright":"CC0","site":"unsplash"}

type Image struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Large string `json:"large_url"`
	Site  string `json:"site"`
	Path  string `json:",omitempty"`
}

func (i *Image) buildPath() (string, error) {

	base := filepath.Base(i.Large)
	// Unsplash sometimes shortcuts endpoints, provide a unique
	if base == "download" || base == "." {
		base = filepath.Base(i.URL)
	}
	// Extensions are truncated
	if ext := filepath.Ext(base); ext == "" {
		ext = filepath.Ext(i.URL)
		if ext == "" {
			// We could check the buffer, or just assume .jpg
			base += ".jpg"
		} else {
			base += ext
		}
	}
	if base == "." || len(base) == 0 {
		return base, ErrResolvePath
	}
	return base, nil
}

func (i *Image) Fetch(params Params) error {
	path, err := i.buildPath()
	if err != nil {
		return err
	}
	if !params.Fetch {
		return nil
	}
	i.Path = filepath.Join(params.Dir, path)
	log.Println("Fetching...", i.Path)
	_, err = os.Stat(i.Path)
	if !os.IsNotExist(err) {
		return os.ErrExist
	}
	log.Printf("New image from %s: %s\n", i.Site, i.Path)
	file, err := os.Create(i.Path)
	if err != nil {
		return err
	}
	url := i.Large
	if len(url) == 0 {
		url = i.URL
	}
	if len(url) == 0 {
		return ErrNoPath
	}
	resp, err := http.Get(i.Large)
	if err != nil {
		log.Printf("% #v\n", i)
		return err
	}

	defer resp.Body.Close()
	for {
		_, err := io.Copy(file, resp.Body)
		if err != nil && err != io.EOF {
			return err
		} else {
			break
		}
	}

	return nil
}

type Images struct {
	Images []Image `json:"images"`
}

// Params for sending options to methods
type Params struct {
	Endpoint string
	// Type indicates struct to Unmarshal to expect
	// 0 Image{}
	// 1 []Image{}
	Type  uint8
	Dir   string
	Fetch bool
}

func get(params Params) (*Images, error) {
	resp, err := http.Get(params.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var imgs Images
	dec := json.NewDecoder(resp.Body)
	switch params.Type {
	case 0:
		img := Image{}
		err := dec.Decode(&img)
		if err != nil {
			return nil, err
		}
		imgs.Images = append(imgs.Images, img)
	case 1:
		for {
			err := dec.Decode(&imgs)
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
		}
	}

	return &imgs, nil
}

func Get(params Params) error {
	dir, err := setupDir()
	if err != nil {
		return err
	}
	if len(dir) == 0 {
		log.Fatal(ErrDir)
	}
	params.Dir = dir

	if len(params.Endpoint) == 0 {
		return errors.New("Endpoint not available")
	}
	switch params.Endpoint {
	case LATEST:
		params.Type = 1
	case RANDOM:
		params.Type = 0
	}
	imgs, err := get(params)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, img := range imgs.Images {
		wg.Add(1)
		go func(img Image) {
			defer wg.Done()
			err := img.Fetch(params)
			if os.IsExist(err) {
				log.Println("File already exists", img.Path)
			}
		}(img)
	}
	wg.Wait()
	return nil
}
