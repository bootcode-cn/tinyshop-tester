package stages

import (
	"github.com/byteforge-run/tester-utils/tester_definition"
)

// GetDefinition returns the TesterDefinition for the tinyshop course.
func GetDefinition() tester_definition.TesterDefinition {
	return tester_definition.TesterDefinition{
		TestCases: []tester_definition.TestCase{
			// Act 1: 开张第一天
			helloShopTestCase(),
			firstProductTestCase(),
			setPriceTestCase(),
			firstSaleTestCase(),
			discountOrNotTestCase(),
			threeTierDiscountTestCase(),
			sellThreeItemsTestCase(),
			sellUntilQuitTestCase(),
			reusableDiscountTestCase(),
			firstBugTestCase(),

			// Act 2: 进货与库存
			productListTestCase(),
			browseMenuTestCase(),

			// Act 3: 第一个顾客（待添加）
			// Act 4: 开始赚钱（待添加）
			// Act 5: 数据驱动决策（待添加）
		},
	}
}

// pythonRule creates a LanguageRule for Python auto-detection.
// detectFile is the source file used for detection (e.g. "work/act1_opening_day/hello_shop.py").
// testDriver is the test script path (e.g. "tests/act1_opening_day/test_hello_shop.py").
func pythonRule(detectFile, testDriver string) tester_definition.LanguageRule {
	return tester_definition.LanguageRule{
		DetectFile: detectFile,
		Language:   "python",
		Source:     detectFile,
		RunCmd:     "python3",
		RunArgs:    []string{testDriver},
	}
}

// autoCompileStep returns a Python-only CompileStep for a tinyshop stage.
func autoCompileStep(detectFile, testDriver string) *tester_definition.CompileStep {
	return &tester_definition.CompileStep{
		Language: "auto",
		AutoDetect: []tester_definition.LanguageRule{
			pythonRule(detectFile, testDriver),
		},
	}
}
