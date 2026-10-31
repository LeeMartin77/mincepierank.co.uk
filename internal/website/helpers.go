package website

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/rs/zerolog/log"
)

// abortStatus writes a plain status response and stops further handling -
// the net/http equivalent of gin's AbortWithStatus.
func abortStatus(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

// abortErr logs err and responds with code - the net/http equivalent of
// gin's AbortWithError, which hid the error in its own internal log.
func abortErr(w http.ResponseWriter, code int, err error) {
	log.Error().Err(err).Int("status", code).Msg("request failed")
	http.Error(w, http.StatusText(code), code)
}

// saveUploadedFile writes the multipart file from formField to
// imgprssrDir/<pathParts...>/<filename> and returns the public path (with
// leading slash), the same convention pie images and brand logos use.
// When the form carries no file for that field it returns ("", false, nil) so
// optional uploads can be detected by the caller.
func saveUploadedFile(imgprssrDir, formField string, r *http.Request, pathParts ...string) (string, bool, error) {
	_, header, err := r.FormFile(formField)
	// A url-encoded form (no multipart at all) is just as much "no file" as an
	// empty file field; only real read/write failures are errors.
	if err == http.ErrMissingFile || err == http.ErrNotMultipart {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	// Only ever use the base name: form filenames can contain path separators.
	name := filepath.Base(header.Filename)
	if name == "." || name == string(filepath.Separator) {
		return "", false, http.ErrMissingFile
	}

	relPath := strings.Join(append(append([]string{}, pathParts...), name), "/")
	dirPath := filepath.Join(append([]string{imgprssrDir}, pathParts...)...)
	if err := os.MkdirAll(dirPath, fs.ModePerm); err != nil {
		return "", false, err
	}

	fl, err := header.Open()
	if err != nil {
		return "", false, err
	}
	defer fl.Close()

	contents, err := io.ReadAll(fl)
	if err != nil {
		return "", false, err
	}
	if err := os.WriteFile(filepath.Join(imgprssrDir, relPath), contents, 0o644); err != nil {
		return "", false, err
	}
	return "/" + relPath, true, nil
}

// snowing reports whether the public site should show snowfall. The snow is
// the signal that the year is active, so it shows exactly while the site is
// not in read-only mode. On a config read failure it fails closed: no snow.
func (wrpr *WebsiteWrapper) snowing(c context.Context) bool {
	ro, err := wrpr.storage.GetReadonly(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly for snow state")
		return false
	}
	return !ro
}

func (wrpr *WebsiteWrapper) getMakerMap(c context.Context, year int64) (map[string]types.Maker, error) {
	mkrs, err := wrpr.storage.GetMakersForYear(c, year)
	if err != nil {
		return nil, err
	}
	mkrMap := map[string]types.Maker{}
	for _, mkr := range *mkrs {
		mkrMap[mkr.Id] = mkr
	}
	return mkrMap, nil
}
