package repositories

import (
	"context"

	cat_entity "github.com/danzbraham/cats-social/internal/domains/entities/cats"
)

type CatRepository interface {
	CreateCat(ctx context.Context, cat *cat_entity.Cat) error
}
