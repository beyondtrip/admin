package admin

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor"
)

// QorAdminSetting admin settings
type QorAdminSetting struct {
	gorm.Model
	Key      string
	Resource string
	UserID   string
	Value    string `gorm:"size:65532"`
}

// LoadAdminSettings load admin settings
func LoadAdminSettings(key string, value interface{}, context *Context) error {
	var (
		settings     = []QorAdminSetting{}
		sqlCondition = "key = ? AND (resource = ? OR resource = ?) AND (user_id = ? OR user_id = ?)"
		resParams    = ""
		userID       = ""
	)

	if context.Resource != nil {
		resParams = context.Resource.ToParam()
	}

	if context.CurrentUser != nil {
		userID = ""
	}

	context.GetDB().Where(sqlCondition, key, resParams, "", userID, "").Order("user_id DESC, resource DESC, id DESC").Find(&settings)

	for _, setting := range settings {
		if err := json.Unmarshal([]byte(setting.Value), value); err != nil {
			return err
		}
	}

	return nil
}

// SaveAdminSettings save admin settings
func SaveAdminSettings(key string, value interface{}, res *Resource, user qor.CurrentUser, context *Context) error {
	var (
		tx          = context.GetDB()
		result, err = json.Marshal(value)
		resParams   = ""
		userID      = ""
	)

	if err != nil {
		return err
	}

	if res != nil {
		resParams = res.ToParam()
	}

	if user != nil {
		userID = ""
	}

	err = tx.Where(QorAdminSetting{
		Key:      key,
		UserID:   userID,
		Resource: resParams,
	}).Assign(QorAdminSetting{Value: string(result)}).FirstOrCreate(&QorAdminSetting{}).Error

	return err
}
