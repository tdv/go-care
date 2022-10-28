package care

type Options struct {
	Switch     Switch
	Methods    Methods
	MetaFilter MetaFilter
	Hash       Hash
	Cache      Cache
}

func NewOptions() *Options {
	opts := Options{
		Switch:     newSwitch(),
		Methods:    newMethodsStorage(),
		MetaFilter: NewZeroMetaFilter(),
		Hash:       newDefaultHash(),
		Cache:      NewInMemoryCache(1024),
	}

	return &opts
}
