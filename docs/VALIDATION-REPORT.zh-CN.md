# Keygate × AutoMate 验收报告

日期：2026-08-11  
范围：本地测试环境，不包含生产部署和真实发卡平台。

## 1. 测试环境

- Keygate：当前工作区源码构建
- Go：1.25.0（项目内便携工具链，仅用于开发验证）
- PostgreSQL：18-alpine，本地 Docker 容器
- AutoMate：`AutoMate.Licensing`，.NET 10
- 测试服务地址：`127.0.0.1:19000`，仅开发模式临时运行

测试过程中使用随机临时 JWT 密钥和许可证签名密钥；未把任何测试密钥写入源码或文档。Keygate 测试服务在验收结束后已停止，PostgreSQL 测试容器保留供后续开发使用。

## 2. 自动化构建与测试

| 项目 | 命令/方式 | 结果 |
|---|---|---|
| Go 服务端 | `go test ./...` | 全部通过 |
| Web 管理端 | `npm run typecheck` | 通过 |
| Web 管理端 | `npm run build` | 通过 |
| AutoMate | `dotnet build AutoMate.sln --no-restore` | 0 错误 |
| 客户端互通 | AutoMate `TokenValidator` 验证真实 Keygate Token | 通过 |

AutoMate 构建仍有两项既有警告：旧 KamiCloud 依赖的 SQLite 原生包安全公告，以及 Windows DPI manifest 提示。它们不是本次 Keygate 接入引入的错误，但正式发布前应单独处理 SQLite 依赖公告。

## 3. 数据库与首次安装

1. 在空 PostgreSQL 数据库执行 22 个迁移，全部成功。
2. 确认 `plans.duration_days`、`license_batches` 和 `licenses.batch_id` 已落库。
3. 首次安装向导实测时发现默认套餐未填写后来变为必填的 `checkout_id`，导致 HTTP 500。
4. 修复后重新构建，首次初始化成功，并创建 owner、AutoMate 桌面产品和默认套餐。

## 4. 批量制卡与导出

### 小批次功能测试

- 月卡批次生成 3 张，TXT 和 CSV 均可重新下载。
- CSV 包含 `license_key`、产品、套餐、`duration_days`、设备数、批次和渠道列。
- 永久版批次生成 2 张。
- 批次卡不绑定客户邮箱，符合“发卡平台只管理库存、Keygate 只管理授权”的边界。

### 1,000 张验收

- 接口一次生成：1,000 张
- 接口返回唯一数：1,000
- TXT 导出行数：1,000
- TXT 导出唯一数：1,000
- 本机实测耗时：约 1,569 ms

因此“单次至少生成并导出 1,000 张互不重复卡密”的验收标准已满足。

## 5. 激活规则

### 月卡

- 套餐 `duration_days=30`、`max_activations=1`。
- 卡密生成时 `valid_until` 为空，不提前消耗期限。
- 首次激活后 `valid_until - valid_from = 30 天`。
- Token 的 `exp - iat = 2,592,000 秒`，即完整 30 天。
- 同一设备再次激活返回 `already_activated`，不增加设备数。
- 第二设备激活返回 `ACTIVATION_LIMIT`。

### 永久版

- 套餐类型为 `perpetual`、`duration_days=0`。
- 首次激活签发的 Token `exp=0`，没有时间到期值。
- 永久离线意味着长期断网设备无法收到服务端撤销，这是产品设计中明确接受的取舍。

### 年卡

- 套餐 `duration_days=365`、`max_activations=1`。
- 首次激活后 Token 的 `exp - iat = 31,536,000 秒`，即完整 365 天。
- 关闭 Keygate 后，AutoMate 仍可本地验签并识别产品、套餐和设备绑定。

## 6. AutoMate 离线互通

取得真实 Keygate 月卡和永久 Token、公钥后，关闭 Keygate 服务，再调用 AutoMate 自己的 `TokenValidator`：

- 月卡和年卡 Token：验签成功，可离线使用至各自到期时间。
- 永久 Token：验签成功，识别为永久。
- 换成其他产品 ID：返回 `product-mismatch`。
- 换成其他设备指纹：返回 `device-mismatch`。
- 修改 Token 签名：被拒绝。

这证明 AutoMate 的 Keygate 两段式 Ed25519 Token 适配并非只在服务端自测通过，而是客户端真实可读、可验、可离线使用。

## 7. 尚未完成

- 正式服务器的 HTTPS、域名、持久化密钥、数据库备份和监控。
- 正式库中的月卡、年卡和永久版套餐配置（本地三档已完成端到端验收）。
- 将测试批次导入实际自动发卡平台沙箱。
- 小批量真实订单验证。
- 生产 AutoMate 仍保持 KamiCloud，不在以上项目完成前切换。

## 8. 上线门槛

正式部署必须持久化并异地备份以下内容：

1. PostgreSQL 数据库；
2. `LICENSE_SIGNING_KEY`；
3. `RELEASE_KEY_ENCRYPTION_KEY`；
4. `JWT_SECRET`；
5. AutoMate 发布包中对应的 Keygate 公钥和产品 ID。

丢失许可证签名私钥会导致新旧客户端信任链断裂，因此不能在每次部署时随机生成。测试环境使用随机密钥只是为了避免把秘密落盘，不可照搬到生产。
