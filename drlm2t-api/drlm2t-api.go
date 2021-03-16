package main

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type ctxKey struct{}

// Http Routing
var routes = []route{
	// File server token protected
	newRoute("GET", "/", staticGet),
	newRoute("GET", "/((images|static|css|js)/[a-zA-Z0-9._/-]+)", staticGet),
	// API functions ////////////////////////////////
	newRoute("GET", "/api/infrastructures", middlewareUserToken(apiGetInfrastructures)),
	newRoute("GET", "/api/infrastructures/([a-zA-Z0-9._-]+)", middlewareUserToken(apiGetInfrastructure)),
	newRoute("GET", "/api/infrastructures/([a-zA-Z0-9._-]+)/hosts/([0-9])+/tests/([0-9-]+)/result", middlewareUserToken(apiGetTestResult)),
	newRoute("POST", "/api/infrastructures/([a-zA-Z0-9._-]+)", middlewareUserToken(apiSetInfrastructure)),
	newRoute("PUT", "/api/infrastructures/([a-zA-Z0-9._-]+)", middlewareUserToken(apiPutInfrastructure)),
	newRoute("DELETE", "/api/infrastructures/([a-zA-Z0-9._-]+)", middlewareUserToken(apiDeleteInfrastructure)),
	newRoute("GET", "/api/running", middlewareUserToken(apiGetRunning)),
	newRoute("GET", "/api/templates", middlewareUserToken(apiGetTemplates)),
	newRoute("POST", "/api/up/([a-zA-Z0-9._-]+)", middlewareUserToken(apiUpTest)),
	newRoute("POST", "/api/down/([a-zA-Z0-9._-]+)", middlewareUserToken(apiDownTest)),
	newRoute("POST", "/api/run/([a-zA-Z0-9._-]+)", middlewareUserToken(apiRunTest)),
	newRoute("POST", "/api/clean/([a-zA-Z0-9._-]+)", middlewareUserToken(apiCleanTest)),

	// User Control Functions ///////////////////////
	newRoute("POST", "/signin", userSignin),
	newRoute("POST", "/logout", userLogout),
	newRoute("POST", "/checkSession", checkSession),
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//http.NotFound(w, r)
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	initDRLM2TConfiguration()
	printDRLMConfiguration()

	// Run HTTPS server with middlewareLog
	log.Println("Running server at https://localhost:8080")
	log.Fatal(http.ListenAndServeTLS(":8080", configDRLM2T.Certificate, configDRLM2T.Key, http.HandlerFunc(middlewareLog(Serve))))
	//log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(middlewareLog(Serve))))
}
