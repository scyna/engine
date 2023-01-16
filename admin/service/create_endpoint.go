package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/proto"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateEndpointHandler(s *scyna.Endpoint, request *proto.CreateEndpointRequest) scyna.Error {
	s.Logger.Info("Receive CreateEndpointRequest")

	if err := validateCreateEndpointRequest(request); err != nil {
		return scyna.REQUEST_INVALID
	}

	if _, context := repository.GetContext(s.Logger, request.Context); context != nil {
		return scyna.REQUEST_INVALID
	}

	endpoint := repository.Endpoint{
		Context: request.Context,
		URL:     request.URL,
		Name:    request.Name,
	}

	if err := repository.CreateEndpoint(s.Logger, &endpoint); err != nil {
		return scyna.SERVER_ERROR
	}
	return scyna.OK
}

func validateCreateEndpointRequest(request *proto.CreateEndpointRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Context, validation.Required),
		validation.Field(&request.URL, validation.Required),
		validation.Field(&request.Name, validation.Required, validation.Length(5, 100)))
}
