import { type UseQueryOptions, type QueryKey } from "@tanstack/react-query";
import { tagsApi } from "./api";
import { type GetTagsResponse } from "../../types";
import { useQueryWithFlash } from "../../hooks";

// Query Keys
export const tagsKeys = {
  all: ["tags"] as const,
  lists: () => [...tagsKeys.all, "list"] as const,
} as const;

// Custom hook for getting all tags
export function useTags(
  options?: Omit<UseQueryOptions<GetTagsResponse, Error, GetTagsResponse, QueryKey>, "queryKey" | "queryFn">
) {
  return useQueryWithFlash({
    queryKey: tagsKeys.lists(),
    queryFn: () => tagsApi.getAll(),
    retry: false,
    refetchOnWindowFocus: false,
    errorMessage: "Failed to load tags",
    ...options,
  });
}
