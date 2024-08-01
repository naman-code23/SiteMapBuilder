package link

import (
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(file io.Reader)([]Link, error) {

	// parse the file
	doc, err := html.Parse(file)
	handleError(err)

	//traverse(doc)
	links := returnLinks(doc)
	return links, nil
}

func returnLinks(n *html.Node) []Link {
	var links []Link

	if n.Type == html.ElementNode {
		if n.Data == "a" {
			l1 := Link{}
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					l1.Href = attr.Val
				}
			}
			s := getData(n)
			l1.Text = s
			links = append(links, l1)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, returnLinks(c)...)
	}

	return links
}
// gets data from node and its children
func getData(n *html.Node) string {
	s := ""
	if n.Type == html.TextNode {
		s = n.Data
		s=strings.TrimSuffix(s, "\n")
		s= strings.TrimPrefix(s, "\n")
		s= strings.TrimSpace(s)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s += getData(c)
	}
	return s
}

// traverse function to print the element nodes
func traverse(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Printf("Element: <%s>, Type: %v\n", n.Data, n.Type)
		if len(n.Attr) > 0 {
			fmt.Println("Attributes:")
			for _, attr := range n.Attr {
				fmt.Printf(" - %s: %s\n", attr.Key, attr.Val)
			}
		}
	}

	// Traverse the child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c)
	}
}
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
