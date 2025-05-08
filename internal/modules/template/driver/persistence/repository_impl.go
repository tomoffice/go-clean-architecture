//go:build template_disable
// +build template_disable

// nolint:glint,style check,unused
package persistence

import (
	"context"
	"module-clean/internal/modules/template/domain/entities"
)

type TemplateRepo struct{}

func (t2 TemplateRepo) Create(ctx context.Context, t *entities.Template) error {
	//TODO implement me
	panic("implement me")
}

func (t2 TemplateRepo) Read(ctx context.Context, id int) (*entities.Template, error) {
	//TODO implement me
	panic("implement me")
}

func (t2 TemplateRepo) Update(ctx context.Context, t *entities.Template) error {
	//TODO implement me
	panic("implement me")
}

func (t2 TemplateRepo) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewTemplateRepo() *TemplateRepo {
	return &TemplateRepo{}
}
