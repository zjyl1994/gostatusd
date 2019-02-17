package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	port := flag.String("port", ":58000", "HTTP listen port")
	flag.Parse()
	http.HandleFunc("/info", getInfo)
	http.HandleFunc("/", index)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, infoJSON())
}

func index(w http.ResponseWriter, r *http.Request) {
	interfaceName := getInterfacesNames()
	infoString := infoJSON()
	t := template.Must(template.ParseFiles("./web.html"))
	t.ExecuteTemplate(w, "web.html", map[string]interface{}{
		"ifName": interfaceName,
		"info":   infoString,
	})
}
