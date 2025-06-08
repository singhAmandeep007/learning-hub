import { type DefaultBodyType, delay, type HttpResponseResolver, type PathParams } from "msw";

export function withDelay<
  Params extends PathParams,
  RequestBodyType extends DefaultBodyType,
  ResponseBodyType extends DefaultBodyType,
>(
  durationMs: number,
  resolver: HttpResponseResolver<Params, RequestBodyType, ResponseBodyType>
): HttpResponseResolver<Params, RequestBodyType, ResponseBodyType> {
  return async (...args) => {
    await delay(durationMs);
    return resolver(...args);
  };
}
