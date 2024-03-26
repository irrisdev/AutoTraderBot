package bot

import (
	"fmt"
	"net/url"
)

var Postcodes = []string{
	"M110AA",
	"BH37AY",
	"HP49DS",
	"TA65FE",
	"M67GU",
	"RH123ED",
	"M201LU",
	"BT490FQ",
	"DE30TD",
	"CB215HQ",
	"BH118RJ",
}

type CarRequest struct {
	Postcode       string
	Make           string
	Model          string
	AggregatedTrim string
	YearFrom       string
	YearTo         string
}

func (r *CarRequest) getCarParams() string {

	params := url.Values{}
	params.Set("postcode", r.Postcode)
	params.Set("make", r.Make)
	params.Set("model", r.Model)
	params.Set("aggregatedTrim", r.AggregatedTrim)
	params.Set("year-from", r.YearFrom)
	params.Set("year-to", r.YearTo)

	return params.Encode()

}

func createURL(req *CarRequest) string {

	urlParams := req.getCarParams()
	urlString := fmt.Sprintf("https://www.autotrader.co.uk/car-search?%s", urlParams)

	return urlString

}
