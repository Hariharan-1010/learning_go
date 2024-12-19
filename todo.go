package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const filePath string = "data.csv"

func init_file() {
	filepath := filePath
	file, ok := os.Open(filepath)

	if ok != nil {
		fmt.Println("File reading error, ensure file exists. Creating new file...")

		file, err := os.Create(filepath)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		return
	}

	defer file.Close()
}

func print_all_data(data [][]string) {
	for sno, row := range data {
		if sno == 0 {
			fmt.Printf("sno, ")
		} else {
			fmt.Printf("%v, ", sno)
		}
		for _, col := range row {
			fmt.Printf("%s, ", col)
		}
		fmt.Println()
	}
}

func fetch_all_note() {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	print_all_data(data)
}

func write_item(item []string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "%v, %v\n", item[0], item[1])
}

func remove_finished_task(sno int) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	temp_file, ok := os.Create("temp_delete.csv")

	i := 0
	if err != nil || ok != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(temp_file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if i != sno {
			fmt.Fprintf(writer, "%v", line)
		}
		i = i + 1
	}
	writer.Flush()
	temp_file.Close()
	file.Close()
	if sno >= i {
		fmt.Println("Serial Number does not exist!")
	} else {
		err := os.Rename("temp_delete.csv", filePath)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	init_file()

	args_len := len(os.Args)

	switch args_len {
	case 1:
		fmt.Printf("Welcome to CLI notes.\nAvailable commands:\n\tTo query all notes: ./todo query-all\n\tTo add a note: ./todo add 'todo to add' 'place'\n\tTo delete a note: ./todo del 'serial number of the note to delete'\n")
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		switch os.Args[1] {
		case "query-all":
			fetch_all_note()
		case "add":
			if args_len < 4 {
				panic("Missing Arguments")
			}
			write_item([]string{os.Args[2], os.Args[3]})
		case "del":
			if args_len < 3 {
				panic("Missing Arguments")
			}
			sno, err := strconv.Atoi(os.Args[2])
			if err != nil {
				panic(err)
			}
			remove_finished_task(sno)
		default:
			panic("Bad Argument")
		}
	default:
		panic("Bad Argument")
	}
}
