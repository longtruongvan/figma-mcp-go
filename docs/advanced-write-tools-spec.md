# Advanced Write Tools Spec

Muc tieu cua bo tool nay la nang tran chat luong UI khi AI ghi vao Figma.
Thay vi chi co primitive co ban nhu frame/text/rectangle, AI se co them
nhung tool de dieu khien layout, typography, effects, style references,
va component instances.

## Nguyen tac thiet ke

- Tap trung vao nhung tool tang chat luong UI ro nhat.
- Input giu don gian, de model goi duoc bang prompt tu nhien.
- Uu tien API Figma on dinh, khong dua vao hack hoac field kho bao tri.
- Tool moi phai bo sung duoc cho flow hien tai, khong pha backward compatibility.

## 1. `set_layout_properties`

Muc dich:
- Dieu khien auto-layout va sizing cua frame hoac child trong auto-layout.

Node ho tro:
- FRAME, COMPONENT, INSTANCE, TEXT va cac node ho tro auto-layout fields.

Input:
- `nodeId` bat buoc
- `layoutMode?`: `HORIZONTAL` | `VERTICAL` | `NONE`
- `layoutWrap?`: `NO_WRAP` | `WRAP`
- `primaryAxisAlignItems?`: `MIN` | `CENTER` | `MAX` | `SPACE_BETWEEN`
- `counterAxisAlignItems?`: `MIN` | `CENTER` | `MAX` | `BASELINE`
- `primaryAxisSizingMode?`: `FIXED` | `AUTO`
- `counterAxisSizingMode?`: `FIXED` | `AUTO`
- `layoutSizingHorizontal?`: `FIXED` | `HUG` | `FILL`
- `layoutSizingVertical?`: `FIXED` | `HUG` | `FILL`
- `layoutAlign?`: `MIN` | `CENTER` | `MAX` | `STRETCH` | `INHERIT`
- `layoutGrow?`: `0` | `1`
- `layoutPositioning?`: `AUTO` | `ABSOLUTE`
- `paddingTop?`, `paddingRight?`, `paddingBottom?`, `paddingLeft?`
- `itemSpacing?`
- `counterAxisSpacing?`
- `clipsContent?`

Output:
- `id`, `name`
- tap field da ap dung

## 2. `set_text_style`

Muc dich:
- Dieu khien typography ngoai text content.

Node ho tro:
- TEXT

Input:
- `nodeId` bat buoc
- `fontFamily?`
- `fontStyle?`
- `fontSize?`
- `textAutoResize?`: `NONE` | `WIDTH_AND_HEIGHT` | `HEIGHT` | `TRUNCATE`
- `width?`
- `height?`
- `textCase?`: `ORIGINAL` | `UPPER` | `LOWER` | `TITLE` | `SMALL_CAPS` | `SMALL_CAPS_FORCED`
- `textAlignHorizontal?`: `LEFT` | `CENTER` | `RIGHT` | `JUSTIFIED`
- `textAlignVertical?`: `TOP` | `CENTER` | `BOTTOM`
- `paragraphSpacing?`
- `lineHeight?`: `{ unit: "AUTO" }` hoac `{ unit: "PIXELS" | "PERCENT" | "FONT_SIZE_%", value: number }`
- `letterSpacing?`: `{ unit: "PIXELS" | "PERCENT", value: number }`

Output:
- `id`, `name`
- cac text style field da ap dung

Ghi chu:
- Tool phai load font truoc khi ghi style.
- Neu text node co mixed fonts, tool van phai co gang load tat ca font dang duoc dung.
- Khi can text wrapping on dinh, truyen `width` va de tool dat `textAutoResize = HEIGHT`.
- Khi can fixed text box, truyen ca `width` va `height` hoac set `textAutoResize = NONE`.

## 3. `set_effects`

Muc dich:
- Them depth cho UI thong qua shadow va blur.

Node ho tro:
- Cac node co `effects`

Input:
- `nodeId` bat buoc
- `effects` bat buoc, cho phep mang rong de clear effect

Effect object ho tro:
- Drop shadow:
  - `type: "DROP_SHADOW"`
  - `color`
  - `x?`, `y?`
  - `radius?`
  - `spread?`
  - `blendMode?`
  - `visible?`
  - `showShadowBehindNode?`
- Inner shadow:
  - `type: "INNER_SHADOW"`
  - `color`
  - `x?`, `y?`
  - `radius?`
  - `spread?`
  - `blendMode?`
  - `visible?`
- Blur:
  - `type: "LAYER_BLUR"` hoac `type: "BACKGROUND_BLUR"`
  - `radius`
  - `visible?`

Output:
- `id`, `name`
- `effectCount`
- danh sach `effectTypes`

## 4. `apply_styles`

Muc dich:
- Ap dung local styles co san de UI dong bo va hien dai hon.

Node ho tro:
- Fill style: node co fills
- Stroke style: node co strokes
- Effect style: node co effects
- Text style: TEXT

Input:
- `nodeId` bat buoc
- it nhat mot trong cac field sau:
  - `fillStyleId?`
  - `strokeStyleId?`
  - `effectStyleId?`
  - `textStyleId?`

Output:
- `id`, `name`
- danh sach style id da ap dung

## 5. `create_instance`

Muc dich:
- Su dung component thay vi tu ve lai moi thu tu dau.

Node ho tro:
- Tao INSTANCE tu local component trong file hien tai.

Input:
- phai co it nhat mot selector:
  - `componentId?`
  - `componentKey?`
  - `name?`
- `parentId?`
- `x?`, `y?`
- `componentProperties?`: object key/value de goi `instance.setProperties()`

Output:
- `id`, `name`, `type`
- `mainComponentId`
- `componentPropertyCount`
- `bounds`

Ghi chu:
- Uu tien resolve theo `componentId`, sau do `componentKey`, sau do `name`.
- `componentProperties` cho phep set variant properties va cac component properties co san.
