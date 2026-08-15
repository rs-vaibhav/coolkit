document.addEventListener('DOMContentLoaded', () => {
  const { getUser } = window.CoolKitAPI;
  
  const user = getUser();
  if (user) {
    document.getElementById('user-name').textContent = user.name;
    document.getElementById('user-avatar').textContent = user.name.charAt(0).toUpperCase();
  }

  const urlParams = new URLSearchParams(window.location.search);
  const clubId = urlParams.get('id');

  if (!clubId) {
    window.location.href = '/dashboard';
    return;
  }

  window._clubId = clubId;
  window._currentUserRole = 'member';
  loadClub(clubId);
  setupTabs();
  setupEventForm(clubId);
  setupAnnouncementForm(clubId);
});

async function loadClub(id) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${id}`);
    const club = res.data;
    
    document.title = `${club.name} — CoolKit`;
    document.getElementById('club-name').textContent = club.name;
    document.getElementById('club-desc').textContent = club.description;
    
    const currentUser = window.CoolKitAPI.getUser();
    const role = club.owner_id === currentUser?.id ? 'owner' : 'member';
    
    const roleBadge = document.getElementById('user-role');
    roleBadge.textContent = role.toUpperCase();
    roleBadge.className = 'badge';
    if (role === 'owner') roleBadge.classList.add('badge-primary');
    else if (role === 'admin') roleBadge.classList.add('badge-warning');
    else roleBadge.classList.add('badge-secondary');
    
    loadMembers(id);
    loadEvents(id);
    loadAnnouncements(id);
  } catch (err) {
    console.error('Failed to load club:', err);
  }
}

// ─── Members ──────────────────────────────
async function loadMembers(id) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${id}/members`);
    const members = res.data;
    const currentUser = window.CoolKitAPI.getUser();
    
    // Find current user's role
    for (const m of members) {
      if (m.user.id === currentUser?.id) {
        window._currentUserRole = m.role;
        break;
      }
    }

    const isAdmin = window._currentUserRole === 'owner' || window._currentUserRole === 'admin';
    
    // Show/hide admin-only buttons
    const createEventBtn = document.getElementById('btn-create-event');
    const createAnnouncementBtn = document.getElementById('btn-create-announcement');
    if (createEventBtn) createEventBtn.style.display = isAdmin ? 'block' : 'none';
    if (createAnnouncementBtn) createAnnouncementBtn.style.display = isAdmin ? 'block' : 'none';
    
    // Show Leave Club button (hide for owners)
    const leaveBtn = document.getElementById('btn-leave-club');
    if (leaveBtn) leaveBtn.style.display = window._currentUserRole !== 'owner' ? 'block' : 'none';
    
    document.getElementById('member-count').textContent = `${members.length} Members`;
    document.getElementById('stat-members').textContent = members.length;
    renderMembers(members, isAdmin);
  } catch (err) {
    console.error('Failed to load members:', err);
  }
}

function renderMembers(members, isAdmin) {
  const list = document.getElementById('members-list');
  list.innerHTML = '';
  const currentUser = window.CoolKitAPI.getUser();
  
  members.forEach(member => {
    const card = document.createElement('div');
    card.className = 'card member-card';
    
    const initial = member.user.name.charAt(0).toUpperCase();
    const joined = new Date(member.joined_at).toLocaleDateString();
    
    let badgeClass = 'badge-secondary';
    if (member.role === 'owner') badgeClass = 'badge-primary';
    else if (member.role === 'admin') badgeClass = 'badge-warning';
    else if (member.role === 'coordinator') badgeClass = 'badge-success';

    let adminControls = '';
    if (isAdmin && member.user.id !== currentUser?.id && member.role !== 'owner') {
      adminControls = `
        <div style="display: flex; gap: var(--spacing-2); margin-top: var(--spacing-2);">
          <select class="form-input" style="padding: 4px 8px; font-size: 12px; width: auto;" onchange="changeRole('${member.user.id}', this.value)">
            <option value="" disabled selected>Change Role</option>
            <option value="admin">Admin</option>
            <option value="coordinator">Coordinator</option>
            <option value="member">Member</option>
          </select>
          <button class="btn btn-ghost btn-sm" style="color: var(--accent-rose); font-size: 12px;" onclick="removeMember('${member.user.id}')">Remove</button>
        </div>
      `;
    }

    card.innerHTML = `
      <div class="avatar">${initial}</div>
      <div class="member-info">
        <div class="member-name">${member.user.name}</div>
        <div class="member-email">${member.user.email}</div>
        <div class="member-joined">Joined ${joined}</div>
        ${adminControls}
      </div>
      <div>
        <span class="badge ${badgeClass}">${member.role}</span>
      </div>
    `;
    list.appendChild(card);
  });
}

async function changeRole(userId, newRole) {
  if (!newRole) return;
  try {
    await window.CoolKitAPI.api(`/clubs/${window._clubId}/members/${userId}/role`, {
      method: 'PUT',
      body: JSON.stringify({ role: newRole })
    });
    loadMembers(window._clubId);
  } catch (err) {
    alert('Failed to change role: ' + (err.message || 'Unknown error'));
  }
}

async function removeMember(userId) {
  if (!confirm('Are you sure you want to remove this member?')) return;
  try {
    await window.CoolKitAPI.api(`/clubs/${window._clubId}/members/${userId}`, {
      method: 'DELETE'
    });
    loadMembers(window._clubId);
  } catch (err) {
    alert('Failed to remove member: ' + (err.message || 'Unknown error'));
  }
}

async function leaveClub() {
  if (!confirm('Are you sure you want to leave this club?')) return;
  try {
    await window.CoolKitAPI.api(`/clubs/${window._clubId}/members/me`, {
      method: 'DELETE'
    });
    window.location.href = '/dashboard';
  } catch (err) {
    alert('Failed to leave club: ' + (err.message || 'Unknown error'));
  }
}

