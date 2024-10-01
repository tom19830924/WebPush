package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
)

var subscriptions []webpush.Subscription

// VAPID keys
const (
	publicKey  = "BFl0cadQjaCPWg6EACPAfBgsDyBDorLNGZM-IEUiz4qCUOMtV611r9A-5RwkQo6yq82NOA7NH93YDKmbrDMh4Bg"
	privateKey = "49qZ6KCEiIsSrjnsr-bIG6xnGUpvY8se_I0frBAixJU"
)

func main() {
	http.HandleFunc("/subscribe", handleSubscribe)
	http.HandleFunc("/options", handleOptions)
	http.HandleFunc("/pushMessage", handlePushMessage) // 新的处理程序

	port := 8081 // 使用不同于网页的端口
	log.Printf("服务器正在运行在 http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	var sub webpush.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	subscriptions = append(subscriptions, sub)
	w.WriteHeader(http.StatusCreated)
	log.Printf("订阅对象已保存: %+v", sub)
}

func handlePushMessage(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	payloadBytes, _ := json.Marshal(payload)

	for _, sub := range subscriptions {
		resp, err := webpush.SendNotification(payloadBytes, &sub, &webpush.Options{
			VAPIDPublicKey:  publicKey,
			VAPIDPrivateKey: privateKey,
			TTL:             30,
		})
		if err != nil {
			log.Printf("推送消息发送失败: %v", err)
			continue
		}
		resp.Body.Close()
		log.Printf("推送消息发送成功: %+v", resp)
	}

	w.WriteHeader(http.StatusOK)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
