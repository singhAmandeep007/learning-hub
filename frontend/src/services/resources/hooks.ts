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

// Query Keys
export const resourcesKeys = {
  all: ["resources"] as const,
  lists: () => [...resourcesKeys.all, "list"] as const,
  list: (params?: GetResourcesParams) => [...resourcesKeys.lists(), params] as const,
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
    queryFn: ({ signal }) => resourcesApi.getAll(params, { signal }),
    retry: false,
    errorMessage: "Failed to load resources",
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
    queryFn: ({ signal }) => resourcesApi.getById(payload, { signal }),
    retry: false,
    errorMessage: "Failed to load resource",
    ...options,
  });
}

// Custom hook for creating a resource
export function useCreateResource(
  options?: Omit<UseMutationOptions<CreateResourceResponse, Error, CreateResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();

  return useMutationWithFlash({
    mutationFn: resourcesApi.create,
    onSuccess: (data, variables, context) => {
      // Invalidate and refetch resources list
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });

      // Call user-provided onSuccess if exists
      options?.onSuccess?.(data, variables, context);
    },
    errorMessage: "Failed to create resource",
    successMessage: "Created resource successfully",
  });
}

// Custom hook for updating a resource
export function useUpdateResource(
  options?: Omit<UseMutationOptions<UpdateResourceResponse, Error, UpdateResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();

  return useMutationWithFlash({
    mutationFn: resourcesApi.update,
    onSuccess: (data, variables, context) => {
      // Invalidate specific resource and lists
      queryClient.invalidateQueries({
        queryKey: resourcesKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });

      // Call user-provided onSuccess if exists
      options?.onSuccess?.(data, variables, context);
    },
    errorMessage: "Failed to update resource",
    successMessage: "Updated resource successfully",
    ...options,
  });
}

// Custom hook for deleting a resource
export function useDeleteResource(
  options?: Omit<UseMutationOptions<void, Error, DeleteResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();

  return useMutationWithFlash({
    mutationFn: resourcesApi.delete,
    onSuccess: (data, variables, context) => {
      // Remove specific resource from cache and invalidate lists
      queryClient.removeQueries({
        queryKey: resourcesKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });

      // Call user-provided onSuccess if exists
      options?.onSuccess?.(data, variables, context);
    },
    errorMessage: "Failed to delete resource",
    successMessage: "Deleted resource successfully",
    ...options,
  });
}
