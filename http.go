package dev

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
)

type Ctx struct {
	W    http.ResponseWriter
	R    *http.Request
	Auth Claims
}

func (c Ctx) Send(w interface{}) {
	switch v := w.(type) {
	// Marsheld json
	case Map:
		jsn, err := json.Marshal(w)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}
		io.WriteString(c.W, string(jsn))
	case []byte:
		c.W.Write(w.([]byte))
	case string:
		io.WriteString(c.W, string(w.(string)))
	default:
		_ = v
		print(v.(string))

	}
}
func (c Ctx) PathEnd() string {
	return path.Base(c.R.URL.Path)
}
func (c Ctx) Param(q string) string {
	return c.R.URL.Query().Get(q)
}
func (c Ctx) Sends(w interface{}) {
	io.WriteString(c.W, string(w.(string)))

}
func (c Ctx) SendErr(w string) {
	io.WriteString(c.W, `{"err":"`+w+`"}`)
}

type Fld struct {
	Err  bool
	Null bool
	Val  string
	Name string
}

func (c Ctx) Field(name string, chk string) Fld {
	v := c.R.FormValue(name)
	var fld = Fld{Name: name, Err: false, Val: v}
	// if len(v) == 0 {
	// 	fld.Null = true
	// }
	if chk == `name` && len(v) < 3 {
		fld.Err = true
		return fld
	}

	if chk == `email` {
		match, err := regexp.MatchString(`^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$`, v)
		if match && err != nil {
			fld.Err = true
			return fld
		}

	}
	if chk == `password` && len(v) < 3 {
		fld.Err = true
		return fld
	}
	return fld
}
func Get(path string, cbk func(c Ctx)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method != http.MethodGet {
			io.WriteString(w, `err`)
			return
		}
		cbk(Ctx{W: w, R: r})
	})
}

//	func Gets(path string, cbk func(c Ctx)) {
//		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
//			if r.Method != http.MethodGet {
//				io.WriteString(w, `err`)
//				return
//			}
//			cbk(Ctx{W: w, R: r})
//		})
//	}
func Post(path string, cbk func(c Ctx)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			io.WriteString(w, `err`)
			return
		}
		cbk(Ctx{W: w, R: r})
	})
}
func Gets(path string, cbk func(c Ctx)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method != http.MethodGet {
			io.WriteString(w, `err`)
			return
		}
		tkn, err := Signed(r)
		if err {
			io.WriteString(w, `Sign err`)
			return
		}
		cbk(Ctx{W: w, R: r, Auth: tkn})
	})
}
func Posts(path string, cbk func(c Ctx)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method != http.MethodPost {
			io.WriteString(w, `err`)
			return
		}
		tkn, err := Signed(r)
		if err {
			io.WriteString(w, `Sign err`)
			return
		}
		cbk(Ctx{W: w, R: r, Auth: tkn})
	})
}
func Cdn(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
	http.Handle(`/`+path+`/`, http.StripPrefix(`/`+path+`/`, http.FileServer(http.Dir(path))))

}

func Serve(path int) {
	fmt.Printf("Serving on port %d\n", path)
	http.ListenAndServe(fmt.Sprintf(":%d", path), nil)
}
