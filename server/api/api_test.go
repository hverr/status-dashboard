package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
)

func (api *API) router() *gin.Engine {
	router := gin.Default()
	if err := api.Install(router); err != nil {
		log.Fatalln("Could not install API:", err)
		return nil
	}
	return router
}

func (api *API) request(method, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, path, body)
	return req
}

func (api *API) call(method, path, clientSecret string, body io.Reader) *httptest.ResponseRecorder {
	router := api.router()
	req := api.request(method, path, body)
	req.Header.Set(server.ClientSecretHeader, clientSecret)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func (api *API) callJSON(method, path, clientSecret string, body interface{}) *httptest.ResponseRecorder {
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatalln("Can't create a JSON call:", err)
	}
	return api.call(method, path, clientSecret, bytes.NewBuffer(b))
}
