package main

import (
	"fmt"
	"log"
)

func main() {
	channels, err := getTwitchChannels("football match")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(channels)
}
