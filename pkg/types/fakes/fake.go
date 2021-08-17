package fakes

import (
	"fmt"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
)

func init() {
	fake.Seed(time.Now().UnixNano())
}

const (
	exampleQuantity = 3
)

// BuildUniqueName builds a name that is hopefully more unique than just one word.
func BuildUniqueName() string {
	return fmt.Sprintf("%s%s%s%s", fake.Word(), fake.Word(), fake.Word(), fake.Word())
}

// BuildFakeSQLQuery builds a fake SQL query and arg pair.
func BuildFakeSQLQuery() (query string, args []interface{}) {
	s := fmt.Sprintf("%s %s WHERE things = ? AND stuff = ?",
		fake.RandomString([]string{"SELECT * FROM", "INSERT INTO", "UPDATE"}),
		fake.Password(true, true, true, false, false, 32),
	)

	return s, []interface{}{"things", "stuff"}
}

// BuildFakeID builds a fake ID.
func BuildFakeID() uint64 {
	return fake.Uint64()
}

// BuildFakeTime builds a fake time.
func BuildFakeTime() uint64 {
	return fake.Uint64()
}
