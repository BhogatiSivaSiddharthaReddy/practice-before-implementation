package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/user", user)
	http.ListenAndServe(":8080", nil)
}

type Request_Body struct {
	Name string `json:"name"`
}

type Response_Body struct {
	Received_name string `json:"received_name"`
	Received_id   string `json:"received_id"`
}

func user(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request_data Request_Body

	err := json.NewDecoder(req.Body).Decode(&request_data)

	if err != nil {
		http.Error(res, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	id := req.URL.Query().Get("id")

	response_data := Response_Body{
		request_data.Name,
		id,
	}

	res.Header().Set("Content-Type",
		"application/json")

	json.NewEncoder(res).Encode(response_data)
}
