package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"baliance.com/gooxml/document"
	bf "github.com/russross/blackfriday"
)

type BlackGoo struct {
	d *document.Document
	TitleStyle string
	Debug bool

	curParagraph *document.Paragraph
	curRun *document.Run
	curListLevel int
	numDef document.NumberingDefinition
}

var _ bf.Renderer = &BlackGoo{}
var _ bf.Renderer = (*BlackGoo)(nil)

//func BlackGooDay() bf.Renderer {
//	return BlackGooDayWithParameters()
//}

//func BlackGooDayWithParameters() bf.Renderer {
//	return &BlackGoo{}
//}

func (b *BlackGoo) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Heading:
		if entering {
			p := b.d.AddParagraph()
			// TODO: Properly handle heading levels: node.HeadingData.Level
			p.Properties().SetStyle(b.TitleStyle)
			b.curParagraph = &p
		} else {
			b.curParagraph = nil
		}
	case bf.Paragraph:
		if entering {
			// We check this because list items seem to nest themselves in a paragraph
			if b.curParagraph == nil {
				p := b.d.AddParagraph()
				b.curParagraph = &p
			}
		} else {
			b.curParagraph = nil
		}
	case bf.Text:
		if b.curRun == nil {
			r := b.curParagraph.AddRun()
			b.curRun = &r
		}
		if node.Literal != nil {
			b.curRun.AddText(string(node.Literal))
		}
		b.curRun = nil
	case bf.Strong:
		if entering {
			if b.curRun == nil {
				r := b.curParagraph.AddRun()
				b.curRun = &r
			}
			b.curRun.Properties().SetBold(true)
		}
	case bf.Emph:
		if entering {
			if b.curRun == nil {
				r := b.curParagraph.AddRun()
				b.curRun = &r
			}
			b.curRun.Properties().SetItalic(true)
		}
	case bf.Link:
		// TODO: Improve this implementation.. hacky
		if entering && node.FirstChild != nil && node.FirstChild.Literal != nil {
			h := b.curParagraph.AddHyperLink()

			if node.LinkData.Destination != nil {
				h.SetTarget(string(node.LinkData.Destination))
			}

			if node.LinkData.Title != nil {
				h.SetToolTip(string(node.LinkData.Title))
			}

			r := h.AddRun()
			r.AddText(string(node.FirstChild.Literal))
		}
		return bf.SkipChildren
	case bf.List:
		if entering {
			b.curListLevel += 1
		} else {
			b.curListLevel -= 1
		}
	case bf.Item:
		if entering {
			p := b.d.AddParagraph()
			p.SetNumberingDefinition(b.numDef)
			p.SetNumberingLevel(b.curListLevel - 1)

			b.curParagraph = &p
		} else {
			b.curParagraph = nil
		}
	default:
		if b.Debug {
			debugNode(node, entering)
		} else {
			fmt.Errorf("error: unknown node type: %s", node.Type.String())
		}
	}
	return bf.GoToNext
}

func debugNode(node *bf.Node, entering bool) {
	child := node.FirstChild
	//var children []*bf.Node
	var children []string
	for child != nil {
		children = append(children, fmt.Sprintf("%+v", child.Type))

		if child.FirstChild != nil {
			child = child.FirstChild
		} else if child.Literal != nil {
			children = append(children, fmt.Sprintf(`"%s"`, child.Literal))
			break
		}
	}
	childrenS := strings.Join(children, "->")
	log.Printf("[debug] %+v (entering=%v, firstChildren=%s)\n", node, entering, childrenS)
}

func (b *BlackGoo) RenderHeader(w io.Writer, ast *bf.Node) {
	log.Printf("RenderHeader is not implemented yet")
}

func (b *BlackGoo) RenderFooter(w io.Writer, ast *bf.Node) {
	log.Printf("RenderFooter is not implemented yet")
}

func RunnyBlackGoo(input []byte, opts ...bf.Option) []byte {
	//renderer := BlackGooDay()
	doc := document.New()
	doc.Styles.InitializeDefault()

	// We need some numbering definitions for lists
	nd := doc.Numbering.Definitions()[0]
	renderer := &BlackGoo{
		d: doc,
		TitleStyle: "Title",
		numDef: nd,
		Debug: true,
	}

	optList := []bf.Option{bf.WithRenderer(renderer), bf.WithExtensions(bf.CommonExtensions)}
	optList = append(optList, opts...)

	// Throw the bytes away since we don't use them, we build the document up as a side effect on the renderer state
	_ = bf.Run(input, optList...)

	buffer := new(bytes.Buffer)
	renderer.d.Save(buffer)

	return buffer.Bytes()
}
