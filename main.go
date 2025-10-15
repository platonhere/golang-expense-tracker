package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Expense struct {
	ID          int
	Date        string
	Description string
	Amount      float64
}

type ExpenseList struct {
	items  []Expense
	nextID int
}

func (e *ExpenseList) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(e.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (e *ExpenseList) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &e.items)
	if err != nil {
		return err
	}

	maxID := 0
	for _, exp := range e.items {
		if exp.ID > maxID {
			maxID = exp.ID
		}

	}
	e.nextID = maxID + 1
	return nil
}

// id e.nextID в параметрах некорректно, потому что e ещё не существует — это приёмник, который создаётся только при вызове метода.
// поэтому id должен быть только внутри метода, а не передаваться в него.
func (e *ExpenseList) NewExpense(date string, description string, amount float64) {
	e.items = append(e.items, Expense{
		ID:          e.nextID,
		Date:        date,
		Description: description,
		Amount:      amount,
	})
	fmt.Printf("Expense added successfully (ID: %d)\n", e.nextID)
	e.nextID++
}

func (e *ExpenseList) DeleteExpense(id int) {
	for i, exp := range e.items {
		if id == exp.ID {
			// Троеточие говорит Go: «распакуй все элементы этого слайса и добавь их по отдельности».
			e.items = append(e.items[:i], e.items[i+1:]...)
		}
	}
	fmt.Println("Expense deleted successfully")

}

func (e *ExpenseList) UpdateExpense(id int, data, discription string, amount float64) {
	for i, exp := range e.items {
		if exp.ID == id {
			if data != "" {
				e.items[i].Date = data
			}
			if discription != "" {
				e.items[i].Description = discription
			}
			if amount != 0 {
				e.items[i].Amount = amount
			}
			fmt.Println("Expense updated successfully")
			return
		}
	}
	fmt.Println("Expense not found")
}

func (e *ExpenseList) List() {
	fmt.Printf("%-4s %-10s %-12s %-6s\n", "ID", "Date", "Description", "Amount")
	for _, e := range e.items {
		fmt.Printf("%-4d %-10s %-12s %-6.2f\n", e.ID, e.Date, e.Description, e.Amount)
	}
}

func GetMouth(data string) int {
	m, err := strconv.Atoi(data[5:7])
	if err != nil {
		return 0
	}
	return m
}

func (e *ExpenseList) Sum(mouth ...int) {
	total := 0.0
	for _, exp := range e.items {
		expMouth := GetMouth(exp.Date)
		if len(mouth) == 0 || expMouth == mouth[0] {
			total += exp.Amount
		}
	}
	if len(mouth) == 0 {
		fmt.Printf("Total expenses: $%.2f\n", total)
	} else {
		mouthNames := map[int]string{
			1: "January", 2: "February", 3: "March", 4: "April",
			5: "May", 6: "June", 7: "July", 8: "August",
			9: "September", 10: "October", 11: "November", 12: "December",
		}
		name := mouthNames[mouth[0]]
		fmt.Printf("Total expenses for %s: $%.2f\n", name, total)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: expense-tracker[add|delete|list|summary|update]")
		return
	}

	command := os.Args[1]
	expenses := &ExpenseList{}
	expenses.LoadFromFile("expenses.json")

	switch command {
	case "add":
		date := flag.String("date", "2025-01-01", "Data in YYYY-MM-DD format")
		desc := flag.String("description", "", "Expense description")
		amount := flag.Float64("amount", 0, "Expense amount")
		flag.CommandLine.Parse(os.Args[2:])
		if *desc == "" || *amount <= 0 {
			fmt.Println("Error: negative amount or empty description")
		}
		expenses.NewExpense(*date, *desc, *amount)
		expenses.SaveToFile("expenses.json")

	case "delete":
		id := flag.Int("id", -1, "Expense ID to delete")
		flag.CommandLine.Parse(os.Args[2:])
		if *id < 0 {
			fmt.Println("Error: id required")
			return
		}
		expenses.DeleteExpense(*id)
		expenses.SaveToFile("expenses.json")

	case "update":
		id := flag.Int("id", -1, "Expense ID to update")
		date := flag.String("date", "", "New date")
		desc := flag.String("description", "", "New description")
		amount := flag.Float64("amount", 0, "Newamount")
		flag.CommandLine.Parse(os.Args[2:])
		if *id < 0 {
			fmt.Println("Error: id required")
			return
		}
		expenses.UpdateExpense(*id, *date, *desc, *amount)
		expenses.SaveToFile("expenses.json")

	case "summary":
		mouth := flag.Int("mouth", 0, "Select a month")
		flag.CommandLine.Parse(os.Args[2:])
		if *mouth > 0 {
			expenses.Sum(*mouth)
		} else {
			expenses.Sum()
		}

	case "list":
		expenses.List()

	default:
		fmt.Println("unknown command:", command)

	}
}
