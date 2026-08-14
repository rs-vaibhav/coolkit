document.addEventListener('DOMContentLoaded', () => {
  const { getUser } = window.CoolKitAPI;
  
  // Load user info
  const user = getUser();
  if (user) {
    document.getElementById('user-name').textContent = user.name;
    document.getElementById('user-avatar').textContent = user.name.charAt(0).toUpperCase();
  }

  const urlParams = new URLSearchParams(window.location.search);
  const eventId = urlParams.get('id');

  if (!eventId) {
    window.location.href = '/dashboard';
    return;
  }

  loadEvent(eventId);
  setupTabs();
  setupAssignRoleForm(eventId);
});

let currentEvent = null;

async function loadEvent(id) {
  try {
    const res = await window.CoolKitAPI.api(`/events/${id}`);
    const event = res.data;
    currentEvent = event;
    
    document.title = `${event.title} — CoolKit`;
    document.getElementById('event-title').textContent = event.title;
    document.getElementById('event-desc').textContent = event.description;
    document.getElementById('event-full-desc').textContent = event.description || "No description provided.";
    document.getElementById('event-date').textContent = new Date(event.date).toLocaleString();
    document.getElementById('event-location').textContent = event.location || "TBA";
    
    document.getElementById('back-to-club').href = `/club?id=${event.club_id}`;
    
    // We also need to fetch club members to check roles and populate the assign dropdown
    const clubRes = await window.CoolKitAPI.api(`/clubs/${event.club_id}/members`);
    const members = clubRes.data || [];
    
    // Check if current user is owner/admin
    const currentUser = window.CoolKitAPI.getUser();
    let isAuthorized = false;
    
    const memberSelect = document.getElementById('role-user-id');
    memberSelect.innerHTML = '<option value="" disabled selected>Select a member...</option>';
    
    members.forEach(m => {
      if (m.user.id === currentUser?.id && (m.role === 'owner' || m.role === 'admin')) {
        isAuthorized = true;
      }
      const opt = document.createElement('option');
      opt.value = m.user.id;
      opt.textContent = `${m.user.name} (${m.user.email})`;
      memberSelect.appendChild(opt);
    });

    if (isAuthorized) {
      document.getElementById('btn-assign-role').style.display = 'block';
    }
    
    loadRoles(id);
  } catch (err) {
    console.error('Failed to load event:', err);
    alert('Failed to load event details.');
  }
}

async function loadRoles(id) {
  try {
    const res = await window.CoolKitAPI.api(`/events/${id}/roles`);
    const roles = res.data || [];
    renderRoles(roles);
  } catch (err) {
    console.error('Failed to load roles:', err);
  }
}

function renderRoles(roles) {
  const list = document.getElementById('roles-list');
  list.innerHTML = '';
  
  if (roles.length === 0) {
    list.innerHTML = `
      <div class="empty-state">
        <div style="font-size: 48px; margin-bottom: var(--spacing-4);">👥</div>
        <h3>No roles assigned yet</h3>
        <p>Assign responsibilities to club members for this event.</p>
      </div>
    `;
    return;
  }
  
  roles.forEach(role => {
    const el = document.createElement('div');
    el.className = 'member-card';
    el.innerHTML = `
      <div class="member-info">
        <div class="avatar">${role.user.name.charAt(0).toUpperCase()}</div>
        <div>
          <div class="member-name">${role.user.name}</div>
          <div class="member-email">${role.user.email}</div>
        </div>
      </div>
      <div class="member-role">
        <span class="badge badge-primary">${role.role_name}</span>
      </div>
    `;
    list.appendChild(el);
  });
}

function setupAssignRoleForm(eventId) {
  const form = document.getElementById('assign-role-form');
  if (!form) return;
  
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const userId = document.getElementById('role-user-id').value;
    const roleName = document.getElementById('role-name').value;
    const errorDiv = document.getElementById('role-error');
    
    try {
      errorDiv.style.display = 'none';
      
      await window.CoolKitAPI.api(`/events/${eventId}/roles`, {
        method: 'POST',
        body: JSON.stringify({
          user_id: userId,
          role_name: roleName
        })
      });
      
      document.getElementById('assign-role-modal').classList.remove('active');
      form.reset();
      loadRoles(eventId);
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to assign role.';
      errorDiv.style.display = 'block';
    }
  });
}

function setupTabs() {
  const buttons = document.querySelectorAll('.tab-btn');
  const contents = document.querySelectorAll('.tab-content');
  
  buttons.forEach(btn => {
    btn.addEventListener('click', () => {
      // Remove active from all
      buttons.forEach(b => b.classList.remove('active'));
      contents.forEach(c => {
        c.classList.remove('active');
        c.style.display = 'none';
      });
      
      // Add active to clicked
      btn.classList.add('active');
      const target = document.getElementById(btn.dataset.tab);
      target.classList.add('active');
      target.style.display = 'block';
    });
  });
}
