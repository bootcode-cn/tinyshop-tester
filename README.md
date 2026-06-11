# TinyShop Tester (tinyshop-tester)

TinyShop 课程自动评测工具。

## 方式一：源码构建

```bash
git clone https://github.com/bootcode-cn/tinyshop-tester
cd tinyshop-tester
go build -o tinyshop-tester .
./tinyshop-tester -s hello-shop -d ~/my-solution
```

**依赖：** Go 1.24+

## 方式二：Docker 镜像

**快速开始**

```bash
cd ~/my-solution  # 你的解答目录（包含 bootcode.yml）
docker pull ghcr.io/bootcode-cn/tinyshop-tester:latest
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/bootcode-cn/tinyshop-tester:latest -s hello-shop -d /workspace
```

**推荐：创建 test.sh 脚本**

在解答目录下创建 `test.sh`：

```bash
#!/bin/bash
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/bootcode-cn/tinyshop-tester:latest \
  -s "${1:-hello-shop}" -d /workspace
```

用法：`chmod +x test.sh && ./test.sh discount-or-not`

**本地构建镜像（可选）**

```bash
git clone https://github.com/bootcode-cn/tinyshop-tester
cd tinyshop-tester
docker build -t my-tester .
docker run --rm --user $(id -u):$(id -g) -v ~/my-solution:/workspace my-tester -s hello-shop -d /workspace
```

## 关卡列表

| #   | Slug              | 关卡                   |
| --- | ----------------- | ---------------------- |
| S01 | `hello-shop`      | 开张第一天：第一句吆喝 |
| S02 | `first-product`   | 上架第一件商品         |
| S03 | `set-price`       | 给拿铁定价             |
| S04 | `first-sale`      | 第一个顾客来了         |
| S05 | `discount-or-not` | 满 100 减 10           |

## License

MIT
