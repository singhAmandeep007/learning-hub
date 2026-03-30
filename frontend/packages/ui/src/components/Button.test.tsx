import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { Button } from "./Button";

describe("Button", () => {
  it("renders children content", () => {
    render(<Button>Save</Button>);

    expect(screen.getByRole("button", { name: "Save" })).toBeInTheDocument();
  });

  it("calls onClick when clicked", async () => {
    const onClick = vi.fn();
    const user = userEvent.setup();

    render(<Button onClick={onClick}>Submit</Button>);
    await user.click(screen.getByRole("button", { name: "Submit" }));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("sets aria-pressed when pressed prop is provided", () => {
    render(<Button pressed>Toggle</Button>);

    expect(screen.getByRole("button", { name: "Toggle" })).toHaveAttribute("aria-pressed", "true");
  });

  it("uses button as the default type", () => {
    render(<Button>Default Type</Button>);

    expect(screen.getByRole("button", { name: "Default Type" })).toHaveAttribute("type", "button");
  });

  it("applies variant classes for intent and size", () => {
    render(
      <Button
        intent="secondary"
        size="lg"
      >
        Styled
      </Button>
    );

    const button = screen.getByRole("button", { name: "Styled" });
    expect(button.className).toContain("lhui-button--intent-secondary");
    expect(button.className).toContain("lhui-button--size-lg");
  });
});
