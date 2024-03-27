package bot

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func ScrapeMakes() []string {

	var makes []string

	res, err := http.Get("https://www.autotrader.co.uk/")
	if err != nil {
		log.Err(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Err(err)
		}
	}(res.Body)
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

	return makes

	//insertMakes(makes)
}

func ScrapeModels(session *UserSession, progress chan<- string, result chan<- []string) {
	manufacture := session.RequestDetails.Make

	progress <- "Getting models..."
	log.Debug().Msg("Checking If Models Exists")
	models, exists := getModels(manufacture)
	if exists {
		log.Debug().Msg("Models Found Returning Models")
		close(progress)
		result <- models
		return
	}

	progress <- "No Existing models... Scraping off AutoTrader"
	log.Info().Msg("Begging Scrape of model")
	newReq := session.RequestDetails
	urlS := createURL(&newReq)
	log.Info().Msg(urlS)

	go ScrapeModelURL(urlS, progress, result, manufacture)
}

func getModels(manufacture string) ([]string, bool) {
	var modelsSlice []string

	db := openDb()
	defer db.Close()

	query := "SELECT COUNT(*) FROM models WHERE make = ?"

	var count int
	err = db.QueryRow(query, manufacture).Scan(&count)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Debug().Msg(fmt.Sprintf("Models found for %s : %d", manufacture, count))
	if count < 1 {
		log.Debug().Msg("Returning not found")
		return make([]string, 0), false
	}
	query = "select model from models where make = ?"

	rows, err := db.Query(query, manufacture)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rowVal string
		err = rows.Scan(&rowVal)
		if err != nil {
			log.Fatal().Err(err)
		}
		modelsSlice = append(modelsSlice, rowVal)
	}

	return modelsSlice, true
}

func ScrapeModelURL(url string, progress chan<- string, result chan<- []string, manufacture string) {
	log.Debug().Msg("Scraping Models")

	browser := rod.New().Timeout(time.Minute).MustConnect()

	page := stealth.MustPage(browser)

	page.MustNavigate(url).MustWindowFullscreen()

	page.MustWaitStable()
	iframe := page.MustWaitLoad().MustElement(`#sp_message_container_1086457 iframe`)
	if iframe == nil {
		log.Debug().Msg("Accept cookies iframe not found")
	} else {
		log.Debug().Msg("Accept cookies iframe found")

		// Switch to the iframe context
		temppage := iframe.MustFrame()

		// Click the reject button using its XPath
		reject := temppage.MustElementX(`/html/body/div/div[2]/div[4]/button[2]`)
		log.Debug().Msg("Reject button found")

		// Wait for the reject button to be clickable
		reject.MustWaitVisible().MustWaitEnabled().MustClick()
		log.Debug().Msg("Reject button clicked")

		temppage.MustWaitStable()

	}

	modelsBtn := page.MustWaitLoad().MustElement(`#content > article > div.at__sc-1okmyrd-3.caRFaz > section > section:nth-child(3) > button`)
	modelsBtn.MustClick()
	models := page.MustElement(`#model-facet-panel`).MustWaitStable()
	html, err := models.HTML()
	if err != nil {
		log.Err(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Err(err)
	}

	var modelsSlice []string

	doc.Find(".at__sc-1n64n0d-9.at__sc-qzn93z-3.bWawNd.eXNMEw").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		modelsSlice = append(modelsSlice, text)
	})

	insertModels(modelsSlice, manufacture)

	log.Debug().Msg("Completed Models Scrape")

	defer browser.MustClose()
	close(progress)
	result <- modelsSlice
}

func insertModels(sliceModel []string, vmake string) {
	log.Debug().Msg("Inserting Scraped Models into DB")
	db := openDb()
	tx := startTx(db)
	defer db.Close()

	stmt, err := tx.Prepare("insert into models values(?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msg("Error while preparing statement")
	}
	defer stmt.Close()

	for _, m := range sliceModel {

		_, err = stmt.Exec(m, vmake)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while executing statement")
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal().Err(err).Msg("Error while committing tx")
	}

}

