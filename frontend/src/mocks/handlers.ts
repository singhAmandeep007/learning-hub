import { http, HttpResponse } from "msw";

import type { Resource, Tag } from "../types";

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
      { status: 400 },
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
      { status: 401 },
    );
  }),

  http.post("api/resources", async ({ request }) => {
    console.log("MSW Intercepted: POST /api/resources");

    const formData = await request.formData();
    console.log(formData);

    return HttpResponse.json({});

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 500 },
    );
  }),

  http.patch("api/resources/:id", async ({ params, request }) => {
    const resource = resources.find((resource) => resource.id === params.id);

    const payload = await request.formData();
    console.log(payload);

    return HttpResponse.json(resource);

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 500 },
    );
  }),

  http.delete("api/resources/:id", async ({ params, request }) => {
    resources = resources.filter((resource) => resource.id === params.id);

    const payload = await request.json();
    console.log(payload);

    return HttpResponse.json({});

    return HttpResponse.json(
      {
        error: "invalid",
        message: "something went wrong",
      },
      { status: 500 },
    );
  }),

  http.get("api/tags", () => {
    return HttpResponse.json(tags);

    // return HttpResponse.json(
    //   {
    //     error: "invalid",
    //     message: "something went wrong",
    //   },
    //   { status: 500 }
    // );
  }),
];

let resources: Resource[] = [
  {
    id: "Axq25c4sGxlkry8ex0Kd",
    title: "Test Video Tutorial 2",
    description: "This is a test video tutorial 2",
    type: "video",
    url: "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/video%2F1748926108571486000_video1.mp4?alt=media",
    tags: ["test", "video", "tutorial"],
    createdAt: "2025-06-03T04:48:28.571485Z",
    updatedAt: "2025-06-03T04:48:28.571485Z",
  },
  {
    id: "ePctsVoTC36XuiGbIEwZ",
    title: "Test Video Tutorial 1",
    description: "This is a test video tutorial 1",
    type: "video",
    url: "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/video%2F1748926108522790000_video1.mp4?alt=media",
    thumbnailUrl:
      "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/image%2F1748926108536108000_image1.webp?alt=media",
    tags: ["test", "video", "tutorial"],
    createdAt: "2025-06-03T04:48:28.522788Z",
    updatedAt: "2025-06-03T04:48:28.522788Z",
  },
  {
    id: "W3ScPIjfOJJm45M8nlWR",
    title: "Test PDF Document 2",
    description: "This is a test PDF document 2",
    type: "pdf",
    url: "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/pdf%2F1748926108498453000_pdf1.pdf?alt=media",
    tags: ["test", "pdf", "documentation", "creativity", "preact", "react"],
    createdAt: "2025-06-03T04:48:28.498451Z",
    updatedAt: "2025-06-03T04:48:28.498451Z",
  },
  {
    id: "GGzOhAHKEJva6L5ieS2j",
    title: "Test PDF Document 1",
    description: "This is a test PDF document 1",
    type: "pdf",
    url: "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/pdf%2F1748926108471917000_pdf2.pdf?alt=media",
    thumbnailUrl:
      "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/image%2F1748926108475623000_image2.webp?alt=media",
    tags: ["test", "pdf", "documentation"],
    createdAt: "2025-06-03T04:48:28.471915Z",
    updatedAt: "2025-06-03T04:48:28.471915Z",
  },
  {
    id: "M2wtDTqlLOcxDyLKt0LE",
    title: "Test Article 3",
    description: "This is a test article 3",
    type: "article",
    url: "https://shorturl.at/jbcrY",
    tags: [
      "test",
      "article",
      "blog",
      "productivity",
      "coding",
      "programming",
      "dev",
      "golang",
    ],
    createdAt: "2025-06-03T04:48:28.444762Z",
    updatedAt: "2025-06-03T04:48:28.444762Z",
  },
  {
    id: "A2prtnNWN1RTxeA3x4dl",
    title: "Test Article 2",
    description: "This is a test article 2",
    type: "article",
    url: "https://shorturl.at/jbcrY",
    thumbnailUrl:
      "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/image%2F1748926108417063000_image2.webp?alt=media",
    tags: ["test", "article", "blog", "productivity", "coding", "programming"],
    createdAt: "2025-06-03T04:48:28.417062Z",
    updatedAt: "2025-06-03T04:48:28.417062Z",
  },
  {
    id: "Gx9Hp9hCp3El5b5BwQqo",
    title: "Test Article 1",
    description: "This is a test article 1",
    type: "article",
    url: "https://shorturl.at/llE4F",
    thumbnailUrl:
      "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/image%2F1748926108385860000_image1.webp?alt=media",
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
