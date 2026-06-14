import process from 'node:process';globalThis._importMeta_=globalThis._importMeta_||{url:"file:///_entry.js",env:process.env};import ft, { Server as Server$1 } from 'node:http';
import Qa, { Server } from 'node:https';
import { EventEmitter } from 'node:events';
import { Buffer as Buffer$1 } from 'node:buffer';
import { promises, existsSync } from 'node:fs';
import { resolve as resolve$1, dirname as dirname$1, join } from 'node:path';
import { createHash } from 'node:crypto';
import { fileURLToPath } from 'node:url';

const suspectProtoRx = /"(?:_|\\u0{2}5[Ff]){2}(?:p|\\u0{2}70)(?:r|\\u0{2}72)(?:o|\\u0{2}6[Ff])(?:t|\\u0{2}74)(?:o|\\u0{2}6[Ff])(?:_|\\u0{2}5[Ff]){2}"\s*:/;
const suspectConstructorRx = /"(?:c|\\u0063)(?:o|\\u006[Ff])(?:n|\\u006[Ee])(?:s|\\u0073)(?:t|\\u0074)(?:r|\\u0072)(?:u|\\u0075)(?:c|\\u0063)(?:t|\\u0074)(?:o|\\u006[Ff])(?:r|\\u0072)"\s*:/;
const JsonSigRx = /^\s*["[{]|^\s*-?\d{1,16}(\.\d{1,17})?([Ee][+-]?\d+)?\s*$/;
function jsonParseTransform(key, value) {
  if (key === "__proto__" || key === "constructor" && value && typeof value === "object" && "prototype" in value) {
    warnKeyDropped(key);
    return;
  }
  return value;
}
function warnKeyDropped(key) {
  console.warn(`[destr] Dropping "${key}" key to prevent prototype pollution.`);
}
function destr(value, options = {}) {
  if (typeof value !== "string") {
    return value;
  }
  if (value[0] === '"' && value[value.length - 1] === '"' && value.indexOf("\\") === -1) {
    return value.slice(1, -1);
  }
  const _value = value.trim();
  if (_value.length <= 9) {
    switch (_value.toLowerCase()) {
      case "true": {
        return true;
      }
      case "false": {
        return false;
      }
      case "undefined": {
        return void 0;
      }
      case "null": {
        return null;
      }
      case "nan": {
        return Number.NaN;
      }
      case "infinity": {
        return Number.POSITIVE_INFINITY;
      }
      case "-infinity": {
        return Number.NEGATIVE_INFINITY;
      }
    }
  }
  if (!JsonSigRx.test(value)) {
    if (options.strict) {
      throw new SyntaxError("[destr] Invalid JSON");
    }
    return value;
  }
  try {
    if (suspectProtoRx.test(value) || suspectConstructorRx.test(value)) {
      if (options.strict) {
        throw new Error("[destr] Possible prototype pollution");
      }
      return JSON.parse(value, jsonParseTransform);
    }
    return JSON.parse(value);
  } catch (error) {
    if (options.strict) {
      throw error;
    }
    return value;
  }
}

const HASH_RE = /#/g;
const AMPERSAND_RE = /&/g;
const SLASH_RE = /\//g;
const EQUAL_RE = /=/g;
const PLUS_RE = /\+/g;
const ENC_CARET_RE = /%5e/gi;
const ENC_BACKTICK_RE = /%60/gi;
const ENC_PIPE_RE = /%7c/gi;
const ENC_SPACE_RE = /%20/gi;
const ENC_SLASH_RE = /%2f/gi;
function encode(text) {
  return encodeURI("" + text).replace(ENC_PIPE_RE, "|");
}
function encodeQueryValue(input) {
  return encode(typeof input === "string" ? input : JSON.stringify(input)).replace(PLUS_RE, "%2B").replace(ENC_SPACE_RE, "+").replace(HASH_RE, "%23").replace(AMPERSAND_RE, "%26").replace(ENC_BACKTICK_RE, "`").replace(ENC_CARET_RE, "^").replace(SLASH_RE, "%2F");
}
function encodeQueryKey(text) {
  return encodeQueryValue(text).replace(EQUAL_RE, "%3D");
}
function decode$2(text = "") {
  try {
    return decodeURIComponent("" + text);
  } catch {
    return "" + text;
  }
}
function decodePath(text) {
  return decode$2(text.replace(ENC_SLASH_RE, "%252F"));
}
function decodeQueryKey(text) {
  return decode$2(text.replace(PLUS_RE, " "));
}
function decodeQueryValue(text) {
  return decode$2(text.replace(PLUS_RE, " "));
}

function parseQuery(parametersString = "") {
  const object = /* @__PURE__ */ Object.create(null);
  if (parametersString[0] === "?") {
    parametersString = parametersString.slice(1);
  }
  for (const parameter of parametersString.split("&")) {
    const s = parameter.match(/([^=]+)=?(.*)/) || [];
    if (s.length < 2) {
      continue;
    }
    const key = decodeQueryKey(s[1]);
    if (key === "__proto__" || key === "constructor") {
      continue;
    }
    const value = decodeQueryValue(s[2] || "");
    if (object[key] === void 0) {
      object[key] = value;
    } else if (Array.isArray(object[key])) {
      object[key].push(value);
    } else {
      object[key] = [object[key], value];
    }
  }
  return object;
}
function encodeQueryItem(key, value) {
  if (typeof value === "number" || typeof value === "boolean") {
    value = String(value);
  }
  if (!value) {
    return encodeQueryKey(key);
  }
  if (Array.isArray(value)) {
    return value.map(
      (_value) => `${encodeQueryKey(key)}=${encodeQueryValue(_value)}`
    ).join("&");
  }
  return `${encodeQueryKey(key)}=${encodeQueryValue(value)}`;
}
function stringifyQuery(query) {
  return Object.keys(query).filter((k) => query[k] !== void 0).map((k) => encodeQueryItem(k, query[k])).filter(Boolean).join("&");
}

const PROTOCOL_STRICT_REGEX = /^[\s\w\0+.-]{2,}:([/\\]{1,2})/;
const PROTOCOL_REGEX = /^[\s\w\0+.-]{2,}:([/\\]{2})?/;
const PROTOCOL_RELATIVE_REGEX = /^([/\\]\s*){2,}[^/\\]/;
const PROTOCOL_SCRIPT_RE = /^[\s\0]*(blob|data|javascript|vbscript):$/i;
const TRAILING_SLASH_RE = /\/$|\/\?|\/#/;
const JOIN_LEADING_SLASH_RE = /^\.?\//;
function hasProtocol(inputString, opts = {}) {
  if (typeof opts === "boolean") {
    opts = { acceptRelative: opts };
  }
  if (opts.strict) {
    return PROTOCOL_STRICT_REGEX.test(inputString);
  }
  return PROTOCOL_REGEX.test(inputString) || (opts.acceptRelative ? PROTOCOL_RELATIVE_REGEX.test(inputString) : false);
}
function isScriptProtocol(protocol) {
  return !!protocol && PROTOCOL_SCRIPT_RE.test(protocol);
}
function hasTrailingSlash(input = "", respectQueryAndFragment) {
  if (!respectQueryAndFragment) {
    return input.endsWith("/");
  }
  return TRAILING_SLASH_RE.test(input);
}
function withoutTrailingSlash(input = "", respectQueryAndFragment) {
  if (!respectQueryAndFragment) {
    return (hasTrailingSlash(input) ? input.slice(0, -1) : input) || "/";
  }
  if (!hasTrailingSlash(input, true)) {
    return input || "/";
  }
  let path = input;
  let fragment = "";
  const fragmentIndex = input.indexOf("#");
  if (fragmentIndex !== -1) {
    path = input.slice(0, fragmentIndex);
    fragment = input.slice(fragmentIndex);
  }
  const [s0, ...s] = path.split("?");
  const cleanPath = s0.endsWith("/") ? s0.slice(0, -1) : s0;
  return (cleanPath || "/") + (s.length > 0 ? `?${s.join("?")}` : "") + fragment;
}
function withTrailingSlash(input = "", respectQueryAndFragment) {
  if (!respectQueryAndFragment) {
    return input.endsWith("/") ? input : input + "/";
  }
  if (hasTrailingSlash(input, true)) {
    return input || "/";
  }
  let path = input;
  let fragment = "";
  const fragmentIndex = input.indexOf("#");
  if (fragmentIndex !== -1) {
    path = input.slice(0, fragmentIndex);
    fragment = input.slice(fragmentIndex);
    if (!path) {
      return fragment;
    }
  }
  const [s0, ...s] = path.split("?");
  return s0 + "/" + (s.length > 0 ? `?${s.join("?")}` : "") + fragment;
}
function hasLeadingSlash(input = "") {
  return input.startsWith("/");
}
function withLeadingSlash(input = "") {
  return hasLeadingSlash(input) ? input : "/" + input;
}
function withBase(input, base) {
  if (isEmptyURL(base) || hasProtocol(input)) {
    return input;
  }
  const _base = withoutTrailingSlash(base);
  if (input.startsWith(_base)) {
    return input;
  }
  return joinURL(_base, input);
}
function withoutBase(input, base) {
  if (isEmptyURL(base)) {
    return input;
  }
  const _base = withoutTrailingSlash(base);
  if (!input.startsWith(_base)) {
    return input;
  }
  const trimmed = input.slice(_base.length);
  return trimmed[0] === "/" ? trimmed : "/" + trimmed;
}
function withQuery(input, query) {
  const parsed = parseURL(input);
  const mergedQuery = { ...parseQuery(parsed.search), ...query };
  parsed.search = stringifyQuery(mergedQuery);
  return stringifyParsedURL(parsed);
}
function getQuery$1(input) {
  return parseQuery(parseURL(input).search);
}
function isEmptyURL(url) {
  return !url || url === "/";
}
function isNonEmptyURL(url) {
  return url && url !== "/";
}
function joinURL(base, ...input) {
  let url = base || "";
  for (const segment of input.filter((url2) => isNonEmptyURL(url2))) {
    if (url) {
      const _segment = segment.replace(JOIN_LEADING_SLASH_RE, "");
      url = withTrailingSlash(url) + _segment;
    } else {
      url = segment;
    }
  }
  return url;
}
function joinRelativeURL(..._input) {
  const JOIN_SEGMENT_SPLIT_RE = /\/(?!\/)/;
  const input = _input.filter(Boolean);
  const segments = [];
  let segmentsDepth = 0;
  for (const i of input) {
    if (!i || i === "/") {
      continue;
    }
    for (const [sindex, s] of i.split(JOIN_SEGMENT_SPLIT_RE).entries()) {
      if (!s || s === ".") {
        continue;
      }
      if (s === "..") {
        if (segments.length === 1 && hasProtocol(segments[0])) {
          continue;
        }
        segments.pop();
        segmentsDepth--;
        continue;
      }
      if (sindex === 1 && segments[segments.length - 1]?.endsWith(":/")) {
        segments[segments.length - 1] += "/" + s;
        continue;
      }
      segments.push(s);
      segmentsDepth++;
    }
  }
  let url = segments.join("/");
  if (segmentsDepth >= 0) {
    if (input[0]?.startsWith("/") && !url.startsWith("/")) {
      url = "/" + url;
    } else if (input[0]?.startsWith("./") && !url.startsWith("./")) {
      url = "./" + url;
    }
  } else {
    url = "../".repeat(-1 * segmentsDepth) + url;
  }
  if (input[input.length - 1]?.endsWith("/") && !url.endsWith("/")) {
    url += "/";
  }
  return url;
}

const protocolRelative = Symbol.for("ufo:protocolRelative");
function parseURL(input = "", defaultProto) {
  const _specialProtoMatch = input.match(
    /^[\s\0]*(blob:|data:|javascript:|vbscript:)(.*)/i
  );
  if (_specialProtoMatch) {
    const [, _proto, _pathname = ""] = _specialProtoMatch;
    return {
      protocol: _proto.toLowerCase(),
      pathname: _pathname,
      href: _proto + _pathname,
      auth: "",
      host: "",
      search: "",
      hash: ""
    };
  }
  if (!hasProtocol(input, { acceptRelative: true })) {
    return parsePath(input);
  }
  const [, protocol = "", auth, hostAndPath = ""] = input.replace(/\\/g, "/").match(/^[\s\0]*([\w+.-]{2,}:)?\/\/([^/@]+@)?(.*)/) || [];
  let [, host = "", path = ""] = hostAndPath.match(/([^#/?]*)(.*)?/) || [];
  if (protocol === "file:") {
    path = path.replace(/\/(?=[A-Za-z]:)/, "");
  }
  const { pathname, search, hash } = parsePath(path);
  return {
    protocol: protocol.toLowerCase(),
    auth: auth ? auth.slice(0, Math.max(0, auth.length - 1)) : "",
    host,
    pathname,
    search,
    hash,
    [protocolRelative]: !protocol
  };
}
function parsePath(input = "") {
  const [pathname = "", search = "", hash = ""] = (input.match(/([^#?]*)(\?[^#]*)?(#.*)?/) || []).splice(1);
  return {
    pathname,
    search,
    hash
  };
}
function stringifyParsedURL(parsed) {
  const pathname = parsed.pathname || "";
  const search = parsed.search ? (parsed.search.startsWith("?") ? "" : "?") + parsed.search : "";
  const hash = parsed.hash || "";
  const auth = parsed.auth ? parsed.auth + "@" : "";
  const host = parsed.host || "";
  const proto = parsed.protocol || parsed[protocolRelative] ? (parsed.protocol || "") + "//" : "";
  return proto + auth + host + pathname + search + hash;
}

function parse$1(str, options) {
  if (typeof str !== "string") {
    throw new TypeError("argument str must be a string");
  }
  const obj = {};
  const opt = {};
  const dec = opt.decode || decode$1;
  let index = 0;
  while (index < str.length) {
    const eqIdx = str.indexOf("=", index);
    if (eqIdx === -1) {
      break;
    }
    let endIdx = str.indexOf(";", index);
    if (endIdx === -1) {
      endIdx = str.length;
    } else if (endIdx < eqIdx) {
      index = str.lastIndexOf(";", eqIdx - 1) + 1;
      continue;
    }
    const key = str.slice(index, eqIdx).trim();
    if (opt?.filter && !opt?.filter(key)) {
      index = endIdx + 1;
      continue;
    }
    if (void 0 === obj[key]) {
      let val = str.slice(eqIdx + 1, endIdx).trim();
      if (val.codePointAt(0) === 34) {
        val = val.slice(1, -1);
      }
      obj[key] = tryDecode$1(val, dec);
    }
    index = endIdx + 1;
  }
  return obj;
}
function decode$1(str) {
  return str.includes("%") ? decodeURIComponent(str) : str;
}
function tryDecode$1(str, decode2) {
  try {
    return decode2(str);
  } catch {
    return str;
  }
}

const fieldContentRegExp = /^[\u0009\u0020-\u007E\u0080-\u00FF]+$/;
function serialize$2(name, value, options) {
  const opt = options || {};
  const enc = opt.encode || encodeURIComponent;
  if (typeof enc !== "function") {
    throw new TypeError("option encode is invalid");
  }
  if (!fieldContentRegExp.test(name)) {
    throw new TypeError("argument name is invalid");
  }
  const encodedValue = enc(value);
  if (encodedValue && !fieldContentRegExp.test(encodedValue)) {
    throw new TypeError("argument val is invalid");
  }
  let str = name + "=" + encodedValue;
  if (void 0 !== opt.maxAge && opt.maxAge !== null) {
    const maxAge = opt.maxAge - 0;
    if (Number.isNaN(maxAge) || !Number.isFinite(maxAge)) {
      throw new TypeError("option maxAge is invalid");
    }
    str += "; Max-Age=" + Math.floor(maxAge);
  }
  if (opt.domain) {
    if (!fieldContentRegExp.test(opt.domain)) {
      throw new TypeError("option domain is invalid");
    }
    str += "; Domain=" + opt.domain;
  }
  if (opt.path) {
    if (!fieldContentRegExp.test(opt.path)) {
      throw new TypeError("option path is invalid");
    }
    str += "; Path=" + opt.path;
  }
  if (opt.expires) {
    if (!isDate(opt.expires) || Number.isNaN(opt.expires.valueOf())) {
      throw new TypeError("option expires is invalid");
    }
    str += "; Expires=" + opt.expires.toUTCString();
  }
  if (opt.httpOnly) {
    str += "; HttpOnly";
  }
  if (opt.secure) {
    str += "; Secure";
  }
  if (opt.priority) {
    const priority = typeof opt.priority === "string" ? opt.priority.toLowerCase() : opt.priority;
    switch (priority) {
      case "low": {
        str += "; Priority=Low";
        break;
      }
      case "medium": {
        str += "; Priority=Medium";
        break;
      }
      case "high": {
        str += "; Priority=High";
        break;
      }
      default: {
        throw new TypeError("option priority is invalid");
      }
    }
  }
  if (opt.sameSite) {
    const sameSite = typeof opt.sameSite === "string" ? opt.sameSite.toLowerCase() : opt.sameSite;
    switch (sameSite) {
      case true: {
        str += "; SameSite=Strict";
        break;
      }
      case "lax": {
        str += "; SameSite=Lax";
        break;
      }
      case "strict": {
        str += "; SameSite=Strict";
        break;
      }
      case "none": {
        str += "; SameSite=None";
        break;
      }
      default: {
        throw new TypeError("option sameSite is invalid");
      }
    }
  }
  if (opt.partitioned) {
    str += "; Partitioned";
  }
  return str;
}
function isDate(val) {
  return Object.prototype.toString.call(val) === "[object Date]" || val instanceof Date;
}

function parseSetCookie(setCookieValue, options) {
  const parts = (setCookieValue || "").split(";").filter((str) => typeof str === "string" && !!str.trim());
  const nameValuePairStr = parts.shift() || "";
  const parsed = _parseNameValuePair(nameValuePairStr);
  const name = parsed.name;
  let value = parsed.value;
  try {
    value = options?.decode === false ? value : (options?.decode || decodeURIComponent)(value);
  } catch {
  }
  const cookie = {
    name,
    value
  };
  for (const part of parts) {
    const sides = part.split("=");
    const partKey = (sides.shift() || "").trimStart().toLowerCase();
    const partValue = sides.join("=");
    switch (partKey) {
      case "expires": {
        cookie.expires = new Date(partValue);
        break;
      }
      case "max-age": {
        cookie.maxAge = Number.parseInt(partValue, 10);
        break;
      }
      case "secure": {
        cookie.secure = true;
        break;
      }
      case "httponly": {
        cookie.httpOnly = true;
        break;
      }
      case "samesite": {
        cookie.sameSite = partValue;
        break;
      }
      default: {
        cookie[partKey] = partValue;
      }
    }
  }
  return cookie;
}
function _parseNameValuePair(nameValuePairStr) {
  let name = "";
  let value = "";
  const nameValueArr = nameValuePairStr.split("=");
  if (nameValueArr.length > 1) {
    name = nameValueArr.shift();
    value = nameValueArr.join("=");
  } else {
    value = nameValuePairStr;
  }
  return { name, value };
}

const NODE_TYPES = {
  NORMAL: 0,
  WILDCARD: 1,
  PLACEHOLDER: 2
};

function createRouter$1(options = {}) {
  const ctx = {
    options,
    rootNode: createRadixNode(),
    staticRoutesMap: {}
  };
  const normalizeTrailingSlash = (p) => options.strictTrailingSlash ? p : p.replace(/\/$/, "") || "/";
  if (options.routes) {
    for (const path in options.routes) {
      insert(ctx, normalizeTrailingSlash(path), options.routes[path]);
    }
  }
  return {
    ctx,
    lookup: (path) => lookup(ctx, normalizeTrailingSlash(path)),
    insert: (path, data) => insert(ctx, normalizeTrailingSlash(path), data),
    remove: (path) => remove(ctx, normalizeTrailingSlash(path))
  };
}
function lookup(ctx, path) {
  const staticPathNode = ctx.staticRoutesMap[path];
  if (staticPathNode) {
    return staticPathNode.data;
  }
  const sections = path.split("/");
  const params = {};
  let paramsFound = false;
  let wildcardNode = null;
  let node = ctx.rootNode;
  let wildCardParam = null;
  for (let i = 0; i < sections.length; i++) {
    const section = sections[i];
    if (node.wildcardChildNode !== null) {
      wildcardNode = node.wildcardChildNode;
      wildCardParam = sections.slice(i).join("/");
    }
    const nextNode = node.children.get(section);
    if (nextNode === void 0) {
      if (node && node.placeholderChildren.length > 1) {
        const remaining = sections.length - i;
        node = node.placeholderChildren.find((c) => c.maxDepth === remaining) || null;
      } else {
        node = node.placeholderChildren[0] || null;
      }
      if (!node) {
        break;
      }
      if (node.paramName) {
        params[node.paramName] = section;
      }
      paramsFound = true;
    } else {
      node = nextNode;
    }
  }
  if ((node === null || node.data === null) && wildcardNode !== null) {
    node = wildcardNode;
    params[node.paramName || "_"] = wildCardParam;
    paramsFound = true;
  }
  if (!node) {
    return null;
  }
  if (paramsFound) {
    return {
      ...node.data,
      params: paramsFound ? params : void 0
    };
  }
  return node.data;
}
function insert(ctx, path, data) {
  let isStaticRoute = true;
  const sections = path.split("/");
  let node = ctx.rootNode;
  let _unnamedPlaceholderCtr = 0;
  const matchedNodes = [node];
  for (const section of sections) {
    let childNode;
    if (childNode = node.children.get(section)) {
      node = childNode;
    } else {
      const type = getNodeType(section);
      childNode = createRadixNode({ type, parent: node });
      node.children.set(section, childNode);
      if (type === NODE_TYPES.PLACEHOLDER) {
        childNode.paramName = section === "*" ? `_${_unnamedPlaceholderCtr++}` : section.slice(1);
        node.placeholderChildren.push(childNode);
        isStaticRoute = false;
      } else if (type === NODE_TYPES.WILDCARD) {
        node.wildcardChildNode = childNode;
        childNode.paramName = section.slice(
          3
          /* "**:" */
        ) || "_";
        isStaticRoute = false;
      }
      matchedNodes.push(childNode);
      node = childNode;
    }
  }
  for (const [depth, node2] of matchedNodes.entries()) {
    node2.maxDepth = Math.max(matchedNodes.length - depth, node2.maxDepth || 0);
  }
  node.data = data;
  if (isStaticRoute === true) {
    ctx.staticRoutesMap[path] = node;
  }
  return node;
}
function remove(ctx, path) {
  let success = false;
  const sections = path.split("/");
  let node = ctx.rootNode;
  for (const section of sections) {
    node = node.children.get(section);
    if (!node) {
      return success;
    }
  }
  if (node.data) {
    const lastSection = sections.at(-1) || "";
    node.data = null;
    if (Object.keys(node.children).length === 0 && node.parent) {
      node.parent.children.delete(lastSection);
      node.parent.wildcardChildNode = null;
      node.parent.placeholderChildren = [];
    }
    success = true;
  }
  return success;
}
function createRadixNode(options = {}) {
  return {
    type: options.type || NODE_TYPES.NORMAL,
    maxDepth: 0,
    parent: options.parent || null,
    children: /* @__PURE__ */ new Map(),
    data: options.data || null,
    paramName: options.paramName || null,
    wildcardChildNode: null,
    placeholderChildren: []
  };
}
function getNodeType(str) {
  if (str.startsWith("**")) {
    return NODE_TYPES.WILDCARD;
  }
  if (str[0] === ":" || str === "*") {
    return NODE_TYPES.PLACEHOLDER;
  }
  return NODE_TYPES.NORMAL;
}

function toRouteMatcher(router) {
  const table = _routerNodeToTable("", router.ctx.rootNode);
  return _createMatcher(table, router.ctx.options.strictTrailingSlash);
}
function _createMatcher(table, strictTrailingSlash) {
  return {
    ctx: { table },
    matchAll: (path) => _matchRoutes(path, table, strictTrailingSlash)
  };
}
function _createRouteTable() {
  return {
    static: /* @__PURE__ */ new Map(),
    wildcard: /* @__PURE__ */ new Map(),
    dynamic: /* @__PURE__ */ new Map()
  };
}
function _matchRoutes(path, table, strictTrailingSlash) {
  if (strictTrailingSlash !== true && path.endsWith("/")) {
    path = path.slice(0, -1) || "/";
  }
  const matches = [];
  for (const [key, value] of _sortRoutesMap(table.wildcard)) {
    if (path === key || path.startsWith(key + "/")) {
      matches.push(value);
    }
  }
  for (const [key, value] of _sortRoutesMap(table.dynamic)) {
    if (path.startsWith(key + "/")) {
      const subPath = "/" + path.slice(key.length).split("/").splice(2).join("/");
      matches.push(..._matchRoutes(subPath, value));
    }
  }
  const staticMatch = table.static.get(path);
  if (staticMatch) {
    matches.push(staticMatch);
  }
  return matches.filter(Boolean);
}
function _sortRoutesMap(m) {
  return [...m.entries()].sort((a, b) => a[0].length - b[0].length);
}
function _routerNodeToTable(initialPath, initialNode) {
  const table = _createRouteTable();
  function _addNode(path, node) {
    if (path) {
      if (node.type === NODE_TYPES.NORMAL && !(path.includes("*") || path.includes(":"))) {
        if (node.data) {
          table.static.set(path, node.data);
        }
      } else if (node.type === NODE_TYPES.WILDCARD) {
        table.wildcard.set(path.replace("/**", ""), node.data);
      } else if (node.type === NODE_TYPES.PLACEHOLDER) {
        const subTable = _routerNodeToTable("", node);
        if (node.data) {
          subTable.static.set("/", node.data);
        }
        table.dynamic.set(path.replace(/\/\*|\/:\w+/, ""), subTable);
        return;
      }
    }
    for (const [childPath, child] of node.children.entries()) {
      _addNode(`${path}/${childPath}`.replace("//", "/"), child);
    }
  }
  _addNode(initialPath, initialNode);
  return table;
}

function isPlainObject(value) {
  if (value === null || typeof value !== "object") {
    return false;
  }
  const prototype = Object.getPrototypeOf(value);
  if (prototype !== null && prototype !== Object.prototype && Object.getPrototypeOf(prototype) !== null) {
    return false;
  }
  if (Symbol.iterator in value) {
    return false;
  }
  if (Symbol.toStringTag in value) {
    return Object.prototype.toString.call(value) === "[object Module]";
  }
  return true;
}

function _defu(baseObject, defaults, namespace = ".", merger) {
  if (!isPlainObject(defaults)) {
    return _defu(baseObject, {}, namespace, merger);
  }
  const object = Object.assign({}, defaults);
  for (const key in baseObject) {
    if (key === "__proto__" || key === "constructor") {
      continue;
    }
    const value = baseObject[key];
    if (value === null || value === void 0) {
      continue;
    }
    if (merger && merger(object, key, value, namespace)) {
      continue;
    }
    if (Array.isArray(value) && Array.isArray(object[key])) {
      object[key] = [...value, ...object[key]];
    } else if (isPlainObject(value) && isPlainObject(object[key])) {
      object[key] = _defu(
        value,
        object[key],
        (namespace ? `${namespace}.` : "") + key.toString(),
        merger
      );
    } else {
      object[key] = value;
    }
  }
  return object;
}
function createDefu(merger) {
  return (...arguments_) => (
    // eslint-disable-next-line unicorn/no-array-reduce
    arguments_.reduce((p, c) => _defu(p, c, "", merger), {})
  );
}
const defu = createDefu();
const defuFn = createDefu((object, key, currentValue) => {
  if (object[key] !== void 0 && typeof currentValue === "function") {
    object[key] = currentValue(object[key]);
    return true;
  }
});

function o(n){throw new Error(`${n} is not implemented yet!`)}let i$1 = class i extends EventEmitter{__unenv__={};readableEncoding=null;readableEnded=true;readableFlowing=false;readableHighWaterMark=0;readableLength=0;readableObjectMode=false;readableAborted=false;readableDidRead=false;closed=false;errored=null;readable=false;destroyed=false;static from(e,t){return new i(t)}constructor(e){super();}_read(e){}read(e){}setEncoding(e){return this}pause(){return this}resume(){return this}isPaused(){return  true}unpipe(e){return this}unshift(e,t){}wrap(e){return this}push(e,t){return  false}_destroy(e,t){this.removeAllListeners();}destroy(e){return this.destroyed=true,this._destroy(e),this}pipe(e,t){return {}}compose(e,t){throw new Error("Method not implemented.")}[Symbol.asyncDispose](){return this.destroy(),Promise.resolve()}async*[Symbol.asyncIterator](){throw o("Readable.asyncIterator")}iterator(e){throw o("Readable.iterator")}map(e,t){throw o("Readable.map")}filter(e,t){throw o("Readable.filter")}forEach(e,t){throw o("Readable.forEach")}reduce(e,t,r){throw o("Readable.reduce")}find(e,t){throw o("Readable.find")}findIndex(e,t){throw o("Readable.findIndex")}some(e,t){throw o("Readable.some")}toArray(e){throw o("Readable.toArray")}every(e,t){throw o("Readable.every")}flatMap(e,t){throw o("Readable.flatMap")}drop(e,t){throw o("Readable.drop")}take(e,t){throw o("Readable.take")}asIndexedPairs(e){throw o("Readable.asIndexedPairs")}};let l$1 = class l extends EventEmitter{__unenv__={};writable=true;writableEnded=false;writableFinished=false;writableHighWaterMark=0;writableLength=0;writableObjectMode=false;writableCorked=0;closed=false;errored=null;writableNeedDrain=false;writableAborted=false;destroyed=false;_data;_encoding="utf8";constructor(e){super();}pipe(e,t){return {}}_write(e,t,r){if(this.writableEnded){r&&r();return}if(this._data===void 0)this._data=e;else {const s=typeof this._data=="string"?Buffer$1.from(this._data,this._encoding||t||"utf8"):this._data,a=typeof e=="string"?Buffer$1.from(e,t||this._encoding||"utf8"):e;this._data=Buffer$1.concat([s,a]);}this._encoding=t,r&&r();}_writev(e,t){}_destroy(e,t){}_final(e){}write(e,t,r){const s=typeof t=="string"?this._encoding:"utf8",a=typeof t=="function"?t:typeof r=="function"?r:void 0;return this._write(e,s,a),true}setDefaultEncoding(e){return this}end(e,t,r){const s=typeof e=="function"?e:typeof t=="function"?t:typeof r=="function"?r:void 0;if(this.writableEnded)return s&&s(),this;const a=e===s?void 0:e;if(a){const u=t===s?void 0:t;this.write(a,u,s);}return this.writableEnded=true,this.writableFinished=true,this.emit("close"),this.emit("finish"),this}cork(){}uncork(){}destroy(e){return this.destroyed=true,delete this._data,this.removeAllListeners(),this}compose(e,t){throw new Error("Method not implemented.")}};const c$1=class c{allowHalfOpen=true;_destroy;constructor(e=new i$1,t=new l$1){Object.assign(this,e),Object.assign(this,t),this._destroy=g(e._destroy,t._destroy);}};function _(){return Object.assign(c$1.prototype,i$1.prototype),Object.assign(c$1.prototype,l$1.prototype),c$1}function g(...n){return function(...e){for(const t of n)t(...e);}}const m=_();class A extends m{__unenv__={};bufferSize=0;bytesRead=0;bytesWritten=0;connecting=false;destroyed=false;pending=false;localAddress="";localPort=0;remoteAddress="";remoteFamily="";remotePort=0;autoSelectFamilyAttemptedAddresses=[];readyState="readOnly";constructor(e){super();}write(e,t,r){return  false}connect(e,t,r){return this}end(e,t,r){return this}setEncoding(e){return this}pause(){return this}resume(){return this}setTimeout(e,t){return this}setNoDelay(e){return this}setKeepAlive(e,t){return this}address(){return {}}unref(){return this}ref(){return this}destroySoon(){this.destroy();}resetAndDestroy(){const e=new Error("ERR_SOCKET_CLOSED");return e.code="ERR_SOCKET_CLOSED",this.destroy(e),this}}class y extends i$1{aborted=false;httpVersion="1.1";httpVersionMajor=1;httpVersionMinor=1;complete=true;connection;socket;headers={};trailers={};method="GET";url="/";statusCode=200;statusMessage="";closed=false;errored=null;readable=false;constructor(e){super(),this.socket=this.connection=e||new A;}get rawHeaders(){const e=this.headers,t=[];for(const r in e)if(Array.isArray(e[r]))for(const s of e[r])t.push(r,s);else t.push(r,e[r]);return t}get rawTrailers(){return []}setTimeout(e,t){return this}get headersDistinct(){return p(this.headers)}get trailersDistinct(){return p(this.trailers)}}function p(n){const e={};for(const[t,r]of Object.entries(n))t&&(e[t]=(Array.isArray(r)?r:[r]).filter(Boolean));return e}class w extends l$1{statusCode=200;statusMessage="";upgrading=false;chunkedEncoding=false;shouldKeepAlive=false;useChunkedEncodingByDefault=false;sendDate=false;finished=false;headersSent=false;strictContentLength=false;connection=null;socket=null;req;_headers={};constructor(e){super(),this.req=e;}assignSocket(e){e._httpMessage=this,this.socket=e,this.connection=e,this.emit("socket",e),this._flush();}_flush(){this.flushHeaders();}detachSocket(e){}writeContinue(e){}writeHead(e,t,r){e&&(this.statusCode=e),typeof t=="string"&&(this.statusMessage=t,t=void 0);const s=r||t;if(s&&!Array.isArray(s))for(const a in s)this.setHeader(a,s[a]);return this.headersSent=true,this}writeProcessing(){}setTimeout(e,t){return this}appendHeader(e,t){e=e.toLowerCase();const r=this._headers[e],s=[...Array.isArray(r)?r:[r],...Array.isArray(t)?t:[t]].filter(Boolean);return this._headers[e]=s.length>1?s:s[0],this}setHeader(e,t){return this._headers[e.toLowerCase()]=t,this}setHeaders(e){for(const[t,r]of Object.entries(e))this.setHeader(t,r);return this}getHeader(e){return this._headers[e.toLowerCase()]}getHeaders(){return this._headers}getHeaderNames(){return Object.keys(this._headers)}hasHeader(e){return e.toLowerCase()in this._headers}removeHeader(e){delete this._headers[e.toLowerCase()];}addTrailers(e){}flushHeaders(){}writeEarlyHints(e,t){typeof t=="function"&&t();}}const E=(()=>{const n=function(){};return n.prototype=Object.create(null),n})();function R(n={}){const e=new E,t=Array.isArray(n)||H(n)?n:Object.entries(n);for(const[r,s]of t)if(s){if(e[r]===void 0){e[r]=s;continue}e[r]=[...Array.isArray(e[r])?e[r]:[e[r]],...Array.isArray(s)?s:[s]];}return e}function H(n){return typeof n?.entries=="function"}function v(n={}){if(n instanceof Headers)return n;const e=new Headers;for(const[t,r]of Object.entries(n))if(r!==void 0){if(Array.isArray(r)){for(const s of r)e.append(t,String(s));continue}e.set(t,String(r));}return e}const S=new Set([101,204,205,304]);async function b(n,e){const t=new y,r=new w(t);t.url=e.url?.toString()||"/";let s;if(!t.url.startsWith("/")){const d=new URL(t.url);s=d.host,t.url=d.pathname+d.search+d.hash;}t.method=e.method||"GET",t.headers=R(e.headers||{}),t.headers.host||(t.headers.host=e.host||s||"localhost"),t.connection.encrypted=t.connection.encrypted||e.protocol==="https",t.body=e.body||null,t.__unenv__=e.context,await n(t,r);let a=r._data;(S.has(r.statusCode)||t.method.toUpperCase()==="HEAD")&&(a=null,delete r._headers["content-length"]);const u={status:r.statusCode,statusText:r.statusMessage,headers:r._headers,body:a};return t.destroy(),r.destroy(),u}async function C(n,e,t={}){try{const r=await b(n,{url:e,...t});return new Response(r.body,{status:r.status,statusText:r.statusText,headers:v(r.headers)})}catch(r){return new Response(r.toString(),{status:Number.parseInt(r.statusCode||r.code)||500,statusText:r.statusText})}}

function hasProp(obj, prop) {
  try {
    return prop in obj;
  } catch {
    return false;
  }
}

class H3Error extends Error {
  static __h3_error__ = true;
  statusCode = 500;
  fatal = false;
  unhandled = false;
  statusMessage;
  data;
  cause;
  constructor(message, opts = {}) {
    super(message, opts);
    if (opts.cause && !this.cause) {
      this.cause = opts.cause;
    }
  }
  toJSON() {
    const obj = {
      message: this.message,
      statusCode: sanitizeStatusCode(this.statusCode, 500)
    };
    if (this.statusMessage) {
      obj.statusMessage = sanitizeStatusMessage(this.statusMessage);
    }
    if (this.data !== void 0) {
      obj.data = this.data;
    }
    return obj;
  }
}
function createError$1(input) {
  if (typeof input === "string") {
    return new H3Error(input);
  }
  if (isError(input)) {
    return input;
  }
  const err = new H3Error(input.message ?? input.statusMessage ?? "", {
    cause: input.cause || input
  });
  if (hasProp(input, "stack")) {
    try {
      Object.defineProperty(err, "stack", {
        get() {
          return input.stack;
        }
      });
    } catch {
      try {
        err.stack = input.stack;
      } catch {
      }
    }
  }
  if (input.data) {
    err.data = input.data;
  }
  if (input.statusCode) {
    err.statusCode = sanitizeStatusCode(input.statusCode, err.statusCode);
  } else if (input.status) {
    err.statusCode = sanitizeStatusCode(input.status, err.statusCode);
  }
  if (input.statusMessage) {
    err.statusMessage = input.statusMessage;
  } else if (input.statusText) {
    err.statusMessage = input.statusText;
  }
  if (err.statusMessage) {
    const originalMessage = err.statusMessage;
    const sanitizedMessage = sanitizeStatusMessage(err.statusMessage);
    if (sanitizedMessage !== originalMessage) {
      console.warn(
        "[h3] Please prefer using `message` for longer error messages instead of `statusMessage`. In the future, `statusMessage` will be sanitized by default."
      );
    }
  }
  if (input.fatal !== void 0) {
    err.fatal = input.fatal;
  }
  if (input.unhandled !== void 0) {
    err.unhandled = input.unhandled;
  }
  return err;
}
function sendError(event, error, debug) {
  if (event.handled) {
    return;
  }
  const h3Error = isError(error) ? error : createError$1(error);
  const responseBody = {
    statusCode: h3Error.statusCode,
    statusMessage: h3Error.statusMessage,
    stack: [],
    data: h3Error.data
  };
  if (debug) {
    responseBody.stack = (h3Error.stack || "").split("\n").map((l) => l.trim());
  }
  if (event.handled) {
    return;
  }
  const _code = Number.parseInt(h3Error.statusCode);
  setResponseStatus(event, _code, h3Error.statusMessage);
  event.node.res.setHeader("content-type", MIMES.json);
  event.node.res.end(JSON.stringify(responseBody, void 0, 2));
}
function isError(input) {
  return input?.constructor?.__h3_error__ === true;
}

function getQuery(event) {
  return getQuery$1(event.path || "");
}
function isMethod(event, expected, allowHead) {
  if (typeof expected === "string") {
    if (event.method === expected) {
      return true;
    }
  } else if (expected.includes(event.method)) {
    return true;
  }
  return false;
}
function assertMethod(event, expected, allowHead) {
  if (!isMethod(event, expected)) {
    throw createError$1({
      statusCode: 405,
      statusMessage: "HTTP method is not allowed."
    });
  }
}
function getRequestHeaders(event) {
  const _headers = {};
  for (const key in event.node.req.headers) {
    const val = event.node.req.headers[key];
    _headers[key] = Array.isArray(val) ? val.filter(Boolean).join(", ") : val;
  }
  return _headers;
}
function getRequestHeader(event, name) {
  const headers = getRequestHeaders(event);
  const value = headers[name.toLowerCase()];
  return value;
}
const getHeader = getRequestHeader;
function getRequestHost(event, opts = {}) {
  if (opts.xForwardedHost) {
    const xForwardedHost = event.node.req.headers["x-forwarded-host"];
    if (xForwardedHost) {
      return xForwardedHost;
    }
  }
  return event.node.req.headers.host || "localhost";
}
function getRequestProtocol(event, opts = {}) {
  if (opts.xForwardedProto !== false && event.node.req.headers["x-forwarded-proto"] === "https") {
    return "https";
  }
  return event.node.req.connection?.encrypted ? "https" : "http";
}
function getRequestURL(event, opts = {}) {
  const host = getRequestHost(event, opts);
  const protocol = getRequestProtocol(event, opts);
  const path = (event.node.req.originalUrl || event.path).replace(
    /^[/\\]+/g,
    "/"
  );
  return new URL(path, `${protocol}://${host}`);
}

const RawBodySymbol = Symbol.for("h3RawBody");
const PayloadMethods$1 = ["PATCH", "POST", "PUT", "DELETE"];
function readRawBody(event, encoding = "utf8") {
  assertMethod(event, PayloadMethods$1);
  const _rawBody = event._requestBody || event.web?.request?.body || event.node.req[RawBodySymbol] || event.node.req.rawBody || event.node.req.body;
  if (_rawBody) {
    const promise2 = Promise.resolve(_rawBody).then((_resolved) => {
      if (Buffer.isBuffer(_resolved)) {
        return _resolved;
      }
      if (typeof _resolved.pipeTo === "function") {
        return new Promise((resolve, reject) => {
          const chunks = [];
          _resolved.pipeTo(
            new WritableStream({
              write(chunk) {
                chunks.push(chunk);
              },
              close() {
                resolve(Buffer.concat(chunks));
              },
              abort(reason) {
                reject(reason);
              }
            })
          ).catch(reject);
        });
      } else if (typeof _resolved.pipe === "function") {
        return new Promise((resolve, reject) => {
          const chunks = [];
          _resolved.on("data", (chunk) => {
            chunks.push(chunk);
          }).on("end", () => {
            resolve(Buffer.concat(chunks));
          }).on("error", reject);
        });
      }
      if (_resolved.constructor === Object) {
        return Buffer.from(JSON.stringify(_resolved));
      }
      if (_resolved instanceof URLSearchParams) {
        return Buffer.from(_resolved.toString());
      }
      if (_resolved instanceof FormData) {
        return new Response(_resolved).bytes().then((uint8arr) => Buffer.from(uint8arr));
      }
      return Buffer.from(_resolved);
    });
    return encoding ? promise2.then((buff) => buff.toString(encoding)) : promise2;
  }
  if (!Number.parseInt(event.node.req.headers["content-length"] || "") && !String(event.node.req.headers["transfer-encoding"] ?? "").split(",").map((e) => e.trim()).filter(Boolean).includes("chunked")) {
    return Promise.resolve(void 0);
  }
  const promise = event.node.req[RawBodySymbol] = new Promise(
    (resolve, reject) => {
      const bodyData = [];
      event.node.req.on("error", (err) => {
        reject(err);
      }).on("data", (chunk) => {
        bodyData.push(chunk);
      }).on("end", () => {
        resolve(Buffer.concat(bodyData));
      });
    }
  );
  const result = encoding ? promise.then((buff) => buff.toString(encoding)) : promise;
  return result;
}
function getRequestWebStream(event) {
  if (!PayloadMethods$1.includes(event.method)) {
    return;
  }
  const bodyStream = event.web?.request?.body || event._requestBody;
  if (bodyStream) {
    return bodyStream;
  }
  const _hasRawBody = RawBodySymbol in event.node.req || "rawBody" in event.node.req || "body" in event.node.req || "__unenv__" in event.node.req;
  if (_hasRawBody) {
    return new ReadableStream({
      async start(controller) {
        const _rawBody = await readRawBody(event, false);
        if (_rawBody) {
          controller.enqueue(_rawBody);
        }
        controller.close();
      }
    });
  }
  return new ReadableStream({
    start: (controller) => {
      event.node.req.on("data", (chunk) => {
        controller.enqueue(chunk);
      });
      event.node.req.on("end", () => {
        controller.close();
      });
      event.node.req.on("error", (err) => {
        controller.error(err);
      });
    }
  });
}

function handleCacheHeaders(event, opts) {
  const cacheControls = ["public", ...opts.cacheControls || []];
  let cacheMatched = false;
  if (opts.maxAge !== void 0) {
    cacheControls.push(`max-age=${+opts.maxAge}`, `s-maxage=${+opts.maxAge}`);
  }
  if (opts.modifiedTime) {
    const modifiedTime = new Date(opts.modifiedTime);
    const ifModifiedSince = event.node.req.headers["if-modified-since"];
    event.node.res.setHeader("last-modified", modifiedTime.toUTCString());
    if (ifModifiedSince && new Date(ifModifiedSince) >= modifiedTime) {
      cacheMatched = true;
    }
  }
  if (opts.etag) {
    event.node.res.setHeader("etag", opts.etag);
    const ifNonMatch = event.node.req.headers["if-none-match"];
    if (ifNonMatch === opts.etag) {
      cacheMatched = true;
    }
  }
  event.node.res.setHeader("cache-control", cacheControls.join(", "));
  if (cacheMatched) {
    event.node.res.statusCode = 304;
    if (!event.handled) {
      event.node.res.end();
    }
    return true;
  }
  return false;
}

const MIMES = {
  html: "text/html",
  json: "application/json"
};

const DISALLOWED_STATUS_CHARS = /[^\u0009\u0020-\u007E]/g;
function sanitizeStatusMessage(statusMessage = "") {
  return statusMessage.replace(DISALLOWED_STATUS_CHARS, "");
}
function sanitizeStatusCode(statusCode, defaultStatusCode = 200) {
  if (!statusCode) {
    return defaultStatusCode;
  }
  if (typeof statusCode === "string") {
    statusCode = Number.parseInt(statusCode, 10);
  }
  if (statusCode < 100 || statusCode > 999) {
    return defaultStatusCode;
  }
  return statusCode;
}

function getDistinctCookieKey(name, opts) {
  return [name, opts.domain || "", opts.path || "/"].join(";");
}

function parseCookies(event) {
  return parse$1(event.node.req.headers.cookie || "");
}
function getCookie(event, name) {
  return parseCookies(event)[name];
}
function setCookie(event, name, value, serializeOptions = {}) {
  if (!serializeOptions.path) {
    serializeOptions = { path: "/", ...serializeOptions };
  }
  const newCookie = serialize$2(name, value, serializeOptions);
  const currentCookies = splitCookiesString(
    event.node.res.getHeader("set-cookie")
  );
  if (currentCookies.length === 0) {
    event.node.res.setHeader("set-cookie", newCookie);
    return;
  }
  const newCookieKey = getDistinctCookieKey(name, serializeOptions);
  event.node.res.removeHeader("set-cookie");
  for (const cookie of currentCookies) {
    const parsed = parseSetCookie(cookie);
    const key = getDistinctCookieKey(parsed.name, parsed);
    if (key === newCookieKey) {
      continue;
    }
    event.node.res.appendHeader("set-cookie", cookie);
  }
  event.node.res.appendHeader("set-cookie", newCookie);
}
function deleteCookie(event, name, serializeOptions) {
  setCookie(event, name, "", {
    ...serializeOptions,
    maxAge: 0
  });
}
function splitCookiesString(cookiesString) {
  if (Array.isArray(cookiesString)) {
    return cookiesString.flatMap((c) => splitCookiesString(c));
  }
  if (typeof cookiesString !== "string") {
    return [];
  }
  const cookiesStrings = [];
  let pos = 0;
  let start;
  let ch;
  let lastComma;
  let nextStart;
  let cookiesSeparatorFound;
  const skipWhitespace = () => {
    while (pos < cookiesString.length && /\s/.test(cookiesString.charAt(pos))) {
      pos += 1;
    }
    return pos < cookiesString.length;
  };
  const notSpecialChar = () => {
    ch = cookiesString.charAt(pos);
    return ch !== "=" && ch !== ";" && ch !== ",";
  };
  while (pos < cookiesString.length) {
    start = pos;
    cookiesSeparatorFound = false;
    while (skipWhitespace()) {
      ch = cookiesString.charAt(pos);
      if (ch === ",") {
        lastComma = pos;
        pos += 1;
        skipWhitespace();
        nextStart = pos;
        while (pos < cookiesString.length && notSpecialChar()) {
          pos += 1;
        }
        if (pos < cookiesString.length && cookiesString.charAt(pos) === "=") {
          cookiesSeparatorFound = true;
          pos = nextStart;
          cookiesStrings.push(cookiesString.slice(start, lastComma));
          start = pos;
        } else {
          pos = lastComma + 1;
        }
      } else {
        pos += 1;
      }
    }
    if (!cookiesSeparatorFound || pos >= cookiesString.length) {
      cookiesStrings.push(cookiesString.slice(start));
    }
  }
  return cookiesStrings;
}

const defer = typeof setImmediate === "undefined" ? (fn) => fn() : setImmediate;
function send(event, data, type) {
  if (type) {
    defaultContentType(event, type);
  }
  return new Promise((resolve) => {
    defer(() => {
      if (!event.handled) {
        event.node.res.end(data);
      }
      resolve();
    });
  });
}
function sendNoContent(event, code) {
  if (event.handled) {
    return;
  }
  if (!code && event.node.res.statusCode !== 200) {
    code = event.node.res.statusCode;
  }
  const _code = sanitizeStatusCode(code, 204);
  if (_code === 204) {
    event.node.res.removeHeader("content-length");
  }
  event.node.res.writeHead(_code);
  event.node.res.end();
}
function setResponseStatus(event, code, text) {
  if (code) {
    event.node.res.statusCode = sanitizeStatusCode(
      code,
      event.node.res.statusCode
    );
  }
  if (text) {
    event.node.res.statusMessage = sanitizeStatusMessage(text);
  }
}
function getResponseStatus(event) {
  return event.node.res.statusCode;
}
function getResponseStatusText(event) {
  return event.node.res.statusMessage;
}
function defaultContentType(event, type) {
  if (type && event.node.res.statusCode !== 304 && !event.node.res.getHeader("content-type")) {
    event.node.res.setHeader("content-type", type);
  }
}
function sendRedirect(event, location, code = 302) {
  event.node.res.statusCode = sanitizeStatusCode(
    code,
    event.node.res.statusCode
  );
  event.node.res.setHeader("location", location);
  const encodedLoc = location.replace(/"/g, "%22");
  const html = `<!DOCTYPE html><html><head><meta http-equiv="refresh" content="0; url=${encodedLoc}"></head></html>`;
  return send(event, html, MIMES.html);
}
function getResponseHeader(event, name) {
  return event.node.res.getHeader(name);
}
function setResponseHeaders(event, headers) {
  for (const [name, value] of Object.entries(headers)) {
    event.node.res.setHeader(
      name,
      value
    );
  }
}
const setHeaders = setResponseHeaders;
function setResponseHeader(event, name, value) {
  event.node.res.setHeader(name, value);
}
function appendResponseHeader(event, name, value) {
  let current = event.node.res.getHeader(name);
  if (!current) {
    event.node.res.setHeader(name, value);
    return;
  }
  if (!Array.isArray(current)) {
    current = [current.toString()];
  }
  event.node.res.setHeader(name, [...current, value]);
}
function removeResponseHeader(event, name) {
  return event.node.res.removeHeader(name);
}
function isStream(data) {
  if (!data || typeof data !== "object") {
    return false;
  }
  if (typeof data.pipe === "function") {
    if (typeof data._read === "function") {
      return true;
    }
    if (typeof data.abort === "function") {
      return true;
    }
  }
  if (typeof data.pipeTo === "function") {
    return true;
  }
  return false;
}
function isWebResponse(data) {
  return typeof Response !== "undefined" && data instanceof Response;
}
function sendStream(event, stream) {
  if (!stream || typeof stream !== "object") {
    throw new Error("[h3] Invalid stream provided.");
  }
  event.node.res._data = stream;
  if (!event.node.res.socket) {
    event._handled = true;
    return Promise.resolve();
  }
  if (hasProp(stream, "pipeTo") && typeof stream.pipeTo === "function") {
    return stream.pipeTo(
      new WritableStream({
        write(chunk) {
          event.node.res.write(chunk);
        }
      })
    ).then(() => {
      event.node.res.end();
    });
  }
  if (hasProp(stream, "pipe") && typeof stream.pipe === "function") {
    return new Promise((resolve, reject) => {
      stream.pipe(event.node.res);
      if (stream.on) {
        stream.on("end", () => {
          event.node.res.end();
          resolve();
        });
        stream.on("error", (error) => {
          reject(error);
        });
      }
      event.node.res.on("close", () => {
        if (stream.abort) {
          stream.abort();
        }
      });
    });
  }
  throw new Error("[h3] Invalid or incompatible stream provided.");
}
function sendWebResponse(event, response) {
  for (const [key, value] of response.headers) {
    if (key === "set-cookie") {
      event.node.res.appendHeader(key, splitCookiesString(value));
    } else {
      event.node.res.setHeader(key, value);
    }
  }
  if (response.status) {
    event.node.res.statusCode = sanitizeStatusCode(
      response.status,
      event.node.res.statusCode
    );
  }
  if (response.statusText) {
    event.node.res.statusMessage = sanitizeStatusMessage(response.statusText);
  }
  if (response.redirected) {
    event.node.res.setHeader("location", response.url);
  }
  if (!response.body) {
    event.node.res.end();
    return;
  }
  return sendStream(event, response.body);
}

const PayloadMethods = /* @__PURE__ */ new Set(["PATCH", "POST", "PUT", "DELETE"]);
const ignoredHeaders = /* @__PURE__ */ new Set([
  "transfer-encoding",
  "accept-encoding",
  "connection",
  "keep-alive",
  "upgrade",
  "expect",
  "host",
  "accept"
]);
async function proxyRequest(event, target, opts = {}) {
  let body;
  let duplex;
  if (PayloadMethods.has(event.method)) {
    if (opts.streamRequest) {
      body = getRequestWebStream(event);
      duplex = "half";
    } else {
      body = await readRawBody(event, false).catch(() => void 0);
    }
  }
  const method = opts.fetchOptions?.method || event.method;
  const fetchHeaders = mergeHeaders$1(
    getProxyRequestHeaders(event, { host: target.startsWith("/") }),
    opts.fetchOptions?.headers,
    opts.headers
  );
  return sendProxy(event, target, {
    ...opts,
    fetchOptions: {
      method,
      body,
      duplex,
      ...opts.fetchOptions,
      headers: fetchHeaders
    }
  });
}
async function sendProxy(event, target, opts = {}) {
  let response;
  try {
    response = await _getFetch(opts.fetch)(target, {
      headers: opts.headers,
      ignoreResponseError: true,
      // make $ofetch.raw transparent
      ...opts.fetchOptions
    });
  } catch (error) {
    throw createError$1({
      status: 502,
      statusMessage: "Bad Gateway",
      cause: error
    });
  }
  event.node.res.statusCode = sanitizeStatusCode(
    response.status,
    event.node.res.statusCode
  );
  event.node.res.statusMessage = sanitizeStatusMessage(response.statusText);
  const cookies = [];
  for (const [key, value] of response.headers.entries()) {
    if (key === "content-encoding") {
      continue;
    }
    if (key === "content-length") {
      continue;
    }
    if (key === "set-cookie") {
      cookies.push(...splitCookiesString(value));
      continue;
    }
    event.node.res.setHeader(key, value);
  }
  if (cookies.length > 0) {
    event.node.res.setHeader(
      "set-cookie",
      cookies.map((cookie) => {
        if (opts.cookieDomainRewrite) {
          cookie = rewriteCookieProperty(
            cookie,
            opts.cookieDomainRewrite,
            "domain"
          );
        }
        if (opts.cookiePathRewrite) {
          cookie = rewriteCookieProperty(
            cookie,
            opts.cookiePathRewrite,
            "path"
          );
        }
        return cookie;
      })
    );
  }
  if (opts.onResponse) {
    await opts.onResponse(event, response);
  }
  if (response._data !== void 0) {
    return response._data;
  }
  if (event.handled) {
    return;
  }
  if (opts.sendStream === false) {
    const data = new Uint8Array(await response.arrayBuffer());
    return event.node.res.end(data);
  }
  if (response.body) {
    for await (const chunk of response.body) {
      event.node.res.write(chunk);
    }
  }
  return event.node.res.end();
}
function getProxyRequestHeaders(event, opts) {
  const headers = /* @__PURE__ */ Object.create(null);
  const reqHeaders = getRequestHeaders(event);
  for (const name in reqHeaders) {
    if (!ignoredHeaders.has(name) || name === "host" && opts?.host) {
      headers[name] = reqHeaders[name];
    }
  }
  return headers;
}
function fetchWithEvent(event, req, init, options) {
  return _getFetch(options?.fetch)(req, {
    ...init,
    context: init?.context || event.context,
    headers: {
      ...getProxyRequestHeaders(event, {
        host: typeof req === "string" && req.startsWith("/")
      }),
      ...init?.headers
    }
  });
}
function _getFetch(_fetch) {
  if (_fetch) {
    return _fetch;
  }
  if (globalThis.fetch) {
    return globalThis.fetch;
  }
  throw new Error(
    "fetch is not available. Try importing `node-fetch-native/polyfill` for Node.js."
  );
}
function rewriteCookieProperty(header, map, property) {
  const _map = typeof map === "string" ? { "*": map } : map;
  return header.replace(
    new RegExp(`(;\\s*${property}=)([^;]+)`, "gi"),
    (match, prefix, previousValue) => {
      let newValue;
      if (previousValue in _map) {
        newValue = _map[previousValue];
      } else if ("*" in _map) {
        newValue = _map["*"];
      } else {
        return match;
      }
      return newValue ? prefix + newValue : "";
    }
  );
}
function mergeHeaders$1(defaults, ...inputs) {
  const _inputs = inputs.filter(Boolean);
  if (_inputs.length === 0) {
    return defaults;
  }
  const merged = new Headers(defaults);
  for (const input of _inputs) {
    const entries = Array.isArray(input) ? input : typeof input.entries === "function" ? input.entries() : Object.entries(input);
    for (const [key, value] of entries) {
      if (value !== void 0) {
        merged.set(key, value);
      }
    }
  }
  return merged;
}

class H3Event {
  "__is_event__" = true;
  // Context
  node;
  // Node
  web;
  // Web
  context = {};
  // Shared
  // Request
  _method;
  _path;
  _headers;
  _requestBody;
  // Response
  _handled = false;
  // Hooks
  _onBeforeResponseCalled;
  _onAfterResponseCalled;
  constructor(req, res) {
    this.node = { req, res };
  }
  // --- Request ---
  get method() {
    if (!this._method) {
      this._method = (this.node.req.method || "GET").toUpperCase();
    }
    return this._method;
  }
  get path() {
    return this._path || this.node.req.url || "/";
  }
  get headers() {
    if (!this._headers) {
      this._headers = _normalizeNodeHeaders(this.node.req.headers);
    }
    return this._headers;
  }
  // --- Respoonse ---
  get handled() {
    return this._handled || this.node.res.writableEnded || this.node.res.headersSent;
  }
  respondWith(response) {
    return Promise.resolve(response).then(
      (_response) => sendWebResponse(this, _response)
    );
  }
  // --- Utils ---
  toString() {
    return `[${this.method}] ${this.path}`;
  }
  toJSON() {
    return this.toString();
  }
  // --- Deprecated ---
  /** @deprecated Please use `event.node.req` instead. */
  get req() {
    return this.node.req;
  }
  /** @deprecated Please use `event.node.res` instead. */
  get res() {
    return this.node.res;
  }
}
function isEvent(input) {
  return hasProp(input, "__is_event__");
}
function createEvent(req, res) {
  return new H3Event(req, res);
}
function _normalizeNodeHeaders(nodeHeaders) {
  const headers = new Headers();
  for (const [name, value] of Object.entries(nodeHeaders)) {
    if (Array.isArray(value)) {
      for (const item of value) {
        headers.append(name, item);
      }
    } else if (value) {
      headers.set(name, value);
    }
  }
  return headers;
}

function defineEventHandler(handler) {
  if (typeof handler === "function") {
    handler.__is_handler__ = true;
    return handler;
  }
  const _hooks = {
    onRequest: _normalizeArray(handler.onRequest),
    onBeforeResponse: _normalizeArray(handler.onBeforeResponse)
  };
  const _handler = (event) => {
    return _callHandler(event, handler.handler, _hooks);
  };
  _handler.__is_handler__ = true;
  _handler.__resolve__ = handler.handler.__resolve__;
  _handler.__websocket__ = handler.websocket;
  return _handler;
}
function _normalizeArray(input) {
  return input ? Array.isArray(input) ? input : [input] : void 0;
}
async function _callHandler(event, handler, hooks) {
  if (hooks.onRequest) {
    for (const hook of hooks.onRequest) {
      await hook(event);
      if (event.handled) {
        return;
      }
    }
  }
  const body = await handler(event);
  const response = { body };
  if (hooks.onBeforeResponse) {
    for (const hook of hooks.onBeforeResponse) {
      await hook(event, response);
    }
  }
  return response.body;
}
const eventHandler = defineEventHandler;
function isEventHandler(input) {
  return hasProp(input, "__is_handler__");
}
function toEventHandler(input, _, _route) {
  if (!isEventHandler(input)) {
    console.warn(
      "[h3] Implicit event handler conversion is deprecated. Use `eventHandler()` or `fromNodeMiddleware()` to define event handlers.",
      _route && _route !== "/" ? `
     Route: ${_route}` : "",
      `
     Handler: ${input}`
    );
  }
  return input;
}
function defineLazyEventHandler(factory) {
  let _promise;
  let _resolved;
  const resolveHandler = () => {
    if (_resolved) {
      return Promise.resolve(_resolved);
    }
    if (!_promise) {
      _promise = Promise.resolve(factory()).then((r) => {
        const handler2 = r.default || r;
        if (typeof handler2 !== "function") {
          throw new TypeError(
            "Invalid lazy handler result. It should be a function:",
            handler2
          );
        }
        _resolved = { handler: toEventHandler(r.default || r) };
        return _resolved;
      });
    }
    return _promise;
  };
  const handler = eventHandler((event) => {
    if (_resolved) {
      return _resolved.handler(event);
    }
    return resolveHandler().then((r) => r.handler(event));
  });
  handler.__resolve__ = resolveHandler;
  return handler;
}
const lazyEventHandler = defineLazyEventHandler;

function createApp(options = {}) {
  const stack = [];
  const handler = createAppEventHandler(stack, options);
  const resolve = createResolver(stack);
  handler.__resolve__ = resolve;
  const getWebsocket = cachedFn(() => websocketOptions(resolve, options));
  const app = {
    // @ts-expect-error
    use: (arg1, arg2, arg3) => use(app, arg1, arg2, arg3),
    resolve,
    handler,
    stack,
    options,
    get websocket() {
      return getWebsocket();
    }
  };
  return app;
}
function use(app, arg1, arg2, arg3) {
  if (Array.isArray(arg1)) {
    for (const i of arg1) {
      use(app, i, arg2, arg3);
    }
  } else if (Array.isArray(arg2)) {
    for (const i of arg2) {
      use(app, arg1, i, arg3);
    }
  } else if (typeof arg1 === "string") {
    app.stack.push(
      normalizeLayer({ ...arg3, route: arg1, handler: arg2 })
    );
  } else if (typeof arg1 === "function") {
    app.stack.push(normalizeLayer({ ...arg2, handler: arg1 }));
  } else {
    app.stack.push(normalizeLayer({ ...arg1 }));
  }
  return app;
}
function createAppEventHandler(stack, options) {
  const spacing = options.debug ? 2 : void 0;
  return eventHandler(async (event) => {
    event.node.req.originalUrl = event.node.req.originalUrl || event.node.req.url || "/";
    const _reqPath = event._path || event.node.req.url || "/";
    let _layerPath;
    if (options.onRequest) {
      await options.onRequest(event);
    }
    for (const layer of stack) {
      if (layer.route.length > 1) {
        if (!_reqPath.startsWith(layer.route)) {
          continue;
        }
        _layerPath = _reqPath.slice(layer.route.length) || "/";
      } else {
        _layerPath = _reqPath;
      }
      if (layer.match && !layer.match(_layerPath, event)) {
        continue;
      }
      event._path = _layerPath;
      event.node.req.url = _layerPath;
      const val = await layer.handler(event);
      const _body = val === void 0 ? void 0 : await val;
      if (_body !== void 0) {
        const _response = { body: _body };
        if (options.onBeforeResponse) {
          event._onBeforeResponseCalled = true;
          await options.onBeforeResponse(event, _response);
        }
        await handleHandlerResponse(event, _response.body, spacing);
        if (options.onAfterResponse) {
          event._onAfterResponseCalled = true;
          await options.onAfterResponse(event, _response);
        }
        return;
      }
      if (event.handled) {
        if (options.onAfterResponse) {
          event._onAfterResponseCalled = true;
          await options.onAfterResponse(event, void 0);
        }
        return;
      }
    }
    if (!event.handled) {
      throw createError$1({
        statusCode: 404,
        statusMessage: `Cannot find any path matching ${event.path || "/"}.`
      });
    }
    if (options.onAfterResponse) {
      event._onAfterResponseCalled = true;
      await options.onAfterResponse(event, void 0);
    }
  });
}
function createResolver(stack) {
  return async (path) => {
    let _layerPath;
    for (const layer of stack) {
      if (layer.route === "/" && !layer.handler.__resolve__) {
        continue;
      }
      if (!path.startsWith(layer.route)) {
        continue;
      }
      _layerPath = path.slice(layer.route.length) || "/";
      if (layer.match && !layer.match(_layerPath, void 0)) {
        continue;
      }
      let res = { route: layer.route, handler: layer.handler };
      if (res.handler.__resolve__) {
        const _res = await res.handler.__resolve__(_layerPath);
        if (!_res) {
          continue;
        }
        res = {
          ...res,
          ..._res,
          route: joinURL(res.route || "/", _res.route || "/")
        };
      }
      return res;
    }
  };
}
function normalizeLayer(input) {
  let handler = input.handler;
  if (handler.handler) {
    handler = handler.handler;
  }
  if (input.lazy) {
    handler = lazyEventHandler(handler);
  } else if (!isEventHandler(handler)) {
    handler = toEventHandler(handler, void 0, input.route);
  }
  return {
    route: withoutTrailingSlash(input.route),
    match: input.match,
    handler
  };
}
function handleHandlerResponse(event, val, jsonSpace) {
  if (val === null) {
    return sendNoContent(event);
  }
  if (val) {
    if (isWebResponse(val)) {
      return sendWebResponse(event, val);
    }
    if (isStream(val)) {
      return sendStream(event, val);
    }
    if (val.buffer) {
      return send(event, val);
    }
    if (val.arrayBuffer && typeof val.arrayBuffer === "function") {
      return val.arrayBuffer().then((arrayBuffer) => {
        return send(event, Buffer.from(arrayBuffer), val.type);
      });
    }
    if (val instanceof Error) {
      throw createError$1(val);
    }
    if (typeof val.end === "function") {
      return true;
    }
  }
  const valType = typeof val;
  if (valType === "string") {
    return send(event, val, MIMES.html);
  }
  if (valType === "object" || valType === "boolean" || valType === "number") {
    return send(event, JSON.stringify(val, void 0, jsonSpace), MIMES.json);
  }
  if (valType === "bigint") {
    return send(event, val.toString(), MIMES.json);
  }
  throw createError$1({
    statusCode: 500,
    statusMessage: `[h3] Cannot send ${valType} as response.`
  });
}
function cachedFn(fn) {
  let cache;
  return () => {
    if (!cache) {
      cache = fn();
    }
    return cache;
  };
}
function websocketOptions(evResolver, appOptions) {
  return {
    ...appOptions.websocket,
    async resolve(info) {
      const url = info.request?.url || info.url || "/";
      const { pathname } = typeof url === "string" ? parseURL(url) : url;
      const resolved = await evResolver(pathname);
      return resolved?.handler?.__websocket__ || {};
    }
  };
}

const RouterMethods = [
  "connect",
  "delete",
  "get",
  "head",
  "options",
  "post",
  "put",
  "trace",
  "patch"
];
function createRouter(opts = {}) {
  const _router = createRouter$1({});
  const routes = {};
  let _matcher;
  const router = {};
  const addRoute = (path, handler, method) => {
    let route = routes[path];
    if (!route) {
      routes[path] = route = { path, handlers: {} };
      _router.insert(path, route);
    }
    if (Array.isArray(method)) {
      for (const m of method) {
        addRoute(path, handler, m);
      }
    } else {
      route.handlers[method] = toEventHandler(handler, void 0, path);
    }
    return router;
  };
  router.use = router.add = (path, handler, method) => addRoute(path, handler, method || "all");
  for (const method of RouterMethods) {
    router[method] = (path, handle) => router.add(path, handle, method);
  }
  const matchHandler = (path = "/", method = "get") => {
    const qIndex = path.indexOf("?");
    if (qIndex !== -1) {
      path = path.slice(0, Math.max(0, qIndex));
    }
    const matched = _router.lookup(path);
    if (!matched || !matched.handlers) {
      return {
        error: createError$1({
          statusCode: 404,
          name: "Not Found",
          statusMessage: `Cannot find any route matching ${path || "/"}.`
        })
      };
    }
    let handler = matched.handlers[method] || matched.handlers.all;
    if (!handler) {
      if (!_matcher) {
        _matcher = toRouteMatcher(_router);
      }
      const _matches = _matcher.matchAll(path).reverse();
      for (const _match of _matches) {
        if (_match.handlers[method]) {
          handler = _match.handlers[method];
          matched.handlers[method] = matched.handlers[method] || handler;
          break;
        }
        if (_match.handlers.all) {
          handler = _match.handlers.all;
          matched.handlers.all = matched.handlers.all || handler;
          break;
        }
      }
    }
    if (!handler) {
      return {
        error: createError$1({
          statusCode: 405,
          name: "Method Not Allowed",
          statusMessage: `Method ${method} is not allowed on this route.`
        })
      };
    }
    return { matched, handler };
  };
  const isPreemptive = opts.preemptive || opts.preemtive;
  router.handler = eventHandler((event) => {
    const match = matchHandler(
      event.path,
      event.method.toLowerCase()
    );
    if ("error" in match) {
      if (isPreemptive) {
        throw match.error;
      } else {
        return;
      }
    }
    event.context.matchedRoute = match.matched;
    const params = match.matched.params || {};
    event.context.params = params;
    return Promise.resolve(match.handler(event)).then((res) => {
      if (res === void 0 && isPreemptive) {
        return null;
      }
      return res;
    });
  });
  router.handler.__resolve__ = async (path) => {
    path = withLeadingSlash(path);
    const match = matchHandler(path);
    if ("error" in match) {
      return;
    }
    let res = {
      route: match.matched.path,
      handler: match.handler
    };
    if (match.handler.__resolve__) {
      const _res = await match.handler.__resolve__(path);
      if (!_res) {
        return;
      }
      res = { ...res, ..._res };
    }
    return res;
  };
  return router;
}
function toNodeListener(app) {
  const toNodeHandle = async function(req, res) {
    const event = createEvent(req, res);
    try {
      await app.handler(event);
    } catch (_error) {
      const error = createError$1(_error);
      if (!isError(_error)) {
        error.unhandled = true;
      }
      setResponseStatus(event, error.statusCode, error.statusMessage);
      if (app.options.onError) {
        await app.options.onError(error, event);
      }
      if (event.handled) {
        return;
      }
      if (error.unhandled || error.fatal) {
        console.error("[h3]", error.fatal ? "[fatal]" : "[unhandled]", error);
      }
      if (app.options.onBeforeResponse && !event._onBeforeResponseCalled) {
        await app.options.onBeforeResponse(event, { body: error });
      }
      await sendError(event, error, !!app.options.debug);
      if (app.options.onAfterResponse && !event._onAfterResponseCalled) {
        await app.options.onAfterResponse(event, { body: error });
      }
    }
  };
  return toNodeHandle;
}

function flatHooks(configHooks, hooks = {}, parentName) {
  for (const key in configHooks) {
    const subHook = configHooks[key];
    const name = parentName ? `${parentName}:${key}` : key;
    if (typeof subHook === "object" && subHook !== null) {
      flatHooks(subHook, hooks, name);
    } else if (typeof subHook === "function") {
      hooks[name] = subHook;
    }
  }
  return hooks;
}
const defaultTask = { run: (function_) => function_() };
const _createTask = () => defaultTask;
const createTask = typeof console.createTask !== "undefined" ? console.createTask : _createTask;
function serialTaskCaller(hooks, args) {
  const name = args.shift();
  const task = createTask(name);
  return hooks.reduce(
    (promise, hookFunction) => promise.then(() => task.run(() => hookFunction(...args))),
    Promise.resolve()
  );
}
function parallelTaskCaller(hooks, args) {
  const name = args.shift();
  const task = createTask(name);
  return Promise.all(hooks.map((hook) => task.run(() => hook(...args))));
}
function callEachWith(callbacks, arg0) {
  for (const callback of [...callbacks]) {
    callback(arg0);
  }
}

class Hookable {
  constructor() {
    this._hooks = {};
    this._before = void 0;
    this._after = void 0;
    this._deprecatedMessages = void 0;
    this._deprecatedHooks = {};
    this.hook = this.hook.bind(this);
    this.callHook = this.callHook.bind(this);
    this.callHookWith = this.callHookWith.bind(this);
  }
  hook(name, function_, options = {}) {
    if (!name || typeof function_ !== "function") {
      return () => {
      };
    }
    const originalName = name;
    let dep;
    while (this._deprecatedHooks[name]) {
      dep = this._deprecatedHooks[name];
      name = dep.to;
    }
    if (dep && !options.allowDeprecated) {
      let message = dep.message;
      if (!message) {
        message = `${originalName} hook has been deprecated` + (dep.to ? `, please use ${dep.to}` : "");
      }
      if (!this._deprecatedMessages) {
        this._deprecatedMessages = /* @__PURE__ */ new Set();
      }
      if (!this._deprecatedMessages.has(message)) {
        console.warn(message);
        this._deprecatedMessages.add(message);
      }
    }
    if (!function_.name) {
      try {
        Object.defineProperty(function_, "name", {
          get: () => "_" + name.replace(/\W+/g, "_") + "_hook_cb",
          configurable: true
        });
      } catch {
      }
    }
    this._hooks[name] = this._hooks[name] || [];
    this._hooks[name].push(function_);
    return () => {
      if (function_) {
        this.removeHook(name, function_);
        function_ = void 0;
      }
    };
  }
  hookOnce(name, function_) {
    let _unreg;
    let _function = (...arguments_) => {
      if (typeof _unreg === "function") {
        _unreg();
      }
      _unreg = void 0;
      _function = void 0;
      return function_(...arguments_);
    };
    _unreg = this.hook(name, _function);
    return _unreg;
  }
  removeHook(name, function_) {
    if (this._hooks[name]) {
      const index = this._hooks[name].indexOf(function_);
      if (index !== -1) {
        this._hooks[name].splice(index, 1);
      }
      if (this._hooks[name].length === 0) {
        delete this._hooks[name];
      }
    }
  }
  deprecateHook(name, deprecated) {
    this._deprecatedHooks[name] = typeof deprecated === "string" ? { to: deprecated } : deprecated;
    const _hooks = this._hooks[name] || [];
    delete this._hooks[name];
    for (const hook of _hooks) {
      this.hook(name, hook);
    }
  }
  deprecateHooks(deprecatedHooks) {
    Object.assign(this._deprecatedHooks, deprecatedHooks);
    for (const name in deprecatedHooks) {
      this.deprecateHook(name, deprecatedHooks[name]);
    }
  }
  addHooks(configHooks) {
    const hooks = flatHooks(configHooks);
    const removeFns = Object.keys(hooks).map(
      (key) => this.hook(key, hooks[key])
    );
    return () => {
      for (const unreg of removeFns.splice(0, removeFns.length)) {
        unreg();
      }
    };
  }
  removeHooks(configHooks) {
    const hooks = flatHooks(configHooks);
    for (const key in hooks) {
      this.removeHook(key, hooks[key]);
    }
  }
  removeAllHooks() {
    for (const key in this._hooks) {
      delete this._hooks[key];
    }
  }
  callHook(name, ...arguments_) {
    arguments_.unshift(name);
    return this.callHookWith(serialTaskCaller, name, ...arguments_);
  }
  callHookParallel(name, ...arguments_) {
    arguments_.unshift(name);
    return this.callHookWith(parallelTaskCaller, name, ...arguments_);
  }
  callHookWith(caller, name, ...arguments_) {
    const event = this._before || this._after ? { name, args: arguments_, context: {} } : void 0;
    if (this._before) {
      callEachWith(this._before, event);
    }
    const result = caller(
      name in this._hooks ? [...this._hooks[name]] : [],
      arguments_
    );
    if (result instanceof Promise) {
      return result.finally(() => {
        if (this._after && event) {
          callEachWith(this._after, event);
        }
      });
    }
    if (this._after && event) {
      callEachWith(this._after, event);
    }
    return result;
  }
  beforeEach(function_) {
    this._before = this._before || [];
    this._before.push(function_);
    return () => {
      if (this._before !== void 0) {
        const index = this._before.indexOf(function_);
        if (index !== -1) {
          this._before.splice(index, 1);
        }
      }
    };
  }
  afterEach(function_) {
    this._after = this._after || [];
    this._after.push(function_);
    return () => {
      if (this._after !== void 0) {
        const index = this._after.indexOf(function_);
        if (index !== -1) {
          this._after.splice(index, 1);
        }
      }
    };
  }
}
function createHooks() {
  return new Hookable();
}

const s$1=globalThis.Headers,i=globalThis.AbortController,l=globalThis.fetch||(()=>{throw new Error("[node-fetch-native] Failed to fetch: `globalThis.fetch` is not available!")});

class FetchError extends Error {
  constructor(message, opts) {
    super(message, opts);
    this.name = "FetchError";
    if (opts?.cause && !this.cause) {
      this.cause = opts.cause;
    }
  }
}
function createFetchError(ctx) {
  const errorMessage = ctx.error?.message || ctx.error?.toString() || "";
  const method = ctx.request?.method || ctx.options?.method || "GET";
  const url = ctx.request?.url || String(ctx.request) || "/";
  const requestStr = `[${method}] ${JSON.stringify(url)}`;
  const statusStr = ctx.response ? `${ctx.response.status} ${ctx.response.statusText}` : "<no response>";
  const message = `${requestStr}: ${statusStr}${errorMessage ? ` ${errorMessage}` : ""}`;
  const fetchError = new FetchError(
    message,
    ctx.error ? { cause: ctx.error } : void 0
  );
  for (const key of ["request", "options", "response"]) {
    Object.defineProperty(fetchError, key, {
      get() {
        return ctx[key];
      }
    });
  }
  for (const [key, refKey] of [
    ["data", "_data"],
    ["status", "status"],
    ["statusCode", "status"],
    ["statusText", "statusText"],
    ["statusMessage", "statusText"]
  ]) {
    Object.defineProperty(fetchError, key, {
      get() {
        return ctx.response && ctx.response[refKey];
      }
    });
  }
  return fetchError;
}

const payloadMethods = new Set(
  Object.freeze(["PATCH", "POST", "PUT", "DELETE"])
);
function isPayloadMethod(method = "GET") {
  return payloadMethods.has(method.toUpperCase());
}
function isJSONSerializable(value) {
  if (value === void 0) {
    return false;
  }
  const t = typeof value;
  if (t === "string" || t === "number" || t === "boolean" || t === null) {
    return true;
  }
  if (t !== "object") {
    return false;
  }
  if (Array.isArray(value)) {
    return true;
  }
  if (value.buffer) {
    return false;
  }
  return value.constructor && value.constructor.name === "Object" || typeof value.toJSON === "function";
}
const textTypes = /* @__PURE__ */ new Set([
  "image/svg",
  "application/xml",
  "application/xhtml",
  "application/html"
]);
const JSON_RE = /^application\/(?:[\w!#$%&*.^`~-]*\+)?json(;.+)?$/i;
function detectResponseType(_contentType = "") {
  if (!_contentType) {
    return "json";
  }
  const contentType = _contentType.split(";").shift() || "";
  if (JSON_RE.test(contentType)) {
    return "json";
  }
  if (textTypes.has(contentType) || contentType.startsWith("text/")) {
    return "text";
  }
  return "blob";
}
function resolveFetchOptions(request, input, defaults, Headers) {
  const headers = mergeHeaders(
    input?.headers ?? request?.headers,
    defaults?.headers,
    Headers
  );
  let query;
  if (defaults?.query || defaults?.params || input?.params || input?.query) {
    query = {
      ...defaults?.params,
      ...defaults?.query,
      ...input?.params,
      ...input?.query
    };
  }
  return {
    ...defaults,
    ...input,
    query,
    params: query,
    headers
  };
}
function mergeHeaders(input, defaults, Headers) {
  if (!defaults) {
    return new Headers(input);
  }
  const headers = new Headers(defaults);
  if (input) {
    for (const [key, value] of Symbol.iterator in input || Array.isArray(input) ? input : new Headers(input)) {
      headers.set(key, value);
    }
  }
  return headers;
}
async function callHooks(context, hooks) {
  if (hooks) {
    if (Array.isArray(hooks)) {
      for (const hook of hooks) {
        await hook(context);
      }
    } else {
      await hooks(context);
    }
  }
}

const retryStatusCodes = /* @__PURE__ */ new Set([
  408,
  // Request Timeout
  409,
  // Conflict
  425,
  // Too Early (Experimental)
  429,
  // Too Many Requests
  500,
  // Internal Server Error
  502,
  // Bad Gateway
  503,
  // Service Unavailable
  504
  // Gateway Timeout
]);
const nullBodyResponses = /* @__PURE__ */ new Set([101, 204, 205, 304]);
function createFetch(globalOptions = {}) {
  const {
    fetch = globalThis.fetch,
    Headers = globalThis.Headers,
    AbortController = globalThis.AbortController
  } = globalOptions;
  async function onError(context) {
    const isAbort = context.error && context.error.name === "AbortError" && !context.options.timeout || false;
    if (context.options.retry !== false && !isAbort) {
      let retries;
      if (typeof context.options.retry === "number") {
        retries = context.options.retry;
      } else {
        retries = isPayloadMethod(context.options.method) ? 0 : 1;
      }
      const responseCode = context.response && context.response.status || 500;
      if (retries > 0 && (Array.isArray(context.options.retryStatusCodes) ? context.options.retryStatusCodes.includes(responseCode) : retryStatusCodes.has(responseCode))) {
        const retryDelay = typeof context.options.retryDelay === "function" ? context.options.retryDelay(context) : context.options.retryDelay || 0;
        if (retryDelay > 0) {
          await new Promise((resolve) => setTimeout(resolve, retryDelay));
        }
        return $fetchRaw(context.request, {
          ...context.options,
          retry: retries - 1
        });
      }
    }
    const error = createFetchError(context);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(error, $fetchRaw);
    }
    throw error;
  }
  const $fetchRaw = async function $fetchRaw2(_request, _options = {}) {
    const context = {
      request: _request,
      options: resolveFetchOptions(
        _request,
        _options,
        globalOptions.defaults,
        Headers
      ),
      response: void 0,
      error: void 0
    };
    if (context.options.method) {
      context.options.method = context.options.method.toUpperCase();
    }
    if (context.options.onRequest) {
      await callHooks(context, context.options.onRequest);
    }
    if (typeof context.request === "string") {
      if (context.options.baseURL) {
        context.request = withBase(context.request, context.options.baseURL);
      }
      if (context.options.query) {
        context.request = withQuery(context.request, context.options.query);
        delete context.options.query;
      }
      if ("query" in context.options) {
        delete context.options.query;
      }
      if ("params" in context.options) {
        delete context.options.params;
      }
    }
    if (context.options.body && isPayloadMethod(context.options.method)) {
      if (isJSONSerializable(context.options.body)) {
        context.options.body = typeof context.options.body === "string" ? context.options.body : JSON.stringify(context.options.body);
        context.options.headers = new Headers(context.options.headers || {});
        if (!context.options.headers.has("content-type")) {
          context.options.headers.set("content-type", "application/json");
        }
        if (!context.options.headers.has("accept")) {
          context.options.headers.set("accept", "application/json");
        }
      } else if (
        // ReadableStream Body
        "pipeTo" in context.options.body && typeof context.options.body.pipeTo === "function" || // Node.js Stream Body
        typeof context.options.body.pipe === "function"
      ) {
        if (!("duplex" in context.options)) {
          context.options.duplex = "half";
        }
      }
    }
    let abortTimeout;
    if (!context.options.signal && context.options.timeout) {
      const controller = new AbortController();
      abortTimeout = setTimeout(() => {
        const error = new Error(
          "[TimeoutError]: The operation was aborted due to timeout"
        );
        error.name = "TimeoutError";
        error.code = 23;
        controller.abort(error);
      }, context.options.timeout);
      context.options.signal = controller.signal;
    }
    try {
      context.response = await fetch(
        context.request,
        context.options
      );
    } catch (error) {
      context.error = error;
      if (context.options.onRequestError) {
        await callHooks(
          context,
          context.options.onRequestError
        );
      }
      return await onError(context);
    } finally {
      if (abortTimeout) {
        clearTimeout(abortTimeout);
      }
    }
    const hasBody = (context.response.body || // https://github.com/unjs/ofetch/issues/324
    // https://github.com/unjs/ofetch/issues/294
    // https://github.com/JakeChampion/fetch/issues/1454
    context.response._bodyInit) && !nullBodyResponses.has(context.response.status) && context.options.method !== "HEAD";
    if (hasBody) {
      const responseType = (context.options.parseResponse ? "json" : context.options.responseType) || detectResponseType(context.response.headers.get("content-type") || "");
      switch (responseType) {
        case "json": {
          const data = await context.response.text();
          const parseFunction = context.options.parseResponse || destr;
          context.response._data = parseFunction(data);
          break;
        }
        case "stream": {
          context.response._data = context.response.body || context.response._bodyInit;
          break;
        }
        default: {
          context.response._data = await context.response[responseType]();
        }
      }
    }
    if (context.options.onResponse) {
      await callHooks(
        context,
        context.options.onResponse
      );
    }
    if (!context.options.ignoreResponseError && context.response.status >= 400 && context.response.status < 600) {
      if (context.options.onResponseError) {
        await callHooks(
          context,
          context.options.onResponseError
        );
      }
      return await onError(context);
    }
    return context.response;
  };
  const $fetch = async function $fetch2(request, options) {
    const r = await $fetchRaw(request, options);
    return r._data;
  };
  $fetch.raw = $fetchRaw;
  $fetch.native = (...args) => fetch(...args);
  $fetch.create = (defaultOptions = {}, customGlobalOptions = {}) => createFetch({
    ...globalOptions,
    ...customGlobalOptions,
    defaults: {
      ...globalOptions.defaults,
      ...customGlobalOptions.defaults,
      ...defaultOptions
    }
  });
  return $fetch;
}

function createNodeFetch() {
  const useKeepAlive = JSON.parse(process.env.FETCH_KEEP_ALIVE || "false");
  if (!useKeepAlive) {
    return l;
  }
  const agentOptions = { keepAlive: true };
  const httpAgent = new ft.Agent(agentOptions);
  const httpsAgent = new Qa.Agent(agentOptions);
  const nodeFetchOptions = {
    agent(parsedURL) {
      return parsedURL.protocol === "http:" ? httpAgent : httpsAgent;
    }
  };
  return function nodeFetchWithKeepAlive(input, init) {
    return l(input, { ...nodeFetchOptions, ...init });
  };
}
const fetch = globalThis.fetch ? (...args) => globalThis.fetch(...args) : createNodeFetch();
const Headers$1 = globalThis.Headers || s$1;
const AbortController = globalThis.AbortController || i;
createFetch({ fetch, Headers: Headers$1, AbortController });

function wrapToPromise(value) {
  if (!value || typeof value.then !== "function") {
    return Promise.resolve(value);
  }
  return value;
}
function asyncCall(function_, ...arguments_) {
  try {
    return wrapToPromise(function_(...arguments_));
  } catch (error) {
    return Promise.reject(error);
  }
}
function isPrimitive(value) {
  const type = typeof value;
  return value === null || type !== "object" && type !== "function";
}
function isPureObject(value) {
  const proto = Object.getPrototypeOf(value);
  return !proto || proto.isPrototypeOf(Object);
}
function stringify(value) {
  if (isPrimitive(value)) {
    return String(value);
  }
  if (isPureObject(value) || Array.isArray(value)) {
    return JSON.stringify(value);
  }
  if (typeof value.toJSON === "function") {
    return stringify(value.toJSON());
  }
  throw new Error("[unstorage] Cannot stringify value!");
}
const BASE64_PREFIX = "base64:";
function serializeRaw(value) {
  if (typeof value === "string") {
    return value;
  }
  return BASE64_PREFIX + base64Encode(value);
}
function deserializeRaw(value) {
  if (typeof value !== "string") {
    return value;
  }
  if (!value.startsWith(BASE64_PREFIX)) {
    return value;
  }
  return base64Decode(value.slice(BASE64_PREFIX.length));
}
function base64Decode(input) {
  if (globalThis.Buffer) {
    return Buffer.from(input, "base64");
  }
  return Uint8Array.from(
    globalThis.atob(input),
    (c) => c.codePointAt(0)
  );
}
function base64Encode(input) {
  if (globalThis.Buffer) {
    return Buffer.from(input).toString("base64");
  }
  return globalThis.btoa(String.fromCodePoint(...input));
}

const storageKeyProperties = [
  "has",
  "hasItem",
  "get",
  "getItem",
  "getItemRaw",
  "set",
  "setItem",
  "setItemRaw",
  "del",
  "remove",
  "removeItem",
  "getMeta",
  "setMeta",
  "removeMeta",
  "getKeys",
  "clear",
  "mount",
  "unmount"
];
function prefixStorage(storage, base) {
  base = normalizeBaseKey(base);
  if (!base) {
    return storage;
  }
  const nsStorage = { ...storage };
  for (const property of storageKeyProperties) {
    nsStorage[property] = (key = "", ...args) => (
      // @ts-ignore
      storage[property](base + key, ...args)
    );
  }
  nsStorage.getKeys = (key = "", ...arguments_) => storage.getKeys(base + key, ...arguments_).then((keys) => keys.map((key2) => key2.slice(base.length)));
  nsStorage.getItems = async (items, commonOptions) => {
    const prefixedItems = items.map(
      (item) => typeof item === "string" ? base + item : { ...item, key: base + item.key }
    );
    const results = await storage.getItems(prefixedItems, commonOptions);
    return results.map((entry) => ({
      key: entry.key.slice(base.length),
      value: entry.value
    }));
  };
  nsStorage.setItems = async (items, commonOptions) => {
    const prefixedItems = items.map((item) => ({
      key: base + item.key,
      value: item.value,
      options: item.options
    }));
    return storage.setItems(prefixedItems, commonOptions);
  };
  return nsStorage;
}
function normalizeKey$1(key) {
  if (!key) {
    return "";
  }
  return key.split("?")[0]?.replace(/[/\\]/g, ":").replace(/:+/g, ":").replace(/^:|:$/g, "") || "";
}
function joinKeys(...keys) {
  return normalizeKey$1(keys.join(":"));
}
function normalizeBaseKey(base) {
  base = normalizeKey$1(base);
  return base ? base + ":" : "";
}
function filterKeyByDepth(key, depth) {
  if (depth === void 0) {
    return true;
  }
  let substrCount = 0;
  let index = key.indexOf(":");
  while (index > -1) {
    substrCount++;
    index = key.indexOf(":", index + 1);
  }
  return substrCount <= depth;
}
function filterKeyByBase(key, base) {
  if (base) {
    return key.startsWith(base) && key[key.length - 1] !== "$";
  }
  return key[key.length - 1] !== "$";
}

function defineDriver$1(factory) {
  return factory;
}

const DRIVER_NAME$2 = "memory";
const memory = defineDriver$1(() => {
  const data = /* @__PURE__ */ new Map();
  return {
    name: DRIVER_NAME$2,
    getInstance: () => data,
    hasItem(key) {
      return data.has(key);
    },
    getItem(key) {
      return data.get(key) ?? null;
    },
    getItemRaw(key) {
      return data.get(key) ?? null;
    },
    setItem(key, value) {
      data.set(key, value);
    },
    setItemRaw(key, value) {
      data.set(key, value);
    },
    removeItem(key) {
      data.delete(key);
    },
    getKeys() {
      return [...data.keys()];
    },
    clear() {
      data.clear();
    },
    dispose() {
      data.clear();
    }
  };
});

function createStorage(options = {}) {
  const context = {
    mounts: { "": options.driver || memory() },
    mountpoints: [""],
    watching: false,
    watchListeners: [],
    unwatch: {}
  };
  const getMount = (key) => {
    for (const base of context.mountpoints) {
      if (key.startsWith(base)) {
        return {
          base,
          relativeKey: key.slice(base.length),
          driver: context.mounts[base]
        };
      }
    }
    return {
      base: "",
      relativeKey: key,
      driver: context.mounts[""]
    };
  };
  const getMounts = (base, includeParent) => {
    return context.mountpoints.filter(
      (mountpoint) => mountpoint.startsWith(base) || includeParent && base.startsWith(mountpoint)
    ).map((mountpoint) => ({
      relativeBase: base.length > mountpoint.length ? base.slice(mountpoint.length) : void 0,
      mountpoint,
      driver: context.mounts[mountpoint]
    }));
  };
  const onChange = (event, key) => {
    if (!context.watching) {
      return;
    }
    key = normalizeKey$1(key);
    for (const listener of context.watchListeners) {
      listener(event, key);
    }
  };
  const startWatch = async () => {
    if (context.watching) {
      return;
    }
    context.watching = true;
    for (const mountpoint in context.mounts) {
      context.unwatch[mountpoint] = await watch(
        context.mounts[mountpoint],
        onChange,
        mountpoint
      );
    }
  };
  const stopWatch = async () => {
    if (!context.watching) {
      return;
    }
    for (const mountpoint in context.unwatch) {
      await context.unwatch[mountpoint]();
    }
    context.unwatch = {};
    context.watching = false;
  };
  const runBatch = (items, commonOptions, cb) => {
    const batches = /* @__PURE__ */ new Map();
    const getBatch = (mount) => {
      let batch = batches.get(mount.base);
      if (!batch) {
        batch = {
          driver: mount.driver,
          base: mount.base,
          items: []
        };
        batches.set(mount.base, batch);
      }
      return batch;
    };
    for (const item of items) {
      const isStringItem = typeof item === "string";
      const key = normalizeKey$1(isStringItem ? item : item.key);
      const value = isStringItem ? void 0 : item.value;
      const options2 = isStringItem || !item.options ? commonOptions : { ...commonOptions, ...item.options };
      const mount = getMount(key);
      getBatch(mount).items.push({
        key,
        value,
        relativeKey: mount.relativeKey,
        options: options2
      });
    }
    return Promise.all([...batches.values()].map((batch) => cb(batch))).then(
      (r) => r.flat()
    );
  };
  const storage = {
    // Item
    hasItem(key, opts = {}) {
      key = normalizeKey$1(key);
      const { relativeKey, driver } = getMount(key);
      return asyncCall(driver.hasItem, relativeKey, opts);
    },
    getItem(key, opts = {}) {
      key = normalizeKey$1(key);
      const { relativeKey, driver } = getMount(key);
      return asyncCall(driver.getItem, relativeKey, opts).then(
        (value) => destr(value)
      );
    },
    getItems(items, commonOptions = {}) {
      return runBatch(items, commonOptions, (batch) => {
        if (batch.driver.getItems) {
          return asyncCall(
            batch.driver.getItems,
            batch.items.map((item) => ({
              key: item.relativeKey,
              options: item.options
            })),
            commonOptions
          ).then(
            (r) => r.map((item) => ({
              key: joinKeys(batch.base, item.key),
              value: destr(item.value)
            }))
          );
        }
        return Promise.all(
          batch.items.map((item) => {
            return asyncCall(
              batch.driver.getItem,
              item.relativeKey,
              item.options
            ).then((value) => ({
              key: item.key,
              value: destr(value)
            }));
          })
        );
      });
    },
    getItemRaw(key, opts = {}) {
      key = normalizeKey$1(key);
      const { relativeKey, driver } = getMount(key);
      if (driver.getItemRaw) {
        return asyncCall(driver.getItemRaw, relativeKey, opts);
      }
      return asyncCall(driver.getItem, relativeKey, opts).then(
        (value) => deserializeRaw(value)
      );
    },
    async setItem(key, value, opts = {}) {
      if (value === void 0) {
        return storage.removeItem(key);
      }
      key = normalizeKey$1(key);
      const { relativeKey, driver } = getMount(key);
      if (!driver.setItem) {
        return;
      }
      await asyncCall(driver.setItem, relativeKey, stringify(value), opts);
      if (!driver.watch) {
        onChange("update", key);
      }
    },
    async setItems(items, commonOptions) {
      await runBatch(items, commonOptions, async (batch) => {
        if (batch.driver.setItems) {
          return asyncCall(
            batch.driver.setItems,
            batch.items.map((item) => ({
              key: item.relativeKey,
              value: stringify(item.value),
              options: item.options
            })),
            commonOptions
          );
        }
        if (!batch.driver.setItem) {
          return;
        }
        await Promise.all(
          batch.items.map((item) => {
            return asyncCall(
              batch.driver.setItem,
              item.relativeKey,
              stringify(item.value),
              item.options
            );
          })
        );
      });
    },
    async setItemRaw(key, value, opts = {}) {
      if (value === void 0) {
        return storage.removeItem(key, opts);
      }
      key = normalizeKey$1(key);
      const { relativeKey, driver } = getMount(key);
      if (driver.setItemRaw) {
        await asyncCall(driver.setItemRaw, relativeKey, value, opts);
      } else if (driver.setItem) {
        await asyncCall(driver.setItem, relativeKey, serializeRaw(value), opts);
      } else {
        return;
      }
      if (!driver.watch) {
        onChange("update", key);
      }
    },
    async removeItem(key, opts = {}) {
      if (typeof opts === "boolean") {
        opts = { removeMeta: opts };
      }
      key = normalizeKey$1(key);
      const { relativeKey, driver } = getMount(key);
      if (!driver.removeItem) {
        return;
      }
      await asyncCall(driver.removeItem, relativeKey, opts);
      if (opts.removeMeta || opts.removeMata) {
        await asyncCall(driver.removeItem, relativeKey + "$", opts);
      }
      if (!driver.watch) {
        onChange("remove", key);
      }
    },
    // Meta
    async getMeta(key, opts = {}) {
      if (typeof opts === "boolean") {
        opts = { nativeOnly: opts };
      }
      key = normalizeKey$1(key);
      const { relativeKey, driver } = getMount(key);
      const meta = /* @__PURE__ */ Object.create(null);
      if (driver.getMeta) {
        Object.assign(meta, await asyncCall(driver.getMeta, relativeKey, opts));
      }
      if (!opts.nativeOnly) {
        const value = await asyncCall(
          driver.getItem,
          relativeKey + "$",
          opts
        ).then((value_) => destr(value_));
        if (value && typeof value === "object") {
          if (typeof value.atime === "string") {
            value.atime = new Date(value.atime);
          }
          if (typeof value.mtime === "string") {
            value.mtime = new Date(value.mtime);
          }
          Object.assign(meta, value);
        }
      }
      return meta;
    },
    setMeta(key, value, opts = {}) {
      return this.setItem(key + "$", value, opts);
    },
    removeMeta(key, opts = {}) {
      return this.removeItem(key + "$", opts);
    },
    // Keys
    async getKeys(base, opts = {}) {
      base = normalizeBaseKey(base);
      const mounts = getMounts(base, true);
      let maskedMounts = [];
      const allKeys = [];
      let allMountsSupportMaxDepth = true;
      for (const mount of mounts) {
        if (!mount.driver.flags?.maxDepth) {
          allMountsSupportMaxDepth = false;
        }
        const rawKeys = await asyncCall(
          mount.driver.getKeys,
          mount.relativeBase,
          opts
        );
        for (const key of rawKeys) {
          const fullKey = mount.mountpoint + normalizeKey$1(key);
          if (!maskedMounts.some((p) => fullKey.startsWith(p))) {
            allKeys.push(fullKey);
          }
        }
        maskedMounts = [
          mount.mountpoint,
          ...maskedMounts.filter((p) => !p.startsWith(mount.mountpoint))
        ];
      }
      const shouldFilterByDepth = opts.maxDepth !== void 0 && !allMountsSupportMaxDepth;
      return allKeys.filter(
        (key) => (!shouldFilterByDepth || filterKeyByDepth(key, opts.maxDepth)) && filterKeyByBase(key, base)
      );
    },
    // Utils
    async clear(base, opts = {}) {
      base = normalizeBaseKey(base);
      await Promise.all(
        getMounts(base, false).map(async (m) => {
          if (m.driver.clear) {
            return asyncCall(m.driver.clear, m.relativeBase, opts);
          }
          if (m.driver.removeItem) {
            const keys = await m.driver.getKeys(m.relativeBase || "", opts);
            return Promise.all(
              keys.map((key) => m.driver.removeItem(key, opts))
            );
          }
        })
      );
    },
    async dispose() {
      await Promise.all(
        Object.values(context.mounts).map((driver) => dispose(driver))
      );
    },
    async watch(callback) {
      await startWatch();
      context.watchListeners.push(callback);
      return async () => {
        context.watchListeners = context.watchListeners.filter(
          (listener) => listener !== callback
        );
        if (context.watchListeners.length === 0) {
          await stopWatch();
        }
      };
    },
    async unwatch() {
      context.watchListeners = [];
      await stopWatch();
    },
    // Mount
    mount(base, driver) {
      base = normalizeBaseKey(base);
      if (base && context.mounts[base]) {
        throw new Error(`already mounted at ${base}`);
      }
      if (base) {
        context.mountpoints.push(base);
        context.mountpoints.sort((a, b) => b.length - a.length);
      }
      context.mounts[base] = driver;
      if (context.watching) {
        Promise.resolve(watch(driver, onChange, base)).then((unwatcher) => {
          context.unwatch[base] = unwatcher;
        }).catch(console.error);
      }
      return storage;
    },
    async unmount(base, _dispose = true) {
      base = normalizeBaseKey(base);
      if (!base || !context.mounts[base]) {
        return;
      }
      if (context.watching && base in context.unwatch) {
        context.unwatch[base]?.();
        delete context.unwatch[base];
      }
      if (_dispose) {
        await dispose(context.mounts[base]);
      }
      context.mountpoints = context.mountpoints.filter((key) => key !== base);
      delete context.mounts[base];
    },
    getMount(key = "") {
      key = normalizeKey$1(key) + ":";
      const m = getMount(key);
      return {
        driver: m.driver,
        base: m.base
      };
    },
    getMounts(base = "", opts = {}) {
      base = normalizeKey$1(base);
      const mounts = getMounts(base, opts.parents);
      return mounts.map((m) => ({
        driver: m.driver,
        base: m.mountpoint
      }));
    },
    // Aliases
    keys: (base, opts = {}) => storage.getKeys(base, opts),
    get: (key, opts = {}) => storage.getItem(key, opts),
    set: (key, value, opts = {}) => storage.setItem(key, value, opts),
    has: (key, opts = {}) => storage.hasItem(key, opts),
    del: (key, opts = {}) => storage.removeItem(key, opts),
    remove: (key, opts = {}) => storage.removeItem(key, opts)
  };
  return storage;
}
function watch(driver, onChange, base) {
  return driver.watch ? driver.watch((event, key) => onChange(event, base + key)) : () => {
  };
}
async function dispose(driver) {
  if (typeof driver.dispose === "function") {
    await asyncCall(driver.dispose);
  }
}

const _assets = {

};

const normalizeKey = function normalizeKey(key) {
  if (!key) {
    return "";
  }
  return key.split("?")[0]?.replace(/[/\\]/g, ":").replace(/:+/g, ":").replace(/^:|:$/g, "") || "";
};

const assets$1 = {
  getKeys() {
    return Promise.resolve(Object.keys(_assets))
  },
  hasItem (id) {
    id = normalizeKey(id);
    return Promise.resolve(id in _assets)
  },
  getItem (id) {
    id = normalizeKey(id);
    return Promise.resolve(_assets[id] ? _assets[id].import() : null)
  },
  getMeta (id) {
    id = normalizeKey(id);
    return Promise.resolve(_assets[id] ? _assets[id].meta : {})
  }
};

function defineDriver(factory) {
  return factory;
}
function createError(driver, message, opts) {
  const err = new Error(`[unstorage] [${driver}] ${message}`, opts);
  if (Error.captureStackTrace) {
    Error.captureStackTrace(err, createError);
  }
  return err;
}
function createRequiredError(driver, name) {
  if (Array.isArray(name)) {
    return createError(
      driver,
      `Missing some of the required options ${name.map((n) => "`" + n + "`").join(", ")}`
    );
  }
  return createError(driver, `Missing required option \`${name}\`.`);
}

const DRIVER_NAME$1 = "memory";
const unstorage_47drivers_47memory = defineDriver(() => {
  const data = /* @__PURE__ */ new Map();
  return {
    name: DRIVER_NAME$1,
    getInstance: () => data,
    hasItem(key) {
      return data.has(key);
    },
    getItem(key) {
      return data.get(key) ?? null;
    },
    getItemRaw(key) {
      return data.get(key) ?? null;
    },
    setItem(key, value) {
      data.set(key, value);
    },
    setItemRaw(key, value) {
      data.set(key, value);
    },
    removeItem(key) {
      data.delete(key);
    },
    getKeys() {
      return [...data.keys()];
    },
    clear() {
      data.clear();
    },
    dispose() {
      data.clear();
    }
  };
});

function ignoreNotfound(err) {
  return err.code === "ENOENT" || err.code === "EISDIR" ? null : err;
}
function ignoreExists(err) {
  return err.code === "EEXIST" ? null : err;
}
async function writeFile(path, data, encoding) {
  await ensuredir(dirname$1(path));
  return promises.writeFile(path, data, encoding);
}
function readFile(path, encoding) {
  return promises.readFile(path, encoding).catch(ignoreNotfound);
}
function unlink(path) {
  return promises.unlink(path).catch(ignoreNotfound);
}
function readdir(dir) {
  return promises.readdir(dir, { withFileTypes: true }).catch(ignoreNotfound).then((r) => r || []);
}
async function ensuredir(dir) {
  if (existsSync(dir)) {
    return;
  }
  await ensuredir(dirname$1(dir)).catch(ignoreExists);
  await promises.mkdir(dir).catch(ignoreExists);
}
async function readdirRecursive(dir, ignore, maxDepth) {
  if (ignore && ignore(dir)) {
    return [];
  }
  const entries = await readdir(dir);
  const files = [];
  await Promise.all(
    entries.map(async (entry) => {
      const entryPath = resolve$1(dir, entry.name);
      if (entry.isDirectory()) {
        if (maxDepth === void 0 || maxDepth > 0) {
          const dirFiles = await readdirRecursive(
            entryPath,
            ignore,
            maxDepth === void 0 ? void 0 : maxDepth - 1
          );
          files.push(...dirFiles.map((f) => entry.name + "/" + f));
        }
      } else {
        if (!(ignore && ignore(entry.name))) {
          files.push(entry.name);
        }
      }
    })
  );
  return files;
}
async function rmRecursive(dir) {
  const entries = await readdir(dir);
  await Promise.all(
    entries.map((entry) => {
      const entryPath = resolve$1(dir, entry.name);
      if (entry.isDirectory()) {
        return rmRecursive(entryPath).then(() => promises.rmdir(entryPath));
      } else {
        return promises.unlink(entryPath);
      }
    })
  );
}

const PATH_TRAVERSE_RE = /\.\.:|\.\.$/;
const DRIVER_NAME = "fs-lite";
const unstorage_47drivers_47fs_45lite = defineDriver((opts = {}) => {
  if (!opts.base) {
    throw createRequiredError(DRIVER_NAME, "base");
  }
  opts.base = resolve$1(opts.base);
  const r = (key) => {
    if (PATH_TRAVERSE_RE.test(key)) {
      throw createError(
        DRIVER_NAME,
        `Invalid key: ${JSON.stringify(key)}. It should not contain .. segments`
      );
    }
    const resolved = join(opts.base, key.replace(/:/g, "/"));
    return resolved;
  };
  return {
    name: DRIVER_NAME,
    options: opts,
    flags: {
      maxDepth: true
    },
    hasItem(key) {
      return existsSync(r(key));
    },
    getItem(key) {
      return readFile(r(key), "utf8");
    },
    getItemRaw(key) {
      return readFile(r(key));
    },
    async getMeta(key) {
      const { atime, mtime, size, birthtime, ctime } = await promises.stat(r(key)).catch(() => ({}));
      return { atime, mtime, size, birthtime, ctime };
    },
    setItem(key, value) {
      if (opts.readOnly) {
        return;
      }
      return writeFile(r(key), value, "utf8");
    },
    setItemRaw(key, value) {
      if (opts.readOnly) {
        return;
      }
      return writeFile(r(key), value);
    },
    removeItem(key) {
      if (opts.readOnly) {
        return;
      }
      return unlink(r(key));
    },
    getKeys(_base, topts) {
      return readdirRecursive(r("."), opts.ignore, topts?.maxDepth);
    },
    async clear() {
      if (opts.readOnly || opts.noClear) {
        return;
      }
      await rmRecursive(r("."));
    }
  };
});

const storage = createStorage({});

storage.mount('/assets', assets$1);

storage.mount('redis', unstorage_47drivers_47memory({"driver":"memory","max":1000}));
storage.mount('data', unstorage_47drivers_47fs_45lite({"driver":"fsLite","base":"./.data/kv"}));

function useStorage(base = "") {
  return base ? prefixStorage(storage, base) : storage;
}

function serialize$1(o){return typeof o=="string"?`'${o}'`:new c().serialize(o)}const c=/*@__PURE__*/function(){class o{#t=new Map;compare(t,r){const e=typeof t,n=typeof r;return e==="string"&&n==="string"?t.localeCompare(r):e==="number"&&n==="number"?t-r:String.prototype.localeCompare.call(this.serialize(t,true),this.serialize(r,true))}serialize(t,r){if(t===null)return "null";switch(typeof t){case "string":return r?t:`'${t}'`;case "bigint":return `${t}n`;case "object":return this.$object(t);case "function":return this.$function(t)}return String(t)}serializeObject(t){const r=Object.prototype.toString.call(t);if(r!=="[object Object]")return this.serializeBuiltInType(r.length<10?`unknown:${r}`:r.slice(8,-1),t);const e=t.constructor,n=e===Object||e===void 0?"":e.name;if(n!==""&&globalThis[n]===e)return this.serializeBuiltInType(n,t);if(typeof t.toJSON=="function"){const i=t.toJSON();return n+(i!==null&&typeof i=="object"?this.$object(i):`(${this.serialize(i)})`)}return this.serializeObjectEntries(n,Object.entries(t))}serializeBuiltInType(t,r){const e=this["$"+t];if(e)return e.call(this,r);if(typeof r?.entries=="function")return this.serializeObjectEntries(t,r.entries());throw new Error(`Cannot serialize ${t}`)}serializeObjectEntries(t,r){const e=Array.from(r).sort((i,a)=>this.compare(i[0],a[0]));let n=`${t}{`;for(let i=0;i<e.length;i++){const[a,l]=e[i];n+=`${this.serialize(a,true)}:${this.serialize(l)}`,i<e.length-1&&(n+=",");}return n+"}"}$object(t){let r=this.#t.get(t);return r===void 0&&(this.#t.set(t,`#${this.#t.size}`),r=this.serializeObject(t),this.#t.set(t,r)),r}$function(t){const r=Function.prototype.toString.call(t);return r.slice(-15)==="[native code] }"?`${t.name||""}()[native]`:`${t.name}(${t.length})${r.replace(/\s*\n\s*/g,"")}`}$Array(t){let r="[";for(let e=0;e<t.length;e++)r+=this.serialize(t[e]),e<t.length-1&&(r+=",");return r+"]"}$Date(t){try{return `Date(${t.toISOString()})`}catch{return "Date(null)"}}$ArrayBuffer(t){return `ArrayBuffer[${new Uint8Array(t).join(",")}]`}$Set(t){return `Set${this.$Array(Array.from(t).sort((r,e)=>this.compare(r,e)))}`}$Map(t){return this.serializeObjectEntries("Map",t.entries())}}for(const s of ["Error","RegExp","URL"])o.prototype["$"+s]=function(t){return `${s}(${t})`};for(const s of ["Int8Array","Uint8Array","Uint8ClampedArray","Int16Array","Uint16Array","Int32Array","Uint32Array","Float32Array","Float64Array"])o.prototype["$"+s]=function(t){return `${s}[${t.join(",")}]`};for(const s of ["BigInt64Array","BigUint64Array"])o.prototype["$"+s]=function(t){return `${s}[${t.join("n,")}${t.length>0?"n":""}]`};return o}();

function isEqual(object1, object2) {
  if (object1 === object2) {
    return true;
  }
  if (serialize$1(object1) === serialize$1(object2)) {
    return true;
  }
  return false;
}

const e=globalThis.process?.getBuiltinModule?.("crypto")?.hash,r="sha256",s="base64url";function digest(t){if(e)return e(r,t,s);const o=createHash(r).update(t);return globalThis.process?.versions?.webcontainer?o.digest().toString(s):o.digest(s)}

const Hasher = /* @__PURE__ */ (() => {
  class Hasher2 {
    buff = "";
    #context = /* @__PURE__ */ new Map();
    write(str) {
      this.buff += str;
    }
    dispatch(value) {
      const type = value === null ? "null" : typeof value;
      return this[type](value);
    }
    object(object) {
      if (object && typeof object.toJSON === "function") {
        return this.object(object.toJSON());
      }
      const objString = Object.prototype.toString.call(object);
      let objType = "";
      const objectLength = objString.length;
      objType = objectLength < 10 ? "unknown:[" + objString + "]" : objString.slice(8, objectLength - 1);
      objType = objType.toLowerCase();
      let objectNumber = null;
      if ((objectNumber = this.#context.get(object)) === void 0) {
        this.#context.set(object, this.#context.size);
      } else {
        return this.dispatch("[CIRCULAR:" + objectNumber + "]");
      }
      if (typeof Buffer !== "undefined" && Buffer.isBuffer && Buffer.isBuffer(object)) {
        this.write("buffer:");
        return this.write(object.toString("utf8"));
      }
      if (objType !== "object" && objType !== "function" && objType !== "asyncfunction") {
        if (this[objType]) {
          this[objType](object);
        } else {
          this.unknown(object, objType);
        }
      } else {
        const keys = Object.keys(object).sort();
        const extraKeys = [];
        this.write("object:" + (keys.length + extraKeys.length) + ":");
        const dispatchForKey = (key) => {
          this.dispatch(key);
          this.write(":");
          this.dispatch(object[key]);
          this.write(",");
        };
        for (const key of keys) {
          dispatchForKey(key);
        }
        for (const key of extraKeys) {
          dispatchForKey(key);
        }
      }
    }
    array(arr, unordered) {
      unordered = unordered === void 0 ? false : unordered;
      this.write("array:" + arr.length + ":");
      if (!unordered || arr.length <= 1) {
        for (const entry of arr) {
          this.dispatch(entry);
        }
        return;
      }
      const contextAdditions = /* @__PURE__ */ new Map();
      const entries = arr.map((entry) => {
        const hasher = new Hasher2();
        hasher.dispatch(entry);
        for (const [key, value] of hasher.#context) {
          contextAdditions.set(key, value);
        }
        return hasher.toString();
      });
      this.#context = contextAdditions;
      entries.sort();
      return this.array(entries, false);
    }
    date(date) {
      return this.write("date:" + date.toJSON());
    }
    symbol(sym) {
      return this.write("symbol:" + sym.toString());
    }
    unknown(value, type) {
      this.write(type);
      if (!value) {
        return;
      }
      this.write(":");
      if (value && typeof value.entries === "function") {
        return this.array(
          [...value.entries()],
          true
          /* ordered */
        );
      }
    }
    error(err) {
      return this.write("error:" + err.toString());
    }
    boolean(bool) {
      return this.write("bool:" + bool);
    }
    string(string) {
      this.write("string:" + string.length + ":");
      this.write(string);
    }
    function(fn) {
      this.write("fn:");
      if (isNativeFunction(fn)) {
        this.dispatch("[native]");
      } else {
        this.dispatch(fn.toString());
      }
    }
    number(number) {
      return this.write("number:" + number);
    }
    null() {
      return this.write("Null");
    }
    undefined() {
      return this.write("Undefined");
    }
    regexp(regex) {
      return this.write("regex:" + regex.toString());
    }
    arraybuffer(arr) {
      this.write("arraybuffer:");
      return this.dispatch(new Uint8Array(arr));
    }
    url(url) {
      return this.write("url:" + url.toString());
    }
    map(map) {
      this.write("map:");
      const arr = [...map];
      return this.array(arr, false);
    }
    set(set) {
      this.write("set:");
      const arr = [...set];
      return this.array(arr, false);
    }
    bigint(number) {
      return this.write("bigint:" + number.toString());
    }
  }
  for (const type of [
    "uint8array",
    "uint8clampedarray",
    "unt8array",
    "uint16array",
    "unt16array",
    "uint32array",
    "unt32array",
    "float32array",
    "float64array"
  ]) {
    Hasher2.prototype[type] = function(arr) {
      this.write(type + ":");
      return this.array([...arr], false);
    };
  }
  function isNativeFunction(f) {
    if (typeof f !== "function") {
      return false;
    }
    return Function.prototype.toString.call(f).slice(
      -15
      /* "[native code] }".length */
    ) === "[native code] }";
  }
  return Hasher2;
})();
function serialize(object) {
  const hasher = new Hasher();
  hasher.dispatch(object);
  return hasher.buff;
}
function hash(value) {
  return digest(typeof value === "string" ? value : serialize(value)).replace(/[-_]/g, "").slice(0, 10);
}

function defaultCacheOptions() {
  return {
    name: "_",
    base: "/cache",
    swr: true,
    maxAge: 1
  };
}
function defineCachedFunction(fn, opts = {}) {
  opts = { ...defaultCacheOptions(), ...opts };
  const pending = {};
  const group = opts.group || "nitro/functions";
  const name = opts.name || fn.name || "_";
  const integrity = opts.integrity || hash([fn, opts]);
  const validate = opts.validate || ((entry) => entry.value !== void 0);
  async function get(key, resolver, shouldInvalidateCache, event) {
    const cacheKey = [opts.base, group, name, key + ".json"].filter(Boolean).join(":").replace(/:\/$/, ":index");
    let entry = await useStorage().getItem(cacheKey).catch((error) => {
      console.error(`[cache] Cache read error.`, error);
      useNitroApp().captureError(error, { event, tags: ["cache"] });
    }) || {};
    if (typeof entry !== "object") {
      entry = {};
      const error = new Error("Malformed data read from cache.");
      console.error("[cache]", error);
      useNitroApp().captureError(error, { event, tags: ["cache"] });
    }
    const ttl = (opts.maxAge ?? 0) * 1e3;
    if (ttl) {
      entry.expires = Date.now() + ttl;
    }
    const expired = shouldInvalidateCache || entry.integrity !== integrity || ttl && Date.now() - (entry.mtime || 0) > ttl || validate(entry) === false;
    const _resolve = async () => {
      const isPending = pending[key];
      if (!isPending) {
        if (entry.value !== void 0 && (opts.staleMaxAge || 0) >= 0 && opts.swr === false) {
          entry.value = void 0;
          entry.integrity = void 0;
          entry.mtime = void 0;
          entry.expires = void 0;
        }
        pending[key] = Promise.resolve(resolver());
      }
      try {
        entry.value = await pending[key];
      } catch (error) {
        if (!isPending) {
          delete pending[key];
        }
        throw error;
      }
      if (!isPending) {
        entry.mtime = Date.now();
        entry.integrity = integrity;
        delete pending[key];
        if (validate(entry) !== false) {
          let setOpts;
          if (opts.maxAge && !opts.swr) {
            setOpts = { ttl: opts.maxAge };
          }
          const promise = useStorage().setItem(cacheKey, entry, setOpts).catch((error) => {
            console.error(`[cache] Cache write error.`, error);
            useNitroApp().captureError(error, { event, tags: ["cache"] });
          });
          if (event?.waitUntil) {
            event.waitUntil(promise);
          }
        }
      }
    };
    const _resolvePromise = expired ? _resolve() : Promise.resolve();
    if (entry.value === void 0) {
      await _resolvePromise;
    } else if (expired && event && event.waitUntil) {
      event.waitUntil(_resolvePromise);
    }
    if (opts.swr && validate(entry) !== false) {
      _resolvePromise.catch((error) => {
        console.error(`[cache] SWR handler error.`, error);
        useNitroApp().captureError(error, { event, tags: ["cache"] });
      });
      return entry;
    }
    return _resolvePromise.then(() => entry);
  }
  return async (...args) => {
    const shouldBypassCache = await opts.shouldBypassCache?.(...args);
    if (shouldBypassCache) {
      return fn(...args);
    }
    const key = await (opts.getKey || getKey)(...args);
    const shouldInvalidateCache = await opts.shouldInvalidateCache?.(...args);
    const entry = await get(
      key,
      () => fn(...args),
      shouldInvalidateCache,
      args[0] && isEvent(args[0]) ? args[0] : void 0
    );
    let value = entry.value;
    if (opts.transform) {
      value = await opts.transform(entry, ...args) || value;
    }
    return value;
  };
}
function cachedFunction(fn, opts = {}) {
  return defineCachedFunction(fn, opts);
}
function getKey(...args) {
  return args.length > 0 ? hash(args) : "";
}
function escapeKey(key) {
  return String(key).replace(/\W/g, "");
}
function defineCachedEventHandler(handler, opts = defaultCacheOptions()) {
  const variableHeaderNames = (opts.varies || []).filter(Boolean).map((h) => h.toLowerCase()).sort();
  const _opts = {
    ...opts,
    getKey: async (event) => {
      const customKey = await opts.getKey?.(event);
      if (customKey) {
        return escapeKey(customKey);
      }
      const _path = event.node.req.originalUrl || event.node.req.url || event.path;
      let _pathname;
      try {
        _pathname = escapeKey(decodeURI(parseURL(_path).pathname)).slice(0, 16) || "index";
      } catch {
        _pathname = "-";
      }
      const _hashedPath = `${_pathname}.${hash(_path)}`;
      const _headers = variableHeaderNames.map((header) => [header, event.node.req.headers[header]]).map(([name, value]) => `${escapeKey(name)}.${hash(value)}`);
      return [_hashedPath, ..._headers].join(":");
    },
    validate: (entry) => {
      if (!entry.value) {
        return false;
      }
      if (entry.value.code >= 400) {
        return false;
      }
      if (entry.value.body === void 0) {
        return false;
      }
      if (entry.value.headers.etag === "undefined" || entry.value.headers["last-modified"] === "undefined") {
        return false;
      }
      return true;
    },
    group: opts.group || "nitro/handlers",
    integrity: opts.integrity || hash([handler, opts])
  };
  const _cachedHandler = cachedFunction(
    async (incomingEvent) => {
      const variableHeaders = {};
      for (const header of variableHeaderNames) {
        const value = incomingEvent.node.req.headers[header];
        if (value !== void 0) {
          variableHeaders[header] = value;
        }
      }
      const reqProxy = cloneWithProxy(incomingEvent.node.req, {
        headers: variableHeaders
      });
      const resHeaders = {};
      let _resSendBody;
      const resProxy = cloneWithProxy(incomingEvent.node.res, {
        statusCode: 200,
        writableEnded: false,
        writableFinished: false,
        headersSent: false,
        closed: false,
        getHeader(name) {
          return resHeaders[name];
        },
        setHeader(name, value) {
          resHeaders[name] = value;
          return this;
        },
        getHeaderNames() {
          return Object.keys(resHeaders);
        },
        hasHeader(name) {
          return name in resHeaders;
        },
        removeHeader(name) {
          delete resHeaders[name];
        },
        getHeaders() {
          return resHeaders;
        },
        end(chunk, arg2, arg3) {
          if (typeof chunk === "string") {
            _resSendBody = chunk;
          }
          if (typeof arg2 === "function") {
            arg2();
          }
          if (typeof arg3 === "function") {
            arg3();
          }
          return this;
        },
        write(chunk, arg2, arg3) {
          if (typeof chunk === "string") {
            _resSendBody = chunk;
          }
          if (typeof arg2 === "function") {
            arg2(void 0);
          }
          if (typeof arg3 === "function") {
            arg3();
          }
          return true;
        },
        writeHead(statusCode, headers2) {
          this.statusCode = statusCode;
          if (headers2) {
            if (Array.isArray(headers2) || typeof headers2 === "string") {
              throw new TypeError("Raw headers  is not supported.");
            }
            for (const header in headers2) {
              const value = headers2[header];
              if (value !== void 0) {
                this.setHeader(
                  header,
                  value
                );
              }
            }
          }
          return this;
        }
      });
      const event = createEvent(reqProxy, resProxy);
      event.fetch = (url, fetchOptions) => fetchWithEvent(event, url, fetchOptions, {
        fetch: useNitroApp().localFetch
      });
      event.$fetch = (url, fetchOptions) => fetchWithEvent(event, url, fetchOptions, {
        fetch: globalThis.$fetch
      });
      event.waitUntil = incomingEvent.waitUntil;
      event.context = incomingEvent.context;
      event.context.cache = {
        options: _opts
      };
      const body = await handler(event) || _resSendBody;
      const headers = event.node.res.getHeaders();
      headers.etag = String(
        headers.Etag || headers.etag || `W/"${hash(body)}"`
      );
      headers["last-modified"] = String(
        headers["Last-Modified"] || headers["last-modified"] || (/* @__PURE__ */ new Date()).toUTCString()
      );
      const cacheControl = [];
      if (opts.swr) {
        if (opts.maxAge) {
          cacheControl.push(`s-maxage=${opts.maxAge}`);
        }
        if (opts.staleMaxAge) {
          cacheControl.push(`stale-while-revalidate=${opts.staleMaxAge}`);
        } else {
          cacheControl.push("stale-while-revalidate");
        }
      } else if (opts.maxAge) {
        cacheControl.push(`max-age=${opts.maxAge}`);
      }
      if (cacheControl.length > 0) {
        headers["cache-control"] = cacheControl.join(", ");
      }
      const cacheEntry = {
        code: event.node.res.statusCode,
        headers,
        body
      };
      return cacheEntry;
    },
    _opts
  );
  return defineEventHandler(async (event) => {
    if (opts.headersOnly) {
      if (handleCacheHeaders(event, { maxAge: opts.maxAge })) {
        return;
      }
      return handler(event);
    }
    const response = await _cachedHandler(
      event
    );
    if (event.node.res.headersSent || event.node.res.writableEnded) {
      return response.body;
    }
    if (handleCacheHeaders(event, {
      modifiedTime: new Date(response.headers["last-modified"]),
      etag: response.headers.etag,
      maxAge: opts.maxAge
    })) {
      return;
    }
    event.node.res.statusCode = response.code;
    for (const name in response.headers) {
      const value = response.headers[name];
      if (name === "set-cookie") {
        event.node.res.appendHeader(
          name,
          splitCookiesString(value)
        );
      } else {
        if (value !== void 0) {
          event.node.res.setHeader(name, value);
        }
      }
    }
    return response.body;
  });
}
function cloneWithProxy(obj, overrides) {
  return new Proxy(obj, {
    get(target, property, receiver) {
      if (property in overrides) {
        return overrides[property];
      }
      return Reflect.get(target, property, receiver);
    },
    set(target, property, value, receiver) {
      if (property in overrides) {
        overrides[property] = value;
        return true;
      }
      return Reflect.set(target, property, value, receiver);
    }
  });
}
const cachedEventHandler = defineCachedEventHandler;

function klona(x) {
	if (typeof x !== 'object') return x;

	var k, tmp, str=Object.prototype.toString.call(x);

	if (str === '[object Object]') {
		if (x.constructor !== Object && typeof x.constructor === 'function') {
			tmp = new x.constructor();
			for (k in x) {
				if (x.hasOwnProperty(k) && tmp[k] !== x[k]) {
					tmp[k] = klona(x[k]);
				}
			}
		} else {
			tmp = {}; // null
			for (k in x) {
				if (k === '__proto__') {
					Object.defineProperty(tmp, k, {
						value: klona(x[k]),
						configurable: true,
						enumerable: true,
						writable: true,
					});
				} else {
					tmp[k] = klona(x[k]);
				}
			}
		}
		return tmp;
	}

	if (str === '[object Array]') {
		k = x.length;
		for (tmp=Array(k); k--;) {
			tmp[k] = klona(x[k]);
		}
		return tmp;
	}

	if (str === '[object Set]') {
		tmp = new Set;
		x.forEach(function (val) {
			tmp.add(klona(val));
		});
		return tmp;
	}

	if (str === '[object Map]') {
		tmp = new Map;
		x.forEach(function (val, key) {
			tmp.set(klona(key), klona(val));
		});
		return tmp;
	}

	if (str === '[object Date]') {
		return new Date(+x);
	}

	if (str === '[object RegExp]') {
		tmp = new RegExp(x.source, x.flags);
		tmp.lastIndex = x.lastIndex;
		return tmp;
	}

	if (str === '[object DataView]') {
		return new x.constructor( klona(x.buffer) );
	}

	if (str === '[object ArrayBuffer]') {
		return x.slice(0);
	}

	// ArrayBuffer.isView(x)
	// ~> `new` bcuz `Buffer.slice` => ref
	if (str.slice(-6) === 'Array]') {
		return new x.constructor(x);
	}

	return x;
}

const inlineAppConfig = {
  "nuxt": {}
};



const appConfig = defuFn(inlineAppConfig);

const NUMBER_CHAR_RE = /\d/;
const STR_SPLITTERS = ["-", "_", "/", "."];
function isUppercase(char = "") {
  if (NUMBER_CHAR_RE.test(char)) {
    return void 0;
  }
  return char !== char.toLowerCase();
}
function splitByCase(str, separators) {
  const splitters = STR_SPLITTERS;
  const parts = [];
  if (!str || typeof str !== "string") {
    return parts;
  }
  let buff = "";
  let previousUpper;
  let previousSplitter;
  for (const char of str) {
    const isSplitter = splitters.includes(char);
    if (isSplitter === true) {
      parts.push(buff);
      buff = "";
      previousUpper = void 0;
      continue;
    }
    const isUpper = isUppercase(char);
    if (previousSplitter === false) {
      if (previousUpper === false && isUpper === true) {
        parts.push(buff);
        buff = char;
        previousUpper = isUpper;
        continue;
      }
      if (previousUpper === true && isUpper === false && buff.length > 1) {
        const lastChar = buff.at(-1);
        parts.push(buff.slice(0, Math.max(0, buff.length - 1)));
        buff = lastChar + char;
        previousUpper = isUpper;
        continue;
      }
    }
    buff += char;
    previousUpper = isUpper;
    previousSplitter = isSplitter;
  }
  parts.push(buff);
  return parts;
}
function kebabCase(str, joiner) {
  return str ? (Array.isArray(str) ? str : splitByCase(str)).map((p) => p.toLowerCase()).join(joiner) : "";
}
function snakeCase(str) {
  return kebabCase(str || "", "_");
}

function getEnv(key, opts) {
  const envKey = snakeCase(key).toUpperCase();
  return destr(
    process.env[opts.prefix + envKey] ?? process.env[opts.altPrefix + envKey]
  );
}
function _isObject(input) {
  return typeof input === "object" && !Array.isArray(input);
}
function applyEnv(obj, opts, parentKey = "") {
  for (const key in obj) {
    const subKey = parentKey ? `${parentKey}_${key}` : key;
    const envValue = getEnv(subKey, opts);
    if (_isObject(obj[key])) {
      if (_isObject(envValue)) {
        obj[key] = { ...obj[key], ...envValue };
        applyEnv(obj[key], opts, subKey);
      } else if (envValue === void 0) {
        applyEnv(obj[key], opts, subKey);
      } else {
        obj[key] = envValue ?? obj[key];
      }
    } else {
      obj[key] = envValue ?? obj[key];
    }
    if (opts.envExpansion && typeof obj[key] === "string") {
      obj[key] = _expandFromEnv(obj[key]);
    }
  }
  return obj;
}
const envExpandRx = /\{\{([^{}]*)\}\}/g;
function _expandFromEnv(value) {
  return value.replace(envExpandRx, (match, key) => {
    return process.env[key] || match;
  });
}

const _inlineRuntimeConfig = {
  "app": {
    "baseURL": "/",
    "buildId": "72f4a96d-66c3-4ff3-8e35-cbf0a02671bc",
    "buildAssetsDir": "/_nuxt/",
    "cdnURL": ""
  },
  "nitro": {
    "envPrefix": "NUXT_",
    "routeRules": {
      "/__nuxt_error": {
        "cache": false
      },
      "/_nuxt/builds/meta/**": {
        "headers": {
          "cache-control": "public, max-age=31536000, immutable"
        }
      },
      "/_nuxt/builds/**": {
        "headers": {
          "cache-control": "public, max-age=1, immutable"
        }
      },
      "/_nuxt/**": {
        "headers": {
          "cache-control": "public, max-age=31536000, immutable"
        }
      }
    }
  },
  "public": {
    "apiBase": "/api",
    "apiServer": "http://localhost:8080/api",
    "ogApiUrl": "/api/og-image",
    "siteUrl": "https://yourdomain.com"
  }
};
const envOptions = {
  prefix: "NITRO_",
  altPrefix: _inlineRuntimeConfig.nitro.envPrefix ?? process.env.NITRO_ENV_PREFIX ?? "_",
  envExpansion: _inlineRuntimeConfig.nitro.envExpansion ?? process.env.NITRO_ENV_EXPANSION ?? false
};
const _sharedRuntimeConfig = _deepFreeze(
  applyEnv(klona(_inlineRuntimeConfig), envOptions)
);
function useRuntimeConfig(event) {
  if (!event) {
    return _sharedRuntimeConfig;
  }
  if (event.context.nitro.runtimeConfig) {
    return event.context.nitro.runtimeConfig;
  }
  const runtimeConfig = klona(_inlineRuntimeConfig);
  applyEnv(runtimeConfig, envOptions);
  event.context.nitro.runtimeConfig = runtimeConfig;
  return runtimeConfig;
}
_deepFreeze(klona(appConfig));
function _deepFreeze(object) {
  const propNames = Object.getOwnPropertyNames(object);
  for (const name of propNames) {
    const value = object[name];
    if (value && typeof value === "object") {
      _deepFreeze(value);
    }
  }
  return Object.freeze(object);
}
new Proxy(/* @__PURE__ */ Object.create(null), {
  get: (_, prop) => {
    console.warn(
      "Please use `useRuntimeConfig()` instead of accessing config directly."
    );
    const runtimeConfig = useRuntimeConfig();
    if (prop in runtimeConfig) {
      return runtimeConfig[prop];
    }
    return void 0;
  }
});

function createContext(opts = {}) {
  let currentInstance;
  let isSingleton = false;
  const checkConflict = (instance) => {
    if (currentInstance && currentInstance !== instance) {
      throw new Error("Context conflict");
    }
  };
  let als;
  if (opts.asyncContext) {
    const _AsyncLocalStorage = opts.AsyncLocalStorage || globalThis.AsyncLocalStorage;
    if (_AsyncLocalStorage) {
      als = new _AsyncLocalStorage();
    } else {
      console.warn("[unctx] `AsyncLocalStorage` is not provided.");
    }
  }
  const _getCurrentInstance = () => {
    if (als) {
      const instance = als.getStore();
      if (instance !== void 0) {
        return instance;
      }
    }
    return currentInstance;
  };
  return {
    use: () => {
      const _instance = _getCurrentInstance();
      if (_instance === void 0) {
        throw new Error("Context is not available");
      }
      return _instance;
    },
    tryUse: () => {
      return _getCurrentInstance();
    },
    set: (instance, replace) => {
      if (!replace) {
        checkConflict(instance);
      }
      currentInstance = instance;
      isSingleton = true;
    },
    unset: () => {
      currentInstance = void 0;
      isSingleton = false;
    },
    call: (instance, callback) => {
      checkConflict(instance);
      currentInstance = instance;
      try {
        return als ? als.run(instance, callback) : callback();
      } finally {
        if (!isSingleton) {
          currentInstance = void 0;
        }
      }
    },
    async callAsync(instance, callback) {
      currentInstance = instance;
      const onRestore = () => {
        currentInstance = instance;
      };
      const onLeave = () => currentInstance === instance ? onRestore : void 0;
      asyncHandlers.add(onLeave);
      try {
        const r = als ? als.run(instance, callback) : callback();
        if (!isSingleton) {
          currentInstance = void 0;
        }
        return await r;
      } finally {
        asyncHandlers.delete(onLeave);
      }
    }
  };
}
function createNamespace(defaultOpts = {}) {
  const contexts = {};
  return {
    get(key, opts = {}) {
      if (!contexts[key]) {
        contexts[key] = createContext({ ...defaultOpts, ...opts });
      }
      return contexts[key];
    }
  };
}
const _globalThis = typeof globalThis !== "undefined" ? globalThis : typeof self !== "undefined" ? self : typeof global !== "undefined" ? global : {};
const globalKey = "__unctx__";
const defaultNamespace = _globalThis[globalKey] || (_globalThis[globalKey] = createNamespace());
const getContext = (key, opts = {}) => defaultNamespace.get(key, opts);
const asyncHandlersKey = "__unctx_async_handlers__";
const asyncHandlers = _globalThis[asyncHandlersKey] || (_globalThis[asyncHandlersKey] = /* @__PURE__ */ new Set());
function executeAsync(function_) {
  const restores = [];
  for (const leaveHandler of asyncHandlers) {
    const restore2 = leaveHandler();
    if (restore2) {
      restores.push(restore2);
    }
  }
  const restore = () => {
    for (const restore2 of restores) {
      restore2();
    }
  };
  let awaitable = function_();
  if (awaitable && typeof awaitable === "object" && "catch" in awaitable) {
    awaitable = awaitable.catch((error) => {
      restore();
      throw error;
    });
  }
  return [awaitable, restore];
}

getContext("nitro-app", {
  asyncContext: false,
  AsyncLocalStorage: void 0
});

const config = useRuntimeConfig();
const _routeRulesMatcher = toRouteMatcher(
  createRouter$1({ routes: config.nitro.routeRules })
);
function createRouteRulesHandler(ctx) {
  return eventHandler((event) => {
    const routeRules = getRouteRules(event);
    if (routeRules.headers) {
      setHeaders(event, routeRules.headers);
    }
    if (routeRules.redirect) {
      let target = routeRules.redirect.to;
      if (target.endsWith("/**")) {
        let targetPath = event.path;
        const strpBase = routeRules.redirect._redirectStripBase;
        if (strpBase) {
          targetPath = withoutBase(targetPath, strpBase);
        }
        target = joinURL(target.slice(0, -3), targetPath);
      } else if (event.path.includes("?")) {
        const query = getQuery$1(event.path);
        target = withQuery(target, query);
      }
      return sendRedirect(event, target, routeRules.redirect.statusCode);
    }
    if (routeRules.proxy) {
      let target = routeRules.proxy.to;
      if (target.endsWith("/**")) {
        let targetPath = event.path;
        const strpBase = routeRules.proxy._proxyStripBase;
        if (strpBase) {
          targetPath = withoutBase(targetPath, strpBase);
        }
        target = joinURL(target.slice(0, -3), targetPath);
      } else if (event.path.includes("?")) {
        const query = getQuery$1(event.path);
        target = withQuery(target, query);
      }
      return proxyRequest(event, target, {
        fetch: ctx.localFetch,
        ...routeRules.proxy
      });
    }
  });
}
function getRouteRules(event) {
  event.context._nitro = event.context._nitro || {};
  if (!event.context._nitro.routeRules) {
    event.context._nitro.routeRules = getRouteRulesForPath(
      withoutBase(event.path.split("?")[0], useRuntimeConfig().app.baseURL)
    );
  }
  return event.context._nitro.routeRules;
}
function getRouteRulesForPath(path) {
  return defu({}, ..._routeRulesMatcher.matchAll(path).reverse());
}

function _captureError(error, type) {
  console.error(`[${type}]`, error);
  useNitroApp().captureError(error, { tags: [type] });
}
function trapUnhandledNodeErrors() {
  process.on(
    "unhandledRejection",
    (error) => _captureError(error, "unhandledRejection")
  );
  process.on(
    "uncaughtException",
    (error) => _captureError(error, "uncaughtException")
  );
}
function joinHeaders(value) {
  return Array.isArray(value) ? value.join(", ") : String(value);
}
function normalizeFetchResponse(response) {
  if (!response.headers.has("set-cookie")) {
    return response;
  }
  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: normalizeCookieHeaders(response.headers)
  });
}
function normalizeCookieHeader(header = "") {
  return splitCookiesString(joinHeaders(header));
}
function normalizeCookieHeaders(headers) {
  const outgoingHeaders = new Headers();
  for (const [name, header] of headers) {
    if (name === "set-cookie") {
      for (const cookie of normalizeCookieHeader(header)) {
        outgoingHeaders.append("set-cookie", cookie);
      }
    } else {
      outgoingHeaders.set(name, joinHeaders(header));
    }
  }
  return outgoingHeaders;
}

function isJsonRequest(event) {
  if (hasReqHeader(event, "accept", "text/html")) {
    return false;
  }
  return hasReqHeader(event, "accept", "application/json") || hasReqHeader(event, "user-agent", "curl/") || hasReqHeader(event, "user-agent", "httpie/") || hasReqHeader(event, "sec-fetch-mode", "cors") || event.path.startsWith("/api/") || event.path.endsWith(".json");
}
function hasReqHeader(event, name, includes) {
  const value = getRequestHeader(event, name);
  return value && typeof value === "string" && value.toLowerCase().includes(includes);
}

const errorHandler$0 = (async function errorhandler(error, event, { defaultHandler }) {
  if (event.handled || isJsonRequest(event)) {
    return;
  }
  const defaultRes = await defaultHandler(error, event, { json: true });
  const statusCode = error.statusCode || 500;
  if (statusCode === 404 && defaultRes.status === 302) {
    setResponseHeaders(event, defaultRes.headers);
    setResponseStatus(event, defaultRes.status, defaultRes.statusText);
    return send(event, JSON.stringify(defaultRes.body, null, 2));
  }
  const errorObject = defaultRes.body;
  const url = new URL(errorObject.url);
  errorObject.url = withoutBase(url.pathname, useRuntimeConfig(event).app.baseURL) + url.search + url.hash;
  errorObject.message ||= "Server Error";
  errorObject.data ||= error.data;
  errorObject.statusMessage ||= error.statusMessage;
  delete defaultRes.headers["content-type"];
  delete defaultRes.headers["content-security-policy"];
  setResponseHeaders(event, defaultRes.headers);
  const reqHeaders = getRequestHeaders(event);
  const isRenderingError = event.path.startsWith("/__nuxt_error") || !!reqHeaders["x-nuxt-error"];
  const res = isRenderingError ? null : await useNitroApp().localFetch(
    withQuery(joinURL(useRuntimeConfig(event).app.baseURL, "/__nuxt_error"), errorObject),
    {
      headers: { ...reqHeaders, "x-nuxt-error": "true" },
      redirect: "manual"
    }
  ).catch(() => null);
  if (event.handled) {
    return;
  }
  if (!res) {
    const { template } = await import('../_/error-500.mjs');
    setResponseHeader(event, "Content-Type", "text/html;charset=UTF-8");
    return send(event, template(errorObject));
  }
  const html = await res.text();
  for (const [header, value] of res.headers.entries()) {
    if (header === "set-cookie") {
      appendResponseHeader(event, header, value);
      continue;
    }
    setResponseHeader(event, header, value);
  }
  setResponseStatus(event, res.status && res.status !== 200 ? res.status : defaultRes.status, res.statusText || defaultRes.statusText);
  return send(event, html);
});

function defineNitroErrorHandler(handler) {
  return handler;
}

const errorHandler$1 = defineNitroErrorHandler(
  function defaultNitroErrorHandler(error, event) {
    const res = defaultHandler(error, event);
    setResponseHeaders(event, res.headers);
    setResponseStatus(event, res.status, res.statusText);
    return send(event, JSON.stringify(res.body, null, 2));
  }
);
function defaultHandler(error, event, opts) {
  const isSensitive = error.unhandled || error.fatal;
  const statusCode = error.statusCode || 500;
  const statusMessage = error.statusMessage || "Server Error";
  const url = getRequestURL(event, { xForwardedHost: true, xForwardedProto: true });
  if (statusCode === 404) {
    const baseURL = "/";
    if (/^\/[^/]/.test(baseURL) && !url.pathname.startsWith(baseURL)) {
      const redirectTo = `${baseURL}${url.pathname.slice(1)}${url.search}`;
      return {
        status: 302,
        statusText: "Found",
        headers: { location: redirectTo },
        body: `Redirecting...`
      };
    }
  }
  if (isSensitive && !opts?.silent) {
    const tags = [error.unhandled && "[unhandled]", error.fatal && "[fatal]"].filter(Boolean).join(" ");
    console.error(`[request error] ${tags} [${event.method}] ${url}
`, error);
  }
  const headers = {
    "content-type": "application/json",
    // Prevent browser from guessing the MIME types of resources.
    "x-content-type-options": "nosniff",
    // Prevent error page from being embedded in an iframe
    "x-frame-options": "DENY",
    // Prevent browsers from sending the Referer header
    "referrer-policy": "no-referrer",
    // Disable the execution of any js
    "content-security-policy": "script-src 'none'; frame-ancestors 'none';"
  };
  setResponseStatus(event, statusCode, statusMessage);
  if (statusCode === 404 || !getResponseHeader(event, "cache-control")) {
    headers["cache-control"] = "no-cache";
  }
  const body = {
    error: true,
    url: url.href,
    statusCode,
    statusMessage,
    message: isSensitive ? "Server Error" : error.message,
    data: isSensitive ? void 0 : error.data
  };
  return {
    status: statusCode,
    statusText: statusMessage,
    headers,
    body
  };
}

const errorHandlers = [errorHandler$0, errorHandler$1];

async function errorHandler(error, event) {
  for (const handler of errorHandlers) {
    try {
      await handler(error, event, { defaultHandler });
      if (event.handled) {
        return; // Response handled
      }
    } catch(error) {
      // Handler itself thrown, log and continue
      console.error(error);
    }
  }
  // H3 will handle fallback
}

const plugins = [
  
];

const assets = {
  "/favicon.ico": {
    "type": "image/vnd.microsoft.icon",
    "etag": "\"10be-rpqBH6lujZRpPwA1hwtSg+BD434\"",
    "mtime": "2026-06-14T09:29:09.652Z",
    "size": 4286,
    "path": "../public/favicon.ico"
  },
  "/_nuxt/0vdbVihP.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"a8f-9EAKKVlcAGw+gWE7AHEvb+Co8L0\"",
    "mtime": "2026-06-14T09:29:09.562Z",
    "size": 2703,
    "path": "../public/_nuxt/0vdbVihP.js"
  },
  "/_nuxt/0yS_RTVI.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"27b8-UFtGlVZNvT/ZpGPu7MvfRar3fFY\"",
    "mtime": "2026-06-14T09:29:09.562Z",
    "size": 10168,
    "path": "../public/_nuxt/0yS_RTVI.js"
  },
  "/_nuxt/15mzXzNh.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"f1ff-Il0hl2OXC/sZoYEp0VgEoa4cP0o\"",
    "mtime": "2026-06-14T09:29:09.563Z",
    "size": 61951,
    "path": "../public/_nuxt/15mzXzNh.js"
  },
  "/_nuxt/40DFha3I.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"defc-gbxJLf3ZdNSB+CmWuqt1R9HWdyI\"",
    "mtime": "2026-06-14T09:29:09.563Z",
    "size": 57084,
    "path": "../public/_nuxt/40DFha3I.js"
  },
  "/_nuxt/6B1W4Q65.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"6f-D0eqBmKp1pyyIlqhD1F7IXUdvMg\"",
    "mtime": "2026-06-14T09:29:09.563Z",
    "size": 111,
    "path": "../public/_nuxt/6B1W4Q65.js"
  },
  "/_nuxt/AdminPageLayout.Br4FUZUG.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"2e-VBMH1XdfWNQZwivH23fr9cF6hl8\"",
    "mtime": "2026-06-14T09:29:09.563Z",
    "size": 46,
    "path": "../public/_nuxt/AdminPageLayout.Br4FUZUG.css"
  },
  "/_nuxt/AppFooter._781YVOA.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"d76-LRKWLbegun3voC49+7TOA828ws8\"",
    "mtime": "2026-06-14T09:29:09.564Z",
    "size": 3446,
    "path": "../public/_nuxt/AppFooter._781YVOA.css"
  },
  "/_nuxt/B-p6aW7q.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"157-5ALo6txPChAD2/oP2oVSEyLbfhw\"",
    "mtime": "2026-06-14T09:29:09.564Z",
    "size": 343,
    "path": "../public/_nuxt/B-p6aW7q.js"
  },
  "/_nuxt/B20Wrg-n.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"256d-H1bLcqyj5JGLZyUXGuBhemhuf+E\"",
    "mtime": "2026-06-14T09:29:09.564Z",
    "size": 9581,
    "path": "../public/_nuxt/B20Wrg-n.js"
  },
  "/_nuxt/B55S-1qy.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"7b4-p0t9cWlPFcA+ORj1VI/gm8ovXgQ\"",
    "mtime": "2026-06-14T09:29:09.565Z",
    "size": 1972,
    "path": "../public/_nuxt/B55S-1qy.js"
  },
  "/_nuxt/B6osxiHy.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"824-121Q1GOWSgK+c2FoRsWHt0LUdc4\"",
    "mtime": "2026-06-14T09:29:09.564Z",
    "size": 2084,
    "path": "../public/_nuxt/B6osxiHy.js"
  },
  "/_nuxt/B7ks24Uf.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1cb0-fbvCJQeqnLqLEo9fZTkEpR9t6aY\"",
    "mtime": "2026-06-14T09:29:09.564Z",
    "size": 7344,
    "path": "../public/_nuxt/B7ks24Uf.js"
  },
  "/_nuxt/B8UDiWez.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"3222-ahCurXrBA2UprgaEsdRuCrGGvLc\"",
    "mtime": "2026-06-14T09:29:09.578Z",
    "size": 12834,
    "path": "../public/_nuxt/B8UDiWez.js"
  },
  "/_nuxt/BADnFztR.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"22fe-LgSLBpmUDFIqM/7lM1Q4TQu8L90\"",
    "mtime": "2026-06-14T09:29:09.566Z",
    "size": 8958,
    "path": "../public/_nuxt/BADnFztR.js"
  },
  "/_nuxt/BD_kU2HR.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"8a1a-QhL8bHqvPH7UOjd+uHCtttkZPUg\"",
    "mtime": "2026-06-14T09:29:09.566Z",
    "size": 35354,
    "path": "../public/_nuxt/BD_kU2HR.js"
  },
  "/_nuxt/BF7UzXbo.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"14463-lj0Dqxbf6sdXM4oIH5TWDTODpHU\"",
    "mtime": "2026-06-14T09:29:09.567Z",
    "size": 83043,
    "path": "../public/_nuxt/BF7UzXbo.js"
  },
  "/_nuxt/BK0zW0Pa.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"8802-sxp50DiYws6lUU8sw/2axHwHnbo\"",
    "mtime": "2026-06-14T09:29:09.567Z",
    "size": 34818,
    "path": "../public/_nuxt/BK0zW0Pa.js"
  },
  "/_nuxt/BMO4Kg4s.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"3b48-adHuwAt6T+R/Pfk/D2aq2cvhw0g\"",
    "mtime": "2026-06-14T09:29:09.567Z",
    "size": 15176,
    "path": "../public/_nuxt/BMO4Kg4s.js"
  },
  "/_nuxt/BP3HD85z.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"7774-omXWKY/x1oMNbAgy2zpGPxYxnXU\"",
    "mtime": "2026-06-14T09:29:09.568Z",
    "size": 30580,
    "path": "../public/_nuxt/BP3HD85z.js"
  },
  "/_nuxt/BPoJEm7E.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2c42-98GzxXdZu7JpOYt/a/ompAzL3LI\"",
    "mtime": "2026-06-14T09:29:09.569Z",
    "size": 11330,
    "path": "../public/_nuxt/BPoJEm7E.js"
  },
  "/_nuxt/BQQJOfr4.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4c2-GBbFYUr34z66B5lyDf2hP8wKIO4\"",
    "mtime": "2026-06-14T09:29:09.568Z",
    "size": 1218,
    "path": "../public/_nuxt/BQQJOfr4.js"
  },
  "/_nuxt/BRxqanVv.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2e81-PfHQQpMB2yDQuO9cJvsdg5l+t/w\"",
    "mtime": "2026-06-14T09:29:09.569Z",
    "size": 11905,
    "path": "../public/_nuxt/BRxqanVv.js"
  },
  "/_nuxt/BSb3thVN.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"cd4-7MOwohPotRojoZW450+LMgO0KB4\"",
    "mtime": "2026-06-14T09:29:09.569Z",
    "size": 3284,
    "path": "../public/_nuxt/BSb3thVN.js"
  },
  "/_nuxt/BTVMoPte.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"658e-QD/CRmbDibjfZVOdl55FHWe88Xc\"",
    "mtime": "2026-06-14T09:29:09.578Z",
    "size": 25998,
    "path": "../public/_nuxt/BTVMoPte.js"
  },
  "/_nuxt/BXPXir9X.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2b40-KubNFAvNROgkpo8SKKggPHw6QW4\"",
    "mtime": "2026-06-14T09:29:09.578Z",
    "size": 11072,
    "path": "../public/_nuxt/BXPXir9X.js"
  },
  "/_nuxt/B_fFakE4.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"164d-3apezGtYYFSVVezwXsuR2Hj6jaw\"",
    "mtime": "2026-06-14T09:29:09.576Z",
    "size": 5709,
    "path": "../public/_nuxt/B_fFakE4.js"
  },
  "/_nuxt/Bfp-MHs6.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2109-/f4IvksNazujJ7w+zfoJJJEtovY\"",
    "mtime": "2026-06-14T09:29:09.579Z",
    "size": 8457,
    "path": "../public/_nuxt/Bfp-MHs6.js"
  },
  "/_nuxt/Bjh6Msyv.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"95-4aPN9Bcq5kWkvmd2/MTFguw7/a4\"",
    "mtime": "2026-06-14T09:29:09.580Z",
    "size": 149,
    "path": "../public/_nuxt/Bjh6Msyv.js"
  },
  "/_nuxt/Bk_rJcZu.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"5b-nOfU5mgxjjfjfGzbkr9VXZzOrp4\"",
    "mtime": "2026-06-14T09:29:09.584Z",
    "size": 91,
    "path": "../public/_nuxt/Bk_rJcZu.js"
  },
  "/_nuxt/Bm4WDvqO.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"7c1b-2sU3UkZvbbgfnRota9vTr3iErSE\"",
    "mtime": "2026-06-14T09:29:09.585Z",
    "size": 31771,
    "path": "../public/_nuxt/Bm4WDvqO.js"
  },
  "/_nuxt/BnsUzXib.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"5d44-r2rxzW/AsIllVOF7gJ8rPmGf+TE\"",
    "mtime": "2026-06-14T09:29:09.584Z",
    "size": 23876,
    "path": "../public/_nuxt/BnsUzXib.js"
  },
  "/_nuxt/Bqfut2EJ.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"212c-dGrCUm51XCXCgAdHmGjALRMbszM\"",
    "mtime": "2026-06-14T09:29:09.585Z",
    "size": 8492,
    "path": "../public/_nuxt/Bqfut2EJ.js"
  },
  "/_nuxt/BzzwCNSt.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"8d-7jOk15jCztKiS+mGpNRBuUCqKr8\"",
    "mtime": "2026-06-14T09:29:09.585Z",
    "size": 141,
    "path": "../public/_nuxt/BzzwCNSt.js"
  },
  "/_nuxt/C-4Oazno.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1ea0-dpXYxUVn4CQ20H+uXK+vKRny2ks\"",
    "mtime": "2026-06-14T09:29:09.585Z",
    "size": 7840,
    "path": "../public/_nuxt/C-4Oazno.js"
  },
  "/_nuxt/C3WqSwOh.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"51a-v06mK8E3O6Nz8RVI8jSs0yu5MQE\"",
    "mtime": "2026-06-14T09:29:09.586Z",
    "size": 1306,
    "path": "../public/_nuxt/C3WqSwOh.js"
  },
  "/_nuxt/C5fjjcKz.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"26d9-tBT7UdqW32iFUkR8Bs/sSYpthts\"",
    "mtime": "2026-06-14T09:29:09.586Z",
    "size": 9945,
    "path": "../public/_nuxt/C5fjjcKz.js"
  },
  "/_nuxt/C6s4YI6l.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1ed3-LLoU6jYsRqL6lE9ZOVQpnwRnBh8\"",
    "mtime": "2026-06-14T09:29:09.586Z",
    "size": 7891,
    "path": "../public/_nuxt/C6s4YI6l.js"
  },
  "/_nuxt/CAag_PRF.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"3665-iFHJn+xXijWmsSlYGGdOi9NZjGw\"",
    "mtime": "2026-06-14T09:29:09.585Z",
    "size": 13925,
    "path": "../public/_nuxt/CAag_PRF.js"
  },
  "/_nuxt/CF5nATxz.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"166d-QEAbKzFTXTktggvMi1reXfDanNI\"",
    "mtime": "2026-06-14T09:29:09.586Z",
    "size": 5741,
    "path": "../public/_nuxt/CF5nATxz.js"
  },
  "/_nuxt/CGIPuhLi.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1c24-Qxz8cZmKV5GO9fQyoxbtfcpY+Q8\"",
    "mtime": "2026-06-14T09:29:09.587Z",
    "size": 7204,
    "path": "../public/_nuxt/CGIPuhLi.js"
  },
  "/_nuxt/CGr2hHf7.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"73-rmf/tM83R+pabL0ezWhMaWdX5Y4\"",
    "mtime": "2026-06-14T09:29:09.586Z",
    "size": 115,
    "path": "../public/_nuxt/CGr2hHf7.js"
  },
  "/_nuxt/CGuFwBVI.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"30e0-S/E371ecWQEUR6gYku62TLI/Kv0\"",
    "mtime": "2026-06-14T09:29:09.587Z",
    "size": 12512,
    "path": "../public/_nuxt/CGuFwBVI.js"
  },
  "/_nuxt/CHJ0O56i.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"277-v7t1z/LT2pop07o+c2vk5aQ3eTg\"",
    "mtime": "2026-06-14T09:29:09.588Z",
    "size": 631,
    "path": "../public/_nuxt/CHJ0O56i.js"
  },
  "/_nuxt/CIHoGVOa.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2a39-79IyEzNEYoLzbhZLVBiYnzNXFkg\"",
    "mtime": "2026-06-14T09:29:09.588Z",
    "size": 10809,
    "path": "../public/_nuxt/CIHoGVOa.js"
  },
  "/_nuxt/CJ_3IvT6.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"52d7-pY01U7hVv8h6DQ3vvBzh0kA1bUQ\"",
    "mtime": "2026-06-14T09:29:09.588Z",
    "size": 21207,
    "path": "../public/_nuxt/CJ_3IvT6.js"
  },
  "/_nuxt/CM8LO42l.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"81-2C8XE1DhYzehJFlimcnTU6C0kGU\"",
    "mtime": "2026-06-14T09:29:09.590Z",
    "size": 129,
    "path": "../public/_nuxt/CM8LO42l.js"
  },
  "/_nuxt/CMXzcO48.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"202a-jcZ2nG1R3oapinJrUHKBMzDWQic\"",
    "mtime": "2026-06-14T09:29:09.591Z",
    "size": 8234,
    "path": "../public/_nuxt/CMXzcO48.js"
  },
  "/_nuxt/CQAWWZv4.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"297b-+8HLX/Zm3tzx1jBlPniUcKyho3A\"",
    "mtime": "2026-06-14T09:29:09.590Z",
    "size": 10619,
    "path": "../public/_nuxt/CQAWWZv4.js"
  },
  "/_nuxt/CTwh3Poe.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"343-EqBXeFkD/QZyC16CoeIGBQC416U\"",
    "mtime": "2026-06-14T09:29:09.590Z",
    "size": 835,
    "path": "../public/_nuxt/CTwh3Poe.js"
  },
  "/_nuxt/CUFXrUjz.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"17154-aahRGbAyOvySzXNzTLhJqroz4HM\"",
    "mtime": "2026-06-14T09:29:09.591Z",
    "size": 94548,
    "path": "../public/_nuxt/CUFXrUjz.js"
  },
  "/_nuxt/CUndgRjY.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"196-v9Z5lCzykb+X8LLGOCqEw7ffUlo\"",
    "mtime": "2026-06-14T09:29:09.591Z",
    "size": 406,
    "path": "../public/_nuxt/CUndgRjY.js"
  },
  "/_nuxt/CXBNmlmj.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"e3f-v4UMGU9M3CrVoyLGg3csg5uBS8Q\"",
    "mtime": "2026-06-14T09:29:09.592Z",
    "size": 3647,
    "path": "../public/_nuxt/CXBNmlmj.js"
  },
  "/_nuxt/CY09zR-O.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"6029-yc8K//UmQL8oErYsKW8lPL0R+SY\"",
    "mtime": "2026-06-14T09:29:09.592Z",
    "size": 24617,
    "path": "../public/_nuxt/CY09zR-O.js"
  },
  "/_nuxt/CYm3WBQu.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"21c2-ZsZjnHJNUeXOSEDuvKXvGjBPFPY\"",
    "mtime": "2026-06-14T09:29:09.592Z",
    "size": 8642,
    "path": "../public/_nuxt/CYm3WBQu.js"
  },
  "/_nuxt/CbvP3fHg.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"260b-WDVfTZd78ZvyGvxNmFAqSNAp9Ds\"",
    "mtime": "2026-06-14T09:29:09.593Z",
    "size": 9739,
    "path": "../public/_nuxt/CbvP3fHg.js"
  },
  "/_nuxt/CgOTg-Ab.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"7958-Fg1/YDI6hOVPUdFpOlalHUC7PMY\"",
    "mtime": "2026-06-14T09:29:09.593Z",
    "size": 31064,
    "path": "../public/_nuxt/CgOTg-Ab.js"
  },
  "/_nuxt/CiG8jYj-.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"509-mZ95i7yD2MiyXuvniAzggpqaaR0\"",
    "mtime": "2026-06-14T09:29:09.593Z",
    "size": 1289,
    "path": "../public/_nuxt/CiG8jYj-.js"
  },
  "/_nuxt/Cj3Mz6S6.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1a8-omObGgrvVoeVO81Gl4g3BGgHDig\"",
    "mtime": "2026-06-14T09:29:09.593Z",
    "size": 424,
    "path": "../public/_nuxt/Cj3Mz6S6.js"
  },
  "/_nuxt/CjDbltun.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"480a-h2mNrDAnfQGKGmgGJL7p3y/NPFQ\"",
    "mtime": "2026-06-14T09:29:09.594Z",
    "size": 18442,
    "path": "../public/_nuxt/CjDbltun.js"
  },
  "/_nuxt/Co6yHwJg.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1961-yEbnNAEDpV5CdOgLvYPlD8uJG/g\"",
    "mtime": "2026-06-14T09:29:09.594Z",
    "size": 6497,
    "path": "../public/_nuxt/Co6yHwJg.js"
  },
  "/_nuxt/CoaUF789.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1f4-ByLvXw69ieTwBWZBkCGYzU+QxKY\"",
    "mtime": "2026-06-14T09:29:09.595Z",
    "size": 500,
    "path": "../public/_nuxt/CoaUF789.js"
  },
  "/_nuxt/Cp81p7Ue.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"dad4-4Wm36etUSfosD+0OUcJE4M0r/K4\"",
    "mtime": "2026-06-14T09:29:09.595Z",
    "size": 56020,
    "path": "../public/_nuxt/Cp81p7Ue.js"
  },
  "/_nuxt/CpgoGwED.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"bb2-ylrry2GPugmVtKkNCK+cbgKUC0o\"",
    "mtime": "2026-06-14T09:29:09.595Z",
    "size": 2994,
    "path": "../public/_nuxt/CpgoGwED.js"
  },
  "/_nuxt/CprGM58k.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"3818-tdBySltsK/s6AqY6Acovvtfl4jY\"",
    "mtime": "2026-06-14T09:29:09.595Z",
    "size": 14360,
    "path": "../public/_nuxt/CprGM58k.js"
  },
  "/_nuxt/CtIrxZNi.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"d54-xoJOf1dBbkzV3a0YDo4k2VQQRNA\"",
    "mtime": "2026-06-14T09:29:09.596Z",
    "size": 3412,
    "path": "../public/_nuxt/CtIrxZNi.js"
  },
  "/_nuxt/CxQ-gEGn.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1ad0-FId8WSarq49xuZVhkgBAmUK/4to\"",
    "mtime": "2026-06-14T09:29:09.596Z",
    "size": 6864,
    "path": "../public/_nuxt/CxQ-gEGn.js"
  },
  "/_nuxt/CzgwBud6.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"588a-fd9y0NFWAx4NEDl665Gu+6hEAjI\"",
    "mtime": "2026-06-14T09:29:09.597Z",
    "size": 22666,
    "path": "../public/_nuxt/CzgwBud6.js"
  },
  "/_nuxt/CzkzlT1V.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1b8b-jZ+T++gFva55KHVkw4LQC1Tng4E\"",
    "mtime": "2026-06-14T09:29:09.596Z",
    "size": 7051,
    "path": "../public/_nuxt/CzkzlT1V.js"
  },
  "/_nuxt/D33QIX7y.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"c345-GhhH4B6ZTZ/bt52CEqDVPBkHOO8\"",
    "mtime": "2026-06-14T09:29:09.598Z",
    "size": 49989,
    "path": "../public/_nuxt/D33QIX7y.js"
  },
  "/_nuxt/D4XAmjRX.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"44f-slNNmbm9WDPOLC6SMD2gzgKqmAs\"",
    "mtime": "2026-06-14T09:29:09.597Z",
    "size": 1103,
    "path": "../public/_nuxt/D4XAmjRX.js"
  },
  "/_nuxt/D7wbzdND.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4b3-zh4Lae7np7HkKRRmtb7YdyMc76I\"",
    "mtime": "2026-06-14T09:29:09.596Z",
    "size": 1203,
    "path": "../public/_nuxt/D7wbzdND.js"
  },
  "/_nuxt/D81Us46Z.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2243-pWCWgo2PFyp1bd3djuR5bxhyh3U\"",
    "mtime": "2026-06-14T09:29:09.597Z",
    "size": 8771,
    "path": "../public/_nuxt/D81Us46Z.js"
  },
  "/_nuxt/D8vZWwr2.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1512-uePaRduHEpZgv+YtgL6uLmR+ig8\"",
    "mtime": "2026-06-14T09:29:09.598Z",
    "size": 5394,
    "path": "../public/_nuxt/D8vZWwr2.js"
  },
  "/_nuxt/D9zP7iJN.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"71f-Ld5c3hT1JEzgXpy4130BMuS6TVs\"",
    "mtime": "2026-06-14T09:29:09.599Z",
    "size": 1823,
    "path": "../public/_nuxt/D9zP7iJN.js"
  },
  "/_nuxt/DEGJpwYi.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1e42-kAOxIzYLvvW3wq5jeQiLZGhNIVA\"",
    "mtime": "2026-06-14T09:29:09.599Z",
    "size": 7746,
    "path": "../public/_nuxt/DEGJpwYi.js"
  },
  "/_nuxt/DEoDX60x.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"420-uKQdHrlqo5jljo8ugVgZBvzqlx0\"",
    "mtime": "2026-06-14T09:29:09.599Z",
    "size": 1056,
    "path": "../public/_nuxt/DEoDX60x.js"
  },
  "/_nuxt/DFUnLHsf.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"186-48doHC8GXjVJ4o+/wvXWUnRzExU\"",
    "mtime": "2026-06-14T09:29:09.600Z",
    "size": 390,
    "path": "../public/_nuxt/DFUnLHsf.js"
  },
  "/_nuxt/DG3JoL_b.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"fdd-Kp1q2MeNVCUtxhXlex1lxC6K0BU\"",
    "mtime": "2026-06-14T09:29:09.600Z",
    "size": 4061,
    "path": "../public/_nuxt/DG3JoL_b.js"
  },
  "/_nuxt/DGsuBzem.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"7acd-JAdlzXxDzomYaMYgMG1Kb3L4mXk\"",
    "mtime": "2026-06-14T09:29:09.601Z",
    "size": 31437,
    "path": "../public/_nuxt/DGsuBzem.js"
  },
  "/_nuxt/DHaV7T_C.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"ecb-c3pvVHoYQPLN48e+kk8XoAvIiP4\"",
    "mtime": "2026-06-14T09:29:09.603Z",
    "size": 3787,
    "path": "../public/_nuxt/DHaV7T_C.js"
  },
  "/_nuxt/DJ5tTi3X.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1620-lrimYiYzD/8ycR5qo9U5NEfliSo\"",
    "mtime": "2026-06-14T09:29:09.600Z",
    "size": 5664,
    "path": "../public/_nuxt/DJ5tTi3X.js"
  },
  "/_nuxt/DLvsLvt-.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2451-A+nRGEtKivgHPZ/GWwUd9UrK7xE\"",
    "mtime": "2026-06-14T09:29:09.601Z",
    "size": 9297,
    "path": "../public/_nuxt/DLvsLvt-.js"
  },
  "/_nuxt/DMTo9cht.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"37ea-5jB53d42giriYCJDky7aeZ7/VHA\"",
    "mtime": "2026-06-14T09:29:09.601Z",
    "size": 14314,
    "path": "../public/_nuxt/DMTo9cht.js"
  },
  "/_nuxt/DO8alW5h.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"547a-YrASWQRHey3ImGVU/nLl3uvtd9o\"",
    "mtime": "2026-06-14T09:29:09.601Z",
    "size": 21626,
    "path": "../public/_nuxt/DO8alW5h.js"
  },
  "/_nuxt/DQAJM8Ob.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"151-qGFYV+5GWZtUWM3jXHMYc7U+Va8\"",
    "mtime": "2026-06-14T09:29:09.602Z",
    "size": 337,
    "path": "../public/_nuxt/DQAJM8Ob.js"
  },
  "/_nuxt/DQDQW_eH.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"225a-Waimas4CIWSWHh2qBkiskMnU++U\"",
    "mtime": "2026-06-14T09:29:09.602Z",
    "size": 8794,
    "path": "../public/_nuxt/DQDQW_eH.js"
  },
  "/_nuxt/DVGuClv8.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4137-WZa+TL4VADRSTD2/gpfAFkwfnYs\"",
    "mtime": "2026-06-14T09:29:09.602Z",
    "size": 16695,
    "path": "../public/_nuxt/DVGuClv8.js"
  },
  "/_nuxt/DVuKBG_m.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4068-gTM5zPYZk4jdVRMqO2oUZQnjWWc\"",
    "mtime": "2026-06-14T09:29:09.602Z",
    "size": 16488,
    "path": "../public/_nuxt/DVuKBG_m.js"
  },
  "/_nuxt/DXCbjL9u.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4514-+G5RKsmc2Sr3PVT5eqyEL2Di9ZA\"",
    "mtime": "2026-06-14T09:29:09.603Z",
    "size": 17684,
    "path": "../public/_nuxt/DXCbjL9u.js"
  },
  "/_nuxt/DYSVe14K.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4f45-XjdYBNaR06iwG54Mwqqv/LJE8BE\"",
    "mtime": "2026-06-14T09:29:09.603Z",
    "size": 20293,
    "path": "../public/_nuxt/DYSVe14K.js"
  },
  "/_nuxt/D__gE1zc.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"98-Y9T1Y3S/NAmBhg0cbMsK4N1Z/NU\"",
    "mtime": "2026-06-14T09:29:09.604Z",
    "size": 152,
    "path": "../public/_nuxt/D__gE1zc.js"
  },
  "/_nuxt/DcofbVqb.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"511e-NsjRe4gy+iR5DW304BMsekRHIAU\"",
    "mtime": "2026-06-14T09:29:09.603Z",
    "size": 20766,
    "path": "../public/_nuxt/DcofbVqb.js"
  },
  "/_nuxt/DdYo3Ixy.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"18e-mA5EE60DhSHO/P5J1nYKL8HG8ew\"",
    "mtime": "2026-06-14T09:29:09.603Z",
    "size": 398,
    "path": "../public/_nuxt/DdYo3Ixy.js"
  },
  "/_nuxt/Dd_fbjQY.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"e22-SNL/SttRP0pDkLGdOsGE7teks8g\"",
    "mtime": "2026-06-14T09:29:09.603Z",
    "size": 3618,
    "path": "../public/_nuxt/Dd_fbjQY.js"
  },
  "/_nuxt/DdfQpC4g.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"7d37-Xqh6IrfwfOr7Oj84gEO6abM917w\"",
    "mtime": "2026-06-14T09:29:09.604Z",
    "size": 32055,
    "path": "../public/_nuxt/DdfQpC4g.js"
  },
  "/_nuxt/DeC0XAZi.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"349e-qdeqWRYfA8/UVbKEMJJj3mSDOh0\"",
    "mtime": "2026-06-14T09:29:09.604Z",
    "size": 13470,
    "path": "../public/_nuxt/DeC0XAZi.js"
  },
  "/_nuxt/Dfqi6zXz.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"12098-Zu4R1g/Vuub1vJj/CgH4m3mtK3c\"",
    "mtime": "2026-06-14T09:29:09.605Z",
    "size": 73880,
    "path": "../public/_nuxt/Dfqi6zXz.js"
  },
  "/_nuxt/DlOKbiQh.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"10e-psLgtmQUXeLwN171qlgRfLYJ07Q\"",
    "mtime": "2026-06-14T09:29:09.604Z",
    "size": 270,
    "path": "../public/_nuxt/DlOKbiQh.js"
  },
  "/_nuxt/Dlfwo7N5.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2ede-shUr+DBp/GOH1KjuC78PCBqbixk\"",
    "mtime": "2026-06-14T09:29:09.605Z",
    "size": 11998,
    "path": "../public/_nuxt/Dlfwo7N5.js"
  },
  "/_nuxt/DmHPR5lg.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4baec-LO1qT18Bnfj7WIj5TTzGnTf1jQo\"",
    "mtime": "2026-06-14T09:29:09.610Z",
    "size": 309996,
    "path": "../public/_nuxt/DmHPR5lg.js"
  },
  "/_nuxt/DoROlVvH.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"7051-/TC7cAfA7lus/wqTk2/2KESqF24\"",
    "mtime": "2026-06-14T09:29:09.606Z",
    "size": 28753,
    "path": "../public/_nuxt/DoROlVvH.js"
  },
  "/_nuxt/DpCb2j92.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"500-y/nh9U0hYF5p2WBZXHfLB13M60U\"",
    "mtime": "2026-06-14T09:29:09.605Z",
    "size": 1280,
    "path": "../public/_nuxt/DpCb2j92.js"
  },
  "/_nuxt/DqVjZ2ke.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"641-QRoBnjbJpPHZAGZErZ/6hsbNSrg\"",
    "mtime": "2026-06-14T09:29:09.605Z",
    "size": 1601,
    "path": "../public/_nuxt/DqVjZ2ke.js"
  },
  "/_nuxt/DrM87O-9.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"79fa-a4TMHLjlIxW+rHfCAf+xByh09Lk\"",
    "mtime": "2026-06-14T09:29:09.608Z",
    "size": 31226,
    "path": "../public/_nuxt/DrM87O-9.js"
  },
  "/_nuxt/Dwg6oJu7.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"8d1-parWx8VOZDqJ7vQ+uCAqTAlqb64\"",
    "mtime": "2026-06-14T09:29:09.607Z",
    "size": 2257,
    "path": "../public/_nuxt/Dwg6oJu7.js"
  },
  "/_nuxt/DwrC6KZY.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"3a1a-ooIkXvk2uL6TxCVq6MDLhsGGeMU\"",
    "mtime": "2026-06-14T09:29:09.607Z",
    "size": 14874,
    "path": "../public/_nuxt/DwrC6KZY.js"
  },
  "/_nuxt/Dy07CK_X.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"29cf-uh5C9Y5IdaCRU2Ati9JyhLqKCBg\"",
    "mtime": "2026-06-14T09:29:09.608Z",
    "size": 10703,
    "path": "../public/_nuxt/Dy07CK_X.js"
  },
  "/_nuxt/Dy_3lkQY.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"368e-pR9V1F3DnHTne2ugXQEwT8fMlIs\"",
    "mtime": "2026-06-14T09:29:09.608Z",
    "size": 13966,
    "path": "../public/_nuxt/Dy_3lkQY.js"
  },
  "/_nuxt/EK524y10.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"67be-sT+cE/ogIie0fs1nS7z+uQ0Bdk0\"",
    "mtime": "2026-06-14T09:29:09.610Z",
    "size": 26558,
    "path": "../public/_nuxt/EK524y10.js"
  },
  "/_nuxt/F4SGeL1Q.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"6a-qdfs8JfquhrXTUtR+bxr+e8i3GE\"",
    "mtime": "2026-06-14T09:29:09.609Z",
    "size": 106,
    "path": "../public/_nuxt/F4SGeL1Q.js"
  },
  "/_nuxt/FiraCode-Regular.CRwVj4V2.woff2": {
    "type": "font/woff2",
    "etag": "\"1bf18-zeT20kIaKCztvuDOx7ICixqHLeA\"",
    "mtime": "2026-06-14T09:29:09.610Z",
    "size": 114456,
    "path": "../public/_nuxt/FiraCode-Regular.CRwVj4V2.woff2"
  },
  "/_nuxt/HCHZ96f3.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"991-/8ZMCKoMp/Dod/d+oRdZ47xSBwk\"",
    "mtime": "2026-06-14T09:29:09.610Z",
    "size": 2449,
    "path": "../public/_nuxt/HCHZ96f3.js"
  },
  "/_nuxt/Hg3Q8c_0.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"182-4npj7UxWcWb3OSorai9O869+ZhQ\"",
    "mtime": "2026-06-14T09:29:09.610Z",
    "size": 386,
    "path": "../public/_nuxt/Hg3Q8c_0.js"
  },
  "/_nuxt/KJZP1L2-.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1eee-QpOZchhsLQgfhl83kjBC5bRAJL8\"",
    "mtime": "2026-06-14T09:29:09.610Z",
    "size": 7918,
    "path": "../public/_nuxt/KJZP1L2-.js"
  },
  "/_nuxt/LatoLatin-Regular.Dmlz1U0B.woff2": {
    "type": "font/woff2",
    "etag": "\"aaf0-VwLLIAV2a5z8Oyr/1z0b/j9PbB4\"",
    "mtime": "2026-06-14T09:29:09.611Z",
    "size": 43760,
    "path": "../public/_nuxt/LatoLatin-Regular.Dmlz1U0B.woff2"
  },
  "/_nuxt/LatoLatin-Semibold.Dbk81p2D.woff2": {
    "type": "font/woff2",
    "etag": "\"ad44-D+QkFKhuev4g+FaJZVKHD2bKXOg\"",
    "mtime": "2026-06-14T09:29:09.612Z",
    "size": 44356,
    "path": "../public/_nuxt/LatoLatin-Semibold.Dbk81p2D.woff2"
  },
  "/_nuxt/QpI9WcJO.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1047-Q1kIEeCqYeMLTp1n0HbTxkAcFnQ\"",
    "mtime": "2026-06-14T09:29:09.611Z",
    "size": 4167,
    "path": "../public/_nuxt/QpI9WcJO.js"
  },
  "/_nuxt/SQnvgvv6.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"8d0-eT7C2b4MSHV2aExvChqveKj90vc\"",
    "mtime": "2026-06-14T09:29:09.611Z",
    "size": 2256,
    "path": "../public/_nuxt/SQnvgvv6.js"
  },
  "/_nuxt/SRgnb1qh.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"245e-alQw3eTMN0LgPBDVxNu1Jr5PI5A\"",
    "mtime": "2026-06-14T09:29:09.613Z",
    "size": 9310,
    "path": "../public/_nuxt/SRgnb1qh.js"
  },
  "/_nuxt/SystemConfigCacheInfo.CIzgsgWP.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"993-HuV3OsX50xAnHGQd5aez6v8umZc\"",
    "mtime": "2026-06-14T09:29:09.612Z",
    "size": 2451,
    "path": "../public/_nuxt/SystemConfigCacheInfo.CIzgsgWP.css"
  },
  "/_nuxt/TRS6STvs.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"ee-eiqC+2erZRYySRpBtCjx8oVaYHo\"",
    "mtime": "2026-06-14T09:29:09.612Z",
    "size": 238,
    "path": "../public/_nuxt/TRS6STvs.js"
  },
  "/_nuxt/WmWtXwec.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1977-etZwN7mBvPCYZSRpJpS/mzHGxrU\"",
    "mtime": "2026-06-14T09:29:09.612Z",
    "size": 6519,
    "path": "../public/_nuxt/WmWtXwec.js"
  },
  "/_nuxt/_key_.DSBzaFa_.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"76c-5yL0v8fmMg3mavrBl57crBcXEDs\"",
    "mtime": "2026-06-14T09:29:09.612Z",
    "size": 1900,
    "path": "../public/_nuxt/_key_.DSBzaFa_.css"
  },
  "/_nuxt/accounts-expansion.CnVVVlsS.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"75-+xf7MUwvgo7SgBtjytvWLqqv6Ys\"",
    "mtime": "2026-06-14T09:29:09.613Z",
    "size": 117,
    "path": "../public/_nuxt/accounts-expansion.CnVVVlsS.css"
  },
  "/_nuxt/accounts.B4OU05cJ.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"e7-NvQaBz9ZiCnBomV6bUPn0crn5rc\"",
    "mtime": "2026-06-14T09:29:09.613Z",
    "size": 231,
    "path": "../public/_nuxt/accounts.B4OU05cJ.css"
  },
  "/_nuxt/add-resource.Dy0fCb-G.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"d4d-biMRDBQmO5vAHqafhIjd6C69RIo\"",
    "mtime": "2026-06-14T09:29:09.613Z",
    "size": 3405,
    "path": "../public/_nuxt/add-resource.Dy0fCb-G.css"
  },
  "/_nuxt/admin.C_FoEPl0.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"81-B5AZbNKNzYXjJkno7Jfb+x5Sqro\"",
    "mtime": "2026-06-14T09:29:09.613Z",
    "size": 129,
    "path": "../public/_nuxt/admin.C_FoEPl0.css"
  },
  "/_nuxt/api-access-logs.DR0svuo8.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"a9-TJGr031oR25ndSYNgwaxCC+4cnY\"",
    "mtime": "2026-06-14T09:29:09.613Z",
    "size": 169,
    "path": "../public/_nuxt/api-access-logs.DR0svuo8.css"
  },
  "/_nuxt/api-docs.a52KE7Hc.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"10c-qgul2la1dSRN0xwsBELmrXlRKPA\"",
    "mtime": "2026-06-14T09:29:09.613Z",
    "size": 268,
    "path": "../public/_nuxt/api-docs.a52KE7Hc.css"
  },
  "/_nuxt/banner.B_89Gk1w.webp": {
    "type": "image/webp",
    "etag": "\"5384-k6IBidwk5vcRDtI29L5H1FrntoY\"",
    "mtime": "2026-06-14T09:29:09.614Z",
    "size": 21380,
    "path": "../public/_nuxt/banner.B_89Gk1w.webp"
  },
  "/_nuxt/bmQVJIR_.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"10744-8m3SqVX1d5RjTJJjMqolDEYGBzo\"",
    "mtime": "2026-06-14T09:29:09.615Z",
    "size": 67396,
    "path": "../public/_nuxt/bmQVJIR_.js"
  },
  "/_nuxt/bot.WqpuTKa5.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"17a-UrHg8pzMqHh3MT9j1X4VLJR1bGg\"",
    "mtime": "2026-06-14T09:29:09.614Z",
    "size": 378,
    "path": "../public/_nuxt/bot.WqpuTKa5.css"
  },
  "/_nuxt/data-transfer.CvbLi8NR.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"26d-n72d4SCRJiR5PyLmAwojG/fHLHI\"",
    "mtime": "2026-06-14T09:29:09.614Z",
    "size": 621,
    "path": "../public/_nuxt/data-transfer.CvbLi8NR.css"
  },
  "/_nuxt/dev-config.CWxgQc5m.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"ad-+en5obCqnIu+5WHrTABPsMK0THk\"",
    "mtime": "2026-06-14T09:29:09.615Z",
    "size": 173,
    "path": "../public/_nuxt/dev-config.CWxgQc5m.css"
  },
  "/_nuxt/error-404.Bduj4t2m.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"de4-62p0OW5S0ek57g2sQ/ogW4J5T8Y\"",
    "mtime": "2026-06-14T09:29:09.615Z",
    "size": 3556,
    "path": "../public/_nuxt/error-404.Bduj4t2m.css"
  },
  "/_nuxt/error-500.D4Hrqg63.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"75c-BB9Fr1bRyPBJQATaxcOS1m2/ES0\"",
    "mtime": "2026-06-14T09:29:09.615Z",
    "size": 1884,
    "path": "../public/_nuxt/error-500.D4Hrqg63.css"
  },
  "/_nuxt/fa-brands-400.D1LuMI3I.ttf": {
    "type": "font/ttf",
    "etag": "\"33768-2m9TfNWXt20mSXDxyV5L7ryC/ek\"",
    "mtime": "2026-06-14T09:29:09.618Z",
    "size": 210792,
    "path": "../public/_nuxt/fa-brands-400.D1LuMI3I.ttf"
  },
  "/_nuxt/fa-brands-400.D_cYUPeE.woff2": {
    "type": "font/woff2",
    "etag": "\"1cf9c-McMDUQhOJnC75cR700LhFUCmhxc\"",
    "mtime": "2026-06-14T09:29:09.617Z",
    "size": 118684,
    "path": "../public/_nuxt/fa-brands-400.D_cYUPeE.woff2"
  },
  "/_nuxt/fa-regular-400.BjRzuEpd.woff2": {
    "type": "font/woff2",
    "etag": "\"6380-6bb1uEK+In4kI3M/65Yd1suAQGY\"",
    "mtime": "2026-06-14T09:29:09.616Z",
    "size": 25472,
    "path": "../public/_nuxt/fa-regular-400.BjRzuEpd.woff2"
  },
  "/_nuxt/fa-regular-400.DZaxPHgR.ttf": {
    "type": "font/ttf",
    "etag": "\"109e0-bn8GUQWe6peD6osoiSssl5wvfzA\"",
    "mtime": "2026-06-14T09:29:09.617Z",
    "size": 68064,
    "path": "../public/_nuxt/fa-regular-400.DZaxPHgR.ttf"
  },
  "/_nuxt/fa-solid-900.CTAAxXor.woff2": {
    "type": "font/woff2",
    "etag": "\"26a0c-fNiv1lAZYv2jXWbw5MO4gVrEcdg\"",
    "mtime": "2026-06-14T09:29:09.618Z",
    "size": 158220,
    "path": "../public/_nuxt/fa-solid-900.CTAAxXor.woff2"
  },
  "/_nuxt/fa-solid-900.D0aA9rwL.ttf": {
    "type": "font/ttf",
    "etag": "\"68080-eAqP8W67aB+F1dMPrJTr6EWhK14\"",
    "mtime": "2026-06-14T09:29:09.621Z",
    "size": 426112,
    "path": "../public/_nuxt/fa-solid-900.D0aA9rwL.ttf"
  },
  "/_nuxt/fa-v4compatibility.C9RhG_FT.woff2": {
    "type": "font/woff2",
    "etag": "\"12bc-T9gvOSN6PlLLTjEdGWG7tHNToUc\"",
    "mtime": "2026-06-14T09:29:09.619Z",
    "size": 4796,
    "path": "../public/_nuxt/fa-v4compatibility.C9RhG_FT.woff2"
  },
  "/_nuxt/fa-v4compatibility.CCth-dXg.ttf": {
    "type": "font/ttf",
    "etag": "\"2a54-rUFfmaq02b8X5SP0hKN3+DnKjDw\"",
    "mtime": "2026-06-14T09:29:09.618Z",
    "size": 10836,
    "path": "../public/_nuxt/fa-v4compatibility.CCth-dXg.ttf"
  },
  "/_nuxt/failed-resources.DxZDdM3b.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"e7-78ZemuwbnQrhiicaQwzW0P0Sov4\"",
    "mtime": "2026-06-14T09:29:09.619Z",
    "size": 231,
    "path": "../public/_nuxt/failed-resources.DxZDdM3b.css"
  },
  "/_nuxt/feature-config.C0Onu6Hy.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"113-J3VOnyvY8kOQOrVVA6wZmgAwqEg\"",
    "mtime": "2026-06-14T09:29:09.619Z",
    "size": 275,
    "path": "../public/_nuxt/feature-config.C0Onu6Hy.css"
  },
  "/_nuxt/files.CxX4u1ik.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"6d6-oNhDs5Zff0rn6KEhy0mO+081Mb0\"",
    "mtime": "2026-06-14T09:29:09.620Z",
    "size": 1750,
    "path": "../public/_nuxt/files.CxX4u1ik.css"
  },
  "/_nuxt/g0YHQayI.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"6b3-Kx6gj43uOTBm8LwGseoGrecbnAM\"",
    "mtime": "2026-06-14T09:29:09.620Z",
    "size": 1715,
    "path": "../public/_nuxt/g0YHQayI.js"
  },
  "/_nuxt/hot-dramas.BAxk1nNo.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"173-kRxcyzGbu4f6tgs4iVP03tgDfno\"",
    "mtime": "2026-06-14T09:29:09.620Z",
    "size": 371,
    "path": "../public/_nuxt/hot-dramas.BAxk1nNo.css"
  },
  "/_nuxt/hot-dramas.D0v_MZaI.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"75-/vgGYoWwEgHjRz5jQWamxUyhNTg\"",
    "mtime": "2026-06-14T09:29:09.622Z",
    "size": 117,
    "path": "../public/_nuxt/hot-dramas.D0v_MZaI.css"
  },
  "/_nuxt/index.Bk_KjbN9.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"48-sApsRydHRVBNLUE5SOnP5Wl+DSM\"",
    "mtime": "2026-06-14T09:29:09.621Z",
    "size": 72,
    "path": "../public/_nuxt/index.Bk_KjbN9.css"
  },
  "/_nuxt/index.CSFv2LIP.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"ca9-fK+dpJtR9yCqR+/fEfxC6dalnNM\"",
    "mtime": "2026-06-14T09:29:09.623Z",
    "size": 3241,
    "path": "../public/_nuxt/index.CSFv2LIP.css"
  },
  "/_nuxt/index.C_p2KSXE.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"455-RsKABPu9WgLc7J3o66/KDbqw+UQ\"",
    "mtime": "2026-06-14T09:29:09.621Z",
    "size": 1109,
    "path": "../public/_nuxt/index.C_p2KSXE.css"
  },
  "/_nuxt/index.TzbN9O-b.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"48-UYNSxz90QU7B/w5J0Axs6ftKrrg\"",
    "mtime": "2026-06-14T09:29:09.622Z",
    "size": 72,
    "path": "../public/_nuxt/index.TzbN9O-b.css"
  },
  "/_nuxt/login.DSwSDMuf.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"3bc-1/Ce0bovr/iaKMc6UZXRmuGQR4I\"",
    "mtime": "2026-06-14T09:29:09.632Z",
    "size": 956,
    "path": "../public/_nuxt/login.DSwSDMuf.css"
  },
  "/_nuxt/mKdXyBVO.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"1dc9-QySJwaVqvY89IQPmC50FnyfxCig\"",
    "mtime": "2026-06-14T09:29:09.623Z",
    "size": 7625,
    "path": "../public/_nuxt/mKdXyBVO.js"
  },
  "/_nuxt/monitor.CC-9iO6V.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"8d-8zNnPFWchQQrzPzXteG1CTEP6MU\"",
    "mtime": "2026-06-14T09:29:09.623Z",
    "size": 141,
    "path": "../public/_nuxt/monitor.CC-9iO6V.css"
  },
  "/_nuxt/pWxRiGYn.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"48cf-RzMhTIBvbAHk/tQVpaL835wRLPo\"",
    "mtime": "2026-06-14T09:29:09.628Z",
    "size": 18639,
    "path": "../public/_nuxt/pWxRiGYn.js"
  },
  "/_nuxt/plugins.73sF14GG.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"1dc-fy/gmYn3UHZZAEiRWebZRzrkLSA\"",
    "mtime": "2026-06-14T09:29:09.623Z",
    "size": 476,
    "path": "../public/_nuxt/plugins.73sF14GG.css"
  },
  "/_nuxt/register.CO8q1D56.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"4b7-r3xn/TI42QCWtmlrloV3bVmMOw8\"",
    "mtime": "2026-06-14T09:29:09.626Z",
    "size": 1207,
    "path": "../public/_nuxt/register.CO8q1D56.css"
  },
  "/_nuxt/resources.u0uIrQ1b.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"e7-RzKOtNk2DTb0N+ButfbW0qxtCZo\"",
    "mtime": "2026-06-14T09:29:09.625Z",
    "size": 231,
    "path": "../public/_nuxt/resources.u0uIrQ1b.css"
  },
  "/_nuxt/sJus4ovQ.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"c9-AGtLIjmph8BjMeEDU6lvXGviPDA\"",
    "mtime": "2026-06-14T09:29:09.625Z",
    "size": 201,
    "path": "../public/_nuxt/sJus4ovQ.js"
  },
  "/_nuxt/seo.Ct71C_v5.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"163-/6ehy2FZVIFzGPoKo/j6IT3A4n0\"",
    "mtime": "2026-06-14T09:29:09.627Z",
    "size": 355,
    "path": "../public/_nuxt/seo.Ct71C_v5.css"
  },
  "/_nuxt/single.DeW7AkxW.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"82-ceXkUPbRp+GtSISu10Z7tGzepdc\"",
    "mtime": "2026-06-14T09:29:09.628Z",
    "size": 130,
    "path": "../public/_nuxt/single.DeW7AkxW.css"
  },
  "/_nuxt/site-config.CyK2VX74.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"d7b-Vg0SK7X7OAWu/obw+1Xr5IZIgKk\"",
    "mtime": "2026-06-14T09:29:09.627Z",
    "size": 3451,
    "path": "../public/_nuxt/site-config.CyK2VX74.css"
  },
  "/_nuxt/system-logs.BdlPZoNU.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"8e-+lVXatpHGsGVUXkBt2AlbsnZEi8\"",
    "mtime": "2026-06-14T09:29:09.630Z",
    "size": 142,
    "path": "../public/_nuxt/system-logs.BdlPZoNU.css"
  },
  "/_nuxt/user.CYUCX3xV.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"74-2UbREVTw9zT1C5dkAy5RpaI3DnY\"",
    "mtime": "2026-06-14T09:29:09.630Z",
    "size": 116,
    "path": "../public/_nuxt/user.CYUCX3xV.css"
  },
  "/_nuxt/users.ByzDyNwU.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"ad-+8wS5cgX9Oph7qyJTzp10i4Uxd8\"",
    "mtime": "2026-06-14T09:29:09.633Z",
    "size": 173,
    "path": "../public/_nuxt/users.ByzDyNwU.css"
  },
  "/_nuxt/usi_DtmX.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2786-NVLkYIgCgnJ2DlEsCEX/zUSTQ2c\"",
    "mtime": "2026-06-14T09:29:09.632Z",
    "size": 10118,
    "path": "../public/_nuxt/usi_DtmX.js"
  },
  "/_nuxt/version.W2N_gfdH.css": {
    "type": "text/css; charset=utf-8",
    "etag": "\"d4-G9OIb3TyS9BX97EwHLNvRffNcVk\"",
    "mtime": "2026-06-14T09:29:09.633Z",
    "size": 212,
    "path": "../public/_nuxt/version.W2N_gfdH.css"
  },
  "/_nuxt/wJIo8u48.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"a38-hVHJ/zeQN7qcUc6txHUjwOWqKvk\"",
    "mtime": "2026-06-14T09:29:09.633Z",
    "size": 2616,
    "path": "../public/_nuxt/wJIo8u48.js"
  },
  "/_nuxt/wlCwR3hZ.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"4c87-d7nuaPI4/AP28qPQCZ86Bl9sjlg\"",
    "mtime": "2026-06-14T09:29:09.634Z",
    "size": 19591,
    "path": "../public/_nuxt/wlCwR3hZ.js"
  },
  "/_nuxt/x_rD_Ya3.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"2b-NQFdZyJZYCba7NdYw0SoR9PvUtA\"",
    "mtime": "2026-06-14T09:29:09.634Z",
    "size": 43,
    "path": "../public/_nuxt/x_rD_Ya3.js"
  },
  "/_nuxt/yb1cLV1l.js": {
    "type": "text/javascript; charset=utf-8",
    "etag": "\"32999-Pa+/A14KcVLNP06jmPf2eHBkxOo\"",
    "mtime": "2026-06-14T09:29:09.635Z",
    "size": 207257,
    "path": "../public/_nuxt/yb1cLV1l.js"
  },
  "/_nuxt/builds/latest.json": {
    "type": "application/json",
    "etag": "\"47-M5Ht5xDNOpp4npX25H4U1Tc+zFo\"",
    "mtime": "2026-06-14T09:29:09.499Z",
    "size": 71,
    "path": "../public/_nuxt/builds/latest.json"
  },
  "/assets/svg/empty.svg": {
    "type": "image/svg+xml",
    "etag": "\"7d67-bIfPbOh6t7bvS2UJStkVBWUThP4\"",
    "mtime": "2026-06-14T09:29:09.647Z",
    "size": 32103,
    "path": "../public/assets/svg/empty.svg"
  },
  "/assets/svg/forbidden.svg": {
    "type": "image/svg+xml",
    "etag": "\"475b-kUZSNbDttfViWs8uJVpfuzUDHjM\"",
    "mtime": "2026-06-14T09:29:09.649Z",
    "size": 18267,
    "path": "../public/assets/svg/forbidden.svg"
  },
  "/assets/svg/qrcode.svg": {
    "type": "image/svg+xml",
    "etag": "\"89b1-GwFMTc1HMS0Gg+ionY8sLXHBbjI\"",
    "mtime": "2026-06-14T09:29:09.652Z",
    "size": 35249,
    "path": "../public/assets/svg/qrcode.svg"
  },
  "/assets/svg/risk.svg": {
    "type": "image/svg+xml",
    "etag": "\"5904-P98w4S2lvJL/q6IrdMb+Uz2N6K8\"",
    "mtime": "2026-06-14T09:29:09.648Z",
    "size": 22788,
    "path": "../public/assets/svg/risk.svg"
  },
  "/assets/images/3dian.png": {
    "type": "image/png",
    "etag": "\"1e1-a9HWaBQgkH2Xg8fH6gs7qetjvO4\"",
    "mtime": "2026-06-14T09:29:09.647Z",
    "size": 481,
    "path": "../public/assets/images/3dian.png"
  },
  "/assets/images/android.png": {
    "type": "image/png",
    "etag": "\"12d0-blPN/4/FoAoGWfuxbSjlJZbdYD0\"",
    "mtime": "2026-06-14T09:29:09.650Z",
    "size": 4816,
    "path": "../public/assets/images/android.png"
  },
  "/assets/images/banner.png": {
    "type": "image/png",
    "etag": "\"7b67-Xl6InP5EAS15gbdonRoX95goOrI\"",
    "mtime": "2026-06-14T09:29:09.649Z",
    "size": 31591,
    "path": "../public/assets/images/banner.png"
  },
  "/assets/images/banner.webp": {
    "type": "image/webp",
    "etag": "\"5384-k6IBidwk5vcRDtI29L5H1FrntoY\"",
    "mtime": "2026-06-14T09:29:09.647Z",
    "size": 21380,
    "path": "../public/assets/images/banner.webp"
  },
  "/assets/images/cover1.webp": {
    "type": "image/webp",
    "etag": "\"19da-I1WRUwCLCCSNFs57oGFxeAy8DOw\"",
    "mtime": "2026-06-14T09:29:09.647Z",
    "size": 6618,
    "path": "../public/assets/images/cover1.webp"
  },
  "/assets/images/cover2.webp": {
    "type": "image/webp",
    "etag": "\"20d0-z/FtoNdSpTEZ+zkAE7Bv6K2W0LY\"",
    "mtime": "2026-06-14T09:29:09.647Z",
    "size": 8400,
    "path": "../public/assets/images/cover2.webp"
  },
  "/assets/images/cover3.webp": {
    "type": "image/webp",
    "etag": "\"1af6-XBrx1DzJXBCPINxatpr/G3Frnq8\"",
    "mtime": "2026-06-14T09:29:09.648Z",
    "size": 6902,
    "path": "../public/assets/images/cover3.webp"
  },
  "/assets/images/cover4.webp": {
    "type": "image/webp",
    "etag": "\"266c-L8R5tJx/daHp7F/lfj5LH1Rk1Qk\"",
    "mtime": "2026-06-14T09:29:09.648Z",
    "size": 9836,
    "path": "../public/assets/images/cover4.webp"
  },
  "/assets/images/cover5.webp": {
    "type": "image/webp",
    "etag": "\"19d4-DJahR7gl1Xzu7Tym5Mkz8kGAPiA\"",
    "mtime": "2026-06-14T09:29:09.651Z",
    "size": 6612,
    "path": "../public/assets/images/cover5.webp"
  },
  "/assets/images/cover6.webp": {
    "type": "image/webp",
    "etag": "\"19f8-H1ogunAqBXv4uFGkz/VRagyzTCQ\"",
    "mtime": "2026-06-14T09:29:09.648Z",
    "size": 6648,
    "path": "../public/assets/images/cover6.webp"
  },
  "/assets/images/cover7.webp": {
    "type": "image/webp",
    "etag": "\"1b84-61n+eReJgH7ca+hBxR/yF92ih48\"",
    "mtime": "2026-06-14T09:29:09.649Z",
    "size": 7044,
    "path": "../public/assets/images/cover7.webp"
  },
  "/assets/images/cover8.webp": {
    "type": "image/webp",
    "etag": "\"172e-YbdavRqsm/zO5WnjZEF4rrPDWJw\"",
    "mtime": "2026-06-14T09:29:09.651Z",
    "size": 5934,
    "path": "../public/assets/images/cover8.webp"
  },
  "/assets/images/footer-banner.webp": {
    "type": "image/webp",
    "etag": "\"9ae-JEtyDR6/OZYNiBOpx7mpmS+zUWc\"",
    "mtime": "2026-06-14T09:29:09.649Z",
    "size": 2478,
    "path": "../public/assets/images/footer-banner.webp"
  },
  "/assets/images/fxt.jpg": {
    "type": "image/jpeg",
    "etag": "\"26549-Hn16yGkqdw0CcC+0k94GjRB3jVI\"",
    "mtime": "2026-06-14T09:29:09.651Z",
    "size": 157001,
    "path": "../public/assets/images/fxt.jpg"
  },
  "/assets/images/iphone.png": {
    "type": "image/png",
    "etag": "\"df6-Gfm/WdGUBckmEDiQVBHMr1Wq5v8\"",
    "mtime": "2026-06-14T09:29:09.650Z",
    "size": 3574,
    "path": "../public/assets/images/iphone.png"
  },
  "/assets/images/logo.png": {
    "type": "image/png",
    "etag": "\"35429-i9JxabxCBiyJzsNvEDWqKv2RYHo\"",
    "mtime": "2026-06-14T09:29:09.651Z",
    "size": 218153,
    "path": "../public/assets/images/logo.png"
  },
  "/assets/images/logo.webp": {
    "type": "image/webp",
    "etag": "\"3b64-fuUpGAD5R/4cyBbN8gCgJFv8Jc0\"",
    "mtime": "2026-06-14T09:29:09.649Z",
    "size": 15204,
    "path": "../public/assets/images/logo.webp"
  },
  "/assets/images/mp.webp": {
    "type": "image/webp",
    "etag": "\"ab7e-zM1T5DiRwafNaZsXcWGmk0SjrNI\"",
    "mtime": "2026-06-14T09:29:09.650Z",
    "size": 43902,
    "path": "../public/assets/images/mp.webp"
  },
  "/assets/images/og.webp": {
    "type": "image/webp",
    "etag": "\"88ea-tgZUsARNHY04Kmh50vfO3mit2Pc\"",
    "mtime": "2026-06-14T09:29:09.651Z",
    "size": 35050,
    "path": "../public/assets/images/og.webp"
  },
  "/assets/images/tg.webp": {
    "type": "image/webp",
    "etag": "\"7e72-CSJDEBF125KMy71EfkSnmIELZhA\"",
    "mtime": "2026-06-14T09:29:09.651Z",
    "size": 32370,
    "path": "../public/assets/images/tg.webp"
  },
  "/_nuxt/builds/meta/72f4a96d-66c3-4ff3-8e35-cbf0a02671bc.json": {
    "type": "application/json",
    "etag": "\"8b-8+/p9SlntlALYNyQLhJ3h9y9Ykg\"",
    "mtime": "2026-06-14T09:29:09.494Z",
    "size": 139,
    "path": "../public/_nuxt/builds/meta/72f4a96d-66c3-4ff3-8e35-cbf0a02671bc.json"
  }
};

