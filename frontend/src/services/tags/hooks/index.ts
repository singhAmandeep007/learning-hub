import { useQuery } from "@tanstack/react-query";
import { tagsApi } from "..";

export const useTags = () => {
  return useQuery({
    queryKey: ["tags"],
    queryFn: tagsApi.getAll,
  });
};