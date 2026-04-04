package elasticsearch

import (
	"im.turms/server/internal/storage/elasticsearch/model"
)

/**
 * @author James Chen
 */
type IndexTextFieldSetting struct {
	FieldToProperty map[string]model.Property
	Analysis        *model.IndexSettingsAnalysis
}
