package escpos

import (
	"fmt"
	"strings"
)

type Tr struct {
	Tds []*Td
}

type Header struct {
	Ths []*Th
}

type Th struct {
	Title string
	width int
	opts  *FillOptions
}

type Td struct {
	Title string
	width int
	opts  *FillOptions
}

type Table struct {
	width  int
	header *Header
	Trs    []*Tr
}

func Row(tds ...TableColumn) TableRow {
	return func() *Tr {
		return newRow(tds...)
	}
}

func Column(title string, width int, opts ...FillOption) TableColumn {
	return func() *Td {
		opt := newFillOptions(opts...)
		td := &Td{
			Title: title,
			width: width,
			opts:  opt,
		}
		return td
	}
}

func HeaderColumn(title string, width int, opts ...FillOption) TableColumnHeader {
	return func() *Th {
		opt := newFillOptions(opts...)
		th := &Th{
			Title: title,
			width: width,
			opts:  opt,
		}
		return th
	}
}

func newTable(cols ...TableRow) *Table {
	t := &Table{}
	for _, c := range cols {
		tr := c()
		t.Trs = append(t.Trs, tr)
	}
	return t
}

func newRow(cols ...TableColumn) *Tr {
	tr := &Tr{}
	for _, c := range cols {
		td := c()
		tr.Tds = append(tr.Tds, td)
	}
	return tr
}

func newHeader(cols ...TableColumnHeader) *Header {
	tr := &Header{}
	for _, c := range cols {
		th := c()
		tr.Ths = append(tr.Ths, th)
	}
	return tr
}

func (t *Table) Print() {
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

func EscTable(header *Header, rows ...TableRow) *Table {
	t := newTable(rows...)
	t.header = header
	t.width = len(header.Ths)
	return t
}
func EscRow(cols ...TableColumn) *Tr {
	t := newRow(cols...)
	return t
}

func EscHeader(cols ...TableColumnHeader) *Header {
	t := newHeader(cols...)
	return t
}

type TableColumnHeader func() *Th
type TableRow func() *Tr
type TableColumn func() *Td

func (e *Escpos) PrintTable(t *Table) {
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
		for _, td := range tr.Tds {
			row = append(row, fillColumn(td.width, td.Title, td.opts.FillWith, td.opts.FontWidth, td.opts.Position))
		}
		rowStr := strings.Join(row, "")
		e.Println(fillColumn(e.opts.MaxChar, rowStr, opt.FillWith, opt.FontWidth, opt.Position))
	}
}
