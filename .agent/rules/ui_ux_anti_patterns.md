# Rule: Fyne UI/UX Anti-Pattern Guardrails

When designing the UI/UX for this Go + Fyne application, avoid outdated, generic, or poorly designed desktop-application patterns.

The application must feel like a **modern 2026 product**, not a traditional legacy desktop application.

---

## 1. Do Not Make It Look Like a Legacy Desktop App

Avoid:
* Windows Forms-style interfaces
* Old-school grey panels everywhere
* Excessive borders
* Beveled buttons
* Heavy gradients
* Skeuomorphic controls
* Dense toolbar-heavy layouts
* Excessive dialog boxes
* Tiny text
* Tiny click targets
* Default-looking Fyne widgets without visual refinement

The interface should feel closer to modern products such as:
* Linear
* Notion
* Vercel
* Raycast
* Slack
* Modern admin dashboards
* Modern Material 3 applications

Do not blindly copy these products; use them as references for modern interaction and visual hierarchy.

---

## 2. Avoid Excessive Borders

Do not put a border around every component.

```text
Bad:
┌────────────────────────────────────────┐
│ ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│ │  Card    │ │  Card    │ │  Card   │ │
│ └──────────┘ └──────────┘ └─────────┘ │
└────────────────────────────────────────┘
```

Prefer hierarchy through:
* Spacing
* Background surfaces
* Typography
* Subtle elevation
* Grouping
* Alignment

Borders should communicate structure, not decorate every element.

---

## 3. Avoid Excessive Rounded Cards

Do not put everything inside a rounded rectangle.

Avoid:
```text
Card
 └── Card
      └── Card
           └── Card
```

Cards should represent meaningful groups of information. Use flat layouts when a card does not add semantic value.

---

## 4. Do Not Overuse Shadows

Avoid heavy drop shadows. Use subtle elevation only when it helps communicate:
* Floating elements
* Dialogs
* Menus
* Popovers
* Important layered surfaces

Most content should rely on spacing and contrast rather than shadows.

---

## 5. Avoid Excessive Color

Do not use a different bright color for every component. Use a restrained design system:
```text
Primary
Secondary
Background
Surface
Text
Muted Text
Border
Success
Warning
Error
Info
```

Semantic colors should communicate meaning rather than decoration.

---

## 6. Establish Strong Visual Hierarchy

Every screen should clearly communicate:
```text
Where am I?
What is this page?
What is important?
What can I do?
What requires my attention?
```

Use typography, spacing, size, weight, contrast, and positioning to establish hierarchy. Do not make every element visually equally important.

---

## 7. Do Not Cram Information Onto the Screen

Avoid trying to show everything simultaneously, especially for dashboards, student records, attendance, complaints, timetables, and administration screens.

Use progressive disclosure where appropriate: show the important information first and allow users to drill into details.

---

## 8. Avoid Tiny Text

Do not compensate for information density by making text extremely small. Prioritize readability. Secondary information can be smaller, but it must remain comfortably readable across desktop and mobile devices.

---

## 9. Avoid Tiny Click Targets

Interactive elements should have comfortable touch/click areas (min 44×44px touch target). The application must run seamlessly on Windows, macOS, Linux, Android, and iOS. Do not design exclusively around a mouse.

---

## 10. Do Not Rely Only on Icons

Avoid ambiguous icon-only controls without context. Use clear text labels (`Edit`, `Delete`, `Download`, `More`) where the action isn't universally obvious. Icon-only controls must have appropriate tooltips/accessibility information.

---

## 11. Avoid Modal Dialog Overuse

Do not turn every interaction into a modal dialog. Prefer inline interactions when possible. Use dialogs only for genuinely interruptive actions:
* Confirmation of high-impact changes
* Destructive actions
* Focused forms
* Critical regulatory warnings

Do not use a dialog merely because it is convenient to implement.

---

## 12. Avoid Nested Dialogs

Never create:
```text
Dialog
 └── Dialog
      └── Dialog
```

This creates cognitive overload. Prefer dedicated screens, sheets, inline editing, or a single focused dialog.

---

## 13. Avoid Infinite Scrolling for Administrative Data

For data-heavy administrative systems, users need to locate specific records, understand position, compare rows, jump pages, and filter. Use appropriate combinations of search, filters, sorting, pagination, and grouping rather than infinite scrolling.

---

## 14. Tables Must Be Scannable

Avoid excessive vertical borders, excessive colors, unnecessary columns, tiny typography, enormous row heights, and random alignment. Use consistent column alignment, row spacing, status indicators, and actions.

---

## 15. Avoid Showing Raw Status as Text Everywhere

Instead of repeatedly displaying:
```text
Status: ACTIVE
Status: PENDING
Status: FAILED
```

Use appropriate visual status indicators:
```text
● Active
● Pending
● Failed
```

