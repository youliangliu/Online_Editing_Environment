package main

import (
	"fmt"
	"log"
	"net/http"
)
func hello(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	http.Error(w, "404 not found.", http.StatusNotFound)
	//	return
	//}
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	log.Println(r.FormValue("content"))

	//switch r.Method {
	//case "GET":
	//	a := "name:"
	//	w.Write([]byte(a))
	//	w.Write([]byte(r.FormValue("name")))
	//default:
	//	fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	//}
}

func update(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	http.Error(w, "404 not found.", http.StatusNotFound)
	//	return
	//}
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		a := "111"
		w.Write([]byte(a))
		//w.Write([]byte(r.FormValue("name")))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/send", hello)
	http.HandleFunc("/update", hello)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}