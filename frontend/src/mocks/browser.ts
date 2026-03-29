import { setupWorker } from "msw/browser";
import { setupHandlers } from "./handlers";

import { db } from "./db";
import { buildScenarios } from "./scenarioBuilder";

buildScenarios(db).withResources(50);

export const worker = setupWorker(...setupHandlers(db));
