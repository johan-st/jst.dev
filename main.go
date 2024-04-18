package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    // serves static files from the provided public dir (if exists)
    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        // Serve the public directory falling back to "index.html" if not found
        e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), true))

        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}