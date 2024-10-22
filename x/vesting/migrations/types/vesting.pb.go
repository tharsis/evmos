// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: eidon-chain/vesting/v1/vesting.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	github_com_cosmos_cosmos_sdk_x_auth_vesting_types "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	types "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ClawbackVestingAccount implements the VestingAccount interface. It provides
// an account that can hold contributions subject to "lockup" (like a
// PeriodicVestingAccount), or vesting which is subject to clawback
// of unvested tokens, or a combination (tokens vest, but are still locked).
type ClawbackVestingAccount struct {
	// base_vesting_account implements the VestingAccount interface. It contains
	// all the necessary fields needed for any vesting account implementation
	*types.BaseVestingAccount `protobuf:"bytes,1,opt,name=base_vesting_account,json=baseVestingAccount,proto3,embedded=base_vesting_account" json:"base_vesting_account,omitempty"`
	// funder_address specifies the account which can perform clawback
	FunderAddress string `protobuf:"bytes,2,opt,name=funder_address,json=funderAddress,proto3" json:"funder_address,omitempty"`
	// start_time defines the time at which the vesting period begins
	StartTime time.Time `protobuf:"bytes,3,opt,name=start_time,json=startTime,proto3,stdtime" json:"start_time"`
	// lockup_periods defines the unlocking schedule relative to the start_time
	LockupPeriods github_com_cosmos_cosmos_sdk_x_auth_vesting_types.Periods `protobuf:"bytes,4,rep,name=lockup_periods,json=lockupPeriods,proto3,castrepeated=github.com/cosmos/cosmos-sdk/x/auth/vesting/types.Periods" json:"lockup_periods"`
	// vesting_periods defines the vesting schedule relative to the start_time
	VestingPeriods github_com_cosmos_cosmos_sdk_x_auth_vesting_types.Periods `protobuf:"bytes,5,rep,name=vesting_periods,json=vestingPeriods,proto3,castrepeated=github.com/cosmos/cosmos-sdk/x/auth/vesting/types.Periods" json:"vesting_periods"`
}

func (m *ClawbackVestingAccount) Reset()      { *m = ClawbackVestingAccount{} }
func (*ClawbackVestingAccount) ProtoMessage() {}
func (*ClawbackVestingAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f1a3c86c0cebe5f, []int{0}
}
func (m *ClawbackVestingAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClawbackVestingAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClawbackVestingAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClawbackVestingAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClawbackVestingAccount.Merge(m, src)
}
func (m *ClawbackVestingAccount) XXX_Size() int {
	return m.Size()
}
func (m *ClawbackVestingAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_ClawbackVestingAccount.DiscardUnknown(m)
}

var xxx_messageInfo_ClawbackVestingAccount proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ClawbackVestingAccount)(nil), "eidon-chain.vesting.v1.ClawbackVestingAccount")
}

func init() { proto.RegisterFile("eidon-chain/vesting/v1/vesting.proto", fileDescriptor_5f1a3c86c0cebe5f) }

