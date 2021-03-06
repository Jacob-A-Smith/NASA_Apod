package main

import (
	"net/url"
	"strconv"

	apodRequester "jsmith/nasa/apodRequester"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()

	w := app.NewWindow("Hello")
	w.Resize(fyne.NewSize(1920, 1080))

	titleLabel := widget.NewLabel("NASA Astronomy Picture of the Day")
	numDaysEntry := widget.NewEntry()
	numDaysEntry.SetPlaceHolder("Number of days...")
	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("api key")
	exitApplication := widget.NewButton("Quit", func() {
		app.Quit()
	})

	var update func([]apodRequester.ApodResponse)
	var showMainMenu func()
	makeForm := func(apod apodRequester.ApodResponse) *widget.Form {
		wid := widget.NewForm(widget.NewFormItem(apod.Date, widget.NewButton(apod.Title, func() {
			img, err := fyne.LoadResourceFromURLString(apod.Hdurl)
			if err != nil {
				u, e := url.Parse(apod.URL)
				if e != nil {
					return
				}
				resource := widget.NewHyperlink(apod.URL, u)
				picWindow := app.NewWindow(apod.Title)
				picWindow.Resize(fyne.NewSize(512, 512))
				picWindow.SetContent(resource)
				picWindow.Show()
				return
			}
			resource := canvas.NewImageFromResource(img)
			picWindow := app.NewWindow(apod.Title)
			picWindow.Resize(fyne.NewSize(1920, 1080))
			picWindow.SetContent(resource)
			picWindow.Show()
		})))
		return wid
	}

	update = func(apods []apodRequester.ApodResponse) {
		bx := widget.NewGroupWithScroller("Astromony Pictures of the Day")
		for _, v := range apods {
			bx.Append(makeForm(v))
		}
		bx.Prepend(widget.NewButton("<- Back to main menu", showMainMenu))
		w.SetContent(bx)
	}

	showMainMenu = func() {
		w.SetContent(widget.NewVBox(
			titleLabel,
			widget.NewForm(widget.NewFormItem("Enter API key: ", apiKeyEntry)),
			widget.NewForm(widget.NewFormItem("Enter number of days: ", numDaysEntry)),
			widget.NewButton("Request photos", func() {
				if len(apiKeyEntry.Text) == 0 {
					return
				}
				num, e := strconv.ParseInt(numDaysEntry.Text, 10, 32)
				if e != nil {
					numDaysEntry.SetText("")
					return
				}
				go apodRequester.GetApodForDateRange(int(num), update, apiKeyEntry.Text)
				w.SetContent(widget.NewVBox(
					titleLabel,
					widget.NewLabel("Processing Request..."),
					exitApplication,
				))
			}),
			exitApplication,
		))
	}

	showMainMenu()
	w.ShowAndRun()
}
