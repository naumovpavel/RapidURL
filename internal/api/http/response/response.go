package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
