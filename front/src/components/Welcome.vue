<template>
  <div class="landing">
    <!-- Hero -->
    <section class="hero">
      <h1>رزرو آنلاین، سریع‌تر از همیشه</h1>
      <p>هر جا که می‌خوای باشی، تنها با چند کلیک جای دلخواهت رو رزرو کن.</p>
      <div class="cta">
        <button class="btn primary" @click="goToRegister">شروع رزرو</button>
        <button class="btn ghost" @click="scrollTo('features')">دیدن امکانات</button>
      </div>
    </section>

    <!-- Features -->
    <section id="features" class="features">
      <h2>چرا ما؟</h2>
      <div class="grid">
        <div class="item">🔒 پرداخت امن</div>
        <div class="item">⚡ رزرو سریع</div>
        <div class="item">⭐ امتیاز کاربران</div>
        <div class="item">📅 یادآوری خودکار</div>
      </div>
    </section>

    <!-- Suggestions -->
    <section class="suggestions">
      <h2>محبوب‌ترین رزروها</h2>
      <div class="cards">
        <div class="card" v-for="(place, i) in places" :key="i">
          <img :src="place.img" :alt="place.name" />
          <h3>{{ place.name }}</h3>
          <p>{{ place.desc }}</p>
          <p class="price">{{ place.price }}</p>
          <button class="btn small" @click="goToRegister">رزرو سریع</button>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section class="how">
      <h2>چطور کار می‌کنه؟</h2>
      <ol>
        <li>انتخاب محل یا سرویس</li>
        <li>پرداخت سریع و امن</li>
        <li>تأیید فوری و دریافت یادآوری</li>
      </ol>
    </section>

    <!-- Call To Action -->
    <section class="cta-section">
      <h2>همین حالا ثبت‌نام کن</h2>
      <p>با عضویت رایگان، تاریخچه رزرو، تخفیف‌ها و مدیریت کامل در اختیارته.</p>
      <button class="btn primary" @click="goToRegister">ثبت‌نام رایگان</button>
    </section>

    <!-- Footer -->
    <footer>
      <p>© 2025 سامانه رزرو آنلاین. همه حقوق محفوظ است.</p>
    </footer>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()

function goToRegister() {
  router.push('/login')
}

function scrollTo(id) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
}

const places = [
  { name: 'هتل ساحلی کیش', desc: 'اقامت لوکس در نزدیکی دریا', price: '۲,۴۰۰,۰۰۰ تومان', img: 'https://source.unsplash.com/400x250/?beach,resort' },
  { name: 'کافه دنج تهران', desc: 'محیط آرام برای دورهمی', price: '۴۵۰,۰۰۰ تومان', img: 'https://source.unsplash.com/400x250/?cafe' },
  { name: 'سالن ورزشی تبریز', desc: 'تمرین حرفه‌ای با تجهیزات کامل', price: '۶۰۰,۰۰۰ تومان', img: 'https://source.unsplash.com/400x250/?gym' },
]
</script>

<style scoped>
.landing {
  font-family: Vazirmatn, sans-serif;
  color: #1f2937;
  line-height: 1.6;
}

/* Hero */
.hero {
  text-align: center;
  padding: 4rem 1rem;
  background: linear-gradient(120deg, #6366f1, #8b5cf6);
  color: #fff;
}
.hero h1 {
  font-size: 2.4rem;
  margin-bottom: 0.8rem;
}
.hero p { margin-bottom: 1.2rem; }
.cta { display: flex; justify-content: center; gap: 1rem; flex-wrap: wrap; }

/* Buttons */
.btn {
  border: none;
  padding: 0.7rem 1.2rem;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s;
}
.btn.primary { background: #fff; color: #4f46e5; }
.btn.ghost { background: transparent; border: 2px solid #fff; color: #fff; }
.btn.small { padding: 0.4rem 0.8rem; }

/* Features */
.features { padding: 3rem 1rem; text-align: center; }
.features h2 { margin-bottom: 1.2rem; }
.grid { display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap; }
.item { background: #f9fafb; padding: 1rem; border-radius: 10px; }

/* Suggestions */
.suggestions { padding: 3rem 1rem; text-align: center; }
.cards { display: flex; gap: 1.5rem; justify-content: center; flex-wrap: wrap; }
.card { width: 280px; background: #fff; border-radius: 12px; box-shadow: 0 6px 20px rgba(0,0,0,0.08); overflow: hidden; }
.card img { width: 100%; height: 180px; object-fit: cover; }
.card h3 { margin: 0.5rem; font-size: 1.2rem; }
.card p { margin: 0.5rem; }
.price { font-weight: bold; color: #4f46e5; }

/* How it works */
.how { background: #f3f4f6; padding: 3rem 1rem; text-align: center; }
.how ol { list-style: decimal; margin: auto; max-width: 400px; text-align: left; }

/* CTA Section */
.cta-section { text-align: center; padding: 3rem 1rem; background: #eef2ff; }
.cta-section h2 { margin-bottom: 0.5rem; }

/* Footer */
footer { text-align: center; padding: 1rem; background: #111827; color: #9ca3af; }
</style>
