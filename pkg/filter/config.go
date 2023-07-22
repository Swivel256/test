package filter

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"os"
)

var Data   dataframe.DataFrame
var Titles []string
var Time   []string
var Source []string
var Url    []string
var Twid    []string
func   SaveToCSV() {
	Data = dataframe.New(
		series.New(Titles, series.String, "title"),
		series.New(Time, series.String, "time"),
		series.New(Source, series.String, "source"),
		series.New(Url, series.String, "url"),
		series.New(Twid, series.String, "twid"),
	)
	file, err := os.Create("news.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	err =  Data.WriteCSV(file)
	if err != nil {
		fmt.Println("Error writing CSV:", err)
	}
}
