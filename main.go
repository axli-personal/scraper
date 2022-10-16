package main

import (
	"github.com/gocolly/colly"
	"log"
	"strings"
	"time"
)

func setCookie(request *colly.Request) {
	request.Headers.Set("Cookie", cookies)
}

func logError(response *colly.Response, err error) {
	log.Print("Fail to send request:", err)
}

var books []Book
var comments = make(map[string][]Comment)

func main() {
	defaultCollector := colly.NewCollector(
		colly.CacheDir("__cache"),
	)

	defaultCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       10 * time.Second,
		RandomDelay: 20 * time.Second,
	})

	bookCollector := defaultCollector.Clone()
	commentCollector := defaultCollector.Clone()

	bookCollector.OnRequest(setCookie)
	commentCollector.OnRequest(setCookie)

	bookCollector.OnError(logError)
	commentCollector.OnError(logError)

	// 1. Define rule.

	bookCollector.OnHTML("#content", func(element *colly.HTMLElement) {
		element.ForEach(".item", func(i int, element *colly.HTMLElement) {
			book := Book{
				Name:   element.ChildAttr("div.pl2 a", "title"),
				Author: extractAuthor(element.ChildText("p.pl")),
				Rate:   element.ChildText(".rating_nums"),
				Quote:  element.ChildText(".quote span"),
				Link:   element.ChildAttr("div.pl2 a", "href"),
			}
			books = append(books, book)
		})

		nextUrl := element.ChildAttr(".paginator .next a", "href")
		if nextUrl != "" && len(books) < 250 {
			element.Request.Visit(nextUrl)
		}
	})

	commentCollector.OnHTML("#comments", func(element *colly.HTMLElement) {
		id := extractBookId(element.Request.URL.Path)
		element.ForEach(".comment", func(i int, element *colly.HTMLElement) {
			comment := Comment{
				User: element.ChildText("a:first-child"),
				Rate: element.ChildAttr(".rating", "title"),
				Text: strings.ReplaceAll(element.ChildText(".comment-content span"), "\n", ""),
			}
			comments[id] = append(comments[id], comment)
		})

		nextUrl := element.ChildAttr("#paginator a[data-page=next]", "href")
		if nextUrl != "" && len(comments[id]) < 200 {
			element.Request.Visit(nextUrl)
		}
	})

	// 2. Start visit.

	bookCollector.Visit("https://book.douban.com/top250")
	for i, book := range books {
		log.Print("Book:", i+1)
		commentCollector.Visit(book.Link + "comments/")
	}

	// 3. Write to file.

	WriteBooks()
	WriteComments()

	log.Print("Done")
}
