package main

import (
	"api-go/internal/entity"
	"encoding/json"
	"fmt"
)

func main() {
	// app.Run()
	var tmp entity.MessangeTypeZiroJson
	var leack entity.MessageTypeOneJSON
	leack.MainMessageJSON.Charge = 1
	leack.MainMessageJSON.Status = 1
	leack.MainMessageJSON.Temperature_MK = 1
	leack.MainMessageJSON.Type = 1
	// leack.MainMessageJSON.Data
	leack.Controlerleack.Leack = 10
	leack.MainMessageJSON.Number = 1
	tmp.One = leack
	t, _ := json.Marshal(&tmp)
	fmt.Printf("%s", t)
	fmt.Println()
}
