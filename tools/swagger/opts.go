package swagger

import (
	"fmt"
	"strings"
)

// Option interface
type Option interface {
	GetParameter() string
}

// QueryOpt param
type QueryOpt Parameter

// GetParameter func
func (q QueryOpt) GetParameter() string {
	var queryArray = ""
	if q.Type == Array {
		queryArray = fmt.Sprintf(parameterArray, q.ItemsType)
	}
	return fmt.Sprintf(parameter, q.Name, query, q.Required, q.Description, q.Type, queryArray)
}

// PathOpt param
type PathOpt Parameter

// GetParameter func
func (q PathOpt) GetParameter() string {
	typeF := q.Type
	if typeF == "" {
		typeF = String
	}
	return fmt.Sprintf(parameter, q.Name, path, q.Required, q.Description, typeF, "")
}

type bodyOpts struct {
	BodySchema string
}

// GetParameter func
func (q bodyOpts) GetParameter() string {
	return fmt.Sprintf(parameterBody, q.BodySchema)
}

// FormDataOpt param
type FormDataOpt Parameter

// GetParameter func
func (q FormDataOpt) GetParameter() string {
	var queryArray = ""
	if q.Type == Array {
		queryArray = fmt.Sprintf(parameterArray, q.ItemsType)
	}
	return fmt.Sprintf(parameterFormData, q.Name, q.Required, q.Description, q.Type, queryArray)
}

// HeaderOpt param
type HeaderOpt Parameter

// GetParameter func
func (q HeaderOpt) GetParameter() string {
	return fmt.Sprintf(parameter, q.Name, header, q.Required, q.Description, q.Type, "")
}

// MergeOptionsJSON func
func MergeOptionsJSON(opts ...Option) ([]string, string) {
	if len(opts) == 0 {
		return []string{}, ""
	}

	var consumes []string

	var params = make([]string, 0, len(opts))
	for _, o := range opts {
		switch o.(type) {
		case FormDataOpt:
			consumes = append(consumes, fmt.Sprintf(requestConsumes, MimeMultipart))
		case bodyOpts:
			consumes = append(consumes, fmt.Sprintf(requestConsumes, MimeJson))
		}
		params = append(params, o.GetParameter())
	}

	return consumes, fmt.Sprintf(parametersRaw, strings.Join(params, ","))
}
