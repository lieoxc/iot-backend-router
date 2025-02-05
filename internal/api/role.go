package api

import (
	"net/http"
	model "project/internal/model"
	service "project/internal/service"
	"project/pkg/errcode"
	utils "project/pkg/utils"

	"github.com/gin-gonic/gin"
)

type RoleApi struct{}

// CreateRole 创建角色管理
// @Router   /api/v1/role [post]
func (*RoleApi) CreateRole(c *gin.Context) {
	var req model.CreateRoleReq
	if !BindAndValidate(c, &req) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.Role.CreateRole(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// UpdateRole 更新角色管理
// @Router   /api/v1/role [put]
func (*RoleApi) UpdateRole(c *gin.Context) {
	var req model.UpdateRoleReq
	if !BindAndValidate(c, &req) {
		return
	}

	if req.Description == nil && req.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "修改内容不能为空"})
		return
	}

	data, err := service.GroupApp.Role.UpdateRole(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// DeleteRole 删除角色管理
// @Router   /api/v1/role/{id} [delete]
func (*RoleApi) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	// 需要角色没有被用户使用
	if service.GroupApp.Casbin.HasRole(id) {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"role_id": id,
			"error":   "Role in use",
		}))
		return
	}

	err := service.GroupApp.Role.DeleteRole(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetRoleListByPage 角色管理分页查询
// @Router   /api/v1/role [get]
func (*RoleApi) HandleRoleListByPage(c *gin.Context) {
	var req model.GetRoleListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	roleList, err := service.GroupApp.Role.GetRoleListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", roleList)
}
