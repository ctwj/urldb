/**
 * TOTP 动态口令纯函数（微信公众号动态口令对接）
 *
 * 与后端 utils.GenerateDynamicCode、《微信公众号动态口令对接说明.md》同算法：
 * HMAC-SHA1（RFC 4226/6238 截断），6 位数字，密钥为共享密钥的 UTF-8 字节（不做 Base32 转换），
 * counter = floor(unix 时间戳 / 时间窗)，输入为 8 字节大端 counter。
 *
 * 自带纯 JS SHA-1/HMAC 实现，不依赖 crypto.subtle，http 环境下后台也能实时预览当期口令。
 * 自检向量（密钥 moss-wxpass-8f4c2a1e9b7d、周期 1800，见对接文档第四节）：
 *   1700000000 → 534579；1700001800 → 079502；1735689600 → 428950
 */

export const TOTP_PERIOD_SECONDS = 1800

/** SHA-1：输入字节数组（0-255），返回 20 字节数组 */
function sha1(msg: number[]): number[] {
  const bitLen = msg.length * 8
  const padded = msg.slice()
  padded.push(0x80)
  while (padded.length % 64 !== 56) padded.push(0)
  padded.push(0, 0, 0, 0) // 64 位长度的高 32 位（消息长度远小于 2^61）
  padded.push((bitLen >>> 24) & 0xff, (bitLen >>> 16) & 0xff, (bitLen >>> 8) & 0xff, bitLen & 0xff)

  let h0 = 0x67452301
  let h1 = 0xefcdab89
  let h2 = 0x98badcfe
  let h3 = 0x10325476
  let h4 = 0xc3d2e1f0

  const w = new Array<number>(80)
  for (let base = 0; base < padded.length; base += 64) {
    for (let i = 0; i < 16; i++) {
      w[i] = ((padded[base + i * 4] << 24) | (padded[base + i * 4 + 1] << 16) | (padded[base + i * 4 + 2] << 8) | padded[base + i * 4 + 3]) >>> 0
    }
    for (let i = 16; i < 80; i++) {
      const n = w[i - 3] ^ w[i - 8] ^ w[i - 14] ^ w[i - 16]
      w[i] = ((n << 1) | (n >>> 31)) >>> 0
    }

    let a = h0
    let b = h1
    let c = h2
    let d = h3
    let e = h4
    for (let i = 0; i < 80; i++) {
      let f: number
      let k: number
      if (i < 20) {
        f = (b & c) | (~b & d)
        k = 0x5a827999
      } else if (i < 40) {
        f = b ^ c ^ d
        k = 0x6ed9eba1
      } else if (i < 60) {
        f = (b & c) | (b & d) | (c & d)
        k = 0x8f1bbcdc
      } else {
        f = b ^ c ^ d
        k = 0xca62c1d6
      }
      const t = (((a << 5) | (a >>> 27)) + f + e + k + w[i]) >>> 0
      e = d
      d = c
      c = ((b << 30) | (b >>> 2)) >>> 0
      b = a
      a = t
    }
    h0 = (h0 + a) >>> 0
    h1 = (h1 + b) >>> 0
    h2 = (h2 + c) >>> 0
    h3 = (h3 + d) >>> 0
    h4 = (h4 + e) >>> 0
  }

  return [
    (h0 >>> 24) & 0xff, (h0 >>> 16) & 0xff, (h0 >>> 8) & 0xff, h0 & 0xff,
    (h1 >>> 24) & 0xff, (h1 >>> 16) & 0xff, (h1 >>> 8) & 0xff, h1 & 0xff,
    (h2 >>> 24) & 0xff, (h2 >>> 16) & 0xff, (h2 >>> 8) & 0xff, h2 & 0xff,
    (h3 >>> 24) & 0xff, (h3 >>> 16) & 0xff, (h3 >>> 8) & 0xff, h3 & 0xff,
    (h4 >>> 24) & 0xff, (h4 >>> 16) & 0xff, (h4 >>> 8) & 0xff, h4 & 0xff,
  ]
}

/** HMAC-SHA1：返回 20 字节数组 */
function hmacSha1(key: number[], msg: number[]): number[] {
  let k = key
  if (k.length > 64) k = sha1(k)
  const block = new Array<number>(64).fill(0)
  for (let i = 0; i < k.length; i++) block[i] = k[i]
  const inner = sha1(block.map((b) => b ^ 0x36).concat(msg))
  return sha1(block.map((b) => b ^ 0x5c).concat(inner))
}

function utf8Bytes(s: string): number[] {
  return Array.from(new TextEncoder().encode(s))
}

/**
 * 计算当期 6 位动态口令。
 * @param secret 共享密钥（原样字符串，内部转 UTF-8 字节）
 * @param unixSeconds unix 时间戳（秒）
 * @param period 时间窗（秒），默认 1800（30 分钟）
 */
export function generateTotpCode(secret: string, unixSeconds: number, period: number = TOTP_PERIOD_SECONDS): string {
  if (period <= 0) period = TOTP_PERIOD_SECONDS
  const counter = Math.floor(unixSeconds / period)
  const hi = Math.floor(counter / 0x100000000)
  const lo = counter >>> 0
  const msg = [
    (hi >>> 24) & 0xff, (hi >>> 16) & 0xff, (hi >>> 8) & 0xff, hi & 0xff,
    (lo >>> 24) & 0xff, (lo >>> 16) & 0xff, (lo >>> 8) & 0xff, lo & 0xff,
  ]
  const sum = hmacSha1(utf8Bytes(secret), msg)
  const off = sum[sum.length - 1] & 0x0f
  const val = (((sum[off] & 0x7f) << 24) | (sum[off + 1] << 16) | (sum[off + 2] << 8) | sum[off + 3]) >>> 0
  return String(val % 1000000).padStart(6, '0')
}

/** 当前时间窗剩余秒数（用于倒计时展示） */
export function totpRemainingSeconds(unixSeconds: number, period: number = TOTP_PERIOD_SECONDS): number {
  if (period <= 0) period = TOTP_PERIOD_SECONDS
  return period - (unixSeconds % period)
}
