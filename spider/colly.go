package spider

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Request struct {
	URL     *url.URL
	Headers *http.Header
	Method  string
	Body    io.Reader

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
		req:     req,
		reqCtx:  req.Ctx,
		ctx:     ctx,
	}
}

func (r *Request) PutReqContextValue(key string, value interface{}) {
	r.reqCtx.Put(key, value)
}

func (r *Request) GetReqContextValue(key string) string {
	return r.reqCtx.Get(key)
}

func (r *Request) GetAnyReqContextValue(key string) interface{} {
	return r.reqCtx.GetAny(key)
}

func (r *Request) SetResponseCharacterEncoding(encoding string) {
	r.req.ResponseCharacterEncoding = encoding
}

func (r *Request) Abort() {
	r.req.Abort()
}

func (r *Request) reqContextClone() *colly.Context {
	newCtx := colly.NewContext()
	r.reqCtx.ForEach(func(k string, v interface{}) interface{} {
		newCtx.Put(k, v)
		return nil
	})

	return newCtx
}

func (r *Request) AbsoluteURL(u string) string {
	return r.req.AbsoluteURL(u)
}

func (r *Request) Visit(URL string) error {
	return r.req.Visit(URL)
}

func (r *Request) VisitForNext(URL string) error {
	return r.ctx.VisitForNext(r.AbsoluteURL(URL))
}

func (r *Request) VisitForNextWithContext(URL string) error {
	return r.ctx.nextC.Request("GET", r.req.AbsoluteURL(URL), nil, r.reqContextClone(), nil)
}

func (r *Request) Post(URL string, requestData map[string]string) error {
	return r.req.Post(URL, requestData)
}

func (r *Request) PostForNext(URL string, requestData map[string]string) error {
	return r.ctx.PostForNext(r.AbsoluteURL(URL), requestData)
}

func (r *Request) PostForNextWithContext(URL string, requestData map[string]string) error {
	return r.ctx.nextC.Request("POST", r.req.AbsoluteURL(URL), createFormReader(requestData), r.reqContextClone(), nil)
}

func (r *Request) PostRaw(URL string, requestData []byte) error {
	return r.req.PostRaw(URL, requestData)
}

func (r *Request) PostRawForNext(URL string, requestData []byte) error {
	return r.ctx.PostRawForNext(r.AbsoluteURL(URL), requestData)
}

func (r *Request) PostRawForNextWithContext(URL string, requestData []byte) error {
	return r.ctx.nextC.Request("POST", r.req.AbsoluteURL(URL), bytes.NewReader(requestData), r.reqContextClone(), nil)
}

func (r *Request) PostMultipart(URL string, requestData map[string][]byte) error {
	return r.req.PostMultipart(URL, requestData)
}

func (r *Request) PostMultipartForNext(URL string, requestData map[string][]byte) error {
	return r.ctx.PostMultipartForNext(r.AbsoluteURL(URL), requestData)
}

func (r *Request) Retry() error {
	return r.req.Retry()
}

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

func (res *Response) Save(fileName string) error {
	return res.res.Save(fileName)
}

func (res *Response) FileName() string {
	return res.res.FileName()
}

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

func (h *HTMLElement) Attr(k string) string {
	return h.el.Attr(k)
}

func (h *HTMLElement) ChildText(goquerySelector string) string {
	return h.el.ChildText(goquerySelector)
}

func (h *HTMLElement) ChildAttr(goquerySelector, attrName string) string {
	return h.el.ChildAttr(goquerySelector, attrName)
}

func (h *HTMLElement) ChildAttrs(goquerySelector, attrName string) []string {
	return h.el.ChildAttrs(goquerySelector, attrName)
}

func (h *HTMLElement) ForEach(goquerySelector string, callback func(int, *HTMLElement)) {
	cb := func(i int, el *colly.HTMLElement) {
		callback(i, newHTMLElement(el, h.Request.ctx))
	}
	h.el.ForEach(goquerySelector, cb)
}

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
	}
}

func (x *XMLElement) Attr(k string) string {
	return x.el.Attr(k)
}

func (x *XMLElement) ChildText(xpathQuery string) string {
	return x.ChildText(xpathQuery)
}

func (x *XMLElement) ChildAttr(xpathQuery, attrName string) string {
	return x.el.ChildAttr(xpathQuery, attrName)
}

func (x *XMLElement) ChildAttrs(xpathQuery, attrName string) []string {
	return x.el.ChildAttrs(xpathQuery, attrName)
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}
