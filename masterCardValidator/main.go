package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type CardRequest struct {
	CardNumber string `json:"card_number"`
}

func (c CardRequest) Validate() bool {
	if len(c.CardNumber) != 16 {
		return false
	} else {
		var sum = 0
		for i := len(c.CardNumber) - 1; i >= 0; i-- {
			digit, _ := strconv.Atoi(string(c.CardNumber[i]))
			if i%2 == 0 {
				digit *= 2
				if digit > 9 {
					digit -= 9
				}
			}
			sum += digit
		}
		return sum%10 == 0
	}
}

func main() {
	port := flag.Int("port", 3000, "server port")
	flag.Parse()

	r := http.NewServeMux()

	r.HandleFunc("/master-card", cardHandler)
	fmt.Printf("Server is running on: http://localhost:%d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

func cardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var cardReq CardRequest
	err = json.Unmarshal(body, &cardReq)
	if err != nil {
		http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
		return
	}

	if cardReq.Validate() {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Valid Card\n")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid Card\n")
	}
}