const _DRIVE_LETTER_START_RE = /^[A-Za-z]:\//;
function normalizeWindowsPath(input = "") {
  if (!input) {
    return input;
  }
  return input.replace(/\\/g, "/").replace(_DRIVE_LETTER_START_RE, (r) => r.toUpperCase());
}
const _IS_ABSOLUTE_RE = /^[/\\](?![/\\])|^[/\\]{2}(?!\.)|^[A-Za-z]:[/\\]/;
const _DRIVE_LETTER_RE = /^[A-Za-z]:$/;
function cwd() {
  if (typeof process !== "undefined" && typeof process.cwd === "function") {
    return process.cwd().replace(/\\/g, "/");
  }
  return "/";
}
const resolve = function(...arguments_) {
  arguments_ = arguments_.map((argument) => normalizeWindowsPath(argument));
  let resolvedPath = "";
  let resolvedAbsolute = false;
  for (let index = arguments_.length - 1; index >= -1 && !resolvedAbsolute; index--) {
    const path = index >= 0 ? arguments_[index] : cwd();
    if (!path || path.length === 0) {
      continue;
    }
    resolvedPath = `${path}/${resolvedPath}`;
    resolvedAbsolute = isAbsolute(path);
  }
  resolvedPath = normalizeString(resolvedPath, !resolvedAbsolute);
  if (resolvedAbsolute && !isAbsolute(resolvedPath)) {
    return `/${resolvedPath}`;
  }
  return resolvedPath.length > 0 ? resolvedPath : ".";
};
function normalizeString(path, allowAboveRoot) {
  let res = "";
  let lastSegmentLength = 0;
  let lastSlash = -1;
  let dots = 0;
  let char = null;
  for (let index = 0; index <= path.length; ++index) {
    if (index < path.length) {
      char = path[index];
    } else if (char === "/") {
      break;
    } else {
      char = "/";
    }
    if (char === "/") {
      if (lastSlash === index - 1 || dots === 1) ; else if (dots === 2) {
        if (res.length < 2 || lastSegmentLength !== 2 || res[res.length - 1] !== "." || res[res.length - 2] !== ".") {
          if (res.length > 2) {
            const lastSlashIndex = res.lastIndexOf("/");
            if (lastSlashIndex === -1) {
              res = "";
              lastSegmentLength = 0;
            } else {
              res = res.slice(0, lastSlashIndex);
              lastSegmentLength = res.length - 1 - res.lastIndexOf("/");
            }
            lastSlash = index;
            dots = 0;
            continue;
          } else if (res.length > 0) {
            res = "";
            lastSegmentLength = 0;
            lastSlash = index;
            dots = 0;
            continue;
          }
        }
        if (allowAboveRoot) {
          res += res.length > 0 ? "/.." : "..";
          lastSegmentLength = 2;
        }
      } else {
        if (res.length > 0) {
          res += `/${path.slice(lastSlash + 1, index)}`;
        } else {
          res = path.slice(lastSlash + 1, index);
        }
        lastSegmentLength = index - lastSlash - 1;
      }
      lastSlash = index;
      dots = 0;
    } else if (char === "." && dots !== -1) {
      ++dots;
    } else {
      dots = -1;
    }
  }
  return res;
}
const isAbsolute = function(p) {
  return _IS_ABSOLUTE_RE.test(p);
};
const dirname = function(p) {
  const segments = normalizeWindowsPath(p).replace(/\/$/, "").split("/").slice(0, -1);
  if (segments.length === 1 && _DRIVE_LETTER_RE.test(segments[0])) {
    segments[0] += "/";
  }
  return segments.join("/") || (isAbsolute(p) ? "/" : ".");
};

