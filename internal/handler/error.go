// Here you can implement custom errors
// and handle them

package handler

import (
	"errors"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/handler/response"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"net/http"
)

var (
	errInvalidArgument = errors.New("invalid argument")
	errIsNotInteger    = errors.New("value is not an integer")
)

func handleErrResponse(rw http.ResponseWriter, err error) {
	defer logger.Info(err)

	switch err {
	default:
		response.Internal(rw, err.Error())
	}
	return
}
