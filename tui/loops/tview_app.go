package loops

import (
	"github.com/manyids2/go-tools/tui/views/layout"
	"github.com/rivo/tview"
)

func Run(ui *layout.UI) {
	app := tview.NewApplication()
	if err := app.SetRoot(ui, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
