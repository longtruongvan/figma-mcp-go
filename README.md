# figma-mcp-go

Figma MCP — Free, No Rate Limits

Open-source Figma MCP server with full read/write access via plugin — no REST API, no rate limits. Turn text into designs and designs into real code. Works with Cursor, Claude, GitHub Copilot, and any MCP-compatible AI tool.

**Highlights**
- No Figma API token required
- No rate limits — free plan friendly
- **Read and Write** live Figma data via plugin bridge
- Supports multiple AI tools simultaneously
- Written in Go, distributed via npm

https://github.com/user-attachments/assets/17bda971-0e83-4f18-8758-8ac2b8dcba62

---

## Why this exists

Most Figma MCP servers rely on the **Figma REST API**.

That sounds fine… until you hit this:

| Plan | Limit |
|------|-------|
| Starter / View / Collab | **6 tool calls/month** |
| Pro / Org (Dev seat) | 200 tool calls/day |
| Enterprise | 600 tool calls/day |

If you're experimenting with AI tools, you'll burn through that in minutes.

I didn't have enough money to pay for higher limits.
So I built something that **doesn't use the API at all**.

---

## Installation & Setup

Install via `npx` — no build step required. Watch the setup video or follow the steps below.

[![Watch the video](https://img.youtube.com/vi/DjqyU0GKv9k/sddefault.jpg)](https://youtu.be/DjqyU0GKv9k)

### Điều kiện cần

- Figma Desktop
- Node.js 18+ cùng `npm`/`npx`
- Một AI client hỗ trợ MCP như Codex, Claude Code, Cursor, VS Code, hoặc GitHub Copilot
- Nếu muốn build plugin local từ repo này: cần chạy được `npm install` trong thư mục `plugin/`
- Nếu muốn build binary Go local từ source: cần cài `go`

### Bổ sung cho Codex local trên macOS

Nếu bạn muốn dùng chính source code trong repo này thay vì bản `npx` đã publish, hãy chuẩn bị như sau:

**Cài Go**
```bash
brew install go
go version
```

**Cài Codex**

- Cài ứng dụng Codex trên macOS.
- Sau khi cài xong, binary CLI thường nằm ở:
  - `/Applications/Codex.app/Contents/Resources/codex`

Nếu terminal báo `zsh: command not found: codex`, hãy export PATH:

```bash
echo 'export PATH="/Applications/Codex.app/Contents/Resources:$PATH"' >> ~/.zshrc
source ~/.zshrc
which codex
```

**Build binary local từ repo này**
```bash
go build -o bin/figma-mcp-go ./cmd/figma-mcp-go
```

**Refresh nhanh sau mỗi lần sửa code**

Repo có sẵn script [`dev-refresh.sh`](dev-refresh.sh) để build lại Go server và plugin chỉ với một lệnh:

```bash
chmod +x dev-refresh.sh
./dev-refresh.sh
```

Nếu shell của bạn không chạy trực tiếp file `.sh`, dùng:

```bash
bash dev-refresh.sh
```

Một vài mode hay dùng:

```bash
./dev-refresh.sh --skip-go
./dev-refresh.sh --skip-plugin
./dev-refresh.sh --restart-codex
```

Ý nghĩa:
- `--skip-go`: chỉ rebuild plugin Figma
- `--skip-plugin`: chỉ rebuild Go MCP server
- `--restart-codex`: build xong thì restart Codex app trên macOS

Sau khi chạy script:
1. Nếu phần Go server thay đổi, hãy mở thread Codex mới hoặc restart Codex.
2. Nếu phần plugin thay đổi, hãy đóng và mở lại plugin `Figma MCP Go` trong Figma.
3. Kiểm tra badge `Connected`.
4. Smoke test bằng `get_metadata`.

**Trỏ Codex sang binary local**
```bash
codex mcp remove figma-mcp-go
codex mcp add figma-mcp-go -- /Users/truongvanlong/Documents/GitHub/figma-mcp-go/bin/figma-mcp-go
```

**Kiểm tra lại MCP server trong Codex**
```bash
codex mcp list
codex mcp get figma-mcp-go
```

Sau đó hãy restart Codex hoặc mở thread mới để Codex nhận config mới.

### 1. Cấu hình AI tool

**Codex CLI**
```bash
codex mcp add figma-mcp-go -- npx -y @vkhanhqui/figma-mcp-go@latest
```

**Claude Code CLI**
```bash
claude mcp add -s project figma-mcp-go -- npx -y @vkhanhqui/figma-mcp-go@latest
```

**.mcp.json** (Claude and other MCP-compatible tools)
```json
{
  "mcpServers": {
    "figma-mcp-go": {
      "command": "npx",
      "args": ["-y", "@vkhanhqui/figma-mcp-go@latest"]
    }
  }
}
```

**.vscode/mcp.json** (Cursor / VS Code / GitHub Copilot)
```json
{
  "servers": {
    "figma-mcp-go": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@vkhanhqui/figma-mcp-go@latest"
      ]
    }
  }
}
```

### 2. Cài plugin Figma

Bạn có 2 cách:

**Cách A — dùng plugin từ release (nhanh nhất)**
1. Trong Figma Desktop: **Plugins → Development → Import plugin from manifest**
2. Chọn `manifest.json` từ [plugin.zip](https://github.com/vkhanhqui/figma-mcp-go/releases)
3. Chạy plugin trong bất kỳ file Figma nào

**Cách B — build plugin local từ repo này**
```bash
cd plugin
npm install
npm run build
```

Sau đó:
1. Trong Figma Desktop: **Plugins → Development → Import plugin from manifest**
2. Chọn [`plugin/manifest.json`](plugin/manifest.json)
3. Chạy plugin trong bất kỳ file Figma nào

Manifest local sẽ trỏ tới:
- `plugin/code.js`
- `plugin/index.html`

### 3. Mở AI client và kiểm tra kết nối

1. Restart hoặc mở lại AI client sau khi add MCP server.
2. Mở bất kỳ file Figma nào trong Figma Desktop.
3. Chạy plugin `Figma MCP Go`.
4. Đợi badge chuyển từ `Disconnected` sang `Connected`.
5. Test bằng một MCP call đơn giản như:

```text
get_metadata
```

Nếu bridge hoạt động đúng, bạn sẽ nhận được tên file hiện tại, tên trang, và ID các page.

### 4. Kiểm tra nhanh bằng local smoke test

Nếu muốn kiểm tra server bên ngoài AI client, chạy:

```bash
npx -y @vkhanhqui/figma-mcp-go@latest
```

Sau đó mở plugin Figma. Plugin sẽ kết nối tới `ws://localhost:1994/ws`.

Flow này đặc biệt hữu ích sau khi:
- chuyển sang máy mới
- cài lại macOS
- reset config của AI client
- debug tình trạng plugin bị `Disconnected`

### Xử lý sự cố

**Plugin vẫn hiện `Disconnected`**

- Hãy chắc chắn AI client của bạn đã thực sự khởi động MCP server. Chỉ thêm config thôi thì chưa đủ.
- Restart AI client hoặc mở session mới sau khi chạy `mcp add`.
- Nếu muốn kiểm tra thủ công, chạy `npx -y @vkhanhqui/figma-mcp-go@latest` trong terminal rồi mở lại plugin.
- Kiểm tra xem có gì đang chặn local port `1994` hay không.

**Dùng plugin local trong repo**

- Chạy `cd plugin && npm install && npm run build` trước khi import manifest.
- Build lại sau mỗi lần sửa file trong `plugin/src/`.
- Hoặc dùng `./dev-refresh.sh --skip-go` để rebuild plugin nhanh hơn.

**Lưu ý riêng cho Codex**

- Sau khi chạy `codex mcp add ...`, các thread cũ có thể chưa nhận MCP server mới ngay.
- Nếu Codex vẫn chưa thấy server, hãy mở thread mới hoặc restart app/CLI session.
- Nếu bạn vừa sửa cả Go server lẫn plugin local, chạy `./dev-refresh.sh --restart-codex` là cách an toàn nhất.

### 5. Dùng Codex để generate UI trực tiếp vào Figma

Làm theo đúng flow này nếu bạn muốn Codex tạo hoặc chỉnh UI trực tiếp trong Figma:

1. Hoàn thành các bước setup ở trên.
2. Mở Figma Desktop và mở file mà bạn muốn làm việc.
3. Chạy plugin `Figma MCP Go` trong file đó.
4. Đợi badge của plugin chuyển sang `Connected`.
5. Nếu bạn vừa mới add MCP server, hãy mở một thread Codex mới.
6. Yêu cầu Codex tạo màn hình trong Figma bằng ngôn ngữ tự nhiên.
7. Xem kết quả trực tiếp trên canvas Figma.
8. Tiếp tục refine bằng cách yêu cầu Codex chỉnh spacing, màu, hierarchy, nội dung, hoặc bố cục.

**Khuôn prompt nên dùng để ra UI đẹp hơn**

Thay vì chỉ nói `tạo login screen` hoặc `có hero, card, button`, hãy mô tả theo khuôn này:

```text
Dùng figma-mcp-go để tạo một màn hình trong Figma.

Loại màn: [mobile / desktop / dashboard / landing page]
Kích thước: [vd 390x844]
Tên frame: [tên màn]

Bối cảnh:
- Đây là màn hình cho [app/web gì]
- Người dùng là [ai]
- Mục tiêu của màn này là [đăng ký / học bài / mua hàng / xem thống kê]

Phong cách:
- [3-5 từ khóa rõ ràng: premium, editorial, modern, calm, playful, fintech...]
- Cảm giác mong muốn: [tin cậy / cao cấp / tối giản / năng động...]
- Nếu có, lấy vibe từ [app/thương hiệu], nhưng không copy y nguyên

Bảng màu:
- Màu chính: [...]
- Màu nền: [...]
- Accent: [...]
- Tránh: [...]

Bố cục:
- Phần quan trọng nhất là [...]
- Section cần có:
  - [...]
  - [...]
  - [...]
- Muốn layout theo hướng [nhiều khoảng thở / card lớn / asymmetrical / very clean]

Nội dung:
- Tiêu đề chính: [...]
- Subtitle: [...]
- CTA chính: [...]
- CTA phụ: [...]

Chi tiết thiết kế:
- Dùng auto-layout
- Spacing rộng và nhất quán
- Bo góc [lớn/vừa/nhỏ]
- Typography rõ hierarchy
- Đặt tên layer sạch

Điều cần tránh:
- Không quá boxy
- Không quá giống template
- Không wireframe look
- Không để card nào cũng giống nhau

Hãy tạo màn hình đủ polished để dùng làm product concept, không chỉ là mockup cơ bản.
```

**Nếu muốn nhờ ChatGPT viết hộ prompt hoàn chỉnh**

Bạn có thể ném đúng khuôn ở trên vào ChatGPT, hoặc dùng prompt meta này:

```text
Hãy viết cho tôi một prompt thật tốt để tôi dùng với figma-mcp-go trong Codex.

Yêu cầu:
- giữ cấu trúc rõ ràng
- mô tả đủ bối cảnh, style, bố cục, section, màu sắc, typography
- có thêm phần "Điều cần tránh"
- prompt đầu ra phải sẵn sàng để copy-paste trực tiếp vào Codex

Thông tin đầu vào của tôi:
- loại màn: [điền vào]
- kích thước: [điền vào]
- chủ đề sản phẩm: [điền vào]
- người dùng mục tiêu: [điền vào]
- phong cách mong muốn: [điền vào]
- bảng màu mong muốn: [điền vào]
- section bắt buộc: [điền vào]
- cảm giác cần tránh: [điền vào]
```

**Ví dụ prompt tốt**

```text
Dùng figma-mcp-go để tạo một màn hình mobile trong Figma.

Loại màn: mobile home screen
Kích thước: 390x844
Tên frame: Life in the UK Home

Bối cảnh:
- Đây là home screen cho app luyện thi quốc tịch Anh
- Người dùng cần cảm thấy bình tĩnh, tin tưởng, có động lực học tiếp

Phong cách:
- premium, editorial, modern, calm
- gợi tinh thần British heritage nhưng tinh tế
- không được nhìn như template giáo dục rẻ tiền

Bảng màu:
- navy sâu
- cream ấm
- red trầm
- gold nhạt
- tránh quá nhiều màu sáng chói

Bố cục:
- hero lớn ở nửa trên
- 1 progress card nổi bật
- 2 card phụ có treatment khác nhau: Mock Test và Study Guide
- bottom nav tối giản
- nhiều khoảng thở, bo góc lớn, hierarchy rõ

Nội dung:
- Title: Life in the UK
- Subtitle: Study with clarity. Pass with confidence.
- Progress: 27 lessons completed
- CTA: Start mock test

Chi tiết thiết kế:
- dùng auto-layout
- spacing rộng
- typography mạnh
- layer name rõ ràng

Điều cần tránh:
- không wireframe look
- không để card nào cũng giống nhau
- không quá phẳng
- không quá boxy
```

**Khi muốn sửa một frame có sẵn**

Chọn frame trong Figma trước, rồi dùng prompt kiểu:

```text
Dùng selection hiện tại trong Figma.

Mục tiêu:
- làm màn hình này premium hơn và hiện đại hơn
- tăng spacing
- cải thiện hierarchy
- giữ nguyên cấu trúc chính

Phong cách:
- [điền style mong muốn]
- [điền bảng màu mong muốn]

Điều cần tránh:
- không làm rối hơn
- không phá layout chính
- không biến thành wireframe hoặc template look
```

**Tóm tắt ngắn cho người mới**

Ít nhất hãy mô tả đủ 6 thứ sau:
- màn gì
- kích thước gì
- dành cho ai
- style gì
- section nào bắt buộc có
- điều gì cần tránh

**Nếu Codex không generate ra gì**

- Kiểm tra xem plugin có còn hiện `Connected` không.
- Mở một thread Codex mới.
- Restart Codex nếu session hiện tại được tạo trước khi chạy `codex mcp add ...`.
- Chạy local smoke test ở bước 4 để xác nhận bridge vẫn đang hoạt động.

---

## Advanced Tool Spec

- Xem chi tiet 5 tool moi phuc vu generate UI hien dai tai [docs/advanced-write-tools-spec.md](docs/advanced-write-tools-spec.md)

---

## Available Tools

### Write — Create

| Tool | Description |
|------|-------------|
| `create_frame` | Create a frame with optional auto-layout, fill, and parent |
| `create_rectangle` | Create a rectangle with optional fill and corner radius |
| `create_ellipse` | Create an ellipse or circle |
| `create_text` | Create a text node with optional font, alignment, spacing, and text-box sizing controls |
| `import_image` | Decode base64 image and place it as a rectangle fill |

### Write — Modify

| Tool | Description |
|------|-------------|
| `set_text` | Update text content of an existing TEXT node |
| `set_fills` | Set solid fill color (hex) on a node |
| `set_strokes` | Set solid stroke color and weight on a node |
| `set_layout_properties` | Update auto-layout, sizing, spacing, and layout-child properties |
| `set_text_style` | Update typography plus TEXT sizing controls such as width, height, and textAutoResize |
| `set_effects` | Replace node effects with shadows and/or blurs |
| `apply_styles` | Apply local fill, stroke, effect, and text styles by style ID |
| `create_instance` | Create an instance from a local component by ID, key, or name |
| `move_nodes` | Move nodes to an absolute x/y position |
| `resize_nodes` | Resize non-TEXT nodes by width and/or height |
| `rename_node` | Rename a node |
| `clone_node` | Clone a node, optionally repositioning or reparenting |

### Write — Delete

| Tool | Description |
|------|-------------|
| `delete_nodes` | Delete one or more nodes permanently |

### Document & Selection

| Tool | Description |
|------|-------------|
| `get_document` | Full current page tree |
| `get_metadata` | File name, pages, current page |
| `get_pages` | All pages (IDs + names) — lightweight, no tree loading |
| `get_selection` | Currently selected nodes |
| `get_node` | Single node by ID |
| `get_nodes_info` | Multiple nodes by ID |
| `get_design_context` | Depth-limited tree with `detail` level (`minimal`/`compact`/`full`) |
| `search_nodes` | Find nodes by name substring and/or type within a subtree |
| `scan_text_nodes` | All text nodes in a subtree |
| `scan_nodes_by_types` | Nodes matching given type list |
| `get_viewport` | Current viewport center, zoom, and visible bounds |

### Styles & Variables

| Tool | Description |
|------|-------------|
| `get_styles` | Paint, text, effect, and grid styles |
| `get_variable_defs` | Variable collections and values |
| `get_local_components` | All components + component sets with variant properties |
| `get_annotations` | Dev-mode annotations |
| `get_fonts` | All fonts used on the current page, sorted by frequency |
| `get_reactions` | Prototype/interaction reactions on a node |

### Export

| Tool | Description |
|------|-------------|
| `get_screenshot` | Base64 image export of any node |
| `save_screenshots` | Export images to disk (server-side, no API call) |

### MCP Prompts

| Prompt | Description |
|--------|-------------|
| `read_design_strategy` | Best practices for reading Figma designs |
| `design_strategy` | Best practices for creating and modifying designs |
| `text_replacement_strategy` | Chunked approach for replacing text across a design |
| `annotation_conversion_strategy` | Convert manual annotations to native Figma annotations |
| `swap_overrides_instances` | Transfer overrides between component instances |
| `reaction_to_connector_strategy` | Map prototype reactions into interaction flow diagrams |

---

## Related Projects

- [magic-spells/figma-mcp-bridge](https://github.com/magic-spells/figma-mcp-bridge)
- [grab/cursor-talk-to-figma-mcp](https://github.com/grab/cursor-talk-to-figma-mcp)
- [gethopp/figma-mcp-bridge](https://github.com/gethopp/figma-mcp-bridge)

---

## Contributing

Issues and PRs are welcome.
