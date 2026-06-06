#!/bin/bash
# 批量测试所有 stage 的 solution（Python）
# 用法: ./scripts/test-solution.sh [stage-slug]
#   不带参数：跑所有已启用的关卡
#   带参数：只跑指定关卡，例如 ./scripts/test-solution.sh hello-shop

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTER_DIR="$(dirname "$SCRIPT_DIR")"
SOLUTION_DIR="${TESTER_DIR}/../solution"

# 构建 tester
cd "$TESTER_DIR"
go build -o tinyshop-tester .

# Stage 列表（按课程顺序，slug 需与 stages.go 中注册的一致）
STAGES=(
	# Act 1: 开张第一天
	"hello-shop"          # S01
	"first-product"       # S02
	"set-price"           # S03
	"first-sale"          # S04
	"discount-or-not"     # S05
	# "three-tier-discount" # S06
	# "sell-three-items"  # S07
	# "sell-until-quit"   # S08
	# "reusable-discount" # S09
	# "first-bug"         # S10

	# Act 2-5: 待添加
)

# 若指定了参数，只跑该关卡
if [ -n "$1" ]; then
	STAGES=("$1")
fi

PASSED=0
FAILED=0

echo "=========================================="
echo "  TinyShop Solution Tester"
echo "=========================================="
echo ""

for stage in "${STAGES[@]}"; do
	printf "🧪 [%-20s python] Testing... " "$stage"

	start_time=$(python3 -c 'import time; print(time.time())')

	if ./tinyshop-tester -d="$SOLUTION_DIR" -s="$stage" > /dev/null 2>&1; then
		end_time=$(python3 -c 'import time; print(time.time())')
		elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
		echo "✅ PASSED (${elapsed}s)"
		((PASSED++))
	else
		end_time=$(python3 -c 'import time; print(time.time())')
		elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
		echo "❌ FAILED (${elapsed}s)"
		./tinyshop-tester -d="$SOLUTION_DIR" -s="$stage" 2>&1 | sed 's/^/    /'
		((FAILED++))
	fi
done

echo ""
echo "=========================================="
echo "  Results: $PASSED passed, $FAILED failed"
echo "=========================================="

if [ $FAILED -gt 0 ]; then
	exit 1
fi
