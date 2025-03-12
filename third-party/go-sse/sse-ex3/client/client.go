package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	b, err := json.Marshal(struct{ Sender, Text string }{
		Sender: "Go client",
		Text:   fmt.Sprintf("The time is %s", time.Now().Format(time.Kitchen)),
	})
	check(err)

	buf := bytes.NewBuffer(b)
	resp, err := http.Post("http://localhost:8888/update", "application/json", buf)
	check(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	fmt.Println(string(body))
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
