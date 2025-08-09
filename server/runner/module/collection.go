package module

import (
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"
	"github.com/SSripilaipong/muto/builtin"
	"github.com/SSripilaipong/muto/builtin/global"
	"github.com/SSripilaipong/muto/builtin/portal"
	"github.com/SSripilaipong/muto/core/module"
	"github.com/SSripilaipong/muto/syntaxtree"
)

type Collection struct {
	modules map[string]module.Module
}

func NewCollection() *Collection {
	globalMod := global.NewModule()
	ports := portal.NewDefaultPortal()
	imported := builtin.NewBuiltinImportMapping(nil).Attach(globalMod, ports)

	mod := module.BuildUserDefinedModule(syntaxtree.NewModule(nil)).
		Attach(module.NewDependency(globalMod, ports, imported))

	return &Collection{
		modules: map[string]module.Module{
			VersionDefault: mod,
		},
	}
}

func (c *Collection) Get(version string) rslt.Of[module.Module] {
	mod, ok := c.modules[version]
	if !ok {
		return rslt.Error[module.Module](fmt.Errorf("unsupported version: %s", version))
	}
	return rslt.Value(mod)
}
