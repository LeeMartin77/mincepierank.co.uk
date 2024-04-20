package website

import (
	"context"
	"io"
	"os"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

type Website struct {
	serverInterface generated.StrictServerInterface
}

func NewWebsite(out io.Writer) (*Website, error) {
	ops, err := storage.NewOperations(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	err = ops.Migrate(context.Background())
	if err != nil {
		return nil, err
	}
	fsi := NewWebsiteWrapper(ops)

	return &Website{
		serverInterface: fsi,
	}, nil
}

type WebsiteWrapper struct {
	storage       storage.Operations
	htmlTemplater templater.Templater
}

// YearAllBrands implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllBrands(ctx context.Context, request generated.YearAllBrandsRequestObject) (generated.YearAllBrandsResponseObject, error) {
	panic("unimplemented")
}

// YearAllCategories implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllCategories(ctx context.Context, request generated.YearAllCategoriesRequestObject) (generated.YearAllCategoriesResponseObject, error) {
	panic("unimplemented")
}

// YearAllPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllPies(ctx context.Context, request generated.YearAllPiesRequestObject) (generated.YearAllPiesResponseObject, error) {
	panic("unimplemented")
}

// YearBrandPie implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearBrandPie(ctx context.Context, request generated.YearBrandPieRequestObject) (generated.YearBrandPieResponseObject, error) {
	panic("unimplemented")
}

// YearBrandPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearBrandPies(ctx context.Context, request generated.YearBrandPiesRequestObject) (generated.YearBrandPiesResponseObject, error) {
	panic("unimplemented")
}

// YearCategoryPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearCategoryPies(ctx context.Context, request generated.YearCategoryPiesRequestObject) (generated.YearCategoryPiesResponseObject, error) {
	panic("unimplemented")
}

// YearPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearPage(ctx context.Context, request generated.YearPageRequestObject) (generated.YearPageResponseObject, error) {
	panic("unimplemented")
}

// YearsPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearsPage(ctx context.Context, request generated.YearsPageRequestObject) (generated.YearsPageResponseObject, error) {
	panic("unimplemented")
}

func NewWebsiteWrapper(storage storage.Operations) generated.StrictServerInterface {

	tmpltr := templater.NewHtmlTemplater()
	return &WebsiteWrapper{
		storage:       storage,
		htmlTemplater: tmpltr,
	}
}
