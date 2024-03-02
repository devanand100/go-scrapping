package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Product struct {
	id    string
	URL   string
	price string
}

func main() {

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))

	defer cancel()
	var nodes []*cdp.Node

	var products []Product

	// var pages string
	if err := chromedp.Run(ctx, chromedp.Navigate("https://www.flipkart.com/search?q=furniture")); err != nil {
		log.Fatal(err)
	}

	var pagesString string

	if err := chromedp.Run(ctx, chromedp.Text("._2MImiq span", &pagesString, chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}

	page, err := getTotalPages(pagesString)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		chromedp.Sleep(2 * time.Second).Do(ctx)

		fmt.Println("visiting page............", i+1, " of ", page)
		chromedp.Run(ctx, chromedp.Nodes("._13oc-S", &nodes, chromedp.ByQueryAll))

		findProduct(ctx, nodes, &products)

		err := chromedp.Run(ctx, chromedp.Click("._1LKTO3", chromedp.NodeVisible))
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("/////////////////////////", len(products))
	for i, product := range products {

		fmt.Printf("Product %d: %s , %s\n", i+1, product.id, product.price)
	}

}

func findProduct(ctx context.Context, nodes []*cdp.Node, products *[]Product) {

	var price string
	product := Product{}
	for _, node := range nodes {

		var divChildrens []*cdp.Node

		if err := chromedp.Run(ctx, chromedp.Nodes("div", &divChildrens, chromedp.ByQueryAll, chromedp.FromNode(node))); err != nil {
			log.Fatal(err)
		}

		for _, child := range divChildrens {
			dataId, found := child.Attribute("data-id")
			if found {
				err := chromedp.Run(ctx, chromedp.Text("._30jeq3", &price, chromedp.ByQuery, chromedp.FromNode(child)))
				if err != nil {
					fmt.Printf("Error getting price: %v\n", err)
					continue
				}

				product.id = dataId
				product.price = price
				*products = append(*products, product)
			}

		}

	}

}

func getTotalPages(inputString string) (int, error) {
	re := regexp.MustCompile(`Page (\d+) of (\d+)`)

	matches := re.FindStringSubmatch(inputString)

	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid format: %s", inputString)
	}

	totalPageNumber := matches[2]

	page, err := strconv.Atoi(totalPageNumber)
	if err != nil {
		return 0, err
	}
	return page, nil
}
