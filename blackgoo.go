package main

import (
	"io"
	"log"

	bf "github.com/russross/blackfriday"
)

type BlackGoo struct{}

var _ bf.Renderer = &BlackGoo{}
var _ bf.Renderer = (*BlackGoo)(nil)

//func BlackGooDay() bf.Renderer {
//	return BlackGooDayWithParameters()
//}

//func BlackGooDayWithParameters() bf.Renderer {
//	return &BlackGoo{}
//}

func (r *BlackGoo) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	default:
		log.Printf("The '%s' node type is not implemented yet (entering=%v)", node.Type, entering)
		//panic("Unknown node type " + node.Type.String())
	}
	return bf.GoToNext
}

func (r *BlackGoo) RenderHeader(w io.Writer, ast *bf.Node) {
	log.Printf("RenderHeader is not implemented yet")
}

func (r *BlackGoo) RenderFooter(w io.Writer, ast *bf.Node) {
	log.Printf("RenderFooter is not implemented yet")
}

func RunnyBlackGoo(input []byte, opts ...bf.Option) []byte {
	//renderer := BlackGooDay()
	renderer := &BlackGoo{}

	optList := []bf.Option{bf.WithRenderer(renderer), bf.WithExtensions(bf.CommonExtensions)}
	optList = append(optList, opts...)

	return bf.Run(input, optList...)
}
