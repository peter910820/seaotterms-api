package teach

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms-db/teach"
)

type linkUpdateSeries struct {
	ArticleAmount uint
	UpdateTime    time.Time
}

func queryArticle(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []teach.Article
	var err error

	id := c.Query("id")
	seriesID := c.Query("series-id")

	if id == "" && seriesID == "" {
		responseData, err = teach.FindAllArticle(db)
		if err != nil {
			return dbErrHandler(c, err)
		}
	} else if seriesID != "" {
		responseData, err = teach.FindArticleBySeriesID(db, seriesID)
		if err != nil {
			return dbErrHandler(c, err)
		}
	} else {
		data, err := teach.FindArticleByID(db, id)
		if err != nil {
			return dbErrHandler(c, err)
		}
		responseData = append(responseData, *data)
	}

	logrus.Info("Query article table success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "查詢Article資料成功",
		"data": responseData,
	})
}

func createArticle(c *fiber.Ctx, db *gorm.DB) error {
	var clientData ArtilceCreateRequest
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// confirm series exists
	var seriesData teach.Series
	seriesID := strconv.Itoa(int(clientData.SeriesID))
	_, err := teach.FindSeriesByID(db, seriesID)
	if err != nil {
		return dbErrHandler(c, err)
	}

	if strings.Trim(clientData.Image, " ") == "" {
		clientData.Image = seriesData.Image
	}

	data := teach.Article{
		Title:       clientData.Title,
		Image:       clientData.Image,
		SeriesID:    clientData.SeriesID,
		Tags:        clientData.Tags,
		Content:     clientData.Content,
		CreatedAt:   time.Now(),
		CreatedName: "Root",
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. create article
		if err := teach.CreateArticle(tx, &data); err != nil {
			return err
		}

		// 2. update series
		err := tx.Model(&teach.Series{}).
			Where("id = ?", clientData.SeriesID).
			Select("article_amount", "updated_at").
			Updates(linkUpdateSeries{
				ArticleAmount: seriesData.ArticleAmount + 1,
				UpdateTime:    time.Now(),
			}).Error

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dbErrHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "新增Article資料成功",
		"data": data,
	})
}

func modifyArticle(c *fiber.Ctx, db *gorm.DB) error {
	var clientData ArtilceModifyRequest
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// confirm series exists
	var seriesData teach.Series
	seriesID := strconv.Itoa(int(clientData.SeriesID))
	_, err := teach.FindSeriesByID(db, seriesID)
	if err != nil {
		return dbErrHandler(c, err)
	}

	if strings.Trim(clientData.Image, " ") == "" {
		clientData.Image = seriesData.Image
	}

	data := teach.Article{
		Title:       clientData.Title,
		Image:       clientData.Image,
		SeriesID:    clientData.SeriesID,
		Tags:        clientData.Tags,
		Content:     clientData.Content,
		UpdatedAt:   time.Now(),
		UpdatedName: "Root",
	}

	err = teach.UpdateArticle(db, c.Params("id"), &data)
	if err != nil {
		return dbErrHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "修改Article資料成功",
		"data": data,
	})
}
