package api_operation

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/services"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) userCreate(c *gin.Context) {

	var params model.User
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	if err := vars.Database.Where("username = ?", params.Username).First(&params).Error; err == nil {
		c.String(500, "duplicated username")
		return
	}

	params.Hashed = utils.Sha512encode("P@ssw0rd")

	err := vars.Database.Create(&params).Error
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(201, "text/plain", nil)

}

func (ctl *controller) userUpdate(c *gin.Context) {

	var params model.User
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	err := vars.Database.Model(&model.User{}).Where("id = ?", params.ID).Update("name", params.Name).Update("user_id", params.UserID).Update("enable", params.Enable).Error
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(204, "text/plain", nil)

}

func (ctl *controller) userPasswordReset(c *gin.Context) {

	var params model.User
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	err := vars.Database.Model(&model.User{}).Where("id = ?", params.ID).Update("hashed", utils.Sha512encode("P@ssw0rd")).Error
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(204, "text/plain", nil)

}

func (ctl *controller) userDelete(c *gin.Context) {

	id := c.Query("id")
	if id == "" {
		c.String(400, "id required.")
	}

	err := vars.Database.Where("id = ?", id).Delete(&model.User{}).Error
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(204, "text/plain", nil)

}

func (ctl *controller) userQuery(c *gin.Context) {

	search := c.Query("search")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	or := c.Query("or")

	offset, e := strconv.Atoi(offsetStr)
	if e != nil {
		c.String(400, "offset is not correct, allow only integer")
		return
	}

	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		c.String(400, "limit is not correct, allow only integer")
		return
	}

	response := []model.User{}
	count, err := services.DBQuery(&response, offset, limit, search, or == "true", false,
		services.DBQueryPreload{
			Name: "Directory",
		},
	)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}
