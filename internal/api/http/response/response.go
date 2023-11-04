package response

type Response struct {
	Error string `json:"error,omitempty"`
}

func Ok() Response {
	return Response{}
}

func Error(err error) Response {
	return Response{
		Error: err.Error(),
	}
}
