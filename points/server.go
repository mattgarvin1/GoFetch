package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type PointServer struct {
	DB     *PointsDB
	logger *LogHandler
}

type LogHandler struct {
	logger *log.Logger
}

// json payload
type TXListJSON struct {
	TXList []*Transaction `json:"txList"`
}

func runServer() {

	// --> TESTING local host port 8080 --> //

	/*
		curl http://localhost:8080/points/_status
			RETURNS
		{
			"status": "healthy"
		}
	*/
	port := flag.Uint("port", 8080, "port on which to expose the API")

	// <-- TESTING local host port 8080 <-- //

	// port := flag.Uint("port", 80, "port on which to expose the API")
	logFlags := log.Ldate | log.Ltime
	logger := log.New(os.Stdout, "", logFlags)
	memDB := &PointsDB{PayerBalances: make(map[string]int)}
	server := server().withDB(memDB).withLogger(logger)
	router := server.makeRouter(os.Stdout)
	addr := fmt.Sprintf(":%d", *port)
	httpLogger := log.New(os.Stdout, "", log.LstdFlags)
	httpServer := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     httpLogger,
		Handler:      router,
	}
	httpLogger.Println(fmt.Sprintf("Fetch Points Server serving at %s", httpServer.Addr))
	httpLogger.Fatal(httpServer.ListenAndServe())
}

func (ps *PointServer) withDB(db *PointsDB) *PointServer {
	ps.DB = db
	return ps
}

func (ps *PointServer) withLogger(logger *log.Logger) *PointServer {
	ps.logger = &LogHandler{logger: logger}
	return ps
}

func server() *PointServer {
	return &PointServer{}
}

// NOTE: could add additional API functionality

// POST /points/addTransactions
//
// curl -X POST -d "@txList.json" http://localhost:8080/points/addTransactions
func (ps *PointServer) handleAddTransactions(w http.ResponseWriter, r *http.Request) {

	txList := unmarshalBody(r, &TXListJSON{}).(*TXListJSON)

	// return error, handle error, write proper HTTP response
	ps.addTransactions(txList.TXList)

	// return summary info of transactions added, dates of first and last transactions (?)
	// -> any useful summary statistics from the operation
	j := map[string]string{
		"result": "success",
		"nTX":    strconv.Itoa(len(txList.TXList)),
	}

	writeJSON(w, j)
}

// POST /points/spend
//
// curl -X POST -d "@spendOrder.json" http://localhost:8080/points/spend
func (ps *PointServer) handleSpendOrder(w http.ResponseWriter, r *http.Request) {

	// again, handle malformed or unexpected JSON or not-JSON etc.
	spendOrder := unmarshalBody(r, &SpendOrder{}).(*SpendOrder)
	j, err := ps.spendPoints(spendOrder)
	if err != nil {
		// handle err , log , http response etc.
		fmt.Println("ERROR:", err)
	}

	writeJSON(w, j)
}

// GET /points/payerBalance // NOTE: extend API - could specify individual payer, get back single payer balance
//
// curl http://localhost:8080/points/payerBalance
func (ps *PointServer) handleFetchBalance(w http.ResponseWriter, r *http.Request) {
	// return error, handle error, etc.
	j := ps.fetchBalance()
	writeJSON(w, j)
}

// GET /points/_status -> server status, health check
//
// curl http://localhost:8080/points/_status
func (ps *PointServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	j := map[string]string{
		"status": "healthy",
	}
	writeJSON(w, j)
}

func (ps *PointServer) makeRouter(out io.Writer) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	/*
		router.HandleFunc("/runs/{runID}", server.handleRunLogGET).Methods("GET")
		router.HandleFunc("/runs/{runID}/cancel", server.handleCancelRunPOST).Methods("POST")
	*/

	router.HandleFunc("/points/addTransactions", ps.handleAddTransactions).Methods("POST")
	router.HandleFunc("/points/spend", ps.handleSpendOrder).Methods("POST")

	// NOTE: could add optional {payerID} var in the URL here, if you want just one balance
	// further: could make this a POST, post a list of payer IDs
	// GET with no params returns all payers' balances
	// POST with list of payerID's, only return those payers' balances
	router.HandleFunc("/points/payerBalance", ps.handleFetchBalance).Methods("GET")
	router.HandleFunc("/points/_status", ps.handleHealthCheck).Methods("GET")

	// router.NotFoundHandler = http.HandlerFunc(handleNotFound) // TODO

	//// add middleware ////

	// set "Content-Type: application/json" header - every endpoint returns JSON
	router.Use(ps.setResponseHeader)

	// todo: add auth middleware ???

	///////////////////////

	// remove trailing slashes sent in URLs
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		router.ServeHTTP(w, r)
	})

	// basic logging solution; keeping it for now
	// see: https://godoc.org/github.com/gorilla/handlers#CombinedLoggingHandler
	return handlers.CombinedLoggingHandler(out, handler)
}

// middleware;
// all endpoints return JSON, so just set that response header here
func (ps *PointServer) setResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// unmarshal the request body to the given go struct
// fixme: return error
func unmarshalBody(r *http.Request, v interface{}) interface{} {
	b := body(r)
	err := json.Unmarshal(b, v)
	if err != nil {
		// fixme
		fmt.Println("error unmarshalling: ", err)
	}
	return v
}

// fixme: return error
func body(r *http.Request) []byte {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// fixme
		fmt.Println("error reading body: ", err)
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return b
}

func writeJSON(w http.ResponseWriter, j interface{}) {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	e.Encode(j)
}
