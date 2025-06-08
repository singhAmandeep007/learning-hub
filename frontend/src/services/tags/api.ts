import { httpClient } from "../httpClient";

import { type GetTagsResponse } from "../../types";

export const tagsApi = {
  // Get all tags
  getAll: async (options?: RequestInit): Promise<GetTagsResponse> => {
    return httpClient.get<GetTagsResponse>("/tags", undefined, options);
  },
};
