package main

import (
	"fmt"
	"github.com/Aleksandr-Rozhok/Blog_Aggregator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(config)
	config.SetUser("Alex")
	fmt.Println(config)
}
