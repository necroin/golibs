package table_tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/necroin/golibs/libs/table"
)

func HeaderFormatter(s string) string {
	start := "\033[4m"
	end := "\033[0m"
	return fmt.Sprintf("%s%s%s", start, s, end)
}

func TestTable(t *testing.T) {
	table := table.New([]string{"ID", "Name", "Timestamp"}, table.WithHeaderFormatter(HeaderFormatter))

	for i := 0; i < 20; i++ {
		table.InsertRow(i, fmt.Sprintf("Name%d", i), time.Now())
	}

	table.Print()
}
