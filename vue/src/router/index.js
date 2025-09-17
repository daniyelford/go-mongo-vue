import { createRouter, createWebHistory } from 'vue-router';
import Login from '@/components/user/Login.vue';
import Register from '@/components/user/Register.vue';
import Dashboard from '@/components/Dashboard.vue';
import Home from '@/components/Home.vue';
import Welcome from '@/components/Welcome.vue';
import NotFind from '@/components/NotFind.vue';
const routes = [
    { 
        path: '/',
        name:'خوش آمدید',
        component: Welcome , 
        meta:{isLogin: false}
    },
    { 
        path: '/login',
        name:'ورود',
        component: Login , 
        meta:{isLogin: false}
    },
    { 
        path: '/register', 
        component: Register ,
        name:'ثبت نام',
        meta:{isLogin: true}
    },
    { 
        path: '/home', 
        component: Home ,
        name:'خانه', 
        meta:{isLogin: true}
    },
    { 
        path: '/dashboard', 
        component: Dashboard ,
        name:'داشبورد',
        meta:{isLogin: true}
    },
    {
        path: '/:pathMatch(.*)*', 
        name: '404', 
        component: NotFind 
    }
];
const router = createRouter({
    history: createWebHistory(),
    routes,
});
router.beforeEach(async (to, from, next) => {
    if (to.name) {
        document.title = to.name;
    }
    const mustBeLogin=to.meta.isLogin
    const userIsLogin=false
    if(userIsLogin && !mustBeLogin){
        return next('/home')
    }        
    if(!userIsLogin && mustBeLogin){
        return next('/login')
    }
    // const security = useSecurityStore()
    // const link = document.querySelector("link[rel~='icon']") || document.createElement('link');
    // link.href = icon;

//     if (to.meta.requiresAuth) {
//         // const ok = await security.checkAuth()
//         // if (!ok) return next('/')
//   }
//     if (to.meta.checkHasMobileId) {
//         // const ok = await security.checkHasMobile()
//         // if (!ok) return next('/')
//     }
//     if (to.meta.onlyAuth) {
//         // const ok = await security.checkOnlyAuth()
//         // if (ok) return next('/dashboard')
//     }
    next()
});
export default router;