// The PDF package of Lexington creates a Screenplay PDF out of the Lex screenplay parsetree. This can be generated with the several other packages, e.g. the fountain package that parses fountain to lex in preparation.
package pdf

import (
	"github.com/lapingvino/lexington/lex"
	"github.com/lapingvino/lexington/rules"
	"github.com/lapingvino/lexington/font"

	"strconv"
	"github.com/phpdave11/gofpdf"
)

type Tree struct {
	PDF   *gofpdf.Fpdf
	Rules rules.Set
	F     lex.Screenplay
}

func (t Tree) pr(a string, text string) {
	line(t.PDF, t.Rules.Get(a), t.Rules.Get(a).Prefix+text+t.Rules.Get(a).Postfix)
}

func pagenumber() {

}

func (t Tree) Render() {
	var block string
	for _, row := range t.F {
		switch row.Type {
		case "newpage":
			block = ""
			t.PDF.SetHeaderFuncMode(func() {
				t.PDF.SetY(0.5)
				t.PDF.SetX(-1)
				t.PDF.Cell(0, 0, strconv.Itoa(t.PDF.PageNo())+".")
			}, true)
			t.PDF.AddPage()
			continue
		case "titlepage":
			block = "title"
			t.PDF.SetY(4)
		case "metasection":
			block = "meta"
			t.PDF.SetY(-2)
		}
		if t.Rules.Get(row.Type).Hide && block == "" {
			continue
		}
		if block != "" {
			row.Type = block
		}
		t.pr(row.Type, row.Contents)
	}
}

func line(pdf *gofpdf.Fpdf, format rules.Format, text string) {
	pdf.SetFont(format.Font, format.Style, format.Size)
	pdf.SetX(format.Left)
	pdf.MultiCell(format.Width, 0.19, text, "", format.Align, false)
}

func Create(file string, format rules.Set, contents lex.Screenplay) {
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.AddUTF8FontFromBytes("CourierPrime", "", font.MustAsset("Courier-Prime.ttf"))
	pdf.AddUTF8FontFromBytes("CourierPrime", "B", font.MustAsset("Courier-Prime-Bold.ttf"))
	pdf.AddUTF8FontFromBytes("CourierPrime", "I", font.MustAsset("Courier-Prime-Italic.ttf"))
	pdf.AddUTF8FontFromBytes("CourierPrime", "BI", font.MustAsset("Courier-Prime-Bold-Italic.ttf"))
	pdf.AddPage()
	pdf.SetMargins(1, 1, 1)
	pdf.SetXY(1, 1)
	f := Tree{
		PDF:   pdf,
		Rules: format,
		F:     contents,
	}
	f.Render()
	err := pdf.OutputFileAndClose(file)
	if err != nil {
		panic(err)
	}
}
