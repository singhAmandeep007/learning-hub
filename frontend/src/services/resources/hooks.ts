import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
  type UseMutationOptions,
  type QueryKey,
} from "@tanstack/react-query";

import { resourcesApi } from "./api";
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
  return useQuery({
    queryKey: resourcesKeys.list(params),
    queryFn: () => resourcesApi.getAll(params),
    ...options,
  });
}

// Custom hook for getting a single resource by ID
export function useResource(
  payload: GetResourcePayload,
  options?: Omit<UseQueryOptions<GetResourceResponse, Error, GetResourceResponse, QueryKey>, "queryKey" | "queryFn">
) {
  return useQuery({
    queryKey: resourcesKeys.detail(payload.id),
    queryFn: () => resourcesApi.getById(payload),
    ...options,
  });
}

// Custom hook for creating a resource
export function useCreateResource(
  options?: Omit<UseMutationOptions<CreateResourceResponse, Error, CreateResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: resourcesApi.create,
    onSuccess: (data, variables, context) => {
      // Invalidate and refetch resources list
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });

      // Call user-provided onSuccess if exists
      options?.onSuccess?.(data, variables, context);
    },
    ...options,
  });
}

// Custom hook for updating a resource
export function useUpdateResource(
  options?: Omit<UseMutationOptions<UpdateResourceResponse, Error, UpdateResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: resourcesApi.update,
    onSuccess: (data, variables, context) => {
      // Invalidate specific resource and lists
      queryClient.invalidateQueries({ queryKey: resourcesKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });

      // Call user-provided onSuccess if exists
      options?.onSuccess?.(data, variables, context);
    },
    ...options,
  });
}

// Custom hook for deleting a resource
export function useDeleteResource(
  options?: Omit<UseMutationOptions<void, Error, DeleteResourcePayload, unknown>, "mutationFn">
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: resourcesApi.delete,
    onSuccess: (data, variables, context) => {
      // Remove specific resource from cache and invalidate lists
      queryClient.removeQueries({ queryKey: resourcesKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: resourcesKeys.lists() });

      // Call user-provided onSuccess if exists
      options?.onSuccess?.(data, variables, context);
    },
    ...options,
  });
}
