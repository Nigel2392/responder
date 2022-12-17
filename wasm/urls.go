//go:build js && wasm && wails
// +build js,wasm,wails

package main

import (
	"github.com/Nigel2392/responder/wasm/views"

	"github.com/Nigel2392/gen/elems"
)

var URLS = elems.URLs{
	elems.URL{
		URL:         "/request",
		Name:        "Request",
		LeftOrRight: false,
		CallBack:    views.MakeRequest,
	},
	elems.URL{
		URL:         "/history",
		Name:        "History",
		LeftOrRight: false,
		CallBack:    views.ViewHistory,
	},
	//elems.URL{
	//	URL:         "/saved",
	//	Name:        "Saved",
	//	LeftOrRight: false,
	//	CallBack:    views.ViewSaves,
	//},
}
