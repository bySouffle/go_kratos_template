package timer

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTimer)
