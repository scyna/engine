package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateContextHandler(s *scyna.Endpoint, request *scyna.CreateContextRequest) {
	s.Logger.Info("Receive CreateContextRequest")

	if err := validateCreateContextRequest(request); err != nil {
		s.Done(scyna.REQUEST_INVALID)
		return
	}

	if _, domain := repository.GetDomain(s.Logger, request.Code); domain != nil {
		s.Error(scyna.REQUEST_INVALID)
		return
	}

	context := repository.Context{
		Domain: request.Domain,
		Code:   request.Code,
		Name:   request.Name,
	}

	if err := repository.CreateContext(s.Logger, &context); err != nil {
		s.Error(err)
		return
	}

	s.Done(scyna.OK)
}

func validateCreateContextRequest(request *scyna.CreateContextRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Domain, validation.Required),
		validation.Field(&request.Code, validation.Required),
		validation.Field(&request.Name, validation.Required, validation.Length(5, 100)))
}
