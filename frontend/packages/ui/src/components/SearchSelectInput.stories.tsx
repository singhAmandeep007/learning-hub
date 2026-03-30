import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, userEvent, within } from "storybook/test";

import { SearchSelectInput, type SearchSelectItem } from "./SearchSelectInput";

const defaultItems: SearchSelectItem[] = [
  { id: "react", name: "React" },
  { id: "typescript", name: "TypeScript" },
  { id: "sass", name: "Sass" },
  { id: "storybook", name: "Storybook" },
];

const meta = {
  title: "Components/SearchSelectInput",
  component: SearchSelectInput,
  args: {
    items: defaultItems,
    onSelectedItemsChange: () => undefined,
    placeholder: "Select technologies...",
    allowNewItems: false,
  },
  parameters: {
    docs: {
      description: {
        component:
          "Searchable multi-select input with keyboard support, outside-click handling, and optional add-new behavior.",
      },
    },
  },
} satisfies Meta<typeof SearchSelectInput>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: (args) => {
    function DefaultExample() {
      const [selectedItems, setSelectedItems] = useState<SearchSelectItem[]>([]);

      return (
        <SearchSelectInput
          {...args}
          initialSelectedItems={selectedItems}
          onSelectedItemsChange={setSelectedItems}
        />
      );
    }

    return <DefaultExample />;
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const input = canvas.getByRole("combobox");

    await userEvent.click(input);
    await userEvent.type(input, "React");
    await userEvent.keyboard("{Enter}");

    await expect(canvas.getByText("React")).toBeInTheDocument();
  },
};

export const AllowNewItems: Story = {
  args: {
    allowNewItems: true,
    placeholder: "Search or add tags...",
  },
  render: (args) => {
    function AllowNewItemsExample() {
      const [selectedItems, setSelectedItems] = useState<SearchSelectItem[]>([]);

      return (
        <SearchSelectInput
          {...args}
          initialSelectedItems={selectedItems}
          onSelectedItemsChange={setSelectedItems}
        />
      );
    }

    return <AllowNewItemsExample />;
  },
};
