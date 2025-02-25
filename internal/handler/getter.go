package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	LimitParam  = "limit"
	OffsetParam = "offset"
)

type Pagination struct {
	Limit  int64
	Offset int64
}

func GetIntQueryParam(r *http.Request, key string) (res int64, err error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return 0, nil
	}
	res, err = strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, errIsNotInteger
	}
	return res, nil
}

func GetIntPathParam(r *http.Request, param string) (id int64, err error) {
	vars := mux.Vars(r)
	idStr, ok := vars[param]
	if !ok {
		return 0, errInvalidArgument
	}

	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errIsNotInteger
	}
	return id, nil
}
