package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// Color constants
const (
	// dark blue
	HeaderColor = "#2c3e50"

	// red
	LowCountColor = "#ff4800"

	// orange
	MediumCountColor = "#f29a4e"

	// green
	HighCountColor = "#4fd246"

	// white
	RegularCellColor = "#ffffff"
)

// CreateHeaderStyle creates and returns a style for the header row
func CreateHeaderStyle(f *excelize.File) (int, error) {
	style, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["` + HeaderColor + `"],
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
                    "color": ["` + RegularCellColor + `"],
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

// FalseCells applies style to cells containing "FALSE"
func FalseCells(f *excelize.File) (int, error) {
	falseStyle, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["` + LowCountColor + `"],
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

// TrueCells applies style to cells containing "TRUE"
func TrueCells(f *excelize.File) (int, error) {
	trueStyle, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["` + HighCountColor + `"],
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
	return trueStyle, err
}

// LowCountStyle creates a style for low count cells
func LowCountStyle(f *excelize.File) (int, error) {
	lowCountStyle, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["` + LowCountColor + `"],
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
	return lowCountStyle, err
}

// MediumCountStyle creates a style for medium count cells (orange)
func MediumCountStyle(f *excelize.File) (int, error) {
	mediumCountStyle, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["` + MediumCountColor + `"],
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
	return mediumCountStyle, err
}

// HighCountStyle creates a style for high count cells (green)
func HighCountStyle(f *excelize.File) (int, error) {
	highCountStyle, err := f.NewStyle(`{
        "fill": {
            "type": "pattern",
            "color": ["` + HighCountColor + `"],
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
	return highCountStyle, err
}

func SetOneHeader(f *excelize.File, sheet string, header string, column rune) {
	headerStyle, _ := CreateHeaderStyle(f)
	cell := fmt.Sprintf("%c1", column)
	f.SetCellValue(sheet, cell, header)
	f.SetCellStyle(sheet, cell, cell, headerStyle)
}
