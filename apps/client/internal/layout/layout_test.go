package layout_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/layout"
)

// 1. Viewport Tests
func TestViewport_Calculations(t *testing.T) {
	bp := layout.DefaultBreakpoints()
	hbp := layout.DefaultHeightBreakpoints()

	tests := []struct {
		name        string
		w, h        float32
		wantSC      layout.SizeClass
		wantHC      layout.HeightClass
		wantOrient  layout.Orientation
		isUltrawide bool
	}{
		{"Small Phone Portrait", 360, 740, layout.Compact, layout.Normal, layout.Portrait, false},
		{"Small Phone Landscape", 740, 360, layout.Medium, layout.Short, layout.Landscape, false},
		{"Tablet Portrait", 768, 1024, layout.Medium, layout.Tall, layout.Portrait, false},
		{"Tablet Landscape", 1024, 768, layout.Expanded, layout.Normal, layout.Landscape, false},
		{"Desktop 1080p", 1920, 1080, layout.Expanded, layout.Tall, layout.Landscape, false},
		{"Ultrawide Monitor", 3440, 1440, layout.Expanded, layout.Tall, layout.Landscape, true},
		{"Awkward Wide Short", 1920, 300, layout.Expanded, layout.Short, layout.Landscape, true},
		{"Awkward Narrow Tall", 320, 960, layout.Compact, layout.Tall, layout.Portrait, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vp := layout.NewViewport(tt.w, tt.h, bp, hbp)
			if vp.SizeClass() != tt.wantSC {
				t.Errorf("SizeClass got %v, want %v", vp.SizeClass(), tt.wantSC)
			}
			if vp.HeightClass() != tt.wantHC {
				t.Errorf("HeightClass got %v, want %v", vp.HeightClass(), tt.wantHC)
			}
			if vp.Orientation() != tt.wantOrient {
				t.Errorf("Orientation got %v, want %v", vp.Orientation(), tt.wantOrient)
			}
			if vp.IsUltrawide() != tt.isUltrawide {
				t.Errorf("IsUltrawide got %v, want %v", vp.IsUltrawide(), tt.isUltrawide)
			}
			if vp.AspectRatio() != tt.w/tt.h {
				t.Errorf("AspectRatio got %v, want %v", vp.AspectRatio(), tt.w/tt.h)
			}
		})
	}
}

// 2. Breakpoint Boundary Tests (Below, At, Above)
func TestBreakpoints_Boundaries(t *testing.T) {
	bp := layout.Breakpoints{Medium: 600, Expanded: 900}
	hbp := layout.HeightBreakpoints{Short: 450, Tall: 900}

	// Horizontal thresholds:
	// Below Medium: 599.9 -> Compact
	vpBelowMed := layout.NewViewport(599.9, 600, bp, hbp)
	if vpBelowMed.SizeClass() != layout.Compact {
		t.Errorf("599.9 expected Compact, got %v", vpBelowMed.SizeClass())
	}
	// At Medium: 600.0 -> Medium
	vpAtMed := layout.NewViewport(600.0, 600, bp, hbp)
	if vpAtMed.SizeClass() != layout.Medium {
		t.Errorf("600.0 expected Medium, got %v", vpAtMed.SizeClass())
	}
	// Above Medium: 600.1 -> Medium
	vpAboveMed := layout.NewViewport(600.1, 600, bp, hbp)
	if vpAboveMed.SizeClass() != layout.Medium {
		t.Errorf("600.1 expected Medium, got %v", vpAboveMed.SizeClass())
	}

	// Below Expanded: 899.9 -> Medium
	vpBelowExp := layout.NewViewport(899.9, 600, bp, hbp)
	if vpBelowExp.SizeClass() != layout.Medium {
		t.Errorf("899.9 expected Medium, got %v", vpBelowExp.SizeClass())
	}
	// At Expanded: 900.0 -> Expanded
	vpAtExp := layout.NewViewport(900.0, 600, bp, hbp)
	if vpAtExp.SizeClass() != layout.Expanded {
		t.Errorf("900.0 expected Expanded, got %v", vpAtExp.SizeClass())
	}
	// Above Expanded: 900.1 -> Expanded
	vpAboveExp := layout.NewViewport(900.1, 600, bp, hbp)
	if vpAboveExp.SizeClass() != layout.Expanded {
		t.Errorf("900.1 expected Expanded, got %v", vpAboveExp.SizeClass())
	}

	// Vertical thresholds:
	// Below Short: 449.9 -> Short
	vpShort := layout.NewViewport(800, 449.9, bp, hbp)
	if vpShort.HeightClass() != layout.Short {
		t.Errorf("449.9 expected Short, got %v", vpShort.HeightClass())
	}
	// At Short: 450.0 -> Normal
	vpNormal := layout.NewViewport(800, 450.0, bp, hbp)
	if vpNormal.HeightClass() != layout.Normal {
		t.Errorf("450.0 expected Normal, got %v", vpNormal.HeightClass())
	}
	// At Tall: 900.0 -> Tall
	vpTall := layout.NewViewport(800, 900.0, bp, hbp)
	if vpTall.HeightClass() != layout.Tall {
		t.Errorf("900.0 expected Tall, got %v", vpTall.HeightClass())
	}
}

