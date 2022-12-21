package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/proto"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateDomainHandler(s *scyna.Endpoint, request *proto.CreateDomainRequest) {
	s.Logger.Info("Receive CreateContextRequest")

	if err := validateCreateDomainRequest(request); err != nil {
		s.Done(scyna.REQUEST_INVALID)
		return
	}

	domain := repository.Domain{
		Code: request.Code,
		Name: request.Name,
	}

	if err := repository.CreateDomain(s.Logger, &domain); err != nil {
		s.Error(err)
		return
	}

	s.Done(scyna.OK)
}

func validateCreateDomainRequest(request *proto.CreateDomainRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Code, validation.Required),
		validation.Field(&request.Name, validation.Required, validation.Length(5, 100)))
}
