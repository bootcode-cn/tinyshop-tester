#!/bin/bash
# scripts/test-pyodide.sh — Pyodide 本地集成测试（模拟前端 Playground）
#
# 用法:
#   ./scripts/test-pyodide.sh                       # 所有关
#   ./scripts/test-pyodide.sh hello-shop            # 指定关卡

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTER_DIR="$(dirname "$SCRIPT_DIR")"

cd "$TESTER_DIR"

# 首次运行装依赖（pyodide npm 包 ~50MB，含 WASM blob）
if [ ! -d node_modules/pyodide ]; then
  echo "📦 Installing pyodide@0.27.0 (first run, ~50MB)..."
  if command -v pnpm >/dev/null 2>&1; then
    pnpm install --silent
  else
    npm install --silent --no-audit --no-fund
  fi
fi

exec node scripts/run-pyodide.mjs "$@"
