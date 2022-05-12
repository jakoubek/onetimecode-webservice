package main

import (
	"expvar"
	"github.com/jakoubek/onetimecode-webservice/internal"
	"net/http"
	"strconv"
	"time"
)

func (app *application) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.onetimecode.net/?ref=apiroot", http.StatusFound)
	}
}

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	requests, _ := strconv.Atoi(expvar.Get("total_requests_received").String())
	data := envelope{
		"version":        version,
		"build_time":     buildTime,
		"server_started": app.startupTime.UTC().String(),
		"requests":       requests,
		"timestamp":      time.Now().Unix(),
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "OK",
		"info":   "API fully operational",
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) numberHandler(w http.ResponseWriter, r *http.Request) {

	rp, err := app.readRequestParameters(r)
	if err != nil {
		app.logError(r, err)
	}

	otc := internal.NewNumericalCode(
		internal.WithLength(rp.length),
		internal.WithMin(rp.min),
		internal.WithMax(rp.max),
		internal.WithGrouping(rp.groupEvery, rp.groupBy),
	)

	data := envelope{
		"result": otc.ResultAsString(),
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) alphanumericHandler(w http.ResponseWriter, r *http.Request) {

	rp, err := app.readRequestParameters(r)
	if err != nil {
		app.logError(r, err)
	}

	otc := internal.NewAlphanumericalCode(
		internal.WithLength(rp.length),
		internal.WithCase(rp.caseStr),
		internal.WithGrouping(rp.groupEvery, rp.groupBy),
	)

	data := envelope{
		"result": otc.ResultAsString(),
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) ksuidHandler(w http.ResponseWriter, r *http.Request) {

	//rp, err := app.readRequestParameters(r)
	//if err != nil {
	//	app.logError(r, err)
	//}

	otc := internal.NewKsuidCode()

	data := envelope{
		"result": otc.ResultAsString(),
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) uuidHandler(w http.ResponseWriter, r *http.Request) {

	rp, err := app.readRequestParameters(r)
	if err != nil {
		app.logError(r, err)
	}

	otc := internal.NewUuidCode(
		internal.WithoutDashesFromBoolean(rp.withoutdashes),
	)

	data := envelope{
		"result": otc.ResultAsString(),
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) ulidHandler(w http.ResponseWriter, r *http.Request) {

	//rp, err := app.readRequestParameters(r)
	//if err != nil {
	//	app.logError(r, err)
	//}

	otc := internal.NewUlidCode()

	data := envelope{
		"result": otc.ResultAsString(),
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) diceHandler(w http.ResponseWriter, r *http.Request) {

	//rp, err := app.readRequestParameters(r)
	//if err != nil {
	//	app.logError(r, err)
	//}

	otc := internal.NewNumericalCode(
		internal.WithMinMax(1, 6),
	)

	data := envelope{
		"result": otc.ResultAsString(),
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) coinHandler(w http.ResponseWriter, r *http.Request) {

	//rp, err := app.readRequestParameters(r)
	//if err != nil {
	//	app.logError(r, err)
	//}

	otc := internal.NewNumericalCode(
		internal.WithMinMax(0, 1),
	)

	var side string
	if otc.NumberCode() == 0 {
		side = "head"
	} else {
		side = "tails"
	}
	data := envelope{
		"result": otc.ResultAsString(),
		"side":   side,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
