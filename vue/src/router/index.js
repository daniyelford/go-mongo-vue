import { createRouter, createWebHistory } from 'vue-router';
import { sendApi } from '@/plugins/api';
import Login from '@/components/user/Login.vue';
import Register from '@/components/user/Register.vue';
import Dashboard from '@/components/Dashboard.vue';
import Home from '@/components/Home.vue';
import Welcome from '@/components/Welcome.vue';
import NotFind from '@/components/NotFind.vue';

const routes = [
  { path: '/', name: 'خوش آمدید', component: Welcome, meta: { isLogin: false } },
  { path: '/login', name: 'ورود', component: Login, meta: { isLogin: false } },
  { path: '/register', component: Register, name: 'ثبت نام', meta: { isLogin: true } },
  { path: '/home', component: Home, name: 'خانه', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/dashboard', component: Dashboard, name: 'داشبورد', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/:pathMatch(.*)*', name: '404', component: NotFind }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

async function checkToken() {
  let token = localStorage.getItem('jwt');
  if (!token) return { loggedIn: false };
  try {
    let res = await sendApi({
      method: 'GET',
      url: '/auth/validate',
      headers: { Authorization: `Bearer ${token}` }
    });
    if (!res.error) return res;
    if (res.expired) {
      const refreshToken = localStorage.getItem('refresh');
      if (!refreshToken) return { loggedIn: false };
      const refreshRes = await sendApi({
        method: 'POST',
        url: '/token/refresh',
        data: { refreshToken }
      });
      if (!refreshRes.error && refreshRes.accessToken) {
        localStorage.setItem('jwt', refreshRes.accessToken);
        const newRes = await sendApi({
          method: 'GET',
          url: '/auth/validate',
          headers: { Authorization: `Bearer ${refreshRes.accessToken}` }
        });
        return newRes;
      } else {
        localStorage.removeItem('jwt');
        localStorage.removeItem('refresh');
        return { loggedIn: false };
      }
    }
  } catch (err) {
    localStorage.removeItem('jwt');
    return { loggedIn: false };
  }
}

router.beforeEach(async (to, from, next) => {
  if (to.name) document.title = to.name;

  const mustBeLogin = to.meta.isLogin;
  const mustHasInfo = to.meta.hasUserInfo;

  const status = await checkToken();
  const userIsLogin = status.loggedIn || false;
  const userInfo = status.userHasInfo || false;

  if (userIsLogin && mustHasInfo && !userInfo) return next('/register');
  if (userIsLogin && !mustBeLogin) return next('/home');
  if (!userIsLogin && mustBeLogin) return next('/login');

  next();
});

export default router;
