# ui-kit-6000-components

## Overview

- Nguon: file Figma `Design System _ UI kit _ +6000 Components (Community)`
- Muc tieu profile: dung de generate UI bam sat visual language cua file nay, nhung chi khi nguoi dung goi ten profile.
- Cach goi lai:
  - `Dung design system ui-kit-6000-components`
- Mac dinh:
  - Khong tu dong ap dung neu nguoi dung khong nhac ten profile nay.

## File Map

### Foundations

- `Layout Grid`
- `Colors and Shadows`
- `Fonts`
- `Typography`
- `Identificadores`

### Iconography

- `Icons`

### Components

- `Buttons`
- `Tag`
- `Inputs`
- `Switch`
- `Text Area`
- `Dropdown`
- `Rating`
- `Checkbox`
- `Menu`
- `Navigation Bars`
- `Cards`
- `Hero`

### Supporting / Showcase

- `Cover`
- `Patterns and Decoration`
- cac page divider `----- Foundations -----`, `----- Atoms -----`, `----- Molecules -----`, `----- Organisms -----`

## Design Tokens

### Color System

- He thong mau duoc to chuc theo scale `0 -> 10`, khong theo token semantic.
- Co 1 palette trung tinh:
  - `Main Colors/Gray/...`
- Co 2 palette chinh mang tinh thuong hieu:
  - `Main Colors/Purple/...`
  - `Main Colors/Green/...`
- Co nhom mau bo tro / accent:
  - `Complementary Colors/Blue/...`
  - `Complementary Colors/Red/...`
  - `Complementary Colors/Yellow/...`
- Mau nen co xu huong rat sach:
  - `White`
  - nhieu card/frame dung nen trang
- Hue chu dao de generate UI:
  - primary brand: purple
  - secondary accent: green
  - supportive accent: blue / red / yellow
  - neutral text/surface: gray scale

### Shadows

- File co 8 effect styles local:
  - `Purple/0`, `Purple/2px`, `Purple/4px`, `Purple/6px`
  - `Gray/0`, `Gray/2px`, `Gray/4px`, `Gray/6px`
- Pattern shadow:
  - radius tang dan 4 / 8 / 12
  - offset y tang dan 2 / 4 / 6
  - opacity rat nhe
- Rule:
  - shadow duoc dung de tao depth mem, khong theo huong heavy neumorphism

### Grid

- Co 3 local grid styles:
  - `Web Layout`
  - `Ipad Layout`
  - `Mobile Layout`
- He thong nay co y thuc responsive ngay tu foundation.

### Variables

- `get_variable_defs` tra ve rong.
- Nghia la file nay dang dua vao local styles + component variants, chua dua vao Figma Variables lam token system chinh.

## Typography

### Font Families

- `Anton`
- `Gotham Black`
- `Gotham`

### Local Text Styles

- `H0/Anton/Regular/Active/Left`
- `H1/Gotham/Black/Active/Center-Left`
- `H2/Gotham/Bold/Active/Center-Left`
- `H3/Gotham/Bold/Active/Center-Left`
- `H4/Gotham/Book/Active/Left`
- `P1/Gotham/Book/Active/Left`
- `P2/Gotham/Light/Active/Left`
- `P3/Gotham/Book/Active/left`
- `P4/Gotham/Book/Active/Left`

### Typography Rules

- Headline display dung `Anton` cho impact manh.
- He thong heading con lai dua vao `Gotham` / `Gotham Black`.
- Letter spacing thuong am nhe `-5%` cho phan Gotham.
- Line-height thuong rat chat, gan `100%`.
- Cam giac chung:
  - modern
  - graphic
  - co tinh editorial / marketing
  - khong phai corporate-neutral

## Iconography

### Icon Library Direction

- Page `Icons` cho thay icon duoc component hoa theo namespace `vuesax/...`.
- Co nhieu style family, it nhat thay ro:
  - `vuesax/linear/...`
  - `vuesax/broken/...`
  - mot so icon `bold` duoc dung lam status / semantic accent

### Icon Style Rules

- Mac dinh icon line rat sach, stroke toi mau `#292d32`.
- Icon duoc dat trong component rieng va duoc reuse trong button/input/card/menu.
- Khi generate UI moi:
  - uu tien `linear` cho default UI
  - dung `bold` cho status success/error/warning neu can nhan manh
  - giu cung mot family trong cung mot man

### Practical Defaults

- default icon style:
  - `vuesax/linear`
- when to use `broken`:
  - khi muon decorative accent hoac can chat minh hoa / ca tinh hon
  - khong dung lam default cho toan bo man
- when to use `bold`:
  - cho status icon, semantic emphasis, hoac diem nhan nho trong form / feedback
  - khong dung `bold` lam family icon chinh cho ca man neu khong co ly do ro rang

## Components

