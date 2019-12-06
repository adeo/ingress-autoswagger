package main

import (
	"encoding/json"
	"github.com/gobuffalo/packr"
	"log"
	"net/http"
	"os"
	"sort"
	"text/template"
)

func main() {
	servicesEnv := os.Getenv("SERVICES")
	if servicesEnv == "" {
		log.Println("Environment variable \"SERVICES\" is empty")
		os.Exit(2)
	}

	services := createServicesJson(servicesEnv)
	if services == "null" {
		log.Println("Incorrect variable \"SERVICES\" or no services with swagger. Exit")
		os.Exit(0)
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
		_ = templateEngine.Execute(w, services)
	})
	_ = http.ListenAndServe(":3000", nil)
}

func createServicesJson(servicesEnv string) string {
	var services map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(servicesEnv), &services); err != nil {
		panic(err)
	}

	var servicesList []map[string]string
	for service, params := range services {
		if (params["swagger"] != false) && (params["skip"] != true) {
			log.Println("Generating service: " + service)
			serviceMap := map[string]string{"name": service, "url": "/" + service + "/v2/api-docs"}

			servicesList = append(servicesList, serviceMap)
		}
	}

	sort.SliceStable(servicesList, func(i, j int) bool {
		return servicesList[i]["name"] < servicesList[j]["name"]
	})
	resultJson, _ := json.Marshal(servicesList)

	return string(resultJson)
}
