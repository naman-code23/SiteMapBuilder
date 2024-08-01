package main

import (
	"SiteMapBuilder/link"
	"SiteMapBuilder/server"
	"encoding/json"
	"flag"
	"fmt"
	_ "log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type linkNode struct {
	Name     string      `json:"name"`
	Url      string      `json:"url,omitempty"`
	Children []*linkNode `json:"children,omitempty"`
}

var (
	visited   = make(map[string]bool)
	visitedMu sync.Mutex
	maxDepth  = 5
)

func main() {
	urlFlag := flag.String("url", "", "URL to start crawling from")
	flag.Parse()
	if *urlFlag == "" {
		fmt.Println("Please provide a URL using the -url flag")
		os.Exit(1)
	}
	fmt.Println("Started...")
	start := time.Now()
	//test url = https://echo.labstack.com/
	url := *urlFlag

	rootNode := &linkNode{"Home", url, []*linkNode{}}

	var wg sync.WaitGroup
	wg.Add(1)
	go traverse(rootNode, 0, &wg)
	wg.Wait()

	jsonData, err := json.MarshalIndent(rootNode, "", "  ")
	handleError(err)
	err = os.WriteFile("tree.json", jsonData, 0644)
	handleError(err)
	fmt.Println("Active on :8080")
	fmt.Println("Time taken: ", time.Since(start))
	server.Serve()
}

func get(url string) []link.Link {
	resp, err := http.Get(url)
	handleError(err)
	defer resp.Body.Close()
	linksTemp, err := link.Parse(resp.Body)
	handleError(err)
	reqUrl := resp.Request.URL
	baseUrl := reqUrl.Scheme + "://" + reqUrl.Host
	links := filter(linksTemp, baseUrl)
	return links
}

func filter(links []link.Link, domain string) []link.Link {
	var filteredLinks []link.Link
	for _, l := range links {
		if strings.HasPrefix(l.Href, domain) {
			filteredLinks = append(filteredLinks, l)
		} else if strings.HasPrefix(l.Href, "/") {
			l.Href = domain + l.Href
			filteredLinks = append(filteredLinks, l)
		} else if strings.HasPrefix(l.Href, "./") {
			l.Href = domain + l.Href[1:]
			filteredLinks = append(filteredLinks, l)
		}
	}
	return filteredLinks
}

func traverse(n *linkNode, depth int, parentWg *sync.WaitGroup) {
	defer parentWg.Done()

	if depth >= maxDepth {
		return
	}

	visitedMu.Lock()
	if visited[n.Url] {
		visitedMu.Unlock()
		return
	}
	visited[n.Url] = true
	visitedMu.Unlock()

	links := get(n.Url)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, l := range links {
		wg.Add(1)
		go func(l link.Link) {
			node := &linkNode{l.Text, l.Href, []*linkNode{}}

			childWg := &sync.WaitGroup{}
			childWg.Add(1)
			go traverse(node, depth+1, childWg)
			childWg.Wait()

			mu.Lock()
			n.Children = append(n.Children, node)
			mu.Unlock()

			wg.Done()
		}(l)
	}

	wg.Wait()
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
