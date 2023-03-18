package main

import (
	"fmt"

	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"

	"strconv"

	"fyne.io/fyne/v2/driver/mobile"
)

type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *numericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

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

func vehicleToString(inputVehicle *vehicle) string {
	registrationPlate := ""
	if len(inputVehicle.vehicle_registration_plate) > 0 {
		registrationPlate = fmt.Sprintf(`,o rejestracji %s`, inputVehicle.vehicle_registration_plate)
	}
	output_string := fmt.Sprintf(
		"%s o numerze %s,"+"%s %s,"+"z roku %s,"+"w posiadaniu %s,"+"z zajezdni %s%s",
		inputVehicle.traction_type,
		inputVehicle.vehicle_number,
		inputVehicle.producer,
		inputVehicle.model,
		inputVehicle.production_year,
		inputVehicle.operator,
		inputVehicle.garage,
		registrationPlate)
	return strings.Replace(output_string, ",", ",\n", -1)
}

func getVehiclesByNum(vehicleNum string) string {
	searchURL := fmt.Sprintf("https://www.ztm.waw.pl/baza-danych-pojazdow/?ztm_traction=&ztm_make=&ztm_model=&ztm_year=&ztm_registration=&ztm_vehicle_number=%s&ztm_carrier=&ztm_depot=", vehicleNum)
	var vehicleURLs []string

	c2 := colly.NewCollector(
		// Visit only domains:
		colly.AllowedDomains("www.ztm.waw.pl"),
	)
	c2.OnHTML(".grid-row-active", func(e *colly.HTMLElement) {
		text := e.Attr("href")
		vehicleURLs = append(vehicleURLs, text)
	})
	c2.Visit(searchURL)
	output_string := ""
	for _, vehicleURL := range vehicleURLs {
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
		output_string += vehicleToString(&retrievedVehicle)
		output_string += "\n"
		output_string += "\n"
	}
	return output_string
}

func main() {
	a := app.New()
	w := a.NewWindow("ZTM vehicles")

	entryLabel := widget.NewLabel("Podaj numer taborowy pojazdu:")
	entry := newNumericalEntry()
	output := widget.NewLabel("")
	output.Wrapping = fyne.TextWrapWord

	//outputContainer := container.NewVScroll(output)
	//I create button earlier in order to reference it later inside function of it so I can easily disable it for the time data is fetching. There is probably better way to do it but I could't find it so
	executeButton := widget.NewButton("Sprawdź numer", func() {
		output.Refresh()
	})
	executeButton = widget.NewButton("Sprawdź numer", func() {
		executeButton.Disable()
		executeButton.SetText("Wczytywanie")
		if len(strings.TrimSpace(entry.Text)) != 0 {
			output_data := getVehiclesByNum(entry.Text)
			if output_data != "" {
				output.Text = output_data
			} else {
				output.Text = "Nie znaleziono pojazdu o podanym numerze taborowym w bazie pojazdów WTP"
			}
		} else {
			output.Text = "Nie wprowadzono prawidłowego numeru"
		}

		output.Refresh()
		executeButton.Enable()
		executeButton.SetText("Sprawdź numer")
	})
	executeButton.Importance = widget.HighImportance

	clearButton := widget.NewButton("Wyczyść", func() {
		entry.Text = ""
		output.Text = ""
		output.Refresh()
		entry.Refresh()
	})
	buttons := container.New(layout.NewGridLayout(2), clearButton, executeButton)
	inputs := container.New(layout.NewGridLayout(1), entryLabel, entry, buttons)

	w.SetContent(container.NewVBox(
		inputs,
		output,
	))

	w.ShowAndRun()
}
