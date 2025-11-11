package blog

import (
	"time"

	model "seaotterms-api/model/blog"
)

type ArticleCreateRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type ArticleUpdateRequest struct {
	Title   string      `gorm:"NOT NULL" json:"title"`
	Content string      `gorm:"NOT NULL" json:"content"`
	Tags    []model.Tag `gorm:"many2many:article_tags" json:"tags"`
}

type TagCreateRequest struct {
	Name     string `json:"name"`
	IconName string `json:"iconName"`
}

type SystemTodoCreateRequest struct {
	SystemName  string     `json:"systemName"`
	Title       string     `json:"title"`
	Detail      string     `json:"detail"`
	Status      uint       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Urgency     uint       `json:"urgency"`
	CreatedName string     `json:"createdName"`
}

type SystemTodoUpdateRequest struct {
	SystemName  string     `json:"systemName"`
	Title       string     `json:"title"`
	Detail      string     `json:"detail"`
	Status      uint       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Urgency     uint       `json:"urgency"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	UpdatedName string     `json:"updatedName"`
}

type QuickSystemTodoUpdateRequest struct {
	Status      uint      `json:"status"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UpdatedName string    `json:"updatedName"`
}

type TodoTopicCreateRequest struct {
	TopicName  string `json:"topicName"`
	TopicOwner string `json:"topicOwner"`
	UpdateName string `json:"updateName"`
}

type TodoUpdateRequest struct {
	Status     uint      `json:"status"`
	UpdatedAt  time.Time `json:"updatedAt"`
	UpdateName string    `json:"updateName"`
}

type BrandCreateRequest struct {
	Name        string `json:"name"`
	WorkAmount  int    `json:"workAmount"`  // 作品數量
	OfficialUrl string `json:"officialUrl"` // 官網URL
	Dissolution bool   `json:"dissolution"` // 解散標記
}

type GameCreateRequest struct {
	Name            string    `json:"name"`
	ChineseName     string    `json:"chineseName"`
	BrandID         int       `json:"brandId"`
	AllAges         bool      `json:"allAges"`
	ReleaseDate     time.Time `json:"releaseDate"`
	OpUrl           string    `json:"opUrl"`
	GameDescription string    `json:"gameDescription"`
}

type PlayRecordCreateRequest struct {
	GameID               int       `json:"gameId"`
	EndPlayDate          time.Time `json:"endPlayDate"`
	OpDisplayScore       *float64  `json:"opDisplayScore"`
	OpSongScore          *float64  `json:"opSongScore"`
	OpCompatibilityScore *float64  `json:"opCompatibilityScore"`
	EdDisplayScore       *float64  `json:"edDisplayScore"`
	EdSongScore          *float64  `json:"edSongScore"`
	MusicScore           *float64  `json:"musicScore"`
	PlotScore            float64   `json:"plotScore"`
	ArtScore             float64   `json:"artScore"`
	SystemScore          float64   `json:"systemScore"`
	ThemeScore           float64   `json:"themeScore"`
	ConclusionScore      float64   `json:"conclusionScore"`
	Category             string    `json:"category"`
	Recommended          int       `json:"recommended"`
	Experience           string    `json:"experience"`
}

type RegisterRequest struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	CheckPassword string `json:"checkPassword"`
}

type UserUpdateRequest struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Exp        int       `json:"exp"`
	Management bool      `json:"management"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdateName string    `json:"update_name"`
	Avatar     string    `json:"avatar"`
}

type SelfBrandUpdateRequest struct {
	Brand       string `json:"brand"`
	Username    string `json:"username"`
	Completed   int    `json:"completed"`
	Total       int    `json:"total"`
	Dissolution bool   `json:"dissolution"`
}

type SelfGameUpdateRequest struct {
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	ReleaseDate time.Time `json:"releaseDate"`
	AllAges     bool      `json:"allAges"`
	EndDate     time.Time `json:"endDate"`
	Username    string    `json:"username"`
}

// login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
