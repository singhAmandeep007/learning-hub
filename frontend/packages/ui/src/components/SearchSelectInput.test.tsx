import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { SearchSelectInput, type SearchSelectItem } from "./SearchSelectInput";

const items: SearchSelectItem[] = [
  { id: "react", name: "React" },
  { id: "typescript", name: "TypeScript" },
  { id: "storybook", name: "Storybook" },
];

describe("SearchSelectInput", () => {
  it("renders placeholder and allows selecting an existing item with keyboard", async () => {
    const onSelectedItemsChange = vi.fn();
    const user = userEvent.setup();

    render(
      <SearchSelectInput
        items={items}
        onSelectedItemsChange={onSelectedItemsChange}
        placeholder="Select technologies..."
      />
    );

    const input = screen.getByRole("combobox");
    expect(input).toHaveAttribute("placeholder", "Select technologies...");

    await user.click(input);
    await user.type(input, "Type");
    await user.keyboard("{Enter}");

    expect(onSelectedItemsChange).toHaveBeenCalledTimes(1);
    expect(onSelectedItemsChange).toHaveBeenLastCalledWith([
      expect.objectContaining({ id: "typescript", name: "TypeScript" }),
    ]);
    expect(screen.getByText("TypeScript")).toBeInTheDocument();
  });

  it("removes selected item when remove button is clicked", async () => {
    const onSelectedItemsChange = vi.fn();
    const user = userEvent.setup();

    render(
      <SearchSelectInput
        items={items}
        onSelectedItemsChange={onSelectedItemsChange}
        initialSelectedItems={[{ id: "react", name: "React" }]}
      />
    );

    await user.click(screen.getByRole("button", { name: "Remove React" }));

    expect(onSelectedItemsChange).toHaveBeenCalledWith([]);
    expect(screen.queryByRole("button", { name: "Remove React" })).not.toBeInTheDocument();
  });

  it("supports creating new items when allowNewItems is enabled", async () => {
    const onSelectedItemsChange = vi.fn();
    const user = userEvent.setup();

    render(
      <SearchSelectInput
        items={items}
        onSelectedItemsChange={onSelectedItemsChange}
        allowNewItems
      />
    );

    const input = screen.getByRole("combobox");
    await user.click(input);
    await user.type(input, "GraphQL");

    await user.click(screen.getByRole("option", { name: 'Add "GraphQL"' }));

    expect(onSelectedItemsChange).toHaveBeenCalledTimes(1);
    expect(onSelectedItemsChange).toHaveBeenLastCalledWith([expect.objectContaining({ name: "GraphQL", isNew: true })]);
    expect(screen.getByText("GraphQL")).toBeInTheDocument();
  });

  it("shows no results message when no item matches and allowNewItems is disabled", async () => {
    const user = userEvent.setup();

    render(
      <SearchSelectInput
        items={items}
        onSelectedItemsChange={() => undefined}
      />
    );

    const input = screen.getByRole("combobox");
    await user.click(input);
    await user.type(input, "NoMatch");

    expect(screen.getByText("No matching items found.")).toBeInTheDocument();
  });

  it("removes the last selected item when backspace is pressed on empty input", async () => {
    const onSelectedItemsChange = vi.fn();
    const user = userEvent.setup();

    render(
      <SearchSelectInput
        items={items}
        onSelectedItemsChange={onSelectedItemsChange}
        initialSelectedItems={[
          { id: "react", name: "React" },
          { id: "typescript", name: "TypeScript" },
        ]}
      />
    );

    const input = screen.getByRole("combobox");
    await user.click(input);
    await user.keyboard("{Backspace}");

    expect(onSelectedItemsChange).toHaveBeenLastCalledWith([expect.objectContaining({ id: "react", name: "React" })]);
    expect(screen.queryByRole("button", { name: "Remove TypeScript" })).not.toBeInTheDocument();
  });
});