function readAsset (id) {
  const serverDir = dirname(fileURLToPath(globalThis._importMeta_.url));
  return promises.readFile(resolve(serverDir, assets[id].path))
}

const publicAssetBases = {"/_nuxt/builds/meta/":{"maxAge":31536000},"/_nuxt/builds/":{"maxAge":1},"/_nuxt/":{"maxAge":31536000}};

function isPublicAssetURL(id = '') {
  if (assets[id]) {
    return true
  }
  for (const base in publicAssetBases) {
    if (id.startsWith(base)) { return true }
  }
  return false
}

function getAsset (id) {
  return assets[id]
}

const METHODS = /* @__PURE__ */ new Set(["HEAD", "GET"]);
const EncodingMap = { gzip: ".gz", br: ".br" };
const _5Sizk1 = eventHandler((event) => {
  if (event.method && !METHODS.has(event.method)) {
    return;
  }
  let id = decodePath(
    withLeadingSlash(withoutTrailingSlash(parseURL(event.path).pathname))
  );
  let asset;
  const encodingHeader = String(
    getRequestHeader(event, "accept-encoding") || ""
  );
  const encodings = [
    ...encodingHeader.split(",").map((e) => EncodingMap[e.trim()]).filter(Boolean).sort(),
    ""
  ];
  if (encodings.length > 1) {
    appendResponseHeader(event, "Vary", "Accept-Encoding");
  }
  for (const encoding of encodings) {
    for (const _id of [id + encoding, joinURL(id, "index.html" + encoding)]) {
      const _asset = getAsset(_id);
      if (_asset) {
        asset = _asset;
        id = _id;
        break;
      }
    }
  }
  if (!asset) {
    if (isPublicAssetURL(id)) {
      removeResponseHeader(event, "Cache-Control");
      throw createError$1({ statusCode: 404 });
    }
    return;
  }
  const ifNotMatch = getRequestHeader(event, "if-none-match") === asset.etag;
  if (ifNotMatch) {
    setResponseStatus(event, 304, "Not Modified");
    return "";
  }
  const ifModifiedSinceH = getRequestHeader(event, "if-modified-since");
  const mtimeDate = new Date(asset.mtime);
  if (ifModifiedSinceH && asset.mtime && new Date(ifModifiedSinceH) >= mtimeDate) {
    setResponseStatus(event, 304, "Not Modified");
    return "";
  }
  if (asset.type && !getResponseHeader(event, "Content-Type")) {
    setResponseHeader(event, "Content-Type", asset.type);
  }
  if (asset.etag && !getResponseHeader(event, "ETag")) {
    setResponseHeader(event, "ETag", asset.etag);
  }
  if (asset.mtime && !getResponseHeader(event, "Last-Modified")) {
    setResponseHeader(event, "Last-Modified", mtimeDate.toUTCString());
  }
  if (asset.encoding && !getResponseHeader(event, "Content-Encoding")) {
    setResponseHeader(event, "Content-Encoding", asset.encoding);
  }
  if (asset.size > 0 && !getResponseHeader(event, "Content-Length")) {
    setResponseHeader(event, "Content-Length", asset.size);
  }
  return readAsset(id);
});

