package stages

import (
	"fmt"
	"time"

	"github.com/bootcode-cn/tester-utils/runner"
	"github.com/bootcode-cn/tester-utils/structured_output"
	"github.com/bootcode-cn/tester-utils/test_case_harness"
	"github.com/bootcode-cn/tester-utils/tester_definition"
)

func firstSaleTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "first-sale",
		Timeout:  10 * time.Second,
		TestFunc: testFirstSale,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/first_sale.py",
			"tests/act1_opening_day/test_first_sale.py",
		),
	}
}

func testFirstSale(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir
	lang := harness.DetectedLang

	r := runner.Run(workDir, lang.RunCmd, lang.RunArgs...).
		WithTimeout(10 * time.Second).
		WithLogger(logger).
		Execute().
		Exit(0)

	if err := r.Error(); err != nil {
		return fmt.Errorf("test driver failed: %v", err)
	}

	results := structured_output.Parse(string(r.Result().Stdout))

	type tc struct {
		name     string
		expected string
		label    string
	}

	tests := []tc{
		{"file_exists", "true", "work/act1_opening_day/first_sale.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "15", "输入 3 杯 → 输出「15」"},
		{"uses_input", "true", "使用了 input() 读取顾客输入"},
		{"uses_int_conversion", "true", "使用了 int(input(...)) 把字符串转成整数"},
		{"uses_count_var", "true", "使用了 count 变量保存杯数"},
		{"total_is_computed", "true", "total 由 price * count 算出（不能直接写死 15）"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
