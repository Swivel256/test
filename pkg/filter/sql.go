package filter

import (
	"github.com/ginvmbot/aitrade/pkg/gojsonq"
	"os"
)

var Sql *gojsonq.JSONQ

type TwitterUser struct {
	UserId int    `json:"userid"`
	Cate   string `json:"cate"`
}

func init() {
	pwd, _ := os.Getwd()
	Sql = gojsonq.New().File(pwd + "/data/twitters.json")
}