// ─── Events ──────────────────────────────
async function loadEvents(id) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${id}/events`);
    const events = res.data || [];
    
    document.getElementById('stat-events').textContent = events.length;
    renderEvents(events);
  } catch (err) {
    console.error('Failed to load events:', err);
  }
}

function renderEvents(events) {
  const list = document.getElementById('events-list');
  list.innerHTML = '';
  
  if (events.length === 0) {
    list.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-4);">No events scheduled yet.</p>';
    return;
  }
  
  events.forEach(event => {
    const card = document.createElement('div');
    card.className = 'card';
    card.style.marginBottom = 'var(--spacing-3)';
    card.style.cursor = 'pointer';
    card.onclick = () => window.location.href = `/event?id=${event.id}`;
    
    const date = new Date(event.date).toLocaleString();
    
    card.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: flex-start;">
        <div>
          <h4 style="margin: 0 0 var(--spacing-2) 0;">${event.title}</h4>
          <p style="margin: 0 0 var(--spacing-3) 0; color: var(--text-secondary);">${event.description || 'No description provided.'}</p>
        </div>
      </div>
      <div style="display: flex; gap: var(--spacing-4); color: var(--text-muted); font-size: var(--text-sm);">
        <span>📅 ${date}</span>
        ${event.location ? `<span>📍 ${event.location}</span>` : ''}
      </div>
    `;
    list.appendChild(card);
  });
}

function setupEventForm(clubId) {
  const form = document.getElementById('create-event-form');
  if (!form) return;
  
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const title = document.getElementById('event-title').value;
    const desc = document.getElementById('event-desc').value;
    const date = document.getElementById('event-date').value;
    const loc = document.getElementById('event-loc').value;
    const errorDiv = document.getElementById('event-error');
    
    try {
      errorDiv.style.display = 'none';
      await window.CoolKitAPI.api(`/clubs/${clubId}/events`, {
        method: 'POST',
        body: JSON.stringify({
          title, description: desc,
          date: new Date(date).toISOString(),
          location: loc
        })
      });
      
      document.getElementById('create-event-modal').classList.remove('active');
      form.reset();
      loadEvents(clubId);
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to create event.';
      errorDiv.style.display = 'block';
    }
  });
}

// ─── Announcements ──────────────────────────────
async function loadAnnouncements(id) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${id}/announcements`);
    const announcements = res.data || [];
    renderAnnouncements(announcements);
  } catch (err) {
    console.error('Failed to load announcements:', err);
  }
}

function renderAnnouncements(announcements) {
  const container = document.getElementById('announcements-feed');
  if (!container) return;
  container.innerHTML = '';
  
  if (announcements.length === 0) {
    container.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-4);">No announcements yet.</p>';
    return;
  }
  
  announcements.forEach(a => {
    const el = document.createElement('div');
    el.className = 'card';
    el.style.marginBottom = 'var(--spacing-3)';
    
    const priorityBadge = a.priority === 'urgent' ? 'badge-danger' : a.priority === 'important' ? 'badge-warning' : 'badge-secondary';
    const date = new Date(a.created_at).toLocaleString();
    const authorName = a.author?.name || 'Unknown';
    const isAdmin = window._currentUserRole === 'owner' || window._currentUserRole === 'admin';
    
    el.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: flex-start;">
        <div>
          <div style="display: flex; align-items: center; gap: var(--spacing-2); margin-bottom: var(--spacing-2);">
            <h4 style="margin: 0;">${a.title}</h4>
            ${a.priority !== 'normal' ? `<span class="badge ${priorityBadge}">${a.priority}</span>` : ''}
          </div>
          <p style="margin: 0 0 var(--spacing-2) 0; color: var(--text-secondary);">${a.content}</p>
          <div style="font-size: var(--text-xs); color: var(--text-muted);">By ${authorName} • ${date}</div>
        </div>
        ${isAdmin ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose);" onclick="deleteAnnouncement('${a.id}')">✕</button>` : ''}
      </div>
    `;
    container.appendChild(el);
  });
}

function setupAnnouncementForm(clubId) {
  const form = document.getElementById('create-announcement-form');
  if (!form) return;
  
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const title = document.getElementById('announcement-title').value;
    const content = document.getElementById('announcement-content').value;
    const priority = document.getElementById('announcement-priority').value;
    
    try {
      await window.CoolKitAPI.api(`/clubs/${clubId}/announcements`, {
        method: 'POST',
        body: JSON.stringify({ title, content, priority })
      });
      document.getElementById('create-announcement-modal').classList.remove('active');
      form.reset();
      loadAnnouncements(clubId);
    } catch (err) {
      alert('Failed to post announcement: ' + (err.message || 'Unknown error'));
    }
  });
}

async function deleteAnnouncement(id) {
  if (!confirm('Delete this announcement?')) return;
  try {
    await window.CoolKitAPI.api(`/announcements/${id}`, { method: 'DELETE' });
    loadAnnouncements(window._clubId);
  } catch (err) {
    alert('Failed to delete announcement.');
  }
}

// ─── Tabs ──────────────────────────────
function setupTabs() {
  const buttons = document.querySelectorAll('.tab-btn');
  const contents = document.querySelectorAll('.tab-content');
  
  buttons.forEach(btn => {
    btn.addEventListener('click', () => {
      buttons.forEach(b => b.classList.remove('active'));
      contents.forEach(c => c.style.display = 'none');
      
      btn.classList.add('active');
      const tabId = `tab-${btn.dataset.tab}`;
      document.getElementById(tabId).style.display = 'block';
    });
  });
}