with text retained for clarity (color + text/icon rather than color alone).

---

## 16. Do Not Hide Primary Actions

Every screen must have an obvious primary action. Do not make users hunt through menus to discover the main action.

---

## 17. Avoid Action Overload

Do not place 15 buttons in a toolbar. Prioritize:
1. Primary action
2. Secondary actions
3. Overflow actions (`⋮`)

---

## 18. Use Consistent Navigation

Do not change navigation behavior from screen to screen. The application should maintain a predictable structure across domains.

---

## 19. Avoid Deep Navigation

Do not force users through deep navigation trees when search or direct filtering can reach the target immediately.

---

## 20. Always Provide Feedback

Every meaningful action must provide feedback:
```text
Saving... → Saved
Failed to save
Syncing... → Synced
Uploading... → Upload complete
```

Never leave users wondering whether their action succeeded or failed.

---

## 21. Avoid Fake Loading

Do not use unnecessary or decorative loading animations. Loading states must communicate actual work. Prefer descriptive status or skeleton placeholders.

---

## 22. Design Proper Empty States

Do not show empty containers without explanation. Empty states must explain:
1. What happened?
2. Why is it empty?
3. What can the user do? (with an actionable button)

---

## 23. Design Proper Error States

Never leave a blank screen after an error. Provide clear messaging, actionable recovery (`[ Try Again ]`), and where appropriate, an error/reference ID.

---

## 24. Do Not Use Color as the Only Indicator

Never communicate meaning exclusively through color. Always pair color with text, badge labels, or vector icons for accessibility.

---

## 25. Maintain Consistent Spacing

Use a disciplined spacing system: `4`, `8`, `12`, `16`, `24`, `32`, `48`, `64`. Avoid arbitrary, one-off pixel values throughout the UI.

---

## 26. Maintain Consistent Typography

Define a clear hierarchy:
* Page Title (Bold, 20–24pt)
* Section Title (Bold, 16–18pt)
* Card Title (Semi-bold, 14–15pt)
* Body (Regular, 13–14pt)
* Secondary / Description (Regular, 12–13pt)
* Caption / Badges (Medium, 10–11pt)

---

## 27. Avoid Unnecessary Animations

Animations should communicate state changes, navigation, hierarchy, progress, or feedback. Avoid bouncing, slow transitions, or decorative animations. The application must feel fast and responsive.

---

## 28. Do Not Make Every Screen Look Identical

Different tasks require tailored compositions:
* Dashboard $\rightarrow$ overview
* Records $\rightarrow$ searchable table
* Attendance $\rightarrow$ data entry + statistics
* Timetable $\rightarrow$ calendar/grid
* Settings $\rightarrow$ grouped forms

Maintain the same design system while adapting layouts to the task.

---

## 29. Avoid Desktop-Only Interaction Patterns

Because the application targets Linux, macOS, Windows, Android, and iOS, essential functionality must remain accessible through touch-friendly controls rather than relying exclusively on right-click context menus, hover states, or drag-and-drop.

---

## 30. Do Not Treat Fyne's Default Appearance as the Final Design

Fyne's default widgets are raw building blocks. Apply the application-level design system: custom tokens, custom padding, high-contrast typography, and Lucide vector icons.

---

## 31. Avoid Inconsistent Component Styling

A button or input must not look different between screens. Use the standard variants (`ButtonDefault`, `ButtonSecondary`, `ButtonOutline`, `ButtonDestructive`, `ButtonGhost`).

---

## 32. Design for Information Density Intentionally

Use task-appropriate density:
* Dashboard $\rightarrow$ comfortable
* Administrative table $\rightarrow$ compact
* Multi-step forms $\rightarrow$ comfortable
* Mobile $\rightarrow$ touch-friendly

---

## 33. Use Progressive Disclosure

Show complexity only when needed. Do not dump every available field onto the first screen. Present summary information first and permit drilling down into details.

---

## 34. Respect Platform Conventions

The UI should feel natural on each platform while maintaining the same design language. Do not force desktop interaction patterns onto mobile or make mobile UI behave like a shrunken desktop application.

---

## 35. Overall Design Principle & Quality Gate

The UI must feel **modern, clean, fast, calm, predictable, accessible, information-efficient, touch-friendly, responsive, and professional**.

### Mandatory Pre-Creation Checklist
Before creating or modifying any UI component or screen, verify:
1. Is the hierarchy immediately clear?
2. Is the primary action obvious?
3. Is there unnecessary visual noise?
4. Is the information density appropriate?
5. Can the interface work comfortably with mouse, keyboard, and touch?
6. Does it work on both desktop and mobile form factors?
7. Is the component consistent with the application's design system?
8. Does every visual element serve a purpose?
9. Does the screen communicate loading, empty, success, and error states?
10. Does this actually improve the user's workflow?
