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
	SelectAndOpenTabs(urls []string) error
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

func (s *SeleniumBrowserAutomater) SelectAndOpenTabs(urls []string) error {

	for _, u := range urls {
		err := s.webDriver.Get(u)
		if err != nil {
			return err
		}
	}
	err := s.webDriver.SwitchWindow(urls[0])
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
