package spider

import (
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Request the object of each request
type Request struct {
	URL     *url.URL
	Headers *http.Header
	Method  string
	Body    io.Reader
	ID      uint32

	req    *colly.Request
	reqCtx *colly.Context

	ctx *Context
}

func newRequest(req *colly.Request, ctx *Context) *Request {
	return &Request{
		URL:     req.URL,
		Headers: req.Headers,
		Method:  req.Method,
		Body:    req.Body,
		ID:      req.ID,
		req:     req,
		reqCtx:  req.Ctx,
		ctx:     ctx,
	}
}

// Response the object of each response
type Response struct {
	StatusCode int
	Body       []byte
	Request    *Request
	Headers    *http.Header

	res *colly.Response
}

func newResponse(res *colly.Response, ctx *Context) *Response {
	return &Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Request:    newRequest(res.Request, ctx),
		Headers:    res.Headers,
		res:        res,
	}
}

// Save save the response to file
func (res *Response) Save(fileName string) error {
	return res.res.Save(fileName)
}

// FileName the filename of response
func (res *Response) FileName() string {
	return res.res.FileName()
}

// HTMLElement the html element object
type HTMLElement struct {
	Name     string
	Text     string
	Request  *Request
	Response *Response
	DOM      *goquery.Selection

	el *colly.HTMLElement
}

func newHTMLElement(el *colly.HTMLElement, ctx *Context) *HTMLElement {
	return &HTMLElement{
		Name:     el.Name,
		Text:     el.Text,
		Request:  newRequest(el.Request, ctx),
		Response: newResponse(el.Response, ctx),
		DOM:      el.DOM,
		el:       el,
	}
}

// Attr return the html element attr value
func (h *HTMLElement) Attr(k string) string {
	return h.el.Attr(k)
}

// ChildText the child text content of h
func (h *HTMLElement) ChildText(goquerySelector string) string {
	return h.el.ChildText(goquerySelector)
}

// ChildAttr the child attr value of h
func (h *HTMLElement) ChildAttr(goquerySelector, attrName string) string {
	return h.el.ChildAttr(goquerySelector, attrName)
}

// ChildAttrs the child attr list of h
func (h *HTMLElement) ChildAttrs(goquerySelector, attrName string) []string {
	return h.el.ChildAttrs(goquerySelector, attrName)
}

// ForEach calls callback on each goquerySelector element
func (h *HTMLElement) ForEach(goquerySelector string, callback func(int, *HTMLElement)) {
	cb := func(i int, el *colly.HTMLElement) {
		callback(i, newHTMLElement(el, h.Request.ctx))
	}
	h.el.ForEach(goquerySelector, cb)
}

// XMLElement the xml element object
type XMLElement struct {
	Name     string
	Text     string
	Request  *Request
	Response *Response
	DOM      interface{}

	el *colly.XMLElement
}

func newXMLElement(el *colly.XMLElement, ctx *Context) *XMLElement {
	return &XMLElement{
		Name:     el.Name,
		Text:     el.Text,
		Request:  newRequest(el.Request, ctx),
		Response: newResponse(el.Response, ctx),
		DOM:      el.DOM,
		el:       el,
	}
}

// Attr return the xml element attr value
func (x *XMLElement) Attr(k string) string {
	return x.el.Attr(k)
}

// ChildText the child text content of x
func (x *XMLElement) ChildText(xpathQuery string) string {
	return x.ChildText(xpathQuery)
}

// ChildAttr the child attr value of x
func (x *XMLElement) ChildAttr(xpathQuery, attrName string) string {
	return x.el.ChildAttr(xpathQuery, attrName)
}

// ChildAttrs the child attr list of x
func (x *XMLElement) ChildAttrs(xpathQuery, attrName string) []string {
	return x.el.ChildAttrs(xpathQuery, attrName)
}
