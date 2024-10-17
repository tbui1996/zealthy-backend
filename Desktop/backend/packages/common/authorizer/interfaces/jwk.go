package interfaces

import (
	"context"

	"github.com/lestrrat-go/jwx/jwk"
)

type Jwk interface {
	Fetch(ctx context.Context, urlstring string, options ...jwk.FetchOption) (jwk.Set, error)
}
