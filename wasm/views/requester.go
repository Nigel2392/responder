//go:build js && wasm && wails
// +build js,wasm,wails

package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/responder/queryvalues"

	"github.com/Nigel2392/gen/client"
	"github.com/Nigel2392/gen/elems"
	"github.com/Nigel2392/gen/predef/loaders"
	"github.com/Nigel2392/gen/predef/messages"
	"github.com/Nigel2392/gen/wailsext"
)

type searchURL struct {
	URL string
}

// MakeRequest is the callback for the request page
func MakeRequest(body *elems.Element, args []js.Value, u *url.URL) {
	var queryList = &queryvalues.QueryValues{}
	GenMakeRequest(body, queryList, "https://httpbin.org/get")
}

const (
	queryType         = "query"
	headerType        = "header"
	formType          = "X_form_urlencoded"
	formTypeJson      = "json"
	formTypeMultipart = "multipart"
)

var (
	getType     = "GET"
	postType    = "POST"
	putType     = "PUT"
	deleteType  = "DELETE"
	headType    = "HEAD"
	patchType   = "PATCH"
	optionsType = "OPTIONS"
	//
	rqTyp   = getType
	addType = queryType
)

func GenMakeRequest(body *elems.Element, queryList *queryvalues.QueryValues, reqURL string) {
	ClearBody()
	var container = elems.Div().Class("container")
	var row = container.Div().Class("row")
	var colOne = row.Div().Class("col-12")
	var colTwo = row.Div().Class("col-6")
	var colThree = row.Div().Class("col-6")
	// var colThree = row.Div().Class("col-6")
	var queryItems = elems.Div().ID("query-items")
	var boxRow = elems.Div().Class("row").Class("mt-2")
	var boxColOne = boxRow.Div().Class("col-8")
	var boxColTwo = boxRow.Div().Class("col-4")
	var responseBox = boxColOne.Div().ID("response-box")
	var headerBox = boxColTwo.Div().ID("header-box")

	responseBox.Style(
		"border: 1px solid #ced4da",
		"border-radius: 5px",
		"padding: 10px",
		"height: 350px",
		"width: 100%",
		"overflow-y: scroll",
		"background-color: #fff",
	)

	headerBox.Style(
		"border: 1px solid #ced4da",
		"border-radius: 5px",
		"padding: 10px",
		"height: 100%",
		"width: 100%",
		"overflow-y: scroll",
		"background-color: #fff",
	)

	var form, rdyForm = makeRequestForm(queryList, responseBox, reqURL)

	var queryValueForm, rdyQV = makeQueryValueForm(queryList, queryItems)

	var typeBox = elems.Select("selectrequestdatatype", []any{
		queryType, headerType,
		formType, formTypeJson,
		formTypeMultipart,
	}, "", addType)
	typeBox.Class("form-control")
	typeBox.Style("width: 100%")
	var _, rdySelect = typeBox.AddEventListener("change", func(this js.Value, args []js.Value) interface{} {
		var selected = this.Get("value").String()
		switch selected {
		case queryType, headerType, formType, formTypeJson, formTypeMultipart:
			addType = selected
		}
		return nil
	})

	var rqtypeBox = elems.Select("selectrequesttype", []any{
		getType, postType,
		putType, deleteType,
		headType, patchType,
		optionsType,
	}, "", rqTyp)
	rqtypeBox.Class("form-control")
	rqtypeBox.Style("width: 100%")
	var _, rdyRQSelect = rqtypeBox.AddEventListener("change", func(this js.Value, args []js.Value) interface{} {
		var selected = this.Get("value").String()
		switch selected {
		case getType, postType, putType, deleteType, headerType, patchType, optionsType:
			rqTyp = selected
		}
		return nil
	})

	colOne.Add(
		elems.Div().Class("row").Add(
			elems.Div().Class("col-8").Add(
				elems.H3("Request URL"),
			),
			elems.Div().Class("col-4").Add(
				rqtypeBox,
			),
		),
		form.Class("mt-2"),
		boxRow,
	)
	colTwo.Add(
		elems.Div().Class("row").Class("mt-2").Add(
			elems.Div().Class("col-6").Add(
				elems.H3("Query Values"),
			),
			elems.Div().Class("col-6").Add(
				typeBox,
			),
		),
		queryValueForm.Class("mt-1"),
	)
	colThree.Add(
		queryItems,
	)

	container.WasmGenerate("content-body", rdyForm, rdyQV, rdySelect, rdyRQSelect) //, rdyHQV)

	addQueryItems(queryList)
}