var fileDescriptor_5f1a3c86c0cebe5f = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x93, 0xcf, 0x6b, 0xd4, 0x40,
	0x14, 0xc7, 0x33, 0x76, 0x15, 0x3a, 0x65, 0x57, 0x0d, 0x45, 0x96, 0x1c, 0x26, 0x8b, 0x28, 0xac,
	0x0b, 0xce, 0xd8, 0x15, 0x05, 0xbd, 0x35, 0x82, 0x78, 0x53, 0x16, 0xf1, 0xe0, 0x25, 0xcc, 0x24,
	0xd3, 0x34, 0xec, 0x26, 0x13, 0x32, 0x93, 0x58, 0xff, 0x03, 0x11, 0xc4, 0x1e, 0x3d, 0xf6, 0x28,
	0x9e, 0xfa, 0x5f, 0xd8, 0xe3, 0x1e, 0x3d, 0xb5, 0xb2, 0x7b, 0xe8, 0xbf, 0x21, 0x99, 0x99, 0x48,
	0xea, 0x8f, 0x6b, 0x2f, 0xc9, 0xcb, 0xf7, 0xbd, 0x7c, 0xe7, 0xf3, 0xf2, 0x5e, 0x20, 0xe2, 0x75,
	0x26, 0x24, 0xa9, 0xb9, 0x54, 0x69, 0x9e, 0x90, 0x7a, 0xa7, 0x0d, 0x71, 0x51, 0x0a, 0x25, 0xdc,
	0x1b, 0x3a, 0x8f, 0x5b, 0xb1, 0xde, 0xf1, 0x6e, 0xd2, 0x2c, 0xcd, 0x05, 0xd1, 0x57, 0x53, 0xe4,
	0xdd, 0x89, 0x84, 0xbc, 0xe8, 0xc2, 0xb8, 0xa2, 0x7f, 0x58, 0x79, 0xdb, 0x89, 0x48, 0x84, 0x0e,
	0x49, 0x13, 0x59, 0xd5, 0x4f, 0x84, 0x48, 0x16, 0x9c, 0xe8, 0x27, 0x56, 0xed, 0x11, 0x95, 0x66,
	0x5c, 0x2a, 0x9a, 0x15, 0xa6, 0xe0, 0xf6, 0xf7, 0x1e, 0xbc, 0xf5, 0x6c, 0x41, 0xdf, 0x31, 0x1a,
	0xcd, 0xdf, 0x18, 0xc3, 0xdd, 0x28, 0x12, 0x55, 0xae, 0x5c, 0x06, 0xb7, 0x19, 0x95, 0x3c, 0xb4,
	0xe7, 0x84, 0xd4, 0xe8, 0x43, 0x30, 0x02, 0xe3, 0xad, 0xe9, 0x04, 0x1b, 0xac, 0x0e, 0xbc, 0xc6,
	0xc2, 0x01, 0x95, 0xfc, 0xa2, 0x53, 0xd0, 0x5b, 0x9e, 0xfa, 0x60, 0xe6, 0xb2, 0xbf, 0x32, 0xee,
	0x5d, 0x38, 0xd8, 0xab, 0xf2, 0x98, 0x97, 0x21, 0x8d, 0xe3, 0x92, 0x4b, 0x39, 0xbc, 0x32, 0x02,
	0xe3, 0xcd, 0x59, 0xdf, 0xa8, 0xbb, 0x46, 0x74, 0x5f, 0x40, 0x28, 0x15, 0x2d, 0x55, 0xd8, 0xe0,
	0x0f, 0x37, 0x34, 0x80, 0x87, 0x4d, 0x6f, 0xb8, 0xed, 0x0d, 0xbf, 0x6e, 0x7b, 0x0b, 0xfa, 0x27,
	0xa7, 0xbe, 0x73, 0x78, 0xe6, 0x83, 0xaf, 0xe7, 0xc7, 0x13, 0x30, 0xdb, 0xd4, 0x2f, 0x37, 0x69,
	0xf7, 0x13, 0x80, 0x83, 0x85, 0x88, 0xe6, 0x55, 0x11, 0x16, 0xbc, 0x4c, 0x45, 0x2c, 0x87, 0xbd,
	0xd1, 0xc6, 0x78, 0x6b, 0x8a, 0xfe, 0xd7, 0xcf, 0x2b, 0x5d, 0x16, 0x3c, 0x6f, 0x2c, 0xbf, 0x9d,
	0xf9, 0x4f, 0x92, 0x54, 0xed, 0x57, 0x0c, 0x47, 0x22, 0x23, 0x76, 0x30, 0xe6, 0x76, 0x5f, 0xc6,
	0x73, 0x72, 0x40, 0x68, 0xa5, 0xf6, 0x7f, 0x8f, 0x4a, 0xbd, 0x2f, 0xb8, 0xb4, 0x0e, 0xd2, 0xb0,
	0xf4, 0xcd, 0xe9, 0x56, 0x73, 0x3f, 0x03, 0x78, 0xbd, 0xfd, 0xc0, 0x2d, 0xd0, 0xd5, 0x4b, 0x05,
	0x1a, 0xd8, 0x9c, 0x15, 0x9f, 0x3e, 0xfe, 0x70, 0xe4, 0x3b, 0x5f, 0x8e, 0x7c, 0xe7, 0xe3, 0xf9,
	0xf1, 0xe4, 0x9e, 0x59, 0xe0, 0x83, 0xee, 0x0a, 0xff, 0x7b, 0x5d, 0x82, 0x97, 0x27, 0x2b, 0x04,
	0x96, 0x2b, 0x04, 0x7e, 0xae, 0x10, 0x38, 0x5c, 0x23, 0x67, 0xb9, 0x46, 0xce, 0x8f, 0x35, 0x72,
	0xde, 0x3e, 0xea, 0x10, 0x1a, 0x3f, 0xfb, 0x5b, 0x4c, 0x1f, 0x74, 0x9c, 0xb3, 0x34, 0x29, 0xa9,
	0x4a, 0x45, 0x2e, 0x0d, 0x25, 0xbb, 0xa6, 0x07, 0xfb, 0xf0, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x19, 0xca, 0xb8, 0xe4, 0x45, 0x03, 0x00, 0x00,
}

