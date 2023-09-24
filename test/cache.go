package test

import (
	"log"
)

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
