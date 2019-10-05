package main

import (
	"flag"
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"regexp"
)

type Instruct struct {
	EntryData EntryData `json:"entry_data"`
}

type EntryData struct {
	PostPage []PostPage `json:"PostPage"`
}

type PostPage struct {
	Graphql Graphql `json:"graphql"`
}

type Graphql struct {
	ShortcodeMedia ShortcodeMedia `json:"shortcode_media"`
}

type ShortcodeMedia struct {
	Shortcode             string                `json:"shortcode"`
	Dimensions            Dimensions            `json:"dimensions"`
	IsVideo               bool                  `json:"is_video"`
	DisplayURL            string                `json:"display_url"`
	EdgeSidecarToChildren EdgeSidecarToChildren `json:"edge_sidecar_to_children"`
}

type Dimensions struct {
	Height int64 `json:"height"`
	Width  int64 `json:"width"`
}

type EdgeSidecarToChildren struct {
	Edges []EdgeSidecarToChildrenEdge `json:"edges"`
}

type EdgeSidecarToChildrenEdge struct {
	Node FluffyNode `json:"node"`
}

type FluffyNode struct {
	Shortcode  string     `json:"shortcode"`
	Dimensions Dimensions `json:"dimensions"`
	IsVideo    bool       `json:"is_video"`
	DisplayURL string     `json:"display_url"`
}

func main() {
	url := flag.String("url", "", "url of instagram photo")
	shortcode := flag.String("shortcode", "", "shortcode of instagram photo")
	onlyUrl := flag.Bool("only-url", false, "show only url")
	flag.Parse()

	if *shortcode == "" && *url == "" {
		fmt.Println("-url flag or -shortcode flag, were defined incorrectly")
	} else {
		if *url == "" {
			*url = "https://www.instagram.com/p/" + *shortcode
		}

		resp, err := http.Get(*url)
		if err != nil {
			fmt.Println("Error[1]: Network error.")
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)

		pattern := regexp.MustCompile("_sharedData = {.*}")
		matches := pattern.FindString(string(body))
		rawJson := matches[14:len(matches)]

		instruct := Instruct{}
		json.Unmarshal([]byte(rawJson), &instruct)

		shortcode := instruct.EntryData.PostPage[0].Graphql.ShortcodeMedia
		if *onlyUrl {
			if shortcode.EdgeSidecarToChildren.Edges == nil {
				fmt.Println(shortcode.DisplayURL)
				return
			}

			for _, edge := range shortcode.EdgeSidecarToChildren.Edges {
				fmt.Println(edge.Node.DisplayURL)
			}
			return
		}

		if shortcode.EdgeSidecarToChildren.Edges == nil {
			resp, err := http.Get(shortcode.DisplayURL)
			if err != nil {
				fmt.Println("Error[2]: Network error.")
				return
			}
			body, _ := ioutil.ReadAll(resp.Body)

			ext := "jpg"
			if shortcode.IsVideo {
				ext = "mp4"
			}

			filename := fmt.Sprintf("%s__%d_%d.%s", shortcode.Shortcode, shortcode.Dimensions.Height, shortcode.Dimensions.Width, ext)
			ioutil.WriteFile(filename, body, 0644)
			return
		}

		os.Mkdir(shortcode.Shortcode, os.ModePerm)
		for _, edge := range shortcode.EdgeSidecarToChildren.Edges {
			resp, err := http.Get(edge.Node.DisplayURL)
			if err != nil {
				fmt.Println("Error[2]: Network error.")
				return
			}
			body, _ := ioutil.ReadAll(resp.Body)

			ext := "jpg"
			if edge.Node.IsVideo {
				ext = "mp4"
			}

			filename := fmt.Sprintf("%s/%s__%d_%d.%s", shortcode.Shortcode, edge.Node.Shortcode, edge.Node.Dimensions.Height, edge.Node.Dimensions.Width, ext)
			ioutil.WriteFile(filename, body, 0644)
		}
	}
}
