package router

import (
	"fmt"

	"net/http"
	"os"

	"github.com/Sicarii-NaumaN/poloniex-gateway-api/middleware"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/swagger"
	"github.com/gorilla/mux"
	v3 "github.com/swaggest/swgui/v3"
)

// Swagger struct
type Swagger struct {
	nameAPI  string
	handlers []swagger.Handler

	JSON string
}

func NewAPI(nameAPI string, appPort int, isSwaggerCreated bool, handlers []swagger.Handler) *mux.Router {
	if len(handlers) == 0 {
		return nil
	}
	r := mux.NewRouter()

	r.Use(middleware.PanicMid)

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.HandlerFunc).Methods(h.Method)
	}

	if !isSwaggerCreated {
		swaggerDoc, err := swagger.GenerateDoc(nameAPI, handlers)
		if err != nil {
			logger.Fatal(err)
		}

		err = os.WriteFile("./internal/router/json/swagger.json", []byte(swaggerDoc), os.ModePerm)
		if err != nil {
			logger.Fatal(err)
		}
	}

	fs := http.FileServer(http.Dir("./internal/router/json"))
	http.Handle("/swagger.json", fs)
	http.Handle("/", v3.NewHandler(nameAPI, "/swagger.json", "/docs"))

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", appPort+1), nil)
	}()

	return r
}
