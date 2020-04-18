package http_errors

import "fmt"

type HttpError struct {
	Err  error
	Code int
}

func (hr *HttpError) Error() string {
	return fmt.Sprintf("Description %v", hr.Err.Error())
}
