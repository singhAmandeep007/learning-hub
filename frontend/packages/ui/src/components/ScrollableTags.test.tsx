import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { ScrollableTags, type ScrollableTagItem } from "./ScrollableTags";

const items: ScrollableTagItem[] = [
  { id: "react", label: "React", count: 12 },
  { id: "typescript", label: "TypeScript", count: 8 },
  { id: "testing", label: "Testing" },
];

describe("ScrollableTags", () => {
  it("renders labels, counts, and toolbar aria label", () => {
    render(
      <ScrollableTags
        items={items}
        selectedItemIds={[]}
        onSelectedItemIdsChange={() => undefined}
        ariaLabel="Skills filter"
      />
    );

    expect(screen.getByRole("toolbar", { name: "Skills filter" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /React/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /TypeScript/i })).toBeInTheDocument();
    expect(screen.getByText("12")).toBeInTheDocument();
    expect(screen.getByText("8")).toBeInTheDocument();
  });

  it("adds an item id when clicking an unselected tag", async () => {
    const onSelectedItemIdsChange = vi.fn();
    const user = userEvent.setup();

    render(
      <ScrollableTags
        items={items}
        selectedItemIds={["react"]}
        onSelectedItemIdsChange={onSelectedItemIdsChange}
      />
    );

    await user.click(screen.getByRole("button", { name: /TypeScript/i }));

    expect(onSelectedItemIdsChange).toHaveBeenCalledWith(["react", "typescript"]);
  });

  it("removes an item id when clicking a selected tag", async () => {
    const onSelectedItemIdsChange = vi.fn();
    const user = userEvent.setup();

    render(
      <ScrollableTags
        items={items}
        selectedItemIds={["react", "typescript"]}
        onSelectedItemIdsChange={onSelectedItemIdsChange}
      />
    );

    await user.click(screen.getByRole("button", { name: /TypeScript/i }));

    expect(onSelectedItemIdsChange).toHaveBeenCalledWith(["react"]);
  });

  it("sets aria-pressed for selected buttons", () => {
    render(
      <ScrollableTags
        items={items}
        selectedItemIds={["react"]}
        onSelectedItemIdsChange={() => undefined}
      />
    );

    expect(screen.getByRole("button", { name: /React/i })).toHaveAttribute("aria-pressed", "true");
    expect(screen.getByRole("button", { name: /TypeScript/i })).toHaveAttribute("aria-pressed", "false");
  });
});
