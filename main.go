package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type month struct {
	Name   string
	Plants []plant
}

type plant struct {
	Name        string
	Description string
	Months      []int
}

type fullData struct {
	Months []month
	Plants []plant
}

var tpl *template.Template
var indexData fullData
var months []month
var plants []plant
var fm = template.FuncMap{
	"cm": checkMonth,
	"ft": firstThree,
	"gm": getMonths,
}

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*"))
	months = []month{
		month{Name: "January"},
		month{Name: "February"},
		month{Name: "March"},
		month{Name: "April"},
		month{Name: "May"},
		month{Name: "Jun"},
		month{Name: "July"},
		month{Name: "August"},
		month{Name: "September"},
		month{Name: "October"},
		month{Name: "November"},
		month{Name: "December"},
	}

	plants = []plant{
		plant{Name: "Asparagus", Description: "Green", Months: []int{0, 11}},
		plant{Name: "Beets", Description: "Red", Months: []int{1, 2, 3, 7, 8, 9}},
		plant{Name: "Broccoli", Description: "Green and tree-like", Months: []int{0, 1, 6, 7, 8, 9, 10, 11}},
		plant{Name: "Cabbage", Description: "Leafy", Months: []int{0, 1, 6, 7, 8, 9, 10}},
		plant{Name: "Beets", Description: "Red", Months: []int{1, 2, 3, 7, 8, 9}},
		plant{Name: "Carrots", Description: "Orange", Months: []int{2, 3, 4, 7, 8, 9, 10, 11}},
		plant{Name: "Cucumbers", Description: "Tubular", Months: []int{3, 4, 5}},
		plant{Name: "Onions", Description: "Like an ogre", Months: []int{7, 8, 9, 10, 11}},
		plant{Name: "Peas", Description: "Bumpy and green", Months: []int{0, 1, 2, 8, 9, 10, 11}},
		plant{Name: "Tomatoes", Description: "Killer", Months: []int{0, 1, 2, 3, 4, 5}},
	}

	assignPlants()

	indexData = fullData{
		Months: months,
		Plants: plants,
	}
}

func firstThree(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 3 {
		s = s[:3]
	}
	return s
}

func checkMonth(index int, plantMonths []int) bool {
	for _, m := range plantMonths {
		if index == m {
			return true
		}
	}
	return false
}

func assignPlants() {
	for _, p := range plants {
		for _, n := range p.Months {
			months[n].Plants = append(months[n].Plants, p)
		}
	}
}

func getMonths() []month {
	return months
}

func main() {
	mux := httprouter.New()
	mux.GET("/", index)
	mux.GET("/plant_detail/:plant", plantDetail)
	mux.GET("/month_detail/:month", monthDetail)

	http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", indexData)
	HandleError(w, err)
}

func plantDetail(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	for _, p := range plants {
		if ps.ByName("plant") == p.Name {
			fmt.Println(p)
			err := tpl.ExecuteTemplate(w, "plant_detail.gohtml", p)
			HandleError(w, err)
		}
	}
}

func monthDetail(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	for _, m := range months {
		if ps.ByName("month") == m.Name {
			err := tpl.ExecuteTemplate(w, "month_detail.gohtml", m)
			HandleError(w, err)
		}
	}
}

// HandleError send error through http and log it out.
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
