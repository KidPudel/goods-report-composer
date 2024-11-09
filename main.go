package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

const (
	userAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"
	userData       = "/Users/iggysleepy/Library/Application Support/Google/Chrome"
	userProfile    = "Profile 1"
	executablePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	url            = "https://komus.ru"
)

const (
	searchBarQuery = "input.input__field.input-search__field.qa-search-field.js-field-input.js-search-input.ui-autocomplete-input"
)

const (
	title = "title"
	price = "price"
	unit  = "unit"
)

var (
	goodsQueries = map[string]string{
		title: "h1.product-details-page__title",
		price: "span.js-current-price",
		unit:  "span.product-price__current-price.product-price__current-price--unit-of-sale",
	}
)

func main() {
	inputReader := bufio.NewReader(os.Stdout)

	goodsNumbers := make([]string, 0)

	for {
		number, _, err := inputReader.ReadLine()
		if err != nil {
			err = fmt.Errorf("failed to read input %d, error: %w", len(goodsNumbers)+1, err)
			log.Fatal(err)
		}
		if string(number) == "-" {
			break
		}
		goodsNumbers = append(goodsNumbers, string(number))
	}

	goods, err := scrapeGoods(goodsNumbers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(goods)

}

type goodsInfo struct {
	title string
	price string
	unit  string
}

func scrapeGoods(goodsNumbers []string) ([]goodsInfo, error) {
	// simplify launching browser with executable path
	launcherURL := launcher.New().
		Headless(false).
		Bin(executablePath).
		Set("user-data-dir", userData).
		Set("profile-directory", userProfile).
		Set("disable-blink-features", "AutomationControlled").
		// Set("--no-sandbox", "true").
		// Set("--disable-extensions", "true").
		MustLaunch()

	browser := rod.New().ControlURL(launcherURL).MustConnect()

	defer browser.Close()

	page := stealth.MustPage(browser)

	sp := page.MustNavigate(url)

	time.Sleep(time.Second * 3)

	goods := make([]goodsInfo, len(goodsNumbers))
	for i, number := range goodsNumbers {
		searchForBarResult, err := sp.Search(searchBarQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to find search bar: %w", err)
		}

		searchBar := searchForBarResult.First

		err = searchBar.WaitWritable()
		if err != nil {
			return nil, fmt.Errorf("failed to wait: %w", err)
		}

		err = searchBar.Hover()
		if err != nil {
			return nil, fmt.Errorf("failed to hover: %w", err)
		}
		searchBar.MustClick()

		err = searchBar.Input(number)
		if err != nil {
			return nil, fmt.Errorf("failed to input into search bar: %w", err)
		}
		time.Sleep(time.Second)
		searchBar.MustKeyActions().Type(input.Enter).MustDo()
		time.Sleep(time.Second * 3)

		good := goodsInfo{}

		for elementName, query := range goodsQueries {
			domSearchResult, err := sp.Search(query)
			if err != nil {
				return nil, fmt.Errorf("failed to find %s in dom: %w", elementName, err)
			}
			fmt.Println(domSearchResult.First.Text())
			switch elementName {
			case title:
				good.title, err = domSearchResult.First.Text()
			case price:
				good.price, err = domSearchResult.First.Text()
			case unit:
				good.unit, err = domSearchResult.First.Text()
			}
			if err != nil {
				return nil, fmt.Errorf("failed to get text from %s: %w", elementName, err)
			}
		}

		goods[i] = good
	}
	return goods, nil
}
