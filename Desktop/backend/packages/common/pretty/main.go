package pretty

import (
	"encoding/json"
	"fmt"
	"log"
)

// helper to pretty print interface as a JSON object
func Print(thing interface{}) {
	thingJSONBytes, err := json.MarshalIndent(thing, "", "  ")
	if err != nil {
		log.Println("Couldn't pretty print: " + err.Error())
		return
	}

	log.Println(string(thingJSONBytes))
}

func Sprint(thing interface{}) string {
	thingJSONBytes, err := json.MarshalIndent(thing, "", "  ")
	if err != nil {
		log.Println("Couldn't pretty print: " + err.Error())
		return ""
	}

	return fmt.Sprint(string(thingJSONBytes))
}
