package blog

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	model "seaotterms-db/blog"
)

// Request
type (
	ArticleCreateRequest struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	ArticleUpdateRequest struct {
		Title   *string      `json:"title"`
		Content *string      `json:"content"`
		Tags    *[]model.Tag `json:"tags"`
	}
)

// query article data (all or use id to query single article data)
func QueryArticle(c fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Article
	var err error

	articleID, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if articleID == "" {
		responseData, err = model.FindAllArticles(db)
	} else {
		var article *model.Article
		article, err = model.FindArticleByID(db, articleID)
		if err == nil {
			responseData = []model.Article{*article}
		}
	}
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Article資料查詢成功")
	response := ResponseFactory(c, fiber.StatusOK, "Article資料查詢成功", &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

// create article data
func CreateArticle(c fiber.Ctx, db *gorm.DB) error {
	var clientData ArticleCreateRequest

	if err := c.Bind().Body(&clientData); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if len(clientData.Tags) > 0 {
		count, err := model.CountTags(db, clientData.Tags)
		if err != nil {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		if count != int64(len(clientData.Tags)) {
			logrus.Error("缺少tags，請先建立tags")
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, "缺少tags，請先建立tags", nil)
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

	if err := model.CreateArticle(db, &data); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Article資料建立成功: " + clientData.Title)
	response := ResponseFactory[any](c, fiber.StatusOK, "Article資料建立成功: "+clientData.Title, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// update article data
func UpdateArticle(c fiber.Ctx, db *gorm.DB) error {
	var clientData ArticleUpdateRequest
	if err := c.Bind().Body(&clientData); err != nil {
		response := ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	id := c.Params("id")
	articleID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || articleID == 0 {
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "文章 ID 無效", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if clientData.Title == nil && clientData.Content == nil && clientData.Tags == nil {
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "沒有提供要更新的資料", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var tags []model.Tag
	if clientData.Tags != nil {
		names := make([]string, 0, len(*clientData.Tags))
		seen := make(map[string]bool)
		for _, tag := range *clientData.Tags {
			if strings.TrimSpace(tag.Name) == "" || seen[tag.Name] {
				response := ResponseFactory[any](c, fiber.StatusBadRequest, "標籤名稱不可為空或重複", nil)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
			seen[tag.Name] = true
			names = append(names, tag.Name)
			tags = append(tags, model.Tag{Name: tag.Name})
		}
		if len(names) > 0 {
			count, err := model.CountTags(db, names)
			if err != nil {
				logrus.Error(err)
				response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}
			if count != int64(len(names)) {
				response := ResponseFactory[any](c, fiber.StatusBadRequest, "缺少tags，請先建立tags", nil)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		article, err := model.FindArticleByID(tx, id)
		if err != nil {
			return err
		}
		if clientData.Title != nil {
			article.Title = *clientData.Title
		}
		if clientData.Content != nil {
			article.Content = *clientData.Content
		}
		if err := model.UpdateArticle(tx, article); err != nil {
			return err
		}
		if clientData.Tags != nil {
			return model.ReplaceArticleTags(tx, article, tags)
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		status := fiber.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
		}
		response := ResponseFactory[any](c, status, err.Error(), nil)
		return c.Status(status).JSON(response)
	}

	response := ResponseFactory[any](c, fiber.StatusOK, "Article資料更新成功: "+id, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// Delete article data
func DeleteArticle(c fiber.Ctx, db *gorm.DB) error {

	// URL decoding
	id, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		article, err := model.FindArticleByID(tx, id)
		if err != nil {
			return err
		}
		if err := model.ClearArticleTags(tx, article); err != nil {
			return err
		}
		return model.DeleteArticle(tx, article)
	})
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("刪除Article成功" + id)
	response := ResponseFactory[any](c, fiber.StatusOK, "刪除Article成功: "+id, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// Query Article data use tag name
func QueryArticleForTag(c fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Article

	// URL decoding
	name, err := url.QueryUnescape(c.Params("name"))
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	responseData, err = model.FindArticlesByTag(db, name)
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	logrus.Info("查詢指定Tag的Article成功: " + name)
	response := ResponseFactory(c, fiber.StatusOK, "查詢指定Tag的Article成功"+name, &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
