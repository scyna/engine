package organization

import (
	scyna "github.com/scyna/core"
	proto "github.com/scyna/go/manager/.proto/generated"
)

func DestroyOrganization(s *scyna.Service, request *proto.DestroyOrganizationRequest) {
	s.Logger.Info("Receive DestroyOrganizationRequest")

	/*TODO*/

	s.Done(scyna.OK)
}
