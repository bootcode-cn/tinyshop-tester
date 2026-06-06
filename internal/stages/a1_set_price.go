package stages

import (
	"fmt"
	"time"

	"github.com/byteforge-run/tester-utils/runner"
	"github.com/byteforge-run/tester-utils/structured_output"
	"github.com/byteforge-run/tester-utils/test_case_harness"
	"github.com/byteforge-run/tester-utils/tester_definition"
)

func setPriceTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:    "set-price",
		Timeout: 10 * time.Second,
		TestFunc: testSetPrice,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/set_price.py",
			"tests/act1_opening_day/test_set_price.py",
		),
	}
}

func testSetPrice(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/set_price.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "3.5", "输出为「3.5」"},
		{"uses_cost_var", "true", "使用了 cost 变量保存进货成本"},
		{"uses_markup_var", "true", "使用了 markup 变量保存加价"},
		{"price_is_computed", "true", "price 由 cost + markup 算出（不能直接写死 3.5）"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
