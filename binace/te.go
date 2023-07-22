package main

import (
	"strings"
)

func containsPhrase(text, phrase string) bool {
	words := strings.Split(phrase, " ")

	for _, word := range words {
		if !strings.Contains(text, word) {
			return false
		}
	}

	return true
}
func main() {

	text := `这里是一段示例文字 Will Listing Trading 内容`

	keywords := []string{"Launchpad", "Launchpool", "Will List Innovation", "20x Leverage",
		"Will Open Trading", "Listing Will", "Will Delist", "Removal of Trading",
		"Will Delist", "Will Support", "Wallet Maintenance", "Completes Integration"}

	for _, keyword := range keywords {
		if containsPhrase(text, keyword) {
			println("Text contains keyword: ", keyword)
			break
		}
	}

}
