package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"

	sw "github.com/Azure-Samples/openhack-devops-team/apis/trips/tripsgo"
)

var (
	webServerPort    = flag.String("webServerPort", getEnv("WEB_PORT", "8080"), "web server port")
	webServerBaseURI = flag.String("webServerBaseURI", getEnv("WEB_SERVER_BASE_URI", "changeme"), "base portion of server uri")
	client           = appinsights.NewTelemetryClient("618f7639-d08b-4f1f-9a64-51df20f19953")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {

	var debug, present = os.LookupEnv("DEBUG_LOGGING")

	client.Context().Tags.Cloud().SetRole("api-trips")

	metriccpu := appinsights.NewMetricTelemetry("CPU Util")
	client.Track(metriccpu)

	// event := appinsights.NewEventTelemetry("button clicked")
	// event.Properties["property"] = "value"
	// client.Track(event)

	if present && debug == "true" {
		sw.InitLogging(os.Stdout, os.Stdout, os.Stdout)
	} else {
		// if debug env is not present or false, do not log debug output to console
		sw.InitLogging(os.Stdout, ioutil.Discard, os.Stdout)
	}

	sw.Info.Println(fmt.Sprintf("%s%s", "Trips Service Server started on port ", *webServerPort))

	router := sw.NewRouter()

	sw.Fatal.Println(http.ListenAndServe(fmt.Sprintf("%s%s", ":", *webServerPort), router))
}
