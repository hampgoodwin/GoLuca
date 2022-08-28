// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: goluca/model/v1/account.proto

package modelv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AccountType int32

const (
	AccountType_ACCOUNT_TYPE_UNSPECIFIED AccountType = 0
	AccountType_ACCOUNT_TYPE_ASSET       AccountType = 1
	AccountType_ACCOUNT_TYPE_LIABILITY   AccountType = 2
	AccountType_ACCOUNT_TYPE_EQUITY      AccountType = 3
	AccountType_ACCOUNT_TYPE_REVENUE     AccountType = 4
	AccountType_ACCOUNT_TYPE_EXPENSE     AccountType = 5
	AccountType_ACCOUNT_TYPE_GAIN        AccountType = 6
	AccountType_ACCOUNT_TYPE_LOSS        AccountType = 7
)

// Enum value maps for AccountType.
var (
	AccountType_name = map[int32]string{
		0: "ACCOUNT_TYPE_UNSPECIFIED",
		1: "ACCOUNT_TYPE_ASSET",
		2: "ACCOUNT_TYPE_LIABILITY",
		3: "ACCOUNT_TYPE_EQUITY",
		4: "ACCOUNT_TYPE_REVENUE",
		5: "ACCOUNT_TYPE_EXPENSE",
		6: "ACCOUNT_TYPE_GAIN",
		7: "ACCOUNT_TYPE_LOSS",
	}
	AccountType_value = map[string]int32{
		"ACCOUNT_TYPE_UNSPECIFIED": 0,
		"ACCOUNT_TYPE_ASSET":       1,
		"ACCOUNT_TYPE_LIABILITY":   2,
		"ACCOUNT_TYPE_EQUITY":      3,
		"ACCOUNT_TYPE_REVENUE":     4,
		"ACCOUNT_TYPE_EXPENSE":     5,
		"ACCOUNT_TYPE_GAIN":        6,
		"ACCOUNT_TYPE_LOSS":        7,
	}
)

func (x AccountType) Enum() *AccountType {
	p := new(AccountType)
	*p = x
	return p
}

func (x AccountType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccountType) Descriptor() protoreflect.EnumDescriptor {
	return file_goluca_model_v1_account_proto_enumTypes[0].Descriptor()
}

func (AccountType) Type() protoreflect.EnumType {
	return &file_goluca_model_v1_account_proto_enumTypes[0]
}

func (x AccountType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccountType.Descriptor instead.
func (AccountType) EnumDescriptor() ([]byte, []int) {
	return file_goluca_model_v1_account_proto_rawDescGZIP(), []int{0}
}

type Basis int32

const (
	Basis_BASIS_UNSPECIFIED Basis = 0
	Basis_BASIS_DEBIT       Basis = 1
	Basis_BASIS_CREDIT      Basis = 2
)

// Enum value maps for Basis.
var (
	Basis_name = map[int32]string{
		0: "BASIS_UNSPECIFIED",
		1: "BASIS_DEBIT",
		2: "BASIS_CREDIT",
	}
	Basis_value = map[string]int32{
		"BASIS_UNSPECIFIED": 0,
		"BASIS_DEBIT":       1,
		"BASIS_CREDIT":      2,
	}
)

func (x Basis) Enum() *Basis {
	p := new(Basis)
	*p = x
	return p
}

func (x Basis) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Basis) Descriptor() protoreflect.EnumDescriptor {
	return file_goluca_model_v1_account_proto_enumTypes[1].Descriptor()
}

func (Basis) Type() protoreflect.EnumType {
	return &file_goluca_model_v1_account_proto_enumTypes[1]
}

func (x Basis) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Basis.Descriptor instead.
func (Basis) EnumDescriptor() ([]byte, []int) {
	return file_goluca_model_v1_account_proto_rawDescGZIP(), []int{1}
}

