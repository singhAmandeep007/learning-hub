import { faker } from "@faker-js/faker";
import { RESOURCE_TYPES, type Resource } from "../types";
import type { TDb } from "./db";
import { seedArticles, seedPdfs, seedResourceThumbnails, seedTags, seedVideos } from "./fixtures";

export const buildScenarios = (db: TDb) => {
  const builder = {
    withResources(n = 10) {
      for (let i = 0; i < n; i++) {
        const type = faker.helpers.arrayElement(Object.values(RESOURCE_TYPES)) as Resource["type"];
        const url =
          type === RESOURCE_TYPES.video
            ? faker.helpers.arrayElement(seedVideos)
            : type === RESOURCE_TYPES.pdf
              ? faker.helpers.arrayElement(seedPdfs)
              : faker.helpers.arrayElement(seedArticles);
        const createdAt = faker.date.past();
        const resource = db.resource.create({
          title: `Resource ${i + 1}`,
          description: `Description for resource ${i + 1}`,
          type,
          url,
          thumbnailUrl: faker.helpers.arrayElement(seedResourceThumbnails),
          tags: [],
          createdAt: createdAt.toISOString(),
          updatedAt: faker.date.soon({ days: 30, refDate: createdAt }).toISOString(),
        });

        const tagCount = faker.number.int({ min: 1, max: Math.min(seedTags.length, 4) });
        const tags = faker.helpers.arrayElements(seedTags, tagCount);
        resource.tags = tags;

        db.resource.update({
          where: {
            id: {
              equals: resource.id,
            },
          },
          data: { tags },
        });

        tags.forEach((tagName) => {
          const tag = db.tag.findFirst({
            where: {
              name: {
                equals: tagName,
              },
            },
          });

          if (!tag) {
            db.tag.create({
              name: tagName,
              usageCount: 1,
            });
            return;
          }

          db.tag.update({
            where: {
              id: {
                equals: tag.id,
              },
            },
            data: {
              usageCount: (tag.usageCount || 0) + 1,
            },
          });
        });
      }

      return builder;
    },
    withTags(n = 10) {
      for (let i = 0; i < n; i++) {
        db.tag.create({
          name: `Tag ${i + 1}`,
        });
      }

      return builder;
    },
  };

  return builder;
};
