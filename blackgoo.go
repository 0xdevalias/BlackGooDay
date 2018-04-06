package main

import (
	"bytes"
	"io"
	"log"

	"baliance.com/gooxml/document"
	bf "github.com/russross/blackfriday"
)

type BlackGoo struct {
	d *document.Document
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
	default:
		log.Printf("The '%s' node type is not implemented yet (entering=%v)", node.Type, entering)
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

	renderer := &BlackGoo{
		d: doc,
	}

	optList := []bf.Option{bf.WithRenderer(renderer), bf.WithExtensions(bf.CommonExtensions)}
	optList = append(optList, opts...)

	// Throw the bytes away since we don't use them, we build the document up as a side effect on the renderer state
	_ = bf.Run(input, optList...)

	buffer := new(bytes.Buffer)
	renderer.d.Save(buffer)

	return buffer.Bytes()
}
