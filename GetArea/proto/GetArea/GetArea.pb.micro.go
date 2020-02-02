// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/GetArea/GetArea.proto

package go_micro_srv_GetArea

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

// Client API for GetArea service

type GetAreaService interface {
	GetArea(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type getAreaService struct {
	c    client.Client
	name string
}

func NewGetAreaService(name string, c client.Client) GetAreaService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.GetArea"
	}
	return &getAreaService{
		c:    c,
		name: name,
	}
}

func (c *getAreaService) GetArea(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "GetArea.GetArea", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GetArea service

type GetAreaHandler interface {
	GetArea(context.Context, *Request, *Response) error
}

func RegisterGetAreaHandler(s server.Server, hdlr GetAreaHandler, opts ...server.HandlerOption) error {
	type getArea interface {
		GetArea(ctx context.Context, in *Request, out *Response) error
	}
	type GetArea struct {
		getArea
	}
	h := &getAreaHandler{hdlr}
	return s.Handle(s.NewHandler(&GetArea{h}, opts...))
}

type getAreaHandler struct {
	GetAreaHandler
}

func (h *getAreaHandler) GetArea(ctx context.Context, in *Request, out *Response) error {
	return h.GetAreaHandler.GetArea(ctx, in, out)
}