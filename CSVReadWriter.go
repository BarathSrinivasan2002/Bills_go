package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"strconv"

	_ "github.com/btnguyen2k/consu/checksum"
	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	Account_Type   string  `json:"account_type"`
	Account_Number string  `json:"account_number"`
	Trans_Date     string  `json:"trans_date"`
	Cheque_no      string  `json:"cheque_no"`
	Description_1  string  `json:"description_1"`
	Description_2  string  `json:"description_2"`
	CAD            float64 `json:"CAD"`
	USD            float64 `json:"USD"`
	checksum       uint32  `json:"Checksum"`
}

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Barath2002@tcp(127.0.0.1:3306)/tax_schema?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db

}

func CSVReadWriter(csvFileName string) {
	var db = getMySQLDB()
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
		cad, err := strconv.ParseFloat(line[6], 64)
		if err != nil {
			fmt.Printf("Error converting string: %v", err)
		}
		usd, err := strconv.ParseFloat("0", 64)
		if err != nil {
			fmt.Printf("Error converting string: %v", err)
		}
		records = append(records, Account{
			Account_Type:   line[0],
			Account_Number: line[1],
			Trans_Date:     line[2],
			Cheque_no:      line[3],
			Description_1:  line[4],
			Description_2:  line[5],
			CAD:            cad,
			USD:            usd,
			checksum:       checksums(line[0] + line[1] + line[2] + line[3] + line[4] + line[5]),
		})

		//checksum:= checksum.Checksum(myval)
		//fmt.Printf(string(checksum))

	}

	for i := 1; i < len(records); i++ {
		_, err := db.Exec("insert into Bill(Account_Type,Account_Number,Trans_Date,Cheque_no,Description_1,Description_2,CAD,USD,checksum) values(?,?,?,?,?,?,?,?,?)",
			records[i].Account_Type, records[i].Account_Number, records[i].Trans_Date, records[i].Cheque_no, records[i].Description_1, records[i].Description_2, records[i].CAD, records[i].USD, records[i].checksum)
		if err != nil {
			fmt.Print(" " + err.Error())
		}
	}
	fmt.Println("all record inserted")

}
func checksums(record1 string) uint32 {
	// The string to generate a checksum for
	data := record1

	// Create a new CRC32 hash
	h := crc32.NewIEEE()

	// Write the data to the hash
	h.Write([]byte(data))

	// Calculate the checksum
	checksum := h.Sum32()

	// Print the checksum
	fmt.Println(checksum)
	return checksum
}
func main() {
	CSVReadWriter("csv44134.csv")
}

// things to do:
//1. change name of package and function of the package
//2. parameterize name of the input and output file
//3. change table name to txns
//4. add column transaction type
//5. derive txn type based on CAD if cad >= 0 then debit else credit
//6. derive check sum and check for duplicates
//7. Store date as date value
//8. Insert and update programs
//9. Path of the input file as parameters.
