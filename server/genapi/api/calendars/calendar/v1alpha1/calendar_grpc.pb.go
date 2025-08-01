// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/calendars/calendar/v1alpha1/calendar.proto

package calendarv1alpha1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	CalendarService_CreateCalendar_FullMethodName = "/api.calendars.calendar.v1alpha1.CalendarService/CreateCalendar"
	CalendarService_ListCalendars_FullMethodName  = "/api.calendars.calendar.v1alpha1.CalendarService/ListCalendars"
	CalendarService_UpdateCalendar_FullMethodName = "/api.calendars.calendar.v1alpha1.CalendarService/UpdateCalendar"
	CalendarService_DeleteCalendar_FullMethodName = "/api.calendars.calendar.v1alpha1.CalendarService/DeleteCalendar"
	CalendarService_GetCalendar_FullMethodName    = "/api.calendars.calendar.v1alpha1.CalendarService/GetCalendar"
)

// CalendarServiceClient is the client API for CalendarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// the calendar service
type CalendarServiceClient interface {
	// create a calendar
	CreateCalendar(ctx context.Context, in *CreateCalendarRequest, opts ...grpc.CallOption) (*Calendar, error)
	// list calendars
	ListCalendars(ctx context.Context, in *ListCalendarsRequest, opts ...grpc.CallOption) (*ListCalendarsResponse, error)
	// update a calendar
	UpdateCalendar(ctx context.Context, in *UpdateCalendarRequest, opts ...grpc.CallOption) (*Calendar, error)
	// delete` a calendar
	DeleteCalendar(ctx context.Context, in *DeleteCalendarRequest, opts ...grpc.CallOption) (*Calendar, error)
	// get a calendar
	GetCalendar(ctx context.Context, in *GetCalendarRequest, opts ...grpc.CallOption) (*Calendar, error)
}

type calendarServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarServiceClient(cc grpc.ClientConnInterface) CalendarServiceClient {
	return &calendarServiceClient{cc}
}

func (c *calendarServiceClient) CreateCalendar(ctx context.Context, in *CreateCalendarRequest, opts ...grpc.CallOption) (*Calendar, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Calendar)
	err := c.cc.Invoke(ctx, CalendarService_CreateCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ListCalendars(ctx context.Context, in *ListCalendarsRequest, opts ...grpc.CallOption) (*ListCalendarsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCalendarsResponse)
	err := c.cc.Invoke(ctx, CalendarService_ListCalendars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) UpdateCalendar(ctx context.Context, in *UpdateCalendarRequest, opts ...grpc.CallOption) (*Calendar, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Calendar)
	err := c.cc.Invoke(ctx, CalendarService_UpdateCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) DeleteCalendar(ctx context.Context, in *DeleteCalendarRequest, opts ...grpc.CallOption) (*Calendar, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Calendar)
	err := c.cc.Invoke(ctx, CalendarService_DeleteCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) GetCalendar(ctx context.Context, in *GetCalendarRequest, opts ...grpc.CallOption) (*Calendar, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Calendar)
	err := c.cc.Invoke(ctx, CalendarService_GetCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServiceServer is the server API for CalendarService service.
// All implementations must embed UnimplementedCalendarServiceServer
// for forward compatibility.
//
// the calendar service
type CalendarServiceServer interface {
	// create a calendar
	CreateCalendar(context.Context, *CreateCalendarRequest) (*Calendar, error)
	// list calendars
	ListCalendars(context.Context, *ListCalendarsRequest) (*ListCalendarsResponse, error)
	// update a calendar
	UpdateCalendar(context.Context, *UpdateCalendarRequest) (*Calendar, error)
	// delete` a calendar
	DeleteCalendar(context.Context, *DeleteCalendarRequest) (*Calendar, error)
	// get a calendar
	GetCalendar(context.Context, *GetCalendarRequest) (*Calendar, error)
	mustEmbedUnimplementedCalendarServiceServer()
}

// UnimplementedCalendarServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCalendarServiceServer struct{}

func (UnimplementedCalendarServiceServer) CreateCalendar(context.Context, *CreateCalendarRequest) (*Calendar, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) ListCalendars(context.Context, *ListCalendarsRequest) (*ListCalendarsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCalendars not implemented")
}
func (UnimplementedCalendarServiceServer) UpdateCalendar(context.Context, *UpdateCalendarRequest) (*Calendar, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) DeleteCalendar(context.Context, *DeleteCalendarRequest) (*Calendar, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) GetCalendar(context.Context, *GetCalendarRequest) (*Calendar, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) mustEmbedUnimplementedCalendarServiceServer() {}
func (UnimplementedCalendarServiceServer) testEmbeddedByValue()                         {}

// UnsafeCalendarServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServiceServer will
// result in compilation errors.
type UnsafeCalendarServiceServer interface {
	mustEmbedUnimplementedCalendarServiceServer()
}

func RegisterCalendarServiceServer(s grpc.ServiceRegistrar, srv CalendarServiceServer) {
	// If the following call pancis, it indicates UnimplementedCalendarServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CalendarService_ServiceDesc, srv)
}

func _CalendarService_CreateCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).CreateCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_CreateCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).CreateCalendar(ctx, req.(*CreateCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ListCalendars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCalendarsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListCalendars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_ListCalendars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListCalendars(ctx, req.(*ListCalendarsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_UpdateCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).UpdateCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_UpdateCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).UpdateCalendar(ctx, req.(*UpdateCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_DeleteCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).DeleteCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_DeleteCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).DeleteCalendar(ctx, req.(*DeleteCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_GetCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).GetCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_GetCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).GetCalendar(ctx, req.(*GetCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CalendarService_ServiceDesc is the grpc.ServiceDesc for CalendarService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CalendarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.calendars.calendar.v1alpha1.CalendarService",
	HandlerType: (*CalendarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCalendar",
			Handler:    _CalendarService_CreateCalendar_Handler,
		},
		{
			MethodName: "ListCalendars",
			Handler:    _CalendarService_ListCalendars_Handler,
		},
		{
			MethodName: "UpdateCalendar",
			Handler:    _CalendarService_UpdateCalendar_Handler,
		},
		{
			MethodName: "DeleteCalendar",
			Handler:    _CalendarService_DeleteCalendar_Handler,
		},
		{
			MethodName: "GetCalendar",
			Handler:    _CalendarService_GetCalendar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/calendars/calendar/v1alpha1/calendar.proto",
}
