package stages

import (
	"fmt"
	"time"

	"github.com/bootcode-cn/tester-utils/runner"
	"github.com/bootcode-cn/tester-utils/structured_output"
	"github.com/bootcode-cn/tester-utils/test_case_harness"
	"github.com/bootcode-cn/tester-utils/tester_definition"
)

func firstBugTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "first-bug",
		Timeout:  10 * time.Second,
		TestFunc: testFirstBug,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/first_bug.py",
			"tests/act1_opening_day/test_first_bug.py",
		),
	}
}

func testFirstBug(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/first_bug.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "[debug] 5 * 3 = 15|15|[debug] 8 * 2 = 16|16", "subtotal(5,3) / subtotal(8,2) 输出 debug + 15 + debug + 16"},
		{"uses_def_subtotal", "true", "用 def 定义了函数 subtotal(price, qty)"},
		{"uses_return", "true", "函数体内使用了 return 返回结果"},
		{"uses_debug_print", "true", "埋了 [debug] 字面量的调试 print"},
		{"calls_subtotal", "true", "在函数外调用了 subtotal(...)"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
