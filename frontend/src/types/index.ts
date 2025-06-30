export type ErrorResponse = {
  error: string;
  message?: string;
};

export type PaginatedResponse<T> = {
  data: T[];
  hasMore: boolean;
  nextCursor?: string;
};

// Products
export const PRODUCTS = {
  ecomm: "ecomm",
  admin: "admin",
  crm: "crm",
} as const;

export type Product = (typeof PRODUCTS)[keyof typeof PRODUCTS];

export const VALID_PRODUCTS = Object.values(PRODUCTS);

export const DEFAULT_PRODUCT: Product = PRODUCTS.ecomm;

// Resource
export const RESOURCE_TYPES = {
  video: "video",
  pdf: "pdf",
  article: "article",
} as const;

export type ResourceType = (typeof RESOURCE_TYPES)[keyof typeof RESOURCE_TYPES];

export type Resource = {
  id: string;
  title: string;
  description: string;
  type: ResourceType;
  url: string;
  thumbnailUrl?: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
};

export type ResourcesFilters = {
  type?: ResourceType | "all";
  tags?: string[];
  search?: string;
};

export type GetResourcesParams = ResourcesFilters & {
  limit?: string;
  cursor?: string;
};

export type GetResourcesResponse = PaginatedResponse<Resource>;

export type GetResourceParams = Pick<Resource, "id">;

export type GetResourceResponse = Resource;

export type CreateResourcePayload = {
  title: string;
  description: string;
  type: ResourceType;
  tags: string;
  url?: string;
  thumbnailUrl?: string;
  file?: File;
  thumbnail?: File;
};

export type CreateResourceResponse = Resource;

export type UpdateResourcePayload = Partial<CreateResourcePayload> & Pick<Resource, "id">;

export type UpdateResourceResponse = Resource;

export type DeleteResourcePayload = Pick<Resource, "id">;

// Tag
export type Tag = {
  name: string;
  usageCount: number;
};

export type GetTagsResponse = Tag[];
