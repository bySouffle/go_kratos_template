// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.6.2
// - protoc             v4.25.0
// source: template/v1/template.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationTemplateCreateTemplate = "/api.template.v1.Template/CreateTemplate"
const OperationTemplateGetTemplate = "/api.template.v1.Template/GetTemplate"
const OperationTemplateWSTemplate = "/api.template.v1.Template/WSTemplate"

type TemplateHTTPServer interface {
	CreateTemplate(context.Context, *CreateTemplateRequest) (*CreateTemplateReply, error)
	GetTemplate(context.Context, *GetTemplateRequest) (*GetTemplateReply, error)
	WSTemplate(context.Context, *WSTemplateRequest) (*WSTemplateReply, error)
}

func RegisterTemplateHTTPServer(s *http.Server, srv TemplateHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/template/create", _Template_CreateTemplate0_HTTP_Handler(srv))
	r.GET("/v1/template/create/{name}", _Template_CreateTemplate1_HTTP_Handler(srv))
	r.POST("/v1/template/get", _Template_GetTemplate0_HTTP_Handler(srv))
	r.GET("/v1/template/get/{name}", _Template_GetTemplate1_HTTP_Handler(srv))
	r.GET("/v1/ws/{name}", _Template_WSTemplate0_HTTP_Handler(srv))
}

func _Template_CreateTemplate0_HTTP_Handler(srv TemplateHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateTemplateRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTemplateCreateTemplate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateTemplate(ctx, req.(*CreateTemplateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateTemplateReply)
		return ctx.Result(200, reply)
	}
}

func _Template_CreateTemplate1_HTTP_Handler(srv TemplateHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateTemplateRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTemplateCreateTemplate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateTemplate(ctx, req.(*CreateTemplateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateTemplateReply)
		return ctx.Result(200, reply)
	}
}

func _Template_GetTemplate0_HTTP_Handler(srv TemplateHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetTemplateRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTemplateGetTemplate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetTemplate(ctx, req.(*GetTemplateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetTemplateReply)
		return ctx.Result(200, reply)
	}
}

func _Template_GetTemplate1_HTTP_Handler(srv TemplateHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetTemplateRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTemplateGetTemplate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetTemplate(ctx, req.(*GetTemplateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetTemplateReply)
		return ctx.Result(200, reply)
	}
}

func _Template_WSTemplate0_HTTP_Handler(srv TemplateHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in WSTemplateRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTemplateWSTemplate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.WSTemplate(ctx, req.(*WSTemplateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*WSTemplateReply)
		return ctx.Result(200, reply)
	}
}

type TemplateHTTPClient interface {
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest, opts ...http.CallOption) (rsp *CreateTemplateReply, err error)
	GetTemplate(ctx context.Context, req *GetTemplateRequest, opts ...http.CallOption) (rsp *GetTemplateReply, err error)
	WSTemplate(ctx context.Context, req *WSTemplateRequest, opts ...http.CallOption) (rsp *WSTemplateReply, err error)
}

type TemplateHTTPClientImpl struct {
	cc *http.Client
}

func NewTemplateHTTPClient(client *http.Client) TemplateHTTPClient {
	return &TemplateHTTPClientImpl{client}
}

func (c *TemplateHTTPClientImpl) CreateTemplate(ctx context.Context, in *CreateTemplateRequest, opts ...http.CallOption) (*CreateTemplateReply, error) {
	var out CreateTemplateReply
	pattern := "/v1/template/create/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationTemplateCreateTemplate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TemplateHTTPClientImpl) GetTemplate(ctx context.Context, in *GetTemplateRequest, opts ...http.CallOption) (*GetTemplateReply, error) {
	var out GetTemplateReply
	pattern := "/v1/template/get/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationTemplateGetTemplate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TemplateHTTPClientImpl) WSTemplate(ctx context.Context, in *WSTemplateRequest, opts ...http.CallOption) (*WSTemplateReply, error) {
	var out WSTemplateReply
	pattern := "/v1/ws/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationTemplateWSTemplate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
