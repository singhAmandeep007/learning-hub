import { http, HttpResponse } from "msw";

import { DEFAULT_PRODUCT, type Resource } from "../types";

import { withDelay } from "./middleware";
import type { TDb } from "./db";

const BASE_URL = `/api/v1/${DEFAULT_PRODUCT}`;

const applyResourceFilters = (resourcesToFilter: Resource[], requestUrl: string): Resource[] => {
  const url = new URL(requestUrl);

  const search = (url.searchParams.get("search") || "").trim().toLowerCase();
  const type = (url.searchParams.get("type") || "").trim().toLowerCase();
  const tagsQuery = (url.searchParams.get("tags") || "").trim();
  const tags = tagsQuery
    ? tagsQuery
        .split(",")
        .map((tag) => tag.trim().toLowerCase())
        .filter(Boolean)
    : [];

  return resourcesToFilter.filter((resource) => {
    if (type && type !== "all" && resource.type !== type) {
      return false;
    }

    if (tags.length > 0 && !tags.every((tag) => resource.tags.map((value) => value.toLowerCase()).includes(tag))) {
      return false;
    }

    if (!search) {
      return true;
    }

    const title = resource.title.toLowerCase();
    const description = resource.description.toLowerCase();
    const resourceTags = resource.tags.join(" ").toLowerCase();

    return title.includes(search) || description.includes(search) || resourceTags.includes(search);
  });
};

const applyPagination = (resourcesToPaginate: Resource[], requestUrl: string) => {
  const url = new URL(requestUrl);

  const cursorValue = Number(url.searchParams.get("cursor") || "0");
  const limitValue = Number(url.searchParams.get("limit") || String(10));

  const cursor = Number.isFinite(cursorValue) && cursorValue > 0 ? Math.floor(cursorValue) : 0;
  const limit = Number.isFinite(limitValue) && limitValue > 0 ? Math.floor(limitValue) : resourcesToPaginate.length;

  const paginatedData = resourcesToPaginate.slice(cursor, cursor + limit);
  const nextCursor = cursor + limit;

  return {
    data: paginatedData,
    hasMore: nextCursor < resourcesToPaginate.length,
    nextCursor: nextCursor < resourcesToPaginate.length ? String(nextCursor) : "",
  };
};

export const setupHandlers = (db: TDb) => {
  return [
    http.get(BASE_URL + "/resources", ({ request }) => {
      const resources = db.resource.getAll();

      const filteredResources = applyResourceFilters(resources, request.url);
      const paginatedResponse = applyPagination(filteredResources, request.url);

      return HttpResponse.json({
        data: paginatedResponse.data,
        hasMore: paginatedResponse.hasMore,
        nextCursor: paginatedResponse.nextCursor,
      });
    }),

    http.get(BASE_URL + "/resources/:id", ({ params }) => {
      const resource = db.resource.findFirst({
        where: {
          id: {
            equals: String(params.id),
          },
        },
      });

      if (!resource) {
        return HttpResponse.json(
          {
            error: "invalid",
            message: "Resource not found",
          },
          { status: 404 }
        );
      }

      return HttpResponse.json(resource);
    }),

    http.post(
      BASE_URL + "/resources",
      withDelay(1000, async ({ request }) => {
        const formData = await request.formData();

        const newResource = {
          title: formData.get("title") as string,
          description: formData.get("description") as string,
          type: formData.get("type") as Resource["type"],
          url: (formData.get("url") as string) || "",
          thumbnailUrl: (formData.get("thumbnailUrl") as string | undefined) || "",
          tags: (formData.get("tags") as string)
            .split(",")
            .map((tag) => tag.trim())
            .filter((tag) => tag.length > 0),
        };

        const createdResource = db.resource.create(newResource);

        return HttpResponse.json(createdResource, { status: 201 });
      })
    ),

    http.patch(BASE_URL + "/resources/:id", async ({ params, request }) => {
      const formData = await request.formData();

      const resource = db.resource.findFirst({
        where: {
          id: {
            equals: String(params.id),
          },
        },
      });

      if (!resource) {
        return HttpResponse.json(
          {
            error: "invalid",
            message: "Resource not found",
          },
          { status: 404 }
        );
      }

      const updatedResource = {
        ...resource,
        ...(formData.get("title") ? { title: formData.get("title") as string } : {}),
        ...(formData.get("description") ? { description: formData.get("description") as string } : {}),
        ...(formData.get("tags")
          ? {
              tags: (formData.get("tags") as string)
                .split(",")
                .map((tag) => tag.trim())
                .filter((tag) => tag.length > 0),
            }
          : {}),
        ...(formData.get("url") ? { url: formData.get("url") as string } : {}),
        ...(formData.get("thumbnailUrl") ? { thumbnailUrl: formData.get("thumbnailUrl") as string } : {}),
      };

      const updated = db.resource.update({
        where: {
          id: {
            equals: String(params.id),
          },
        },
        data: updatedResource,
      });

      return HttpResponse.json(updated);
    }),

    http.delete(BASE_URL + "/resources/:id", async ({ params }) => {
      const resource = db.resource.findFirst({
        where: {
          id: {
            equals: String(params.id),
          },
        },
      });

      if (!resource) {
        return HttpResponse.json(
          {
            error: "invalid",
            message: "Resource not found",
          },
          { status: 404 }
        );
      }

      db.resource.delete({
        where: {
          id: {
            equals: String(params.id),
          },
        },
      });

      return HttpResponse.json({}, { status: 200 });
    }),

    http.get(BASE_URL + "/tags", () => {
      const tags = db.tag.getAll();

      return HttpResponse.json(tags);
    }),
  ];
};
