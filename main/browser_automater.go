package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"os"
)

type Browser int

const (
	Chrome Browser = iota
	Firefox
)

func (b Browser) String() string {
	return []string{"chrome", "firefox"}[b]
}

const (
	seleniumPath    = "selenium-server-standalone-3.5.3.jar"
	port            = 8080
	geckoDriverPath = "geckodriver.exe"
)

type BrowserAutomater interface {
	StartSession(browser Browser) error
	SelectAndOpenTabs(mainUrl string, candidateUrls []string) error
	ReplaceTab(url string) error
	EndSession() error
}

type SeleniumBrowserAutomater struct {
	webDriver selenium.WebDriver
}

func NewSeleniumBrowserAutomater() *SeleniumBrowserAutomater {
	return &SeleniumBrowserAutomater{}
}

func startSeleniumService() error {
	//TODO: add firefox integration through gecko driver
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver("chromedriver.exe"),
		//selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	_, err := selenium.NewSeleniumService(seleniumPath, port, opts...)

	return err
}

func (s *SeleniumBrowserAutomater) StartSession(browser Browser) error {
	capabilities := selenium.Capabilities{"browserName": browser.String()}
	err := startSeleniumService()
	if err != nil {
		return err
	}

	wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	s.webDriver = wd
	return err
}

func (s *SeleniumBrowserAutomater) SelectAndOpenTabs(mainCandidate string, candidateUrls []string, limit int) error {

	err := s.webDriver.Get(mainCandidate)
	if err != nil {
		return err
	}
	handle, err := s.webDriver.CurrentWindowHandle()
	if err != nil {
		return err
	}

	for i, u := range candidateUrls {
		_, err := s.webDriver.ExecuteScript(fmt.Sprintf("window.open('%q','_blank');", u), nil)
		if err != nil {
			return err
		}
		if i >= limit {
			break
		}
	}
	err = s.webDriver.SwitchWindow(handle)
	return err
}

func (s *SeleniumBrowserAutomater) ReplaceTab(url string) error {
	err := s.webDriver.Close()
	if err != nil {
		return err
	}
	err = s.webDriver.SwitchWindow(url)
	if err != nil {
		return err
	}
	return nil
}

func (s *SeleniumBrowserAutomater) EndSession() error {
	err := s.webDriver.Close()
	if err != nil {
		return err
	}
	return nil
}
