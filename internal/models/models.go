package models

// Region - регион
type Region struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UsageType - тип использования травы
type UsageType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Usage - описание того как можно использовать траву
type Usage struct {
	ID          int    `json:"id"`
	HerbID      int    `json:"herb_id"`
	UsageTypeID int    `json:"usage_type_id"`
	Description string `json:"description"`

	HerbName      string `json:"herb_name,omitempty"`
	UsageTypeName string `json:"usage_type_name,omitempty"`
}

// HerbRegion - связь многое ко многому
type HerbRegion struct {
	HerbID   int `json:"herb_id"`
	RegionID int `json:"region_id"`

	HerbName   string `json:"herb_name,omitempty"`
	RegionName string `json:"region_name,omitempty"`
}

// Структура для ответа
type HerbWithDetails struct {
	Herb    Herb     `json:"herb"`
	Regions []Region `json:"regions"`
	Usages  []Usage  `json:"usages"`
}
