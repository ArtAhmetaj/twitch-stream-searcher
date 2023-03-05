package main

import (
	"fmt"
	"github.com/tebeka/selenium"
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
	DoesVideoExistInPage() bool
	CloseCurrentTab() error
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
	}
	_, err := selenium.NewSeleniumService(seleniumPath, port, opts...)

	return err
}

func (s *SeleniumBrowserAutomater) openNewTab(url string) error {
	_, err := s.webDriver.ExecuteScript(fmt.Sprintf("window.open(%q,'_blank');", url), nil)
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
		err = s.openNewTab(u)
		if err != nil {
			return err
		}
		if i+1 >= limit {
			break
		}
	}
	err = s.webDriver.SwitchWindow(handle)
	if err != nil {
		return err
	}
	err = s.webDriver.MaximizeWindow(handle)
	return err
}

func (s *SeleniumBrowserAutomater) CloseCurrentTab() error {
	err := s.webDriver.Close()
	if err != nil {
		return err
	}
	handles, err := s.webDriver.WindowHandles()
	if err != nil {
		return err
	}
	return s.webDriver.SwitchWindow(handles[0])
}

func (s *SeleniumBrowserAutomater) DoesVideoExistInPage() bool {
	_, err := s.webDriver.FindElement(selenium.ByTagName, "video")
	fmt.Println(err)
	return err == nil
}

func (s *SeleniumBrowserAutomater) EndSession() error {
	err := s.webDriver.Close()
	if err != nil {
		return err
	}
	return nil
}
