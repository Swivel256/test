package config

type TickerMessage struct {
	Data Ticker `json:"data"`
}
type CanOrder struct {
	Symbol     string
	Exchange   string
	RuelName   string
	PlaceOrder bool
	Info      map[string]string
}
type Ticker struct {
	E            string `json:"e"`
	E1           int    `json:"E"`
	Symbol       string `json:"s"`
	P            string `json:"p"`
	P1           string `json:"P"`
	W            string `json:"w"`
	CurrentPrice string `json:"c"`
	Quantity     string `json:"Q"`
	O            string `json:"o"`
	H            string `json:"h"`
	L            string `json:"l"`
	V            string `json:"v"`
	Q1           string `json:"q"`
	O1           int    `json:"O"`
	C1           int    `json:"C"`
	F            int    `json:"F"`
	L1           int    `json:"L"`
	N            int    `json:"n"`
}

type Symbol struct {
	Symbol string `json:"symbol"`
}
type KlineMessage struct {
	Data Kline `json:"data"`
}
type Kline struct {
	E      string `json:"e"`
	E1     int    `json:"E"`
	Symbol string `json:"s"`
	Kline  struct {
		T               int    `json:"t"`
		T1              int    `json:"T"`
		S               string `json:"s"`
		I               string `json:"i"`
		F               int    `json:"f"`
		L               int    `json:"L"`
		Open            string `json:"o"`
		Close           string `json:"c"`
		High            string `json:"h"`
		Low             string `json:"l"`
		Volume          string `json:"v"`
		Number          int    `json:"n"`
		X               bool   `json:"x"`
		Quantity        string `json:"q"`
		TakeBuyVolume   string `json:"V"`
		TakeBuyQuantity string `json:"Q"`
		B               string `json:"B"`
	} `json:"k"`
}
