// Copyright (c) 2012 Ernest Micklei. All rights reserved.

package restful

import (
	"net/http"
)

type Dispatcher interface {
	Routes() []Route
	RootPath() string
	//	rootRegEx
}

// Collection of registered Dispatchers that can handle Http requests
var webServices = []Dispatcher{}
var isRegisteredOnRoot = false

// Register a new Dispatcher add it to the http listeners.
// Check its root path to see if 
func Add(service Dispatcher) {
	webServices = append(webServices, service)
	if !isRegisteredOnRoot {
		http.HandleFunc("/", Dispatch)
		isRegisteredOnRoot = true
	}
}

// Dispatch the incoming Http Request to a matching Dispatcher.
// Matching algorithm is conform http://jsr311.java.net/nonav/releases/1.1/spec/spec.html, see jsr311.go
func Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// step 1. Identify the root resource class (Dispatcher)
	dispatcher, finalMatch, err := detectDispatcher(httpRequest.URL.Path, webServices)
	if err != nil {
		httpWriter.WriteHeader(http.StatusNotFound)
		return
	}
	// step 2. Obtain the object (dispatcher) that will handle the request and a set of candidate methods
	routes := selectRoutes(dispatcher, finalMatch)
	// step 3. Identify the method (Route) that will handle the request
	route, detected := detectRoute(routes, httpWriter, httpRequest)
	if detected {
		route.dispatch(httpWriter, httpRequest)
	}
	// a response has already been written
}
