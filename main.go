package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	config := InitConfig()

	memories, err := config.Session.GetMemories()
	if err != nil {
		fmt.Println("Fetching memories")
		panic(err)
	}
	res, err := json.Marshal(memories.Data[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}
