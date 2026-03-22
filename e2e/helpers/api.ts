import { expect, type APIRequestContext, type Page, type Request } from "@playwright/test";
import { z } from "zod";

import { apiBaseURL } from "./config";

type JsonResponseLike = {
  status(): number;
  json(): Promise<unknown>;
};

export async function parseMultipartRequestFields(request: Request): Promise<Record<string, string>> {
  const contentType = (await request.headerValue("content-type")) ?? "";
  const boundaryMatch = contentType.match(/boundary=(.+)$/i);

  if (!boundaryMatch?.[1]) {
    return {};
  }

  const boundary = boundaryMatch[1].trim();
  const bodyBuffer = request.postDataBuffer();

  if (!bodyBuffer) {
    return {};
  }

  const rawBody = bodyBuffer.toString("utf8");
  const parts = rawBody.split(`--${boundary}`);
  const fields: Record<string, string> = {};

  for (const part of parts) {
    const nameMatch = part.match(/name="([^"]+)"/);
    if (!nameMatch?.[1]) {
      continue;
    }

    const valueMatch = part.match(/\r\n\r\n([\s\S]*?)\r\n$/);
    if (!valueMatch?.[1]) {
      continue;
    }

    fields[nameMatch[1]] = valueMatch[1].trim();
  }

  return fields;
}

export function waitForApiRequest(page: Page, method: string, path: string, timeout = 15_000) {
  return page.waitForRequest((request) => request.method() === method && new URL(request.url()).pathname === path, {
    timeout,
  });
}

export function waitForApiResponse(page: Page, method: string, path: string, timeout = 15_000) {
  return page.waitForResponse(
    (response) => response.request().method() === method && new URL(response.url()).pathname === path,
    { timeout }
  );
}

export async function expectJsonResponse<T>(response: JsonResponseLike, expectedStatus: number): Promise<T> {
  expect(response.status()).toBe(expectedStatus);
  return (await response.json()) as T;
}

export async function expectJsonResponseWithSchema<TSchema extends z.ZodTypeAny>(
  response: JsonResponseLike,
  expectedStatus: number,
  schema: TSchema
): Promise<z.infer<TSchema>> {
  const body = await expectJsonResponse<unknown>(response, expectedStatus);
  return schema.parse(body);
}

export async function expectDeleteSucceeded(request: APIRequestContext, path: string) {
  const response = await request.delete(`${apiBaseURL}${path}`);
  expect([200, 204, 404]).toContain(response.status());
}
