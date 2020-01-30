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
	namespaceHost := os.Getenv("NAMESPACE_HOST")
	oasVersionEnv, exists := os.LookupEnv("OAS_VERSION")
	if !exists {
		oasVersionEnv = "v2"
	}
	log.Println("Using OpenAPI version " + oasVersionEnv)
	log.Println("Namespace host " + namespaceHost)
	//servicesEnv := "[\"artmagrepository\",\"complements-generator\",\"eligibility-calculator\",\"family\",\"maskrepository\",\"mediarepository\",\"offerorchestrator\",\"pricerepository\",\"productrepository\",\"reportpriceftp\",\"reportproductga\",\"reportstockftp\",\"search-engine\",\"search-suggestions\",\"stockrepository\",\"storerepository\",\"substitutes-generator\",\"transliteration\",\"variants\",\"visibility\"]"
	if servicesEnv == "" {
		log.Println("Environment variable \"SERVICES\" is empty")
		os.Exit(2)
	}
	services := make([]string, 0)
	parsed := Map(strings.Split(servicesEnv[1:len(servicesEnv)-1], ","), func(s string) interface{} {
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
		for _, service := range services {
			wg.Add(1)
			go checkService(service, oasVersionEnv, serviceAvailability, &wg)
		}
		wg.Wait()
		availableServices := make([]string, 0, len(services))
		for service, available := range serviceAvailability {
			if available {
				availableServices = append(availableServices, service)
			}
		}
		log.Println("Available services: " + strings.Join(availableServices, ", "))
		resultJson, _ := json.Marshal(Map(availableServices, func(service string) interface{} {
			return map[string]string{
				"name": service,
				"url":  "/" + service + "/" + oasVersionEnv + "/api-docs",
			}
		}))
		_ = templateEngine.Execute(w, string(resultJson))
	})
	_ = http.ListenAndServe(":3000", nil)
}

func checkService(service string, oasVersion string, availability map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "http://" + service + "/" + oasVersion + "/api-docs"
	_, err := http.Get(url)
	availability[service] = err == nil
}

func Map(vs []string, f func(string) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