const _SxA8c9 = defineEventHandler(() => {});

const _lazy_kaP1kM = () => import('../routes/api/pans.get.mjs');
const _lazy_Dyu8za = () => import('../routes/api/resources.get.mjs');
const _lazy_TOD0pv = () => import('../routes/api/stats.get.mjs');
const _lazy_zGRi9J = () => import('../routes/api/system-config.get.mjs');
const _lazy_alMESp = () => import('../routes/api/version.get.mjs');
const _lazy_antXf3 = () => import('../routes/renderer.mjs').then(function (n) { return n.r; });

const handlers = [
  { route: '', handler: _5Sizk1, lazy: false, middleware: true, method: undefined },
  { route: '/api/pans', handler: _lazy_kaP1kM, lazy: true, middleware: false, method: "get" },
  { route: '/api/resources', handler: _lazy_Dyu8za, lazy: true, middleware: false, method: "get" },
  { route: '/api/stats', handler: _lazy_TOD0pv, lazy: true, middleware: false, method: "get" },
  { route: '/api/system-config', handler: _lazy_zGRi9J, lazy: true, middleware: false, method: "get" },
  { route: '/api/version', handler: _lazy_alMESp, lazy: true, middleware: false, method: "get" },
  { route: '/__nuxt_error', handler: _lazy_antXf3, lazy: true, middleware: false, method: undefined },
  { route: '/__nuxt_island/**', handler: _SxA8c9, lazy: false, middleware: false, method: undefined },
  { route: '/**', handler: _lazy_antXf3, lazy: true, middleware: false, method: undefined }
];

