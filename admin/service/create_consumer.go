package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/engine/admin/proto"
	"github.com/scyna/go/engine/admin/repository"
)

func CreateConsumerHandler(s *scyna.Endpoint, request *proto.CreateConsumerRequest) {
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

	/*TODO: Create JetStream Consumer*/

	s.Done(scyna.OK)
}

func validateCreateConsumerRequest(request *proto.CreateConsumerRequest) error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Sender, validation.Required),
		validation.Field(&request.Receiver, validation.Required))
}
