package asn

import (
	"errors"

	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/utils"
)

type Resolver struct {
	key string
}

func New(key string) *Resolver {
	return &Resolver{
		key: key,
	}
}

func (r *Resolver) Type() resolver.ResolverType {
	return resolver.ResolverTypeASN
}

func (r *Resolver) Resolve(ip string) ([]utils.ResolvedSubnet, error) {
	return nil, errors.New("not implemented")
}
