package care

type Options struct {
	Switch     Switch
	Methods    Methods
	MetaFilter MetaFilter
	Hash       Hash
	Cache      Cache
}

func NewDefaultOptions() (*Options, error) {
	opts := Options{
		Switch:     newSwitch(),
		Methods:    newMethodsStorage(),
		MetaFilter: NewZeroMetaFilter(),
		Hash:       newDefaultHash(),
		Cache:      NewInMemoryLRUCache(1024),
	}

	return &opts, nil
}