// 3. Grid Columns Calculation Tests
func TestGrid_ColumnsCalculation(t *testing.T) {
	minItemW := float32(280)
	gap := float32(16)

	tests := []struct {
		name       string
		availWidth float32
		wantCols   int
	}{
		{"Tiny 240px", 240, 1},
		{"Narrow 320px", 320, 1},
		{"Single Column 500px", 500, 1},
		{"Two Columns 600px", 600, 2},
		{"Two Columns 800px", 800, 2},
		{"Three Columns 900px", 900, 3},
		{"Four Columns 1200px", 1200, 4},
		{"Five Columns 1500px", 1500, 5},
		{"Ultrawide 3440px", 3440, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := layout.Columns(tt.availWidth, minItemW, gap)
			if cols != tt.wantCols {
				t.Errorf("availWidth %.0f: got %d cols, want %d", tt.availWidth, cols, tt.wantCols)
			}
		})
	}
}

// 4. Constraints Clamping Tests
func TestConstraints_Clamping(t *testing.T) {
	c := layout.Constraints{
		MinWidth:  200,
		MaxWidth:  800,
		MinHeight: 100,
		MaxHeight: 500,
		Padding:   16,
	}

	// Below Min:
	if w := c.ClampWidth(150); w != 200 {
		t.Errorf("ClampWidth(150) got %v, want 200", w)
	}
	// Within Range:
	if w := c.ClampWidth(500); w != 500 {
		t.Errorf("ClampWidth(500) got %v, want 500", w)
	}
	// Above Max:
	if w := c.ClampWidth(1200); w != 800 {
		t.Errorf("ClampWidth(1200) got %v, want 800", w)
	}

	// Effective Width with padding:
	// avail = 1000 -> 1000 - 32 = 968 -> clamped to 800
	if effW := c.EffectiveWidth(1000); effW != 800 {
		t.Errorf("EffectiveWidth(1000) got %v, want 800", effW)
	}
	// avail = 200 -> 200 - 32 = 168 -> clamped to 200
	if effW := c.EffectiveWidth(200); effW != 200 {
		t.Errorf("EffectiveWidth(200) got %v, want 200", effW)
	}
}

