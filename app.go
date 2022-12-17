package main

import (
	"context"

	"github.com/Nigel2392/responder/queryvalues"
)

var Hist = queryvalues.History{Saves: make([]queryvalues.Save, 0), Fname: "history.json", MaximumItems: 200}
var Saved = queryvalues.Saves{Saves: make([]queryvalues.Save, 0), Fname: "saved.json", MaximumItems: 0}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	Hist.Save()
	// var sav = queryvalues.History(Saved)
	// sav.Save()
}

// func (a *App) PushToHistory(qv *queryvalues.QueryValues, method string, url string) {
// Hist.Add(qv, method, url)
// }
// func (a *App) PushToHistory(qv queryvalues.Save) {
// Hist.Add(qv.QueryValues, qv.Method, qv.URL)
// }
//
// func (a *App) PushToHistory(data string) {
// var s = &queryvalues.Save{}
//
// gobber := gob.NewDecoder(strings.NewReader(data))
// err := gobber.Decode(s)
// if err != nil {
// panic(err)
// }
//
// Hist.Add(s.QueryValues, s.Method, s.URL)
// }

func (a *App) PushToHistory(jsonVals string) {
	var s = &queryvalues.Save{}
	// glob
	err := queryvalues.WailsDecode(jsonVals, s)
	if err != nil {
		panic(err)
	}

	Hist.Add(s.QueryValues, s.Method, s.URL, s.RsHeaders, s.RsBody)
}

func (a *App) SaveHistory() {
	Hist.Save()
}

var histWasLoaded = false

func (a *App) LoadHistory() string {
	if !histWasLoaded {
		var s = Hist
		var s_ptr = &s
		s_ptr.Load("history.json")
		Hist = *s_ptr
		histWasLoaded = true
	}
	var data, err = queryvalues.WailsEncode(Hist)
	if err != nil {
		panic(err)
	}
	return data
}

var savedWasLoaded = false

func (a *App) LoadSaves() string {
	if !savedWasLoaded {
		var s = queryvalues.History(Saved)
		var s_ptr = &s
		s_ptr.Load("saved.json")
		Saved = queryvalues.Saves(*s_ptr)
		savedWasLoaded = true
	}
	var data, err = queryvalues.WailsEncode(Saved)
	if err != nil {
		panic(err)
	}
	return data
}

func (a *App) PushToSaved(jsonVals string) {
	var s = &queryvalues.Save{}
	// glob
	err := queryvalues.WailsDecode(jsonVals, s)
	if err != nil {
		panic(err)
	}

	var sav = queryvalues.History(Saved)
	sav.Add(s.QueryValues, s.Method, s.URL, s.RsHeaders, s.RsBody)
	Saved = queryvalues.Saves(sav)
}

func (a *App) RemoveFromSaved(jsonVals string) {
	var s = &queryvalues.Save{}
	// glob
	err := queryvalues.WailsDecode(jsonVals, s)
	if err != nil {
		panic(err)
	}

	var sav = queryvalues.History(Saved)
	sav.Remove(s.QueryValues, s.Method, s.URL, s.Timestamp)
	Saved = queryvalues.Saves(sav)
}
