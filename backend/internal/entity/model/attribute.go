package model

type Attributes []*Attribute

func (a Attributes) GetAttribute(key string) string {
	for _, attribute := range a {
		if attribute.Key == key {
			return attribute.Value
		}
	}

	return ""
}

type Attribute struct {
	ID    int    `json:"-" gorm:"primarykey"`
	Key   string `json:"key" gorm:"column:key_name"`
	Value string `json:"value"`
}

func (Attribute) TableName() string {
	return "authz_attributes"
}
