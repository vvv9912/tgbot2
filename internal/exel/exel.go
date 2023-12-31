package exel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

// todo
type Excel struct {
	nameFile string
}

func NewExcel(nameFile string) *Excel {
	return &Excel{nameFile: nameFile}
}

func (e *Excel) Read() *[]ProductsPars {
	//f, err := excelize.OpenFile("test_bd/База данных.xlsx")
	f, err := excelize.OpenFile(e.nameFile)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.GetRows("Лист1")
	products := make([]ProductsPars, len(rows)-1)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for i := 1; i < len(rows); i++ {
		fmt.Println(len(rows[i]))
		if len(rows[i]) > 10 || len(rows) < 9 {
			fmt.Println(err)
			break
		}
		products[i-1].Article, err = strconv.Atoi(rows[i][0])
		if err != nil {
			fmt.Println("strconv  Atoi stroka %v, err: %v", i, err)
			return nil
		}
		products[i-1].Catalog = rows[i][1]
		fmt.Println(rows[i][1])
		fmt.Println(products[i-1].Catalog)
		products[i-1].Name = rows[i][2]
		products[i-1].Description = rows[i][3]

		products[i-1].Price, err = strconv.ParseFloat(rows[i][4], 64)
		if err != nil {
			fmt.Println("ParseFloat stroka %v, err: %v", i, err)
			return nil
		}
		products[i-1].Length, err = strconv.Atoi(rows[i][5])
		if err != nil {
			fmt.Println("strconv  Atoi stroka %v, err: %v", i, err)
			return nil
		}
		products[i-1].Width, err = strconv.Atoi(rows[i][6])
		if err != nil {
			fmt.Println("strconv  Atoi stroka %v, err: %v", i, err)
			return nil
		}
		products[i-1].Height, err = strconv.Atoi(rows[i][7])
		if err != nil {
			fmt.Println("strconv  Atoi stroka %v, err: %v", i, err)
			return nil
		}
		products[i-1].Weight, err = strconv.Atoi(rows[i][8])
		if err != nil {
			fmt.Println("strconv  Atoi stroka %v, err: %v", i, err)
			return nil
		}
		if len(rows[i]) > 9 {
			a := rows[i][9]
			b := strings.Split(a, ",")
			products[i-1].PhotoUrl = b
		}
	}
	//for _, row := range rows {
	//	for _, colCell := range row {
	//		fmt.Print(colCell, "\t")
	//	}
	//	fmt.Println()
	//}
	return &products
}

type ProductsPars struct {
	Article     int      `json:"article,omitempty"`
	Catalog     string   `json:"catalog,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	PhotoUrl    []string `json:"photo_url,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Length      int      `json:"length"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Weight      int      `json:"weight"`
}
