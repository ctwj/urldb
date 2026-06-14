import { d as defineEventHandler, g as getQuery, u as useRuntimeConfig, c as createError } from '../../nitro/nitro.mjs';
import 'node:http';
import 'node:https';
import 'node:events';
import 'node:buffer';
import 'node:fs';
import 'node:path';
import 'node:crypto';
import 'node:url';

const resources_get = defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const query = getQuery(event);
  try {
    const response = await $fetch("/resources", {
      baseURL: String(true ? config.public.apiServer : config.public.apiBase),
      query,
      headers: {
        "Content-Type": "application/json"
      }
    });
    return response;
  } catch (error) {
    console.error("\u670D\u52A1\u7AEF\u83B7\u53D6\u8D44\u6E90\u5931\u8D25:", error);
    throw createError({
      statusCode: 500,
      statusMessage: "\u83B7\u53D6\u8D44\u6E90\u5931\u8D25"
    });
  }
});

export { resources_get as default };
//# sourceMappingURL=resources.get.mjs.map
