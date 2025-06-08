import { http, HttpResponse } from "msw";

import type { Resource, Tag } from "../types";

import { withDelay } from "./middleware";

export const handlers = [
  http.get("api/resources", () => {
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

  http.get("api/resources/:id", ({ params }) => {
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
    "api/resources",
    withDelay(1000, async ({ request }) => {
      const formData = await request.formData();
      console.log("Create resource formdata", formData);

      return HttpResponse.json(resources[0]);

      return HttpResponse.json(
        {
          error: "invalid",
          message: "something went wrong",
        },
        { status: 500 }
      );
    })
  ),

  http.patch("api/resources/:id", async ({ params, request }) => {
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

  http.delete("api/resources/:id", async ({ params }) => {
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

  http.get("api/tags", () => {
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
    thumbnailUrl: "https://storage.googleapis.com/gtv-videos-bucket/sample/images/ForBiggerBlazes.jpg",
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
    thumbnailUrl: "https://t4.ftcdn.net/jpg/01/67/19/37/360_F_167193773_nI3NaWJMBdTTvz1EcBmqvjoeAW0WGzlu.jpg",
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
