package website

import (
	"context"

	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

// AboutPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) AboutPage(c context.Context, request generated.AboutPageRequestObject) (generated.AboutPageResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "About",
			Description: "All about what Mince Pie Rank is",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("static:about", vals)
	if err != nil {
		return nil, err
	}
	return generated.AboutPage200TexthtmlResponse{
		Body: rdr,
	}, nil
}

// CookiePage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) CookiePage(c context.Context, request generated.CookiePageRequestObject) (generated.CookiePageResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Cookies",
			Description: "Mince Pie Rank cookie policy",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("static:cookies", vals)
	if err != nil {
		return nil, err
	}
	return generated.CookiePage200TexthtmlResponse{
		Body: rdr,
	}, nil
}

// PrivacyPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) PrivacyPage(c context.Context, request generated.PrivacyPageRequestObject) (generated.PrivacyPageResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Privacy",
			Description: "Mince Pie Rank's privacy policy",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("static:privacy", vals)
	if err != nil {
		return nil, err
	}
	return generated.PrivacyPage200TexthtmlResponse{
		Body: rdr,
	}, nil
}
