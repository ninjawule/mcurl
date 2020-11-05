package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//---------------------------------------------------------------------------------------------------------------------
// main function : reading the options, then sending requests (almost) in parallel
//---------------------------------------------------------------------------------------------------------------------

func main() {

	// how much time does this all take ?
	globalStart := time.Now()

	// what are the specs for our respects here ?
	opt := getArgs()

	// a bit of logging
	println("")
	log.Printf("-----------------------------------------------------------------------------------------------")
	log.Printf("--- STARTING the testing of %d 'parallel' requests @ %s (%d times)", opt.nbReqs, opt.url, opt.repeat)
	log.Printf("-----------------------------------------------------------------------------------------------")

	// a bit more of logging
	if opt.verbose {
		log.Printf(opt.String())
	}

	// running the main loop for a number of times, and gathering avg and amx times
	var avgTime, maxTime time.Duration
	for try := 1; try <= opt.repeat; try++ {

		currentAvg, currentMax := mainLoop(try, opt)
		avgTime = avgTime + currentAvg
		if currentMax > maxTime {
			maxTime = currentMax
		}

		// allowing all the systems to rest a bit
		time.Sleep(1 * time.Second)
	}
	avgTime = time.Duration(float64(avgTime) / float64(opt.repeat))

	log.Printf("-----------------------------------------------------------------------------------------------")
	log.Printf("--- Done in %s. Average / Max request time: %s / %s", time.Since(globalStart), avgTime, maxTime)
	log.Printf("-----------------------------------------------------------------------------------------------")
}

// the main loop, launched as many as the desired number of tries
func mainLoop(try int, opt *options) (avg, max time.Duration) {

	loopStart := time.Now()

	log.Printf("--- Doing try n°%d", try)

	// everybody will wait that everybody finishes
	wg := new(sync.WaitGroup)

	// we'll keep track of all the times spent for each request
	times := make(chan time.Duration, opt.nbReqs)

	// performing n requests
	for i := 0; i < opt.nbReqs; i++ {
		wg.Add(1)
		go sendRequest(i, wg, opt, times)
	}

	// let's all wait here for all the requests to be finished
	wg.Wait()

	// computing the average time spent per request
	close(times)
	var avgTime, maxTime time.Duration
	for spent := range times {
		if spent > maxTime {
			maxTime = spent
		}
		avgTime = avgTime + spent
	}
	avgTime = time.Duration(float64(avgTime) / float64(opt.nbReqs))

	log.Printf("--- --> Done in %s. Average / Max request time: %s / %s", time.Since(loopStart), avgTime, maxTime)

	return avgTime, maxTime
}

//---------------------------------------------------------------------------------------------------------------------
// actually performing 1 request
//---------------------------------------------------------------------------------------------------------------------

func sendRequest(num int, wg *sync.WaitGroup, opt *options, times chan time.Duration) {

	reqID := num + 1

	defer wg.Done()

	if opt.verbose {
		log.Printf("Starting request %d", reqID)
	}

	// preparing the request
	var body io.Reader
	if opt.data != "" {
		body = strings.NewReader(opt.data)
	}
	httpRequest, errReq := http.NewRequest(opt.method, opt.url, body)
	if errReq != nil {
		panic("Issue while initialising a request")
	}

	// setting the headers
	for key, value := range opt.getHeadersMap() {
		httpRequest.Header.Set(key, value)
	}

	// running
	client := &http.Client{}

	// getting the time before we actually request the external service
	requestStartTime := time.Now()

	// processing the request by calling the remote URL using our client (and timing it)
	resp, errResponse := client.Do(httpRequest)

	// showing the request time
	spent := time.Since(requestStartTime)
	if opt.verbose {
		log.Printf("Request %d done in %s", reqID, spent)
	}
	times <- spent

	// handling errors such as timeout
	if errResponse != nil {
		err("Issue while running a request %s", errResponse)
	}

	// very important to prevent leaks !
	// cf. https://husobee.github.io/golang/memory/leak/2016/02/11/go-mem-leak.html
	// has to be done right before the read, since before that the body is nil
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			err("Could not close the body properly")
		}
	}()

	// reading the response body
	if _, errRead := ioutil.ReadAll(resp.Body); errRead != nil {
		err("Issue while reading the response body: %s", errRead)
	}
}

//---------------------------------------------------------------------------------------------------------------------
// mcurl options
//---------------------------------------------------------------------------------------------------------------------

// they're all here
type options struct {
	method         string
	url            string
	nbReqs         int
	repeat         int
	data           string
	headersStrings []string
	headersMap     map[string]string
	verbose        bool
}

// showing off
func (opt *options) String() string {
	data := opt.data
	if len(data) > 128 {
		data = data[0:128] + " [...]"
	}
	var headers []string
	for key, value := range opt.getHeadersMap() {
		headers = append(headers, key+": "+value)
	}
	return fmt.Sprintf(
		"Running options =\n"+
			"method  : %s\n"+
			"url     : %s\n"+
			"nbReqs  : %d\n"+
			"repeat  : %d\n"+
			"data    : %s\n"+
			"headers : %s\n",
		opt.method, opt.url, opt.nbReqs, opt.repeat, data, strings.Join(headers, ", "),
	)
}

// returns the headers as a key-value map
func (opt *options) getHeadersMap() map[string]string {
	if opt.headersMap == nil {
		opt.headersMap = map[string]string{}
		for _, header := range opt.headersStrings {
			sepIndex := strings.Index(header, ":")
			opt.headersMap[header[0:sepIndex]] = strings.TrimSpace(header[sepIndex+1:])
		}
	}
	return opt.headersMap
}

// getArgs parses the args for our parallel testing
func getArgs() (opt *options) {

	opt = &options{}

	var file string // is there a file to read the data from ?

	pflag.StringVarP(&opt.method, "method", "X", "", "the HTTP method to use (GET, POST, etc)")
	pflag.StringVarP(&opt.url, "url", "U", "", "the URL to work with")
	pflag.IntVarP(&opt.nbReqs, "requests", "n", 10, "the number of requests to send in 'parallel'")
	pflag.IntVarP(&opt.repeat, "repeat", "r", 1, "the number of tries to run the main loop, to gain some confidence on the results")
	pflag.StringVarP(&opt.data, "data", "d", "", "the data to send within the request body")
	pflag.StringVarP(&file, "file", "i", "", "the file to fetch the data from, to send within the request body")
	pflag.StringArrayVarP(&opt.headersStrings, "header", "H", nil, "the request headers")
	pflag.BoolVarP(&opt.verbose, "verbose", "v", false, "verbose + debug mode")
	pflag.Parse()

	if opt.method == "" {
		err("No HTTP method provided!")
	}

	if opt.url == "" {
		err("No URL provided!")
	}

	if file != "" {
		fileContent, errRead := ioutil.ReadFile(file)
		if errRead != nil {
			err("Could not read the given file '%'. Cause: %s", file, errRead)
		}
		opt.data = string(fileContent)
	}

	return
}

// if we need a 'clean' exit
func err(format string, params ...interface{}) {
	log.Println(fmt.Sprintf("\n\nERROR: "+format, params...))
	os.Exit(1)
}
