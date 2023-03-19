package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type Lock struct {
	Resource string `json:"resource"`
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

	RESOURCE := "Alpha"
	LOCK_SECONDS := time.Second * 5

	json, err := json.Marshal(Lock{
		Resource: RESOURCE,
		UserID:   "admin",
		Duration: LOCK_SECONDS.String(),
	})
	if err != nil {
		fmt.Println(err)
	}

	TEST_KEY := makeKey(NAMESPACE, RESOURCE)

	fmt.Printf("TEST_KEY: %s \n", TEST_KEY)

	err = client.Set(TEST_KEY, json, LOCK_SECONDS).Err()
	if err != nil {
		fmt.Println(err)
	}

	// loop until the key expires
	for range time.Tick(time.Second * 1) {
		// Test for existance of the key
		exists, err := client.Exists(TEST_KEY).Result()
		fmt.Printf("\nEXISTS? key %s exists? %v ... err: %v \n", TEST_KEY, exists, err)
		if err != nil {
			// An unexpected error occured
			fmt.Printf("ERROR [Exists]: %v \n", err)
		}
		// Get the key
		val, err := client.Get(TEST_KEY).Result()
		if len(val) == 0 { // or val == ""
			// if an empty sting was retuned, the key was not found
			fmt.Println("--- key not found ---")
		}
		if err != nil {
			if err == redis.Nil {
				// If the error was redis.Nil, the key was not found
				fmt.Printf("--- GET returned redis.Nil, err: %v ---\n", err)
			} else {
				// Otherwise an unexpected error occurred
				fmt.Printf("ERROR [Get]: %v \n", err)
			}

			break
		}

		fmt.Println(val)
	}

}
