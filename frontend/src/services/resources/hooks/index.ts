import { useQuery } from "@tanstack/react-query";
import { resourcesApi } from "..";
import type { GetResourcesParams, Resource } from "../../../types";

export const useResources = (params?: GetResourcesParams) => {
  return useQuery({
    queryKey: ["resources", params],
    queryFn: () => resourcesApi.getAll(params),
  });
};

export const useResource = (id: string) => {
  return useQuery({
    queryKey: ["resource", id],
    queryFn: () => resourcesApi.getById({ id }),
    enabled: !!id,
  });
};