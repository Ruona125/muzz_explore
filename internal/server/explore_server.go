package server

import (
	"context"
	"strconv"

	"explore_service/internal/storage"
	explorepb "explore_service/proto"
)

// ExploreServer implements the ExploreService gRPC service.
type ExploreServer struct {
	explorepb.UnimplementedExploreServiceServer
	store *storage.Store
	// pageSize controls the number of results returned per call to
	// ListLikedYou and ListNewLikedYou.  The token returned to the
	// client encodes the next offset.  This value can be tuned
	// depending on expected client consumption patterns.
	pageSize int
}

// NewExploreServer constructs a new ExploreServer with the given
// storage backend.  pageSize controls the default number of
// likers returned per page.  A sensible default of 50 is used if
// pageSize is less than or equal to zero.
func NewExploreServer(store *storage.Store, pageSize int) *ExploreServer {
	if pageSize <= 0 {
		pageSize = 50
	}
	return &ExploreServer{store: store, pageSize: pageSize}
}

// PutDecision records a decision and returns whether the like is mutual.
func (s *ExploreServer) PutDecision(ctx context.Context, req *explorepb.PutDecisionRequest) (*explorepb.PutDecisionResponse, error) {
	mutual, err := s.store.PutDecision(ctx, req.GetActorUserId(), req.GetRecipientUserId(), req.GetLikedRecipient())
	if err != nil {
		return nil, err
	}
	return &explorepb.PutDecisionResponse{MutualLikes: mutual}, nil
}

// ListLikedYou returns all actors who have liked the recipient.  The
// pagination token, if present, is interpreted as the numeric offset
// into the result set.  A new token is returned if additional
// results are available.
func (s *ExploreServer) ListLikedYou(ctx context.Context, req *explorepb.ListLikedYouRequest) (*explorepb.ListLikedYouResponse, error) {
	offset := 0
	if tok := req.GetPaginationToken(); tok != "" {
		// parse the offset encoded as a string.  Ignore errors and
		// fall back to zero.
		if o, err := strconv.Atoi(tok); err == nil && o >= 0 {
			offset = o
		}
	}
	likers, next, err := s.store.ListLikedYou(ctx, req.GetRecipientUserId(), offset, s.pageSize)
	if err != nil {
		return nil, err
	}
	resp := &explorepb.ListLikedYouResponse{Likers: make([]*explorepb.ListLikedYouResponse_Liker, len(likers))}
	for i, l := range likers {
		resp.Likers[i] = &explorepb.ListLikedYouResponse_Liker{
			ActorId:       l.ActorID,
			UnixTimestamp: l.Unix,
		}
	}
	if next != nil {
		resp.NextPaginationToken = next
	}
	return resp, nil
}

// ListNewLikedYou returns actors who like the recipient but have not
// been liked back.  Pagination works in the same way as
// ListLikedYou.
func (s *ExploreServer) ListNewLikedYou(ctx context.Context, req *explorepb.ListLikedYouRequest) (*explorepb.ListLikedYouResponse, error) {
	offset := 0
	if tok := req.GetPaginationToken(); tok != "" {
		if o, err := strconv.Atoi(tok); err == nil && o >= 0 {
			offset = o
		}
	}
	likers, next, err := s.store.ListNewLikedYou(ctx, req.GetRecipientUserId(), offset, s.pageSize)
	if err != nil {
		return nil, err
	}
	resp := &explorepb.ListLikedYouResponse{Likers: make([]*explorepb.ListLikedYouResponse_Liker, len(likers))}
	for i, l := range likers {
		resp.Likers[i] = &explorepb.ListLikedYouResponse_Liker{
			ActorId:       l.ActorID,
			UnixTimestamp: l.Unix,
		}
	}
	if next != nil {
		resp.NextPaginationToken = next
	}
	return resp, nil
}

// CountLikedYou returns the total number of actors who liked the
// recipient.  No pagination is required for counts.
func (s *ExploreServer) CountLikedYou(ctx context.Context, req *explorepb.CountLikedYouRequest) (*explorepb.CountLikedYouResponse, error) {
	count, err := s.store.CountLikedYou(ctx, req.GetRecipientUserId())
	if err != nil {
		return nil, err
	}
	return &explorepb.CountLikedYouResponse{Count: count}, nil
}
