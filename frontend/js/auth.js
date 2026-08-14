document.addEventListener('DOMContentLoaded', () => {
  const { getToken, api, setToken, setUser, removeToken, removeUser } = window.CoolKitAPI;
  
  const path = window.location.pathname;
  const hasToken = !!getToken();

  // Route protection
  if ((path.includes('dashboard') || path.includes('club')) && !hasToken) {
    window.location.href = '/';
    return;
  }
  
  if (path === '/' || path === '/index.html') {
    if (hasToken) {
      window.location.href = '/dashboard';
      return;
    }
  }

  // Setup Auth Modal if it exists
  const loginForm = document.getElementById('login-form');
  const registerForm = document.getElementById('register-form');
  const errorDiv = document.getElementById('auth-error');

  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('login-email').value;
      const password = document.getElementById('login-password').value;
      
      try {
        if (errorDiv) errorDiv.style.display = 'none';
        const res = await api('/auth/login', {
          method: 'POST',
          body: JSON.stringify({ email, password })
        });
        
        setToken(res.data.token);
        setUser(res.data.user);
        window.location.href = '/dashboard';
        
      } catch (err) {
        showError(err.message || 'Login failed. Please check your credentials.');
      }
    });
  }

  if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const name = document.getElementById('register-name').value;
      const email = document.getElementById('register-email').value;
      const password = document.getElementById('register-password').value;
      const confirm = document.getElementById('register-confirm').value;

      if (password !== confirm) {
        return showError('Passwords do not match');
      }
      
      if (password.length < 6) {
        return showError('Password must be at least 6 characters');
      }

      try {
        if (errorDiv) errorDiv.style.display = 'none';
        const res = await api('/auth/register', {
          method: 'POST',
          body: JSON.stringify({ name, email, password })
        });
        
        setToken(res.data.token);
        setUser(res.data.user);
        window.location.href = '/dashboard';
      } catch (err) {
        showError(err.message || 'Registration failed. Please try again.');
      }
    });
  }

  function showError(msg) {
    if (errorDiv) {
      errorDiv.textContent = msg;
      errorDiv.style.display = 'block';
    } else {
      alert(msg);
    }
  }
});

function openAuthModal(tab = 'login') {
  const modal = document.getElementById('auth-modal');
  if (modal) {
    modal.classList.add('active');
    switchAuthTab(tab);
  }
}

function closeAuthModal() {
  const modal = document.getElementById('auth-modal');
  if (modal) {
    modal.classList.remove('active');
    // Clear forms
    document.getElementById('login-form')?.reset();
    document.getElementById('register-form')?.reset();
    const errorDiv = document.getElementById('auth-error');
    if (errorDiv) errorDiv.style.display = 'none';
  }
}

function switchAuthTab(tab) {
  const loginForm = document.getElementById('login-form');
  const registerForm = document.getElementById('register-form');
  const tabLogin = document.getElementById('tab-login');
  const tabRegister = document.getElementById('tab-register');

  if (tab === 'login') {
    loginForm.style.display = 'block';
    registerForm.style.display = 'none';
    tabLogin.style.color = 'white';
    tabRegister.style.color = 'var(--text-secondary)';
  } else {
    loginForm.style.display = 'none';
    registerForm.style.display = 'block';
    tabLogin.style.color = 'var(--text-secondary)';
    tabRegister.style.color = 'white';
  }
}

function logout() {
  window.CoolKitAPI.removeToken();
  window.CoolKitAPI.removeUser();
  window.location.href = '/';
}

window.CoolKitAuth = { openAuthModal, closeAuthModal, switchAuthTab, logout };
