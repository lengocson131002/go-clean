package validation

// interface for validation
type Validator interface {
	Validate(i interface{}) error
}
