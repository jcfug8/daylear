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
// source: api/users/user/v1alpha1/public_user.proto

package userv1alpha1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/fieldmaskpb"
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

// the main public user object
type PublicUser struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the name of the user
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// username
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	// the given name of the user
	GivenName string `protobuf:"bytes,3,opt,name=given_name,json=givenName,proto3" json:"given_name,omitempty"`
	// the family name of the user
	FamilyName    string `protobuf:"bytes,4,opt,name=family_name,json=familyName,proto3" json:"family_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PublicUser) Reset() {
	*x = PublicUser{}
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PublicUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicUser) ProtoMessage() {}

func (x *PublicUser) ProtoReflect() protoreflect.Message {
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicUser.ProtoReflect.Descriptor instead.
func (*PublicUser) Descriptor() ([]byte, []int) {
	return file_api_users_user_v1alpha1_public_user_proto_rawDescGZIP(), []int{0}
}

func (x *PublicUser) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PublicUser) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *PublicUser) GetGivenName() string {
	if x != nil {
		return x.GivenName
	}
	return ""
}

func (x *PublicUser) GetFamilyName() string {
	if x != nil {
		return x.FamilyName
	}
	return ""
}

// the request to list public users
type ListPublicUsersRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// returned page
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// used to specify the page token
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// used to specify the filter
	Filter        string `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListPublicUsersRequest) Reset() {
	*x = ListPublicUsersRequest{}
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListPublicUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPublicUsersRequest) ProtoMessage() {}

