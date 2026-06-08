package stages

import (
	"fmt"
	"time"

	"github.com/byteforge-run/tester-utils/runner"
	"github.com/byteforge-run/tester-utils/structured_output"
	"github.com/byteforge-run/tester-utils/test_case_harness"
	"github.com/byteforge-run/tester-utils/tester_definition"
)

func browseMenuTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "browse-menu",
		Timeout:  10 * time.Second,
		TestFunc: testBrowseMenu,
		CompileStep: autoCompileStep(
			"work/act2_stock_inventory/browse_menu.py",
			"tests/act2_stock_inventory/test_browse_menu.py",
		),
	}
}

func testBrowseMenu(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act2_stock_inventory/browse_menu.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "1. 可乐|2. 雪碧|3. 红茶|4. 咖啡|5. 矿泉水", "输出 5 行带 1./2./... 编号的饮料"},
		{"uses_list_literal", "true", "用 list 字面量 products = [...] 定义"},
		{"uses_enumerate", "true", "用 enumerate(...) 拿到下标（不允许 range(len(...)) 绕过）"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