// 5. Fluid Max-Width Container Tests
func TestContainer_FluidMaxWidthAndCentering(t *testing.T) {
	child := widget.NewLabel("Fluid Content")
	cont := layout.Container(child, layout.ContainerOptions{
		MinWidth:   300,
		MaxWidth:   1000,
		Padding:    20,
		AutoCenter: true,
	})

	// Case A: Screen wider than MaxWidth (1600px) -> Child width capped at 1000px, centered
	cont.Resize(fyne.NewSize(1600, 800))
	childSize := child.Size()
	childPos := child.Position()

	if childSize.Width != 1000 {
		t.Errorf("Expected child width capped at 1000, got %v", childSize.Width)
	}
	expectedX := float32(1600-1000) / 2 // 300
	if childPos.X != expectedX {
		t.Errorf("Expected centered X %v, got %v", expectedX, childPos.X)
	}

	// Case B: Screen smaller than MaxWidth (800px) -> Child width is 800 - 2*20 = 760px, posX = 20
	cont.Resize(fyne.NewSize(800, 600))
	childSize = child.Size()
	childPos = child.Position()

	if childSize.Width != 760 {
		t.Errorf("Expected fluid child width 760, got %v", childSize.Width)
	}
	if childPos.X != 20 {
		t.Errorf("Expected left-aligned padding X 20, got %v", childPos.X)
	}

	// Case C: Narrow screen smaller than MinWidth (320px) -> Child width clamped to MinWidth 300
	cont.Resize(fyne.NewSize(320, 480))
	childSize = child.Size()
	if childSize.Width != 300 {
		t.Errorf("Expected clamped MinWidth 300, got %v", childSize.Width)
	}
}

// 6. Responsive Switcher & Live Transitions
func TestResponsive_Transitions(t *testing.T) {
	test.NewApp()

	compactObj := widget.NewLabel("Compact View")
	mediumObj := widget.NewLabel("Medium View")
	expandedObj := widget.NewLabel("Expanded View")

	transitions := []string{}
	rc := layout.NewResponsive(layout.ResponsiveOptions{
		Breakpoints: layout.Breakpoints{Medium: 600, Expanded: 900},
		Compact:     compactObj,
		Medium:      mediumObj,
		Expanded:    expandedObj,
		OnChange: func(from, to layout.SizeClass, vp layout.Viewport) {
			transitions = append(transitions, to.String())
		},
	})

	// Initial layout at Desktop size (1200x800) -> Expanded
	rc.Resize(fyne.NewSize(1200, 800))
	if rc.CurrentSizeClass() != layout.Expanded {
		t.Fatalf("Expected Expanded, got %v", rc.CurrentSizeClass())
	}

	// Continuous resize within Expanded: 1200 -> 1100 -> 950
	rc.Resize(fyne.NewSize(1100, 800))
	rc.Resize(fyne.NewSize(950, 800))
	if rc.CurrentSizeClass() != layout.Expanded {
		t.Fatalf("Expected still Expanded, got %v", rc.CurrentSizeClass())
	}

	// Cross into Medium (800x600)
	rc.Resize(fyne.NewSize(800, 600))
	if rc.CurrentSizeClass() != layout.Medium {
		t.Fatalf("Expected Medium, got %v", rc.CurrentSizeClass())
	}

	// Continuous resize within Medium: 800 -> 700 -> 650
	rc.Resize(fyne.NewSize(700, 600))
	rc.Resize(fyne.NewSize(650, 600))
	if rc.CurrentSizeClass() != layout.Medium {
		t.Fatalf("Expected still Medium, got %v", rc.CurrentSizeClass())
	}

	// Cross into Compact (400x700)
	rc.Resize(fyne.NewSize(400, 700))
	if rc.CurrentSizeClass() != layout.Compact {
		t.Fatalf("Expected Compact, got %v", rc.CurrentSizeClass())
	}

	// Reverse transition: Compact -> Medium (750x600)
	rc.Resize(fyne.NewSize(750, 600))
	if rc.CurrentSizeClass() != layout.Medium {
		t.Fatalf("Expected Medium after resize up, got %v", rc.CurrentSizeClass())
	}

	// Reverse transition: Medium -> Expanded (1440x900)
	rc.Resize(fyne.NewSize(1440, 900))
	if rc.CurrentSizeClass() != layout.Expanded {
		t.Fatalf("Expected Expanded after resize up, got %v", rc.CurrentSizeClass())
	}

	// Verify transitions fired correctly
	if len(transitions) < 4 {
		t.Errorf("Expected at least 4 transitions, got %d: %v", len(transitions), transitions)
	}
}