function createNitroApp() {
  const config = useRuntimeConfig();
  const hooks = createHooks();
  const captureError = (error, context = {}) => {
    const promise = hooks.callHookParallel("error", error, context).catch((error_) => {
      console.error("Error while capturing another error", error_);
    });
    if (context.event && isEvent(context.event)) {
      const errors = context.event.context.nitro?.errors;
      if (errors) {
        errors.push({ error, context });
      }
      if (context.event.waitUntil) {
        context.event.waitUntil(promise);
      }
    }
  };
  const h3App = createApp({
    debug: destr(false),
    onError: (error, event) => {
      captureError(error, { event, tags: ["request"] });
      return errorHandler(error, event);
    },
    onRequest: async (event) => {
      event.context.nitro = event.context.nitro || { errors: [] };
      const fetchContext = event.node.req?.__unenv__;
      if (fetchContext?._platform) {
        event.context = {
          _platform: fetchContext?._platform,
          // #3335
          ...fetchContext._platform,
          ...event.context
        };
      }
      if (!event.context.waitUntil && fetchContext?.waitUntil) {
        event.context.waitUntil = fetchContext.waitUntil;
      }
      event.fetch = (req, init) => fetchWithEvent(event, req, init, { fetch: localFetch });
      event.$fetch = (req, init) => fetchWithEvent(event, req, init, {
        fetch: $fetch
      });
      event.waitUntil = (promise) => {
        if (!event.context.nitro._waitUntilPromises) {
          event.context.nitro._waitUntilPromises = [];
        }
        event.context.nitro._waitUntilPromises.push(promise);
        if (event.context.waitUntil) {
          event.context.waitUntil(promise);
        }
      };
      event.captureError = (error, context) => {
        captureError(error, { event, ...context });
      };
      await nitroApp$1.hooks.callHook("request", event).catch((error) => {
        captureError(error, { event, tags: ["request"] });
      });
    },
    onBeforeResponse: async (event, response) => {
      await nitroApp$1.hooks.callHook("beforeResponse", event, response).catch((error) => {
        captureError(error, { event, tags: ["request", "response"] });
      });
    },
    onAfterResponse: async (event, response) => {
      await nitroApp$1.hooks.callHook("afterResponse", event, response).catch((error) => {
        captureError(error, { event, tags: ["request", "response"] });
      });
    }
  });
  const router = createRouter({
    preemptive: true
  });
  const nodeHandler = toNodeListener(h3App);
  const localCall = (aRequest) => b(nodeHandler, aRequest);
  const localFetch = (input, init) => {
    if (!input.toString().startsWith("/")) {
      return globalThis.fetch(input, init);
    }
    return C(
      nodeHandler,
      input,
      init
    ).then((response) => normalizeFetchResponse(response));
  };
  const $fetch = createFetch({
    fetch: localFetch,
    Headers: Headers$1,
    defaults: { baseURL: config.app.baseURL }
  });
  globalThis.$fetch = $fetch;
  h3App.use(createRouteRulesHandler({ localFetch }));
  for (const h of handlers) {
    let handler = h.lazy ? lazyEventHandler(h.handler) : h.handler;
    if (h.middleware || !h.route) {
      const middlewareBase = (config.app.baseURL + (h.route || "/")).replace(
        /\/+/g,
        "/"
      );
      h3App.use(middlewareBase, handler);
    } else {
      const routeRules = getRouteRulesForPath(
        h.route.replace(/:\w+|\*\*/g, "_")
      );
      if (routeRules.cache) {
        handler = cachedEventHandler(handler, {
          group: "nitro/routes",
          ...routeRules.cache
        });
      }
      router.use(h.route, handler, h.method);
    }
  }
  h3App.use(config.app.baseURL, router.handler);
  const app = {
    hooks,
    h3App,
    router,
    localCall,
    localFetch,
    captureError
  };
  return app;
}
function runNitroPlugins(nitroApp2) {
  for (const plugin of plugins) {
    try {
      plugin(nitroApp2);
    } catch (error) {
      nitroApp2.captureError(error, { tags: ["plugin"] });
      throw error;
    }
  }
}
const nitroApp$1 = createNitroApp();
function useNitroApp() {
  return nitroApp$1;
}
runNitroPlugins(nitroApp$1);

