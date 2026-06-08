package stages

import (
	"fmt"
	"time"

	"github.com/byteforge-run/tester-utils/runner"
	"github.com/byteforge-run/tester-utils/structured_output"
	"github.com/byteforge-run/tester-utils/test_case_harness"
	"github.com/byteforge-run/tester-utils/tester_definition"
)

func threeTierDiscountTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "three-tier-discount",
		Timeout:  10 * time.Second,
		TestFunc: testThreeTierDiscount,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/three_tier_discount.py",
			"tests/act1_opening_day/test_three_tier_discount.py",
		),
	}
}

func testThreeTierDiscount(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/three_tier_discount.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output_high", "220", "输入 250 → 输出「220」（高档 -30）"},
		{"output_mid", "110", "输入 120 → 输出「110」（中档 -10）"},
		{"output_low", "80", "输入 80 → 输出「80」（不到 100，原价）"},
		{"uses_elif", "true", "使用了 elif 关键字做多分支判断"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
