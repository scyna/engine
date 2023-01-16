package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/proto"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateClientHandler(s *scyna.Endpoint, request *proto.CreateClientRequest) scyna.Error {
	s.Logger.Info("Receive CreateClientRequest")

	if err := validateCreateClientRequest(request); err != nil {
		return scyna.REQUEST_INVALID
	}

	if _, domain := repository.GetDomain(s.Logger, request.Domain); domain != nil {
		return scyna.REQUEST_INVALID
	}

	client := repository.Client{
		Domain: request.Domain,
		ID:     request.ID,
		Secret: request.Secret,
	}

	if err := repository.CreateClient(s.Logger, &client); err != nil {
		return err
	}

	return scyna.OK
}

func validateCreateClientRequest(request *proto.CreateClientRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Domain, validation.Required),
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.Secret, validation.Required))
}
