package blog

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/blog"
	model "seaotterms-api/model/galgame"
	utils "seaotterms-api/utils/blog"
)

type BrandRecordForUpdate struct {
	Brand       string
	Completed   int
	Total       int
	Annotation  string
	Dissolution bool
	UpdateName  string
	UpdateTime  time.Time
}

// query all galgamebrand data
func QueryAllGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	var data []model.SelfBrand

	err := db.Order("update_time DESC").Find(&data).Error
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	logrus.Info("GalgameBrand全部資料查詢成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "GalgameBrand全部資料查詢成功", &data)
	return c.Status(fiber.StatusOK).JSON(response)
}

// use brand name to query single galgamebrand data
func QueryGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	var data []model.SelfBrand
	// URL decoding
	brand, err := url.QueryUnescape(c.Params("brand"))
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = db.Where("brand = ?", brand).First(&data).Error
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到Brand資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Info("GalgameBrand單筆資料查詢成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "GalgameBrand單筆資料查詢成功", &data)
	return c.Status(fiber.StatusOK).JSON(response)
}

// insert data to galgamebrand
func CreateGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData dto.SelfBrandUpdateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	annotation := "待攻略"
	if clientData.Completed == clientData.Total {
		annotation = "制霸"
	}

	data := model.SelfBrand{
		Brand:       clientData.Brand,
		Completed:   clientData.Completed,
		Total:       clientData.Total,
		Annotation:  annotation,
		Dissolution: clientData.Dissolution,
		InputName:   clientData.Username,
		UpdateName:  clientData.Username,
	}
	err := db.Create(&data).Error
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	logrus.Infof("資料 %s 創建成功", clientData.Brand)
	response := utils.ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("資料 %s 創建成功", clientData.Brand), nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// update single galgamebrand data
func UpdateGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData dto.SelfBrandUpdateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// URL decoding
	brand, err := url.QueryUnescape(c.Params("brand"))
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	annotation := "待攻略"
	if clientData.Completed == clientData.Total {
		annotation = "制霸"
	}

	// gorm:"autoUpdateTime" can not update, so manual update update_time
	err = db.Model(&model.SelfBrand{}).Where("brand = ?", brand).
		Select("brand", "completed", "total", "annotation", "dissolution", "update_name", "update_time").
		Updates(BrandRecordForUpdate{
			Brand:       clientData.Brand,
			Completed:   clientData.Completed,
			Total:       clientData.Total,
			Annotation:  annotation,
			Dissolution: clientData.Dissolution,
			UpdateName:  clientData.Username,
			UpdateTime:  time.Now(),
		}).Error
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到Brand資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Infof("資料 %s 更新成功", brand)
	response := utils.ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("資料 %s 更新成功", brand), nil)
	return c.Status(fiber.StatusOK).JSON(response)

}
