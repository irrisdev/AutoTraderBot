package bot

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	_ "github.com/mattn/go-sqlite3"
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

	insertMakes(makes)
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
	iframe := page.MustWaitLoad().MustElement(`#sp_message_container_1086457 iframe`)
	if iframe == nil {
		fmt.Println("Accept cookies iframe not found")
	} else {
		fmt.Println("Accept cookies iframe found")

		// Switch to the iframe context
		temppage := iframe.MustFrame()

		// Click the reject button using its XPath
		reject := temppage.MustElementX(`/html/body/div/div[2]/div[4]/button[2]`)
		fmt.Println("Reject button found")

		// Wait for the reject button to be clickable
		reject.MustWaitVisible().MustWaitEnabled().MustClick()
		fmt.Println("Reject button clicked")

		temppage.MustWaitStable()

	}
	page.MustScreenshot("afterCookies.png")

	modelsBtn := page.MustWaitLoad().MustElement(`#content > article > div.at__sc-1okmyrd-3.caRFaz > section > section:nth-child(3) > button`)
	modelsBtn.MustClick()
	models := page.MustElement(`#model-facet-panel`).MustWaitStable()

	fmt.Println(models.HTML())
	page.MustScreenshot("after.png")

	fmt.Println("Complete")

	defer browser.MustClose()
}

func insertMakes(makes []string) {

	db, err := sql.Open("sqlite3", "database/autotraderDB.db")
	if err != nil {
		log.Err(err).Msg("Error while opening database")
	}
	defer db.Close()

	//Begin Transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal().Err(err).Msg("Error while beginning tx")
	}

	//Prepare SQL Statement
	stmt, err := tx.Prepare("insert into manufacturer(make) values(?)")
	if err != nil {
		log.Fatal().Err(err).Msg("Error while preparing statement")
	}
	defer stmt.Close()

	for _, m := range makes {

		_, err = stmt.Exec(m)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while executing statement")
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal().Err(err).Msg("Error while committing tx")
	}

}
