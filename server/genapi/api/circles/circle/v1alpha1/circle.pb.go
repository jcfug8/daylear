// (-- api-linter: core::0191::java-outer-classname=disabled
//     aip.dev/not-precedent: I don't care about java. --)
// (-- api-linter: core::0191::java-multiple-files=disabled
//     aip.dev/not-precedent: I don't care about java. --)
// (-- api-linter: core::0191::java-package=disabled
//     aip.dev/not-precedent: I don't care about java. --)

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: api/circles/circle/v1alpha1/circle.proto

package circlev1alpha1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// the main user circle
type Circle struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the name of the circle
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// the public name of the circle
	PublicName string `protobuf:"bytes,2,opt,name=public_name,json=publicName,proto3" json:"public_name,omitempty"`
	// the title of the circle
	Title string `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	// if the circle is public
	IsPublic      bool `protobuf:"varint,4,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Circle) Reset() {
	*x = Circle{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Circle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Circle) ProtoMessage() {}

func (x *Circle) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Circle.ProtoReflect.Descriptor instead.
func (*Circle) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{0}
}

func (x *Circle) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Circle) GetPublicName() string {
	if x != nil {
		return x.PublicName
	}
	return ""
}

func (x *Circle) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Circle) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

// the request to create a circle
type CreateCircleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the circle to create
	Circle *Circle `protobuf:"bytes,2,opt,name=circle,proto3" json:"circle,omitempty"`
	// the id of the circle
	CircleId      string `protobuf:"bytes,3,opt,name=circle_id,json=circleId,proto3" json:"circle_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateCircleRequest) Reset() {
	*x = CreateCircleRequest{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateCircleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCircleRequest) ProtoMessage() {}

func (x *CreateCircleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCircleRequest.ProtoReflect.Descriptor instead.
func (*CreateCircleRequest) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{1}
}

func (x *CreateCircleRequest) GetCircle() *Circle {
	if x != nil {
		return x.Circle
	}
	return nil
}

func (x *CreateCircleRequest) GetCircleId() string {
	if x != nil {
		return x.CircleId
	}
	return ""
}

// the request to list circles
type ListCirclesRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the page size
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// the page token
	PageToken string `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// used to specify the filter
	Filter        string `protobuf:"bytes,4,opt,name=filter,proto3" json:"filter,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCirclesRequest) Reset() {
	*x = ListCirclesRequest{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCirclesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCirclesRequest) ProtoMessage() {}

func (x *ListCirclesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCirclesRequest.ProtoReflect.Descriptor instead.
func (*ListCirclesRequest) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{2}
}

func (x *ListCirclesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListCirclesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListCirclesRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

// the response to list circles
type ListCirclesResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the circles
	Circles []*Circle `protobuf:"bytes,1,rep,name=circles,proto3" json:"circles,omitempty"`
	// the next page token
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCirclesResponse) Reset() {
	*x = ListCirclesResponse{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCirclesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCirclesResponse) ProtoMessage() {}

func (x *ListCirclesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCirclesResponse.ProtoReflect.Descriptor instead.
func (*ListCirclesResponse) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{3}
}

func (x *ListCirclesResponse) GetCircles() []*Circle {
	if x != nil {
		return x.Circles
	}
	return nil
}

func (x *ListCirclesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// the request to update a circle
type UpdateCircleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the circle to update
	Circle *Circle `protobuf:"bytes,1,opt,name=circle,proto3" json:"circle,omitempty"`
	// the fields to update
	UpdateMask    *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCircleRequest) Reset() {
	*x = UpdateCircleRequest{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCircleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCircleRequest) ProtoMessage() {}

func (x *UpdateCircleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCircleRequest.ProtoReflect.Descriptor instead.
func (*UpdateCircleRequest) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateCircleRequest) GetCircle() *Circle {
	if x != nil {
		return x.Circle
	}
	return nil
}

func (x *UpdateCircleRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

// the request to delete a circle
type DeleteCircleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the name of the circle
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCircleRequest) Reset() {
	*x = DeleteCircleRequest{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCircleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCircleRequest) ProtoMessage() {}

func (x *DeleteCircleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCircleRequest.ProtoReflect.Descriptor instead.
func (*DeleteCircleRequest) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteCircleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// the request to get a circle
type GetCircleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the name of the circle
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCircleRequest) Reset() {
	*x = GetCircleRequest{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCircleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCircleRequest) ProtoMessage() {}

func (x *GetCircleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCircleRequest.ProtoReflect.Descriptor instead.
func (*GetCircleRequest) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{6}
}

func (x *GetCircleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// the request to share a circle
type ShareCircleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the name of the circle
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShareCircleRequest) Reset() {
	*x = ShareCircleRequest{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShareCircleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShareCircleRequest) ProtoMessage() {}

func (x *ShareCircleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShareCircleRequest.ProtoReflect.Descriptor instead.
func (*ShareCircleRequest) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{7}
}

func (x *ShareCircleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// the response to share a circle
type ShareCircleResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the circle that was shared
	Circle        *Circle `protobuf:"bytes,1,opt,name=circle,proto3" json:"circle,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShareCircleResponse) Reset() {
	*x = ShareCircleResponse{}
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShareCircleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShareCircleResponse) ProtoMessage() {}

