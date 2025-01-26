package pkg

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func ReadFromFile(file string) (rows [][]string, err error) {

	f, err := excelize.OpenFile(file)
	if err != nil {
		err = errors.Wrap(err, "[pkg.ReadFromFile] OpenFile")
		return
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			err = errors.Wrap(err, "[pkg.ReadFromFile]")
			fmt.Println(err.Error())
		}
	}()

	rows, err = f.GetRows("Sheet1")
	if err != nil {
		err = errors.Wrap(err, "[pkg.ReadFromFile] GetRows")
	}

	return
}
