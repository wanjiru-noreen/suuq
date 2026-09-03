import { api } from './api.js';

export const auth = {
  checkAuth() {
    const token = api.getToken();
    const publicPages = ['/login.html', '/register.html', '/index.html', '/'];
    const currentPath = window.location.pathname;

    if (!token && !publicPages.includes(currentPath)) {
      window.location.href = '/login.html';
    } else if (token && (currentPath === '/login.html' || currentPath === '/register.html')) {
      window.location.href = '/dashboard.html';
    }
  },

  async login(email, password) {
    try {
      const data = await api.post('/auth/login', { email, password });
      if (data.token) {
        api.setToken(data.token);
        window.location.href = '/dashboard.html';
      }
      return data;
    } catch (error) {
      throw error;
    }
  },

  async register(name, email, password) {
    try {
      const data = await api.post('/auth/register', { name, email, password });
      if (data.token) {
        api.setToken(data.token);
        window.location.href = '/dashboard.html';
      }
      return data;
    } catch (error) {
      throw error;
    }
  },

  logout() {
    api.clearToken();
    window.location.href = '/login.html';
  }
};

document.addEventListener('DOMContentLoaded', () => {
  auth.checkAuth();
});