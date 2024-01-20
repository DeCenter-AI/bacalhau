/* tslint:disable */
/* eslint-disable */
import { http, HttpResponse, RequestHandler, RequestHandlerOptions } from "msw";
import { TestData } from "../../basic/msw.tests";

const BASE_URL = "https://localhost:1234/"

// export const fetchTasksIncompleteTaskResponse = http.get(
//   BASE_URL,
//   async (req, res, ctx) =>
//     res(
//       ctx.status(200),
//       ctx.json([
//         {
//           id: "1",
//           name: "Finish course",
//           createdOn: Date.now(),
//           status: TaskStatus.INCOMPLETE,
//         },
//       ])
//     ) as any
// )

export const sampleQueryResponse = http.get('/sampleQuery', ({ cookies }) => {
  // Placeholders for messing around with cookies
  const { v } = cookies

  const mockTestDataArray: TestData[] = [
    {
      "userId": 1234,
      "id": 1,
      "date": new Date("1970-01-01"),
      "bool": true
    },
    {
      "userId": 9876,
      "id": 2,
      "date": new Date("2023-12-31"),
      "bool": false
    },
  ]


  return HttpResponse.json(mockTestDataArray);
})


export const jobsDashboardResponse = http.get('http://localhost:1234/api/v1/orchestrator/jobs', ({ cookies }) => {
  // Placeholders for messing around with cookies
  const { v } = cookies

  return HttpResponse.json(v === 'a' ? { foo: 'a' } : { bar: 'b' })
})

export const rootResponse = http.get('http://localhost:1234/', ({ cookies }) => {
  // Placeholders for messing around with cookies
  const { v } = cookies

  return HttpResponse.json(v === 'a' ? { foo: 'a' } : { bar: 'b' })
})

export const handlers: RequestHandler<any, any, any, RequestHandlerOptions>[] = [sampleQueryResponse, rootResponse, jobsDashboardResponse]

// export const sampResp = http.get<never, RequestBody, { foo: 'a' } | { bar: 'b' }>('/', resolver)

// export const fetchTasksEmptyResponse: HttpResponseResolver = async (_req: MockedRequest, res: ResponseComposition, ctx: Context) => await res(ctx.status(200), ctx.json([]))

// export const saveTasksEmptyResponse: HttpResponseResolver = async (_req: http.MockedRequest, res: http.ResponseComposition, ctx: http.Context) => await res(ctx.status(200), ctx.json([]))

// export const handlers = [
//   fetchTasksEmptyResponse,
//   saveTasks_empty_response,
// ]
// export const loadOneJob = http.get(BASE_URL, async (req, res, ctx) =>
//   res(ctx.status(200), ctx.json([]))
// )

// export const handlers = [
//   http.get("http://localhost:1234/api/v1/*", () => passthrough()),
// ]
