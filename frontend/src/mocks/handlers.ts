import { http, HttpResponse } from "msw";

import type { Resource, Tag } from "../types";

import { withDelay } from "./middleware";

const BASE_URL = "api/v1/resources";

export const handlers = [
  http.get(BASE_URL, () => {
    return HttpResponse.json({
      data: resources,
      hasMore: false,
    });

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 400 }
    );
  }),

  http.get(BASE_URL + "/:id", ({ params }) => {
    const resource = resources.find((resource) => resource.id === params.id);

    return HttpResponse.json(resource);

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 401 }
    );
  }),

  http.post(
    BASE_URL,
    withDelay(1000, async ({ request }) => {
      const formData = await request.formData();

      const newResource: Resource = {
        id: crypto.randomUUID(),

        title: formData.get("title") as string,
        description: formData.get("description") as string,

        type: formData.get("type") as Resource["type"],
        url: (formData.get("url") as string) || "",
        thumbnailUrl: (formData.get("thumbnailUrl") as string | undefined) || "",
        tags: (formData.get("tags") as string)
          .split(",")
          .map((tag) => tag.trim())
          .filter((tag) => tag.length > 0),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      resources.push(newResource);

      return HttpResponse.json(newResource, { status: 201 });

      return HttpResponse.json(
        {
          error: "invalid",
          message: "something went wrong",
        },
        { status: 500 }
      );
    })
  ),

  http.patch(BASE_URL + "/:id", async ({ params, request }) => {
    const resource = resources.find((resource) => resource.id === params.id);

    const formData = await request.formData();
    console.log("Update resource formdata", formData);

    return HttpResponse.json(resource);

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 500 }
    );
  }),

  http.delete(BASE_URL + "/:id", async ({ params }) => {
    resources = resources.filter((resource) => resource.id !== params.id);

    return HttpResponse.json({}, { status: 200 });

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 500 }
    );
  }),

  http.get(BASE_URL + "/tags", () => {
    return HttpResponse.json(tags);

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 500 }
    );
  }),
];

let resources: Resource[] = [
  {
    id: "Axq25c4sGxlkry8ex0Kd",
    title: "Test Video Tutorial 2",
    description: "This is a test video tutorial 2",
    type: "video",
    url: "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
    tags: ["test", "video", "tutorial"],
    createdAt: "2025-06-03T04:48:28.571485Z",
    updatedAt: "2025-06-03T04:48:28.571485Z",
  },
  {
    id: "ePctsVoTC36XuiGbIEwZ",
    title: "Test Video Tutorial 1",
    description: "This is a test video tutorial 1",
    type: "video",
    url: "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
    thumbnailUrl: "https://picsum.photos/id/237/536/354",
    tags: ["test", "video", "tutorial"],
    createdAt: "2025-06-03T04:48:28.522788Z",
    updatedAt: "2025-06-03T04:48:28.522788Z",
  },
  {
    id: "W3ScPIjfOJJm45M8nlWR",
    title: "Test PDF Document 2",
    description: "This is a test PDF document 2",
    type: "pdf",
    url: "https://s24.q4cdn.com/216390268/files/doc_downloads/test.pdf",
    tags: ["test", "pdf", "documentation", "creativity", "preact", "react"],
    createdAt: "2025-06-03T04:48:28.498451Z",
    updatedAt: "2025-06-03T04:48:28.498451Z",
  },
  {
    id: "GGzOhAHKEJva6L5ieS2j",
    title: "Test PDF Document 1",
    description: "This is a test PDF document 1",
    type: "pdf",
    url: "https://s24.q4cdn.com/216390268/files/doc_downloads/test.pdf",
    thumbnailUrl: "https://picsum.photos/id/237/536/354",
    tags: ["test", "pdf", "documentation"],
    createdAt: "2025-06-03T04:48:28.471915Z",
    updatedAt: "2025-06-03T04:48:28.471915Z",
  },
  {
    id: "M2wtDTqlLOcxDyLKt0LE",
    title: "Test Article 3",
    description: "This is a test article 3",
    type: "article",
    url: "https://gist.github.com/jsturgis/3b19447b304616f18657?permalink_comment_id=3658531",
    tags: ["test", "article", "blog", "productivity", "coding", "programming", "dev", "golang"],
    createdAt: "2025-06-03T04:48:28.444762Z",
    updatedAt: "2025-06-03T04:48:28.444762Z",
  },
  {
    id: "A2prtnNWN1RTxeA3x4dl",
    title: "Test Article 2",
    description: "This is a test article 2",
    type: "article",
    url: "https://gist.github.com/jsturgis/3b19447b304616f18657?permalink_comment_id=3658531",
    thumbnailUrl: "https://storage.googleapis.com/gtv-videos-bucket/sample/images/BigBuckBunny.jpg",
    tags: ["test", "article", "blog", "productivity", "coding", "programming"],
    createdAt: "2025-06-03T04:48:28.417062Z",
    updatedAt: "2025-06-03T04:48:28.417062Z",
  },
  {
    id: "Gx9Hp9hCp3El5b5BwQqo",
    title: "Test Article 1",
    description: "This is a test article 1",
    type: "article",
    url: "https://gist.github.com/jsturgis/3b19447b304616f18657?permalink_comment_id=3658531",
    thumbnailUrl: "https://storage.googleapis.com/gtv-videos-bucket/sample/images/BigBuckBunny.jpg",
    tags: ["test", "article", "blog"],
    createdAt: "2025-06-03T04:48:28.385855Z",
    updatedAt: "2025-06-03T04:48:28.385855Z",
  },
];

// eslint-disable-next-line prefer-const
let tags: Tag[] = [
  {
    name: "test",
    usageCount: 7,
  },
  {
    name: "blog",
    usageCount: 3,
  },
  {
    name: "article",
    usageCount: 3,
  },
  {
    name: "video",
    usageCount: 2,
  },
  {
    name: "tutorial",
    usageCount: 2,
  },
  {
    name: "programming",
    usageCount: 2,
  },
  {
    name: "productivity",
    usageCount: 2,
  },
  {
    name: "pdf",
    usageCount: 2,
  },
  {
    name: "documentation",
    usageCount: 2,
  },
  {
    name: "coding",
    usageCount: 2,
  },
  {
    name: "react",
    usageCount: 1,
  },
  {
    name: "preact",
    usageCount: 1,
  },
  {
    name: "golang",
    usageCount: 1,
  },
  {
    name: "dev",
    usageCount: 1,
  },
  {
    name: "creativity",
    usageCount: 1,
  },
];
