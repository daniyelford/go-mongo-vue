<template>
  <BaseModal v-model="show" size="lg" title="Edit Post">
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
            <p class="text-muted small mt-2 mb-0">Click to select media</p>
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
  </BaseModal>
</template>
<script setup>
import { ref, watch, defineEmits, defineProps } from 'vue'
import { sendApi } from '@/plugins/api'
import BaseModal from '@/components/tooles/BaseModal.vue'
const props = defineProps({
  modelValue: Boolean,
  postData: Object
})
const emit = defineEmits(['update:show', 'updated'])
const show = ref(props.modelValue)
watch(() => props.modelValue, v => show.value = v)
watch(show, v => emit('update:show', v))
const post = ref({ ...props.postData })
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
      emit('updated')
      show.value = false
    } else {
      alert(res.error || 'Update failed')
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>
