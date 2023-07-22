package main

import (
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/gojsonq"
)

func main() {
	//const json = `{"city":"dhaka","type":"weekly","temperatures":[30,39.9,35.4,33.5,31.6,33.2,30.7]}`
	//avg := gojsonq.New().FromString(json).From("temperatures").Avg()
	//fmt.Printf("Average temperature: %.2f", avg) // 33.471428571428575

	res := gojsonq.New().File("./aaa.json").From("items").Select("userid", "cate").Where("userid", "=", "180505807").Get()
	fmt.Println(res)
}
