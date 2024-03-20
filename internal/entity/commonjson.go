package entity

type MainMessageJSON struct {
	Type           int         `json:"type" `
	Number         int         `json:"number"`
	Status         int         `json:"status"`
	Charge         int         `json:"charge"`
	Temperature_MK int         `json:"temperature_MK"`
	Data           DataTypeMsg `json:"data"`
}

type MessageTypeOneJSON struct {
	MainMessageJSON
	Controlerleack ControlerLeackDataMessange `json:"controlerleack"`
}

type MessageTypeTwoJSON struct {
	MainMessageJSON
	Controlermodule ControlerModuleDataMessange `json:"controlermodule"`
}

type MessageTypeThreeJSON struct {
	MainMessageJSON
	Controlerenviroment ControlerEnviromentDataMessange `json:"controlerenviroment"`
}

type MessangeTypeZiroJson struct {
	One   MessageTypeOneJSON   `json:"one,omitempty"`
	Two   MessageTypeTwoJSON   `json:"two,omitempty"`
	Three MessageTypeThreeJSON `json:"three,omitempty"`
}

type DataTypeMsg struct {
	Second int `json:"second"`
	Minute int `json:"minute"`
	Hour   int `json:"hour"`
	Day    int `json:"day"`
	Month  int `json:"month"`
	Year   int `json:"year"`
}

type ControlerLeackDataMessange struct {
	Leack int `json:"leack"`
}

type ControlerModuleDataMessange struct {
	Temperature int `json:"temperature"`
	Humidity    int `json:"humidity"`
	Pressure    int `json:"pressure"`
	Gas         int `json:"gas"`
}

type ControlerEnviromentDataMessange struct {
	Temperature int `json:"temperature"`
	Humidity    int `json:"humidity"`
	Pressure    int `json:"pressure"`
	Voc         int `json:"voc"`
	Gas1        int `json:"gas1"`
	Gas2        int `json:"gas2"`
	Gas3        int `json:"gas3"`
	Pm1         int `json:"pm1"`
	Pm25        int `json:"pm25"`
	Pm10        int `json:"pm10"`
	Fire        int `json:"fire"`
	Smoke       int `json:"smoke"`
}
