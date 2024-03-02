package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type Product struct {
	id    string
	URL   string
	price string
}

func main() {
	var totalPages = 0

	hostURL := "https://www.flipkart.com"

	products := []Product{}

	c := colly.NewCollector(colly.Async(true))

	c.Limit(&colly.LimitRule{Parallelism: 4})
	c.SetRequestTimeout(30 * time.Second)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Body....................", string(r.Body))
		log.Println("StatusCode..................", r.StatusCode)
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("._2MImiq", func(p *colly.HTMLElement) {
		if totalPages == 0 {
			pagesString := strings.Split(p.ChildText("span"), " ")
			pageVisitUrl := p.ChildAttr("a", "href")

			pageVisitUrl = pageVisitUrl[:len(pageVisitUrl)-1]

			asd := pagesString[len(pagesString)-1]
			first, _ := strconv.Atoi(strings.TrimSuffix(asd, "Next"))
			totalPages = first

			for i := 2; i <= totalPages; i++ {
				c.Visit(hostURL + pageVisitUrl + strconv.Itoa(i))
			}
		}
	})

	c.OnHTML("._1AtVbE ._13oc-S ", func(h *colly.HTMLElement) {
		product := Product{}

		h.ForEach("div", func(i int, div *colly.HTMLElement) {

			id := div.Attr("data-id")
			if len(id) > 0 {
				product.id = id
				product.price = div.ChildText("._30jeq3")
				products = append(products, product)
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, "scraped")
	})

	c.Visit(hostURL)
	c.Wait()

	for i, product := range products {
		fmt.Printf("Product %d: %s , %s\n", i+1, product.id, product.price)
	}
}
