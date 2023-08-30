package session

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	"google.golang.org/protobuf/proto"
)

func TestCreateSession(t *testing.T) {

	request := scyna_proto.CreateSessionRequest{
		Module: "scyna_test",
		Secret: "123456",
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		t.Fatal("Can not marshal request")
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8081"+scyna_const.SESSION_CREATE_URL, bytes.NewBuffer(data))
	if err != nil {
		t.Fatal("Error in create http request:", err)
	}

	res, err := scyna.HttpClient().Do(req)
	if err != nil {
		t.Fatal("Error in send http request:", err)
	}

	if res.StatusCode != 200 {
		t.Fatal("Error in autheticate")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Can not read response body:", err)
	}

	var response scyna_proto.CreateSessionResponse
	if err := proto.Unmarshal(resBody, &response); err != nil {
		t.Fatal("Authenticate error")
	}
}
