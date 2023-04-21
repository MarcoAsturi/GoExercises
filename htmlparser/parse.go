package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// struct for representing a link in html document
type Link struct {
	Href string
	Text string
}

// reading html content and extracting link as slices
func extractLinks(n *html.Node) ([]Link, error) {

	var links []Link

	// ausiliary func for link extraction
	var extractLinksFromNode func(*html.Node)

	extractLinksFromNode = func(n *html.Node) {
		// if <a> tag is present, extracts link
		if n.Type == html.ElementNode && n.Data == "a" {
			link := buildLinkFromNode(n)
			links = append(links, link)
		}

		// for each child of the current html node, call the previous func
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinksFromNode(c)
		}
	}

	// call the func for the entire html document
	extractLinksFromNode(n)

	// checkin errors
	if len(links) == 0 {
		return nil, fmt.Errorf("no links found in document")
	}

	// return the list of extract links
	return links, nil
}

// building a link from a html node <a>
func buildLinkFromNode(n *html.Node) Link {
	var link Link

	// iterating all attributes of html node
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}

	// initializing a new string
	text := ""

	// extraxcting html text and concatenating to the String
	extractTextFromNode(n, &text)

	// remove empty space (first and last)
	link.Text = strings.TrimSpace(text)
	return link

}

// extract text contained in html node and his childs recorsively
func extractTextFromNode(n *html.Node, text *string) {

	// if node is text type, then add the content in output string
	if n.Type == html.TextNode {
		*text = n.Data
	}

	// check recorively all the childs of current node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTextFromNode(c, text)
	}
}
