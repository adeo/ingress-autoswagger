package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr"
	"log"
	"net/http"
	"os"
	"strings"
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
		available := Filter(services, func(service string) bool {
			url := "http://" + namespaceHost + "/" + service + "/" + oasVersionEnv + "/api-docs"
			log.Println("Requesting: " + url)
			_, err := http.Get(url)
			return err == nil
		})

		resultJson, _ := json.Marshal(Map(available, func(service string) interface{} {
			return map[string]string{
				"name": service,
				"url":  "/" + namespaceHost + "/" + service + "/" + oasVersionEnv + "/api-docs",
			}
		}))
		_ = templateEngine.Execute(w, string(resultJson))
	})
	_ = http.ListenAndServe(":3000", nil)
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map(vs []string, f func(string) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
