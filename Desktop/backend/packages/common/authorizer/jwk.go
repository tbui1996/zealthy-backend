package authorizer

import (
	"context"

	"github.com/lestrrat-go/jwx/jwk"
)

type Jwk struct{}

func (j *Jwk) Fetch(ctx context.Context, urlstring string, options ...jwk.FetchOption) (jwk.Set, error) {
	return jwk.Fetch(ctx, urlstring, options...)
}
