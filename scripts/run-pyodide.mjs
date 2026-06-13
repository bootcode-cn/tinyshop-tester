// scripts/run-pyodide.mjs — 用 Node + Pyodide 跑 tinyshop-python sample_test_*.py
//
// 与 dsa-core-python 完全对齐：
//   - SUBMISSION/tests/_sample_runner.py 写入 /sp/_sample_runner.py
//   - sample_test_*.py 写入 /sp/<phase>/sample_test_<slug>.py
//   - STUDENT_SOURCE 全局注入
//   - TEST:/RESULT: 解析协议
//
// 用法：
//   node scripts/run-pyodide.mjs                       # 所有已启用关
//   node scripts/run-pyodide.mjs hello-shop            # 指定关卡
//
// 环境变量：
//   SUBMISSION_DIR  默认 ../solution（需同时含 tests/_sample_runner.py、
//                   tests/<phase>/sample_test_*.py、work/<phase>/*.py）

import { loadPyodide } from "pyodide";
import { readFileSync, existsSync } from "node:fs";
import { resolve, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const TESTER_DIR = resolve(__dirname, "..");
const SUBMISSION = process.env.SUBMISSION_DIR
  ? resolve(process.env.SUBMISSION_DIR)
  : resolve(TESTER_DIR, "../solution");

// Stage 注册表 —— 与 internal/stages/stages.go 顺序一致
const STAGES = [
  // Act 1: 开张第一天
  { slug: "hello-shop",           phase: "act1_opening_day",      file: "hello_shop"           },
  { slug: "first-product",        phase: "act1_opening_day",      file: "first_product"        },
  { slug: "set-price",            phase: "act1_opening_day",      file: "set_price"            },
  { slug: "first-sale",           phase: "act1_opening_day",      file: "first_sale"           },
  { slug: "discount-or-not",      phase: "act1_opening_day",      file: "discount_or_not"      },
  { slug: "three-tier-discount",  phase: "act1_opening_day",      file: "three_tier_discount"  },
  { slug: "sell-three-items",     phase: "act1_opening_day",      file: "sell_three_items"     },
  { slug: "sell-until-quit",      phase: "act1_opening_day",      file: "sell_until_quit"      },
  { slug: "reusable-discount",    phase: "act1_opening_day",      file: "reusable_discount"    },
  { slug: "first-bug",            phase: "act1_opening_day",      file: "first_bug"            },
  { slug: "import-stdlib",        phase: "act1_opening_day",      file: "import_stdlib"        },
  // Act 2: 进货与库存
  { slug: "product-list",         phase: "act2_stock_inventory",  file: "product_list"         },
  { slug: "browse-menu",          phase: "act2_stock_inventory",  file: "browse_menu"          },
];

const filterSlug = process.argv[2];
const stages = filterSlug ? STAGES.filter((s) => s.slug === filterSlug) : STAGES;
if (stages.length === 0) {
  console.error(`No matching stage for slug: ${filterSlug}`);
  console.error(`Available: ${STAGES.map((s) => s.slug).join(", ")}`);
  process.exit(2);
}

console.log("==========================================");
console.log("  tinyshop-python Pyodide Tester (Node)");
console.log("==========================================\n");

// ── 启动 Pyodide ──
let stdoutBuf = [];
let stderrBuf = [];
const tBoot = Date.now();
const pyodide = await loadPyodide({
  stdout: (text) => stdoutBuf.push(text + "\n"),
  stderr: (text) => stderrBuf.push(text + "\n"),
});
const bootSecs = ((Date.now() - tBoot) / 1000).toFixed(2);
console.log(`Pyodide ${pyodide.version} loaded in ${bootSecs}s`);
console.log(`SUBMISSION_DIR = ${SUBMISSION}\n`);

// ── 把 _sample_runner.py 写入 /sp（一次，所有 stage 共用）──
try { pyodide.FS.mkdir("/sp"); } catch (_) {}
const sampleRunnerPath = resolve(SUBMISSION, "tests/_sample_runner.py");
if (!existsSync(sampleRunnerPath)) {
  console.error(`Fatal: missing ${sampleRunnerPath}`);
  process.exit(2);
}
pyodide.FS.writeFile("/sp/_sample_runner.py", readFileSync(sampleRunnerPath, "utf-8"));
await pyodide.runPythonAsync(
  "import sys\nif '/sp' not in sys.path: sys.path.insert(0, '/sp')"
);

// ── 逐关跑 ──
let passed = 0;
let failed = 0;

for (const stage of stages) {
  const entryPath = resolve(SUBMISSION, `tests/${stage.phase}/sample_test_${stage.file}.py`);
  const studentPath = resolve(SUBMISSION, `work/${stage.phase}/${stage.file}.py`);

  process.stdout.write(`🧪 [${stage.slug.padEnd(25)} pyodide] Testing... `);

  if (!existsSync(entryPath)) {
    console.log(`⚠️  SKIP (missing tests/${stage.phase}/sample_test_${stage.file}.py)`);
    failed++;
    continue;
  }
  if (!existsSync(studentPath)) {
    console.log(`⚠️  SKIP (missing work/${stage.phase}/${stage.file}.py)`);
    failed++;
    continue;
  }

  stdoutBuf = [];
  stderrBuf = [];
  const tStage = Date.now();

  // 清理上一关残留的模块缓存
  await pyodide.runPythonAsync(
    "import sys\nfor k in [m for m in list(sys.modules) if m.startswith('sample_test_') or m == '_sample_runner']:\n    del sys.modules[k]"
  );

  // 注入学员源码
  const studentCode = readFileSync(studentPath, "utf-8");
  pyodide.globals.set("STUDENT_SOURCE", studentCode);

  // 把 sample_test 写入 /sp/<phase>/sample_test_<file>.py 并设 __file__
  const entryCode = readFileSync(entryPath, "utf-8");
  const virtualDir = `/sp/${stage.phase}`;
  const virtualPath = `${virtualDir}/sample_test_${stage.file}.py`;
  try { pyodide.FS.mkdir(virtualDir); } catch (_) {}
  pyodide.FS.writeFile(virtualPath, entryCode);
  pyodide.globals.set("__file__", virtualPath);

  // 执行 sample_test 主体
  let runErr = null;
  try {
    await pyodide.runPythonAsync(entryCode);
  } catch (e) {
    runErr = e?.message || String(e);
    // 裁剪巨大 asm.js 堆栈跟踪
    if (runErr.length > 2000) runErr = runErr.slice(0, 2000) + "\n...(truncated)";
  } finally {
    pyodide.globals.delete("__file__");
  }

  const elapsed = ((Date.now() - tStage) / 1000).toFixed(2);
  const stdoutText = stdoutBuf.join("");
  const assertions = parseAssertions(stdoutText);
  const allPassed =
    !runErr && assertions.length > 0 && assertions.every((a) => a.passed);

  if (allPassed) {
    console.log(`✅ PASSED (${elapsed}s, ${assertions.length} assertions)`);
    passed++;
  } else {
    console.log(`❌ FAILED (${elapsed}s)`);
    if (runErr) {
      console.log(`    runtime error:`);
      for (const l of runErr.split("\n").slice(0, 15)) console.log(`      ${l}`);
    }
    for (const a of assertions) {
      if (!a.passed) {
        const exp = a.expected !== undefined ? `expected=${a.expected} ` : "";
        console.log(`    ✗ ${a.name}: ${exp}actual=${a.actual}`);
      }
    }
    const stderrText = stderrBuf.join("").trim();
    if (stderrText) {
      const lines = stderrText.split("\n").slice(0, 6);
      console.log(`    stderr:`);
      for (const l of lines) console.log(`      ${l}`);
    }
    failed++;
  }
}

console.log("\n==========================================");
console.log(`  Results: ${passed} passed, ${failed} failed`);
console.log("==========================================");
process.exit(failed > 0 ? 1 : 0);

// ── parseAssertions ──
function parseAssertions(stdout) {
  const lines = stdout.split("\n");
  const pairs = {};
  for (let i = 0; i < lines.length - 1; i++) {
    if (lines[i].startsWith("TEST:") && lines[i + 1].startsWith("RESULT:")) {
      pairs[lines[i].slice(5)] = lines[i + 1].slice(7);
      i++;
    }
  }
  const results = [];
  const seen = new Set();
  for (const key of Object.keys(pairs)) {
    if (key.endsWith(".actual")) {
      const base = key.slice(0, -7);
      if (seen.has(base)) continue;
      seen.add(base);
      const actual = pairs[key];
      const expected = pairs[`${base}.expected`];
      results.push({ name: base, actual, expected, passed: actual === expected });
    } else if (!key.endsWith(".expected") && !seen.has(key)) {
      seen.add(key);
      const v = pairs[key];
      results.push({ name: key, actual: v, expected: undefined, passed: v === "true" });
    }
  }
  return results;
}
