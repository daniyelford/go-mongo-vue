export function base64UrlToUint8Array(base64UrlString) {
  const base64 = base64UrlString
    .replace(/-/g, '+')
    .replace(/_/g, '/')
  const padLength = (4 - (base64.length % 4)) % 4
  const padded = base64 + '='.repeat(padLength)
  const binary = atob(padded)
  const array = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) {
    array[i] = binary.charCodeAt(i)
  }
  return array
}