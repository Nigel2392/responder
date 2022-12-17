package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	var middleWare assetserver.Middleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Middleware: ", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	// Create application with options
	var err = wails.Run(&options.App{
		Title:     "Remoter",
		Width:     1024,
		Height:    858,
		MinWidth:  380,
		MinHeight: 460,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: assetserver.ChainMiddleware(middleWare),
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Frameless: true,
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// Get the handler for the static/media files
func getHandler(url string, dir string, Fs fs.FS) (http.Handler, error) {
	fsdirs, err := fs.Sub(Fs, dir)
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(fsdirs)), nil
}
