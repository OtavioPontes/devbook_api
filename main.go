package main

import (
	"devbook_api/src/config"
	"devbook_api/src/router"
	"fmt"
	"log"
	"net/http"
)

// func init() {
// 	key := make([]byte, 64)

// 	if _, err := rand.Read(key); err != nil {
// 		log.Fatal(err)
// 	}

// 	stringBase64 := base64.StdEncoding.EncodeToString(key)
// 	fmt.Println(stringBase64)
// }

func main() {
	config.Load()

	r := router.Generate()

	fmt.Printf("Listening to port %d", config.Port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
