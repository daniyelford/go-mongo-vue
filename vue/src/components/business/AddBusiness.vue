<template>
    <BaseModal title="Create Business" size="lg" v-model="showCreateModal">
        <b-form @submit.prevent="createBusiness">
            <b-form-group label="Brand name" label-for="name">
            <b-form-input
                id="name"
                v-model="form.name"
                placeholder="like Amazon"
                required
            />
            </b-form-group>
            <b-form-group label="Field of activity" label-for="category">
                <b-form-input
                    id="category"
                    v-model="form.category"
                    placeholder="like beauty, restaurant ..."
                />
            </b-form-group>
            <b-form-group label="Description" label-for="description">
                <b-form-textarea
                    id="description"
                    v-model="form.description"
                    rows="3"
                    placeholder="about service"
                />
            </b-form-group>
            <b-form-group label="Media (Images & Videos)" label-for="media">
                <input
                id="media"
                type="file"
                multiple
                accept="image/*,video/*"
                class="form-control"
                @change="onFilesSelected" />
                <div class="mt-3 d-flex flex-wrap gap-3">
                    <div
                    v-for="(file, index) in selectedFiles"
                    :key="index"
                    class="border rounded p-2 position-relative"
                    style="width: 120px;">
                        <button
                            class="btn-close position-absolute top-0 end-0"
                            @click.prevent="removeFile(index)">
                        </button>
                        <img
                        v-if="file.type.startsWith('image/')"
                        :src="file.preview"
                        class="img-fluid rounded"
                        style="max-height: 100px; object-fit: cover;"/>
                        <video
                        v-else-if="file.type.startsWith('video/')"
                        :src="file.preview"
                        class="rounded"
                        controls
                        style="width: 100%; height: 100px; object-fit: cover;" />
                    </div>
                </div>
            </b-form-group>
            <div class="text-end mt-3">
                <b-button type="submit" variant="primary">Save</b-button>
            </div>
        </b-form>
    </BaseModal>
</template>
<script setup>
    import { sendApi } from '@/plugins/api';
    import BaseModal from '@/components/tooles/BaseModal.vue';
    import { BForm,BFormGroup,BFormInput,BFormTextarea } from 'bootstrap-vue-3';
    import { ref, watch } from 'vue'
    const previewImages = ref([])
    const businesses = ref([])
    const selectedFiles = ref([])
    const props = defineProps({
        showCreate: { type: Boolean, required: true },
    })
    const form = ref({
        name: '',
        category: '',
        description: '',
        media: []
    })
    const showCreateModal = ref(props.showCreate)
    const emit = defineEmits(['update:showCreate', 'businessCreated'])
    const onFilesSelected = (e) => {
        const files = Array.from(e.target.files)
        selectedFiles.value = files.map((f) => ({
            file: f,
            type: f.type,
            preview: URL.createObjectURL(f),
        }))
    }
    const removeFile = (index) => {
        selectedFiles.value.splice(index, 1)
    }
    const createBusiness = async () => {
        if (!form.value.name) return alert('نام الزامی است')
        const fd = new FormData()
        fd.append('name', form.value.name)
        fd.append('category', form.value.category)
        fd.append('description', form.value.description)
        for (const f of selectedFiles.value) {
            fd.append('media[]', f.file)
        }
        const res = await sendApi({
            url:'/user/business/add', 
            method: 'POST',
            data: fd,
            autoCheckToken: true
        })
        emit('businessCreated', res)        
        form.value = { name: '', category: '', description: '', media: [] }
        previewImages.value = []
        showCreateModal.value = false
    }
    watch(() => props.showCreate, val => (showCreateModal.value = val))
    watch(showCreateModal, val => emit('update:showCreate', val))
</script>