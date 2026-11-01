package website

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
)

// BulkCreateMakerPieAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) BulkCreateMakerPieAdmin(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	tx, err := wrpr.storage.GetTransaction(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	defer func() {
		tx.Rollback(r.Context())
	}()
	q := wrpr.storage.GetQuerierWithTx(tx)

	dest := "/tmp/" + uuid.NewString()
	_, hdr, err := r.FormFile("zipFile")
	if err != nil {
		AbortErr(w, 400, err)
		return
	}
	err = unzip(hdr, dest)
	if err != nil {
		AbortErr(w, 400, err)
		return
	}

	mappedCsv := []map[string]string{}

	file, _ := os.ReadFile(dest + "/piedata.csv")
	csvReader := csv.NewReader(bytes.NewReader(file))

	headers, err := csvReader.Read()
	if err != nil {
		AbortErr(w, 400, err)
		return
	}

	for {
		nxt, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			AbortErr(w, 400, err)
			return
		}
		if nxt != nil {
			if len(nxt) == 0 {
				//skip empty rows
				continue
			}
			mp := map[string]string{}
			for i, hdr := range headers {
				if i > len(nxt)-1 {

					AbortErr(w, 400, err)
					return
				}
				mp[hdr] = nxt[i]
			}
			mappedCsv = append(mappedCsv, mp)
		}
	}

	for _, pie := range mappedCsv {

		rowLabel := "Pie " + pie["year"] + "/" + pie["makerid"] + "/" + pie["id"] + ": "
		year, yearErr := parseFormYear(pie["year"])
		packCount, packCountErr := parseFormInt("Pack count", pie["pack_count"])
		packPrice, packPriceErr := parseFormInt("Pack price in pence", pie["pack_price_in_pence"])
		if errs := formFieldErrors(yearErr, packCountErr, packPriceErr); len(errs) > 0 {
			piesPage(wrpr, w, r, rowLabel+strings.Join(errs, "; "))
			return
		}

		dirpath := strings.Join([]string{
			wrpr.config.ImgprssrDir,
			pie["year"],
			pie["makerid"],
			pie["id"],
		}, "/")

		img, err := os.Open(dest + "/" + pie["image_file"])
		if err != nil {

			AbortErr(w, 400, err)
			return
		}
		defer img.Close()

		rf, err := os.ReadFile(dest + "/" + pie["image_file"])
		if err != nil {

			AbortErr(w, 400, err)
			return
		}

		prts := strings.Split(img.Name(), "/")
		path := strings.Join([]string{
			pie["year"],
			pie["makerid"],
			pie["id"],
			prts[len(prts)-1],
		}, "/")

		fresh := pie["fresh"] == "yes"
		oid, err := q.UpsertMakerPie(r.Context(), sqlcgen.UpsertMakerPieParams{
			Year:    year,
			Makerid: pie["makerid"],
			ID:      pie["id"],

			Displayname:      pie["displayname"],
			Fresh:            fresh,
			ImageFile:        "/" + path,
			WebLink:          pie["web_link"],
			PackCount:        packCount,
			PackPriceInPence: packPrice,
			Validated:        true,
		})
		if err != nil {
			AbortErr(w, 400, err)
			return
		}
		err = wrpr.storage.SetPieCategoriesTx(r.Context(), tx, uuid.MustParse(oid), strings.Split(pie["category_slugs"], ","))
		if errors.Is(err, storage.ErrUnknownCategories) {
			piesPage(wrpr, w, r, "Pie "+pie["year"]+"/"+pie["makerid"]+"/"+pie["id"]+": "+err.Error())
			return
		}
		if err != nil {
			AbortErr(w, 400, err)
			return
		}

		err = os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			AbortErr(w, 500, err)
			return
		}

		err = os.WriteFile(wrpr.config.ImgprssrDir+"/"+path, rf, fs.ModeAppend)
		if err != nil {
			AbortErr(w, 500, err)
			return
		}

	}
	tx.Commit(r.Context())

	http.Redirect(w, r, "/admin/makerpies/", 303)
}

func unzip(zipfile *multipart.FileHeader, dest string) error {
	fl, err := zipfile.Open()
	if err != nil {
		return err
	}
	byteContainer, err := io.ReadAll(fl)
	if err != nil {
		return err
	}

	r, err := zip.NewReader(bytes.NewReader(byteContainer), int64(len(byteContainer)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(dest + "/piedata.csv")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
