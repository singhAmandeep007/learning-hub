import { httpClient } from "../httpClient";

import { type GetTagsResponse } from "../../types";

export const tagsApi = {
  // Get all tags
  getAll: async (): Promise<GetTagsResponse> => {
    return httpClient.get<GetTagsResponse>("/tags");
  },
};
