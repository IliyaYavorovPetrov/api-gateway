package test

import (
	"encoding/json"
	"fmt"
	"log"
)

type ItemTest struct {
	ID      int     `redis:"id" json:"id"`
	Content string  `redis:"content" json:"content"`
	Size    float64 `redis:"size" json:"size"`
}

func ToString(item ItemTest) string {
	data, err := json.Marshal(item)
	if err != nil {
		fmt.Println("error encoding to json:", err)
		return ""
	}

	return string(data)
}

func ClearLocalCache() {
	err := loc.Flush(ctx)
	if err != nil {
		log.Fatal("failed to clear local storage")
	}
}

func ClearDistributedCache() {
	err := dist.Flush(ctx)
	if err != nil {
		log.Fatal("failed to clear local storage")
	}
}

func ClearCaches() {
	ClearLocalCache()
	ClearDistributedCache()
}
