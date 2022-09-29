package helpers

import (
	"encoding/json"
	"fmt"
	"log"
)

func JSONEncode(data interface{}) string {
	jsonResult, _ := json.Marshal(data)

	return string(jsonResult)
}

func ErrorMessage(message string, err interface{}) {
	fmt.Println()
	log.Println("-------------------- Start Error Message --------------------")
	log.Println("Message => " + message + ".")

	if err != nil {
		log.Println("Error => ", err)
	}

	log.Println("-------------------- End Of Error Message --------------------")
	fmt.Println()
}
