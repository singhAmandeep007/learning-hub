# @learning-hub/ui

A minimal, accessible React UI component library for Learning Hub.

## Installation (local workspace)

```bash
npm install
```

Then consume in the frontend app:

```tsx
import { Button } from "@learning-hub/ui";
import "@learning-hub/ui/styles.css";
```

## Scripts

- `npm run build` - Build CSS + TypeScript bundles
- `npm run check:sass-tokens` - Fail if component styles use `var(--lhui-...)` instead of Sass tokens
- `npm run lint` - Run ESLint checks with zero warnings
- `npm run lint:fix` - Auto-fix lint issues where possible
- `npm run test` - Run Vitest unit tests once
- `npm run test:watch` - Run Vitest unit tests in watch mode
- `npm run test:coverage` - Run Vitest unit tests with coverage output
- `npm run storybook` - Start Storybook
- `npm run build-storybook` - Generate static Storybook docs
- `npm run test:storybook` - Build Storybook in test mode (runs interaction-ready stories)
- `npm run typecheck` - Run TypeScript checks
- `npm run validate` - Run lint + Sass token checks + typecheck + tests + build

## Release flow (Changesets)

From `frontend/`:

```bash
npm run ui:changeset
npm run ui:version
npm run ui:build
npm run ui:publish
```

Use `npm run ui:publish:dry` before publishing to verify package contents.

## Accessibility baseline

- Uses semantic HTML first (`button` element)
- Keyboard and focus-visible support by default
- Supports ARIA state (`aria-pressed`) for toggle-like button use cases
