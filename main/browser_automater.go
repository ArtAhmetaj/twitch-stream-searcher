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
	seleniumPath = "./selenium-server-4.8.0.jar"
	port         = 8080
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
		//selenium.StartFrameBuffer(), // Start an X frame buffer for the browser to run in.
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

	return nil
}

func (s *SeleniumBrowserAutomater) SelectAndOpenTabs(mainCandidate string, candidateUrls []string) error {

	err := s.webDriver.Get(mainCandidate)
	if err != nil {
		return err
	}

	for _, u := range candidateUrls {
		err := s.webDriver.Get(u)
		if err != nil {
			return err
		}
	}
	err = s.webDriver.SwitchWindow(mainCandidate)
	if err != nil {
		return err
	}
	return nil
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
