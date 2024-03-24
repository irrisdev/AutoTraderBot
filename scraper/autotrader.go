package scraper

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
)

const Postcode string = "M110AA"

type carRequest struct {
	Postcode       string
	Make           string
	Model          string
	AggregatedTrim string
	YearFrom       string
	YearTo         string
}

func (r *carRequest) getCarParams() string {

	params := url.Values{}
	params.Set("postcode", r.Postcode)
	params.Set("make", r.Make)
	params.Set("model", r.Model)
	params.Set("aggregatedTrim", r.AggregatedTrim)
	params.Set("year-from", r.YearFrom)
	params.Set("year-to", r.YearTo)

	return params.Encode()

}

func Scrape() {

	req := carRequest{
		Postcode:       Postcode,
		Make:           "Volkswagen",
		Model:          "Polo",
		AggregatedTrim: "",
		YearFrom:       "2010",
		YearTo:         "2020",
	}

	urlParams := req.getCarParams()
	urlString := fmt.Sprintf("https://www.autotrader.co.uk/car-search?%s", urlParams)
	urlParsed, err := url.Parse(urlString)

	if err != nil {
		log.Err(err).Msg("Failed to parse target URL")
		return
	}

	fmt.Println(urlParsed)

}
