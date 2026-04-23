package runningProcess

import (
	"archive/zip"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mskKandula/oes/api/handler"
	ds "github.com/mskKandula/oes/dataSources"
)

// UnzipFile receives exam proof zip paths from the channel, extracts each zip
// into the media directory, records the proof paths in the DB, then removes
// the original zip.
func UnzipFile(ctx context.Context, resultPaths <-chan handler.ProofData, db *ds.DataSources) {

	for result := range resultPaths {
		dir := filepath.Join("../media/examProofs", result.ClientId, result.ExamId, result.UserId)

		reader, err := zip.OpenReader(result.ZipFilePath)
		if err != nil {
			log.Println(err)
			continue
		}

		var vals []interface{}

		for _, file := range reader.File {
			fpath := filepath.Join(dir, file.Name)

			if err := extractFile(file, fpath); err != nil {
				log.Println(err)
				continue
			}

			vals = append(vals, result.UserId, result.ExamId, fpath)
		}

		if err := reader.Close(); err != nil {
			log.Println(err)
		}

		// Nothing was extracted successfully; skip DB insert and zip removal
		if len(vals) == 0 {
			continue
		}

		if err := StudentExamProofsInsertion(ctx, db, vals); err != nil {
			log.Println(err)
			continue
		}

		if err := os.Remove(result.ZipFilePath); err != nil {
			log.Println(err)
		}
	}
}

// extractFile extracts a single zip entry to fpath.
func extractFile(file *zip.File, fpath string) error {
	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	inFile, err := file.Open()
	if err != nil {
		return err
	}
	defer inFile.Close()

	_, err = io.Copy(outFile, inFile)
	return err
}

// StudentExamProofsInsertion bulk-inserts proof records for a student's exam.
func StudentExamProofsInsertion(ctx context.Context, db *ds.DataSources, vals []interface{}) error {

	insertQuery := "INSERT INTO StudentExamProofs(studentId, examId, proofPath) VALUES "

	totalVals := len(vals) / 3
	for i := 0; i < totalVals; i++ {
		insertQuery += "(?,?,?),"
	}

	// Trim the trailing comma
	insertQuery = insertQuery[:len(insertQuery)-1]

	_, err := db.MySQLDB.ExecContext(ctx, insertQuery, vals...)
	return err
}
