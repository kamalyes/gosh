/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-11-16 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 12:16:02
 * @FilePath: \gosh\response.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// ResponseOption 是用于构建返回响应的结构体
type ResponseOption struct {
	Data      interface{}
	SceneCode SceneCode
	HttpCode  StatusCode
	Message   string
}

// convertToSceneCode 辅助函数用于将输入值转换为 SceneCode 类型
func convertToSceneCode(val interface{}) SceneCode {
	code, ok := val.(SceneCode)
	if !ok {
		// 如果类型断言失败，可以执行适当的错误处理逻辑
		// 这里简单地返回一个默认值
		return SceneCode(Success)
	}
	return SceneCode(code)
}

// convertToHttpStatusCode 辅助函数用于将输入值转换为 StatusCode 类型
func convertToHttpStatusCode(val interface{}) StatusCode {
	statusCode, ok := val.(StatusCode)
	if !ok {
		return StatusCode(StatusOK)
	}
	return StatusCode(statusCode)
}

// NewResponseOption 用于创建 ResponseOption 实例
func NewResponseOption(data interface{}, options ...interface{}) *ResponseOption {
	response := &ResponseOption{
		Data: data,
	}

	for _, option := range options {
		switch opt := option.(type) {
		case int:
			response.HttpCode = StatusCode(opt)
			response.SceneCode = SceneCode(opt)
		case SceneCode:
			response.SceneCode = opt
		case StatusCode:
			response.HttpCode = opt
		case string:
			response.Message = opt
		}
	}

	return response
}

// Merge 用于处理 ResponseOption 实例的属性值
func (o *ResponseOption) Merge() *ResponseOption {
	// 将 o.Code 的值根据条件进行转换
	o.SceneCode = convertToSceneCode(ternary(o.SceneCode == 0, Success, o.SceneCode))

	// 根据条件将 o.StatusCode 的值进行转换
	o.HttpCode = convertToHttpStatusCode(ternary(o.HttpCode == 0, StatusOK, o.HttpCode))

	// 根据条件设置消息内容
	if o.Message == "" {
		o.Message = GetSceneCodeText(o.SceneCode)
	}
	if o.Message == "" {
		o.Message = GetStatusCodeText(o.HttpCode)
	}

	return o
}

// ternary 函数实现三元运算
func ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// SendJSONResponse 生成 JSON 格式的响应
func SendJSONResponse(c *Context, respOption *ResponseOption) error {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Merge()

	// 创建一个map来存储不包含HttpStatusCode和Language的字段
	cleanedResp := map[string]interface{}{
		"data":    respOption.Data,
		"code":    respOption.SceneCode,
		"message": respOption.Message,
	}

	c.WriteJSONResponse(int(respOption.HttpCode), cleanedResp)
	return nil
}

// Gen400xResponse 生成 HTTP 400x 错误响应
func Gen400xResponse(ctx *Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.SceneCode = BadRequest
	respOption.HttpCode = StatusBadRequest
	SendJSONResponse(ctx, respOption)
}

// Gen500xResponse 生成 HTTP 500 错误响应
func Gen500xResponse(ctx *Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.SceneCode = Fail
	respOption.HttpCode = StatusInternalServerError
	SendJSONResponse(ctx, respOption)
}

// ValidatorError 处理字段校验异常
func ValidatorError(ctx *Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		Gen400xResponse(ctx, &ResponseOption{
			SceneCode: ValidateError,
			Data:      RemoveTopStruct(errs.Translate(ctx.Engine.Config.Trans)),
		})
		return
	}

	Gen400xResponse(ctx, &ResponseOption{Message: err.Error()})
}

// RemoveTopStruct 定义一个去掉结构体名称前缀的自定义方法：
func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
