package v2ray

import (
	"context"
	"github.com/devfeel/mapper"
	"github.com/v2fly/v2ray-core/v4/app/observatory"
	pb "github.com/v2fly/v2ray-core/v4/app/observatory/command"
	"github.com/v2rayA/v2rayA-lib4/db/configure"
	"github.com/v2rayA/v2rayA-lib4/pkg/util/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
	"time"
)

var (
	ApiProducts = []string{
		"observatory",
	}
	ApiFeed *Feed
)

const (
	ApiFeedBoxSize  = 10
	ApiFeedInterval = 1 * time.Second
)

type OutboundStatus struct {
	Alive           bool             `json:"alive"`
	Delay           int64            `json:"delay"`
	LastErrorReason string           `json:"last_error_reason"`
	OutboundTag     string           `json:"outbound_tag"`
	Which           *configure.Which `json:"which"`
	LastSeenTime    int64            `json:"last_seen_time"`
	LastTryTime     int64            `json:"last_try_time"`
}

func init() {
	mapper.Register(&observatory.OutboundStatus{})

	ApiFeed = NewSubscriptions(ApiFeedBoxSize)
	for _, product := range ApiProducts {
		ApiFeed.RegisterProduct(product)
	}
}

func getObservatoryResponse(conn *grpc.ClientConn) (r *pb.GetOutboundStatusResponse, err error) {
	c := pb.NewObservatoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err = c.GetOutboundStatus(ctx, &pb.GetOutboundStatusRequest{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func ObservatoryProducer(apiPort int) (closeFunc func()) {
	closed := make(chan struct{})
	go func() {
		const product = "observatory"
		var conn *grpc.ClientConn
	nextLoop:
		for {
			select {
			case <-closed:
				return
			default:
			}
			p := ProcessManager.Process()
			if p == nil {
				time.Sleep(ApiFeedInterval)
				continue
			}
			// Set up a connection to the server.
			if conn == nil {
				ctx, cancel := context.WithTimeout(context.Background(), ApiFeedInterval)
				defer cancel()
				c, err := grpc.DialContext(
					ctx,
					net.JoinHostPort("127.0.0.1", strconv.Itoa(apiPort)),
					grpc.WithInsecure(),
					grpc.WithBlock(),
				)
				if err != nil {
					log.Warn("ObservatoryProducer: did not connect: %v", err)
					continue nextLoop
				}
				defer c.Close()
				conn = c
			}
			r, err := getObservatoryResponse(conn)
			if err != nil {
				if status.Code(err) == codes.Unavailable {
					// the connection is reliable, and reconnect
					conn = nil
					continue nextLoop
				}
				log.Warn("ObservatoryProducer: %v", err)
			} else {
				outboundStatus := r.GetStatus().GetStatus()
				os := make([]OutboundStatus, len(outboundStatus))
				css := configure.GetConnectedServers()
				for i := range outboundStatus {
					_ = mapper.AutoMapper(outboundStatus[i], &os[i])
					index := p.tag2WhichIndex[os[i].OutboundTag]
					if index >= css.Len() {
						continue nextLoop
					}
					os[i].Which = css.Get()[index]
					var w []configure.Which
					for _, v := range css.Get() {
						w = append(w, *v)
					}
				}
				ApiFeed.ProductMessage(product, os)
			}
			time.Sleep(ApiFeedInterval)
		}
	}()
	return func() {
		close(closed)
	}
}
