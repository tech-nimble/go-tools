package errors_test

import (
	"encoding/json"

	base "github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tech-nimble/go-i18n"
	"github.com/tech-nimble/go-tools/helpers/errors"
	"golang.org/x/text/language"
)

var _ = Describe("Errors Translate", func() {
	Describe("Translate", func() {
		It("T", func() {
			russian := `{
			  "error.remain_days_from_day": {
				"one": "остался {{.Days}} день",
				"few": "осталось {{.Days}} дня",
				"many": "осталось {{.Days}} дней",
				"other": "other {{.Days}} days"
			  }
			}`
			kz := `{
			  "error.remain_days_from_day": {
				"one": "{{.Days}} күн қалды",
				"few": "{{.Days}} күн қалды",
				"many": "{{.Days}} күн қалды",
				"other": "{{.Days}} күн қалды"
			  }
			}`

			bundle := base.NewBundle(language.Russian)

			bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

			bundle.MustParseMessageFileBytes([]byte(russian), "ru.json")
			bundle.MustParseMessageFileBytes([]byte(kz), "kk.json")

			responder := i18n.NewResponder(bundle)

			localizeRu := i18n.NewLocalize(responder.I18nBundle, "ru")
			localizeKz := i18n.NewLocalize(responder.I18nBundle, "kk")

			err := errors.Domain.NewTf(
				"остался день",
				"error.remain_days_from_day",
				map[string]any{"Days": "1"},
				1,
			)

			err1 := errors.AddErrLocalize(err, localizeKz)
			err2 := errors.AddErrLocalize(err, localizeRu)

			Expect(err1.Error()).To(Equal("1 күн қалды"))
			Expect(err2.Error()).To(Equal("остался 1 день"))
			Expect(err.Error()).To(Equal("остался день"))

			err3 := errors.Domain.NewT("Остался мало дней", "error.remain_days_from_day")
			err4 := errors.AddErrLocalizePluralCount(errors.AddErrLocalizeData(err3, map[string]any{"Days": 2}), 2)

			err5 := errors.AddErrLocalize(err4, localizeKz)
			err6 := errors.AddErrLocalize(err4, localizeRu)

			Expect(err5.Error()).To(Equal("2 күн қалды"))
			Expect(err6.Error()).To(Equal("осталось 2 дня"))
			Expect(err3.Error()).To(Equal("Остался мало дней"))
		})
	})
})
