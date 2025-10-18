<template>
    <b-form @submit.prevent="submitEdit" enctype="multipart/form-data">
      <b-form-group label="Media" class="mt-3">
        <div
          class="border rounded-3 p-4 text-center"
          style="cursor: pointer; background: #f8f9fa;"
          @click="mediaFile.click()"
        >
          <input
            style="display:none"
            ref="mediaFile"
            type="file"
            multiple
            accept="image/*,video/*"
            @change="handleFiles"
          />
            <div v-if="!files.length && !post.media?.length">
                <svg aria-label="Icon to represent media such as images or videos" class="x1lliihq x1n2onr6 x5n08af" fill="currentColor" height="77" role="img" viewBox="0 0 97.6 77.3" width="96"><title>Icon to represent media such as images or videos</title><path d="M16.3 24h.3c2.8-.2 4.9-2.6 4.8-5.4-.2-2.8-2.6-4.9-5.4-4.8s-4.9 2.6-4.8 5.4c.1 2.7 2.4 4.8 5.1 4.8zm-2.4-7.2c.5-.6 1.3-1 2.1-1h.2c1.7 0 3.1 1.4 3.1 3.1 0 1.7-1.4 3.1-3.1 3.1-1.7 0-3.1-1.4-3.1-3.1 0-.8.3-1.5.8-2.1z" fill="currentColor"></path><path d="M84.7 18.4 58 16.9l-.2-3c-.3-5.7-5.2-10.1-11-9.8L12.9 6c-5.7.3-10.1 5.3-9.8 11L5 51v.8c.7 5.2 5.1 9.1 10.3 9.1h.6l21.7-1.2v.6c-.3 5.7 4 10.7 9.8 11l34 2h.6c5.5 0 10.1-4.3 10.4-9.8l2-34c.4-5.8-4-10.7-9.7-11.1zM7.2 10.8C8.7 9.1 10.8 8.1 13 8l34-1.9c4.6-.3 8.6 3.3 8.9 7.9l.2 2.8-5.3-.3c-5.7-.3-10.7 4-11 9.8l-.6 9.5-9.5 10.7c-.2.3-.6.4-1 .5-.4 0-.7-.1-1-.4l-7.8-7c-1.4-1.3-3.5-1.1-4.8.3L7 49 5.2 17c-.2-2.3.6-4.5 2-6.2zm8.7 48c-4.3.2-8.1-2.8-8.8-7.1l9.4-10.5c.2-.3.6-.4 1-.5.4 0 .7.1 1 .4l7.8 7c.7.6 1.6.9 2.5.9.9 0 1.7-.5 2.3-1.1l7.8-8.8-1.1 18.6-21.9 1.1zm76.5-29.5-2 34c-.3 4.6-4.3 8.2-8.9 7.9l-34-2c-4.6-.3-8.2-4.3-7.9-8.9l2-34c.3-4.4 3.9-7.9 8.4-7.9h.5l34 2c4.7.3 8.2 4.3 7.9 8.9z" fill="currentColor"></path><path d="M78.2 41.6 61.3 30.5c-2.1-1.4-4.9-.8-6.2 1.3-.4.7-.7 1.4-.7 2.2l-1.2 20.1c-.1 2.5 1.7 4.6 4.2 4.8h.3c.7 0 1.4-.2 2-.5l18-9c2.2-1.1 3.1-3.8 2-6-.4-.7-.9-1.3-1.5-1.8zm-1.4 6-18 9c-.4.2-.8.3-1.3.3-.4 0-.9-.2-1.2-.4-.7-.5-1.2-1.3-1.1-2.2l1.2-20.1c.1-.9.6-1.7 1.4-2.1.8-.4 1.7-.3 2.5.1L77 43.3c1.2.8 1.5 2.3.7 3.4-.2.4-.5.7-.9.9z" fill="currentColor"></path></svg>
            </div>
          <div class="d-flex flex-wrap gap-3 justify-content-center mt-3">
            <div v-for="(file, index) in files" :key="index" class="position-relative">
              <img
                v-if="file.type.startsWith('image/')"
                :src="file.preview"
                class="rounded shadow-sm"
                style="width:100px;height:100px;object-fit:cover;"
              />
              <video
                v-else-if="file.type.startsWith('video/')"
                :src="file.preview"
                class="rounded shadow-sm"
                style="width:100px;height:100px;object-fit:cover;"
                controls />
              <button type="button" class="btn btn-sm btn-danger position-absolute top-0 end-0 translate-middle rounded-circle" 
                style="width:24px;height:24px;font-size:12px;" @click.stop="removeFile(index)">
                ×
              </button>
            </div>
            <div v-for="(m, idx) in post.media" :key="'old'+idx" class="position-relative">
              <img
                v-if="m.type.startsWith('image/')"
                :src="m.url"
                class="rounded shadow-sm"
                style="width:100px;height:100px;object-fit:cover;"
              />
              <video
                v-else-if="m.type.startsWith('video/')"
                :src="m.url"
                class="rounded shadow-sm"
                style="width:100px;height:100px;object-fit:cover;"
                controls
              />
              <button type="button" class="btn btn-sm btn-danger position-absolute top-0 end-0 translate-middle rounded-circle" 
                style="width:24px;height:24px;font-size:12px;" @click.stop="removeOldMedia(idx)">
                ×
              </button>
            </div>
          </div>
        </div>
      </b-form-group>
      <b-form-group label="Title" class="mt-3">
        <b-form-input v-model="post.title" required />
      </b-form-group>
      <b-form-group label="Content" class="mt-3">
        <b-form-textarea v-model="post.content" rows="5" required />
      </b-form-group>
      <div class="mt-4 d-flex justify-content-end">
        <b-button type="submit" variant="success" :disabled="loading">
          {{ loading ? 'Updating...' : 'Update' }}
        </b-button>
      </div>
    </b-form>
</template>
<script setup>
import { ref, watch } from 'vue'
import { sendApi } from '@/plugins/api'
const props = defineProps({
  postData: Object
})
const emit = defineEmits(['edit'])
const post = ref({ title: '', content: '', media: [] })
watch(() => props.postData, val => {
  if (val) post.value = { ...val }
})
const files = ref([])
const loading = ref(false)
const mediaFile = ref(null)
const handleFiles = (e) => {
  const selected = Array.from(e.target.files)
  selected.forEach(f => f.preview = URL.createObjectURL(f))
  files.value = [...files.value, ...selected]
}
const removeFile = (index) => {
  URL.revokeObjectURL(files.value[index].preview)
  files.value.splice(index, 1)
}
const removeOldMedia = (index) => {
  post.value.media.splice(index, 1)
}
const submitEdit = async () => {
  if (!post.value.title || !post.value.content) return
  loading.value = true
  const formData = new FormData()
  formData.append('id', post.value._id)
  formData.append('title', post.value.title)
  formData.append('content', post.value.content)
  files.value.forEach(f => formData.append('media[]', f))
  formData.append('oldMedia', JSON.stringify(post.value.media))
  try {
    const res = await sendApi({
      url: '/posts/edit',
      data: formData,
      method: 'PUT',
      autoCheckToken: true,
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (res.success) {
      emit('edit')
    } else {
      alert(res.error.message || 'Update failed')
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>
