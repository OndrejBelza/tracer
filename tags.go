package tracer

type Tag string

// Common Tag values
const (
	TagHTTPMethod       Tag = "http.method"
	TagHTTPPath         Tag = "http.path"
	TagHTTPUrl          Tag = "http.url"
	TagHTTPRoute        Tag = "http.route"
	TagHTTPStatusCode   Tag = "http.status_code"
	TagHTTPRequestSize  Tag = "http.request.size"
	TagHTTPResponseSize Tag = "http.response.size"
	TagGRPCStatusCode   Tag = "grpc.status_code"
	TagSQLQuery         Tag = "sql.query"
	TagError            Tag = "error"
)
