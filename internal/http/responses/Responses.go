package responses

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var statusOK = []byte("( ͡ᵔ ͜ʖ ͡ᵔ ) 200 OK")
var statusCreated = []byte("( ͡ᵔ ͜ʖ ͡ᵔ ) 201 CREATED")
var statusBadRequest = []byte("( ͠° ͟ʖ ͡°) 400 BAD REQUEST")
var statusUnauthorized = []byte("( ͠° ͟ʖ ͡°) 401 UNAUTHORIZED")
var statusForbidden = []byte("( ͠° ͟ʖ ͡°) 403 FORBIDDEN")
var statusNotFound = []byte("( ͡° ʖ̯ ͡°) 404 NOT FOUND")
var statusMethodNotAllowed = []byte("( ͠° ͟ʖ ͡°) 405 METHOD NOT ALLOWED")
var statusInternalServerError = []byte("( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR")

func ResponseOK(w http.ResponseWriter, s interface{}) {
	jsonData, _ := json.Marshal(s)
	w.WriteHeader(http.StatusOK)
	var length int
	if s == nil {
		length, _ = w.Write(statusOK)
	} else {
		length, _ = w.Write(jsonData)
	}
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
func ResponseCreated(w http.ResponseWriter, s interface{}) {
	jsonData, _ := json.Marshal(s)
	w.WriteHeader(http.StatusCreated)
	var length int
	if s == nil {
		length, _ = w.Write(statusCreated)
	} else {
		length, _ = w.Write(jsonData)
	}
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
func ResponseNoContent(w http.ResponseWriter, s interface{}) {
	jsonData, _ := json.Marshal(s)
	w.WriteHeader(http.StatusNoContent)
	var length int
	if s == nil {
		length, _ = w.Write(statusCreated)
	} else {
		length, _ = w.Write(jsonData)
	}
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
func ResponseBadRequest(w http.ResponseWriter, s interface{}) {
	jsonData, _ := json.Marshal(s)
	w.WriteHeader(http.StatusBadRequest)
	var length int
	if s == nil {
		length, _ = w.Write(statusBadRequest)
	} else {
		length, _ = w.Write(jsonData)
	}
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
func ResponseUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	length, _ := w.Write(statusUnauthorized)
	w.Header().Set("Content-Length", strconv.Itoa(length))
}

func ResponseForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	length, _ := w.Write(statusForbidden)
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
func ResponseNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	length, _ := w.Write(statusNotFound)
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
func ResponseMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	length, _ := w.Write(statusMethodNotAllowed)
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
func ResponseInternalServerError(w http.ResponseWriter, s interface{}) {
	jsonData, _ := json.Marshal(s)
	w.WriteHeader(http.StatusInternalServerError)
	var length int
	if s == nil {
		length, _ = w.Write(statusInternalServerError)
	} else {
		length, _ = w.Write(jsonData)
	}
	w.Header().Set("Content-Length", strconv.Itoa(length))
}
