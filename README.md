# Ingress Autoswagger
**Why:** Automatization is good. Human manual job should disappear. There are no need to create swagger-ui for each microservice.

**What:** Generates UI for all services in provided environment variable.

## Summary

The main reason for this tool is using it for group of microservices launched via kubernetes and exposed with Ingress. 
So, you shold have list of microservices launched in different paths and each instance should expose /v2/api-docs Swagger annotation.
After that you can run this application on Ingress root level (/) and this tool will start one Swagger UI for all microservices.

![Main window screen](https://github.com/adeo/ingress-autoswagger/raw/master/docs/main_window.png)

## Usage

### With docker
docker run -it -e SERVICES="[\"plaster-calculator\",\"product-binder\"]" docker-devops.art.lmru.tech/bricks/ingress-autoswagger:3.0

### Without docker
SERVICES="[\"plaster-calculator\",\"product-binder\"]" go run ingress-autoswagger.go 

After run you can open http://localhost:3000 in browser.

## Development & Build

0. The tool written in simple Go language, so one that you need it to have installed Go.
1. Install dependencies
go get -u github.com/gobuffalo/packr/packr
2. Build with packr (syntax the same with typical 'go build' command)
packr build .

## Maintainers

Dmitrii Sugrobov dmitrii.sugrobov@leroymerlin.ru
