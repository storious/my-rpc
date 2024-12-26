package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Req struct {
	Parameter any `json:"parameter"`
}

func main() {
	req := Req{Parameter: "req"}
	body, _ := json.Marshal(req)
	resp, err := http.Post("url", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	resBody, _ := io.ReadAll(resp.Body)
	var id int
	err = json.Unmarshal(resBody, &id)
	if err != nil {
		return
	}
}
