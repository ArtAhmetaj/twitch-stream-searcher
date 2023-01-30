package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	twitchClient := NewTwitchClient()
	seleniumBrowserAutomater := NewSeleniumBrowserAutomater()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter text to search for on twitch: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	channels, err := twitchClient.getTwitchChannels(input)
	if err != nil {
		return
	}

	err = seleniumBrowserAutomater.StartSession(Chrome)
	if err != nil {
		return
	}

	var links []string
	for k := range channels {
		links = append(links, k.channelLink)
	}
	candidateLink := links[0]
	links = links[1:]
	err = seleniumBrowserAutomater.SelectAndOpenTabs(candidateLink, links)
	if err != nil {
		return
	}

	for {
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		switch char {
		case 'n':
			if len(links) == 0 {
				//TODO: getChannels again and then refresh link list
				fmt.Println("Could not find any more channels")
				os.Exit(1)
			}
			newCandidateLink := links[0]
			links = links[1:]
			err := seleniumBrowserAutomater.ReplaceTab(newCandidateLink)
			if err != nil {
				return
			}

		case 'z':
			os.Exit(0)
		}
	}
}
