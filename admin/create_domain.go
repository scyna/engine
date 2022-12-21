package admin

import (
	scyna "github.com/scyna/core"
)

func CreateDomainHandler(s *scyna.Endpoint, request *scyna.CreateDomainRequest) {
	s.Logger.Info("Receive CreateContextRequest")

	/*TODO: check request*/
	/*TODO: check if domain exists*/
	/*TODO: save domain*/

	s.Done(scyna.OK)
}
