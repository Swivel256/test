package filter

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func GetTwitter(id string) string {

	v := []TwitterUser{}
	Sql.From("items").Select("userid", "cate").Where("userid", "=", id).Out(&v)
	if len(v) > 0 {
		return v[0].Cate
	}
	return ""
}
func TwitterFilter(json []byte) bool {
	data := gjson.ParseBytes(json)
	title := data.Get("title").String()
	if data.Get("en").Exists() {
		title = data.Get("en").String()
	}
	if data.Get("body").Exists() {
		title = fmt.Sprintf("%s: %s", title, data.Get("body").String())
	}

	twid := data.Get("info.twitterId").String()

	var maySend bool
	//@cz_binance
	if twid == "902926941413453824" {
		maySend = CzBinance(title)
		return maySend
	}
	//elonmusk

	if twid == "44196397" {
		fmt.Println("elonmusk")
		maySend = Elonmusk(title)
		return maySend
	}
	//s:=[]int64{1391538435261861894}
	/**
	tier10k 2361601055
	Lookonchain	1462727797135216641

	zoomer 1391538435261861894
	snailnews 1634057220206874624
	Tree News 1282727055604486148
	*/

	if  twid == "2361601055" ||
		twid == "1462727797135216641" ||
		twid == "1391538435261861894" ||
		twid == "1634057220206874624" ||
		twid == "1282727055604486148" {

		return true
	}

	cate := GetTwitter(twid)

	switch cate {
	case "news-outlets-twitter":
	case "news-relays":
		maySend = BlogFilter(title)
		break
	case "coins":
	case "spot-coins":
		maySend = CoinsFilter(title)
		break
	case "companies":
		maySend = CompaniesFilter(title)
		break

	case "wublockchain":
		maySend = true
		break

	case "coinbase-listings":

		maySend = true
		break
	case "news-of-alpha":
		maySend = true
		break
	default:
		maySend = false
		break
	}
	return maySend

}
