import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

import { FlashProvider, ReactQueryFlashProvider } from "./FlashProvider";
import { useFlash, useReactQueryFlash } from "./useFlash";

function TriggerButtons() {
  const flash = useFlash();

  return (
    <div>
      <button
        type="button"
        onClick={() => flash.showSuccess("Saved", 10)}
      >
        success
      </button>
      <button
        type="button"
        onClick={() => flash.showMutationError({ message: "Mutation broke" }, undefined, 1500)}
      >
        mutation-error
      </button>
      <button
        type="button"
        onClick={() => flash.showQueryError({ response: { data: { message: "Query broke" } } }, undefined, 1500)}
      >
        query-error
      </button>
      <button
        type="button"
        onClick={() => flash.addNotification("Be careful", "warning")}
      >
        warning
      </button>
    </div>
  );
}

function TriggerAliasHook() {
  const flash = useReactQueryFlash();

  return (
    <button
      type="button"
      onClick={() => flash.showInfo("Alias works")}
    >
      alias
    </button>
  );
}

function HookOutsideProvider() {
  useFlash();
  return null;
}

describe("FlashProvider and hooks", () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it("throws when useFlash is called outside provider", () => {
    expect(() => render(<HookOutsideProvider />)).toThrowError("useFlash must be used within a FlashProvider");
  });

  it("renders and dismisses success notifications", async () => {
    const user = userEvent.setup();

    render(
      <FlashProvider>
        <TriggerButtons />
      </FlashProvider>
    );

    await user.click(screen.getByRole("button", { name: "success" }));
    expect(screen.getByText("Saved")).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.queryByText("Saved")).not.toBeInTheDocument();
    });
  });

  it("extracts error messages for query and mutation helpers", async () => {
    const user = userEvent.setup();

    render(
      <FlashProvider>
        <TriggerButtons />
      </FlashProvider>
    );

    await user.click(screen.getByRole("button", { name: "query-error" }));
    await user.click(screen.getByRole("button", { name: "mutation-error" }));

    expect(screen.getByText("Query broke")).toBeInTheDocument();
    expect(screen.getByText("Mutation broke")).toBeInTheDocument();
  });

  it("allows manual dismissal through close action", async () => {
    const user = userEvent.setup();

    render(
      <FlashProvider>
        <TriggerButtons />
      </FlashProvider>
    );

    await user.click(screen.getByRole("button", { name: "warning" }));
    expect(screen.getByText("Be careful")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Close notification" }));

    expect(screen.queryByText("Be careful")).not.toBeInTheDocument();
  });

  it("supports alias provider and alias hook", async () => {
    const user = userEvent.setup();

    render(
      <ReactQueryFlashProvider>
        <TriggerAliasHook />
      </ReactQueryFlashProvider>
    );

    await user.click(screen.getByRole("button", { name: "alias" }));

    expect(screen.getByText("Alias works")).toBeInTheDocument();
  });
});