func makeRequestForm(queryList *queryvalues.QueryValues, responseBox *elems.Element, sURL string) (*elems.Element, elems.Ready) {
	var requestSearch = searchURL{
		URL: sURL,
	}
	var form = elems.StructToForm(&requestSearch, "rq-form-label", "form-control")
	var submit = elems.Button("Submit")
	submit.A_Type("submit")
	submit.Name("action")
	submit.A_Value("submit")
	submit.Class("mt-3 mb-3 btn btn-primary")
	// var saveSubmit = elems.Button("Save")
	// saveSubmit.A_Type("submit")
	// saveSubmit.Name("action")
	// saveSubmit.A_Value("save")
	// saveSubmit.Class("mt-3 mb-3 btn btn-secondary")

	form.Add(
		elems.Div().Class("row").Add(
			elems.Div().Class("col-6").Add(
				submit,
			),
			//elems.Div().Class("col-6").Style("text-align:right").Add(
			//	saveSubmit.Style("width:50%"),
			//),
		),
	)

	var _, rdy = form.WasmFormSubmit(func(data map[string]string, jsElements []js.Value) {
		var searchURL = searchURL{}
		// Parse the form data into the searchURL struct
		elems.FormDataToStruct(data, &searchURL)
		// if data["action"] == "save" {
		// save := queryvalues.Save{
		// URL:         searchURL.URL,
		// Method:      rqTyp,
		// QueryValues: queryList,
		// }
		// Save the query values to json
		// jsondata, err := queryvalues.WailsEncode(save)
		// if err != nil {
		// return
		// }
		// Save the query values to the history
		// go wailsext.GetStructure("main", "App").Call("PushToSaved", jsondata)
		// return
		// }
		// Parse the URL
		var searchURLURL, _ = url.Parse(searchURL.URL)
		// Make the request
		var cli = client.NewAPIClient().WithLoader(
			loaders.NewLoader("", "loader", true, loaders.LoaderRotatingBlock),
		)
		println(fmt.Sprintf("Request Type: %s", rqTyp))
		switch rqTyp {
		case "GET":
			cli.Get(searchURLURL.String())
		case "POST":
			cli.Post(searchURLURL.String())
		case "PUT":
			cli.Put(searchURLURL.String())
		case "DELETE":
			cli.Delete(searchURLURL.String())
		case "PATCH":
			cli.Patch(searchURLURL.String())
		case "HEAD":
			cli.Head(searchURLURL.String())
		case "OPTIONS":
			cli.Options(searchURLURL.String())
		}
		cli.ChangeRequest(func(rq *http.Request) {
			// Set the query values
			rq.URL.RawQuery = queryList.StringByType(queryType)
			println(fmt.Sprintf("Request Queries: %s, URL: %s", queryList.StringByType(queryType), rq.URL.String()))
			var formPresentJson, formPresent, formPresentMultipart = false, false, false
			for _, qv := range queryList.Values {
				if qv.Type == headerType {
					for _, v := range qv.V {
						rq.Header.Add(qv.Name, v)
					}
				} else if qv.Type == formType || qv.Type == formTypeJson || qv.Type == formTypeMultipart {
					mapped := make(url.Values)
					for _, v := range qv.V {
						mapped[qv.Name] = append(mapped[qv.Name], v)
					}
					if qv.Type == formTypeJson {
						jsonData, _ := json.Marshal(mapped)
						rq.Body = io.NopCloser(bytes.NewReader(jsonData))
						formPresentJson = true
					} else if qv.Type == formType {
						rq.Body = io.NopCloser(bytes.NewReader([]byte(rq.Form.Encode())))
						formPresent = true
					} else if qv.Type == formTypeMultipart {
						formPresentMultipart = true
						var b bytes.Buffer
						var w = multipart.NewWriter(&b)
						for k, v := range mapped {
							for _, vv := range v {
								w.WriteField(k, vv)
							}
						}
						w.Close()
						rq.Body = io.NopCloser(&b)
					}
				}
			}
			if formPresentJson {
				rq.Header.Add("Content-Type", "application/json")
			} else if formPresent {
				rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			} else if formPresentMultipart {
				rq.Header.Add("Content-Type", "multipart/form-data")
			}

		}).OnError(func(err error) bool {
			messages.ActiveMessages.NewError(err.Error(), 5000, 150)
			return false
		}).Do(func(resp *http.Response) {
			var bodyBytes, _ = io.ReadAll(resp.Body)

			// Check if the response is JSON
			if json.Valid(bodyBytes) {
				var prettyJSON bytes.Buffer
				json.Indent(&prettyJSON, bodyBytes, "", "    ")
				bodyBytes = prettyJSON.Bytes()
			}

			var el = elems.Div().InnerHTML(string(bodyBytes)).Style("white-space: pre")
			responseBox.WasmClearInnerHTML()
			el.WasmGenerate("response-box")
			js.Global().Get("document").Call("getElementById", "header-box").Set("innerHTML", "")
			for k, v := range resp.Header {
				helem := elems.Div().InnerHTML(fmt.Sprintf("%s: %s", k, strings.Join(v, ", ")))
				helem.WasmGenerate("header-box")
			}
			save := queryvalues.Save{
				URL:         searchURL.URL,
				Method:      rqTyp,
				QueryValues: queryList,
				RsHeaders:   resp.Header,
				RsBody:      bodyBytes,
			}
			// Save the query values to json
			jsondata, err := queryvalues.WailsEncode(save)
			if err != nil {
				return
			}
			// Save the query values to the history
			go wailsext.GetStructure("main", "App").Call("PushToHistory", jsondata)
		})
	})
	return form, rdy
}

