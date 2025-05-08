//go:build template_disable
// +build template_disable

// nolint:glint,style check,unused
package controller

import "module-clean/internal/modules/member/usecase"

type TemplateController struct {
	memberUseCase *usecase.MemberUseCase
}

func NewTemplateController(memberUseCase *usecase.MemberUseCase) *TemplateController {
	return &TemplateController{
		memberUseCase: memberUseCase,
	}
}
