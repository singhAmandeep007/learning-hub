import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, fn, userEvent, within } from "storybook/test";

import { Button } from "./Button";

const meta = {
  title: "Components/Button",
  component: Button,
  args: {
    children: "Action",
    intent: "primary",
    size: "md",
    disabled: false,
    onClick: fn(),
  },
  parameters: {
    docs: {
      description: {
        component:
          "Accessible button built on semantic `<button>` with variant styling powered by class-variance-authority.",
      },
    },
  },
} satisfies Meta<typeof Button>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Clickable: Story = {
  play: async ({ canvasElement, args }) => {
    const canvas = within(canvasElement);
    const button = canvas.getByRole("button", { name: "Action" });

    await userEvent.click(button);
    await expect(args.onClick).toHaveBeenCalledTimes(1);
  },
};

export const Secondary: Story = {
  args: {
    intent: "secondary",
  },
};

export const Ghost: Story = {
  args: {
    intent: "ghost",
  },
};

export const Large: Story = {
  args: {
    size: "lg",
    children: "Large Action",
  },
};

export const TogglePressed: Story = {
  args: {
    pressed: true,
    children: "Toggled",
  },
};
