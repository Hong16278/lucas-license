# Lucas License

Lucas License 是面向独立软件作者的自托管激活码管理平台。本仓库基于 [Keygate](https://keygate.app) 修改，服务于 AutoMate 及后续桌面软件，销售发货与授权验证保持分离。

> 顾客在自动发卡平台付款并收到卡密；软件首次联网兑换卡密，随后使用服务端签名的本机许可证离线运行。

## 当前能力

- 按产品和套餐批量生成最多 5,000 张卡密，支持 TXT / CSV 导出。
- 月卡和年卡从第一次成功激活开始计时，而不是从制卡时间开始。
- 永久版签发 `exp=0` 的 Ed25519 令牌，首次激活后可永久离线验证。
- 默认一机一码；相同设备重复激活不重复占用名额，第二台设备会被拒绝。
- 卡密原文加密保存，查询用哈希；签名私钥仅存在于服务端。
- AutoMate 客户端本地验证产品、设备、到期时间和签名，并使用 AES-GCM 保存令牌。
- 管理登录使用邮箱 OTP，数据库定时备份，公网入口使用 Cloudflare Tunnel。

## 业务边界

```text
Lucas License 批量制卡
        ↓ TXT / CSV
自动发卡平台库存
        ↓ 顾客付款后自动发卡
AutoMate 首次联网激活
        ↓ Ed25519 签名令牌
本机离线使用
```

Lucas License 不处理淘宝订单、支付、聊天、退款或人工发货。永久离线设备无法实时收到服务端撤销状态，这是离线授权的必然取舍。

## 文档

- [产品规格](docs/PRODUCT-SPEC.zh-CN.md)
- [实施与验收进度](docs/IMPLEMENTATION-PLAN.zh-CN.md)
- [端到端验收报告](docs/VALIDATION-REPORT.zh-CN.md)
- [生产运维手册](docs/OPERATIONS.zh-CN.md)
- [OpenAPI](docs/openapi.yaml)

## 开发

后端需要 Go 1.25 和 PostgreSQL，前端位于 `web/`，使用 TypeScript、React 和 Vite。

```bash
go test ./...
cd web
npm ci
npm run build
```

生产部署必须自行生成 JWT、数据库、许可证签名和数据加密密钥；不要把 `.env`、卡密、SMTP 授权码或 Tunnel 凭据提交到仓库。

## 开源与归属

本项目是 Keygate 的衍生作品，使用不同的名称和 Logo，并保留界面中的 [Powered by Keygate](https://keygate.app) 归属。许可证和附加署名要求见 [LICENSE](LICENSE) 与 [NOTICE](NOTICE)。

本仓库按 GNU AGPL v3（含原项目 Section 7(b) 附加条款）公开。任何通过网络运行修改版的部署者，都应向该服务的用户提供相应源码。

