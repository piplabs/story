# DKG 造景脚本（Scenarios）

当你有 **validator 与 TEE 的控制权限** 时，可通过本目录下的场景脚本在测试前“造景”、测试后恢复，从而自动化运行 IT-E2E-03、IT-E2E-04、IT-REG-04、IT-DL-02/03、IT-FN-02、IT-ACT-02/03/04、IT-SKIP-01/02 等用例。

## 用法

1. 设置环境变量启用脚本驱动：
   ```bash
   export DKG_DRIVER=script
   export DKG_SCENARIO_SCRIPTS_DIR=/path/to/story/scripts/dkg_e2e/scenarios
   export STORY_GRPC_URL=localhost:9090
   ```
2. 按你的部署方式编辑各场景下的 `pre.sh` / `post.sh`（见下表）。
3. 运行集成测试；对绑定了场景的用例，测试会先执行 `pre.sh`，再跑断言，最后执行 `post.sh`。
   ```bash
   go test -tags=integration -v -run TestDKG_P1 ./tests/integration/dkg/
   ```

## 场景列表

| 目录名 | 说明 | 对应用例 |
|--------|------|----------|
| **insufficient_verified** | 只保留 1 个 validator 开 DKG（或 TEE），其余关闭 | IT-E2E-03, IT-DL-02, IT-SKIP-01 |
| **insufficient_finalized** | Dealing 结束后只让 1 个能 finalize | IT-E2E-04, IT-ACT-02, IT-ACT-03, IT-SKIP-02 |
| **dkg_disabled_one_node** | 关闭其中 1 个节点的 DKG（config + 重启） | IT-REG-04, IT-DL-03, IT-FN-02, IT-ACT-04 |
| **tee_down_one_node** | 停止其中 1 个节点的 story-kernel | IT-REG-10/11, IT-DL-08/10, IT-FN-07/08, IT-ACT-12, IT-RES-*, IT-EDGE-06 |
| **height_below_dkg_start** | 当前高度 < dkgStartBlock（通常需新链或快照） | IT-BB-01 |
| **upgrade_scheduled** | 调用 UpgradeScheduled(activationHeight, version) | IT-BB-03/09, IT-REG-03, IT-DL-09, IT-ACT-05, IT-E2E-06, IT-UPG-* |
| **resharing_second_round** | 第一轮 Active 结束后等周期结束触发新轮 | IT-REG-02, IT-BB-07, IT-E2E-02 |
| **complaint_justification** | TEE 或数据层面产生无效 deal 触发 complaint 路径 | IT-E2E-05 |

每个场景目录下应有可执行的 `pre.sh`（造景）和 `post.sh`（恢复）；若不存在则跳过，不报错。
