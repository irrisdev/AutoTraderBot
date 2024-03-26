package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ScrapeMakes() {

	var makes []string

	res, err := http.Get("https://www.autotrader.co.uk/")
	if err != nil {
		log.Err(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal().Msg(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Err(err)
	}

	makeSelection := doc.Find("#make")

	optgroup := makeSelection.Find("optgroup[label='All makes']")

	optgroup.Find("option").Each(func(i int, selection *goquery.Selection) {
		value, _ := selection.Attr("value")

		makes = append(makes, value)
	})

	var result string
	for _, value := range makes {
		result += fmt.Sprintf("\"%s\": struct{}{},\n", value)
	}
	fmt.Println(result)
}