function defineRenderHandler(render) {
  const runtimeConfig = useRuntimeConfig();
  return eventHandler(async (event) => {
    const nitroApp = useNitroApp();
    const ctx = { event, render, response: void 0 };
    await nitroApp.hooks.callHook("render:before", ctx);
    if (!ctx.response) {
      if (event.path === `${runtimeConfig.app.baseURL}favicon.ico`) {
        setResponseHeader(event, "Content-Type", "image/x-icon");
        return send(
          event,
          "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
        );
      }
      ctx.response = await ctx.render(event);
      if (!ctx.response) {
        const _currentStatus = getResponseStatus(event);
        setResponseStatus(event, _currentStatus === 200 ? 500 : _currentStatus);
        return send(
          event,
          "No response returned from render handler: " + event.path
        );
      }
    }
    await nitroApp.hooks.callHook("render:response", ctx.response, ctx);
    if (ctx.response.headers) {
      setResponseHeaders(event, ctx.response.headers);
    }
    if (ctx.response.statusCode || ctx.response.statusMessage) {
      setResponseStatus(
        event,
        ctx.response.statusCode,
        ctx.response.statusMessage
      );
    }
    return ctx.response.body;
  });
}

function parse(str, options) {
  if (typeof str !== "string") {
    throw new TypeError("argument str must be a string");
  }
  const obj = {};
  const opt = options || {};
  const dec = opt.decode || decode;
  let index = 0;
  while (index < str.length) {
    const eqIdx = str.indexOf("=", index);
    if (eqIdx === -1) {
      break;
    }
    let endIdx = str.indexOf(";", index);
    if (endIdx === -1) {
      endIdx = str.length;
    } else if (endIdx < eqIdx) {
      index = str.lastIndexOf(";", eqIdx - 1) + 1;
      continue;
    }
    const key = str.slice(index, eqIdx).trim();
    if (opt?.filter && !opt?.filter(key)) {
      index = endIdx + 1;
      continue;
    }
    if (void 0 === obj[key]) {
      let val = str.slice(eqIdx + 1, endIdx).trim();
      if (val.codePointAt(0) === 34) {
        val = val.slice(1, -1);
      }
      obj[key] = tryDecode(val, dec);
    }
    index = endIdx + 1;
  }
  return obj;
}
function decode(str) {
  return str.includes("%") ? decodeURIComponent(str) : str;
}
function tryDecode(str, decode2) {
  try {
    return decode2(str);
  } catch {
    return str;
  }
}

