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
		if entering &&
			node.FirstChild != nil &&
			node.FirstChild.Type == bf.Text &&
			node.FirstChild.Literal != nil {
			p := b.d.AddParagraph()
			p.Properties().SetStyle(b.TitleStyle)

			r := p.AddRun()
			r.AddText(string(node.FirstChild.Literal))
		}
		return bf.SkipChildren
	case bf.Paragraph:
		if entering {
			p := b.d.AddParagraph()
			b.curParagraph = &p
		} else {
			b.curParagraph = nil
		}
	default:
		if b.Debug {
			child := node.FirstChild
			//var children []*bf.Node
			var children []string
			for child != nil {
				children = append(children, fmt.Sprintf("%+v", child.Type))
				child = child.FirstChild
			}
			childrenS := strings.Join(children, "->")
			log.Printf("[debug] %+v (entering=%v, firstChildren=%s)\n", node, entering, childrenS)
		}
		//log.Printf("The '%s' node type is not implemented yet (entering=%v)", node.Type, entering)
		//panic("Unknown node type " + node.Type.String())
	}
	return bf.GoToNext
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

	renderer := &BlackGoo{
		d: doc,
		TitleStyle: "Title",
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
