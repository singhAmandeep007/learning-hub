import { faker } from "@faker-js/faker";
import { factory, primaryKey, drop } from "@mswjs/data";

import { type Resource, RESOURCE_TYPES } from "../types";

export const db = factory({
  resource: {
    id: primaryKey(faker.string.uuid),
    title: String,
    description: String,
    type: () => RESOURCE_TYPES.article as Resource["type"],
    url: String,
    thumbnailUrl: () => "",
    tags: () => [] as string[],
    createdAt: String,
    updatedAt: String,
  },
  tag: {
    id: primaryKey(faker.string.uuid),
    name: String,
    usageCount: Number,
  },
});

export type TDb = typeof db;

export const dropDb = (db: TDb) => {
  drop(db);
};
