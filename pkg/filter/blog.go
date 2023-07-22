package filter

import (
	"regexp"
)

func BlogFilter(title string) bool {
	p1, _ := regexp.MatchString(`(?i).*(eth.*etf)|(btc.*etf)|(sec.*etf)|(Bitcoin.*etf).*`, title)

	p2, _ := regexp.MatchString("(?i).*sec.*bnb.*", title)
	//p2 := regexp.MustCompile("(?i)sec").MatchString(title) && regexp.MustCompile("(?i)bnb").MatchString(title)
	p3 := regexp.MustCompile("(?i)binance").MatchString(title)
	return p1 || p2 || p3

}
