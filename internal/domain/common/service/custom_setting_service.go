package service

// CustomSettingService maps to CustomSettingService.java
// It's responsible for managing and parsing custom settings properties.
// @MappedFrom CustomSettingService
type CustomSettingService struct {
	*CustomValueService // Inherits from CustomValueService

	// TODO: implement maps for parsed and cached properties
}

func NewCustomSettingService() *CustomSettingService {
	return &CustomSettingService{
		CustomValueService: NewCustomValueService(
			"The value of the setting \"",
			"The string value of the setting \"",
			"The string value length of the setting \"",
			"The array value of the setting \""),
	}
}

// TODO: implement `updateGlobalProperties(CustomSettingsProperties properties)` and `parseSettings(boolean ignoreUnknownSettings, map<string, value> inputSettings)`
