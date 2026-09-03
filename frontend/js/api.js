const API_BASE_URL = '/api';

export const api = {
  getToken() {
    return localStorage.getItem('suuq_token');
  },

  setToken(token) {
    localStorage.setItem('suuq_token', token);
  },

  clearToken() {
    localStorage.removeItem('suuq_token');
  },

  async request(endpoint, options = {}) {
    const token = this.getToken();
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const config = {
      ...options,
      headers,
    };

    if (config.body && typeof config.body === 'object') {
      config.body = JSON.stringify(config.body);
    }

    try {
      const response = await fetch(`${API_BASE_URL}${endpoint}`, config);

      if (response.status === 401) {
        this.clearToken();
        window.location.href = '/login.html';
        return;
      }

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'An error occurred during the request');
      }

      return data;
    } catch (error) {
      console.error('API Request Failed:', error);
      throw error;
    }
  },

  get(endpoint) {
    return this.request(endpoint, { method: 'GET' });
  },

  post(endpoint, body) {
    return this.request(endpoint, { method: 'POST', body });
  },

  put(endpoint, body) {
    return this.request(endpoint, { method: 'PUT', body });
  },

  delete(endpoint) {
    return this.request(endpoint, { method: 'DELETE' });
  }
};