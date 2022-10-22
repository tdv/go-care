package care

type Headers struct {
	Allowed    []string
	Disallowed []string
}
type MetaFilter interface {
	Allowed(key string, val []string) bool
}

type metaFilter struct {
	allowedHeaders    map[string]struct{}
	disallowedHeaders map[string]struct{}
}

func (this *metaFilter) Allowed(key string, val []string) bool {
	if _, ok := this.disallowedHeaders[key]; ok {
		return false
	}

	_, ok := this.allowedHeaders[key]
	return ok
}

func NewMetaFilter(headers Headers) MetaFilter {
	inst := metaFilter{
		allowedHeaders:    make(map[string]struct{}),
		disallowedHeaders: make(map[string]struct{}),
	}

	if len(headers.Allowed) > 0 {
		for _, h := range headers.Allowed {
			inst.allowedHeaders[h] = struct{}{}
		}
	}

	if len(headers.Disallowed) > 0 {
		for _, h := range headers.Disallowed {
			inst.disallowedHeaders[h] = struct{}{}
		}
	}

	return &inst
}
