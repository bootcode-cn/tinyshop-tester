package stages

import (
	"fmt"
	"time"

	"github.com/bootcode-cn/tester-utils/runner"
	"github.com/bootcode-cn/tester-utils/structured_output"
	"github.com/bootcode-cn/tester-utils/test_case_harness"
	"github.com/bootcode-cn/tester-utils/tester_definition"
)

func importStdlibTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "import-stdlib",
		Timeout:  10 * time.Second,
		TestFunc: testImportStdlib,
		CompileStep: autoCompileStep(
			"work/act1_opening_day/import_stdlib.py",
			"tests/act1_opening_day/test_import_stdlib.py",
		),
	}
}

func testImportStdlib(harness *test_case_harness.TestCaseHarness) error {
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
		{"file_exists", "true", "work/act1_opening_day/import_stdlib.py 文件存在"},
		{"exit_ok", "true", "脚本正常退出（exit code 0）"},
		{"output", "4|169", "math.ceil(7/2)=4、floor(199*0.85)=169 两行输出正确"},
		{"uses_import_math", "true", "用 `import math` 引入了 math 模块"},
		{"uses_from_math_import", "true", "用 `from math import xxx` 单独借出函数"},
		{"uses_math_call", "true", "真的调用了 math.ceil / math.floor / ceil / floor"},
	}

	for _, t := range tests {
		if err := structured_output.AssertEqual(results, t.name, t.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", t.label)
	}

	return nil
}
