package fetchhtml

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func PollUrlForID(url string, id string) (tagContents string, err error) {
	resp, getErr := http.Get(url)
	if getErr != nil {
		err = fmt.Errorf("Failed to get given url. URL: %v; ERR: %v", url, getErr)
		return "", err
	}

	defer resp.Body.Close()
	httpBody := resp.Body
	node, parseErr := html.Parse(httpBody)
	if parseErr != nil {
		err = fmt.Errorf("Failed to parse html document. ERR: %v", parseErr)
		return "", err
	}

	document := goquery.NewDocumentFromNode(node)
	djString, queryErr := document.Find("#" + id).Html()
	if queryErr != nil {
		err = fmt.Errorf("Failed to find correct html tag: %v; ERR: %v", id, queryErr)
		return "", err
	}

	return djString, nil
}
