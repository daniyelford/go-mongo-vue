<template>
    <div class="container py-4">
    <b-navbar type="dark" variant="primary" class="rounded-3 mb-4">
      <b-navbar-brand href="#">go mongo vue</b-navbar-brand>
      <b-navbar-nav class="ms-auto">
        <RouterLink class="nav-link text-white" :to="{ path: '/dashboard' }">dashboard</RouterLink>
        <RouterLink class="nav-link text-white" :to="{ path: '/setting' }">setting</RouterLink>
        <Logout class="nav-link text-white"/>
        <Logout />
      </b-navbar-nav>
    </b-navbar>
    <b-card class="shadow-sm rounded-4">
      <b-card-title>📄 posts</b-card-title>
      <b-list-group flush>
        <b-list-group-item
          v-for="post in posts"
          :key="post._id"
          class="d-flex justify-content-between align-items-center"
        >
          <div>
            <h5 class="mb-1">{{ post.title }}</h5>
            <small class="text-muted">{{ post.content }}</small>
          </div>
          <div>
            <b-button size="sm" variant="outline-warning" @click="editPost(post)">edit</b-button>
            <b-button size="sm" variant="outline-danger" class="ms-2" @click="deletePost(post._id)">delete</b-button>
          </div>
        </b-list-group-item>
      </b-list-group>

      <div class="text-center mt-3">
        <b-button
          variant="primary"
          @click="loadMore"
          :disabled="loading || finished"
        >
          {{ loading ? 'loading...' : finished ? 'end' : 'more' }}
        </b-button>
      </div>
    </b-card>
    </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import Logout from '@/components/Logout.vue';
import { sendApi } from '@/plugins/api';
const posts = ref([])
const page = ref(1)
const limit = 10
const loading = ref(false)
const finished = ref(false)
const loadPosts = async () => {
  if (loading.value || finished.value) return
  loading.value = true

  try {
    const res = await sendApi({
        url:'/posts/all',
        method:'POST',
        data:{
            page: page.value,
            limit
        },
        autoCheckToken:true
    })
    const newPosts = res.data || []

    if (newPosts.length < limit) finished.value = true
    posts.value.push(...newPosts)
    page.value++
  } catch (e) {
    console.error('خطا در دریافت پست‌ها:', e)
  } finally {
    loading.value = false
  }
}

const loadMore = () => loadPosts()

const deletePost = async (id) => {
  if (!confirm('آیا مطمئنی می‌خواهی این پست را حذف کنی؟')) return

  try {
    await sendApi('/api/posts/delete', { id })
    posts.value = posts.value.filter(p => p._id !== id)
  } catch (e) {
    console.error('خطا در حذف پست:', e)
  }
}

// ویرایش پست
const editPost = (post) => {
  alert(`ویرایش پست "${post.title}" (در مرحله بعد پیاده می‌کنیم)`)
}

onMounted(() => {
  loadPosts()
  setInterval(refreshPosts, 15000)
})

const refreshPosts = async () => {
  try {
    const res = await sendApi('/api/posts', { page: 1, limit: page.value * limit })
    posts.value = res.data || []
  } catch (e) {
    console.error('خطا در رفرش پست‌ها:', e)
  }
}
</script>