package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Info struct {
	Page   string
	Values []string
}

var distInfo Info = Info{"distance", []string{"inch", "foot", "mile", "meter", "kilometer"}}

var distInMeters map[string]float64 = map[string]float64{
	"inch":      .0254,
	"foot":      .3048,
	"mile":      1609.34,
	"meter":     1,
	"kilometer": 1000,
}

var weightInGrams map[string]float64 = map[string]float64{
	"gram":     1,
	"kilogram": 1000,
	"ounce":    28.35,
	"pound":    453.592,
}

var weightInfo Info = Info{"weight", []string{"gram", "kilogram", "ounce", "pound"}}

var tempInCelsius map[string]func(float64, bool) float64 = map[string]func(float64, bool) float64{
	"fahrenheit": func(num float64, from bool) float64 {
		if from {
			return (num - 32) * (5.0 / 9.0)
		}
		return num*(9.0/5.0) + 32
	},
	"celsius": func(num float64, from bool) float64 {
		return num
	},
	"kelvin": func(num float64, from bool) float64 {
		if from {
			return num - 273.15
		}
		return num + 273.15
	},
}

var tempInfo Info = Info{"temperature", []string{"fahrenheit", "celsius", "kelvin"}}

var resultTemplate *template.Template
var errorTemplate *template.Template
var formTemplate *template.Template

func main() {

	errorTemplate = template.Must(template.ParseFiles("./files/templates/error.html", "./files/templates/header.html", "./files/templates/footer.html"))
	resultTemplate = template.Must(template.ParseFiles("./files/templates/result.html", "./files/templates/header.html", "./files/templates/footer.html"))
	formTemplate = template.Must(template.ParseFiles("./files/templates/form.html", "./files/templates/header.html", "./files/templates/footer.html"))

	http.HandleFunc("GET /distance", handleDistForm)
	http.HandleFunc("POST /distance", handleDistResult)
	http.HandleFunc("GET /weight", handleWeightForm)
	http.HandleFunc("POST /weight", handleWeightResult)
	http.HandleFunc("GET /temperature", handleTempForm)
	http.HandleFunc("POST /temperature", handleTempResult)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files/static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleDistForm(w http.ResponseWriter, r *http.Request) {
	err := formTemplate.Execute(w, distInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleWeightForm(w http.ResponseWriter, r *http.Request) {
	err := formTemplate.Execute(w, weightInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleTempForm(w http.ResponseWriter, r *http.Request) {
	err := formTemplate.Execute(w, tempInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleDistResult(w http.ResponseWriter, r *http.Request) {
	// parse form values
	if r.FormValue("num") == "" {
		errorPage(w, errors.New("please enter a value"))
		return
	}
	num, err := strconv.ParseFloat(r.FormValue("num"), 64)
	if err != nil {
		// should do error page instead
		errorPage(w, err)
		return
	}
	to := r.FormValue("to")
	from := r.FormValue("from")

	num *= distInMeters[from]
	num /= distInMeters[to]
	vals := struct {
		Num  float64
		Page string
	}{num, "distance"}
	err = resultTemplate.Execute(w, vals)
	if err != nil {
		errorPage(w, err)
	}
}

func handleWeightResult(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("num") == "" {
		errorPage(w, errors.New("please enter a value"))
		return
	}
	num, err := strconv.ParseFloat(r.FormValue("num"), 64)
	if err != nil {
		errorPage(w, err)
		return
	}
	to := r.FormValue("to")
	from := r.FormValue("from")

	num *= weightInGrams[from]
	num /= weightInGrams[to]
	vals := struct {
		Num  float64
		Page string
	}{num, "weight"}
	err = resultTemplate.Execute(w, vals)
	if err != nil {
		errorPage(w, err)
	}
}

func handleTempResult(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("num") == "" {
		errorPage(w, errors.New("please enter a value"))
		return
	}
	num, err := strconv.ParseFloat(r.FormValue("num"), 64)
	if err != nil {
		errorPage(w, err)
		return
	}
	to := r.FormValue("to")
	from := r.FormValue("from")

	num = tempInCelsius[from](num, true)
	num = tempInCelsius[to](num, false)
	vals := struct {
		Num  float64
		Page string
	}{num, "temperature"}
	err = resultTemplate.Execute(w, vals)
	if err != nil {
		errorPage(w, err)
	}
}

func errorPage(w http.ResponseWriter, err error) {
	err = errorTemplate.Execute(w, err)
	if err != nil {
		log.Fatal(err)
	}
}