### Buttons

- Page `Buttons` rat day variant.
- Co `COMPONENT_SET` lon ten `Buttons`.
- Naming pattern:
  - `Content=...`
  - `Size=Small|Medium|Large`
  - `Type=Primary|Secondary|Outliine|Text-only|Rounded|Square|Negative`
  - `State=Default|Hover|Pressed|Focused|Disabled`
  - `Style=Light|Dark`
- Kieu content:
  - `Text`
  - `Text-icon`
  - `Icon`
- Rule:
  - khi generate, uu tien dung button variants thay vi tu ve rectangle + text
  - voi man sang, uu tien `Style=Light`
  - voi hero dark, co the dung `Style=Dark`

### Inputs

- Page `Inputs` co he component rat ky luong.
- Co `COMPONENT_SET` lon ten `Input Text`.
- Axes/props thay ro:
  - `State=Default|Hover|Focused|Disabled|Filled`
  - `Status=Default|Error|Success`
  - `Icon=On|Off`
  - `Helper Text=On|Off`
  - `Label=On|Off`
  - `Second Label` hoac `Secondary Label` = `On|Off`
- Input dung icon suffix/prefix theo context:
  - `eye`
  - `tick-circle`
  - `info-circle`
  - `warning-2`
  - `tick-square`
- Rule:
  - neu can form polished, uu tien dung variants san co cua `Input Text`
  - khong nen tu ve input bang frame/rectangle neu profile nay da duoc goi

### Cards

- Page `Cards` co it nhat cac mau:
  - product card
  - invitation / form card
  - login card
  - vertical menu card
  - card thong tin nho
- Card thuong:
  - nen trang
  - radius lon
  - icon vuesax
  - dung tag/button/input co san ben trong
- Rule:
  - card cua he nay co xu huong “clean marketing card”, khong qua data-dense

### Navigation Bars

- Duoc dung lai trong `Hero`.
- Co cac sub-instance:
  - `Logo`
  - `Options`
  - `Buttons`
- Rule:
  - top nav mang tinh landing page / website nhieu hon mobile app tab bar

### Hero

- Page `Hero` cho thay day la he design nghiêng ve marketing / landing hero.
- Hero component thuong compose tu:
  - `Navigation Bars`
  - visual pattern hoac astronaut / sculpture
  - large headline
  - primary CTA button
- Rule:
  - neu can tao hero section theo profile nay, uu tien visual composition manh, khong lam qua flat

## Layout Rules

- Card/frame thuong nen trang, radius lon, spacing rong.
- Co phan biet ro `Light` va `Dark` style o level component.
- He thong thich hop cho:
  - landing page
  - marketing site
  - polished SaaS sections
  - premium concept screens
- He thong khong nghieng ve:
  - iOS-native minimal
  - enterprise dashboard khac nghiet
  - brutalism

## Naming Conventions

- Text styles dat ten rat cau truc theo:
  - `Level/Font/Weight/State/Alignment`
- Buttons va inputs dat ten theo variant props thay vi ten chung chung.
- Icon components theo namespace:
  - `vuesax/<style>/<name>`

## Prompt Tips

- Khi dung profile nay, nhac ro:
  - su dung `Anton` cho display headline neu hop ngu canh
  - dung `Gotham` cho heading/body
  - uu tien palette purple / green / gray
  - dung components san co cho button, input, nav, card
  - giu icon family `vuesax/linear` dong nhat

- Prompt mau:
  - `Dung design system ui-kit-6000-components de tao mot hero landing page sang trong, dung heading display Anton, body Gotham, palette purple va green, icon family vuesax/linear, va uu tien component buttons/cards/navigation bars co san.`

## Do / Don't

### Do

- Dung component variants truoc khi tu ve tay.
- Dung local styles cho text, paint, effect neu co the.
- Dung palette purple / green / neutral lam cot song.
- Dung shadows nhe va radius lon de giu chat polished.
- Dung icons cung family trong cung mot man.

### Don't

- Khong tron nhieu icon family khac nhau.
- Khong bien style nay thanh dark cyberpunk neu khong co ly do ro rang.
- Khong dung typography system default `Inter everywhere` khi da goi profile nay.
- Khong tu ve input/button generic neu variant san co da du.
- Khong mac dinh dung profile nay neu user khong nhac ten.

## Operational Notes

- `get_local_components` toan file hien dang loi vi mot `component set` trong file co san loi variant metadata:
  - `in get_variantProperties: Component set for node has existing errors`
- Tuy vay, profile nay van su dung duoc tot thong qua:
  - page map
  - local styles
  - scan page theo component category
- Neu can generate man bam sat profile:
  - uu tien doc lai cac page `Buttons`, `Inputs`, `Cards`, `Navigation Bars`, `Hero`, `Icons` truoc khi viet vao Figma
