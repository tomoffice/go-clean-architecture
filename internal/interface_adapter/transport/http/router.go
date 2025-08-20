package http

type Router interface {
	GET(path string, h HandlerFunc)
	POST(path string, h HandlerFunc)
	PUT(path string, h HandlerFunc)
	PATCH(path string, h HandlerFunc)
	DELETE(path string, h HandlerFunc)
}

type HandlerFunc func(Context)
