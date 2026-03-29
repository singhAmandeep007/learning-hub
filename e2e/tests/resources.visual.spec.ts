import { type Page } from "@playwright/test";
import { expect, test } from "@playwright/test";

import {
  clearAllResourcesViaApi,
  fillCreateArticleForm,
  product,
  resourcesCollectionPath,
  waitForApiResponse,
} from "../helpers";

const screenshotOptions = {
  animations: "disabled" as const,
  caret: "hide" as const,
  scale: "css" as const,
};

async function captureResourcesScreenshot(page: Page, fileName: string) {
  const resourcesRoot = page.locator(".resources");
  await expect(resourcesRoot).toBeVisible();
  await expect(resourcesRoot).toHaveScreenshot(fileName, screenshotOptions);
}

test.describe("Learning Hub visual regression", () => {
  test.describe.configure({ mode: "serial" });

  test.beforeEach(async ({ page, request }) => {
    await clearAllResourcesViaApi(request);
    await page.goto(`/${product}/resources`);
    await expect(page.getByRole("heading", { name: "Learning Hub" })).toBeVisible();

    // Hide volatile date text to keep snapshots stable across days/timezones.
    await page.addStyleTag({ content: ".resource-card-dates { visibility: hidden !important; }" });
  });

  test.afterEach(async ({ request }) => {
    await clearAllResourcesViaApi(request);
  });

  test("captures resource CRUD visual states", async ({ page }) => {
    const initialTitle = "Visual Regression Article";
    const updatedTitle = "Visual Regression Article Updated";
    const tag = "vr-tag";

    await expect(page.getByRole("heading", { name: "No resources found" })).toBeVisible();
    await captureResourcesScreenshot(page, "resources-01-empty.png");

    await fillCreateArticleForm(page, initialTitle, tag);
    const createResponsePromise = waitForApiResponse(page, "POST", resourcesCollectionPath);
    await page.getByRole("button", { name: "Create Resource" }).click();

    const createResponse = await createResponsePromise;
    expect(createResponse.status()).toBe(201);
    const createdBody = (await createResponse.json()) as { id: string };

    await expect(page.locator(".resource-card-title", { hasText: initialTitle })).toBeVisible();
    await captureResourcesScreenshot(page, "resources-02-after-create.png");

    const createdCard = page.locator(".resource-card", {
      has: page.locator(".resource-card-title", { hasText: initialTitle }),
    });

    await createdCard.locator('button[title="View Resource"]').click();
    await expect(page.locator(".resource-details-title", { hasText: initialTitle })).toBeVisible();
    await expect(page.locator(".resource-details")).toHaveScreenshot(
      "resources-03-view-details.png",
      screenshotOptions,
    );
    await page.getByRole("button", { name: "Close preview" }).click();

    await createdCard.locator('button[title="Edit Resource"]').click();
    await expect(page.getByRole("heading", { name: "Edit Resource" })).toBeVisible();
    await page.locator('input[name="title"]').fill(updatedTitle);

    const updateResponsePromise = waitForApiResponse(page, "PATCH", `${resourcesCollectionPath}/${createdBody.id}`);
    await page.getByRole("button", { name: "Update Resource" }).click();

    const updateResponse = await updateResponsePromise;
    expect(updateResponse.status()).toBe(200);

    await expect(page.locator(".resource-card-title", { hasText: updatedTitle })).toBeVisible();
    await captureResourcesScreenshot(page, "resources-04-after-update.png");

    const updatedCard = page.locator(".resource-card", {
      has: page.locator(".resource-card-title", { hasText: updatedTitle }),
    });

    page.once("dialog", (dialog) => dialog.accept());
    const deleteResponsePromise = waitForApiResponse(page, "DELETE", `${resourcesCollectionPath}/${createdBody.id}`);
    await updatedCard.locator('button[title="Delete Resource"]').click();

    const deleteResponse = await deleteResponsePromise;
    expect([200, 204]).toContain(deleteResponse.status());

    await expect(page.getByRole("heading", { name: "No resources found" })).toBeVisible();
    await captureResourcesScreenshot(page, "resources-05-after-delete.png");
  });
});
