package website

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
)

// BulkCreateMakerPieAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) BulkCreateMakerPieAdmin(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	tx, err := wrpr.storage.GetTransaction(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer func() {
		tx.Rollback(c)
	}()
	q := wrpr.storage.GetQuerierWithTx(tx)

	dest := "/tmp/" + uuid.NewString()
	hdr, err := c.FormFile("zipFile")
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	err = unzip(hdr, dest)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	mappedCsv := []map[string]string{}

	file, _ := os.ReadFile(dest + "/piedata.csv")
	csvReader := csv.NewReader(bytes.NewReader(file))

	headers, err := csvReader.Read()
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	for {
		nxt, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.AbortWithError(400, err)
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

					c.AbortWithError(400, err)
					return
				}
				mp[hdr] = nxt[i]
			}
			mappedCsv = append(mappedCsv, mp)
		}
	}

	for _, pie := range mappedCsv {

		dirpath := strings.Join([]string{
			wrpr.config.ImgprssrDir,
			pie["year"],
			pie["makerid"],
			pie["id"],
		}, "/")

		img, err := os.Open(dest + "/" + pie["image_file"])
		if err != nil {

			c.AbortWithError(400, err)
			return
		}
		defer img.Close()

		rf, err := os.ReadFile(dest + "/" + pie["image_file"])
		if err != nil {

			c.AbortWithError(400, err)
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

		oid, err := q.CreateMakerPie(c, sqlcgen.CreateMakerPieParams{
			Year:    mustParseInt(pie["year"]),
			Makerid: pie["makerid"],
			ID:      pie["id"],

			Displayname:      pie["displayname"],
			Fresh:            fresh,
			ImageFile:        "/" + path,
			WebLink:          pie["web_link"],
			PackCount:        mustParseInt(pie["pack_count"]),
			PackPriceInPence: mustParseInt(pie["pack_price_in_pence"]),
			Validated:        true,
		})
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = wrpr.storage.SetPieCategoriesTx(c, tx, uuid.MustParse(oid), strings.Split(pie["category_slugs"], ","))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		err = os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		err = os.WriteFile(wrpr.config.ImgprssrDir+"/"+path, rf, fs.ModeAppend)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

	}
	tx.Commit(c)

	c.Redirect(303, "/admin/makerpies/")
}

// Unzip function extracts a ZIP archive to a target directory
func unzip(zipfile *multipart.FileHeader, dest string) error {
	// Open the zip archive for reading
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

	// Iterate through each file in the archive
	for _, f := range r.File {
		// Create the full path for the extracted file
		fpath := filepath.Join(dest, f.Name)

		// Check for directory traversal vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create directories if necessary
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create the file's parent directories if necessary
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Open the file inside the zip archive
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		// Create the output file on disk
		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Copy the contents of the file in the archive to the new file on disk
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
	}

	// Check for CSV
	_, err = os.Stat(dest + "/piedata.csv")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
