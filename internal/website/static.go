package website

import (
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
	"net/http"
)

// AboutPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) AboutPage(w http.ResponseWriter, r *http.Request) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       "About",
			Description: "All about what Mince Pie Rank is",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        signedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
	}

	wrpr.renderPage(w, 200, "static:about", vals)
}

// CookiePage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) CookiePage(w http.ResponseWriter, r *http.Request) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       "Cookies",
			Description: "Mince Pie Rank cookie policy",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        signedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
	}

	wrpr.renderPage(w, 200, "static:cookies", vals)
}

// PrivacyPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) PrivacyPage(w http.ResponseWriter, r *http.Request) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       "Privacy",
			Description: "Mince Pie Rank's privacy policy",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        signedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
	}

	wrpr.renderPage(w, 200, "static:privacy", vals)
}