type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ParentId  *string                `protobuf:"bytes,2,opt,name=parent_id,json=parentId,proto3,oneof" json:"parent_id,omitempty"`
	Name      string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Type      AccountType            `protobuf:"varint,4,opt,name=type,proto3,enum=goluca.model.v1.AccountType" json:"type,omitempty"`
	Basis     Basis                  `protobuf:"varint,5,opt,name=basis,proto3,enum=goluca.model.v1.Basis" json:"basis,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goluca_model_v1_account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_goluca_model_v1_account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_goluca_model_v1_account_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Account) GetParentId() string {
	if x != nil && x.ParentId != nil {
		return *x.ParentId
	}
	return ""
}

func (x *Account) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Account) GetType() AccountType {
	if x != nil {
		return x.Type
	}
	return AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

func (x *Account) GetBasis() Basis {
	if x != nil {
		return x.Basis
	}
	return Basis_BASIS_UNSPECIFIED
}

func (x *Account) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_goluca_model_v1_account_proto protoreflect.FileDescriptor

var file_goluca_model_v1_account_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x67, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76,
	0x31, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0f, 0x67, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xf8, 0x01, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a,
	0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x62, 0x61, 0x73, 0x69, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x73, 0x52, 0x05, 0x62, 0x61,
	0x73, 0x69, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x0c,
	0x0a, 0x0a, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x2a, 0xda, 0x01, 0x0a,
	0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x18,
	0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x41, 0x43,
	0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x53, 0x53, 0x45, 0x54,
	0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4c, 0x49, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x10, 0x02, 0x12, 0x17,
	0x0a, 0x13, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45,
	0x51, 0x55, 0x49, 0x54, 0x59, 0x10, 0x03, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x43, 0x43, 0x4f, 0x55,
	0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x56, 0x45, 0x4e, 0x55, 0x45, 0x10,
	0x04, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x45, 0x58, 0x50, 0x45, 0x4e, 0x53, 0x45, 0x10, 0x05, 0x12, 0x15, 0x0a, 0x11, 0x41,
	0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x41, 0x49, 0x4e,
	0x10, 0x06, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4c, 0x4f, 0x53, 0x53, 0x10, 0x07, 0x2a, 0x41, 0x0a, 0x05, 0x42, 0x61, 0x73,
	0x69, 0x73, 0x12, 0x15, 0x0a, 0x11, 0x42, 0x41, 0x53, 0x49, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x42, 0x41, 0x53,
	0x49, 0x53, 0x5f, 0x44, 0x45, 0x42, 0x49, 0x54, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x42, 0x41,
	0x53, 0x49, 0x53, 0x5f, 0x43, 0x52, 0x45, 0x44, 0x49, 0x54, 0x10, 0x02, 0x42, 0xc5, 0x01, 0x0a,
	0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x68, 0x61, 0x6d, 0x70, 0x67, 0x6f, 0x6f, 0x64, 0x77, 0x69, 0x6e, 0x2f, 0x47, 0x6f, 0x4c,
	0x75, 0x63, 0x61, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f,
	0x2f, 0x67, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31,
	0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x47, 0x4d, 0x58, 0xaa, 0x02,
	0x0f, 0x47, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0f, 0x47, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x5c, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x1b, 0x47, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x5c, 0x4d, 0x6f, 0x64, 0x65,
	0x6c, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x11, 0x47, 0x6f, 0x6c, 0x75, 0x63, 0x61, 0x3a, 0x3a, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_goluca_model_v1_account_proto_rawDescOnce sync.Once
	file_goluca_model_v1_account_proto_rawDescData = file_goluca_model_v1_account_proto_rawDesc
)

func file_goluca_model_v1_account_proto_rawDescGZIP() []byte {
	file_goluca_model_v1_account_proto_rawDescOnce.Do(func() {
		file_goluca_model_v1_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_goluca_model_v1_account_proto_rawDescData)
	})
	return file_goluca_model_v1_account_proto_rawDescData
}

var file_goluca_model_v1_account_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_goluca_model_v1_account_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_goluca_model_v1_account_proto_goTypes = []interface{}{
	(AccountType)(0),              // 0: goluca.model.v1.AccountType
	(Basis)(0),                    // 1: goluca.model.v1.Basis
	(*Account)(nil),               // 2: goluca.model.v1.Account
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_goluca_model_v1_account_proto_depIdxs = []int32{
	0, // 0: goluca.model.v1.Account.type:type_name -> goluca.model.v1.AccountType
	1, // 1: goluca.model.v1.Account.basis:type_name -> goluca.model.v1.Basis
	3, // 2: goluca.model.v1.Account.created_at:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_goluca_model_v1_account_proto_init() }
func file_goluca_model_v1_account_proto_init() {
	if File_goluca_model_v1_account_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_goluca_model_v1_account_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_goluca_model_v1_account_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_goluca_model_v1_account_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_goluca_model_v1_account_proto_goTypes,
		DependencyIndexes: file_goluca_model_v1_account_proto_depIdxs,
		EnumInfos:         file_goluca_model_v1_account_proto_enumTypes,
		MessageInfos:      file_goluca_model_v1_account_proto_msgTypes,
	}.Build()
	File_goluca_model_v1_account_proto = out.File
	file_goluca_model_v1_account_proto_rawDesc = nil
	file_goluca_model_v1_account_proto_goTypes = nil
	file_goluca_model_v1_account_proto_depIdxs = nil
}