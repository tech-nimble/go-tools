package errors_test

import (
	errorsBase "errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tech-nimble/go-tools/helpers/errors"
)

var _ = Describe("Errors", func() {
	Describe("NoType", func() {
		It("New", func() {
			err := errors.New("Test error")

			Expect(err.Error()).To(Equal("Test error"))
		})
		It("Newf", func() {
			err := errors.Newf("Test error %d", 1)

			Expect(err.Error()).To(Equal("Test error 1"))
		})
		It("Wrap", func() {
			bErr := errorsBase.New("base errors")
			err := errors.Wrap(bErr, "Test error")

			Expect(err.Error()).To(Equal("Test error: base errors"))
		})
		It("Wrapf", func() {
			bErr := errorsBase.New("base errors")
			err := errors.Wrapf(bErr, "Test error %d", 1)

			Expect(err.Error()).To(Equal("Test error 1: base errors"))
		})
		It("GetType", func() {
			err := errors.New("Test error")
			errType := errors.GetType(err)

			Expect(errType).To(Equal(errors.NoType))
		})
	})
	Describe("Domain", func() {
		It("New", func() {
			err := errors.Domain.New("Test error")

			Expect(err.Error()).To(Equal("Test error"))
		})
		It("Newf", func() {
			err := errors.Domain.Newf("Test error %d", 1)

			Expect(err.Error()).To(Equal("Test error 1"))
		})
		It("Wrap", func() {
			bErr := errorsBase.New("base errors")
			err := errors.Domain.Wrap(bErr, "Test error")

			Expect(err.Error()).To(Equal("Test error"))
		})
		It("Wrapf", func() {
			bErr := errorsBase.New("base errors")
			err := errors.Domain.Wrapf(bErr, "Test error %d", 1)

			Expect(err.Error()).To(Equal("Test error 1"))
		})
		It("GetType", func() {
			err := errors.Domain.New("Test error")
			errType := errors.GetType(err)

			Expect(errType).To(Equal(errors.Domain))
		})
	})
	Describe("Extended", func() {
		It("Added error code", func() {
			err := errors.Domain.New("Test error")
			err = errors.AddErrorCode(err, 100)

			Expect(err.Error()).To(Equal("Test error"))
			Expect(errors.GetErrorCode(err)).To(Equal(100))
		})
		It("Get empty error code", func() {
			err := errors.Domain.New("Test error")

			Expect(err.Error()).To(Equal("Test error"))
			Expect(errors.GetErrorCode(err)).To(Equal(0))
		})
		It("Added error data", func() {
			err := errors.Domain.New("Test error")
			err = errors.AddErrorData(err, "testField", "testMessage1")
			err = errors.AddErrorData(err, "testField2", 2)

			Expect(err.Error()).To(Equal("Test error"))
			Expect(errors.GetErrorData(err, "testField")).To(Equal("testMessage1"))
			Expect(errors.GetErrorData(err, "testField2")).To(Equal(2))
		})
		It("Is", func() {
			err := errors.Domain.New("Test error")
			err2 := err

			Expect(errorsBase.Is(err2, err)).To(BeTrue())
		})
		It("Is fail", func() {
			err := errors.Domain.New("Test error")
			err2 := errors.Domain.New("Test error 2")
			err3 := errors.Domain.New("Test error")

			Expect(errorsBase.Is(err2, err)).To(BeFalse())
			Expect(errorsBase.Is(err3, err)).To(BeFalse())
		})
		It("Is by errors with code", func() {
			err := errors.Domain.New("Test error")
			err2 := errors.AddErrorCode(err, 100)

			Expect(errorsBase.Is(err2, err)).To(BeTrue())
		})
		It("Is by errors with data", func() {
			err := errors.Domain.New("Test error")
			err2 := errors.AddErrorData(err, "test", 100)

			Expect(errorsBase.Is(err2, err)).To(BeTrue())
		})
		It("Вложенная дата", func() {
			err := errors.Domain.New("Test error")
			err2 := errors.AddErrorData(err, "test", 100)
			err3 := errors.Domain.Wrap(err2, "NewError")
			err4 := errors.AddErrorData(err3, "test2", 200)

			Expect(errors.GetErrorData(err4, "test")).To(Equal(100))
			Expect(errors.GetErrorData(err4, "test2")).To(Equal(200))
		})
	})
})
