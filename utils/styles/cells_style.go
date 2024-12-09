package styles

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// CreateHeaderStyle creates and returns a style for the header row
func CreateHeaderStyle(f *excelize.File) (int, error) {
	style, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["#ff4800"],
            "pattern": 1
        },
        "font": {
            "bold": true,
            "color": "#ffffff"
        },
        "alignment": {
            "horizontal": "center",
            "vertical": "center"
        },
        "border": [
            {
                "type": "left",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "top",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "bottom",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "right",
                "color": "#000000",
                "style": 1
            }
        ]
    }`)
	if err != nil {
		return 0, err
	}
	return style, nil
}

func CreateCellStyle(f *excelize.File) (int, error) {
	// Define the style for the other cells
	cellStyle, err := f.NewStyle(`{
                "fill": {
                    "type": "pattern",
                    "color": ["#ffe9e0"],
                    "pattern": 1
                },
				"alignment": {
                    "horizontal": "center",
                    "vertical": "center"
                },
				"border": [
                    {
                        "type": "left",
                        "color": "#000000",
                        "style": 1
                    },
                    {
                        "type": "top",
                        "color": "#000000",
                        "style": 1
                    },
                    {
                        "type": "bottom",
                        "color": "#000000",
                        "style": 1
                    },
                    {
                        "type": "right",
                        "color": "#000000",
                        "style": 1
                    }
                ]
            }`)
	if err != nil {
		return 0, err
	}
	return cellStyle, nil
}

// FalseCells applies conditional formatting to cells containing "FALSE"
func FalseCells(f *excelize.File) (int, error) {
	// Define the style for the other cells
	falseStyle, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["#f29a4e"],
            "pattern": 1
        },
        "alignment": {
            "horizontal": "center",
            "vertical": "center"
        },
        "border": [
            {
                "type": "left",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "top",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "bottom",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "right",
                "color": "#000000",
                "style": 1
            }
        ]
    }`)
	if err != nil {
		return 0, err
	}
	return falseStyle, err
}

func SetOneHeader(f *excelize.File, sheet string, header string, column rune) {
	headerStyle, _ := CreateHeaderStyle(f)
	cell := fmt.Sprintf("%c1", column)
	f.SetCellValue(sheet, cell, header)
	f.SetCellStyle(sheet, cell, cell, headerStyle)

}
