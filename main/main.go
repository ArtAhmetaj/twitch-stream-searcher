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

	err = seleniumBrowserAutomater.SelectAndOpenTabs(links[0], links[1:])
	if err != nil {
		return
	}

	for {
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		switch char {
		case 'n': // select new tab
		}
	}
}
