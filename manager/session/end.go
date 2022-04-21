package session

import (
	scyna "github.com/scyna/go"
	"google.golang.org/protobuf/proto"
)

func End(data []byte) {
	var signal scyna.EndSessionSignal
	if err := proto.Unmarshal(data, &signal); err != nil {
		scyna.LOG.Error("Can not parse EndSessionSignal")
		return
	}

	/*TODO*/

}
