package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

const appName = "Gostatusd"

//go:embed web.zip
var staticZip []byte

var listenAddr string
var recordFile string
var themeZip string

func main() {
	flag.StringVar(&listenAddr, "listen", ":9900", "listen address")
	flag.StringVar(&recordFile, "record", filepath.Join(os.TempDir(), "gostatusd.json"), "monthly traffic record file")
	flag.StringVar(&themeZip, "theme", "", "custom theme in zip file")
	flag.Parse()
	// load theme
	if len(themeZip) > 0 {
		if zipdata, err := ioutil.ReadFile(themeZip); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		} else {
			staticZip = zipdata
		}
	}
	themeData, err := zip.NewReader(bytes.NewReader(staticZip), int64(len(staticZip)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	// init
	if err := initWorker(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	go backgroundWorker()

	// run web ui
	app := fiber.New(fiber.Config{
		ServerHeader: appName,
		AppName:      appName,
	})
	app.Get("/stat", getNetStatHandler)
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(themeData),
	}))
	app.Listen(listenAddr)
}

func getNetStatHandler(c *fiber.Ctx) error {
	systemInfo, err := getSystemInfo()
	if err != nil {
		return err
	}
	return c.JSON(systemInfo)
}
