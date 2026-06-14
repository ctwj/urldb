import { d as defineEventHandler, u as useRuntimeConfig, c as createError } from '../../nitro/nitro.mjs';
import 'node:http';
import 'node:https';
import 'node:events';
import 'node:buffer';
import 'node:fs';
import 'node:path';
import 'node:crypto';
import 'node:url';

const stats_get = defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  try {
    const response = await $fetch("/stats", {
      baseURL: String(true ? config.public.apiServer : config.public.apiBase),
      headers: {
        "Content-Type": "application/json"
      }
    });
    return response;
  } catch (error) {
    console.error("\u670D\u52A1\u7AEF\u83B7\u53D6\u7EDF\u8BA1\u6570\u636E\u5931\u8D25:", error);
    throw createError({
      statusCode: 500,
      statusMessage: "\u83B7\u53D6\u7EDF\u8BA1\u6570\u636E\u5931\u8D25"
    });
  }
});

export { stats_get as default };
//# sourceMappingURL=stats.get.mjs.map
