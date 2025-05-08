//go:build template_disable
// +build template_disable

// nolint:glint,style check,unused
package usecase

import (
	"module-clean/internal/modules/template/domain/entities"
	"module-clean/internal/modules/template/domain/repository"
)

type TemplateUseCase struct {
	TemplateRepo repository.TemplateRepository
}

func NewTemplateUseCase(templateRepo repository.TemplateRepository) *TemplateUseCase {
	return &TemplateUseCase{
		TemplateRepo: templateRepo,
	}
}
func (t *TemplateUseCase) CreateTemplate(ctx context.Context, template *entities.Template) error {
	err := t.TemplateRepo.Create(ctx, template)
	if err != nil {
		return MapInfraErrorToUseCaseError(err)
	}
	return nil
}
func (t *TemplateUseCase) GetTemplateByID(ctx context.Context, id int) (*entities.Template, error) {
	template, err := t.TemplateRepo.Read(ctx, id)
	if err != nil {
		return nil, MapInfraErrorToUseCaseError(err)
	}
	return template, nil
}
