# DKG E2E 脚本目录

DKG 的**自动化集成测试**已迁移到 Go 标准测试框架，请使用：

**`go test -tags=integration -v ./tests/integration/dkg/`**

详见 [tests/integration/dkg/README.md](../../tests/integration/dkg/README.md)。  
环境变量：`STORY_GRPC_URL`（必填）、`DKG_POLL_INTERVAL`、`DKG_MAX_WAIT`（可选）。

本目录保留供后续放置与 DKG E2E 相关的辅助脚本（如生成报告、清理数据等），如需可在此新增。
