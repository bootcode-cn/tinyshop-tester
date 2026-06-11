package stages

import (
	"fmt"
	"time"

	"github.com/bootcode-cn/tester-utils/runner"
	"github.com/bootcode-cn/tester-utils/structured_output"
	"github.com/bootcode-cn/tester-utils/test_case_harness"
	"github.com/bootcode-cn/tester-utils/tester_definition"
)

func discountOrNotTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "discount-or-not",
		Timeout:  10 * time.Second,
		TestFunc: testDiscountOrNot,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/discount_or_not.py",
			"tests/act1_opening_day/test_discount_or_not.py",
		),
	}
}

func testDiscountOrNot(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/discount_or_not.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output_high", "110", "输入 120 → 输出「110」（满 100 减 10）"},
		{"output_low", "80", "输入 80 → 输出「80」（不到 100，原价）"},
		{"uses_if", "true", "使用了 if 关键字做分支判断"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
