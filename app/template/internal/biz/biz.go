package biz

import (
	"github.com/google/wire"
	"go_kratos_template/pkg/ws"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewTemplateUseCase, ws.NewNewClientManagerWithRun)
