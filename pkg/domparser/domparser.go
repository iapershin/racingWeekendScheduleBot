package domparser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var url = _url(os.Getenv("URL"), "https://vk.com/world_of_speed")

func BuildMessage(url string) (string, error) {
	html, err := makeRequest(url)
	if err != nil {
		log.Println(err)
		return "", err
	}

	parsedPage, err := parsePage(html)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return formatter(parsedPage), nil
}

func makeRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ERROR Make request to %s : %s", url, err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ERROR Reading page %s : %s", url, err.Error())
	}
	return string(body), nil
}

func parsePage(html string) (string, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("ERROR Parsing DOM : %s", err.Error())
	}
	result := dom.Find(".pi_text").First().Text()
	return result, nil
}

func formatter(content string) string {
	expanderText := [2]string{"Показать полностью...", "See more"}
	for _, text := range expanderText {
		content = strings.ReplaceAll(content, text, "")
	}
	formatted := strings.ReplaceAll(content, "#", "\n\n#")
	return formatted
}

func _url(u, default_value string) string {
	if u == "" {
		return default_value
	} else {
		return u
	}
}
