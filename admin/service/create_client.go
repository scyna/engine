package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateClientHandler(s *scyna.Endpoint, request *scyna.CreateClientRequest) {
	s.Logger.Info("Receive CreateClientRequest")

	if err := validateCreateClientRequest(request); err != nil {
		s.Done(scyna.REQUEST_INVALID)
		return
	}

	if _, domain := repository.GetDomain(s.Logger, request.Domain); domain != nil {
		s.Error(scyna.REQUEST_INVALID)
		return
	}

	client := repository.Client{
		Domain: request.Domain,
		ID:     request.ID,
		Secret: request.Secret,
	}

	if err := repository.CreateClient(s.Logger, &client); err != nil {
		s.Error(err)
		return
	}

	s.Done(scyna.OK)
}

func validateCreateClientRequest(request *scyna.CreateClientRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Domain, validation.Required),
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.Secret, validation.Required))
}
