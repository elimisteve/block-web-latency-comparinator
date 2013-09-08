// Steve Phillips / elimisteve
// 2013.09.07

package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	DEFAULT_LISTEN_ADDR = ":8080"
)

var (
	// `PORT` environment var used by Heroku
	LISTEN_ADDR = ":" + os.Getenv("PORT")

	router = mux.NewRouter()
)

func init() {
	if LISTEN_ADDR == ":" {
		LISTEN_ADDR = DEFAULT_LISTEN_ADDR
	}

	router.HandleFunc("/", ShowSchema).Methods("OPTIONS")
	router.HandleFunc("/", MeasureLatency).Methods("POST")

	http.Handle("/", router)
}

func main() {
	serve(router)
}

func ShowSchema(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(BLOCK_DEFINITION)
}

// MeasureLatency expects {"inputs": {"urls": ["http://..."]}} and
// returns the time taken to perform a HEAD request to the given URL
// in the form {"outputs": [{"url": "http://...", "latency": 100}]}
func MeasureLatency(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read Request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Bad Request
		log.Printf("Error at %s: %v\n", r.URL, err)
		http.Error(w, `{"outputs": []}`, 400)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	input := InputURLs{}
	if err := json.Unmarshal(body, &input); err != nil {
		// Bad Request
		log.Printf("Error at %s: %v\n", r.URL, err)
		http.Error(w, `{"outputs": []}`, 400)
		return
	}

	// Perform HEAD requests and record reponse times
	ch := make(chan *Stopwatch)
	for _, url := range input.Inputs.URLs {
		go timeHead(ch, url)
	}

	// Collect response times
	latencies := make([]*Stopwatch, len(input.Inputs.URLs))
	for i := 0; i < len(input.Inputs.URLs); i++ {
		latencies[i] = <-ch
	}

	// Prepare JSON response
	output := OutputLatencies{Outputs: latencies}

	jsonData, err := json.Marshal(&output)
	if err != nil {
		// Unable to marshal output
		log.Printf("Error at %s: %v\n", r.URL, err)
		http.Error(w, `{"outputs": []}`, 500)
		return
	}

	w.Write(jsonData)
}

func timeHead(ch chan *Stopwatch, url string) {
	stopwatch := &Stopwatch{URL: url}

	// Send `stopwatch` to `ch` when `timeHead` returns
	defer func() {
		ch <- stopwatch
	}()

	// Prepare HEAD request
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Printf("Error creating request to %s: %v\n", url, err)
		stopwatch.Latency = -1
		return
	}

	// Measure latency
	start := time.Now()
	_, err = new(http.Client).Do(req)
	latencyInNanos := time.Since(start)

	if err != nil {
		log.Printf("Error retrieving %s: %v\n", url, err)
		stopwatch.Latency = -1
		return
	}

	// Convert nanoseconds to milliseconds
	stopwatch.Latency = int64(latencyInNanos / 1e6)
}

//
// Types
//

type InputURLs struct {
	Inputs struct {
		URLs []string `json:"urls"`
	} `json:"inputs"`
}

type OutputLatencies struct {
	Outputs []*Stopwatch `json:"outputs"`
}

type Stopwatch struct {
	URL     string `json:"url"`
	Latency int64  `json:"latency"`  // Milliseconds
}

//
// Miscellaneous
//

func serve(h http.Handler) {
	server := &http.Server{
		Addr:           LISTEN_ADDR,
		Handler:        h,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("HTTP server trying to listen on %s...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP listen failed: %v\n", err)
	}
}

var BLOCK_DEFINITION = []byte(`{
  "name": "Web Latency Comparinator",
  "url": "http://web-latency-comparinator.herokuapp.com",
  "description": "Performs a HEAD request to all given URLs in parallel and returns the time taken to receive a response from each.",
  "inputs": {
      "name": "urls",
      "type": "Array",
      "description": "URLs to be visited"
  },
  "outputs": [
    {
      "name": "url",
      "type": "String",
      "description": "Visited URL"
    },
    {
      "name": "latency",
      "type": "Number",
      "description": "URL response time (in milliseconds)"
    }
  ]
}
`)
