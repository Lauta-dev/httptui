package ui

/*
CREATE TABLE request_history (
  id TEXT PRIMARY KEY,
  url TEXT,
  method TEXT,
  status_code INTEGER,
  content_type TEXT,
  response_body TEXT,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
DROP TABLE IF EXISTS request_history;
*/

import (
	"fmt"
	colors "http_client/const/color_wrapper"
	"http_client/logic"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func statusCodesColors(code int) string {
	var color string
	switch {

	case code >= 100 && code < 200:
		color = fmt.Sprintf("[%s::b]", colors.ColorTextPrimary.String())

	case code >= 200 && code < 300:
		color = fmt.Sprintf("[%s::b]", colors.ColorSuccess.String())

	case code >= 300 && code < 400:
		color = fmt.Sprintf("[%s::b]", colors.ColorWarning.String())

	case code >= 400:
		color = fmt.Sprintf("[%s::b]", colors.ColorError.String())

	}

	return fmt.Sprintf("%s[::B]", color)

}

func History(app *tview.Application) *tview.Flex {
	list := tview.NewList()

	list.SetBorder(true)
	list.SetTitle(" > Historial de Request (F2: volver, u: actualizar lista)")
	list.ShowSecondaryText(false)

	logic.ApplySelectedBackgroundIfSupported(list, colors.ColorHighlight.TrueColor())

	flex := tview.NewFlex()
	responseView := ResponseView()

	data, err := logic.GetAllItems()

	if err != nil {
		return flex
	}

	for _, v := range data {
		code := strconv.Itoa(v.StatusCode)

		color := statusCodesColors(v.StatusCode)
		mainText := fmt.Sprintf("%s %s, %s [white]- %s", color, v.Method, code, v.URL)
		list.AddItem(mainText, v.ID, 0, nil)
	}

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		responseView.Clear()

		row, err := logic.GetItemById(secondaryText)

		if err != nil {
			fmt.Fprintf(responseView, "[red]%s", err.Error())
			return
		}

		fmt.Fprintf(responseView,
			"URL: %s\nMethod: %s\nCode: %s\nContent Type: %s\n\n%s",
			row.URL,
			row.Method,
			strconv.Itoa(row.StatusCode),
			row.ContentType,
			row.ResponseBody,
		)

	})

	flex.AddItem(list, 0, 1, true)
	flex.AddItem(responseView, 0, 1, false)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'k':
			app.SetFocus(responseView)

		case 'r':
			i := list.GetCurrentItem()
			_, id := list.GetItemText(i)

			_, err := logic.DelItems(id)

			if err != nil {
				responseView.Clear()
				fmt.Fprintf(responseView, "[red]%s", err.Error())
				return nil
			}

			list.RemoveItem(i)
			return nil

		case 'u':
			list.Clear()
			data, _ := logic.GetAllItems()
			for _, v := range data {
				code := strconv.Itoa(v.StatusCode)

				color := statusCodesColors(v.StatusCode)
				mainText := fmt.Sprintf("%s %s, %s [white]- %s", color, v.Method, code, v.URL)
				list.AddItem(mainText, v.ID, 0, nil)
			}
		}

		return event
	})

	responseView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			app.SetFocus(list)
		}

		return event
	})
	return flex
}
