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

**Lưu ý riêng cho Codex**

- Sau khi chạy `codex mcp add ...`, các thread cũ có thể chưa nhận MCP server mới ngay.
- Nếu Codex vẫn chưa thấy server, hãy mở thread mới hoặc restart app/CLI session.

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

**Prompt mẫu — tạo một màn hình mobile mới**
```text
Dùng figma-mcp-go để tạo một màn hình mobile trong Figma.
Kích thước: 390x844.
Phong cách: premium, nhẹ nhàng, hiện đại.
Chủ đề: app luyện thi Life in the UK.
Bao gồm:
- hero section
- progress card
- mock test card
- study guide card
- bottom navigation
Dùng bảng màu navy, cream, red và gold.
Đặt tên layer rõ ràng.
```

**Prompt mẫu — chỉnh một frame có sẵn**
```text
Dùng selection hiện tại trong Figma.
Biến màn hình này thành một mobile UI sạch hơn và premium hơn.
Tăng spacing, cải thiện hierarchy, giữ nguyên cấu trúc chính,
và dùng bảng màu cream ấm cùng navy.
```

**Mẹo viết prompt**

- Luôn nói rõ loại màn hình: mobile, desktop, dashboard, landing page, v.v.
- Nếu biết kích thước thì nên ghi rõ, ví dụ `390x844` hoặc `1440x1024`.
- Nêu rõ phong cách hình ảnh: premium, editorial, fintech, education, playful, minimal.
- Liệt kê các section hoặc component mà bạn muốn xuất hiện trên màn hình.
- Nếu muốn Codex sửa một thiết kế có sẵn, hãy chọn frame đó trong Figma trước rồi nói rõ là `dùng selection hiện tại`.

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
| `create_text` | Create a text node (font loaded automatically) |
| `import_image` | Decode base64 image and place it as a rectangle fill |

### Write — Modify

| Tool | Description |
|------|-------------|
| `set_text` | Update text content of an existing TEXT node |
| `set_fills` | Set solid fill color (hex) on a node |
| `set_strokes` | Set solid stroke color and weight on a node |
| `set_layout_properties` | Update auto-layout, sizing, spacing, and layout-child properties |
| `set_text_style` | Update typography properties such as font, line height, spacing, and alignment |
| `set_effects` | Replace node effects with shadows and/or blurs |
| `apply_styles` | Apply local fill, stroke, effect, and text styles by style ID |
| `create_instance` | Create an instance from a local component by ID, key, or name |
| `move_nodes` | Move nodes to an absolute x/y position |
| `resize_nodes` | Resize nodes by width and/or height |
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
