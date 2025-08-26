package explorepb

import (
	context "context"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Ensure the generated code is sufficiently recent.
const _ = proto.ProtoPackageIsVersion4

// ListLikedYouRequest is the input message for ListLikedYou and
// ListNewLikedYou RPCs.
type ListLikedYouRequest struct {
	RecipientUserId string  `protobuf:"bytes,1,opt,name=recipient_user_id,json=recipientUserId,proto3" json:"recipient_user_id,omitempty"`
	PaginationToken *string `protobuf:"bytes,2,opt,name=pagination_token,json=paginationToken,proto3,oneof" json:"pagination_token,omitempty"`
}

func (m *ListLikedYouRequest) Reset()         { *m = ListLikedYouRequest{} }
func (m *ListLikedYouRequest) String() string { return proto.CompactTextString(m) }
func (*ListLikedYouRequest) ProtoMessage()    {}

func (m *ListLikedYouRequest) GetRecipientUserId() string {
	if m != nil {
		return m.RecipientUserId
	}
	return ""
}

func (m *ListLikedYouRequest) GetPaginationToken() string {
	if m != nil && m.PaginationToken != nil {
		return *m.PaginationToken
	}
	return ""
}

// ListLikedYouResponse_Liker represents a single like.
type ListLikedYouResponse_Liker struct {
	ActorId       string `protobuf:"bytes,1,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
	UnixTimestamp uint64 `protobuf:"varint,2,opt,name=unix_timestamp,json=unixTimestamp,proto3" json:"unix_timestamp,omitempty"`
}

func (m *ListLikedYouResponse_Liker) Reset()         { *m = ListLikedYouResponse_Liker{} }
func (m *ListLikedYouResponse_Liker) String() string { return proto.CompactTextString(m) }
func (*ListLikedYouResponse_Liker) ProtoMessage()    {}

func (m *ListLikedYouResponse_Liker) GetActorId() string {
	if m != nil {
		return m.ActorId
	}
	return ""
}

func (m *ListLikedYouResponse_Liker) GetUnixTimestamp() uint64 {
	if m != nil {
		return m.UnixTimestamp
	}
	return 0
}

// ListLikedYouResponse is the response for list RPCs.
type ListLikedYouResponse struct {
	Likers              []*ListLikedYouResponse_Liker `protobuf:"bytes,1,rep,name=likers,proto3" json:"likers,omitempty"`
	NextPaginationToken *string                       `protobuf:"bytes,2,opt,name=next_pagination_token,json=nextPaginationToken,proto3,oneof" json:"next_pagination_token,omitempty"`
}

func (m *ListLikedYouResponse) Reset()         { *m = ListLikedYouResponse{} }
func (m *ListLikedYouResponse) String() string { return proto.CompactTextString(m) }
func (*ListLikedYouResponse) ProtoMessage()    {}

func (m *ListLikedYouResponse) GetLikers() []*ListLikedYouResponse_Liker {
	if m != nil {
		return m.Likers
	}
	return nil
}

func (m *ListLikedYouResponse) GetNextPaginationToken() string {
	if m != nil && m.NextPaginationToken != nil {
		return *m.NextPaginationToken
	}
	return ""
}

// CountLikedYouRequest is the input for CountLikedYou RPC.
type CountLikedYouRequest struct {
	RecipientUserId string `protobuf:"bytes,1,opt,name=recipient_user_id,json=recipientUserId,proto3" json:"recipient_user_id,omitempty"`
}

func (m *CountLikedYouRequest) Reset()         { *m = CountLikedYouRequest{} }
func (m *CountLikedYouRequest) String() string { return proto.CompactTextString(m) }
func (*CountLikedYouRequest) ProtoMessage()    {}

func (m *CountLikedYouRequest) GetRecipientUserId() string {
	if m != nil {
		return m.RecipientUserId
	}
	return ""
}

// CountLikedYouResponse returns a count.
type CountLikedYouResponse struct {
	Count uint64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (m *CountLikedYouResponse) Reset()         { *m = CountLikedYouResponse{} }
func (m *CountLikedYouResponse) String() string { return proto.CompactTextString(m) }
func (*CountLikedYouResponse) ProtoMessage()    {}

func (m *CountLikedYouResponse) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

// PutDecisionRequest records a like or pass.
type PutDecisionRequest struct {
	ActorUserId     string `protobuf:"bytes,1,opt,name=actor_user_id,json=actorUserId,proto3" json:"actor_user_id,omitempty"`
	RecipientUserId string `protobuf:"bytes,2,opt,name=recipient_user_id,json=recipientUserId,proto3" json:"recipient_user_id,omitempty"`
	LikedRecipient  bool   `protobuf:"varint,3,opt,name=liked_recipient,json=likedRecipient,proto3" json:"liked_recipient,omitempty"`
}

func (m *PutDecisionRequest) Reset()         { *m = PutDecisionRequest{} }
func (m *PutDecisionRequest) String() string { return proto.CompactTextString(m) }
func (*PutDecisionRequest) ProtoMessage()    {}

func (m *PutDecisionRequest) GetActorUserId() string {
	if m != nil {
		return m.ActorUserId
	}
	return ""
}

func (m *PutDecisionRequest) GetRecipientUserId() string {
	if m != nil {
		return m.RecipientUserId
	}
	return ""
}

func (m *PutDecisionRequest) GetLikedRecipient() bool {
	if m != nil {
		return m.LikedRecipient
	}
	return false
}

// PutDecisionResponse returns whether the like is mutual.
type PutDecisionResponse struct {
	MutualLikes bool `protobuf:"varint,1,opt,name=mutual_likes,json=mutualLikes,proto3" json:"mutual_likes,omitempty"`
}

func (m *PutDecisionResponse) Reset()         { *m = PutDecisionResponse{} }
func (m *PutDecisionResponse) String() string { return proto.CompactTextString(m) }
func (*PutDecisionResponse) ProtoMessage()    {}

func (m *PutDecisionResponse) GetMutualLikes() bool {
	if m != nil {
		return m.MutualLikes
	}
	return false
}

// ExploreServiceServer defines the server API for ExploreService.
// All implementations must embed UnimplementedExploreServiceServer for
// forward compatibility.
type ExploreServiceServer interface {
	ListLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error)
	ListNewLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error)
	CountLikedYou(context.Context, *CountLikedYouRequest) (*CountLikedYouResponse, error)
	PutDecision(context.Context, *PutDecisionRequest) (*PutDecisionResponse, error)
}

// UnimplementedExploreServiceServer can be embedded to have forward
// compatible implementations.
type UnimplementedExploreServiceServer struct{}

func (*UnimplementedExploreServiceServer) ListLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLikedYou not implemented")
}
func (*UnimplementedExploreServiceServer) ListNewLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNewLikedYou not implemented")
}
func (*UnimplementedExploreServiceServer) CountLikedYou(context.Context, *CountLikedYouRequest) (*CountLikedYouResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountLikedYou not implemented")
}
func (*UnimplementedExploreServiceServer) PutDecision(context.Context, *PutDecisionRequest) (*PutDecisionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutDecision not implemented")
}

// RegisterExploreServiceServer registers the service implementation with a gRPC server.
func RegisterExploreServiceServer(s *grpc.Server, srv ExploreServiceServer) {
	s.RegisterService(&ExploreService_ServiceDesc, srv)
}

// ExploreService_ServiceDesc is the grpc.ServiceDesc for ExploreService service.
var ExploreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "explore.ExploreService",
	HandlerType: (*ExploreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListLikedYou",
			Handler:    _ExploreService_ListLikedYou_Handler,
		},
		{
			MethodName: "ListNewLikedYou",
			Handler:    _ExploreService_ListNewLikedYou_Handler,
		},
		{
			MethodName: "CountLikedYou",
			Handler:    _ExploreService_CountLikedYou_Handler,
		},
		{
			MethodName: "PutDecision",
			Handler:    _ExploreService_PutDecision_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/explore-service.proto",
}

func _ExploreService_ListLikedYou_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLikedYouRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).ListLikedYou(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/explore.ExploreService/ListLikedYou",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).ListLikedYou(ctx, req.(*ListLikedYouRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExploreService_ListNewLikedYou_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLikedYouRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).ListNewLikedYou(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/explore.ExploreService/ListNewLikedYou",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).ListNewLikedYou(ctx, req.(*ListLikedYouRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExploreService_CountLikedYou_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountLikedYouRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).CountLikedYou(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/explore.ExploreService/CountLikedYou",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).CountLikedYou(ctx, req.(*CountLikedYouRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExploreService_PutDecision_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutDecisionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).PutDecision(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/explore.ExploreService/PutDecision",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).PutDecision(ctx, req.(*PutDecisionRequest))
	}
	return interceptor(ctx, in, info, handler)
}
