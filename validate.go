package validate

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	localTime "github.com/go-tron/local-time"
	"reflect"
	"strings"
)

var Validate validate

func init() {
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := fld.Tag.Get("comment"); name != "" {
			return name
		}
		return fld.Tag.Get("json")
	})
	v.RegisterCustomTypeFunc(LocalTime, localTime.Time{})

	zh_translations.RegisterDefaultTranslations(v, trans)
	Validate.Validate = v
	Validate.Trans = trans
}

type validate struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

func (v *validate) Struct(p interface{}) error {
	err := v.Validate.Struct(p)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	var errList = make([]string, len(errs))
	for i, err := range errs {
		errList[i] = err.Translate(v.Trans)
	}
	return errors.New(strings.Join(errList, ","))
}

func LocalTime(field reflect.Value) interface{} {
	if value, ok := field.Interface().(localTime.Time); ok {
		return value.String()
	}
	return nil
}
