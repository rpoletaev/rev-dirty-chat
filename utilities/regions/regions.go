package regions

import (
	"container/list"
	"github.com/tealeg/xlsx"
)

func GetRegionsFromFile(excelFileName string) (*list.List, error) {
	regs := list.New()
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return nil, err
	}

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			regs.PushBack(row.Cells[1].Value)
		}
	}
	return regs, nil
}
