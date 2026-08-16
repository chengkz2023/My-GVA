# 部署约定

客户项目的部署形态由实际项目决定（Q12 结论）；本文件记录脚手架强制一致的**部署约定**，形态可变的产物（compose/systemd/nginx 具体配置）不入脚手架。

## 同源代理（Same-origin proxy）

前端静态资源由 Nginx 托管，`/api` 反向代理到后端。浏览器始终同源，后端**不启用 CORS**（配置中已移除 cors 段）。

```nginx
server {
    listen 80;
    server_name example.com;

    root /usr/share/nginx/html;   # web/dist
    index index.html;

    location /api/ {
        proxy_pass http://127.0.0.1:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

- `web/nginx.conf` 是该约定的镜像构建用示例。
- 若未来确有跨域直连需求，必须按「扩展接缝」新增 CORS 中间件并恢复配置段；在此之前配置中没有 cors 段，属于有意删除。

## 参考形态（公司现状，随项目复刻）

- MySQL、Redis：Docker 容器。
- 后端：裸机 + systemd（环境变量注入 `ADMIN_INITIAL_PASSWORD` 等，见 `server/.env.example`）。
- 前端：Nginx 托管静态资源 + 上述反向代理。

## 迁移分工

系统表 AutoMigrate / 业务表 SQL，见 `docs/adr/0003-migration-split.md`。
