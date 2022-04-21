package call

import (
	scyna "github.com/scyna/go"
	"google.golang.org/protobuf/proto"
)

func Write(data []byte) {
	var signal scyna.WriteCallSignal
	if err := proto.Unmarshal(data, &signal); err != nil {
		scyna.LOG.Error("Can not parse WriteCallSignal")
		return
	}

	/*TODO*/
}
