import router from '@/router';
import axios from 'axios';
// const apiGo = axios.create({ baseURL: '/api', timeout: 5000 })
const apiLaravel = axios.create({ baseURL: '/laravel-api/api', timeout: 5000 })
const api = axios.create({
  baseURL: apiLaravel,
  timeout: 5000,
});
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      console.error("API Error Response:", {
        status: error.response.status,
        data: error.response.data,
        url: error.config.url,
        method: error.config.method,
      });
    } else if (error.request) {
      console.error("API No Response:", error.request);
    } else {
      console.error("API Setup Error:", error.message);
    }
    return Promise.reject(error);
  }
);
function getToken() {
  return localStorage.getItem('jwt');
}
function getRefreshToken() {
  return localStorage.getItem('refresh');
}
function setToken(token) {
  localStorage.setItem('jwt', token);
}
function removeTokens() {
  localStorage.removeItem('jwt');
  localStorage.removeItem('refresh');
  router.push({path:'/login'})
}
export async function sendApi({ method = 'get', url = '', data = {}, headers = {}, autoCheckToken = false }) {
  try {
    let token = getToken();
    if (autoCheckToken && token && getRefreshToken()) {
      try {
        const res = await api.get('/auth/validate', {
          headers: { Authorization: `Bearer ${token}` }
        });
  console.log(res);

        if (res.data?.expired) {
          const refreshRes = await api.post('/token/refresh', 
            { refreshToken: getRefreshToken() },
            { headers: { 'Content-Type': 'application/json' }}
          );
          if (refreshRes.data?.accessToken) {
            token = refreshRes.data.accessToken;
            setToken(token);
          } else {
            removeTokens();
            throw new Error('Session expired');
          }
        }
      } catch (err) {
        console.warn('Token validation failed:', err.message);
        removeTokens();
        return { error: true, message: 'Token validation failed' };
      }
    }
    let hdr=autoCheckToken && token ? { ...headers, Authorization: `Bearer ${token}` } : headers
    const response = await api({
      method,
      url,
      data,
      headers: hdr
    });
    return response.data;
  } catch (error) {
    console.error('API Error:', error);
    return { error: true, message: error.response?.data || error.message };
  }
}
export function removeTokensOut() {
  removeTokens()
}
export default api;