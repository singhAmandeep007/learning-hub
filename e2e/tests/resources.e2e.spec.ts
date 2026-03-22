import { test, expect } from "@playwright/test";

import {
  assertCreateRequestFields,
  assertCreatedResourceResponse,
  createArticleResourceViaApi,
  deleteResourceViaApi,
  expectJsonResponse,
  fillCreateArticleForm,
  parseMultipartRequestFields,
  product,
  resourcesCollectionPath,
  waitForApiRequest,
  waitForApiResponse,
} from "../helpers";

test.describe("Learning Hub end-to-end journeys", () => {
  test("creates article resource and syncs payload, response, and UI", async ({ page, request }) => {
    const title = `e2e-create-${Date.now()}`;
    const tag = `tag-${Date.now()}`;

    await page.goto(`/${product}/resources`);
    await expect(page.getByRole("heading", { name: "Learning Hub" })).toBeVisible();

    await fillCreateArticleForm(page, title, tag);

    const createRequestPromise = waitForApiRequest(page, "POST", resourcesCollectionPath);
    const createResponsePromise = waitForApiResponse(page, "POST", resourcesCollectionPath);

    await page.getByRole("button", { name: "Create Resource" }).click();

    const [createRequest, createResponse] = await Promise.all([createRequestPromise, createResponsePromise]);
    const requestFields = await parseMultipartRequestFields(createRequest);

    assertCreateRequestFields(requestFields, title, tag);

    const createdResourceRaw = await expectJsonResponse<unknown>(createResponse, 201);
    const createdResource = assertCreatedResourceResponse(createdResourceRaw, title, tag);

    await expect(page.locator(".resource-card-title", { hasText: title })).toBeVisible();

    await deleteResourceViaApi(request, createdResource.id);
  });

  test("deletes resource and reflects backend + UI state", async ({ page, request }) => {
    const title = `e2e-delete-${Date.now()}`;
    const tag = `cleanup-${Date.now()}`;
    const seeded = await createArticleResourceViaApi(request, title, tag);

    await page.goto(`/${product}/resources`);
    await expect(page.getByRole("heading", { name: "Learning Hub" })).toBeVisible();

    const seededCard = page.locator(".resource-card", {
      has: page.locator(".resource-card-title", { hasText: title }),
    });

    await expect(seededCard).toBeVisible();

    page.once("dialog", (dialog) => dialog.accept());

    const deleteResponsePromise = waitForApiResponse(page, "DELETE", `${resourcesCollectionPath}/${seeded.id}`);

    await seededCard.locator('button[title="Delete Resource"]').click();

    const deleteResponse = await deleteResponsePromise;
    expect([200, 204]).toContain(deleteResponse.status());

    await expect(page.locator(".resource-card-title", { hasText: title })).toHaveCount(0);
  });
});
