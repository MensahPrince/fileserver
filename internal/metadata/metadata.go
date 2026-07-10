package metadata

import (
	"context"

	"github.com/MensahPrince/fileserver/types"
)


type Metadata interface {
	Save(ctx context.Context, meta types.FileMeta) error
	Get(ctx context.Context, id string) (types.FileMeta, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]types.FileMeta, error)

}