// 7. Dynamic Grid Rendering & Column Wrapping
func TestGrid_LayoutItems(t *testing.T) {
	c1 := widget.NewLabel("Item 1")
	c2 := widget.NewLabel("Item 2")
	c3 := widget.NewLabel("Item 3")
	c4 := widget.NewLabel("Item 4")

	grid := layout.Grid(layout.GridOptions{
		MinItemWidth: 200,
		Gap:          10,
	}, c1, c2, c3, c4)

	// Available width = 830 -> (830 + 10) / (200 + 10) = 4 columns
	grid.Resize(fyne.NewSize(830, 400))
	if c1.Position().Y != c4.Position().Y {
		t.Errorf("Expected 4 items on same row at width 830, but Y pos differs: c1=%v, c4=%v",
			c1.Position().Y, c4.Position().Y)
	}

	// Shrink width to 410 -> (410 + 10) / 210 = 2 columns
	grid.Resize(fyne.NewSize(410, 400))
	if c3.Position().Y <= c1.Position().Y {
		t.Errorf("Expected c3 to wrap to row 2 at width 410, c1.Y=%v, c3.Y=%v",
			c1.Position().Y, c3.Position().Y)
	}

	// Shrink width to 200 -> 1 column
	grid.Resize(fyne.NewSize(200, 600))
	if c2.Position().Y <= c1.Position().Y || c3.Position().Y <= c2.Position().Y {
		t.Errorf("Expected single column vertical stacking at width 200")
	}
}

// 8. Row, Column, Flex, and Spacer Tests
func TestRow_FlexDistribution(t *testing.T) {
	fixed1 := canvas.NewText("Fixed", nil)
	fixed1.TextSize = 12

	flex1 := canvas.NewText("Flex1", nil)
	flex2 := canvas.NewText("Flex2", nil)

	wFixed := layout.FixedSpacer(100, 30)
	wFlex1 := layout.FlexItem(flex1, 1.0)
	wFlex2 := layout.FlexItem(flex2, 2.0)

	row := layout.Row(10, wFixed, wFlex1, wFlex2)

	// Total width 430:
	// gaps = 2 * 10 = 20
	// fixed = 100
	// flex available = 430 - 20 - 100 = 310
	// flex1 (weight 1/3) = ~103.3
	// flex2 (weight 2/3) = ~206.6
	row.Resize(fyne.NewSize(430, 40))

	if wFixed.Size().Width != 100 {
		t.Errorf("Fixed item expected width 100, got %v", wFixed.Size().Width)
	}
	f1Width := wFlex1.Size().Width
	f2Width := wFlex2.Size().Width

	if f2Width <= f1Width {
		t.Errorf("Expected Flex2 (weight 2) to be larger than Flex1 (weight 1), got f1=%v, f2=%v", f1Width, f2Width)
	}
}

// 9. Flow Wrapping Tests
func TestFlow_Wrapping(t *testing.T) {
	tag1 := layout.FixedSpacer(80, 30)
	tag2 := layout.FixedSpacer(80, 30)
	tag3 := layout.FixedSpacer(80, 30)
	tag4 := layout.FixedSpacer(80, 30)

	flow := layout.Flow(10, tag1, tag2, tag3, tag4)

	// Width 400: All 4 fit horizontally (80*4 + 30 = 350)
	flow.Resize(fyne.NewSize(400, 200))
	if tag1.Position().Y != tag4.Position().Y {
		t.Errorf("Expected all tags in row 1 at width 400")
	}

	// Width 180: Only 2 fit per row (80*2 + 10 = 170) -> tag3 wraps
	flow.Resize(fyne.NewSize(180, 200))
	if tag3.Position().Y <= tag1.Position().Y {
		t.Errorf("Expected tag3 to wrap to row 2 at width 180")
	}
}

