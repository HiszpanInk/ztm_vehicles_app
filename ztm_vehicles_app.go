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

func getVehicleByNum(vehicleNum int) vehicle {
	var retrievedData [10]string

	//get data from website and insert it into array
	vehicleURL := fmt.Sprintf("https://www.ztm.waw.pl/baza-danych-pojazdow/?ztm_mode=2&ztm_vehicle=%d", vehicleNum)
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

	return retrievedVehicle
}

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	output := widget.NewLabel("Numer taborowy pojazdu: ")
	form := widget.NewForm(
		widget.NewFormItem("search_number", widget.NewEntry()),
	)
	form.OnSubmit = func() {
		output.Text = "submitted"
		output.Refresh()
	}

	w.SetContent(container.NewVBox(
		output,
		form,
	))

	w.ShowAndRun()
}
