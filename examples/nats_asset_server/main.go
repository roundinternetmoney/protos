// example server client using nats/micro and protoc-gen-go

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	basev1 "roundinternet.money/protos/gen/dex/base/v1"
	pb "roundinternet.money/protos/gen/dex/base/v1"
)

var ErrMatch = errors.New("invalid request: dex does not match")

type DexAssetService struct {
	basev1.Dex
	streams []basev1.PayloadKind

	assets []*basev1.Asset
}

func (s *DexAssetService) matchDex(ds []basev1.Dex) bool {
	if s.Dex == basev1.Dex_DEX_UNSPECIFIED {
		return true
	}
	for _, opt := range ds {
		if opt == s.Dex {
			return true
		}
	}
	return false
}

func (s *DexAssetService) matchKind(ds basev1.PayloadKind) bool {
	for _, opt := range s.streams {
		if opt == ds {
			return true
		}
	}
	return false
}

func (s *DexAssetService) matchRelatedAsset(related string, known []*pb.Asset) *pb.Asset {
	for _, k := range known {
		if k.Ticker == related {
			return k
		}
	}
	return nil
}

func (s *DexAssetService) matchRelatedAssets(related []string, known []*pb.Asset) []*pb.Asset {
	out := make([]*pb.Asset, 0)
	for _, r := range related {
		if found := s.matchRelatedAsset(r, known); found != nil {
			out = append(out, found)
		}
	}
	return out
}

func (s *DexAssetService) StreamAsset(ctx context.Context, req *basev1.StreamAssetRequest, stream *basev1.ReaderService_StreamAsset_Stream,
) error {
	return nil
}

func (s *DexAssetService) StreamMargin(ctx context.Context, req *basev1.StreamMarginRequest, stream *basev1.ReaderService_StreamMargin_Stream,
) error {
	return nil
}

func (s *DexAssetService) StreamPosition(ctx context.Context, req *basev1.StreamPositionRequest, stream *basev1.ReaderService_StreamPosition_Stream,
) error {
	return nil
}

func (s *DexAssetService) StreamDiff(ctx context.Context, req *basev1.StreamDiffRequest, stream *basev1.ReaderService_StreamDiff_Stream,
) error {
	return nil
}

func (s *DexAssetService) StreamMids(ctx context.Context, req *basev1.StreamMidsRequest, stream *basev1.ReaderService_StreamMids_Stream,
) error {
	return nil
}

func (s *DexAssetService) DexAsset(ctx context.Context, req *basev1.DexAssetRequest,
) (*basev1.DexAssetResponse, error) {
	if !s.matchDex(req.Dexes) {
		return nil, ErrMatch
	}
	if len(req.RelatedAssets) > 0 {
		return &basev1.DexAssetResponse{A: s.matchRelatedAssets(req.RelatedAssets, s.assets)}, nil
	}
	fmt.Println(req.String())
	fmt.Println(s.assets)
	return &basev1.DexAssetResponse{A: s.assets}, nil
}

func (s *DexAssetService) Stream(ctx context.Context, req *basev1.StreamRequest, stream *basev1.ReaderService_Stream_Stream,
) error {
	if !s.matchDex(req.Dexes) {
		return ErrMatch
	}
	for _, kind := range req.Kinds {
		if !s.matchKind(kind) {
			continue
		}

		for i := 0; i < 5; i++ {
			if err := stream.Send(&basev1.StreamResponse{
				Kind: kind,
				Payload: &basev1.StreamResponse_A{A: &basev1.Asset{
					Ticker:      fmt.Sprintf("EXAMPLEUSD-%d", i),
					Id:          fmt.Sprintf("example-%d", i),
					MinQuantity: 0.001,
					MaxQuantity: 100,
					LotSize:     0.01,
					TickSize:    0.1,
				}},
			}); err != nil {
				return err
			}
			time.Sleep(200 * time.Millisecond) // "work"
		}
	}
	return ErrMatch
}

func (s *DexAssetService) Upstream(ctx context.Context, stream *basev1.ReaderService_Upstream_Stream,
) (*basev1.UpstreamResponse, error) {
	return nil, nil
}

func (s *DexAssetService) Downstream(ctx context.Context, stream *basev1.ReaderService_Downstream_Stream,
) error {
	return nil
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	svc := &DexAssetService{
		streams: []basev1.PayloadKind{
			basev1.PayloadKind_PAYLOAD_KIND_ASSET,
		},
		assets: []*basev1.Asset{
			{
				Ticker:      "EXAMPLEUSD",
				Id:          "example",
				MinQuantity: 0.001,
				MaxQuantity: 100,
				LotSize:     0.01,
				TickSize:    0.1,
			},
		},
	}

	_, err = basev1.RegisterReaderServiceHandlers(nc, svc)
	if err != nil {
		log.Fatal(err)
	}

	select {} // block
}
