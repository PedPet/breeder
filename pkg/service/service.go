package service

import (
	"context"

	"github.com/PedPet/breeder/model"
)

// Breeder service
type Breeder interface {
	CreateBreeder(ctx context.Context, affix, shortAffix, website string, owners []model.Owner) (*model.Breeder, error)
	GetBreeder(ctx context.Context, id int) (*model.Breeder, error)
	UpdateBreeder(ctx context.Context, id int, affix string, shortAffix string, owners []model.Owner) (*model.Breeder, error)
	DeleteBreeder(ctx context.Context, id int) error
}

type service struct {
    repository
}