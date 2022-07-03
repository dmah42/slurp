package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	pb "github.com/dominichamon/slurp/internal/api/slurp"
)

var (
	port = flag.Int("port", 3333, "The port on which to listen for HTTP requests")

	client pb.SlurpClient
)

func serve(c pb.SlurpClient) error {
	client = c

	r := mux.NewRouter()
	r.HandleFunc("/", root)
	r.HandleFunc("/server/{s}", server)
	//r.HandleFunc("/server/{s}/group/{g}", group)
	//r.HandleFunc("/server/{s}/group/{g}/article/{a}", article)

	log.Printf("listening on :%d", *port)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), r)
}

func root(w http.ResponseWriter, r *http.Request) {
	addrs, err := client.Addresses(r.Context(), &pb.AddressesRequest{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// TODO: template
	w.Write([]byte("<html><body>\n"))
	for _, a := range addrs.GetAddress() {
		w.Write([]byte(fmt.Sprintf("<p><a href=\"/server/%s\">%s</a></p>", a, a)))
	}
	w.Write([]byte("</body></html>\n"))
}

func server(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, ok := vars["s"]
	if !ok {
		w.Write([]byte("s is missing in parameters"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	gs, err := client.Groups(r.Context(), &pb.GroupsRequest{Server: s})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// TODO: template
	w.Write([]byte("<html><body>\n"))
	w.Write([]byte(fmt.Sprintf("<h1>%s</h1>", s)))
	for _, g := range gs.GetGroup() {
		w.Write([]byte(fmt.Sprintf("<p><a href=\"/server/%s/group/%s\">%s</a></p>", s, g, g)))
	}
	w.Write([]byte("</body></html>\n"))
}

/*
func group(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, ok := vars["s"]
	if !ok {
		w.Write([]byte("s is missing in parameters"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	g, ok := vars["g"]
	if !ok {
		w.Write([]byte("g is missing in parameters"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	as, err := n.Articles(s, g)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// TODO: template
	w.Write([]byte("<html><body>\n"))
	w.Write([]byte(fmt.Sprintf("<h1>%s</h1>", s)))
	w.Write([]byte(fmt.Sprintf("<h2>%s</h2>", g)))
	for _, a := range as {
		w.Write([]byte(fmt.Sprintf("<p>%+v</p>\n", a)))
	}
	w.Write([]byte("</body></html>\n"))
}*/
