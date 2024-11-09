package main

import (
	"fmt"
	"log"
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

func main() {
	// simplify launching browser with executable path
	launcherURL := launcher.New().
		Headless(false).
		Bin(executablePath).
		Set("user-data-dir", userData).
		Set("profile-directory", userProfile).
		Set("disable-blink-features", "AutomationControlled").
		Set("--no-sandbox", "true").
		Set("--disable-extensions", "true").
		MustLaunch()

	browser := rod.New().ControlURL(launcherURL).MustConnect()

	defer browser.Close()

	page := stealth.MustPage(browser)

	sp := page.MustNavigate(url)

	time.Sleep(time.Second * 5)

	searchForElementResult, err := sp.Search("input.input__field.input-search__field.qa-search-field.js-field-input.js-search-input.ui-autocomplete-input")
	if err != nil {
		log.Fatal("failed to get search bar: ", err.Error())
	}

	searchBar := searchForElementResult.First

	err = searchBar.WaitWritable()
	if err != nil {
		log.Fatal("failed to get wait: ", err.Error())
	}

	fmt.Println("end writable")

	err = searchBar.Hover()
	if err != nil {
		log.Fatal("failed to hover: ", err.Error())
	}
	searchBar.MustClick()

	err = searchBar.Input("1308488")
	time.Sleep(time.Second * 3)
	searchBar.MustKeyActions().Press(input.Enter)

	time.Sleep(time.Second * 30)
}