func (x *ShareCircleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_circles_circle_v1alpha1_circle_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShareCircleResponse.ProtoReflect.Descriptor instead.
func (*ShareCircleResponse) Descriptor() ([]byte, []int) {
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP(), []int{8}
}

func (x *ShareCircleResponse) GetCircle() *Circle {
	if x != nil {
		return x.Circle
	}
	return nil
}

var File_api_circles_circle_v1alpha1_circle_proto protoreflect.FileDescriptor

const file_api_circles_circle_v1alpha1_circle_proto_rawDesc = "" +
	"\n" +
	"(api/circles/circle/v1alpha1/circle.proto\x12\x1bapi.circles.circle.v1alpha1\x1a\x1cgoogle/api/annotations.proto\x1a\x17google/api/client.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x19google/api/resource.proto\x1a google/protobuf/field_mask.proto\"\xd0\x01\n" +
	"\x06Circle\x12\x17\n" +
	"\x04name\x18\x01 \x01(\tB\x03\xe0A\bR\x04name\x12$\n" +
	"\vpublic_name\x18\x02 \x01(\tB\x03\xe0A\x03R\n" +
	"publicName\x12\x19\n" +
	"\x05title\x18\x03 \x01(\tB\x03\xe0A\x02R\x05title\x12 \n" +
	"\tis_public\x18\x04 \x01(\bB\x03\xe0A\x01R\bisPublic:J\xeaAG\n" +
	"\"api.circles.circle.v1alpha1/Circle\x12\x10circles/{circle}*\acircles2\x06circle\"y\n" +
	"\x13CreateCircleRequest\x12@\n" +
	"\x06circle\x18\x02 \x01(\v2#.api.circles.circle.v1alpha1.CircleB\x03\xe0A\x02R\x06circle\x12 \n" +
	"\tcircle_id\x18\x03 \x01(\tB\x03\xe0A\x02R\bcircleId\"w\n" +
	"\x12ListCirclesRequest\x12 \n" +
	"\tpage_size\x18\x02 \x01(\x05B\x03\xe0A\x01R\bpageSize\x12\"\n" +
	"\n" +
	"page_token\x18\x03 \x01(\tB\x03\xe0A\x01R\tpageToken\x12\x1b\n" +
	"\x06filter\x18\x04 \x01(\tB\x03\xe0A\x01R\x06filter\"|\n" +
	"\x13ListCirclesResponse\x12=\n" +
	"\acircles\x18\x01 \x03(\v2#.api.circles.circle.v1alpha1.CircleR\acircles\x12&\n" +
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\"\x99\x01\n" +
	"\x13UpdateCircleRequest\x12@\n" +
	"\x06circle\x18\x01 \x01(\v2#.api.circles.circle.v1alpha1.CircleB\x03\xe0A\x02R\x06circle\x12@\n" +
	"\vupdate_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskB\x03\xe0A\x01R\n" +
	"updateMask\"U\n" +
	"\x13DeleteCircleRequest\x12>\n" +
	"\x04name\x18\x01 \x01(\tB*\xe0A\x02\xfaA$\n" +
	"\"api.circles.circle.v1alpha1/CircleR\x04name\"R\n" +
	"\x10GetCircleRequest\x12>\n" +
	"\x04name\x18\x01 \x01(\tB*\xe0A\x02\xfaA$\n" +
	"\"api.circles.circle.v1alpha1/CircleR\x04name\"T\n" +
	"\x12ShareCircleRequest\x12>\n" +
	"\x04name\x18\x01 \x01(\tB*\xe0A\x02\xfaA$\n" +
	"\"api.circles.circle.v1alpha1/CircleR\x04name\"R\n" +
	"\x13ShareCircleResponse\x12;\n" +
	"\x06circle\x18\x01 \x01(\v2#.api.circles.circle.v1alpha1.CircleR\x06circle2\xed\a\n" +
	"\rCircleService\x12\xa3\x01\n" +
	"\fCreateCircle\x120.api.circles.circle.v1alpha1.CreateCircleRequest\x1a#.api.circles.circle.v1alpha1.Circle\"<\xdaA\x10circle,circle_id\x82\xd3\xe4\x93\x02#:\x06circle\"\x19/circles/v1alpha1/circles\x12\x93\x01\n" +
	"\vListCircles\x12/.api.circles.circle.v1alpha1.ListCirclesRequest\x1a0.api.circles.circle.v1alpha1.ListCirclesResponse\"!\x82\xd3\xe4\x93\x02\x1b\x12\x19/circles/v1alpha1/circles\x12\xb5\x01\n" +
	"\fUpdateCircle\x120.api.circles.circle.v1alpha1.UpdateCircleRequest\x1a#.api.circles.circle.v1alpha1.Circle\"N\xdaA\x12circle,update_mask\x82\xd3\xe4\x93\x023:\x06circle2)/circles/v1alpha1/{circle.name=circles/*}\x12\x98\x01\n" +
	"\fDeleteCircle\x120.api.circles.circle.v1alpha1.DeleteCircleRequest\x1a#.api.circles.circle.v1alpha1.Circle\"1\xdaA\x04name\x82\xd3\xe4\x93\x02$*\"/circles/v1alpha1/{name=circles/*}\x12\x92\x01\n" +
	"\tGetCircle\x12-.api.circles.circle.v1alpha1.GetCircleRequest\x1a#.api.circles.circle.v1alpha1.Circle\"1\xdaA\x04name\x82\xd3\xe4\x93\x02$\x12\"/circles/v1alpha1/{name=circles/*}\x12\xb7\x01\n" +
	"\vShareCircle\x12/.api.circles.circle.v1alpha1.ShareCircleRequest\x1a0.api.circles.circle.v1alpha1.ShareCircleResponse\"E\xdaA\x0fname,recipients\x82\xd3\xe4\x93\x02-:\x01*\"(/circles/v1alpha1/{name=circles/*}:shareB\x91\x02\n" +
	"\x1fcom.api.circles.circle.v1alpha1B\vCircleProtoP\x01ZRgithub.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1;circlev1alpha1\xa2\x02\x03ACC\xaa\x02\x1bApi.Circles.Circle.V1alpha1\xca\x02\x1bApi\\Circles\\Circle\\V1alpha1\xe2\x02'Api\\Circles\\Circle\\V1alpha1\\GPBMetadata\xea\x02\x1eApi::Circles::Circle::V1alpha1b\x06proto3"

var (
	file_api_circles_circle_v1alpha1_circle_proto_rawDescOnce sync.Once
	file_api_circles_circle_v1alpha1_circle_proto_rawDescData []byte
)

func file_api_circles_circle_v1alpha1_circle_proto_rawDescGZIP() []byte {
	file_api_circles_circle_v1alpha1_circle_proto_rawDescOnce.Do(func() {
		file_api_circles_circle_v1alpha1_circle_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_circles_circle_v1alpha1_circle_proto_rawDesc), len(file_api_circles_circle_v1alpha1_circle_proto_rawDesc)))
	})
	return file_api_circles_circle_v1alpha1_circle_proto_rawDescData
}

