package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	// take stdin input
	var input any
	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&input)
	if err != nil {
		log.Println("ERROR Failed to decode plugin input:", err)
		return
	}

	// process the input and produce output
	output := map[string]any{
		"status":  "success",
		"message": "Plugin executed successfully",
		"input":   input,
	}

	// output the result as JSON to stdout
	encoder := json.NewEncoder(os.Stdout)
	err = encoder.Encode(output)
	if err != nil {
		log.Println("ERROR Failed to encode plugin output:", err)
		return
	}
}
