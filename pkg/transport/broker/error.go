package broker

import "fmt"

type EmptyRequestError struct{}

func (e EmptyRequestError) Error() string {
	return fmt.Sprintf("Empty broker request")
}

type InvalidDataFormatError struct{}

func (e InvalidDataFormatError) Error() string {
	return fmt.Sprintf("Invalid data format")
}