func (m *ClawbackVestingAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClawbackVestingAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClawbackVestingAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.VestingPeriods) > 0 {
		for iNdEx := len(m.VestingPeriods) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.VestingPeriods[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintVesting(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.LockupPeriods) > 0 {
		for iNdEx := len(m.LockupPeriods) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LockupPeriods[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintVesting(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.StartTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.StartTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintVesting(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x1a
	if len(m.FunderAddress) > 0 {
		i -= len(m.FunderAddress)
		copy(dAtA[i:], m.FunderAddress)
		i = encodeVarintVesting(dAtA, i, uint64(len(m.FunderAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.BaseVestingAccount != nil {
		{
			size, err := m.BaseVestingAccount.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintVesting(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintVesting(dAtA []byte, offset int, v uint64) int {
	offset -= sovVesting(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ClawbackVestingAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BaseVestingAccount != nil {
		l = m.BaseVestingAccount.Size()
		n += 1 + l + sovVesting(uint64(l))
	}
	l = len(m.FunderAddress)
	if l > 0 {
		n += 1 + l + sovVesting(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.StartTime)
	n += 1 + l + sovVesting(uint64(l))
	if len(m.LockupPeriods) > 0 {
		for _, e := range m.LockupPeriods {
			l = e.Size()
			n += 1 + l + sovVesting(uint64(l))
		}
	}
	if len(m.VestingPeriods) > 0 {
		for _, e := range m.VestingPeriods {
			l = e.Size()
			n += 1 + l + sovVesting(uint64(l))
		}
	}
	return n
}

func sovVesting(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVesting(x uint64) (n int) {
	return sovVesting(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClawbackVestingAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVesting
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClawbackVestingAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClawbackVestingAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseVestingAccount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVesting
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthVesting
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVesting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BaseVestingAccount == nil {
				m.BaseVestingAccount = &types.BaseVestingAccount{}
			}
			if err := m.BaseVestingAccount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FunderAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVesting
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthVesting
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVesting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FunderAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVesting
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthVesting
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVesting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.StartTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockupPeriods", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVesting
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthVesting
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVesting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LockupPeriods = append(m.LockupPeriods, types.Period{})
			if err := m.LockupPeriods[len(m.LockupPeriods)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VestingPeriods", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVesting
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthVesting
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVesting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VestingPeriods = append(m.VestingPeriods, types.Period{})
			if err := m.VestingPeriods[len(m.VestingPeriods)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVesting(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVesting
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipVesting(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVesting
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowVesting
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowVesting
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthVesting
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVesting
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVesting
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVesting        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVesting          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVesting = fmt.Errorf("proto: unexpected end of group")
)
