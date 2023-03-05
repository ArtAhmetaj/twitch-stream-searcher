package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var (
	connectionsLimit = 10
)

func main() {
	twitchClient := NewTwitchClient()
	seleniumBrowserAutomater := NewSeleniumBrowserAutomater()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter text to search for on twitch: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	channels, err := twitchClient.getTwitchChannels(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	links := getLinks(channels)
	if len(links) == 0 {
		fmt.Println("No links available")
		return
	}
	candidateLink := links[0]
	links = links[1:]

	err = seleniumBrowserAutomater.StartSession(Chrome)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = seleniumBrowserAutomater.SelectAndOpenTabs(candidateLink, links, connectionsLimit)
	if err != nil {
		fmt.Println(err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			videoPlaying := seleniumBrowserAutomater.DoesVideoExistInPage()
			if !videoPlaying {
				if len(links) == 0 {
					err = seleniumBrowserAutomater.EndSession()
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("No more links available")
					return
				}
				candidateLink = links[0]
				links = links[1:]
				err = seleniumBrowserAutomater.CloseCurrentTab()
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}

}

func getLinks(channels []TwitchChannel) []string {

	var links []string
	for _, v := range channels {
		links = append(links, v.channelLink)
	}

	return links

}
