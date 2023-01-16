package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/proto"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateConsumerHandler(s *scyna.Endpoint, request *proto.CreateConsumerRequest) scyna.Error {
	s.Logger.Info("Receive CreateConsumerRequest")

	if err := validateCreateConsumerRequest(request); err != nil {
		return scyna.REQUEST_INVALID
	}

	if _, context := repository.GetContext(s.Logger, request.Sender); context != nil {
		return scyna.REQUEST_INVALID
	}

	if _, context := repository.GetContext(s.Logger, request.Receiver); context != nil {
		return scyna.REQUEST_INVALID
	}

	/*TODO: Create JetStream Consumer*/

	return scyna.OK
}

func validateCreateConsumerRequest(request *proto.CreateConsumerRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Sender, validation.Required),
		validation.Field(&request.Receiver, validation.Required))
}
