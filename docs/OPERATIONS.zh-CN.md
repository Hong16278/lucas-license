# Lucas License 生产运维手册

## 1. 生产结构

```text
顾客 / AutoMate
  -> https://license.lucas.dpdns.org
  -> Cloudflare Edge
  -> cloudflared（服务器主动向外建立连接）
  -> 127.0.0.1:19000
  -> Lucas License / Keygate
  -> PostgreSQL
```

Cloudflare Tunnel 解决大陆云服务器对未备案域名的 HTTPS 阻断问题。隧道是出站连接，Keygate 仍只监听服务器本机端口。

## 2. 日常销售

1. 在管理后台选择产品、套餐、数量、批次名和渠道。
2. 生成批次并下载 TXT 或 CSV。
3. 把文件导入自动发卡平台库存。
4. 发卡平台负责付款和发货，Lucas License 只负责首次兑换、设备绑定和签名。

不需要逐单在授权后台操作。只有补卡、解绑或调查异常时才进入后台。

## 3. 套餐规则

| 套餐 | 起算时间 | 设备数 | 离线行为 |
|---|---|---:|---|
| 月卡 | 首次激活后 30 天 | 1 | 到期前可离线 |
| 年卡 | 首次激活后 365 天 | 1 | 到期前可离线 |
| 永久版 | 首次激活 | 1 | 当前大版本永久离线 |

永久离线令牌发出后，服务端无法主动让一台永不联网的电脑失效。后台撤销只会阻止后续联网激活或验证。

## 4. 健康检查

```bash
curl -fsS https://license.lucas.dpdns.org/health
systemctl is-active cloudflared
docker ps --filter name=keygate-prod
```

正常健康接口应返回 `status: ok` 且数据库检查为 `ok`。

## 5. 备份与恢复

- `keygate-backup.timer` 每天执行一次 PostgreSQL 自定义格式备份。
- 备份脚本会用 `pg_restore --list` 验证备份可读取，再原子改名。
- 默认保留约 14 天，目录位于 `/opt/keygate/backups/`。
- `.env.production`、Tunnel 凭据和数据库备份必须一同放入独立的加密离线备份。

恢复前先停止 Keygate 写入，创建新 PostgreSQL 数据卷，用 `pg_restore` 导入最近验证通过的 dump，再启动应用做健康检查。不要在未验证备份的情况下删除旧数据卷。

## 6. 发布与回滚

生产代码放在 `/opt/keygate/releases/<版本>/`，`/opt/keygate/current` 是当前版本软链接。新版本应先完整上传并检查文件，再原子切换软链接并重建 Keygate 容器。

若新版本健康检查失败，将软链接切回上一个 release，重新创建 Keygate 容器。数据库迁移前必须先完成备份；涉及不可逆迁移时要准备对应的回滚方案。

## 7. 密钥边界

- AutoMate 只内置产品公钥、产品 ID、产品 slug 和 API 地址。
- 许可证签名私钥、数据库密码、JWT 密钥、SMTP 授权码和 Cloudflare 凭据只保存在服务器。
- 不在 GitHub、网盘安装包、聊天记录或自动发卡平台中存放服务端密钥。
- 卡密可以导入发卡平台；卡密是兑换凭证，不是签名私钥。

## 8. Git 仓库分工

- `Hong16278/AutoMate`：私有，保存商业客户端源码，不作为顾客下载渠道。
- `Hong16278/lucas-license`：公开，满足 AGPL 对服务端对应源码的提供义务。
- 顾客安装包继续通过网盘发放；GitHub 仓库是否可访问不影响软件下载和激活。

