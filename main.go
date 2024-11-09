package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
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
		MustLaunch()

	browser := rod.New().ControlURL(launcherURL).MustConnect()

	defer browser.Close()

	page := stealth.MustPage(browser)

	sp := page.MustNavigate(url)

	fmt.Println(sp.MustHTML())

	time.Sleep(time.Second * 30)

}
