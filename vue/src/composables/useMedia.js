import { ref } from 'vue'
export function useMedia() {
  const media = ref({})
  const loaded = ref(false)
  async function loadMedia() {
    const modules = import.meta.glob('../assets/media/*.{jpg,jpeg,png,jfif,gif,mp4,webm,mp3,wav}')
    const entries = Object.entries(modules)
    const result = {}
    await Promise.all(
      entries.map(async ([key, importer]) => {
        const file = await importer()
        const shortKey = key.split('/').pop()
        result[shortKey] = file.default
      })
    )
    media.value = result
    loaded.value = true
  }
  return { media, loaded, loadMedia }
}