// 10. Awkward Dimensions Matrix
func TestAwkwardDimensions_Matrix(t *testing.T) {
	engine := layout.DefaultEngine

	awkwardSizes := []fyne.Size{
		fyne.NewSize(320, 240),
		fyne.NewSize(480, 320),
		fyne.NewSize(600, 400),
		fyne.NewSize(800, 480),
		fyne.NewSize(1024, 600),
		fyne.NewSize(1280, 400),
		fyne.NewSize(1280, 720),
		fyne.NewSize(1366, 768),
		fyne.NewSize(1920, 1080),
		fyne.NewSize(2560, 1440),
		fyne.NewSize(3440, 1440),
		fyne.NewSize(5120, 1440),
		fyne.NewSize(1137, 713),
		fyne.NewSize(927, 587),
		fyne.NewSize(741, 512),
		fyne.NewSize(583, 731),
		fyne.NewSize(421, 612),
	}

	for _, sz := range awkwardSizes {
		vp := engine.Viewport(sz)
		if vp.Width <= 0 || vp.Height <= 0 {
			t.Errorf("Invalid viewport size for %v", sz)
		}
		if vp.AspectRatio() <= 0 {
			t.Errorf("Invalid aspect ratio for %v: %v", sz, vp.AspectRatio())
		}
		// Confirm layout utilities never produce negative or NaN values
		w := layout.FluidWidth(sz.Width, 200, 1200, 16)
		if w < 200 && sz.Width >= 232 {
			t.Errorf("FluidWidth below minimum for size %v: %v", sz, w)
		}
		cols := layout.Columns(sz.Width, 250, 16)
		if cols < 1 {
			t.Errorf("Columns must be at least 1 for size %v, got %d", sz, cols)
		}
	}
}

// 11. Responsive App Shell
func TestAppShell_TransitionsAndDrawer(t *testing.T) {
	header := widget.NewLabel("Header")
	sidebar := widget.NewLabel("Sidebar")
	content := widget.NewLabel("Content")
	bottomNav := widget.NewLabel("BottomNav")
	drawer := widget.NewLabel("Drawer")

	shell := layout.NewResponsiveAppShell(layout.AppShellOptions{
		Header:    header,
		Sidebar:   sidebar,
		Content:   content,
		BottomNav: bottomNav,
		Drawer:    drawer,
		Breakpoints: layout.Breakpoints{
			Medium:   600,
			Expanded: 900,
		},
	})

	root := shell.Root()

	// Initial resize to Desktop
	root.Resize(fyne.NewSize(1200, 800))
	if shell.CurrentSizeClass() != layout.Expanded {
		t.Errorf("Expected Expanded shell, got %v", shell.CurrentSizeClass())
	}

	// Resize to Tablet
	root.Resize(fyne.NewSize(750, 600))
	if shell.CurrentSizeClass() != layout.Medium {
		t.Errorf("Expected Medium shell, got %v", shell.CurrentSizeClass())
	}

	// Toggle drawer on and off
	shell.ToggleDrawer()
	if !drawer.Visible() {
		t.Errorf("Expected drawer visible after ToggleDrawer")
	}
	shell.CloseDrawer()
	if drawer.Visible() {
		t.Errorf("Expected drawer hidden after CloseDrawer")
	}

	// Resize to Mobile
	root.Resize(fyne.NewSize(400, 700))
	if shell.CurrentSizeClass() != layout.Compact {
		t.Errorf("Expected Compact shell, got %v", shell.CurrentSizeClass())
	}
}

