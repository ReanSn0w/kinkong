package resolver

import "github.com/ReanSn0w/kincong/internal/utils"

const (
	ResolverTypeIP ResolverType = iota
	ResolverTypeDNS
	ResolverTypeASN
)

type ResolverType int

type ResolverItem interface {
	Type() ResolverType
	Resolve(string) ([]utils.ResolvedSubnet, error)
}
