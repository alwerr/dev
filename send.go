package dev

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Map map[string]interface{}

func Send(w http.ResponseWriter, s interface{}) {
	io.WriteString(
		w, string(s.(string)),
	)
}
func SendOk(w http.ResponseWriter, s string) {
	io.WriteString(w, `{"ok":"`+s+`"}`)
}

func SendErr(w http.ResponseWriter, s string) {
	io.WriteString(w, `{"err":"`+s+`"}`)
}
func SendJsn(w http.ResponseWriter, jsns any) {
	jsn, err := json.Marshal(jsns)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	io.WriteString(w, string(jsn))
}
func SendByte(w http.ResponseWriter, jsn interface{}) {
	io.WriteString(w, string(jsn.(string)))
}
