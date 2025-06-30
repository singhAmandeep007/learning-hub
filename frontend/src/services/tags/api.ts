import { httpClient } from "../httpClient";
import { getProductFromUrl } from "../utils";

import { type GetTagsResponse } from "../../types";

export const tagsApi = {
  // Get all tags
  getAll: async (options?: RequestInit): Promise<GetTagsResponse> => {
    const product = getProductFromUrl();
    return httpClient.get<GetTagsResponse>(`/${product}/tags`, undefined, options);
  },
};
