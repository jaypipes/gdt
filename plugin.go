// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	gdtplugin "github.com/jaypipes/gdt-core/plugin"
	gdttypes "github.com/jaypipes/gdt-core/types"
)

var (
	knownPlugins = gdtplugin.NewRegistry()
)

// RegisterPlugin registers a plugin with gdt's set of known plugins.
//
// Generally only plugin authors will ever need to call this function. It is
// not required for normal use of gdt or any known plugin.
func RegisterPlugin(p gdttypes.Plugin) {
	knownPlugins.Add(p)
}
