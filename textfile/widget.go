package textfile

import (
	//"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	//"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	//"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/wtf/wtf"
)

const HelpText = `
  Keyboard commands for Textfile:

    /: Show/hide this help window
    o: Open the text file in the operating system
`

type Widget struct {
	wtf.HelpfulWidget
	wtf.TextWidget

	filePath string
}

func NewWidget(app *tview.Application, pages *tview.Pages) *Widget {
	widget := Widget{
		HelpfulWidget: wtf.NewHelpfulWidget(app, pages, HelpText),
		TextWidget:    wtf.NewTextWidget("TextFile", "textfile", true),

		filePath: wtf.Config.UString("wtf.mods.textfile.filePath"),
	}

	widget.HelpfulWidget.SetView(widget.View)

	widget.View.SetWrap(true)
	widget.View.SetWordWrap(true)
	widget.View.SetInputCapture(widget.keyboardIntercept)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.UpdateRefreshedAt()
	widget.View.SetTitle(widget.ContextualTitle(widget.fileName()))

	filePath, _ := wtf.ExpandHomeDir(widget.filePath)
	file, err := os.Open(filePath)

	if err != nil {
		widget.View.SetText(err.Error())
	} else {
		lexer := lexers.Match(filePath)
		if lexer == nil {
			lexer = lexers.Fallback
		}

		style := styles.Get(wtf.Config.UString("wtf.mods.textfile.formatStyle", "vim"))
		if style == nil {
			style = styles.Fallback
		}
		formatter := formatters.Get("terminal256")
		if formatter == nil {
			formatter = formatters.Fallback
		}

		contents, _ := ioutil.ReadAll(file)
		iterator, _ := lexer.Tokenise(nil, string(contents))

		var buf bytes.Buffer
		formatter.Format(&buf, style, iterator)

		formatted := tview.TranslateANSI(buf.String())
		widget.View.SetText(formatted)
	}
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) fileName() string {
	return filepath.Base(widget.filePath)
}

func (widget *Widget) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	switch string(event.Rune()) {
	case "/":
		widget.ShowHelp()
		return nil
	case "o":
		wtf.OpenFile(widget.filePath)
		return nil
	}

	return event
}
