package forms

// Errors type to represent form field errors
type Errors map[string][]string

// NewErrors creates a new Errors struct
func NewErrors() Errors {
	return make(Errors)
}

// Add adds an error message for a given form field
func (e Errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns first error message for a field
func (e Errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0] // return the first error message
}

// Has checks if a field has any errors
func (e Errors) Has(field string) bool {
	return len(e[field]) > 0
}
