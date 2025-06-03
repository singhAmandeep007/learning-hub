import { httpClient } from "../httpClient";

import {
  type PaginatedResponse,
  type Resource,
  type GetResourcesParams,
  type GetResourcePayload,
  type GetResourceResponse,
  type CreateResourcePayload,
  type CreateResourceResponse,
  type UpdateResourcePayload,
  type UpdateResourceResponse,
  type DeleteResourcePayload,
} from "../../types";

export const resourcesApi = {
  // Get all resources with optional pagination and filtering
  getAll: async (params?: GetResourcesParams): Promise<PaginatedResponse<Resource>> => {
    return httpClient.get<PaginatedResponse<Resource>>("/resources", params);
  },

  getById: async (payload: GetResourcePayload): Promise<GetResourceResponse> => {
    return httpClient.get<GetResourceResponse>(`/resources/${payload.id}`);
  },

  create: async (payload: CreateResourcePayload): Promise<CreateResourceResponse> => {
    return httpClient.post<CreateResourceResponse>("/resources", payload);
  },

  update: async (payload: UpdateResourcePayload): Promise<UpdateResourceResponse> => {
    const { id, ...data } = payload;
    return httpClient.patch<UpdateResourceResponse>(`/resources/${id}`, data);
  },

  // Delete resource
  delete: async (payload: DeleteResourcePayload): Promise<void> => {
    return httpClient.delete<void>(`/resources/${payload.id}`);
  },
};
