package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	Account_Type   string `json:"account_type"`
	Account_Number string `json:"account_number"`
	Trans_Date     string `json:"trans_date"`
	Cheque_no      string `json:"cheque_no"`
	Description_1  string `json:"description_1"`
	Description_2  string `json:"description_2"`
	CAD            string `json:"CAD"`
	USD            string `json:"USD"`
}

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Barath2002@tcp(127.0.0.1:3306)/tax_schema?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db

}

func CSVReadWriter() {
	var db = getMySQLDB()
	csvFileName := "csv44134.csv"
	csvFileRead, err := os.Open(csvFileName)

	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(csvFileRead)
	reader.FieldsPerRecord = -1
	var records []Account
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		records = append(records, Account{
			Account_Type:   line[0],
			Account_Number: line[1],
			Trans_Date:     line[2],
			Cheque_no:      line[3],
			Description_1:  line[4],
			Description_2:  line[5],
			CAD:            line[6],
			USD:            line[7],
		})

	}

	for i := 1; i < len(records); i++ {
		_, err := db.Exec("insert into Bill(Account_Type,Account_Number,Trans_Date,Cheque_no,Description_1,Description_2,CAD,USD) values(?,?,?,?,?,?,?,?)", records[i].Account_Type, records[i].Account_Number, records[i].Trans_Date, records[i].Cheque_no, records[i].Description_1, records[i].Description_2, records[i].CAD, records[i].USD)
		if err != nil {
			fmt.Print(" " + err.Error())
		}
	}
	fmt.Println("all record inserted")

}
func main() {
	CSVReadWriter()
}

// things to do:
//1. change name of package and function of the package
//2. parameterize name of the input and output file
