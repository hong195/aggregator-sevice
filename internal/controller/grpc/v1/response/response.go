package response

import (
	v1 "github.com/hong195/aggregator-sevice/docs/proto/v1"
	"github.com/hong195/aggregator-sevice/internal/usecase/query"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewFindByIdResponse -.
func NewFindByIdResponse(packet query.DataPacketView) *v1.FindPacketByIDResponse {

	return &v1.FindPacketByIDResponse{
		Packet: &v1.DataPacket{
			Id:        packet.ID,
			Timestamp: timestamppb.New(packet.Timestamp.UTC()),
			MaxValue:  int64(packet.MaxValue),
		},
	}
}

// NewListPacketsByPeriodResponse -
func NewListPacketsByPeriodResponse(packets []query.DataPacketView) *v1.ListPacketsByPeriodResponse {
	items := make([]*v1.DataPacket, 0, len(packets))
	for i := range packets {
		p := packets[i]
		items = append(items, &v1.DataPacket{
			Id:        p.ID,
			Timestamp: timestamppb.New(p.Timestamp.UTC()),
			MaxValue:  int64(p.MaxValue),
		})
	}
	return &v1.ListPacketsByPeriodResponse{Items: items}
}