func makeQueryValueForm(queryList *queryvalues.QueryValues, queryItems *elems.Element) (*elems.Element, elems.Ready) {
	type keyVal struct {
		Key   string
		Value string
	}

	var kv = keyVal{}
	var queryValueForm = elems.StructToForm(&kv, "rq-form-label", "form-control").ID("query-value-form")
	var qvSubmit = elems.Button("Add")
	qvSubmit.A_Type("submit")
	qvSubmit.Class("mt-3 mb-3 btn btn-primary")
	queryValueForm.Add(qvSubmit)

	var _, rdyQV = queryValueForm.WasmFormSubmit(func(data map[string]string, jsElements []js.Value) {
		for _, e := range jsElements {
			e.Set("value", "")
		}
		kv := keyVal{}
		elems.FormDataToStruct(data, &kv)
		queryList.Add(kv.Key, kv.Value, addType)

		println(fmt.Sprintf("Query Value: %v, data: %v, list: %v", kv, data, queryList))

		queryItems.WasmClearInnerHTML()
		addQueryItems(queryList)
	})
	return queryValueForm, rdyQV
}

func addQueryItems(queryList *queryvalues.QueryValues) {
	var qI = js.Global().Get("document").Call("getElementById", "query-items")
	qI.Set("innerHTML", "")
	for _, qv := range queryList.Values {
		var paragraph = elems.PF("%s, %s: ", strings.ToUpper(qv.Type), qv.Name).ID(addType+"-item-"+qv.Name).Style(
			"padding: 4px",
			"background-color: #fff",
			"border-radius: 4px",
			"border: 1px solid #ced4da",
			"color: black",
			"display: inline-block",
			"width: 100%",
		).Class("mt-2")
		spans := make(elems.Elements, 0)
		for index, v := range qv.V {
			var span = elems.Span(v).Class("badge badge-secondary mr-1 query-list-item")
			span.ID(qv.HashString() + "-" + strconv.Itoa(index))
			span.Style(
				"cursor: pointer",
				"padding: 4px",
				"color: black",
			)
			spans = append(spans, span)
		}
		var _, srdy = spans.AddEventListeners("click", func(this js.Value, args []js.Value) interface{} {
			var id = this.Get("id").String()
			var parts = strings.Split(id, "-")
			var hash = parts[0]
			var index, _ = strconv.Atoi(parts[1])
			var singleRemoved, allGone = queryList.DelByHash(hash, index)

			if allGone {
				this.Get("parentElement").Call("remove")
			} else if singleRemoved {
				this.Call("remove")
				addQueryItems(queryList)
			}

			return nil
		})
		// Add the spans to the paragraph
		paragraph.Add(spans...)
		// Generate the paragraph
		paragraph.WasmGenerate("query-items", srdy)
	}
}
