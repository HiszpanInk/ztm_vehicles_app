package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly/v2"
)

type vehicle struct {
	producer                   string
	model                      string
	production_year            string
	traction_type              string
	vehicle_registration_plate string
	vehicle_number             string
	operator                   string
	garage                     string
	ticket_machine             string
	equipment                  string
}

func getVehicleByNum(vehicleNum string) string {
	searchURL := fmt.Sprintf("https://www.ztm.waw.pl/baza-danych-pojazdow/?ztm_traction=&ztm_make=&ztm_model=&ztm_year=&ztm_registration=&ztm_vehicle_number=%s&ztm_carrier=&ztm_depot=", vehicleNum)
	vehicleURL := ""
	c2 := colly.NewCollector(
		// Visit only domains:
		colly.AllowedDomains("www.ztm.waw.pl"),
	)
	c2.OnHTML(".grid-row-active", func(e *colly.HTMLElement) {
		text := e.Attr("href")
		vehicleURL = text
	})
	c2.Visit(searchURL)
	if searchURL == "" {
		return ""
	} else {
		var retrievedData [10]string
		// Instantiate default collector
		c := colly.NewCollector(
			// Visit only domains:
			colly.AllowedDomains("www.ztm.waw.pl"),
		)
		dataIndex := 0
		c.OnHTML(".vehicle-details-entry-value", func(e *colly.HTMLElement) {
			text := e.Text
			retrievedData[dataIndex] = text
			dataIndex++
		})
		c.Visit(vehicleURL)
		retrievedVehicle := vehicle{
			producer:                   retrievedData[0],
			model:                      retrievedData[1],
			production_year:            retrievedData[2],
			traction_type:              retrievedData[3],
			vehicle_registration_plate: retrievedData[4],
			vehicle_number:             retrievedData[5],
			operator:                   retrievedData[6],
			garage:                     retrievedData[7],
			ticket_machine:             retrievedData[8],
			equipment:                  retrievedData[9],
		}

		output_string := fmt.Sprintf(
			`%s %s
		 	z roku %s,
	 	 	w posiadaniu %s,
		 	z zajezdni %s,
		 	o rejestracji %s`, retrievedVehicle.producer, retrievedVehicle.model, retrievedVehicle.production_year, retrievedVehicle.operator, retrievedVehicle.garage, retrievedVehicle.vehicle_registration_plate)
		return output_string
	}

}

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	output := widget.NewLabel("")
	entry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Podaj numer taborowy pojazdu:", Widget: entry}},
	}
	form.OnSubmit = func() {
		output_data := getVehicleByNum(entry.Text)
		if output_data != "" {
			output.Text = output_data
		} else {
			output.Text = "Nie znaleziono pojazdu o podanym numerze taborowym w bazie pojazdów WTP"
		}
		output.Refresh()
	}
	clearButton := widget.NewButton("Wyczyść", func() {
		entry.Text = ""
		output.Text = ""
	})
	w.SetContent(container.NewVBox(
		form,
		clearButton,
		output,
	))

	w.ShowAndRun()
}
