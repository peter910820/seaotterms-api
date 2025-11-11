package gal

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"

	dto "seaotterms-api/dto/gal"
	model "seaotterms-api/model/galgame"
	utils "seaotterms-api/utils/gal"
)

func Login(c *fiber.Ctx, db *gorm.DB) error {
	var clientData dto.LoginRequest
	var responseData dto.CommonResponse
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responseData)
	}
	// check user agent
	if !utils.CheckUserAgent(c) {
		return c.Status(fiber.StatusForbidden).JSON(responseData)
	}
	// 檢查ip 決定要不要發信警告
	// ip := c.IP()

	return c.Status(fiber.StatusOK).JSON(responseData)
}

func Register(c *fiber.Ctx, db *gorm.DB) error {
	var clientData dto.RegisterRequest
	var responseData dto.CommonResponse
	timeNow := time.Now()
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responseData)
	}

	var find model.User
	err := db.Where("email = ? OR user_name = ?", clientData.Email, clientData.UserName).First(&find).Error
	if err == nil {
		responseData.ErrMsg = "帳號或使用者名稱已被使用"
		return c.Status(fiber.StatusInternalServerError).JSON(responseData)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Unexpected Error
		logrus.Println("查詢錯誤:", err)
		responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
		return c.Status(fiber.StatusInternalServerError).JSON(responseData)
	} else {
		// can register
		// generate email check key
		key, err := utils.GenerateRandomKey(16)
		if err != nil {
			logrus.Println("錯誤:", err)
			responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
			return c.Status(fiber.StatusInternalServerError).JSON(responseData)
		}
		err = utils.WriteTmpData("RegisterKey", key, timeNow.Add(24*time.Hour), db)
		if err != nil {
			logrus.Println("錯誤:", err)
			responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
			return c.Status(fiber.StatusInternalServerError).JSON(responseData)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(clientData.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Println("錯誤:", err)
			responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
			return c.Status(fiber.StatusInternalServerError).JSON(responseData)
		}
		// write user data to db
		clientData.Password = string(hashedPassword)
		dataCreate := model.User{
			Email:    clientData.Email,
			UserName: clientData.UserName,
			Password: clientData.Password,
			SignupIP: c.IP(),
		}
		err = db.Create(&dataCreate).Error
		if err != nil {
			logrus.Println("錯誤:", err)
			responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
			return c.Status(fiber.StatusInternalServerError).JSON(responseData)
		}
		go utils.SendRegisterEmail(clientData.Email, key)
		responseData.InfoMsg = "註冊申請已送出，請檢查電子郵件並且點擊認證連結完成註冊(信件有時會被判定為垃圾郵件，若未收到請檢查垃圾郵件)"
		return c.Status(fiber.StatusOK).JSON(responseData)
	}
}

func RegisterKeyCheck(c *fiber.Ctx, db *gorm.DB) error {
	var responseData dto.CommonResponse
	// URL decoding
	mailName, err := url.QueryUnescape(c.Params("mail_name"))
	if err != nil || strings.TrimSpace(mailName) == "" {
		logrus.Error(err)
		responseData.ErrMsg = "Not allowed"
		return c.Status(fiber.StatusBadRequest).JSON(responseData)
	}
	registerKey, err := url.QueryUnescape(c.Params("register_key"))
	if err != nil || strings.TrimSpace(mailName) == "" {
		logrus.Error(err)
		responseData.ErrMsg = "Not allowed"
		return c.Status(fiber.StatusBadRequest).JSON(responseData)
	}

	var find model.TmpData
	err = db.Where("content = ?", registerKey).First(&find).Error
	if err == nil {
		var userData model.User
		err = db.Where("email = ? and management = -2", mailName+"@gmail.com").First(&userData).Error
		if err == nil {
			err = updateUserForRegister(mailName+"@gmail.com", db)
			if err != nil {
				logrus.Error(err)
				responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
				return c.Status(fiber.StatusInternalServerError).JSON(responseData)
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Unexpected Error
			logrus.Error(err)
			responseData.ErrMsg = "Not allowed"
			return c.Status(fiber.StatusBadRequest).JSON(responseData)
		} else {
			// Unexpected Error
			logrus.Error(err)
			responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
			return c.Status(fiber.StatusInternalServerError).JSON(responseData)
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		responseData.ErrMsg = "Not allowed"
		return c.Status(fiber.StatusBadRequest).JSON(responseData)
	} else {
		// Unexpected Error
		logrus.Error(err)
		responseData.ErrMsg = "未預期的錯誤，請聯繫管理員"
		return c.Status(fiber.StatusInternalServerError).JSON(responseData)
	}

	responseData.InfoMsg = "註冊成功！"
	return c.Status(fiber.StatusOK).JSON(responseData)
}

func updateUserForRegister(userEmail string, db *gorm.DB) error {
	type UpdateUser struct {
		Management int
		UpdatedAt  time.Time
	}
	updateData := UpdateUser{
		Management: 0,
		UpdatedAt:  time.Now(),
	}

	err := db.Model(&model.User{}).Where("email = ?", userEmail).
		Select("management", "updated_at").
		Updates(updateData).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
