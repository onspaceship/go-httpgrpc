package httpgrpc

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// Assert *Client implements ClientConnInterface.
var _ grpc.ClientConnInterface = (*ClientConn)(nil)

type ClientConn struct {
	BaseURI            string
	AuthorizationToken string
}

func (client *ClientConn) Invoke(ctx context.Context, method string, in interface{}, out interface{}, _ ...grpc.CallOption) error {
	msg, err := proto.Marshal(in.(proto.Message))
	body := bytes.NewBuffer(msg)

	req, err := http.NewRequestWithContext(ctx, "POST", client.BaseURI+method, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/grpc")
	if client.AuthorizationToken != "" {
		req.Header.Set("Authorization", "Bearer "+client.AuthorizationToken)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	responseBody, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New(string(responseBody))
	}

	err = proto.Unmarshal(responseBody, out.(proto.Message))

	return err
}

func (client *ClientConn) NewStream(_ context.Context, _ *grpc.StreamDesc, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, errors.New("streaming not implemented")
}
