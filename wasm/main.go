//go:build js && wasm && wails
// +build js,wasm,wails

package main

import (
	"strconv"
	"syscall/js"

	"github.com/Nigel2392/gen/elems"
	"github.com/Nigel2392/gen/predef/navbars"
	"github.com/Nigel2392/gen/wailsext"
)

var WAITER = make(chan struct{})
var isPinnedToTop = false

func main() {
	var logoUrl = js.Global().Get("document").Call("getElementById", "LogoURL").Get("href").String()
	// var logoUrl = js.Global().Get("window").Get("LogoURL").String()
	var header, controls, urls = navbars.WailsFrame(logoUrl, URLS)
	var style, offsetHeight = navbars.WailsFrameCss(25, "#333333", "#9200ff")
	var base = elems.Div()
	var head = base.Section()
	var body = base.Section().ID("content-body").Style("margin-top: " + strconv.Itoa(offsetHeight+10) + "px;")

	var pinToWindowButton = elems.A().Class(navbars.CLASS_NAVIGATION_LINK+"-control", "remoter-pin")
	header.StyleBlock(`.remoter-pin {
		cursor: pointer;
		border-radius: 50% 50% 50% 0;
		border: 1px solid #fff;
		width: 20px;
		height: 20px !important;
		transform: rotate(-45deg);
		margin-bottom: 5px;
		margin-left: 5px;
	  }
	  .remoter-pin::after {
		position: absolute;
		content: '';
		width: 8px;
		height: 8px;
		border-radius: 50%;
		top: 50%;
		left: 50%;
		margin-left: -4px;
		margin-top: -4px;
		background-color: #fff;
	  }`)
	var _, rdyPin = pinToWindowButton.AddEventListener("click", func(this js.Value, args []js.Value) any {

		if isPinnedToTop {
			isPinnedToTop = false
			this.Get("classList").Call("remove", navbars.CLASS_ACTIVE)
		} else {
			isPinnedToTop = true
			this.Get("classList").Call("add", navbars.CLASS_ACTIVE)
		}

		wailsext.WindowSetAlwaysOnTop(isPinnedToTop)

		return nil
	})
	var left_box = header.GetByClassname("gohtml-window-control-left")[0]
	left_box.Add(pinToWindowButton)

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
	base.WasmGenerate("app", append(navbars.WailsControlListeners(controls), rdy, rdyPin)...)
	<-WAITER
}
