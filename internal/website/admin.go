package website

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func piesPage(wrpr *WebsiteWrapper, c *gin.Context) {
	pies, err := wrpr.storage.GetQuerier().GetAllMakerPies(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Admin Pies",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear: ay,
			},
		},
		PageData: map[string]interface{}{
			"Pies": pies,
		},
	}

	c.HTML(http.StatusOK, "page:admin:pies", vals)
}

func piePage(wrpr *WebsiteWrapper, c *gin.Context, oid string) {
	pie, err := wrpr.storage.GetQuerier().GetMakerPieByOid(c, oid)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Admin Pie",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear: ay,
			},
		},
		PageData: map[string]interface{}{
			"Pie": pie,
		},
	}

	c.HTML(http.StatusOK, "page:admin:pie", vals)
}

// GetMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetMakerPieAdmin(c *gin.Context, oid string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	piePage(wrpr, c, oid)
}

// GetMakerPiesAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetMakerPiesAdmin(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	piesPage(wrpr, c)
}

// CreateMakerPieAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateMakerPieAdmin(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	img, err := c.FormFile("img")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	path := strings.Join([]string{
		c.Request.FormValue("year"),
		c.Request.FormValue("makerid"),
		c.Request.FormValue("id"),
		img.Filename,
	}, "/")

	dirpath := strings.Join([]string{
		wrpr.config.ImgprssrDir,
		c.Request.FormValue("year"),
		c.Request.FormValue("makerid"),
		c.Request.FormValue("id"),
	}, "/")

	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	fl, err := img.Open()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	byteContainer, err := io.ReadAll(fl)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	err = os.WriteFile(wrpr.config.ImgprssrDir+"/"+path, byteContainer, fs.ModeAppend)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	oid, err := wrpr.storage.GetQuerier().CreateMakerPie(c, sqlcgen.CreateMakerPieParams{
		Year:    mustParseInt(c.Request.FormValue("year")),
		Makerid: c.Request.FormValue("makerid"),
		ID:      c.Request.FormValue("id"),

		Displayname:      c.Request.FormValue("displayname"),
		Fresh:            mustParseBool(c.Request.FormValue("fresh")),
		ImageFile:        "/" + path,
		WebLink:          c.Request.FormValue("web_link"),
		PackCount:        mustParseInt(c.Request.FormValue("pack_count")),
		PackPriceInPence: mustParseInt(c.Request.FormValue("pack_price_in_pence")),
		Validated:        mustParseBool(c.Request.FormValue("validated")),
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Redirect(301, "/admin/makerpies/"+oid)
}

// DeleteMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteMakerPie(c *gin.Context, oid string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	err := wrpr.storage.GetQuerier().DeleteMakerPieByOid(c, oid)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesPage(wrpr, c)
}

// UpdateMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateMakerPie(c *gin.Context, oid string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	err := wrpr.storage.GetQuerier().UpdateMakerPieByOid(c, sqlcgen.UpdateMakerPieByOidParams{
		Displayname:      c.Request.FormValue("displayname"),
		Fresh:            mustParseBool(c.Request.FormValue("fresh")),
		WebLink:          c.Request.FormValue("web_link"),
		PackCount:        mustParseInt(c.Request.FormValue("pack_count")),
		PackPriceInPence: mustParseInt(c.Request.FormValue("pack_price_in_pence")),
		Validated:        mustParseBool(c.Request.FormValue("validated")),
		Oid:              oid,
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piePage(wrpr, c, oid)
}

func mustParseInt(str string) int32 {
	i, _ := strconv.ParseInt(str, 10, 32)
	return int32(i)
}

func mustParseBool(str string) bool {
	i, _ := strconv.ParseBool(str)
	return i
}