var file_api_circles_circle_v1alpha1_circle_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_circles_circle_v1alpha1_circle_proto_goTypes = []any{
	(*Circle)(nil),                // 0: api.circles.circle.v1alpha1.Circle
	(*CreateCircleRequest)(nil),   // 1: api.circles.circle.v1alpha1.CreateCircleRequest
	(*ListCirclesRequest)(nil),    // 2: api.circles.circle.v1alpha1.ListCirclesRequest
	(*ListCirclesResponse)(nil),   // 3: api.circles.circle.v1alpha1.ListCirclesResponse
	(*UpdateCircleRequest)(nil),   // 4: api.circles.circle.v1alpha1.UpdateCircleRequest
	(*DeleteCircleRequest)(nil),   // 5: api.circles.circle.v1alpha1.DeleteCircleRequest
	(*GetCircleRequest)(nil),      // 6: api.circles.circle.v1alpha1.GetCircleRequest
	(*ShareCircleRequest)(nil),    // 7: api.circles.circle.v1alpha1.ShareCircleRequest
	(*ShareCircleResponse)(nil),   // 8: api.circles.circle.v1alpha1.ShareCircleResponse
	(*fieldmaskpb.FieldMask)(nil), // 9: google.protobuf.FieldMask
}
var file_api_circles_circle_v1alpha1_circle_proto_depIdxs = []int32{
	0,  // 0: api.circles.circle.v1alpha1.CreateCircleRequest.circle:type_name -> api.circles.circle.v1alpha1.Circle
	0,  // 1: api.circles.circle.v1alpha1.ListCirclesResponse.circles:type_name -> api.circles.circle.v1alpha1.Circle
	0,  // 2: api.circles.circle.v1alpha1.UpdateCircleRequest.circle:type_name -> api.circles.circle.v1alpha1.Circle
	9,  // 3: api.circles.circle.v1alpha1.UpdateCircleRequest.update_mask:type_name -> google.protobuf.FieldMask
	0,  // 4: api.circles.circle.v1alpha1.ShareCircleResponse.circle:type_name -> api.circles.circle.v1alpha1.Circle
	1,  // 5: api.circles.circle.v1alpha1.CircleService.CreateCircle:input_type -> api.circles.circle.v1alpha1.CreateCircleRequest
	2,  // 6: api.circles.circle.v1alpha1.CircleService.ListCircles:input_type -> api.circles.circle.v1alpha1.ListCirclesRequest
	4,  // 7: api.circles.circle.v1alpha1.CircleService.UpdateCircle:input_type -> api.circles.circle.v1alpha1.UpdateCircleRequest
	5,  // 8: api.circles.circle.v1alpha1.CircleService.DeleteCircle:input_type -> api.circles.circle.v1alpha1.DeleteCircleRequest
	6,  // 9: api.circles.circle.v1alpha1.CircleService.GetCircle:input_type -> api.circles.circle.v1alpha1.GetCircleRequest
	7,  // 10: api.circles.circle.v1alpha1.CircleService.ShareCircle:input_type -> api.circles.circle.v1alpha1.ShareCircleRequest
	0,  // 11: api.circles.circle.v1alpha1.CircleService.CreateCircle:output_type -> api.circles.circle.v1alpha1.Circle
	3,  // 12: api.circles.circle.v1alpha1.CircleService.ListCircles:output_type -> api.circles.circle.v1alpha1.ListCirclesResponse
	0,  // 13: api.circles.circle.v1alpha1.CircleService.UpdateCircle:output_type -> api.circles.circle.v1alpha1.Circle
	0,  // 14: api.circles.circle.v1alpha1.CircleService.DeleteCircle:output_type -> api.circles.circle.v1alpha1.Circle
	0,  // 15: api.circles.circle.v1alpha1.CircleService.GetCircle:output_type -> api.circles.circle.v1alpha1.Circle
	8,  // 16: api.circles.circle.v1alpha1.CircleService.ShareCircle:output_type -> api.circles.circle.v1alpha1.ShareCircleResponse
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_circles_circle_v1alpha1_circle_proto_init() }
func file_api_circles_circle_v1alpha1_circle_proto_init() {
	if File_api_circles_circle_v1alpha1_circle_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_circles_circle_v1alpha1_circle_proto_rawDesc), len(file_api_circles_circle_v1alpha1_circle_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_circles_circle_v1alpha1_circle_proto_goTypes,
		DependencyIndexes: file_api_circles_circle_v1alpha1_circle_proto_depIdxs,
		MessageInfos:      file_api_circles_circle_v1alpha1_circle_proto_msgTypes,
	}.Build()
	File_api_circles_circle_v1alpha1_circle_proto = out.File
	file_api_circles_circle_v1alpha1_circle_proto_goTypes = nil
	file_api_circles_circle_v1alpha1_circle_proto_depIdxs = nil
}