const debug = (...args) => {
};
function GracefulShutdown(server, opts) {
  opts = opts || {};
  const options = Object.assign(
    {
      signals: "SIGINT SIGTERM",
      timeout: 3e4,
      development: false,
      forceExit: true,
      onShutdown: (signal) => Promise.resolve(signal),
      preShutdown: (signal) => Promise.resolve(signal)
    },
    opts
  );
  let isShuttingDown = false;
  const connections = {};
  let connectionCounter = 0;
  const secureConnections = {};
  let secureConnectionCounter = 0;
  let failed = false;
  let finalRun = false;
  function onceFactory() {
    let called = false;
    return (emitter, events, callback) => {
      function call() {
        if (!called) {
          called = true;
          return Reflect.apply(callback, this, arguments);
        }
      }
      for (const e of events) {
        emitter.on(e, call);
      }
    };
  }
  const signals = options.signals.split(" ").map((s) => s.trim()).filter((s) => s.length > 0);
  const once = onceFactory();
  once(process, signals, (signal) => {
    debug("received shut down signal", signal);
    shutdown(signal).then(() => {
      if (options.forceExit) {
        process.exit(failed ? 1 : 0);
      }
    }).catch((error) => {
      debug("server shut down error occurred", error);
      process.exit(1);
    });
  });
  function isFunction(functionToCheck) {
    const getType = Object.prototype.toString.call(functionToCheck);
    return /^\[object\s([A-Za-z]+)?Function]$/.test(getType);
  }
  function destroy(socket, force = false) {
    if (socket._isIdle && isShuttingDown || force) {
      socket.destroy();
      if (socket.server instanceof ft.Server) {
        delete connections[socket._connectionId];
      } else {
        delete secureConnections[socket._connectionId];
      }
    }
  }
  function destroyAllConnections(force = false) {
    debug("Destroy Connections : " + (force ? "forced close" : "close"));
    let counter = 0;
    let secureCounter = 0;
    for (const key of Object.keys(connections)) {
      const socket = connections[key];
      const serverResponse = socket._httpMessage;
      if (serverResponse && !force) {
        if (!serverResponse.headersSent) {
          serverResponse.setHeader("connection", "close");
        }
      } else {
        counter++;
        destroy(socket);
      }
    }
    debug("Connections destroyed : " + counter);
    debug("Connection Counter    : " + connectionCounter);
    for (const key of Object.keys(secureConnections)) {
      const socket = secureConnections[key];
      const serverResponse = socket._httpMessage;
      if (serverResponse && !force) {
        if (!serverResponse.headersSent) {
          serverResponse.setHeader("connection", "close");
        }
      } else {
        secureCounter++;
        destroy(socket);
      }
    }
    debug("Secure Connections destroyed : " + secureCounter);
    debug("Secure Connection Counter    : " + secureConnectionCounter);
  }
  server.on("request", (req, res) => {
    req.socket._isIdle = false;
    if (isShuttingDown && !res.headersSent) {
      res.setHeader("connection", "close");
    }
    res.on("finish", () => {
      req.socket._isIdle = true;
      destroy(req.socket);
    });
  });
  server.on("connection", (socket) => {
    if (isShuttingDown) {
      socket.destroy();
    } else {
      const id = connectionCounter++;
      socket._isIdle = true;
      socket._connectionId = id;
      connections[id] = socket;
      socket.once("close", () => {
        delete connections[socket._connectionId];
      });
    }
  });
  server.on("secureConnection", (socket) => {
    if (isShuttingDown) {
      socket.destroy();
    } else {
      const id = secureConnectionCounter++;
      socket._isIdle = true;
      socket._connectionId = id;
      secureConnections[id] = socket;
      socket.once("close", () => {
        delete secureConnections[socket._connectionId];
      });
    }
  });
  process.on("close", () => {
    debug("closed");
  });
  function shutdown(sig) {
    function cleanupHttp() {
      destroyAllConnections();
      debug("Close http server");
      return new Promise((resolve, reject) => {
        server.close((err) => {
          if (err) {
            return reject(err);
          }
          return resolve(true);
        });
      });
    }
    debug("shutdown signal - " + sig);
    if (options.development) {
      debug("DEV-Mode - immediate forceful shutdown");
      return process.exit(0);
    }
    function finalHandler() {
      if (!finalRun) {
        finalRun = true;
        if (options.finally && isFunction(options.finally)) {
          debug("executing finally()");
          options.finally();
        }
      }
      return Promise.resolve();
    }
    function waitForReadyToShutDown(totalNumInterval) {
      debug(`waitForReadyToShutDown... ${totalNumInterval}`);
      if (totalNumInterval === 0) {
        debug(
          `Could not close connections in time (${options.timeout}ms), will forcefully shut down`
        );
        return Promise.resolve(true);
      }
      const allConnectionsClosed = Object.keys(connections).length === 0 && Object.keys(secureConnections).length === 0;
      if (allConnectionsClosed) {
        debug("All connections closed. Continue to shutting down");
        return Promise.resolve(false);
      }
      debug("Schedule the next waitForReadyToShutdown");
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve(waitForReadyToShutDown(totalNumInterval - 1));
        }, 250);
      });
    }
    if (isShuttingDown) {
      return Promise.resolve();
    }
    debug("shutting down");
    return options.preShutdown(sig).then(() => {
      isShuttingDown = true;
      cleanupHttp();
    }).then(() => {
      const pollIterations = options.timeout ? Math.round(options.timeout / 250) : 0;
      return waitForReadyToShutDown(pollIterations);
    }).then((force) => {
      debug("Do onShutdown now");
      if (force) {
        destroyAllConnections(force);
      }
      return options.onShutdown(sig);
    }).then(finalHandler).catch((error) => {
      const errString = typeof error === "string" ? error : JSON.stringify(error);
      debug(errString);
      failed = true;
      throw errString;
    });
  }
  function shutdownManual() {
    return shutdown("manual");
  }
  return shutdownManual;
}

