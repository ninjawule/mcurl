package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
	"strings"
)

//---------------------------------------------------------------------------------------------------------------------
// handling the options passed to the
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
	duration       int64
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
			"method        : %s\n"+
			"url           : %s\n"+
			"nbReqs        : %d\n"+
			"repeat        : %d\n"+
			"duration (ms) : %d\n"+
			"data          : %s\n"+
			"headers       : %s\n",
		opt.method, opt.url, opt.nbReqs, opt.repeat, opt.duration, data, strings.Join(headers, ", "),
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
	pflag.IntVarP(&opt.repeat, "repeat", "r", 1, "the number of tries to run the main loop, to gain some confidence on the results; 0 means endless repeats")
	pflag.StringVarP(&opt.data, "data", "d", "", "the data to send within the request body")
	pflag.StringVarP(&file, "file", "i", "", "the file to fetch the data from, to send within the request body")
	pflag.StringArrayVarP(&opt.headersStrings, "header", "H", nil, "the request headers")
	pflag.BoolVarP(&opt.verbose, "verbose", "v", false, "verbose + debug mode")
	pflag.Int64VarP(&opt.duration, "duration", "t", 0, "the total duration in milliseconds; 0 means no limitation")
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