// 12. Fluid Container Text-Wrap Pre-Resize & No Left-Clipping
func TestContainer_SymmetricCenteringAndNoLeftClipping(t *testing.T) {
	lbl := widget.NewLabel("Header Title")
	cont := layout.Container(lbl, layout.ContainerOptions{
		MinWidth:      240,
		MaxWidth:      600,
		Padding:       16,
		AutoCenter:    true,
		ResponsivePad: true,
	})

	// Scenario A: Screen narrower than MinWidth (e.g. 200px)
	// Must NOT produce negative posX or clip left side of element
	cont.Resize(fyne.NewSize(200, 400))
	if lbl.Position().X < 0 {
		t.Errorf("Expected posX >= 0 on narrow screen, got %v", lbl.Position().X)
	}

	// Scenario B: Ultrawide screen (2560px)
	// Must be symmetrically centered with equal left and right margins
	cont.Resize(fyne.NewSize(2560, 1080))
	if lbl.Size().Width != 600 {
		t.Errorf("Expected child width clamped to MaxWidth 600, got %v", lbl.Size().Width)
	}
	expectedX := float32(2560-600) / 2
	if lbl.Position().X != expectedX {
		t.Errorf("Expected symmetrical posX %v, got %v", expectedX, lbl.Position().X)
	}
	rightMargin := 2560 - lbl.Position().X - lbl.Size().Width
	if lbl.Position().X != rightMargin {
		t.Errorf("Expected left margin == right margin, got left=%v, right=%v", lbl.Position().X, rightMargin)
	}
}

// 13. Dynamic Grid Height Reflow inside Scrollable/Vertical Containers
func TestGrid_DynamicHeightReflow(t *testing.T) {
	c1 := widget.NewLabel("Card 1")
	c2 := widget.NewLabel("Card 2")
	c3 := widget.NewLabel("Card 3")
	c4 := widget.NewLabel("Card 4")

	grid := layout.Grid(layout.GridOptions{
		MinItemWidth:  260,
		Gap:           16,
		UniformHeight: true,
	}, c1, c2, c3, c4)

	// Desktop width 1200px -> 4 items fit in 1 row (4 columns)
	grid.Resize(fyne.NewSize(1200, 400))
	desktopMin := grid.MinSize()

	// Mobile width 360px -> 4 items wrap into 4 rows (1 column)
	grid.Resize(fyne.NewSize(360, 800))
	mobileMin := grid.MinSize()

	if mobileMin.Height <= desktopMin.Height {
		t.Errorf("Expected mobile 4-row height to be strictly greater than desktop 1-row height, got mobile=%v, desktop=%v",
			mobileMin.Height, desktopMin.Height)
	}
}

// 14. Flow Layout Multi-Row MinSize Calculation
func TestFlow_MultiRowMinSize(t *testing.T) {
	b1 := layout.FixedSpacer(120, 36)
	b2 := layout.FixedSpacer(120, 36)
	b3 := layout.FixedSpacer(120, 36)

	flow := layout.Flow(8, b1, b2, b3)

	// Wide viewport: all fit in 1 row
	flow.Resize(fyne.NewSize(500, 100))
	singleRowMin := flow.MinSize()

	// Narrow viewport: wraps into multiple rows
	flow.Resize(fyne.NewSize(200, 200))
	multiRowMin := flow.MinSize()

	if multiRowMin.Height <= singleRowMin.Height {
		t.Errorf("Expected multi-row wrapped height > single-row height, got multi=%v, single=%v",
			multiRowMin.Height, singleRowMin.Height)
	}
}

// 15. FixedWidth Layout Container
func TestFixedWidth_FluidHeight(t *testing.T) {
	content := widget.NewLabel("Sidebar Nav")
	side := layout.FixedWidth(240, content)

	side.Resize(fyne.NewSize(240, 900))
	if content.Size().Width != 240 {
		t.Errorf("Expected content width 240, got %v", content.Size().Width)
	}
	if content.Size().Height != 900 {
		t.Errorf("Expected content to fluidly expand to full height 900, got %v", content.Size().Height)
	}
}
