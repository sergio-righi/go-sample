package models

type Document struct {
	Key       string `json:"key" bson:"key"`
	Size      *int64 `json:"size,omitempty" bson:"size,omitempty"`
	Count     *int64 `json:"count,omitempty" bson:"count,omitempty"`
	Delimiter *bool  `json:"delimiter,omitempty" bson:"delimiter,omitempty"`
}
