package helper

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	FieldErrors map[string]string
}

func (e *ValidationError) Error() string {
	return "Validasi gagal"
}

func ValidateStruct(s any) *ValidationError {
	if err := validate.Struct(s); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			fieldErrors := make(map[string]string)
			t := reflect.TypeOf(s)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			for _, v := range verr {
				key := v.Field()
				if field, ok := t.FieldByName(v.Field()); ok {
					if tag := field.Tag.Get("json"); tag != "" {
						key = strings.Split(tag, ",")[0]
					}
				}
				fieldErrors[key] = translateFieldError(v)
			}
			return &ValidationError{FieldErrors: fieldErrors}
		}
	}
	return nil
}

func translateFieldError(v validator.FieldError) string {
	label := fieldLabel(v.Field())
	switch v.Tag() {
	case "required":
		if v.Field() == "Title" {
			return "Judul tugas tidak boleh kosong"
		}
		if v.Field() == "IDPriorities" {
			return "Prioritas harus dipilih antara skala 1 sampai 4"
		}
		return label + " wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "numeric":
		return label + " harus berupa angka"
	case "len":
		return label + " harus " + v.Param() + " karakter"
	case "min", "max":
		if v.Field() == "IDPriorities" {
			return "Prioritas harus dipilih antara skala 1 sampai 4"
		}
		if v.Tag() == "min" {
			return label + " minimal " + v.Param() + " karakter"
		}
		return label + " maksimal " + v.Param() + " karakter"
	default:
		return label + " tidak valid"
	}
}

func fieldLabel(field string) string {
	labels := map[string]string{
		"Nama":        "Nama",
		"Email":       "Email",
		"OtpCode":     "Kode OTP",
		"Title":       "Judul tugas",
		"IDPriorities": "Prioritas",
	}
	if label, ok := labels[field]; ok {
		return label
	}
	return field
}