function getGracefulShutdownConfig() {
  return {
    disabled: !!process.env.NITRO_SHUTDOWN_DISABLED,
    signals: (process.env.NITRO_SHUTDOWN_SIGNALS || "SIGTERM SIGINT").split(" ").map((s) => s.trim()),
    timeout: Number.parseInt(process.env.NITRO_SHUTDOWN_TIMEOUT || "", 10) || 3e4,
    forceExit: !process.env.NITRO_SHUTDOWN_NO_FORCE_EXIT
  };
}
function setupGracefulShutdown(listener, nitroApp) {
  const shutdownConfig = getGracefulShutdownConfig();
  if (shutdownConfig.disabled) {
    return;
  }
  GracefulShutdown(listener, {
    signals: shutdownConfig.signals.join(" "),
    timeout: shutdownConfig.timeout,
    forceExit: shutdownConfig.forceExit,
    onShutdown: async () => {
      await new Promise((resolve) => {
        const timeout = setTimeout(() => {
          console.warn("Graceful shutdown timeout, force exiting...");
          resolve();
        }, shutdownConfig.timeout);
        nitroApp.hooks.callHook("close").catch((error) => {
          console.error(error);
        }).finally(() => {
          clearTimeout(timeout);
          resolve();
        });
      });
    }
  });
}

const cert = process.env.NITRO_SSL_CERT;
const key = process.env.NITRO_SSL_KEY;
const nitroApp = useNitroApp();
const server = cert && key ? new Server({ key, cert }, toNodeListener(nitroApp.h3App)) : new Server$1(toNodeListener(nitroApp.h3App));
const port = destr(process.env.NITRO_PORT || process.env.PORT) || 3e3;
const host = process.env.NITRO_HOST || process.env.HOST;
const path = process.env.NITRO_UNIX_SOCKET;
const listener = server.listen(path ? { path } : { port, host }, (err) => {
  if (err) {
    console.error(err);
    process.exit(1);
  }
  const protocol = cert && key ? "https" : "http";
  const addressInfo = listener.address();
  if (typeof addressInfo === "string") {
    console.log(`Listening on unix socket ${addressInfo}`);
    return;
  }
  const baseURL = (useRuntimeConfig().app.baseURL || "").replace(/\/$/, "");
  const url = `${protocol}://${addressInfo.family === "IPv6" ? `[${addressInfo.address}]` : addressInfo.address}:${addressInfo.port}${baseURL}`;
  console.log(`Listening on ${url}`);
});
trapUnhandledNodeErrors();
setupGracefulShutdown(listener, nitroApp);
const nodeServer = {};

export { createRouter$1 as A, defu as B, parseQuery as C, withTrailingSlash as D, withoutTrailingSlash as E, klona as F, parse as G, getRequestHeader as H, isEqual as I, setCookie as J, getCookie as K, deleteCookie as L, nodeServer as M, getResponseStatusText as a, getResponseStatus as b, createError$1 as c, defineEventHandler as d, defineRenderHandler as e, getRouteRules as f, getQuery as g, useNitroApp as h, hasProtocol as i, joinRelativeURL as j, isScriptProtocol as k, joinURL as l, getContext as m, createHooks as n, executeAsync as o, withBase as p, destr as q, getHeader as r, sanitizeStatusCode as s, getRequestURL as t, useRuntimeConfig as u, setResponseStatus as v, withQuery as w, setResponseHeader as x, send as y, toRouteMatcher as z };
//# sourceMappingURL=nitro.mjs.map
