package handlers

import (
	"encoding/json"
	"net/http"
)

// var (
// 	DecodeJWTFunc = helpers.DecodeJWT
// 	FetchSecrets  = helpers.FetchSecrets
// 	SendToKafka   = helpers.SendToKafka
// )

// writeJSONResponse writes a JSON response with a given status code and message
func writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	responseData := map[string]interface{}{
		"status":  statusCode,
		"message": message,
	}
	responseJSON, _ := json.Marshal(responseData)
	w.Write(responseJSON)
}

// SendMessageHandler handles the sending of a message to Kafka
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {

}
