# 11 — Fully Adaptive Responsive Layout Engine for Go + Fyne

**Package:** `github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/layout`  
**Status:** Authoritative Repository Standard • Universal Across Desktop, Tablet, and Mobile

---

## 1. Core Philosophy: Available Space Over Device Classes

Legacy applications often rely on rigid device checks:
```go
// ANTI-PATTERN: DO NOT DO THIS
if isMobile { ... } else if isDesktop { ... }
```

The Campus OS Responsive Layout Engine abandons device-specific classifications in favor of **geometry-driven constraint evaluation**:

> *"Given the current width and height available to the container or window, what is the most ergonomic and functional arrangement of these components?"*

The layout continuously answers this question during:
* Live desktop window border dragging and resizing
* Tablet split-screen and stage manager transitions
* Mobile phone orientation rotations (portrait $\leftrightarrow$ landscape)
* Ultrawide desktop displays ($> 21:9$)
* Awkward or constrained aspect ratios (e.g. $1920 \times 300$, $421 \times 612$, $320 \times 240$)

---

## 2. Architecture Overview

```text
Fyne Canvas / Window
        │
        ▼
   Viewport (Width, Height, Breakpoints)
        │
        ├── SizeClass:   Compact (<600) | Medium (600..900) | Expanded (>=900)
        ├── HeightClass: Short (<450)   | Normal (450..900) | Tall (>=900)
        └── Orientation: Portrait       | Landscape
        │
        ▼
Responsive Engine
        ├── Fluid Max-Width Containers (Automatic Centering, Max-Width Caps)
        ├── Dynamic Adaptive Grid (Mathematical Column Derivation)
        ├── Flex Row & Column (Weight-based Proportional Distribution)
        ├── Flow Layout (Dynamic Horizontal Wrapping)
        └── Responsive Switcher / Shell (State-retaining Mode Transitions)
        │
        ▼
Fyne Containers & Standard Widgets
```

---

## 3. Viewport & Breakpoints System

### 3.1 Viewport Abstraction
Every layout decision originates from the `Viewport` struct:
```go
type Viewport struct {
    Width             float32
    Height            float32
    Breakpoints       Breakpoints
    HeightBreakpoints HeightBreakpoints
}
```

Methods available on `Viewport`:
* `vp.SizeClass()`: `Compact`, `Medium`, or `Expanded`.
* `vp.HeightClass()`: `Short`, `Normal`, or `Tall`.
* `vp.Orientation()`: `Portrait` ($H > W$) or `Landscape` ($W \ge H$).
* `vp.AspectRatio()`: $W / H$.
* `vp.IsUltrawide()`: Returns `true` if aspect ratio $\ge 2.1$.

### 3.2 Configurable Breakpoints
Default thresholds are centralized in `layout.DefaultBreakpoints()`:
* **Horizontal (`Breakpoints`):**
  * `Medium`: `600dp`
  * `Expanded`: `900dp`
* **Vertical (`HeightBreakpoints`):**
  * `Short`: `450dp`
  * `Tall`: `900dp`

---

## 4. Reusable Layout Primitives

### 4.1 Fluid Max-Width Container (`layout.Container`)
Prevents wide monitors from stretching cards or forms into absurdly long lines:
```go
content := widget.NewCard("Admission Form", "", formFields)

fluidWrapper := layout.Container(content, layout.ContainerOptions{
    MinWidth:      320,  // Prevents collapsing below usable width
    MaxWidth:      1200, // Capped maximum width on desktop / ultrawide
    Padding:       16,
    AutoCenter:    true, // Horizontally centers content when width > MaxWidth
    ResponsivePad: true, // Scales padding (Compact: 8dp, Expanded: 24dp)
})
```

### 4.2 Dynamic Adaptive Grid (`layout.Grid`)
Calculates column count mathematically using $N = \lfloor\frac{W + \text{gap}}{\text{minItemWidth} + \text{gap}}\rfloor$. No hardcoded column assumptions:
```go
grid := layout.Grid(layout.GridOptions{
    MinItemWidth:  280, // Each card must maintain at least 280dp
    Gap:           16,  // Gap between cards
    UniformHeight: true,// All cards in a row match the tallest item
}, card1, card2, card3, card4)
```

**Live Behavior:**
* $360\text{px} \rightarrow 1\text{ column}$
* $640\text{px} \rightarrow 2\text{ columns}$
* $960\text{px} \rightarrow 3\text{ columns}$
* $1400\text{px} \rightarrow 4\text{ columns}$
* $3440\text{px} \rightarrow \text{capped fluid columns}$

### 4.3 Flex Row & Column (`layout.Row`, `layout.Column`)
Enables flex weighting and alignment:
```go
// Fixed label with auto-expanding input field and submit button
searchBar := layout.Row(8,
    widget.NewLabel("Search:"),
    layout.FlexItem(searchEntry, 1.0), // Expands to occupy available space
    submitButton,
)
```

### 4.4 Flow Layout (`layout.Flow`)
Wraps pills, badges, or chips horizontally onto subsequent rows:
```go
tagList := layout.Flow(8,
    badgeCSE,
    badgeSemester3,
    badgeLateralEntry,
    badgeHostelBlockA,
)
```

---

## 5. Responsive Switcher (`layout.Responsive`)

Separates **layout calculation** from **widget creation** for optimal performance:
* Pixel-level resizes (e.g., $1920 \rightarrow 1919 \rightarrow 1918$) perform lightweight fluid adjustments without widget re-allocations.
* Structural shifts (e.g., `Expanded` $\rightarrow$ `Medium`) transition views cleanly with cached instances.

```go
responsiveArea := layout.Responsive(layout.ResponsiveOptions{
    Compact:  compactView,   // Used when < 600dp
    Medium:   mediumView,    // Used when 600..900dp
    Expanded: expandedView,  // Used when >= 900dp
    OnChange: func(from, to layout.SizeClass, vp layout.Viewport) {
        log.Printf("Layout transition: %s -> %s (Viewport: %s)", from, to, vp)
    },
})
```

---

## 6. Adaptive Application Shell (`layout.ResponsiveAppShell`)

Provides radical structural adaptation while sharing identical state:

| Viewport | Primary Structure | Navigation Mechanism |
| :--- | :--- | :--- |
| **Expanded ($\ge 900\text{dp}$)** | Top Header + Fixed Left Sidebar ($220\text{dp}$) + Centered Content | Dedicated persistent Sidebar |
| **Medium ($600..900\text{dp}$)** | Top Header + Centered Content | Left-docked sliding Drawer with backdrop overlay |
| **Compact ($< 600\text{dp}$)** | Top Header + Centered Content + Sticky Bottom Navigation Bar | Bottom Tab Bar + Left Drawer |

```go
shell := layout.NewResponsiveAppShell(layout.AppShellOptions{
    Header:       topHeader,
    Sidebar:      desktopSidebar,
    Content:      scrollableMainContent,
    BottomNav:    mobileBottomNav,
    Drawer:       slideInDrawer,
    SidebarWidth: 220,
})

win.SetContent(shell.Root())
```

---

## 7. Testing & Verification

Run the full layout test suite:
```bash
go test -v ./apps/client/internal/layout/...
```

Launch the interactive live demonstration:
```bash
./scripts/demo.sh
```
Use the window border or quick preset buttons (`Phone`, `Tablet`, `Laptop`, `Awkward`) to inspect live layout transitions and real-time HUD stats.
