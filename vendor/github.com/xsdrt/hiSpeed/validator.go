package hiSpeed

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

func (h *HiSpeed) Validator(data url.Values) *Validation {
	return &Validation{
		Errors: make(map[string]string),
		Data:   data,
	}
}

func (v *Validation) Valid() bool {
	return len(v.Errors) == 0 //So if no errors than returns true otherwise if an error has been mapped returns false...
}

func (v *Validation) AddError(key, message string) { //Add errors to the map...
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validation) Has(field string, r *http.Request) bool { //true if in the request , false if it is not...
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true //Otherwise return true
}

func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "This field cannot be a blank or empty field")
		}
	}
}

func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "Invalid email address")
	}
}

func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "This field must be an integer(whole number)")
	}
}

func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "This field must be a floating point number")
	}
}

func (v *Validation) IsDateISO(field, value string) {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		v.AddError(field, "This field must be a date in the form of YYYY-MM-DD")
	}
}

func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "Spaces are not permitted in this field")
	}
}
