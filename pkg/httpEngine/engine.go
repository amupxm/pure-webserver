package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/amupxm/pure-webserver/config"
	"github.com/amupxm/pure-webserver/constants"
)

type (
	Server interface {
		// AddHandler adds a new handler to the server
		AddHandler(path, method string, handler func(c *ServerContext))
		//mainEngineHandler is the main handler which calls on every request to find the right handler
		mainEngineHandler(w http.ResponseWriter, r *http.Request)
		// filterRoutesByPath is a helper function to filter routes by path
		filterRoutesByPath(path string) []serverRoutes
		// StartServer starts the server
		StartServer(port string)
		//filterMatchedRoutesByMethod is a helper function to filter matched routes by method
		filterMatchedRoutesByMethod(method string, mc []serverRoutes) []serverRoutes
		// extractURLParams is a helper function to extract url params
		extractURLParams(master, slave string) map[string]string
	}
	server struct {
		Port string
	}
	serverRoutes struct {
		Path          string
		RequestMethod string
		Handler       func(c *ServerContext)
		URLParams     *map[string]string
	}

	ServerContext struct {
		Response  http.ResponseWriter
		Request   *http.Request
		URLParams map[string]string
	}
	ServerContextInterface interface {
		// ErrorHandler is a helper function to handle errors and return them to the client
		ErrorHandler(code int, err error)
		// GetURLParam is a helper function to get url param
		GetURLParam(param string) (string, error)
		// JSON is a helper function to return json response
		JSON(core int, response interface{})
		// BindToJson is a helper function to bind struct to json
		BindToJson(c interface{}) error
	}
)

var serverRoutesMap []serverRoutes

// NewServer creates new server instance with port defined in config
func NewServer() Server {
	serverAbstract := &server{
		Port: config.AppConf.Http.Port,
	}
	return serverAbstract
}

// StartServer starts the server
func (s *server) StartServer(port string) {
	router := http.NewServeMux()
	router.HandleFunc("/", s.mainEngineHandler)
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	httpServer.ListenAndServe()
}

// AddHandler adds a new handler to the server
func (s *server) AddHandler(path, method string, handler func(c *ServerContext)) {
	serverRoutesMap = append(serverRoutesMap, serverRoutes{
		Path:          path,
		RequestMethod: method,
		Handler:       handler,
		URLParams:     nil,
	})
}

//mainEngineHandler is the main handler which calls on every request to find the right handler
func (s *server) mainEngineHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	// to Check if is path allowed
	matchedRoutes := s.filterRoutesByPath(r.RequestURI)
	if len(matchedRoutes) == 0 {
		http.NotFound(w, r)
		return
	}

	// to check if the method is allowed
	matchedRoutes = s.filterMatchedRoutesByMethod(method, matchedRoutes)
	if len(matchedRoutes) != 1 {
		http.NotFound(w, r)
		return
	}
	matchedRoutes[0].Handler(
		&ServerContext{
			Response:  w,
			Request:   r,
			URLParams: s.extractURLParams(matchedRoutes[0].Path, r.RequestURI),
		},
	)
}

// extractURLParams is a helper function to extract url params
func (s *server) extractURLParams(master, slave string) map[string]string {
	var res = map[string]string{}
	masterArr := strings.Split(master, "/")
	slaveArr := strings.Split(slave, "/")
	for i, str := range masterArr {
		if strings.Contains(str, ":") {
			// to compare if be mistake in arr
			// FIXME : may cause error on master or slave size difference
			res[strings.ReplaceAll(str, ":", "")] = slaveArr[i]
		}
	}
	return res
}

// filterRoutesByPath is a helper function to filter routes by path
func (s *server) filterRoutesByPath(path string) []serverRoutes {
	var matchedRoutes []serverRoutes
	for _, route := range serverRoutesMap {
		// add extra slash to the end of pathes
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}
		if !strings.HasSuffix(route.Path, "/") {
			route.Path = route.Path + "/"
		}
		splittedMasterRoute := strings.Split(route.Path, "/") // it will be simething like ["s","s2","s3",":id"] which id is structure of route
		splittedSlaveRoute := strings.Split(path, "/")        // it will be something like ["s","s2","s3","1212"] which 1212 is id in request

		if len(splittedMasterRoute) == len(splittedSlaveRoute) {
			matched := true
			for i, c := range splittedMasterRoute {
				if c != splittedSlaveRoute[i] && !strings.Contains(c, ":") {
					matched = false
					break
				}
			}
			if matched {
				matchedRoutes = append(matchedRoutes, route)
			}
		}
	}
	return matchedRoutes
}

//filterMatchedRoutesByMethod is a helper function to filter matched routes by method
func (s *server) filterMatchedRoutesByMethod(method string, mc []serverRoutes) []serverRoutes {
	var matchedRoutes []serverRoutes
	for _, route := range mc {
		if route.RequestMethod == method {
			matchedRoutes = append(matchedRoutes, route)
		}
	}
	return matchedRoutes
}

// GetURLParam is a helper function to get url param
func (s *ServerContext) GetURLParam(param string) (string, error) {
	if s.URLParams[param] == "" {
		return "", errors.New(constants.NoParam)
	}
	return s.URLParams[param], nil
}

// ErrorHandler is a helper function to handle errors and return them to the client
func (s *ServerContext) ErrorHandler(code int, err error) {
	s.JSON(code, map[string]string{
		"error": err.Error(),
	})
}

// JSON is a helper function to return json response
func (s *ServerContext) JSON(core int, response interface{}) {
	s.Response.Header().Set("Content-Type", "application/json")
	s.Response.WriteHeader(core)
	jsoned, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	s.Response.Write(
		[]byte(
			jsoned,
		),
	)
}

// BindToJson is a helper function to bind struct to json
func (s *ServerContext) BindToJson(c interface{}) error {
	body, err := ioutil.ReadAll(s.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(
		[]byte(
			body,
		),
		c,
	)
	if err != nil {
		return err
	}
	return nil
}
