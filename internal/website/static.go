package website

import (
	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// AboutPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) AboutPage(c *gin.Context) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "About",
			Description: "All about what Mince Pie Rank is",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
	}

	c.HTML(200, "page:static:about", vals)
}

// CookiePage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) CookiePage(c *gin.Context) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Cookies",
			Description: "Mince Pie Rank cookie policy",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
	}

	c.HTML(200, "page:static:cookies", vals)
}

// PrivacyPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) PrivacyPage(c *gin.Context) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Privacy",
			Description: "Mince Pie Rank's privacy policy",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
	}

	c.HTML(200, "page:static:privacy", vals)
}
