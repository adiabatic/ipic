// © 2018 Nathan Galt
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var html = `<!DOCTYPE HTML>
<html>
<head>
	<meta charset='UTF-8'>
	<title>Pictures for {{.Title}}</title>
	<style>
		html {
			font-family: system-ui, sans-serif;
		}

		@media (prefers-color-scheme: dark) {
			html {
				color: white;
				background: black;
			}
		}

		h1, p {
			text-align: center;
		}
	</style>
</head>
<body>
<h1>Pictures for “{{.Title}}”</h1>
{{ range .Results }}
	<a href='{{.ArtworkURLBig}}'><img src='{{.ArtworkURL100}}' alt='' title='{{.Name}}'></a>
{{ else }}
	<p>No results.</p>
{{ end }}
</body>
</html>
`

// SearchResponse is what comes from Apple’s search response. It’s also used for populating the spat-out HTML.
type SearchResponse struct {
	Title   string
	Results []struct {
		ArtworkURL100  string `json:"artworkUrl100"`
		ArtworkURLBig  string
		Name           string
		TrackName      string `json:"trackName"`
		CollectionName string `json:"collectionName"`
	} `json:"results"`
}

var iOS = flag.Bool("i", false, "iOS app")
var macOS = flag.Bool("m", false, "macOS app")
var album = flag.Bool("a", false, "album")
var film = flag.Bool("f", false, "film")
var tvShow = flag.Bool("t", false, "TV show")
var book = flag.Bool("b", false, "book")
var help = flag.Bool("h", false, "show this help message")

func main() {
	flag.Parse()

	v := url.Values{}
	v.Set("size", "4096")      // might as well ask for the biggest thing imaginable, right?
	v.Set("name", "trackName") // most of these ask for trackName
	v.Set("term", strings.Join(flag.Args(), " "))

	if *help {
		fmt.Fprintf(os.Stderr, "Usage: ipic (-i | -m | -a | -f | -t | -b | -h) SEARCH_TERM\n\nGenerates a page of thumbnails and links to larger images for items in the iTunes/App/macOS App Stores.\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nOnly one option is allowed. The HTML file for the generated webpage is saved to ~/Desktop.\n")
		return
	} else if *iOS {
		// 512
		v.Set("media", "softare")
		v.Set("entity", "software")
	} else if *macOS {
		//size = 512
		v.Set("media", "software")
		v.Set("entity", "macSoftware")
	} else if *album {
		v.Set("media", "music")
		v.Set("entity", "album")
		v.Set("name", "collectionName")
	} else if *film {
		v.Set("media", "movie")
		v.Set("entity", "movie")
	} else if *tvShow {
		v.Set("media", "tvShow")
		v.Set("entity", "tvSeason")
		v.Set("name", "collectionName")
	} else if *book {
		v.Set("media", "ebook")
		v.Set("entity", "ebook")
	}

	var buf bytes.Buffer
	r, err := http.NewRequest("GET", "https://itunes.apple.com/search", &buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn’t create request: %s.\n", err)
		os.Exit(1)
	}
	r.URL.RawQuery = v.Encode()

	r.Form = v

	c := http.Client{Timeout: 3 * time.Second}

	resp, err := c.Do(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong while trying to get Apple data: %s\n", err)
		os.Exit(2)
	}
	defer r.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/javascript") {
		fmt.Fprintf(os.Stderr, "Response doesn’t even claim to be text/javascript. Suspicious!\n")
	}

	sr := SearchResponse{}

	err = json.NewDecoder(resp.Body).Decode(&sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn’t decode the presumably-JSON response: %s\n", err)
		os.Exit(3)
	}

	for i := 0; i != len(sr.Results); i++ {
		result := sr.Results[i]
		size := fmt.Sprintf("%sx%s", v.Get("size"), v.Get("size"))
		sr.Results[i].ArtworkURLBig = strings.Replace(result.ArtworkURL100, "100x100", size, -1)

		// coalesce different name types into .Name
		sr.Results[i].Name = result.CollectionName
		if result.TrackName != "" {
			sr.Results[i].Name = result.TrackName
		}
	}

	sr.Title = v.Get("term")

	fn := filepath.Join(os.Getenv("HOME"), "Desktop", "“"+sr.Title+"” Images.html")

	file, err := os.Create(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn’t open «%s» for writing: %s", fn, err)
	}
	defer file.Close()

	// template stuff
	tmpl := template.Must(template.New("html").Parse(html))

	err = tmpl.Execute(file, sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn’t whip up the HTML: %s\n", err)
	}

	if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", fn)
		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn’t automatically open %s in a browser: %s", fn, err)
		}
	}
}
