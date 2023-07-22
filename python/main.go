package main

import "github.com/ginvmbot/aitrade/python/client"

func main() {
	err := client.NewClient("Examples").Start()

	panic(err)
}
