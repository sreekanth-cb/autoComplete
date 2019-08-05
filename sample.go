package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	gocb "gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocb.v1/cbft"
)

type Movie struct {
	Id         string  `json:"id"`
	Score      float64 `json:"score"`
	MovieTitle string  `json:"movie_title"`
}

var bucket *gocb.Bucket
var STATIC_DIR = "/static/"

func main() {
	var err error
	cluster, _ := gocb.Connect("http://localhost:9000")
	cluster.Authenticate(gocb.PasswordAuthenticator{Username: "Administrator", Password: "asdasd"})
	bucket, err = cluster.OpenBucket("movies", "")
	if err != nil {
		fmt.Printf("cluster.OpenBucket err %+v", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/search", SearchEndpoint).Methods("GET")
	router.PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))
	http.ListenAndServe(":12345", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
}

func SearchEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := request.URL.Query()
	subStrings := strings.Split(params.Get("query"), " ")
	var conjunctionQuery *cbft.ConjunctionQuery
	for _, sq := range subStrings {
		sq := cbft.NewMatchQuery(sq).Field("movie_title").Analyzer("simple") //.Fuzziness(1)
		if conjunctionQuery != nil {
			conjunctionQuery.And(sq)
			continue
		}
		conjunctionQuery = cbft.NewConjunctionQuery(sq)
	}

	query := gocb.NewSearchQuery("FTS", conjunctionQuery).Limit(20)
	query.Fields("movie_title")
	//query := gocb.NewSearchQuery("FTS", cbft.NewMatchQuery(params.Get("query")).Analyzer("standard")).Limit(20)

	result, err := bucket.ExecuteSearchQuery(query)
	if err != nil {
		fmt.Printf("err %+v", err)
	}
	//fmt.Printf("result.Hits() %+v\n\n", result.Hits())
	hits := make([]Movie, 20)
	for _, hit := range result.Hits() {
		hits = append(hits, Movie{
			Id:         hit.Id,
			Score:      hit.Score,
			MovieTitle: hit.Fields["movie_title"],
		})
	}
	json.NewEncoder(response).Encode(hits)
}
