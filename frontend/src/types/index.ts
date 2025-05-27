export const RESOURCE_TYPES = ["video", "pdf", "article"] as const;

export type ResourceType = (typeof RESOURCE_TYPES)[number];

export interface Resource {
  id: string;
  title: string;
  description: string;
  type: ResourceType;
  url: string;
  thumbnailUrl?: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Tag {
  name: string;
  usageCount: number;
}

export interface ResourceFilters {
  type?: ResourceType | "all";
  tags?: string[];
  search?: string;
  limit?: number;
  cursor?: number;
}

export interface ErrorResponse {
  error: string;
  message: string;
}

export interface PaginatedResponse {
  data: Resource[];
  nextCursor: string;
  hasMore: boolean;
  total: number;
}
