package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ImageProcessingRequest struct {
	ImageURL string `json:"imageUrl"`
}

func main() {
	http.HandleFunc("/process", imageProcessingHandler)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func processImage(url string) {
	fmt.Println("Processing image:", url)
	// TODO 非同期処理
	time.Sleep(3 * time.Second)

	fmt.Println("Image processed:", url)

	//TODO 完了通知　Firebase Cloud Messaging?
	// Node.js に通知を送る
	http.Post("http://node-backend:3000/notify", "application/json",
		strings.NewReader(fmt.Sprintf(`{"imageUrl": "%s", "status": "done"}`, url)))

}

func imageProcessingHandler(w http.ResponseWriter, r *http.Request) {
	var req ImageProcessingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// ゴルーチンで非同期実行
	go processImage(req.ImageURL)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Processing started")
}
