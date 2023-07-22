package filter

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func CoinsFilter(title string) bool {

	p1 := regexp.MustCompile("(?i)invest").MatchString(title)
	p2 := regexp.MustCompile("(?i)financing").MatchString(title)
	p3 := regexp.MustCompile("(?i)Google").MatchString(title)
	p4 := regexp.MustCompile("(?i)apple").MatchString(title)
	return p1 || p2 || p3 || p4

}
func CompaniesFilter(title string) bool {

	p1 := regexp.MustCompile("(?i)btc").MatchString(title)
	p2 := regexp.MustCompile("(?i)eth").MatchString(title)
	p3 := regexp.MustCompile("(?i)crypto").MatchString(title)

	return p1 || p2 || p3

}
func CzBinance(title string) bool {

	p1 := regexp.MustCompile("(?i)bnb").MatchString(title)

	return p1

}
func Elonmusk(title string) bool {

	p1 := regexp.MustCompile("(?i)doge").MatchString(title)

	return p1

}
func Upbit(title string) bool {

	p1 := regexp.MustCompile("(?i)\\[Trading|Transaction\\]").MatchString(title) &&
		regexp.MustCompile("(?i)Add").MatchString(title) &&
		regexp.MustCompile("(?i)KRW").MatchString(title)
	return p1

}
func Coinbase(title string) bool {

	p1 := regexp.MustCompile("(?i)Coinbase").MatchString(title)
	p2 := regexp.MustCompile("(?i)Roadmap").MatchString(title)

	return p1 && p2

}
func containsPhrase(text, phrase string) bool {
	text = strings.ToLower(text)
	phrase = strings.ToLower(phrase)
	words := strings.Split(phrase, " ")

	for _, word := range words {
		if !strings.Contains(text, word) {
			return false
		}
	}

	return true
}
func BinanceEN(title string) bool {
	//
	//p1 := regexp.MustCompile("(?i)Launchpad").MatchString(title) && regexp.MustCompile("(?i)Token Sale").MatchString(title)
	//p2 := regexp.MustCompile("(?i)Launchpool").MatchString(title) && regexp.MustCompile("(?i)Staking BNB").MatchString(title)
	//p3 := regexp.MustCompile("(?i)Launchpad").MatchString(title) && regexp.MustCompile("(?i)Now Open ").MatchString(title)
	//
	//p4 := regexp.MustCompile("(?i)Binance Will List").MatchString(title)
	//
	//return p1 || p2 || p3 || p4

	keywords := []string{"Launchpad", "Launchpool", "Will List Innovation", "20x Leverage",
		"Will Open Trading", "Listing Will", "Will Delist", "Removal of Trading",
		"Will Delist", "Will Support", "Wallet Maintenance", "Completes Integration"}

	for _, keyword := range keywords {
		if containsPhrase(title, keyword) {
			log.Info("Text contains keyword: ", keyword)
			return true
			break
		}
	}
	return false
}
