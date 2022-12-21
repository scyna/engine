package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateConsumerHandler(s *scyna.Endpoint, request *scyna.CreateConsumerRequest) {
	s.Logger.Info("Receive CreateConsumerRequest")

	if err := validateCreateConsumerRequest(request); err != nil {
		s.Done(scyna.REQUEST_INVALID)
		return
	}

	if _, context := repository.GetContext(s.Logger, request.Sender); context != nil {
		s.Error(scyna.REQUEST_INVALID)
		return
	}

	if _, context := repository.GetContext(s.Logger, request.Receiver); context != nil {
		s.Error(scyna.REQUEST_INVALID)
		return
	}

	/*TODO: create JetStream*/

	s.Done(scyna.OK)
}

func validateCreateConsumerRequest(request *scyna.CreateConsumerRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Sender, validation.Required),
		validation.Field(&request.Receiver, validation.Required))
}
