export type ErrorResponse = {
  error: string;
  message?: string;
};

export type PaginatedResponse<T> = {
  data: T[];
  hasMore: boolean;
};

// Resource
export const RESOURCE_TYPES = ["video", "pdf", "article"] as const;

export type ResourceType = (typeof RESOURCE_TYPES)[number];

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
  limit?: number;
  cursor?: number;
};

export type GetResourcesResponse = PaginatedResponse<Resource>;

export type GetResourcePayload = Pick<Resource, "id">;

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

// Tag
export type Tag = {
  name: string;
  usageCount: number;
};

export type GetTagsResponse = Tag[];
