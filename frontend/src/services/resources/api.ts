import { httpClient } from "../httpClient";
import { getProductFromUrl } from "../utils";

import {
  type PaginatedResponse,
  type Resource,
  type GetResourcesParams,
  type GetResourceParams,
  type GetResourceResponse,
  type CreateResourcePayload,
  type CreateResourceResponse,
  type UpdateResourcePayload,
  type UpdateResourceResponse,
  type DeleteResourcePayload,
} from "../../types";

const toFormData = (payload: Partial<CreateResourcePayload>): FormData => {
  const formData = new FormData();

  if (payload.title) formData.append("title", payload.title);
  if (payload.description) formData.append("description", payload.description);
  if (payload.type) formData.append("type", payload.type);
  if (payload.tags) formData.append("tags", payload.tags);
  if (payload.url) formData.append("url", payload.url);
  if (payload.thumbnailUrl) formData.append("thumbnailUrl", payload.thumbnailUrl);
  if (payload.file) formData.append("file", payload.file);
  if (payload.thumbnail) formData.append("thumbnail", payload.thumbnail);

  return formData;
};

const adminSecret = import.meta.env["VITE_ADMIN_SECRET"] || "your-admin-secret-key";

export const resourcesApi = {
  // Get all resources with optional pagination and filtering
  getAll: async (params?: GetResourcesParams, options?: RequestInit): Promise<PaginatedResponse<Resource>> => {
    const product = getProductFromUrl();
    return httpClient.get<PaginatedResponse<Resource>>(`/${product}/resources`, params, options);
  },

  getById: async (params: GetResourceParams, options?: RequestInit): Promise<GetResourceResponse> => {
    const product = getProductFromUrl();
    return httpClient.get<GetResourceResponse>(`/${product}/resources/${params.id}`, undefined, options);
  },

  create: async (payload: CreateResourcePayload): Promise<CreateResourceResponse> => {
    const product = getProductFromUrl();
    const formData = toFormData(payload);

    return httpClient.postFormData<CreateResourceResponse>(`/${product}/resources`, {
      body: formData,
      headers: {
        AdminSecret: `${adminSecret}`,
      },
    });
  },

  update: async (payload: UpdateResourcePayload): Promise<UpdateResourceResponse> => {
    const product = getProductFromUrl();
    const { id, ...data } = payload;
    const formData = toFormData(data);

    return httpClient.patchFormData<UpdateResourceResponse>(`/${product}/resources/${id}`, {
      body: formData,
      headers: {
        AdminSecret: `${adminSecret}`,
      },
    });
  },

  // Delete resource
  delete: async (payload: DeleteResourcePayload): Promise<void> => {
    const product = getProductFromUrl();
    return httpClient.delete<void>(`/${product}/resources/${payload.id}`, {
      headers: {
        AdminSecret: `${adminSecret}`,
      },
    });
  },
};
