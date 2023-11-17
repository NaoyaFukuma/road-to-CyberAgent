package response

import (
	"encoding/json"
	"net/http"
)

func SetStatusAndJson(writer http.ResponseWriter, statusCode int, body interface{}) {
	writer.WriteHeader(statusCode)
	jsonData, err := json.Marshal(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(jsonData)
}
