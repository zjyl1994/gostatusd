package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var tpl string

func main() {
	var err error
	tpl, err = getTemplate()
	if err != nil {
		log.Fatalln("Unpack web panel failed", err)
	}
	port := flag.String("port", ":58000", "HTTP listen port")
	flag.Parse()
	http.HandleFunc("/info", getInfo)
	http.HandleFunc("/", index)
	err = http.ListenAndServe(*port, nil)
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
	t := template.Must(template.New("tpl").Parse(tpl))
	t.ExecuteTemplate(w, "tpl", map[string]interface{}{
		"ifName": interfaceName,
		"info":   infoString,
	})
}
