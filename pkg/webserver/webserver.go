package wb

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const defaultPort = "8080"

var maxRequestSize, _ = strconv.Atoi(os.Getenv("MAX_REQ_SIZE"))

type handlerFunc func(ctx *fasthttp.RequestCtx)
type middlewareFunc func(next fasthttp.RequestHandler) fasthttp.RequestHandler

// WebServer ...
type WebServer struct {
	Handlers []*Handlers
}

//Handlers ...
type Handlers struct {
	Handler    func(ctx *fasthttp.RequestCtx)
	Middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler
	Pattern    string
	HasMidd    bool
	Method     string
}

//AddHandler ....
func (w *WebServer) AddHandler(pattern, method string, handler handlerFunc) {

	exist := false

	for _, h := range w.Handlers {

		if h.Pattern == pattern && h.Method == method {
			fmt.Println("The pattern", method, ": ", pattern, "already exist therefore was dropped")
			exist = true
			break
		}
	}

	if !exist {

		w.Handlers = append(w.Handlers, &Handlers{Handler: handler, Pattern: pattern, Method: method, HasMidd: false})
	}

}

//AddHandlerWithMiddleware ....
func (w *WebServer) AddHandlerWithMiddleware(pattern, method string, handler handlerFunc, middleware middlewareFunc) {

	exist := false

	for _, h := range w.Handlers {

		if h.Pattern == pattern && h.Method == method {
			fmt.Println("The pattern", method, ": ", pattern, "already exist therefore was dropped")
			exist = true
			break
		}
	}

	if !exist {

		w.Handlers = append(w.Handlers, &Handlers{Handler: handler, Pattern: pattern, Method: method, HasMidd: true, Middleware: middleware})
	}
}

//StartUp ...
func (w *WebServer) StartUp(port string) {

	if maxRequestSize <= 1 {

		maxRequestSize = 5
	}
	if strings.TrimSpace(port) == "" {

		port = defaultPort
	}

	router := fasthttprouter.New()

	for _, handler := range w.Handlers {

		if handler.HasMidd {

			router.Handle(handler.Method, handler.Pattern, handler.Middleware(handler.Handler))

		} else {

			router.Handle(handler.Method, handler.Pattern, handler.Handler)

		}
	}

	fmt.Println("Http server running at ", port)
	fmt.Println("ctrl + C to interrupt")

	server := fasthttp.Server{
		MaxRequestBodySize: maxRequestSize * 1024 * 1024,
		Name:               "cm-http-server",
		ReadBufferSize:     4096 * 3,
		Handler:            router.Handler,
	}

	if err := server.ListenAndServe(":" + port); err != nil {

		fmt.Println("Error starting up the http server")
	}
}
