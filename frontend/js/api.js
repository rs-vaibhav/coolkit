const API_BASE = '/api/v1';

function getToken() {
  return localStorage.getItem('coolkit_token');
}

function setToken(token) {
  localStorage.setItem('coolkit_token', token);
}

function removeToken() {
  localStorage.removeItem('coolkit_token');
}

function getUser() {
  const user = localStorage.getItem('coolkit_user');
  return user ? JSON.parse(user) : null;
}

function setUser(user) {
  localStorage.setItem('coolkit_user', JSON.stringify(user));
}

function removeUser() {
  localStorage.removeItem('coolkit_user');
}

async function api(endpoint, options = {}) {
  const url = API_BASE + endpoint;
  const token = getToken();

  const defaultHeaders = {
    'Content-Type': 'application/json'
  };

  if (token) {
    defaultHeaders['Authorization'] = `Bearer ${token}`;
  }

  const config = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers
    }
  };

  try {
    const response = await fetch(url, config);

    if (response.status === 401) {
      removeToken();
      removeUser();
      window.location.href = '/';
      return;
    }

    // Try to parse JSON, handle empty responses
    let data = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.indexOf("application/json") !== -1) {
      data = await response.json();
    }

    if (!response.ok) {
      const errorMsg = data && data.message ? data.message : `HTTP Error ${response.status}`;
      throw new Error(errorMsg);
    }

    return data;
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
}

window.CoolKitAPI = { api, getToken, setToken, removeToken, getUser, setUser, removeUser };
