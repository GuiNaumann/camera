package utils

import (
	"bear/infrastructure/modules/impl/http_error"
	"encoding/csv"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"log"
)

type GridFileDecoder struct {
}

func (d GridFileDecoder) DecodeCSV(file io.Reader) ([][]string, error) {
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Println("[DecodeXLSX] Error ReadAll", err)
		err = http_error.NewBadRequestError(http_error.Unexpected)
		return nil, err
	}

	return records, nil
}

func (d GridFileDecoder) DecodeXLSX(file io.Reader) ([][]string, error) {
	grid := make([][]string, 0)

	xlsxFile, err := excelize.OpenReader(file)
	if err != nil {
		log.Println("[Decode] Error OpenReader", err)
		return nil, err
	}

	sheets := xlsxFile.GetSheetMap()
	if len(sheets) > 1 {
		log.Println("[Decode] Error GetSheetMap", err)
		return nil, http_error.NewBadRequestError(http_error.MustHaveOnlyOneSheet)

	}

	var rows *excelize.Rows
	rows, err = xlsxFile.Rows(sheets[1])
	if err != nil {
		log.Println("[DecodeCSV] Error Rows", err)
		http_error.NewBadRequestError(http_error.Unexpected)
		return nil, err
	}

	for rows.Next() {
		grid = append(grid, rows.Columns())
	}

	return grid, nil
}
