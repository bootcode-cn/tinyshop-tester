# tinyshop-tester

ByteForge 平台用的 **tinyshop** 课程 Go 评测器，基于 [`tester-utils`](https://github.com/byteforge-run/tester-utils) 框架。

## 职责

- 监听用户提交，针对每个关卡运行对应的 Python `test_*.py` 测试驱动
- 解析驱动输出的 `TEST:`/`RESULT:` 结构化协议
- 校验每条断言并返回通过/失败

## 关卡组织

```
internal/stages/
├── stages.go            ← 注册所有 TestCase
├── a1_hello_shop.go     ← Act 1 第 1 关
├── a1_first_product.go  ← Act 1 第 2 关（待添加）
└── ...                  ← 命名约定: a{ActN}_{stage_slug}.go
```

每个关卡文件导出一个 `xxxTestCase()` 函数，在 `stages.go` 的 `GetDefinition()` 中注册。

## 开发

```bash
# 1. 构建
make build

# 2. 跑所有已启用关卡的 solution 验证
make test-solution

# 3. 只跑指定关卡
./scripts/test-solution.sh hello-shop

# 4. 手动跑（针对任意提交目录）
./tinyshop-tester -d=../solution -s=hello-shop
```

## 添加新关卡

1. 在 `internal/stages/aN_{slug}.go` 写 `{slug}TestCase()` + `test{Slug}()`
2. 在 `stages.go` 的 `GetDefinition()` 取消注释/添加到 `TestCases` 列表
3. 在 `scripts/test-solution.sh` 的 `STAGES` 数组取消注释
4. 在 `../solution/work/actN_xxx/` 写参考实现
5. 在 `../starter/tests/actN_xxx/test_{slug}.py` 写测试驱动（同步给 starter / solution）
6. `make test-solution` 验证

## 相关目录

- `../course/` — 课程内容（stage.yml / README.md / LEARNING.md）
- `../starter/` — 用户起始模板（work/ 空骨架 + tests/ 测试驱动）
- `../solution/` — 参考实现（starter + 通关代码）
