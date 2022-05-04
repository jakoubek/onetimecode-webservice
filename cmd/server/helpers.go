package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(status)
	w.Write(js)

	return nil
}

type reqParams struct {
	length        int
	min           int
	max           int
	caseStr       string
	withoutdashes bool
	groupBy       string
	groupEvery    int
}

func (app *application) readRequestParameters(r *http.Request) (reqParams, error) {

	rp := reqParams{
		length:        -1,
		min:           -1,
		max:           -1,
		caseStr:       "",
		withoutdashes: false,
		groupBy:       "",
		groupEvery:    -1,
	}

	if r.URL.Query().Has("length") {
		length := r.URL.Query().Get("length")
		lengthVal, err := strconv.Atoi(length)
		if err != nil {
			return rp, errors.New(fmt.Sprintf("invalid length parameter (%s)", length))
		}
		rp.length = lengthVal
	}

	if r.URL.Query().Has("min") {
		min := r.URL.Query().Get("min")
		minVal, err := strconv.Atoi(min)
		if err != nil {
			return rp, errors.New(fmt.Sprintf("invalid min parameter (%s)", min))
		}
		rp.min = minVal
	}

	if r.URL.Query().Has("max") {
		max := r.URL.Query().Get("max")
		maxVal, err := strconv.Atoi(max)
		if err != nil {
			return rp, errors.New(fmt.Sprintf("invalid max parameter (%s)", max))
		}
		rp.max = maxVal
	}

	if r.URL.Query().Has("case") {
		caseStr := r.URL.Query().Get("case")
		rp.caseStr = strings.TrimSpace(strings.ToLower(caseStr))
	}

	if r.URL.Query().Has("withoutdashes") {
		rp.withoutdashes = true
	}

	if r.URL.Query().Has("group_by") {
		groupByStr := r.URL.Query().Get("group_by")
		rp.groupBy = strings.TrimSpace(strings.ToLower(groupByStr))
	}

	if r.URL.Query().Has("group_every") {
		groupEveryStr := r.URL.Query().Get("group_every")
		groupEvery, err := strconv.Atoi(groupEveryStr)
		if err != nil {
			return rp, errors.New(fmt.Sprintf("invalid group_every parameter (%s)", groupEveryStr))
		}
		rp.groupEvery = groupEvery
	}

	return rp, nil
}
