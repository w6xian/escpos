package escpos

import (
	"fmt"
	"strings"
)

type EscTr struct {
	Tds []*EscTd
}

type EscHeader struct {
	Ths []*EscTh
}

type EscTh struct {
	Key   string
	Title string
	width int
	opts  *FillOptions
}

type EscTd struct {
	Key   string
	Title string
	width int
	opts  *FillOptions
}

type EscTable struct {
	width  int
	header *EscHeader
	Trs    []*EscTr
}

func Row(tds ...TableColumn) TableRow {
	return func() *EscTr {
		return newEscRow(tds...)
	}
}

func ColumnData(title string, width int, opts ...FillOption) TableColumn {
	return func() *EscTd {
		opt := newFillOptions(opts...)
		td := &EscTd{
			Title: title,
			width: width,
			opts:  opt,
		}
		return td
	}
}

func ColumnHeader(title string, width int, opts ...FillOption) TableColumnHeader {
	return func() *EscTh {
		opt := newFillOptions(opts...)
		th := &EscTh{
			Title: title,
			width: width,
			opts:  opt,
		}
		return th
	}
}

func newEscTable(cols ...TableRow) *EscTable {
	t := &EscTable{}
	for _, c := range cols {
		tr := c()
		t.Trs = append(t.Trs, tr)
	}
	return t
}

func newEscRow(cols ...TableColumn) *EscTr {
	tr := &EscTr{}
	for _, c := range cols {
		td := c()
		tr.Tds = append(tr.Tds, td)
	}
	return tr
}

func newEscHeader(cols ...TableColumnHeader) *EscHeader {
	tr := &EscHeader{}
	for _, c := range cols {
		th := c()
		tr.Ths = append(tr.Ths, th)
	}
	return tr
}

func (t *EscTable) Print() {
	// 头
	for _, th := range t.header.Ths {
		fmt.Printf("%s", th.Title)
	}
	fmt.Println()
	for _, tr := range t.Trs {
		for _, td := range tr.Tds {
			fmt.Printf("%s", td.Title)
			fmt.Printf("%d", td.width)
		}
		fmt.Println()
	}
}

func Table(header *EscHeader, rows ...TableRow) *EscTable {
	t := newEscTable(rows...)
	t.header = header
	t.width = len(header.Ths)
	return t
}
func EscRow(cols ...TableColumn) *EscTr {
	t := newEscRow(cols...)
	return t
}

func Header(cols ...TableColumnHeader) *EscHeader {
	t := newEscHeader(cols...)
	return t
}

type TableColumnHeader func() *EscTh
type TableRow func() *EscTr
type TableColumn func() *EscTd

func (e *Escpos) PrintTable(t *EscTable) {
	e.Font(FontA)
	e.FontAlign(AlignLeft)
	e.FontSize(1, 1)
	e.FontBold(false)

	opt := newFillOptions()
	header := []string{}
	for _, th := range t.header.Ths {
		header = append(header, fillColumn(th.width, th.Title, th.opts.FillWith, th.opts.FontWidth, th.opts.Position))
	}
	headerStr := strings.Join(header, "")
	e.Println(fillColumn(e.opts.MaxChar, headerStr, opt.FillWith, opt.FontWidth, opt.Position))
	for _, tr := range t.Trs {
		row := []string{}
		trWidth := 0
		totalWidth := len(tr.Tds)
		for i, td := range tr.Tds {
			w := td.width
			// 一行打印
			if w == -2 || w == e.opts.MaxChar {
				w = e.opts.MaxChar
			}
			if w < e.opts.MaxChar {
				// 最后一列，宽度自适应
				if i == totalWidth-1 {
					w = e.opts.MaxChar - trWidth
				} else {
					trWidth += w
				}
			}
			row = append(row, fillColumn(w, td.Title, td.opts.FillWith, td.opts.FontWidth, td.opts.Position))
		}
		rowStr := strings.Join(row, "")
		e.Println(fillColumn(e.opts.MaxChar, rowStr, opt.FillWith, opt.FontWidth, opt.Position))
	}
}
