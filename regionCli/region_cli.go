package region_cli

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type Branch struct {
	ID        int
	Region    string
	Name      string
	Address   string
	Phone     string
	CreatedAt string
	EmpCount  string
}

func RegionCmd() {
	// ./cli --command=list --region=tehran --> list of branches in tehran
	// ./cli --command=get   --> enter branch id --> get branch info
	// ./cli --command=create --region=tehran --> enter name,address,phone,created_at,empCount
	// ./cli --command=edit --region=tehran --> enter branch id
	// ./cli --command=status --region=tehran --> count of branches in region and count of employeas
	command := flag.String("command", "list", "its can be:list,get,create,edit,status")
	region := flag.String("region", "tehran", " your region ")
	flag.Parse()

	for {
		runCommand(*command, *region)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("plz enter new command :")
		scanner.Scan()
		*command = scanner.Text()

		fmt.Println("plz enter new region :")
		scanner.Scan()
		*region = scanner.Text()

	}
}

func runCommand(command string, region string) {
	switch command {
	case "list":
		listOfBranches(region)
	case "get":
		getBranchInfo()
	case "create":
		createBranch(region)
	case "edit":
		editBranch()
	case "status":
		getStatus(region)
	case "exit":
		os.Exit(0)
	default:
		os.Exit(1)

	}
}
func createBranch(region string) {
	fmt.Printf("create new branch for %s\n", region)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("plz enter name :")
	scanner.Scan()
	name := scanner.Text()

	fmt.Println("plz enter address :")
	scanner.Scan()
	address := scanner.Text()

	fmt.Println("plz enter date :")
	scanner.Scan()
	createdat := scanner.Text()

	fmt.Println("plz enter phone :")
	scanner.Scan()
	phone := scanner.Text()

	fmt.Println("plz enter empCount :")
	scanner.Scan()
	empCount := scanner.Text()

	branch := Branch{
		rand.Int(),
		region,
		name,
		address,
		phone,
		createdat,
		empCount,
	}
	writeToCsv(branch)
	fmt.Println("the new branch created id is:", branch.ID)
}

func listOfBranches(region string) {

	branches := readFromCsv()
	var branchesForRegion []Branch
	for _, branch := range branches {
		if branch.Region == region {
			branchesForRegion = append(branchesForRegion, branch)
		}
	}
	fmt.Printf("\nthe  branches of %s are:%+v\n", region, branchesForRegion)
}

func getBranchInfo() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("plz enter branch id :")
	scanner.Scan()
	id := scanner.Text()
	branches := readFromCsv()
	ID, _ := strconv.Atoi(id)
	for _, b := range branches {
		if b.ID == ID {
			fmt.Printf("branch:%+v\n", b)
			return
		}
	}
	fmt.Printf("dont have this branch\n")
}

func editBranch() {
	branches := readFromCsv()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("plz enter branch id :")
	scanner.Scan()
	id := scanner.Text()
	ID, _ := strconv.Atoi(id)

	fmt.Println("plz enter name :")
	scanner.Scan()
	name := scanner.Text()

	fmt.Println("plz enter address :")
	scanner.Scan()
	address := scanner.Text()

	fmt.Println("plz enter region :")
	scanner.Scan()
	region := scanner.Text()

	fmt.Println("plz enter phone :")
	scanner.Scan()
	phone := scanner.Text()

	fmt.Println("plz enter empCount :")
	scanner.Scan()
	empCount := scanner.Text()

	clearCsvFile()
	for _, branch := range branches {
		if branch.ID == ID {
			if name != "" {
				branch.Name = name
			}
			if address != "" {
				branch.Address = address
			}
			if region != "" {
				branch.Region = region
			}
			if phone != "" {
				branch.Phone = phone
			}
			if empCount != "" {
				branch.EmpCount = empCount
			}
		}
		writeToCsv(branch)
	}
}

func getStatus(region string) {
	branches := readFromCsv()
	empCount := 0
	branchCount := 0
	for _, branch := range branches {
		if branch.Region == region {
			branchCount++
			bempCount, _ := strconv.Atoi(branch.EmpCount)
			empCount += bempCount
		}
	}
	fmt.Printf("for %s the empCount is %d and branchCount is %d\n", region, empCount, branchCount)
}
func writeToCsv(branch Branch) {
	// Check if the CSV file already exists
	fileName := "output.csv"
	fileExists := fileExists(fileName)

	// Open the CSV file in append mode or create a new one if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// If the file is new, write the header row
	if !fileExists {
		header := []string{"ID", "Region", "Branch Name", "Address", "Phone", "Created At", "Employee Count"}
		writer.Write(header)
	}

	// Write data from the branches slice to the CSV file
	record := []string{strconv.Itoa(branch.ID), branch.Region, branch.Name, branch.Address, branch.Phone, branch.CreatedAt, branch.EmpCount}
	writer.Write(record)

	// Flush any buffered data to the file
	writer.Flush()

	// Check for writing errors
	if err = writer.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data has been written to", fileName)
}

func readFromCsv() []Branch {
	// Open the CSV file for reading
	fileName := "output.csv"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read and process each row from the CSV file
	var branches []Branch
	for {
		record, err := reader.Read()
		if err != nil {
			// End of the file
			break
		}

		// Parse the CSV record into a Branch struct
		id, _ := strconv.Atoi(record[0])
		branch := Branch{
			ID:        id,
			Region:    record[1],
			Name:      record[2],
			Address:   record[3],
			Phone:     record[4],
			CreatedAt: record[5],
			EmpCount:  record[6],
		}

		branches = append(branches, branch)

	}
	fmt.Println("read all records from csv")
	return branches
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func clearCsvFile() {
	fileName := "output.csv"

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Content of %s has been truncated\n", fileName)
}
