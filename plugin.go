// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"strings"

	gdttypes "github.com/jaypipes/gdt-core/types"
	gdthttp "github.com/jaypipes/gdt-http"
)

var (
	knownPlugins = []gdttypes.Plugin{
		gdthttp.Plugin(),
	}
)

// RegisterPlugin registers a plugin with gdt's set of known plugins.
//
// Generally only plugin authors will ever need to call this function. It is
// not required for normal use of gdt or any known plugin.
func RegisterPlugin(plugin gdttypes.Plugin) {
	for _, p := range knownPlugins {
		if strings.EqualFold(p.Info().Name, plugin.Info().Name) {
			return
		}
	}
	knownPlugins = append(knownPlugins, plugin)
}
