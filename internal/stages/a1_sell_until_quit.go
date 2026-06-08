package stages

import (
	"fmt"
	"time"

	"github.com/byteforge-run/tester-utils/runner"
	"github.com/byteforge-run/tester-utils/structured_output"
	"github.com/byteforge-run/tester-utils/test_case_harness"
	"github.com/byteforge-run/tester-utils/tester_definition"
)

func sellUntilQuitTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "sell-until-quit",
		Timeout:  10 * time.Second,
		TestFunc: testSellUntilQuit,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/sell_until_quit.py",
			"tests/act1_opening_day/test_sell_until_quit.py",
		),
	}
}

func testSellUntilQuit(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/sell_until_quit.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output_sum", "15", "输入 5/3/7/done → 输出「15」"},
		{"output_empty", "0", "输入仅 done → 输出「0」（什么都没买）"},
		{"uses_while", "true", "使用了 while 关键字"},
		{"uses_break", "true", "使用了 break 跳出循环"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
