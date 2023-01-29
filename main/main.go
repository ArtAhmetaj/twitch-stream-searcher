package main

import (
	"fmt"
)

func main() {

	channels, err := getTwitchChannels("big brother")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(channels)
}
