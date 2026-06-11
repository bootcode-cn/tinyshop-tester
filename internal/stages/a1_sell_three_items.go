package stages

import (
	"fmt"
	"time"

	"github.com/bootcode-cn/tester-utils/runner"
	"github.com/bootcode-cn/tester-utils/structured_output"
	"github.com/bootcode-cn/tester-utils/test_case_harness"
	"github.com/bootcode-cn/tester-utils/tester_definition"
)

func sellThreeItemsTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "sell-three-items",
		Timeout:  10 * time.Second,
		TestFunc: testSellThreeItems,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/sell_three_items.py",
			"tests/act1_opening_day/test_sell_three_items.py",
		),
	}
}

func testSellThreeItems(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/sell_three_items.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "18", "卖 3 件 × 6 元 → 输出「18」"},
		{"uses_for_range", "true", "使用了 for ... in range(...) 循环"},
		{"uses_price_var", "true", "使用了 price 变量保存单价"},
		{"total_is_computed", "true", "total 由循环累加得出（不能直接写死 18）"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
