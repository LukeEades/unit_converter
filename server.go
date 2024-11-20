package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

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

var distTemplate *template.Template
var resultTemplate *template.Template
var errorTemplate *template.Template
var weightTemplate *template.Template
var tempTemplate *template.Template

func main() {

	distTemplate = template.Must(template.ParseFiles("./files/templates/distance.html", "./files/templates/header.html", "./files/templates/footer.html"))
	errorTemplate = template.Must(template.ParseFiles("./files/templates/error.html", "./files/templates/header.html", "./files/templates/footer.html"))
	resultTemplate = template.Must(template.ParseFiles("./files/templates/result.html", "./files/templates/header.html", "./files/templates/footer.html"))
	weightTemplate = template.Must(template.ParseFiles("./files/templates/weight.html", "./files/templates/header.html", "./files/templates/footer.html"))
	tempTemplate = template.Must(template.ParseFiles("./files/templates/temperature.html", "./files/templates/header.html", "./files/templates/footer.html"))

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
	err := distTemplate.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleWeightForm(w http.ResponseWriter, r *http.Request) {
	err := weightTemplate.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleTempForm(w http.ResponseWriter, r *http.Request) {
	err := tempTemplate.Execute(w, nil)
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
