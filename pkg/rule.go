package pkg

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"regexp"
	"time"
)

type RuleEngine struct {
	Rules  map[string]func(data map[string]interface{}) bool
	Rows   [][]interface{}
	Data   dataframe.DataFrame
	Titles []string
	Time   []string
	Source []string
	Url    []string
}

func  NewRule( ) RuleEngine {
	engine := RuleEngine{Rules: make(map[string]func(data map[string]interface{}) bool)}

	engine.AddRule("Upbit", Upbit1)
	engine.AddRule("Coinbase", Coinbase1)
	engine.AddRule("Binance_Bnb", Binance1)
	engine.AddRule("Binance_Alt", Binance2)
	return engine
}


func (e *RuleEngine) AddRule(name string, fn func(data map[string]interface{}) bool) {
	e.Rules[name] = fn
}

func (e *RuleEngine) RunRules(data map[string]interface{}) (bool, string) {

	for name, fn := range e.Rules {
		if fn(data) == true {

			e.Titles = append(e.Titles, data["en"].(string))
			e.Time = append(e.Time, data["time"].(string))
			e.Source = append(e.Source, data["source"].(string))
			e.Url = append(e.Url, data["url"].(string))
			return true, name
		}
	}

	return false, ""
}
func (e *RuleEngine) processNews(data gjson.Result) (bool, string) {
	res := map[string]interface{}{}
	res["title"] = data.Get("title").String()
	res["en"] = data.Get("en").String()
	res["source"] = data.Get("source").String()
	res["url"] = data.Get("url").String()
	t := data.Get("time").Int()
	unixTimeStamp := time.Unix(t/1000, 0)
	res["time"] = unixTimeStamp.Format("2006-01-02 15:04:05")
	res["suggestions"] = data.Get("suggestions").Array()
	if res["en"] == "" {
		res["en"] = res["title"]
	}

	return e.RunRules(res)
}
func (e *RuleEngine) SaveToCSV() {
	e.Data = dataframe.New(
		series.New(e.Titles, series.String, "title"),
		series.New(e.Time, series.String, "time"),
		series.New(e.Source, series.String, "source"),
		series.New(e.Url, series.String, "url"),
	)
	file, err := os.Create("news.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	 
	err = e.Data.WriteCSV(file)
	if err != nil {
		fmt.Println("Error writing CSV:", err)
	}
}

func Upbit1(data map[string]interface{}) bool {
	title := data["en"].(string)
	//fmt.Println(title)
	canTrade := regexp.MustCompile("(?i)\\[Trading|Transaction\\]").MatchString(title) &&
		regexp.MustCompile("(?i)Add").MatchString(title) &&
		regexp.MustCompile("(?i)KRW").MatchString(title)
	//if canTrade == true {
	//	fmt.Println(canTrade, title)
	//}
	return canTrade
}

func Coinbase1(data map[string]interface{}) bool {
	title := data["en"].(string)

	p1 := regexp.MustCompile("(?i)Coinbase").MatchString(title)
	p2 := regexp.MustCompile("(?i)Roadmap").MatchString(title)
	source := data["source"].(string) == "Coinbase"

	return source && p1 && p2

}

func Binance1(data map[string]interface{}) bool {
	title := data["en"].(string)
	p1 := regexp.MustCompile("(?i)Launchpad").MatchString(title) && regexp.MustCompile("(?i)Token Sale").MatchString(title)
	p2 := regexp.MustCompile("(?i)Launchpool").MatchString(title) && regexp.MustCompile("(?i)Staking BNB").MatchString(title)
	p3 := regexp.MustCompile("(?i)Launchpad").MatchString(title) && regexp.MustCompile("(?i)Now Open ").MatchString(title)
	source := data["source"].(string) == "Binance EN"
	return source && (p1 || p2 || p3)
}
func Binance2(data map[string]interface{}) bool {
	title := data["en"].(string)

	p1 := regexp.MustCompile("(?i)Binance Will List").MatchString(title)
	source := data["source"].(string) == "Binance EN"
	return source && p1
}

func test() {
	engine := RuleEngine{Rules: make(map[string]func(data map[string]interface{}) bool)}

	engine.AddRule("Upbit", Upbit1)
	engine.AddRule("Coinbase", Coinbase1)
	engine.AddRule("Binance_Bnb", Binance1)
	engine.AddRule("Binance_Alt", Binance2)

	// ... Load JSON data
	content, err := os.ReadFile("./allNews.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	json := string(content)
	value := gjson.Parse(json)

	for _, item := range value.Array() {
		engine.processNews(item)
	}

	engine.SaveToCSV()
}
