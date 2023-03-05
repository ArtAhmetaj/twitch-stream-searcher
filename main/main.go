package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	twitchClient := NewTwitchClient()
	seleniumBrowserAutomater := NewSeleniumBrowserAutomater()
	defer seleniumBrowserAutomater.EndSession()
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

	err = seleniumBrowserAutomater.SelectAndOpenTabs(candidateLink, links, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Possible args: \nn=Switch to new channel when the channel breaks\nz=stop execution")
		switch char {
		case 'n':
			if len(links) == 0 {
				channels, err = twitchClient.getTwitchChannels(input)
				links = getLinks(channels)
			}
			if len(links) == 0 {
				fmt.Println("No links available")
				return
			}
			newCandidateLink := links[0]
			links = links[1:]
			err := seleniumBrowserAutomater.ReplaceTab(newCandidateLink)
			if err != nil {
				fmt.Println(err)
				return
			}
			break
		case 'z':
			os.Exit(0)
		}
	}
}

func getLinks(channels map[TwitchChannel]bool) []string {

	var links []string
	for k := range channels {
		links = append(links, k.channelLink)
	}

	return links

}
