package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type User struct {
	Location string `json:"location"`
	UserID   string `json:"userId"`
	Duration string `json:"duration"`
}

func makeKey(namespace string, location string) string {
	return fmt.Sprintf(
		"%s:%s",
		strings.ToLower(namespace),
		strings.ToLower(location),
	)
}

func main() {
	fmt.Println("Testing Golang Redis")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	NAMESPACE := "lock"

	LOCATION := "Alpha"
	LOCK_SECONDS := time.Second * 5

	json, err := json.Marshal(User{
		Location: LOCATION,
		UserID:   "admin",
		Duration: LOCK_SECONDS.String(),
	})
	if err != nil {
		fmt.Println(err)
	}

	TEST_KEY := makeKey(NAMESPACE, LOCATION)

	fmt.Printf("TEST_KEY: %s \n", TEST_KEY)

	err = client.Set(TEST_KEY, json, LOCK_SECONDS).Err()
	if err != nil {
		fmt.Println(err)
	}

	for range time.Tick(time.Second * 1) {
		val, err := client.Get(TEST_KEY).Result()
		if len(val) == 0 {
			fmt.Println("--- key not found ---")
		}
		if err != nil {
			fmt.Printf("ERROR: %v \n", err)
			panic(err)
		}

		fmt.Println(val)
	}

}
