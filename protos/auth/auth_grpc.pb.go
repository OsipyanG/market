// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.19.6
// source: auth/auth.proto

package auth

import (
	context "context"
	jwt "github.com/OsipyanG/market/protos/jwt"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Auth_NewUser_FullMethodName        = "/auth.Auth/NewUser"
	Auth_Login_FullMethodName          = "/auth.Auth/Login"
	Auth_UpdateTokens_FullMethodName   = "/auth.Auth/UpdateTokens"
	Auth_UpdatePassword_FullMethodName = "/auth.Auth/UpdatePassword"
	Auth_Logout_FullMethodName         = "/auth.Auth/Logout"
	Auth_GetJWTClaims_FullMethodName   = "/auth.Auth/GetJWTClaims"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	NewUser(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Tokens, error)
	Login(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Tokens, error)
	UpdateTokens(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*Tokens, error)
	UpdatePassword(ctx context.Context, in *RequestUpdatePassword, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Logout(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetJWTClaims(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*jwt.JWTClaims, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) NewUser(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Tokens, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tokens)
	err := c.cc.Invoke(ctx, Auth_NewUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Login(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Tokens, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tokens)
	err := c.cc.Invoke(ctx, Auth_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UpdateTokens(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*Tokens, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tokens)
	err := c.cc.Invoke(ctx, Auth_UpdateTokens_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UpdatePassword(ctx context.Context, in *RequestUpdatePassword, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Auth_UpdatePassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Logout(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Auth_Logout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) GetJWTClaims(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*jwt.JWTClaims, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(jwt.JWTClaims)
	err := c.cc.Invoke(ctx, Auth_GetJWTClaims_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	NewUser(context.Context, *UserCredentials) (*Tokens, error)
	Login(context.Context, *UserCredentials) (*Tokens, error)
	UpdateTokens(context.Context, *RefreshToken) (*Tokens, error)
	UpdatePassword(context.Context, *RequestUpdatePassword) (*emptypb.Empty, error)
	Logout(context.Context, *RefreshToken) (*emptypb.Empty, error)
	GetJWTClaims(context.Context, *AccessToken) (*jwt.JWTClaims, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) NewUser(context.Context, *UserCredentials) (*Tokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewUser not implemented")
}
func (UnimplementedAuthServer) Login(context.Context, *UserCredentials) (*Tokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServer) UpdateTokens(context.Context, *RefreshToken) (*Tokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTokens not implemented")
}
func (UnimplementedAuthServer) UpdatePassword(context.Context, *RequestUpdatePassword) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePassword not implemented")
}
func (UnimplementedAuthServer) Logout(context.Context, *RefreshToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAuthServer) GetJWTClaims(context.Context, *AccessToken) (*jwt.JWTClaims, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJWTClaims not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_NewUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCredentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).NewUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_NewUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).NewUser(ctx, req.(*UserCredentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCredentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Login(ctx, req.(*UserCredentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UpdateTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UpdateTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_UpdateTokens_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UpdateTokens(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUpdatePassword)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_UpdatePassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UpdatePassword(ctx, req.(*RequestUpdatePassword))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Logout(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_GetJWTClaims_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).GetJWTClaims(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_GetJWTClaims_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).GetJWTClaims(ctx, req.(*AccessToken))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewUser",
			Handler:    _Auth_NewUser_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Auth_Login_Handler,
		},
		{
			MethodName: "UpdateTokens",
			Handler:    _Auth_UpdateTokens_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _Auth_UpdatePassword_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Auth_Logout_Handler,
		},
		{
			MethodName: "GetJWTClaims",
			Handler:    _Auth_GetJWTClaims_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth/auth.proto",
}

const (
	AuthAdmin_DeleteUser_FullMethodName           = "/auth.AuthAdmin/DeleteUser"
	AuthAdmin_SetAccessLevel_FullMethodName       = "/auth.AuthAdmin/SetAccessLevel"
	AuthAdmin_GetAllUsersWithLevel_FullMethodName = "/auth.AuthAdmin/GetAllUsersWithLevel"
)

// AuthAdminClient is the client API for AuthAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthAdminClient interface {
	DeleteUser(ctx context.Context, in *RequestByUserID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetAccessLevel(ctx context.Context, in *SetAccessLevelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetAllUsersWithLevel(ctx context.Context, in *RequestByLevel, opts ...grpc.CallOption) (*UsersInfoResponse, error)
}

type authAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthAdminClient(cc grpc.ClientConnInterface) AuthAdminClient {
	return &authAdminClient{cc}
}

func (c *authAdminClient) DeleteUser(ctx context.Context, in *RequestByUserID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AuthAdmin_DeleteUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authAdminClient) SetAccessLevel(ctx context.Context, in *SetAccessLevelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AuthAdmin_SetAccessLevel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authAdminClient) GetAllUsersWithLevel(ctx context.Context, in *RequestByLevel, opts ...grpc.CallOption) (*UsersInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UsersInfoResponse)
	err := c.cc.Invoke(ctx, AuthAdmin_GetAllUsersWithLevel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthAdminServer is the server API for AuthAdmin service.
// All implementations must embed UnimplementedAuthAdminServer
// for forward compatibility
type AuthAdminServer interface {
	DeleteUser(context.Context, *RequestByUserID) (*emptypb.Empty, error)
	SetAccessLevel(context.Context, *SetAccessLevelRequest) (*emptypb.Empty, error)
	GetAllUsersWithLevel(context.Context, *RequestByLevel) (*UsersInfoResponse, error)
	mustEmbedUnimplementedAuthAdminServer()
}

// UnimplementedAuthAdminServer must be embedded to have forward compatible implementations.
type UnimplementedAuthAdminServer struct {
}

func (UnimplementedAuthAdminServer) DeleteUser(context.Context, *RequestByUserID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedAuthAdminServer) SetAccessLevel(context.Context, *SetAccessLevelRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAccessLevel not implemented")
}
func (UnimplementedAuthAdminServer) GetAllUsersWithLevel(context.Context, *RequestByLevel) (*UsersInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUsersWithLevel not implemented")
}
func (UnimplementedAuthAdminServer) mustEmbedUnimplementedAuthAdminServer() {}

// UnsafeAuthAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthAdminServer will
// result in compilation errors.
type UnsafeAuthAdminServer interface {
	mustEmbedUnimplementedAuthAdminServer()
}

func RegisterAuthAdminServer(s grpc.ServiceRegistrar, srv AuthAdminServer) {
	s.RegisterService(&AuthAdmin_ServiceDesc, srv)
}

func _AuthAdmin_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestByUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthAdminServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthAdmin_DeleteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthAdminServer).DeleteUser(ctx, req.(*RequestByUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthAdmin_SetAccessLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAccessLevelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthAdminServer).SetAccessLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthAdmin_SetAccessLevel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthAdminServer).SetAccessLevel(ctx, req.(*SetAccessLevelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthAdmin_GetAllUsersWithLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestByLevel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthAdminServer).GetAllUsersWithLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthAdmin_GetAllUsersWithLevel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthAdminServer).GetAllUsersWithLevel(ctx, req.(*RequestByLevel))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthAdmin_ServiceDesc is the grpc.ServiceDesc for AuthAdmin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthAdmin",
	HandlerType: (*AuthAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteUser",
			Handler:    _AuthAdmin_DeleteUser_Handler,
		},
		{
			MethodName: "SetAccessLevel",
			Handler:    _AuthAdmin_SetAccessLevel_Handler,
		},
		{
			MethodName: "GetAllUsersWithLevel",
			Handler:    _AuthAdmin_GetAllUsersWithLevel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth/auth.proto",
}