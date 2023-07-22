package news

import (
	"bytes"
	"encoding/json"
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	gtranslate "github.com/gilang-as/google-translate"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func DeepTranslate(text string) string {

	data := map[string]interface{}{

		"text":        []string{text},
		"target_lang": "ZH",
	}

	jsonStr, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key a27ae0bd-4f79-2fd6-6c84-038da52ecf9b:fx")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	value := gjson.Parse(string(body))
	res := value.Get("translations.0.text")

	return res.String()
}
func GoogleTranslate(text string) string {
	result, err := gt.Translate(text, "en", "zh")
	//fmt.Println(result)
	if err != nil {
		fmt.Println("google", err)
		return ""
	}
	return result
}

func Translate(text string) string {
	//fmt.Println("deeprrr", text)
	//title := DeepTranslate(text)
	//fmt.Println("deep", title)
	//if title != "" {
	//
	//	fmt.Println("google", title)
	//
	//	title = GoogleTranslate(title)
	//}

	value := gtranslate.Translate{
		Text: text,
		From: "en",
		To:   "zh",
	}
	translated, err := gtranslate.Translator(value)
	if err != nil {

	} else {
		return translated.Text
	}
	return text

	//title := GoogleTranslate(text)
	////fmt.Println(result)
	//return title

}
