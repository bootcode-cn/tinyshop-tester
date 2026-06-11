package stages

import (
	"fmt"
	"time"

	"github.com/bootcode-cn/tester-utils/runner"
	"github.com/bootcode-cn/tester-utils/structured_output"
	"github.com/bootcode-cn/tester-utils/test_case_harness"
	"github.com/bootcode-cn/tester-utils/tester_definition"
)

func firstProductTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:    "first-product",
		Timeout: 10 * time.Second,
		TestFunc: testFirstProduct,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/first_product.py",
			"tests/act1_opening_day/test_first_product.py",
		),
	}
}

func testFirstProduct(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/first_product.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "拿铁", "输出为「拿铁」"},
		{"uses_product_var", "true", "使用了 product 变量保存商品名（不能直接 print 字面量）"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
