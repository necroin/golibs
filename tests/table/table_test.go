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
	tb := table.New("ID", "Name", "Timestamp")
	tb.SetOptions(table.WithHeaderFormatter(HeaderFormatter), table.WithPadding(4), table.WithPadchar('|'))
	for i := 0; i < 20; i++ {
		tb.InsertRow(i, fmt.Sprintf("Name%d", i), time.Now())
	}
	tb.InsertRow(20, fmt.Sprintf("Name%d", 20))
	tb.InsertRow(21, "", time.Now())

	tb.Print()
}
