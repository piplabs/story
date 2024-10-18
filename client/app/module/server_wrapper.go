package module

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/grpc"

	googlegrpc "google.golang.org/grpc"
)

// serverWrapper wraps the grpc.Server for registering a service but includes
// logic to extract all the sdk.Msg types that the service declares in its
// methods and fires a callback to add them to the configurator. This allows us
// to create a map of which messages are accepted across which versions.
type serverWrapper struct {
	addMessages func(msgs []string)
	msgServer   grpc.Server
}

func (s *serverWrapper) RegisterService(sd *googlegrpc.ServiceDesc, v interface{}) {
	msgs := make([]string, len(sd.Methods))
	for idx, method := range sd.Methods {
		// we execute the handler to extract the message type
		_, _ = method.Handler(nil, context.Background(), func(i interface{}) error {
			msg, ok := i.(sdk.Msg)
			if !ok {
				panic(fmt.Errorf("unable to register service method %s/%s: %T does not implement sdk.Msg", sd.ServiceName, method.MethodName, i))
			}
			msgs[idx] = sdk.MsgTypeURL(msg)
			return nil
		}, noopInterceptor)
	}
	s.addMessages(msgs)
	// call the underlying msg server to actually register the grpc server
	s.msgServer.RegisterService(sd, v)
}

func noopInterceptor(_ context.Context, _ interface{}, _ *googlegrpc.UnaryServerInfo, _ googlegrpc.UnaryHandler) (interface{}, error) {
	return nil, nil
}
