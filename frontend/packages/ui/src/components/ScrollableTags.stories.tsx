import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, userEvent, within } from "storybook/test";

import { ScrollableTags, type ScrollableTagItem } from "./ScrollableTags";

const defaultItems: ScrollableTagItem[] = [
  { id: "react", label: "React", count: 12 },
  { id: "typescript", label: "TypeScript", count: 8 },
  { id: "sass", label: "Sass", count: 6 },
  { id: "storybook", label: "Storybook", count: 3 },
  { id: "testing", label: "Testing", count: 10 },
  { id: "accessibility", label: "Accessibility", count: 5 },
];

const meta = {
  title: "Components/ScrollableTags",
  component: ScrollableTags,
  args: {
    items: defaultItems,
    selectedItemIds: [],
    onSelectedItemIdsChange: () => undefined,
    ariaLabel: "Tag filters",
  },
  parameters: {
    docs: {
      description: {
        component:
          "Horizontally scrollable filter chips with keyboard-accessible toggle buttons and optional item counts.",
      },
    },
  },
} satisfies Meta<typeof ScrollableTags>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    selectedItemIds: ["react"],
  },
  render: (args) => {
    function DefaultExample() {
      const [selectedIds, setSelectedIds] = useState<string[]>(["react"]);

      return (
        <ScrollableTags
          {...args}
          selectedItemIds={selectedIds}
          onSelectedItemIdsChange={setSelectedIds}
        />
      );
    }

    return <DefaultExample />;
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const button = canvas.getByRole("button", { name: /TypeScript/i });

    await userEvent.click(button);
    await expect(button).toHaveAttribute("aria-pressed", "true");
  },
};
