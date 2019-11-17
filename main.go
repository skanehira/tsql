package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

func saveBuffer(b *femto.Buffer, path string) error {
	return ioutil.WriteFile(path, []byte(b.String()), 0600)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: femto [filename]\n")
		os.Exit(1)
	}
	path := os.Args[1]

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("could not read %v: %v", path, err)
	}

	var colorscheme femto.Colorscheme
	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}

	app := tview.NewApplication()

	buffer := femto.NewBufferFromString(string(content), path)
	editor := femto.NewView(buffer)
	editor.SetRuntimeFiles(runtime.Files)
	editor.SetColorscheme(colorscheme)
	editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			saveBuffer(buffer, path)
			return nil
		case tcell.KeyCtrlQ:
			app.Stop()
			return nil
		}
		return event
	})

	editor.SetBorder(true).SetTitle("editor")

	grid := tview.NewGrid().SetColumns(0, 0).
		AddItem(editor, 0, 0, 1, 1, 0, 0, true)

	app.SetRoot(grid, true)

	if err := app.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
