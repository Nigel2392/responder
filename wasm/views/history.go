//go:build js && wasm && wails
// +build js,wasm,wails

package views

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/responder/queryvalues"

	"github.com/Nigel2392/gen/elems"
	"github.com/Nigel2392/gen/predef/messages"
	"github.com/Nigel2392/gen/wailsext"
	"golang.org/x/exp/slices"
)

func ViewHistory(body *elems.Element, args []js.Value, u *url.URL) {
	body.WasmClearInnerHTML()
	var base = body.Div().Class("container")
	var Header = base.Div().Style(
		"background-color", "#fff",
		"padding", "10px",
		"border-radius", "5px",
		"margin-bottom", "10px",
	)
	Header.H2("History").Class("text-center")
	histDiv := base.Div().ID("history")

	wailsext.GetStructure("main", "App").Call("LoadHistory").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		histDiv.WasmClearInnerHTML()
		jsondata := args[0].String()
		println(fmt.Sprintf("%v", args))
		hist := &queryvalues.History{}
		err := json.Unmarshal([]byte(jsondata), hist)
		if err != nil {
			messages.ActiveMessages.NewError("Error loading history: "+err.Error(), 5000, 150)
			return nil
		}
		slices.SortFunc(hist.Saves, func(i, j queryvalues.Save) bool {
			return i.Timestamp.After(j.Timestamp)
		})
		for _, v := range hist.Saves {
			var card = elems.Div().Class("card mb-4")
			card.Div().Class("card-header").InnerText(v.URL).TextAfter().B(v.Method + ": ")
			if len(v.QueryValues.Values) > 0 {
				var cardBody = card.Div().Class("card-body")
				var cardText = cardBody.Div().Class("card-text")
				var cardList = cardText.Ul().Class("list-group list-group-flush")
				for _, v := range v.QueryValues.Values {
					cardList.Li().Class("list-group-item").InnerText(v.Name + ": " + strings.Join(v.V, ", ")).TextAfter().B(v.Type + ": ")
				}
			}
			var data, _ = queryvalues.WailsEncodeB64(v)
			var cardFooter = card.Div().Class("card-footer")
			var anchor = elems.A().Class("btn", "btn-primary").Href("data:text/plain;charset=utf-8," + data).InnerText("Use")
			cardFooter.Div().Class("row").Add(
				elems.Div().Class("col-6").Add(
					anchor,
				),
				elems.Div().Class("col-6").Style("text-align:right;").Add(
					elems.H6().Class("text-muted").InnerText(v.Timestamp.Format("2006-01-02 15:04:05")),
				),
			)
			var _, rdyAnchor = anchor.AddEventListener("click", func(this js.Value, args []js.Value) any {
				var data = args[0].Get("target").Get("href").String()
				data = strings.TrimPrefix(data, "data:text/plain;charset=utf-8,")
				var save = &queryvalues.Save{}
				err := queryvalues.WailsDecodeB64(data, save)
				if err != nil {
					messages.ActiveMessages.NewError("Error loading history: "+err.Error(), 5000, 150)
					return nil
				}
				rqTyp = save.Method
				GenMakeRequest(body, save.QueryValues, save.URL)
				return nil
			})
			card.WasmGenerate("history", rdyAnchor)
		}
		return nil
	}))

	body.WasmGenerate("content-body")
}
