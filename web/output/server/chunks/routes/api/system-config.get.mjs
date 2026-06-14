import { d as defineEventHandler, u as useRuntimeConfig } from '../../nitro/nitro.mjs';
import 'node:http';
import 'node:https';
import 'node:events';
import 'node:buffer';
import 'node:fs';
import 'node:path';
import 'node:crypto';
import 'node:url';

const systemConfig_get = defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  try {
    const response = await $fetch("/system/config", {
      baseURL: String(true ? config.public.apiServer : config.public.apiBase),
      headers: {
        "Content-Type": "application/json"
      }
    });
    return response;
  } catch (error) {
    console.error("\u670D\u52A1\u7AEF\u83B7\u53D6\u7CFB\u7EDF\u914D\u7F6E\u5931\u8D25:", error);
    return {
      site_title: "\u8001\u4E5D\u7F51\u76D8\u8D44\u6E90\u6570\u636E\u5E93",
      site_description: "\u4E00\u4E2A\u73B0\u4EE3\u5316\u7684\u8D44\u6E90\u7BA1\u7406\u7CFB\u7EDF",
      keywords: "\u7F51\u76D8\u8D44\u6E90,\u8D44\u6E90\u7BA1\u7406,\u6570\u636E\u5E93",
      author: "\u8001\u4E5D",
      copyright: "\xA9 2025 \u8001\u4E5D\u7F51\u76D8\u8D44\u6E90\u6570\u636E\u5E93"
    };
  }
});

export { systemConfig_get as default };
//# sourceMappingURL=system-config.get.mjs.map
