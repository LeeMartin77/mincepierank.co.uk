package website

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
)

// will abort with a 403 if not admin
func validateAdmin(wrpr *WebsiteWrapper, c *gin.Context) bool {
	userId := c.Keys["userid"].(string)
	isAdmin, err := wrpr.storage.GetQuerier().IsAdminId(c, userId)
	if err != nil {
		c.AbortWithError(500, err)
		return false
	}
	if !isAdmin {
		c.AbortWithStatus(403)
		return false
	}
	return true
}

// GetAdminIndex implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminIndex(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Admin Index",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear: ay,
			},
		},
		PageData: map[string]interface{}{},
	}

	c.HTML(http.StatusOK, "page:admin:index", vals)
}

func configPage(wrpr *WebsiteWrapper, c *gin.Context) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	config, err := wrpr.storage.GetQuerier().GetAllConfig(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Admin Config",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear: ay,
			},
		},
		PageData: map[string]interface{}{
			"Config": config,
		},
	}

	c.HTML(http.StatusOK, "page:admin:config", vals)
}

// GetAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminConfig(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	configPage(wrpr, c)
}

// CreateAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateAdminConfig(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	err := wrpr.storage.GetQuerier().InsertConfig(c, sqlcgen.InsertConfigParams{
		Key:   c.Request.FormValue("key"),
		Value: c.Request.FormValue("value"),
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	configPage(wrpr, c)
}

// DeleteAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminConfig(c *gin.Context, key string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	err := wrpr.storage.GetQuerier().DeleteConfig(c, key)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	configPage(wrpr, c)
}

// UpdateAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateAdminConfig(c *gin.Context, key string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	err := wrpr.storage.GetQuerier().UpdateConfig(c, sqlcgen.UpdateConfigParams{
		Key:   key,
		Value: c.Request.FormValue("value"),
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	configPage(wrpr, c)
}
