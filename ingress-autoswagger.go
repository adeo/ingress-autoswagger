package main

import (
	"encoding/json"
	"github.com/gobuffalo/packr"
	"gopkg.in/robfig/cron.v3"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
)

//service to oas version
var cachedAvailableServices = make([]map[string]string, 0)
var versions = make([]string, 0)
var versionsFormat = ""

func main() {
	refreshCron, exists := os.LookupEnv("REFRESH_CRON")
	if !exists {
		refreshCron = "@every 1m"
	}

	servicesEnv := os.Getenv("SERVICES")
	if servicesEnv == "" {
		log.Println("Environment variable \"SERVICES\" is empty")
		os.Exit(2)
	}
	services := make([]string, 0)
	services = mapValues(strings.Split(servicesEnv[1:len(servicesEnv)-1], ","), func(s string) string {
		return s[1 : len(s)-1]
	})
	sort.Strings(services)

	//set versions
	versionsEnv, versionsEnvExists := os.LookupEnv("VERSIONS")
	
	if versionsEnvExists {
		versions = mapValues(strings.Split(versionsEnv[1:len(versionsEnv)-1], ","), func(s string) string {
			return s[1 : len(s)-1]
		})
	} else {
		versions = []string{"v2", "v3"}
	}

	log.Println("Server started on 3000 port!")
	log.Println("Services:", services)
	log.Println("Discovering versions:", versions)
	html, err := packr.NewBox("./templates").FindString("index.html")
	if err != nil {
		panic(err)
	}

	templateEngine, err := template.New("index").Parse(html)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		refreshCache(services)
		w.WriteHeader(200)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resultJson, _ := json.Marshal(cachedAvailableServices)
		_ = templateEngine.Execute(w, string(resultJson))
	})
	refreshCache(services)

	c := cron.New()
	c.AddFunc(refreshCron, func() {
		log.Println("Cron init")
		refreshCache(services)
		log.Println("Cron has been finished")
	})
	c.Start()

	_ = http.ListenAndServe(":3000", nil)
}

func checkService(service string) {
	passedVersion := ""
	passedFormat := ""
	versionsFormat, versionsFormatExists := os.LookupEnv("VERSION_FORMAT")

	for _, ver := range versions {

		if versionsFormatExists {
			log.Println("Trying swagger format: " + versionsFormat)
			if versionsFormat != "json" { passedFormat = "." + versionsFormat }
		}
		url := "http://" + service + "/" + ver + "/api-docs" + passedFormat
		log.Println("Trying url: " + url)
		resp, err := http.Get(url)

		if err == nil && strings.Contains(resp.Status, "200") {
			passedVersion = ver
		}
		if resp != nil {
			log.Println("for version " + ver + " status code is " + resp.Status)
			resp.Body.Close()
		}
	}

	log.Println("for " + service + " version is '" + passedVersion + "'")
	if passedVersion != "" {
		cachedAvailableServices = append(cachedAvailableServices, map[string]string{
			"name": service,
			"url":  "/" + service + "/" + passedVersion + "/api-docs" + passedFormat,
		})
	}
}

func mapValues(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func refreshCache(services []string) {
	log.Println("Refresh start")
	cachedAvailableServices = cachedAvailableServices[:0]
	for _, service := range services {
		checkService(service)
	}
	log.Println("Refresh finish")
}
