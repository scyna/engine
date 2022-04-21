package session

import (
	scyna "github.com/scyna/go"
	"google.golang.org/protobuf/proto"
)

func Update(data []byte) {
	var signal scyna.UpdateSessionSignal
	if err := proto.Unmarshal(data, &signal); err != nil {
		scyna.LOG.Error("Can not parse UpdateSessionSignal")
		return
	}

	/*TODO*/

}
