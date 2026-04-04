package throttle

// AdminApiRateLimitingManager maps to AdminApiRateLimitingManager.java
// It inherits BaseAdminApiRateLimitingManager in Java to restrict admin API rate.
// @MappedFrom AdminApiRateLimitingManager
type AdminApiRateLimitingManager struct {
	// Inherit fields/methods from BaseAdminApiRateLimitingManager here when available.
}

// @MappedFrom AdminApiRateLimitingManager(TaskManager taskManager, TurmsPropertiesManager propertiesManager)
func NewAdminApiRateLimitingManager() *AdminApiRateLimitingManager {
	return &AdminApiRateLimitingManager{}
}
