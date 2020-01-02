package crawler

import "testing"

func TestRunner(t *testing.T) {
	testDataSet := []string{
		"HOz-9FzIDf0",
		"8Z3TbMBfDM0",
		"ARt2fVT33Lw",
	}

	runner := NewRunner(3, testDataSet)

	runner.startRunner()
}
