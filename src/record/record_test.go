package record

import (
	"fmt"
	"testing"
)

func TestRecord(t *testing.T) {
	NewRecord("1", 10, 100, 0, 0)
	NewRecord("1", 2, 50, 10, 0)
	NewRecord("1", 1, 40, 10, 0)
	NewRecord("1", 10, 30, 10, 0)
	fmt.Println(GetPrinciple("1", 10, 100))
	NewRecord("1", -10, 100, 0, 0)
	fmt.Println(GetPrinciple("1", 3, 40))
	NewRecord("1", 10, 0, 0, 0)
	fmt.Println(GetPrinciple("1", 10, 0))
}
