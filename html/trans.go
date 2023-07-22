package main

import (
	"encoding/json"
	"fmt"
	gtranslate "github.com/gilang-as/google-translate"
)

func main() {
	value := gtranslate.Translate{
		Text: "DODO: The DODO Trading Widget: DeFi Wherever You Need It",
		From: "en",
		To:   "zh",
	}
	translated, err := gtranslate.Translator(value)
	if err != nil {
		panic(err)
	} else {
		prettyJSON, err := json.MarshalIndent(translated, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Println(translated.Text, prettyJSON)
	}
}
