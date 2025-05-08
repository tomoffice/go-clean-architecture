package repository

import (
	"context"
	"module-clean/internal/modules/template/domain/entities"
)

type TemplateRepository interface {
	Create(ctx context.Context, t *entities.Template) error
	Read(ctx context.Context, id int) (*entities.Template, error)
	Update(ctx context.Context, t *entities.Template) error
	Delete(ctx context.Context, id int) error
}
