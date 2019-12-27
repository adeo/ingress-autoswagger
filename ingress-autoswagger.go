package main

import (
	"encoding/json"
	"github.com/gobuffalo/packr"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	servicesEnv := os.Getenv("SERVICES")
	//servicesEnv := "app1,app2,app3"
	if servicesEnv == "" {
		log.Println("Environment variable \"SERVICES\" is empty")
		os.Exit(2)
	}
	services := strings.Split(servicesEnv, ",")
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
			_, err := http.Get("http://" + service + "/" + service + "/v2/api-docs")
			return err == nil
		})

		resultJson, _ := json.Marshal(Map(available, func(service string) interface{} {
			return map[string]string{
				"name": service,
				"url":  "/" + service + "/v2/api-docs",
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
