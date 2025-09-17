<template>
  <div class="landing">

    <!-- Hero -->
    <b-container fluid class="text-center py-5 text-white" style="background: linear-gradient(120deg, #6366f1, #8b5cf6);">
      <b-row class="justify-content-center">
        <b-col md="8">
          <h1 class="display-4 mb-3">Online Booking, Faster Than Ever</h1>
          <p class="lead mb-4">Book your favorite spot anywhere with just a few clicks.</p>
          <b-button variant="light" class="me-2 mb-2" @click="goToRegister">Start Booking</b-button>
          <b-button variant="outline-light" class="mb-2" @click="scrollTo('features')">See Features</b-button>
        </b-col>
      </b-row>
    </b-container>

    <!-- Features -->
    <b-container id="features" class="text-center py-5">
      <h2 class="mb-4">Why Us?</h2>
      <b-row class="g-3 justify-content-center">
        <b-col cols="6" md="3" class="p-3 border rounded bg-light">🔒 Secure Payment</b-col>
        <b-col cols="6" md="3" class="p-3 border rounded bg-light">⚡ Fast Booking</b-col>
        <b-col cols="6" md="3" class="p-3 border rounded bg-light">⭐ User Ratings</b-col>
        <b-col cols="6" md="3" class="p-3 border rounded bg-light">📅 Automatic Reminders</b-col>
      </b-row>
    </b-container>

    <!-- Suggestions -->
    <b-container class="text-center py-5">
      <h2 class="mb-4">Most Popular Bookings</h2>
      <b-row class="g-4 justify-content-center">
        <b-col cols="12" md="4" v-for="(place, i) in places" :key="i">
          <b-card
            :title="place.name"
            :img-src="media[place.img]"
            :img-alt="place.name"
            img-top
            class="h-100"
          >
            <b-card-text>{{ place.desc }}</b-card-text>
            <p class="fw-bold text-primary">{{ place.price }}</p>
            <b-button variant="primary" size="sm" @click="goToRegister">Quick Book</b-button>
          </b-card>
        </b-col>
      </b-row>
    </b-container>

    <!-- How it Works -->
    <b-container class="py-5 bg-light text-center">
      <h2 class="mb-4">How It Works</h2>
      <b-list-group flush class="mx-auto" style="max-width:400px; text-align:left;">
        <b-list-group-item>1️⃣ Choose a place or service</b-list-group-item>
        <b-list-group-item>2️⃣ Fast and secure payment</b-list-group-item>
        <b-list-group-item>3️⃣ Instant confirmation and reminders</b-list-group-item>
      </b-list-group>
    </b-container>

    <!-- Call To Action -->
    <b-container class="text-center py-5" style="background:#eef2ff;">
      <h2 class="mb-2">Sign Up Now</h2>
      <p class="mb-3">With free membership, access booking history, discounts, and full management.</p>
      <b-button variant="primary" @click="goToRegister">Free Sign Up</b-button>
    </b-container>

    <!-- Footer -->
    <b-container fluid class="text-center py-3 bg-dark text-light">
      <small>© 2025 Online Booking Platform. All rights reserved.</small>
    </b-container>

  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'
import { useMedia } from '@/composables/useMedia'

const { media, loaded, loadMedia } = useMedia()
onMounted(() => loadMedia())

const router = useRouter()

function goToRegister() {
  router.push('/login')
}

function scrollTo(id) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
}

const places = [
  { name: 'Kish Beach Hotel', desc: 'Luxury stay near the sea', price: '2,400,000 IRR', img: 'hotel.jfif' },
  { name: 'Cozy Café Tehran', desc: 'Peaceful environment for gatherings', price: '450,000 IRR', img: 'cafe.jfif' },
  { name: 'Tabriz Gym', desc: 'Professional training with full equipment', price: '600,000 IRR', img: 'gym.jpg' },
]
</script>
