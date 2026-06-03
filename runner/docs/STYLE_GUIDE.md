# DaoForge Cultivation Style Guide

## Design Intent

DaoForge should feel like a quiet cultivation sect library for algorithm practice. The interface should be calm, focused, and atmospheric: ink-wash paper, jade, pine, mist, gold accents, soft depth, and language around realms, qi, manuals, meditation, and breakthroughs.

The theme should support concentration. Avoid loud arcade styling, heavy outlines, and offset-shadow visual language.

## Visual Principles

- Calm clarity before decoration.
- Soft cards, light borders, subtle shadows, and misty backgrounds.
- Jade and pine for primary structure, gold for progress and highlights, cinnabar only for errors or danger.
- Use cultivation flavor in labels, badges, and progress states without obscuring core programming tasks.
- Keep code editor surfaces practical and high contrast.

## Color Palette

| Token | Hex | Usage |
| --- | --- | --- |
| Ink | `#18231F` | Primary text |
| Pine | `#173F35` | Headings, navigation, strong surfaces |
| Jade | `#3F8F72` | Primary actions, success, solved states |
| Deep Jade | `#1F5F50` | Gradients and active states |
| Moss | `#8AA66D` | Secondary accents |
| Mist | `#EEF3EC` | Main page background |
| Paper | `#FBF7EA` | Panels and cards |
| Silk | `#FFFDF6` | Inputs and bright surfaces |
| Gold | `#C99A35` | Progress, realm badges, highlights |
| Cinnabar | `#B94F3B` | Error and hard difficulty |
| Stone | `#6F7B73` | Muted text |
| Night Ink | `#10201D` | Code editor shell |

## Typography

- Use a readable sans-serif for UI and a monospace font for code.
- Headings should feel refined and grounded, not shouted.
- Do not use all-caps everywhere.
- Letter spacing remains `0`.
- Keep body text at 16px or larger.

## Component Style

- Buttons are rounded pills with soft hover lift.
- Cards use `20-24px` radius, `1px` translucent borders, and soft shadows.
- Inputs use silk/paper fills, jade focus rings, and concise labels.
- Difficulty badges map to realms:
  - Easy: Foundation Realm
  - Medium: Qi Condensation
  - Hard: Core Formation
- Progress currency is Spirit Stones.
- Streaks are Meditation Streaks.
- Problems are Manuals.
- Dashboard is the Sect Hall.
- Submission history is the Karma Ledger.
- Running visible tests is Meditate.
- Submitting all tests is Attempt Breakthrough.

## Page Notes

### Landing

The first screen should immediately establish DaoForge as a cultivation coding platform. Use calm copy, realm progression, and a dark pine cultivation board.

### Dashboard

The Sect Hall should show mastered manuals, meditation streak, spirit stones, accuracy, best language, recent attempts, and recommended manuals.

### Manual Library

The catalog should feel like a searchable sect library. Filters should be simple and persistent.

### Problem Detail

The problem page should be practical: manual prompt, examples, constraints, editor, and result panel. Theme language should appear in headings and actions, while technical content stays precise.
