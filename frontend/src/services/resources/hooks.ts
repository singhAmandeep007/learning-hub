import { useQueryClient, type UseQueryOptions, type UseMutationOptions, type QueryKey } from "@tanstack/react-query";

import { resourcesApi } from "./api";
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

import { useMutationWithFlash, useQueryWithFlash } from "../../hooks";
import { tagsKeys } from "../tags/hooks";

// Query Keys
export const resourcesKeys = {
  all: ["resources"] as const,
  lists: () => [...resourcesKeys.all, "list"] as const,
  list: (params?: GetResourcesParams) => [...resourcesKeys.lists(), JSON.stringify(params)] as const,
  details: () => [...resourcesKeys.all, "detail"] as const,
  detail: (id: string | number) => [...resourcesKeys.details(), id] as const,
} as const;

// Custom hook for getting all resources
export function useResources(
  params?: GetResourcesParams,
  options?: Omit<
    UseQueryOptions<PaginatedResponse<Resource>, Error, PaginatedResponse<Resource>, QueryKey>,
    "queryKey" | "queryFn"
  >
) {
  return useQueryWithFlash({
    queryKey: resourcesKeys.list(params),
    queryFn: () => resourcesApi.getAll(params),
    retry: false,
    errorMessage: "Failed to load resources",
    refetchOnWindowFocus: false,
    ...options,
  });
}

// Custom hook for getting a single resource by ID
export function useResource(
  payload: GetResourceParams,
  options?: Omit<UseQueryOptions<GetResourceResponse, Error, GetResourceResponse, QueryKey>, "queryKey" | "queryFn">
) {
  return useQueryWithFlash({
    queryKey: resourcesKeys.detail(payload.id),
    queryFn: () => resourcesApi.getById(payload),
    retry: false,
    errorMessage: "Failed to load resource",
    refetchOnWindowFocus: false,
    ...options,
  });
}

// Custom hook for creating a resource
export function useCreateResource(
  options?: Omit<UseMutationOptions<CreateResourceResponse, Error, CreateResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();
  const { onSuccess, ...restOptions } = options || {};

  return useMutationWithFlash({
    mutationFn: resourcesApi.create,
    onSuccess: (data, variables, context) => {
      // Invalidate and refetch resources list
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagsKeys.lists() });

      // Call user-provided onSuccess if exists
      onSuccess?.(data, variables, context);
    },
    errorMessage: "Failed to create resource",
    successMessage: "Created resource successfully",
    ...restOptions,
  });
}

// Custom hook for updating a resource
export function useUpdateResource(
  options?: Omit<UseMutationOptions<UpdateResourceResponse, Error, UpdateResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();
  const { onSuccess, ...restOptions } = options || {};

  return useMutationWithFlash({
    mutationFn: resourcesApi.update,
    onSuccess: (data, variables, context) => {
      // Invalidate specific resource and lists
      queryClient.invalidateQueries({
        queryKey: resourcesKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagsKeys.lists() });

      // Call user-provided onSuccess if exists
      onSuccess?.(data, variables, context);
    },
    errorMessage: "Failed to update resource",
    successMessage: "Updated resource successfully",
    ...restOptions,
  });
}

// Custom hook for deleting a resource
export function useDeleteResource(
  options?: Omit<UseMutationOptions<void, Error, DeleteResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();
  const { onSuccess, ...restOptions } = options || {};

  return useMutationWithFlash({
    mutationFn: resourcesApi.delete,
    onSuccess: (data, variables, context) => {
      // Call user-provided onSuccess if exists
      onSuccess?.(data, variables, context);
      // Remove specific resource from cache and invalidate lists
      queryClient.removeQueries({
        queryKey: resourcesKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagsKeys.lists() });
    },
    errorMessage: "Failed to delete resource",
    successMessage: "Deleted resource successfully",
    ...restOptions,
  });
}
