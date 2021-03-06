// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: userBaseOperation.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for UserBaseOps service

type UserBaseOpsService interface {
	RegisterUser(ctx context.Context, in *RegRequest, opts ...client.CallOption) (*RegResponse, error)
	UnRegisterUser(ctx context.Context, in *UnRegRequest, opts ...client.CallOption) (*UnRegResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...client.CallOption) (*LogoutResponse, error)
}

type userBaseOpsService struct {
	c    client.Client
	name string
}

func NewUserBaseOpsService(name string, c client.Client) UserBaseOpsService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &userBaseOpsService{
		c:    c,
		name: name,
	}
}

func (c *userBaseOpsService) RegisterUser(ctx context.Context, in *RegRequest, opts ...client.CallOption) (*RegResponse, error) {
	req := c.c.NewRequest(c.name, "UserBaseOps.RegisterUser", in)
	out := new(RegResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userBaseOpsService) UnRegisterUser(ctx context.Context, in *UnRegRequest, opts ...client.CallOption) (*UnRegResponse, error) {
	req := c.c.NewRequest(c.name, "UserBaseOps.UnRegisterUser", in)
	out := new(UnRegResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userBaseOpsService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "UserBaseOps.Login", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userBaseOpsService) Logout(ctx context.Context, in *LogoutRequest, opts ...client.CallOption) (*LogoutResponse, error) {
	req := c.c.NewRequest(c.name, "UserBaseOps.Logout", in)
	out := new(LogoutResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserBaseOps service

type UserBaseOpsHandler interface {
	RegisterUser(context.Context, *RegRequest, *RegResponse) error
	UnRegisterUser(context.Context, *UnRegRequest, *UnRegResponse) error
	Login(context.Context, *LoginRequest, *LoginResponse) error
	Logout(context.Context, *LogoutRequest, *LogoutResponse) error
}

func RegisterUserBaseOpsHandler(s server.Server, hdlr UserBaseOpsHandler, opts ...server.HandlerOption) error {
	type userBaseOps interface {
		RegisterUser(ctx context.Context, in *RegRequest, out *RegResponse) error
		UnRegisterUser(ctx context.Context, in *UnRegRequest, out *UnRegResponse) error
		Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error
		Logout(ctx context.Context, in *LogoutRequest, out *LogoutResponse) error
	}
	type UserBaseOps struct {
		userBaseOps
	}
	h := &userBaseOpsHandler{hdlr}
	return s.Handle(s.NewHandler(&UserBaseOps{h}, opts...))
}

type userBaseOpsHandler struct {
	UserBaseOpsHandler
}

func (h *userBaseOpsHandler) RegisterUser(ctx context.Context, in *RegRequest, out *RegResponse) error {
	return h.UserBaseOpsHandler.RegisterUser(ctx, in, out)
}

func (h *userBaseOpsHandler) UnRegisterUser(ctx context.Context, in *UnRegRequest, out *UnRegResponse) error {
	return h.UserBaseOpsHandler.UnRegisterUser(ctx, in, out)
}

func (h *userBaseOpsHandler) Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.UserBaseOpsHandler.Login(ctx, in, out)
}

func (h *userBaseOpsHandler) Logout(ctx context.Context, in *LogoutRequest, out *LogoutResponse) error {
	return h.UserBaseOpsHandler.Logout(ctx, in, out)
}
