package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	translations "github.com/go-playground/validator/v10/translations/ru"
)

type AppValidator struct {
	universalTranslator *ut.UniversalTranslator
	LocaleTranslator    ut.Translator
	*validator.Validate
}

func NewValidator() (appValidator *AppValidator, err error) {
	const locale = "ru"

	translator := ru.New()
	appValidator = &AppValidator{
		universalTranslator: ut.New(translator, translator),
	}

	var isFound bool
	appValidator.LocaleTranslator, isFound = appValidator.universalTranslator.GetTranslator(locale)
	if isFound == false {
		err = fmt.Errorf("universalTranslator not found: %s", locale)
		return
	}

	appValidator.Validate = validator.New()
	err = translations.RegisterDefaultTranslations(appValidator.Validate, appValidator.LocaleTranslator)

	return
}

func (appValidator AppValidator) ErrorTranslated(validatorError error) (err error) {
	errs, ok := validatorError.(validator.ValidationErrors)
	if !ok {
		log.Println("appValidator.ErrorTranslated error: can't cast error")
		return validatorError
	}

	var errString string
	for _, e := range errs {
		errString += e.StructNamespace() + " " + e.Translate(appValidator.LocaleTranslator) + "; "
	}

	err = errors.New(errString)
	return
}
