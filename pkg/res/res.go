package res

import (
	"encoding/json"
	"log"
	"net/http"
)

func JsonResp(writer http.ResponseWriter, res any, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	err := json.NewEncoder(writer).Encode(res)
	if err != nil {
		log.Println(err, "error encode response")
	}
}
