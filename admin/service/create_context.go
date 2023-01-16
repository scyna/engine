package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/proto"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateContextHandler(s *scyna.Endpoint, request *proto.CreateContextRequest) scyna.Error {
	s.Logger.Info("Receive CreateContextRequest")

	if err := validateCreateContextRequest(request); err != nil {
		return scyna.REQUEST_INVALID
	}

	if _, domain := repository.GetDomain(s.Logger, request.Code); domain != nil {
		return scyna.REQUEST_INVALID
	}

	context := repository.Context{
		Domain: request.Domain,
		Code:   request.Code,
		Name:   request.Name,
	}

	if err := repository.CreateContext(s.Logger, &context); err != nil {
		return err
	}

	/*TODO: Create Stream*/

	return scyna.OK
}

func validateCreateContextRequest(request *proto.CreateContextRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Domain, validation.Required),
		validation.Field(&request.Code, validation.Required),
		validation.Field(&request.Name, validation.Required, validation.Length(5, 100)))
}
