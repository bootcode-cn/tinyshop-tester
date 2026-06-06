package stages

import (
	"fmt"
	"time"

	"github.com/byteforge-run/tester-utils/runner"
	"github.com/byteforge-run/tester-utils/structured_output"
	"github.com/byteforge-run/tester-utils/test_case_harness"
	"github.com/byteforge-run/tester-utils/tester_definition"
)

func helloShopTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:    "hello-shop",
		Timeout: 10 * time.Second,
		TestFunc: testHelloShop,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/hello_shop.py",
			"tests/act1_opening_day/test_hello_shop.py",
		),
	}
}

func testHelloShop(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/hello_shop.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "欢迎光临 TinyShop", "输出为「欢迎光临 TinyShop」"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
