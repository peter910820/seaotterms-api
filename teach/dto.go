package teach

import "github.com/lib/pq"

type SeriesCreateRequest struct {
	Title string `gorm:"NOT NULL" json:"title"`
	Image string `json:"image"`
}

type SeriesModifyRequest struct {
	Title string `gorm:"NOT NULL" json:"title"`
	Image string `json:"image"`
}

type ArtilceCreateRequest struct {
	Title    string         `gorm:"NOT NULL" json:"title"`
	SeriesID uint           `gorm:"NOT NULL" json:"seriesId"`
	Image    string         `json:"image"`
	Tags     pq.StringArray `gorm:"type:text[]" json:"tags"`
	Content  string         `gorm:"NOT NULL" json:"content"`
}

type ArtilceModifyRequest struct {
	Title    string         `gorm:"NOT NULL" json:"title"`
	SeriesID uint           `gorm:"NOT NULL" json:"seriesId"`
	Image    string         `json:"image"`
	Tags     pq.StringArray `gorm:"type:text[]" json:"tags"`
	Content  string         `gorm:"NOT NULL" json:"content"`
}
