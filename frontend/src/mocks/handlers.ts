import { http } from "msw";

export const handlers = [
  http.get("/resource", ({ request }) => {
    console.log(request.method, request.url);
  }),
];
