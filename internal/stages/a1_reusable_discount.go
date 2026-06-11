package stages

import (
	"fmt"
	"time"

	"github.com/bootcode-cn/tester-utils/runner"
	"github.com/bootcode-cn/tester-utils/structured_output"
	"github.com/bootcode-cn/tester-utils/test_case_harness"
	"github.com/bootcode-cn/tester-utils/tester_definition"
)

func reusableDiscountTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "reusable-discount",
		Timeout:  10 * time.Second,
		TestFunc: testReusableDiscount,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/reusable_discount.py",
			"tests/act1_opening_day/test_reusable_discount.py",
		),
	}
}

func testReusableDiscount(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/reusable_discount.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "110|80|90", "discount(120)/discount(80)/discount(100) 分别输出 110 / 80 / 90"},
		{"uses_def_discount", "true", "用 def 定义了函数 discount(...)"},
		{"uses_return", "true", "函数体内使用了 return 返回结果"},
		{"calls_discount", "true", "在函数外调用了 discount(...)"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
