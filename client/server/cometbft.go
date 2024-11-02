package server

import (
	"errors"
	"net/http"
	"slices"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	"github.com/piplabs/story/client/server/utils"
	liberrors "github.com/piplabs/story/lib/errors"
)

func (s *Server) initCometBFTRoute() {
	s.httpMux.HandleFunc("/cometbft/block_events", utils.AutoWrap(s.aminoCodec, s.GetCometbftBlockEvents))
}

func (s *Server) GetCometbftBlockEvents(req *getCometbftBlockEventsRequest, r *http.Request) (resp any, err error) {
	if req.To-req.From > 100 {
		return nil, errors.New("search max 100 blocks")
	}

	if len(req.EventTypeFilter) == 0 {
		return nil, errors.New("event filter empty")
	}

	curBlock, err := s.cl.Block(r.Context(), nil)
	if err != nil {
		return nil, liberrors.Wrap(err, "failed to get the current block")
	}

	allRetBlock := make([]*getCometbftBlockEventsBlockResults, 0)
	for i := req.From; i < min(req.To, curBlock.Block.Height); i++ {
		results, err := s.cl.BlockResults(r.Context(), &i)
		if err != nil {
			return nil, liberrors.Wrap(err, "failed to get block result")
		}

		events := slices.DeleteFunc(results.FinalizeBlockEvents, func(event abcitypes.Event) bool {
			return !slices.Contains(req.EventTypeFilter, event.Type)
		})

		if len(events) > 0 {
			allRetBlock = append(allRetBlock, &getCometbftBlockEventsBlockResults{
				Height: results.Height,
				FinalizeBlockEvents: slices.DeleteFunc(results.FinalizeBlockEvents, func(event abcitypes.Event) bool {
					return !slices.Contains(req.EventTypeFilter, event.Type)
				}),
			})
		}
	}

	return allRetBlock, nil
}
