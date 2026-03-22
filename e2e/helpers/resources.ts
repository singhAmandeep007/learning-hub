import { expect, type APIRequestContext, type Page } from "@playwright/test";
import { z } from "zod";

import { apiBaseURL, resourceByIdPath, resourcesCollectionPath } from "./config";
import { expectDeleteSucceeded, expectJsonResponseWithSchema } from "./api";

export const resourceResponseSchema = z.object({
  id: z.string().min(1),
  title: z.string().min(1),
  tags: z.array(z.string()),
  type: z.enum(["article", "video", "pdf"]),
  url: z.string().url(),
});

export type ResourceResponse = z.infer<typeof resourceResponseSchema>;

export async function fillCreateArticleForm(page: Page, title: string, tag: string) {
  await page.getByRole("button", { name: "Create" }).click();
  await page.locator('input[name="title"]').fill(title);
  await page.locator(".ProseMirror").first().fill(`${title} rich description`);

  await page
    .locator(".create-update-resource-form")
    .getByRole("button", { name: /^article$/i })
    .click();

  await page.locator('input[name="url"]').fill(`https://example.com/${title}`);

  const tagsInput = page.locator(".search-select-input-wrapper input.search-input");
  await tagsInput.fill(tag);
  await tagsInput.press("Enter");
}

export function assertCreateRequestFields(requestFields: Record<string, string>, title: string, tag: string) {
  expect(requestFields.title).toBe(title);
  expect(requestFields.type).toBe("article");
  expect(requestFields.tags).toContain(tag);
  expect(requestFields.url).toBe(`https://example.com/${title}`);
}

export function assertCreatedResourceResponse(body: unknown, title: string, tag: string): ResourceResponse {
  const resource = resourceResponseSchema.parse(body);
  expect(resource.title).toBe(title);
  expect(resource.type).toBe("article");
  expect(resource.tags).toContain(tag);
  return resource;
}

export async function createArticleResourceViaApi(request: APIRequestContext, title: string, tag: string) {
  const response = await request.post(`${apiBaseURL}${resourcesCollectionPath}`, {
    multipart: {
      title,
      description: `<p>${title} description</p>`,
      type: "article",
      tags: tag,
      url: `https://example.com/${encodeURIComponent(title)}`,
    },
  });

  return await expectJsonResponseWithSchema(
    response,
    201,
    resourceResponseSchema.pick({ id: true, title: true, tags: true })
  );
}

export async function deleteResourceViaApi(request: APIRequestContext, id: string) {
  await expectDeleteSucceeded(request, resourceByIdPath(id));
}
