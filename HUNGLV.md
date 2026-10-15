# Fork vhung/vikunja — bản vá cho pm.hunglv.com

Fork của [go-vikunja/vikunja](https://github.com/go-vikunja/vikunja) (AGPL-3.0). Mã nguồn của bản đang chạy
trên pm.hunglv.com nằm ở nhánh `hunglv/<tag>` tương ứng với image `ghcr.io/vhung/vikunja:<tag>-hunglv`.

## Bản vá: giới hạn người tạo project gốc

- Cấu hình `service.projectcreators` (env `VIKUNJA_SERVICE_PROJECTCREATORS`): danh sách username, cách nhau
  bằng dấu phẩy hoặc khoảng trắng, so chính xác từng ký tự (kể cả hoa thường). **Rỗng = hành vi gốc** (ai cũng tạo được).
- User không có tên: không tạo project gốc qua `POST /projects`, nhân bản ra gốc, tách project con ra gốc,
  hay import/migration. Vẫn tạo project con dưới project mình có quyền ghi.
- Project "Inbox" khi tạo user (đăng ký, CLI, OpenID) vẫn được tạo.
- `GET /api/v2/user` có thêm `can_create_projects`; trang danh sách project ẩn nút "New project" khi false.
- Không đổi schema DB: quay về image gốc chỉ cần đổi tag image.

## Tự theo bản mới

Workflow `.github/workflows/hunglv-sync.yml` chạy mỗi ngày: thấy tag `vX.Y.Z` mới ở upstream thì cherry-pick
các commit từ tag gốc của `hunglv/main` lên tag mới, chạy test, build image rồi đẩy nhánh `hunglv/<tag>`.
Xung đột hoặc test đỏ thì dừng, không build. Không tự deploy.
