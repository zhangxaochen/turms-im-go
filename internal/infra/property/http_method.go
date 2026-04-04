package property

type HttpMethod string

const (
	HttpMethod_GET     HttpMethod = "GET"
	HttpMethod_POST    HttpMethod = "POST"
	HttpMethod_PUT     HttpMethod = "PUT"
	HttpMethod_DELETE  HttpMethod = "DELETE"
	HttpMethod_HEAD    HttpMethod = "HEAD"
	HttpMethod_PATCH   HttpMethod = "PATCH"
	HttpMethod_OPTIONS HttpMethod = "OPTIONS"
	HttpMethod_TRACE   HttpMethod = "TRACE"
	HttpMethod_CONNECT HttpMethod = "CONNECT"
)
