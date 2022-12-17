//go:build js && wasm && wails
// +build js,wasm,wails

package main

import (
	"strconv"

	"github.com/Nigel2392/gen/elems"
	"github.com/Nigel2392/gen/predef/navbars"
)

var WAITER = make(chan struct{})

func main() {
	var header, controls, urls = navbars.WailsFrame("", URLS)
	var style, offsetHeight = navbars.WailsFrameCss(30, "#333333", "#9200ff")
	var base = elems.Div()
	var head = base.Section()
	var body = base.Section().ID("content-body").Style("margin-top: " + strconv.Itoa(offsetHeight) + "px;")

	var rdy = urls.ActiveToggleListener(elems.SimplePaginator(body, URLS, func(err error, app *elems.Element) {
		body.WasmClearInnerHTML()
		var errElem = elems.Div().Class("alert", "alert-danger").InnerText(err.Error())
		errElem.WasmGenerate("content-body")
	}))

	header.Add(style)
	head.Add(header)
	// var loader = loader.LoaderQuadSquares("ld-1", "ldcont")
	// header.Add(loader.Style(fmt.Sprintf("top: %vpx;", offsetHeight)))
	// header.WasmGenerate("app", navbars.WailsControlListeners(controls)...)
	base.WasmGenerate("app", append(navbars.WailsControlListeners(controls), rdy)...)
	<-WAITER
}
