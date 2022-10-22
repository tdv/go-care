package care

type zeroMetaFilter struct {
}

func (*zeroMetaFilter) Allowed(string, []string) bool {
	return true
}

func NewZeroMetaFilter() MetaFilter {
	return &zeroMetaFilter{}
}
