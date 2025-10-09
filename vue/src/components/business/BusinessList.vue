<template>
    <div class="container py-4">
        <div class="d-flex justify-content-between align-items-center mb-4">
            <h4 class="fw-bold">Business List</h4>
            <b-button variant="primary" @click="showCreateModal = true">
                <i class="bi bi-plus-lg"></i> Create Business
            </b-button>
        </div>
        <b-row v-if="businesses.length">
            <b-col
                v-for="b in businesses"
                :key="b.id"
                cols="12"
                md="6"
                lg="4"
                class="mb-3"
            >
                <b-card
                :title="b.name"
                class="h-100 shadow-sm"
                @click="goToBusiness(b.id)"
                style="cursor: pointer;"
                >
                    <Swiper
                    :modules="[Navigation]"
                    navigation
                    class="mb-2"
                    style="--swiper-navigation-color:#000;">
                    <SwiperSlide v-for="(media, i) in b.media" :key="i">
                        <template v-if="media.type === 'image'">
                        <img :src="media.url" class="img-fluid rounded" alt="" />
                        </template>
                        <template v-else-if="media.type === 'video'">
                        <video controls class="w-100 rounded">
                            <source :src="media.url" type="video/mp4" />
                        </video>
                        </template>
                    </SwiperSlide>
                    </Swiper>
                    <p class="mb-1 text-muted">{{ b.category }}</p>
                    <small>{{ b.description }}</small>
                    <template #footer>
                        <b-badge :variant="b.status === 'active' ? 'success' : 'secondary'">
                            {{ b.status === 'active' ? 'فعال' : 'غیرفعال' }}
                        </b-badge>
                    </template>
                </b-card>
            </b-col>
        </b-row>
        <div v-else class="text-center text-muted mt-5">
            you have not any business
        </div>
        <AddBusiness :showCreate="showCreateModal" @businessCreated="businesses.push($event)"/>
    </div>
</template>

<script setup>
    import { ref, onMounted } from 'vue'
    import { useRouter } from 'vue-router'
    import { sendApi } from '@/plugins/api'
    import AddBusiness from './AddBusiness.vue'
    import { Swiper, SwiperSlide } from 'swiper/vue'
    import { Navigation } from 'swiper/modules'
    import 'swiper/css'
    import 'swiper/css/navigation'
    const router = useRouter()
    const businesses = ref([])
    const showCreateModal = ref(false)
    const goToBusiness = (id) => router.push(`/business/${id}`)
    onMounted(async () => {
        try {
            const res = await sendApi({
                url:'/user/business',
                method:'GET',
                autoCheckToken: true,
                headers:{ 'Content-Type': 'application/json' }
            })
            if(res){
                businesses.value=res
            }else{
                console.log(res);
            }
        } catch (error) {
            console.log(error);
        }
    })
</script>
<style scoped>
    .card:hover {
        transform: translateY(-2px);
        transition: 0.2s;
    }
</style>