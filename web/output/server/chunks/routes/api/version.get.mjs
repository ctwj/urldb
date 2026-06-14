import { d as defineEventHandler, u as useRuntimeConfig, c as createError } from '../../nitro/nitro.mjs';
import 'node:http';
import 'node:https';
import 'node:events';
import 'node:buffer';
import 'node:fs';
import 'node:path';
import 'node:crypto';
import 'node:url';

const version_get = defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const apiBase = String(config.public.apiServer );
  try {
    const response = await $fetch(`${apiBase}/version`);
    return response;
  } catch (error) {
    throw createError({
      statusCode: error.statusCode || 500,
      statusMessage: error.statusMessage || "\u83B7\u53D6\u7248\u672C\u4FE1\u606F\u5931\u8D25"
    });
  }
});

export { version_get as default };
//# sourceMappingURL=version.get.mjs.map
