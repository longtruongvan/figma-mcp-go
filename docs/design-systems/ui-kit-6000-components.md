# Design System: UI Kit +6000 Components (Community)

## Brand Colors
- Purple: `#925FF0`
- Green: `#A3FDA1`
- Black/Gray: `#0B0B0B`
- Yellow: `#FFE74C`
- Red: `#FF5964`
- Blue: `#35A7FF`

## Color Styles (key values)
- `Main Colors/Purple/5` → `#925FF0` (primary)
- `Main Colors/Green/5` → `#A3FDA1` (secondary)
- `Main Colors/Gray/10` → `#000000` (darkest)
- `Main Colors/Gray/-1` → `#E6E6E6` (lightest)
- `White` → `#FFFFFF`

## Typography
| Style | Font | Size | Weight |
|---|---|---|---|
| H0 | Anton | 48px | Regular |
| H1 | Gotham | 36px | Black |
| H2 | Gotham | 28px | Bold |
| H3 | Gotham | 22px | Bold |
| H4 | Gotham | 20px | Book |
| P1 | Gotham | 14px | Book |
| P2 | Gotham | 12px | Light |
| P3 | Gotham | 16px | Book |
| P4 | Gotham | 10px | Book |

## Shadow Styles
- `Purple/6px` — drop shadow y:6 radius:12
- `Purple/4px` — drop shadow y:4 radius:8
- `Purple/2px` — drop shadow y:2 radius:4
- `Gray/6px` — drop shadow y:6 radius:12
- `Gray/4px` — drop shadow y:4 radius:8
- `Gray/2px` — drop shadow y:2 radius:4

## Components

### Buttons — `58:1436`
Variants: Style (Primary/Secondary/Outline/Text/Rounded/Square) × State (Default/Hover/Pressed/Focused/Disabled)

### Navigation Bars — `67:3651`
Variants: Size (complete/rounded/round) × Mode (both) × Color (purple/purple-light/white/purple-dark/black/green-light/outline) × distribution (center/left/right/icon)

### Cards (Product) — `252:9415`
Variants: Property 1 (Default/Variant2/Variant3)
- `252:9414` Default — 200×286
- `252:9416` Variant2 — 240×286
- `252:9493` Variant3 — 320×235

Other card components:
- `249:9383` Send Invitation — 370×429
- `248:8765` LogIn — 310×359
- `237:9349` Vertical Menu — 260×306

### Input Text — `147:4375`
Variants: State (Default/Hover/Focused/Filled/Disabled) × Status (Default/Error/Success) × Icon (On/Off) × Helper Text (On/Off) × Label (On/Off) × Second Label (On/Off)

### Input Tags — `151:8720`
Same variants as Input Text

### Input Big — `153:12268`
Variants: State × Status × Icon × Helper Text × Label × Secondary Label

### Tags — `150:8530`
Variants: Icon (On/Off) × Style (Disabled/Primary/Danger/Info/Success/Warning) × Border (On/Off)

### Progress Tags — `94:4202`
Variants: State (on/off) × progress (not-started/done/in-progress/Delayed) × border (on/off)

### Menu Vertical (with icon) — `237:8239`
Variants: State (Default/Hover/Focused/Pressed/Disabled/Active) × icon (on/off) × arrow (on/off)

### Menu Vertical (no icon) — `237:8929`
Variants: State × arrow (on/off)

### Checkbox — `236:8243`
Variants: Status (Active/Inactive) × State (Default/Hover/Focused/Disabled) × Label (on/off)

### Rating — `164:10449`
Variants: Number (1-5) × Style (Hearts/Stars) × Label (On/Off) × Unselected Items (On/Off)

### Switch — `131:4240`
Variants: State (on/off) × Style (main/hover/border/inactive) × Label (on/off) × Inner Label (On/Off)

### Switch Mode — `158:6795`
Variants: Mode (Light/Dark) × Style (Default/Hover/Border/Disabled) × Icon × Inner Icon × Label × Outside Icons

### Dropdown — `236:8114`
Variants: State (Default/Hover/Focused/Pressed/Disabled) × Open (True/False)

### Text Area — `164:6981`
Variants: State (Default/Hover/Focused/Disabled) × Label (On/Off) × Helper Text (On/Off) × Text Options (On/Off)

### Logo — `22:55`
Variants: Negative/Green/Purple/Positive/White

### Logo (icon) — `22:56`
Variants: Icon × (Positive/Negative/White/Purple-Only/Green-Only)

### Patterns — `305:8238`
Single pattern component

### Astronauts — `104:4196`
Variants: Default/Astronaut 2/Astronaut 3

### Escultures — `113:4175`
Variants: Number 1-12

## Icons
6 styles: Linear, Outline, Bold, Two Tone, Bulk, Broken
Categories: Money, Arrow, Weather, Users, Files, Emails, Design Tools, Location, Shop, Delivery, Security, Content, Notifications, Settings, Essential, Business, Building, Astrology, Search, Call, Programing, Archive, Time, Grid, Crypto Currency, School, Car, Computers/Devices

## Layout Grids
- Desktop: 1512px, 12 columns, 20px gutter, 100px offset
- iPad: 1024px, 8 columns, 20px gutter, 80px offset
- Mobile: 430px, 4 columns, 8px gutter, 22px offset