func (x *ListPublicUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPublicUsersRequest.ProtoReflect.Descriptor instead.
func (*ListPublicUsersRequest) Descriptor() ([]byte, []int) {
	return file_api_users_user_v1alpha1_public_user_proto_rawDescGZIP(), []int{1}
}

func (x *ListPublicUsersRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListPublicUsersRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListPublicUsersRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

// the response to list users
type ListPublicUsersResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the users
	PublicUsers []*PublicUser `protobuf:"bytes,1,rep,name=public_users,json=publicUsers,proto3" json:"public_users,omitempty"`
	// the next page token
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListPublicUsersResponse) Reset() {
	*x = ListPublicUsersResponse{}
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListPublicUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPublicUsersResponse) ProtoMessage() {}

func (x *ListPublicUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPublicUsersResponse.ProtoReflect.Descriptor instead.
func (*ListPublicUsersResponse) Descriptor() ([]byte, []int) {
	return file_api_users_user_v1alpha1_public_user_proto_rawDescGZIP(), []int{2}
}

func (x *ListPublicUsersResponse) GetPublicUsers() []*PublicUser {
	if x != nil {
		return x.PublicUsers
	}
	return nil
}

func (x *ListPublicUsersResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// the request to get a public user
type GetPublicUserRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the name of the public user to get
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPublicUserRequest) Reset() {
	*x = GetPublicUserRequest{}
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPublicUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPublicUserRequest) ProtoMessage() {}

func (x *GetPublicUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_users_user_v1alpha1_public_user_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPublicUserRequest.ProtoReflect.Descriptor instead.
func (*GetPublicUserRequest) Descriptor() ([]byte, []int) {
	return file_api_users_user_v1alpha1_public_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetPublicUserRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_api_users_user_v1alpha1_public_user_proto protoreflect.FileDescriptor

const file_api_users_user_v1alpha1_public_user_proto_rawDesc = "" +
	"\n" +
	")api/users/user/v1alpha1/public_user.proto\x12\x17api.users.user.v1alpha1\x1a\x1cgoogle/api/annotations.proto\x1a\x17google/api/client.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x19google/api/resource.proto\x1a google/protobuf/field_mask.proto\"\xed\x01\n" +
	"\n" +
	"PublicUser\x12\x17\n" +
	"\x04name\x18\x01 \x01(\tB\x03\xe0A\bR\x04name\x12\x1f\n" +
	"\busername\x18\x02 \x01(\tB\x03\xe0A\x03R\busername\x12\"\n" +
	"\n" +
	"given_name\x18\x03 \x01(\tB\x03\xe0A\x03R\tgivenName\x12$\n" +
	"\vfamily_name\x18\x04 \x01(\tB\x03\xe0A\x03R\n" +
	"familyName:[\xeaAX\n" +
	"\"api.users.user.v1alpha1/PublicUser\x12\x19publicUsers/{public_user}*\vpublicUsers2\n" +
	"publicUser\"{\n" +
	"\x16ListPublicUsersRequest\x12 \n" +
	"\tpage_size\x18\x01 \x01(\x05B\x03\xe0A\x01R\bpageSize\x12\"\n" +
	"\n" +
	"page_token\x18\x02 \x01(\tB\x03\xe0A\x01R\tpageToken\x12\x1b\n" +
	"\x06filter\x18\x03 \x01(\tB\x03\xe0A\x01R\x06filter\"\x89\x01\n" +
	"\x17ListPublicUsersResponse\x12F\n" +
	"\fpublic_users\x18\x01 \x03(\v2#.api.users.user.v1alpha1.PublicUserR\vpublicUsers\x12&\n" +
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\"V\n" +
	"\x14GetPublicUserRequest\x12>\n" +
	"\x04name\x18\x01 \x01(\tB*\xe0A\x02\xfaA$\n" +
	"\"api.users.user.v1alpha1/PublicUserR\x04name2\xca\x02\n" +
	"\x11PublicUserService\x12\x99\x01\n" +
	"\x0fListPublicUsers\x12/.api.users.user.v1alpha1.ListPublicUsersRequest\x1a0.api.users.user.v1alpha1.ListPublicUsersResponse\"#\x82\xd3\xe4\x93\x02\x1d\x12\x1b/users/v1alpha1/publicUsers\x12\x98\x01\n" +
	"\rGetPublicUser\x12-.api.users.user.v1alpha1.GetPublicUserRequest\x1a#.api.users.user.v1alpha1.PublicUser\"3\xdaA\x04name\x82\xd3\xe4\x93\x02&\x12$/users/v1alpha1/{name=publicUsers/*}B\xfb\x01\n" +
	"\x1bcom.api.users.user.v1alpha1B\x0fPublicUserProtoP\x01ZLgithub.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1;userv1alpha1\xa2\x02\x03AUU\xaa\x02\x17Api.Users.User.V1alpha1\xca\x02\x17Api\\Users\\User\\V1alpha1\xe2\x02#Api\\Users\\User\\V1alpha1\\GPBMetadata\xea\x02\x1aApi::Users::User::V1alpha1b\x06proto3"

var (
	file_api_users_user_v1alpha1_public_user_proto_rawDescOnce sync.Once
	file_api_users_user_v1alpha1_public_user_proto_rawDescData []byte
)

func file_api_users_user_v1alpha1_public_user_proto_rawDescGZIP() []byte {
	file_api_users_user_v1alpha1_public_user_proto_rawDescOnce.Do(func() {
		file_api_users_user_v1alpha1_public_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_users_user_v1alpha1_public_user_proto_rawDesc), len(file_api_users_user_v1alpha1_public_user_proto_rawDesc)))
	})
	return file_api_users_user_v1alpha1_public_user_proto_rawDescData
}

var file_api_users_user_v1alpha1_public_user_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_users_user_v1alpha1_public_user_proto_goTypes = []any{
	(*PublicUser)(nil),              // 0: api.users.user.v1alpha1.PublicUser
	(*ListPublicUsersRequest)(nil),  // 1: api.users.user.v1alpha1.ListPublicUsersRequest
	(*ListPublicUsersResponse)(nil), // 2: api.users.user.v1alpha1.ListPublicUsersResponse
	(*GetPublicUserRequest)(nil),    // 3: api.users.user.v1alpha1.GetPublicUserRequest
}
var file_api_users_user_v1alpha1_public_user_proto_depIdxs = []int32{
	0, // 0: api.users.user.v1alpha1.ListPublicUsersResponse.public_users:type_name -> api.users.user.v1alpha1.PublicUser
	1, // 1: api.users.user.v1alpha1.PublicUserService.ListPublicUsers:input_type -> api.users.user.v1alpha1.ListPublicUsersRequest
	3, // 2: api.users.user.v1alpha1.PublicUserService.GetPublicUser:input_type -> api.users.user.v1alpha1.GetPublicUserRequest
	2, // 3: api.users.user.v1alpha1.PublicUserService.ListPublicUsers:output_type -> api.users.user.v1alpha1.ListPublicUsersResponse
	0, // 4: api.users.user.v1alpha1.PublicUserService.GetPublicUser:output_type -> api.users.user.v1alpha1.PublicUser
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_users_user_v1alpha1_public_user_proto_init() }
func file_api_users_user_v1alpha1_public_user_proto_init() {
	if File_api_users_user_v1alpha1_public_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_users_user_v1alpha1_public_user_proto_rawDesc), len(file_api_users_user_v1alpha1_public_user_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_users_user_v1alpha1_public_user_proto_goTypes,
		DependencyIndexes: file_api_users_user_v1alpha1_public_user_proto_depIdxs,
		MessageInfos:      file_api_users_user_v1alpha1_public_user_proto_msgTypes,
	}.Build()
	File_api_users_user_v1alpha1_public_user_proto = out.File
	file_api_users_user_v1alpha1_public_user_proto_goTypes = nil
	file_api_users_user_v1alpha1_public_user_proto_depIdxs = nil
}
