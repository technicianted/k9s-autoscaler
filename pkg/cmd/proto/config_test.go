package proto

import (
	"fmt"
	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/providers/proto"
	storageproto "k9s-autoscaler/pkg/providers/storage/proto"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestProto(t *testing.T) {
	storageConfig := storageproto.InlineStorageConfig{
		Autoscalers: []*prototypes.Autoscaler{
			{
				Name: "first",
			},
		},
	}
	storageAny, err := anypb.New(&storageConfig)
	require.NoError(t, err)

	config := ControllerConfig{
		StorageClient: &proto.ProviderConfig{
			Name:   "inline",
			Config: storageAny,
		},
		ResyncPeriod: durationpb.New(100 * time.Millisecond),
	}

	bytes, err := protojson.Marshal(&config)
	require.NoError(t, err)
	fmt.Printf("out: %s\n", string(bytes))
}
