# webserver

## Usage
- install the webserver dependency ```go go get -u gitlab.com/mainmast/microservices/cm-http.git```
-  import the webserber module `wb "gitlab.com/mainmast/cm-http.git/pkg/webserver"`
- initialize the webserver in this way

```go

	server := &wb.WebServer{}

	server.AddHandler("/rating", "POST", handler.MyHandler) // you can use a simple handler
	server.AddHandlerWithMiddleware("/rating", "POST", handler.MyHandler, handler.MyMiddleware) // or you can use a handler with a middleware
	server.StartUp("8080")

```
