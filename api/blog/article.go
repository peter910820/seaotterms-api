package blog

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/blog"
	model "seaotterms-api/model/blog"
	utils "seaotterms-api/utils/blog"
)

// query article data (all or use id to query single article data)
func QueryArticle(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Article
	var err error

	articleID, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if articleID != "" {
		err = db.Preload("Tags").First(&responseData, articleID).Error
	} else {
		err = db.Preload("Tags").Order("created_at desc").Find(&responseData).Error
	}
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Article資料查詢成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "Article資料查詢成功", &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

// create article data
func CreateArticle(c *fiber.Ctx, db *gorm.DB) error {
	var clientData dto.ArticleCreateRequest

	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if len(clientData.Tags) > 0 {
		var count int64
		db.Model(&model.Tag{}).Where("name IN ?", clientData.Tags).Count(&count)
		if count != int64(len(clientData.Tags)) {
			logrus.Error("缺少tags，請先建立tags")
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, "缺少tags，請先建立tags", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	data := model.Article{
		Title:   clientData.Title,
		Content: clientData.Content,
		Tags:    []model.Tag{},
	}
	for _, tag := range clientData.Tags {
		if !(strings.TrimSpace(tag) == "") {
			data.Tags = append(data.Tags, model.Tag{Name: tag})
		}
	}

	if err := db.Create(&data).Error; err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Article資料建立成功: " + clientData.Title)
	response := utils.ResponseFactory[any](c, fiber.StatusOK, "Article資料建立成功: "+clientData.Title, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// func UpdateArticle(c *fiber.Ctx, db *gorm.DB) error {
// 	// load client data
// 	var clientData
// 	if err := c.BodyParser(&clientData); err != nil {
// 		logrus.Error(err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"msg": err.Error(),
// 		})
// 	}
// 	response := dto.CreateDefalutCommonResponse[[]model.Article]()

// }

// Delete article data
func DeleteArticle(c *fiber.Ctx, db *gorm.DB) error {

	// URL decoding
	id, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var article model.Article
	db.Preload("Tags").First(&article, id)

	db.Model(&article).Association("Tags").Clear()

	db.Delete(&article)

	logrus.Info("刪除Article成功" + id)
	response := utils.ResponseFactory[any](c, fiber.StatusOK, "刪除Article成功: "+id, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// Query Article data use tag name
func QueryArticleForTag(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Article

	// URL decoding
	name, err := url.QueryUnescape(c.Params("name"))
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	err = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Joins("JOIN tags ON tags.name = article_tags.tag_name").
		Where("tags.name = ?", name).
		Find(&responseData).Error
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	logrus.Info("查詢指定Tag的Article成功: " + name)
	response := utils.ResponseFactory(c, fiber.StatusOK, "查詢指定Tag的Article成功"+name, &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
