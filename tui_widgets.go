package main

import (
	ui "github.com/gizak/termui"
)

type CustomWidget struct {
	ui.Block
	Data string
}

func NewCustomWidget(data string) *CustomWidget {
	cw := &CustomWidget{
		Block: *ui.NewBlock(),
		Data:  data,
	}
	cw.BorderLabel = " Messages "
	return cw
}

func (cw *CustomWidget) Buffer() ui.Buffer {
	buf := cw.Block.Buffer()
	for i, c := range cw.Data {
		// NewCell takes a rune and 2 Attributes
		cell := ui.Cell{c, ui.ColorBlue, ui.ColorDefault}

		// Set takes x, y, and a Cell
		buf.Set(i+5, 2, cell)
	}

	return buf
}

// func main() {
// 	err := ui.Init()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer ui.Close()

// 	cw := NewCustomWidget("Hello World!")
// 	cw.Width = 25
// 	cw.Height = 5
// 	ui.Render(cw)

// 	uiEvents := ui.PollEvents()
// 	for {
// 		e := <-uiEvents
// 		if e.Type == ui.KeyboardEvent {
// 			return
// 		}
// 	}
// }
