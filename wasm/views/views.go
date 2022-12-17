//go:build js && wasm && wails
// +build js,wasm,wails

package views

import "syscall/js"

func ClearBody() {
	js.Global().Get("document").Call("getElementById", "content-body").Set("innerHTML", "")
}
