package main

import "marluxgithub/muehle/pkg/muehle/interfaces"

func main() {
	client := interfaces.NewClient()
	client.Start()

}