func insertMakes(makes []string) {

	db := openDb()
	tx := startTx(db)
	defer db.Close()

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

// Opens Database connection
func openDb() *sql.DB {

	db, err := sql.Open("sqlite3", "database/autotraderDB.db")
	if err != nil {
		log.Err(err).Msg("Error while opening database")
	}

	return db
}

// Takes db as parameter and starts a new db transaction
func startTx(db *sql.DB) *sql.Tx {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal().Err(err).Msg("Error while beginning tx")
	}
	return tx
}

func InsertAll(carMakes []string, wg *sync.WaitGroup) {
	pst := Postcodes[rand.IntN(len(Postcodes))]

	browser := rod.New().Timeout(time.Minute * 5).MustConnect()

	page := stealth.MustPage(browser).Timeout(time.Minute * 5)

	page.MustNavigate("https://www.autotrader.co.uk/").MustWindowFullscreen()

	page.MustWaitStable()
	iframe := page.MustWaitLoad().MustElement(`#sp_message_container_1086457 iframe`)
	if iframe == nil {
		log.Debug().Msg("Accept cookies iframe not found")
	} else {
		log.Debug().Msg("Accept cookies iframe found")

		// Switch to the iframe context
		temppage := iframe.MustFrame()

		// Click the reject button using its XPath
		reject := temppage.MustElementX(`/html/body/div/div[2]/div[4]/button[2]`)
		log.Debug().Msg("Reject button found")

		// Wait for the reject button to be clickable
		reject.MustWaitVisible().MustWaitEnabled().MustClick()
		log.Debug().Msg("Reject button clicked")

		temppage.MustWaitStable()

	}

	for _, m := range carMakes {
		log.Info().Msg(fmt.Sprintf("Getting %s", m))
		params := url.Values{}
		params.Set("postcode", pst)
		params.Set("make", m)
		urlS := fmt.Sprintf("https://www.autotrader.co.uk/car-search?%s", params.Encode())
		log.Info().Msg(urlS)
		page.Timeout(time.Minute * 5)
		page.MustNavigate(urlS)
		page.MustWaitStable()

		modelsBtn := page.MustElement(`#content > article > div.at__sc-1okmyrd-3.caRFaz > section > section:nth-child(3) > button`)

		//fmt.Println(modelsBtn.HTML())
		if modelsBtn == nil || isDisabled(modelsBtn) {
			continue
		}
		modelsBtn.MustWaitStable().MustClick()

		models := page.MustElement(`#model-facet-panel`).MustWaitStable()
		html, err := models.HTML()
		if err != nil {
			log.Err(err)
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Err(err)
		}

		var modelsSlice []string
		db := openDb()
		defer db.Close()
		query := "SELECT EXISTS(SELECT 1 FROM models WHERE model = '%s' AND make = '%s')"

		doc.Find(".at__sc-1n64n0d-9.at__sc-qzn93z-3.bWawNd.eXNMEw").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()

			sqlStatement := fmt.Sprintf(query, text, m)
			fmt.Println("SQL statement:", sqlStatement)

			var exists int
			err = db.QueryRow(sqlStatement).Scan(&exists)
			if err != nil {
				log.Err(err)
			}
			fmt.Println(exists)

			if exists == 1 {
				log.Info().Msg("Already exists... Skipping")
				return
			}

			modelsSlice = append(modelsSlice, text)
		})

		log.Info().Msg(fmt.Sprintf("Results of %s %v", m, modelsSlice))

		fmt.Println(modelsSlice)
		insertModels(modelsSlice, m)

	}

	defer browser.MustClose()
	wg.Done()
}

func isDisabled(element *rod.Element) bool {
	attr, _ := element.Attribute("disabled")

	newString := fmt.Sprintf("%V", attr)

	if !strings.Contains(newString, "nil") {
		log.Info().Msg("button disabled")
		return true
	}
	return false
}
