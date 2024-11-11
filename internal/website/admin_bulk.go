package website

import (
	"archive/zip"
	"bytes"
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

	dirpath := strings.Join([]string{
		wrpr.config.ImgprssrDir,
		c.Request.FormValue("year"),
		c.Request.FormValue("makerid"),
		c.Request.FormValue("id"),
	}, "/")

	path := strings.Join([]string{
		c.Request.FormValue("year"),
		c.Request.FormValue("makerid"),
		c.Request.FormValue("id"),
		img.Filename,
	}, "/")

	oid, err := q.CreateMakerPie(c, sqlcgen.CreateMakerPieParams{
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

	err = wrpr.storage.SetPieCategoriesTx(c, tx, uuid.MustParse(oid), c.Request.Form["categories"])
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = os.WriteFile(wrpr.config.ImgprssrDir+"/"+path, byteContainer, fs.ModeAppend)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(303, "/admin/makerpies/")
}

// Unzip function extracts a ZIP archive to a target directory
func unzip(src *multipart.FileHeader, dest string) error {
	// Open the zip archive for reading
	zipfile, err := c.FormFile("zipFile")
	if err != nil {
		return err
	}

	fl, err := zipfile.Open()
	if err != nil {
		return err
	}
	defer zipfile.Close()
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
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {

	}
	return nil
}
