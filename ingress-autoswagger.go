package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
)

func main() {
	servicesEnv := os.Getenv("SERVICES")
	oasVersionEnv, exists := os.LookupEnv("OAS_VERSION")
	if !exists {
		oasVersionEnv = "v2"
	}
	log.Println("Using OpenAPI version " + oasVersionEnv)
	if servicesEnv == "" {
		log.Println("Environment variable \"SERVICES\" is empty")
		os.Exit(2)
	}
	services := make([]string, 0)
	parsed := mapValues(strings.Split(servicesEnv[1:len(servicesEnv)-1], ","), func(s string) interface{} {
		return s[1 : len(s)-1]
	})

	for _, str := range parsed {
		services = append(services, fmt.Sprintf("%v", str))
	}

	log.Println("Server started on 3000 port!")
	log.Println(services)
	html, err := packr.NewBox("./templates").FindString("index.html")
	if err != nil {
		panic(err)
	}

	templateEngine, err := template.New("index").Parse(html)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serviceAvailability := make(map[string]bool)
		var wg sync.WaitGroup
		var mtx = sync.Mutex{}
		for _, service := range services {
			wg.Add(1)
			go checkService(service, oasVersionEnv, serviceAvailability, &wg, &mtx)
		}
		wg.Wait()
		availableServices := make([]string, 0, len(services))
		for service, available := range serviceAvailability {
			if available {
				availableServices = append(availableServices, service)
			}
		}
		log.Println("Available services: " + strings.Join(availableServices, ", "))
		resultJson, _ := json.Marshal(mapValues(availableServices, func(service string) interface{} {
			return map[string]string{
				"name": service,
				"url":  "/" + service + "/" + oasVersionEnv + "/api-docs",
			}
		}))
		_ = templateEngine.Execute(w, string(resultJson))
	})
	_ = http.ListenAndServe(":3000", nil)
}

func checkService(service string, oasVersion string, availability map[string]bool, wg *sync.WaitGroup, m *sync.Mutex) {
	url := "http://" + service + "/" + oasVersion + "/api-docs"
	_, err := http.Get(url)
	m.Lock()
	availability[service] = err == nil
	m.Unlock()
	wg.Done()
}

func mapValues(vs []string, f func(string) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
