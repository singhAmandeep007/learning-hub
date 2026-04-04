import type { components, operations } from "./api-contract.generated";

export type ErrorResponse = components["schemas"]["ErrorResponse"];

export type PaginatedResourceResponse = components["schemas"]["PaginatedResourceResponse"];

// Products
const parseProductsList = (rawValue: string | undefined): string[] => {
  const values = (rawValue ?? "")
    .split(",")
    .map((value) => value.trim())
    .filter(Boolean);

  if (values.length === 0) {
    throw new Error("VITE_VALID_PRODUCTS is required and must contain at least one product");
  }

  return [...new Set(values)];
};

export const VALID_PRODUCTS = parseProductsList(import.meta.env.VITE_VALID_PRODUCTS);

export type Product = (typeof VALID_PRODUCTS)[number];

export const DEFAULT_PRODUCT: Product = VALID_PRODUCTS[0];

// Resource
export const RESOURCE_TYPES = {
  video: "video",
  pdf: "pdf",
  article: "article",
} as const;

export type ResourceType = (typeof RESOURCE_TYPES)[keyof typeof RESOURCE_TYPES];

export type Resource = components["schemas"]["Resource"];

export type GetResourcesParams = operations["listResources"]["parameters"]["query"];

export type ResourcesFiltersType = Resource["type"] | "all";

export type GetResourcesResponse = PaginatedResourceResponse;

export type GetResourceParams = Pick<Resource, "id">;

export type GetResourceResponse = Resource;

export type CreateResourcePayload = components["schemas"]["CreateResourceRequest"];

export type CreateResourceResponse = Resource;

export type UpdateResourcePayload = Partial<CreateResourcePayload> & Pick<Resource, "id">;

export type UpdateResourceResponse = Resource;

export type DeleteResourcePayload = Pick<Resource, "id">;

// Tag
export type Tag = components["schemas"]["Tag"];

export type GetTagsResponse = Tag[];
