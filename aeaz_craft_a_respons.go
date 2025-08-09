package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Analyzer struct {
	client http.Client
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		client: http.Client{},
	}
}

func (a *Analyzer) Analyze(url string) (*AnalysisResult, error) {
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	metrics := &AnalysisResult{}

	// Get page title
	metrics.Title = doc.Find("title").Text()

	// Get meta viewport
	metrics.Viewport = doc.Find("meta[name=viewport]").AttrOr("content", "")

	// Get media queries
	metrics.MediaQueries = make([]string, 0)
	doc.Find("style, link[rel=stylesheet]").Each(func(i int, s *goquery.Selection) {
		matches := regexFindAll(s.Text(), `@media ([^{}]+) {`)
		metrics.MediaQueries = append(metrics.MediaQueries, matches...)
	})

	// Get responsive images
	metrics.ResponsiveImages = make([]string, 0)
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		srcset, exists := s.Attr("srcset")
		if exists {
			metrics.ResponsiveImages = append(metrics.ResponsiveImages, srcset)
		}
	})

	return metrics, nil
}

type AnalysisResult struct {
	Title         string
	Viewport      string
	MediaQueries  []string
	ResponsiveImages []string
}

func regexFindAll(text string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)
	result := make([]string, 0)
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func main() {
	analyzer := NewAnalyzer()
	result, err := analyzer.Analyze("https://example.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %s\n", result.Title)
	fmt.Printf("Viewport: %s\n", result.Viewport)
	fmt.Println("Media Queries:")
	for _, mq := range result.MediaQueries {
		fmt.Println(mq)
	}
	fmt.Println("Responsive Images:")
	for _, ri := range result.ResponsiveImages {
		fmt.Println(ri)
	}
}