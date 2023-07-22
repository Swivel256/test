package pkg

import (
	hashmap "github.com/duke-git/lancet/v2/datastructure/hashmap"
	"github.com/ginvmbot/aitrade/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"os"
)

type RuelEngine struct {
	Keys []string
	RuleEngine
	CanOrder chan *config.CanOrder
	Hashmap  *hashmap.HashMap
}

func NewRuelEngine(co chan *config.CanOrder) RuelEngine {
	content, err := os.ReadFile("data/key.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	json := string(content)
	value := gjson.Parse(json)
	engine := RuelEngine{
		CanOrder: co,
	}
	for _, item := range value.Get("keys").Array() {
		engine.Keys = append(engine.Keys, item.String())
	}

	return engine

}
func (e RuelEngine) ParseNews(json []byte) {
	e.OrdreNews(json)
	TeleNews(json)

}
func (e RuelEngine) OrdreNews(json []byte) {
	canTrade, name := e.processNews(gjson.ParseBytes(json))
	value := gjson.ParseBytes(json)
	data := value.Get("suggestions.0.symbols")
	//hm := hashmap.NewHashMapWithCapacity(uint64(100), uint64(1000))
	hm := make(map[string]string)

	for _, v := range data.Array() {
		//fmt.Println(v.Get("symbol").String(), v.Get("exchange").String())
		hm[v.Get("exchange").String()] = v.Get("symbol").String()
	}
	e.CanOrder <- &config.CanOrder{

		PlaceOrder: canTrade,
		RuelName:   name,
		Info:       hm,
	}

}
