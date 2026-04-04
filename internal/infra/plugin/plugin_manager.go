package plugin

import "context"

// PluginManager is a stub for managing plugins.
type PluginManager struct{}

func NewPluginManager() *PluginManager {
	return &PluginManager{}
}

func (m *PluginManager) HasRunningExtensions(extensionPointClass string) bool {
	// TODO: implement
	return false
}

func (m *PluginManager) InvokeExtensionPoints(ctx context.Context, extensionPointClass string, method string, args ...interface{}) (bool, error) {
	// TODO: implement
	return false, nil
}
