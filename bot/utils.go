package bot

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
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

func ScrapeModels(session *UserSession, complete chan<- bool) {

	log.Info().Msg("Begging Scrape")
	newReq := session.RequestDetails
	url := createURL(&newReq)
	log.Info().Msg(url)

	complete <- true
}

func ScrapeModelTest() {

	url := "https://www.autotrader.co.uk/car-search?postcode=BT490FQ&make=BMW&model="

	browser := rod.New().Timeout(time.Minute).MustConnect()

	page := stealth.MustPage(browser)

	page.MustNavigate(url).MustWindowFullscreen()

	page.MustWaitStable()
	iframe := page.MustElement(`#sp_message_container_1086457 iframe`)
	if iframe == nil {
		fmt.Println("Accept cookies iframe not found")

		page.MustWaitLoad().MustElement(`[data-testid="toggle-facet-button"]`).MustClick()
	} else {
		fmt.Println("Accept cookies iframe found")

		// Switch to the iframe context
		page = iframe.MustFrame()

		// Click the reject button using its XPath
		reject := page.MustElementX(`/html/body/div/div[2]/div[4]/button[2]`)
		fmt.Println("Reject button found")

		// Wait for the reject button to be clickable
		reject.MustWaitVisible().MustWaitEnabled().MustClick()
		fmt.Println("Reject button clicked")

		page.MustWaitStable()

	}

	models := page.MustWaitLoad().MustElement(`section.at__sc-n1gtx5-0:nth-child(3)`)

	fmt.Println(models.HTML())

	page.MustScreenshot("after.png")

	defer browser.MustClose()
}
