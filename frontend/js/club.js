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
  setupResourceForm(clubId);
  setupBookResourceForm(clubId);
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
    
    // Load requests if admin
    if (isAdmin) {
      document.getElementById('tab-btn-requests').style.display = 'inline-block';
      loadJoinRequests(id);
    }
  } catch (err) {
    console.error('Failed to load members:', err);
  }
}

// ─── Join Requests ──────────────────────────────
async function loadJoinRequests(id) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${id}/requests`);
    renderJoinRequests(res.data || []);
  } catch (err) {
    console.error('Failed to load requests:', err);
  }
}

function renderJoinRequests(requests) {
  const list = document.getElementById('requests-list');
  list.innerHTML = '';
  
  if (requests.length === 0) {
    list.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-4);">No pending requests.</p>';
    return;
  }
  
  requests.forEach(req => {
    const card = document.createElement('div');
    card.className = 'card member-card';
    
    const initial = req.user.name.charAt(0).toUpperCase();
    const requestedAt = new Date(req.created_at).toLocaleDateString();
    
    card.innerHTML = `
      <div class="avatar">${initial}</div>
      <div class="member-info">
        <div class="member-name">${req.user.name}</div>
        <div class="member-email">${req.user.email}</div>
        <div class="member-joined">Requested on ${requestedAt}</div>
      </div>
      <div style="display: flex; gap: var(--spacing-2);">
        <button class="btn btn-primary btn-sm" onclick="approveRequest('${req.id}')">Approve</button>
        <button class="btn btn-secondary btn-sm" onclick="rejectRequest('${req.id}')">Reject</button>
      </div>
    `;
    list.appendChild(card);
  });
}

async function approveRequest(requestId) {
  try {
    await window.CoolKitAPI.api(`/clubs/${window._clubId}/requests/${requestId}/approve`, { method: 'POST' });
    loadJoinRequests(window._clubId);
    loadMembers(window._clubId); // reload members to show the new member
  } catch (err) {
    alert('Failed to approve request: ' + (err.message || 'Unknown error'));
  }
}

async function rejectRequest(requestId) {
  try {
    await window.CoolKitAPI.api(`/clubs/${window._clubId}/requests/${requestId}/reject`, { method: 'POST' });
    loadJoinRequests(window._clubId);
  } catch (err) {
    alert('Failed to reject request: ' + (err.message || 'Unknown error'));
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

      if (btn.dataset.tab === 'resources') {
        loadResources(window._clubId);
      } else if (btn.dataset.tab === 'analytics') {
        loadAnalytics(window._clubId);
      }
    });
  });
}

// ─── Resources & Bookings ──────────────────────────────
async function loadResources(clubId) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${clubId}/resources`);
    const resources = res.data || [];
    
    const bookingsRes = await window.CoolKitAPI.api(`/clubs/${clubId}/bookings`);
    const bookings = bookingsRes.data || [];
    
    const isAdmin = window._currentUserRole === 'owner' || window._currentUserRole === 'admin';
    
    const createResBtn = document.getElementById('btn-create-resource');
    if (createResBtn) createResBtn.style.display = isAdmin ? 'block' : 'none';
    
    renderResources(resources, bookings, isAdmin);
    
    if (isAdmin) {
      document.getElementById('admin-bookings-section').style.display = 'block';
      renderAdminBookings(bookings);
    } else {
      document.getElementById('admin-bookings-section').style.display = 'none';
    }
  } catch (err) {
    console.error('Failed to load resources:', err);
  }
}

function renderResources(resources, bookings, isAdmin) {
  const container = document.getElementById('resources-list');
  container.innerHTML = '';
  
  if (resources.length === 0) {
    container.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-4);">No resources/assets registered for this club.</p>';
    return;
  }
  
  resources.forEach(r => {
    const card = document.createElement('div');
    card.className = 'card';
    card.style.marginBottom = 'var(--spacing-4)';
    card.style.padding = 'var(--spacing-4)';
    
    const resBookings = bookings.filter(b => b.resource_id === r.id);
    let bookingsHtml = '';
    
    if (resBookings.length > 0) {
      bookingsHtml = `
        <div style="margin-top: var(--spacing-4); border-top: 1px solid var(--border-color); padding-top: var(--spacing-3);">
          <strong style="font-size: var(--text-sm); color: var(--text-secondary);">Bookings Schedule:</strong>
          <div style="margin-top: var(--spacing-2); display: flex; flex-direction: column; gap: var(--spacing-2);">
            ${resBookings.map(b => {
              const start = new Date(b.start_time).toLocaleString();
              const end = new Date(b.end_time).toLocaleString();
              const statusClass = b.status === 'approved' ? 'badge-success' : b.status === 'rejected' ? 'badge-danger' : 'badge-secondary';
              const user = b.booked_by?.name || 'Someone';
              const isCurrentUserBooker = b.booked_by_id === window.CoolKitAPI.getUser()?.id;
              
              const cancelBtn = (isCurrentUserBooker || isAdmin) ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose); font-size: 11px; padding: 2px 6px;" onclick="cancelBooking('${b.id}')">Cancel</button>` : '';

              return `
                <div style="display: flex; justify-content: space-between; align-items: center; background-color: var(--bg-muted); padding: var(--spacing-2) var(--spacing-3); border-radius: var(--radius-md); font-size: var(--text-sm);">
                  <div>
                    <strong>${user}</strong>: ${start} - ${end}
                    <div style="color: var(--text-muted); font-size: var(--text-xs); margin-top: 2px;">Purpose: ${b.purpose}</div>
                  </div>
                  <div style="display: flex; align-items: center; gap: var(--spacing-2);">
                    <span class="badge ${statusClass}">${b.status}</span>
                    ${cancelBtn}
                  </div>
                </div>
              `;
            }).join('')}
          </div>
        </div>
      `;
    } else {
      bookingsHtml = `
        <div style="margin-top: var(--spacing-4); border-top: 1px solid var(--border-color); padding-top: var(--spacing-3); color: var(--text-muted); font-size: var(--text-sm);">
          No bookings scheduled yet.
        </div>
      `;
    }

    card.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: flex-start;">
        <div>
          <h4 style="margin: 0 0 var(--spacing-1) 0; font-size: var(--text-lg);">${r.name}</h4>
          <p style="margin: 0 0 var(--spacing-2) 0; color: var(--text-secondary);">${r.description || 'No description provided.'}</p>
          <span class="badge badge-secondary">Qty Available: ${r.quantity}</span>
        </div>
        <div style="display: flex; gap: var(--spacing-2);">
          <button class="btn btn-primary btn-sm" onclick="openBookModal('${r.id}', '${r.name}')">Book</button>
          ${isAdmin ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose);" onclick="deleteResource('${r.id}')">Delete</button>` : ''}
        </div>
      </div>
      ${bookingsHtml}
    `;
    container.appendChild(card);
  });
}

function renderAdminBookings(bookings) {
  const tbody = document.getElementById('admin-bookings-list');
  tbody.innerHTML = '';
  
  const pendingBookings = bookings.filter(b => b.status === 'pending');
  
  if (pendingBookings.length === 0) {
    tbody.innerHTML = `<tr><td colspan="6" style="text-align: center; color: var(--text-muted); padding: var(--spacing-4);">No pending booking requests.</td></tr>`;
    return;
  }
  
  pendingBookings.forEach(b => {
    const tr = document.createElement('tr');
    tr.style.borderBottom = '1px solid var(--border-color)';
    tr.style.fontSize = 'var(--text-sm)';
    
    const start = new Date(b.start_time).toLocaleString();
    const end = new Date(b.end_time).toLocaleString();
    const resName = b.resource?.name || 'Unknown';
    const user = b.booked_by?.name || 'Unknown';

    tr.innerHTML = `
      <td style="padding: var(--spacing-3);"><strong>${resName}</strong></td>
      <td style="padding: var(--spacing-3);">${user}</td>
      <td style="padding: var(--spacing-3);">${start} - ${end}</td>
      <td style="padding: var(--spacing-3); color: var(--text-secondary);">${b.purpose}</td>
      <td style="padding: var(--spacing-3);"><span class="badge badge-secondary">${b.status}</span></td>
      <td style="padding: var(--spacing-3); text-align: right;">
        <button class="btn btn-primary btn-sm" style="margin-right: 4px; padding: 4px 8px; font-size: 11px;" onclick="approveBooking('${b.id}')">Approve</button>
        <button class="btn btn-secondary btn-sm" style="padding: 4px 8px; font-size: 11px;" onclick="rejectBooking('${b.id}')">Reject</button>
      </td>
    `;
    tbody.appendChild(tr);
  });
}

function setupResourceForm(clubId) {
  const form = document.getElementById('create-resource-form');
  if (!form) return;
  
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = document.getElementById('resource-name').value;
    const desc = document.getElementById('resource-desc').value;
    const qty = parseInt(document.getElementById('resource-qty').value, 10);
    
    try {
      await window.CoolKitAPI.api(`/clubs/${clubId}/resources`, {
        method: 'POST',
        body: JSON.stringify({ name, description: desc, quantity: qty })
      });
      document.getElementById('create-resource-modal').classList.remove('active');
      form.reset();
      loadResources(clubId);
    } catch (err) {
      alert('Failed to add resource: ' + (err.message || 'Unknown error'));
    }
  });
}

function setupBookResourceForm(clubId) {
  const form = document.getElementById('book-resource-form');
  if (!form) return;
  
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const resId = document.getElementById('book-resource-id').value;
    const start = document.getElementById('book-start-time').value;
    const end = document.getElementById('book-end-time').value;
    const purpose = document.getElementById('book-purpose').value;
    const errorDiv = document.getElementById('booking-error');
    
    try {
      errorDiv.style.display = 'none';
      await window.CoolKitAPI.api(`/resources/${resId}/bookings`, {
        method: 'POST',
        body: JSON.stringify({
          club_id: clubId,
          start_time: new Date(start).toISOString(),
          end_time: new Date(end).toISOString(),
          purpose
        })
      });
      document.getElementById('book-resource-modal').classList.remove('active');
      form.reset();
      loadResources(clubId);
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to submit booking request.';
      errorDiv.style.display = 'block';
    }
  });
}

function openBookModal(resId, resName) {
  document.getElementById('book-resource-id').value = resId;
  document.getElementById('book-resource-name-display').value = resName;
  document.getElementById('booking-error').style.display = 'none';
  
  // Set default start/end times
  const now = new Date();
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset() + 30); // 30 mins from now
  const startStr = now.toISOString().slice(0, 16);
  now.setHours(now.getHours() + 1); // 1 hour duration
  const endStr = now.toISOString().slice(0, 16);
  
  document.getElementById('book-start-time').value = startStr;
  document.getElementById('book-end-time').value = endStr;
  
  document.getElementById('book-resource-modal').classList.add('active');
}

async function deleteResource(resId) {
  if (!confirm('Are you sure you want to delete this resource? All associated bookings will be lost.')) return;
  try {
    await window.CoolKitAPI.api(`/resources/${resId}`, { method: 'DELETE' });
    loadResources(window._clubId);
  } catch (err) {
    alert('Failed to delete resource.');
  }
}

async function approveBooking(bookingId) {
  try {
    await window.CoolKitAPI.api(`/bookings/${bookingId}/approve`, { method: 'POST' });
    loadResources(window._clubId);
  } catch (err) {
    alert('Failed to approve booking: ' + (err.message || ''));
  }
}

async function rejectBooking(bookingId) {
  try {
    await window.CoolKitAPI.api(`/bookings/${bookingId}/reject`, { method: 'POST' });
    loadResources(window._clubId);
  } catch (err) {
    alert('Failed to reject booking: ' + (err.message || ''));
  }
}

async function cancelBooking(bookingId) {
  if (!confirm('Cancel this booking request?')) return;
  try {
    await window.CoolKitAPI.api(`/bookings/${bookingId}`, { method: 'DELETE' });
    loadResources(window._clubId);
  } catch (err) {
    alert('Failed to cancel booking.');
  }
}

// ─── Analytics ──────────────────────────────
async function loadAnalytics(clubId) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${clubId}/analytics`);
    const analytics = res.data;
    
    const finalMemberCount = analytics.member_growth && analytics.member_growth.length > 0 
      ? analytics.member_growth[analytics.member_growth.length - 1].count 
      : 0;
    
    document.getElementById('analytics-total-members').textContent = finalMemberCount;
    document.getElementById('analytics-total-events').textContent = analytics.event_stats.total_events;
    
    const taskRate = analytics.task_stats.total > 0 
      ? Math.round((analytics.task_stats.done / analytics.task_stats.total) * 100) 
      : 0;
    document.getElementById('analytics-task-rate').textContent = `${taskRate}%`;
    
    const balance = analytics.finance_stats.balance || 0;
    document.getElementById('analytics-balance').textContent = `₹${balance.toLocaleString()}`;
    document.getElementById('analytics-balance').style.color = balance >= 0 ? 'var(--accent-green)' : 'var(--accent-rose)';

    renderMemberGrowthChart(analytics.member_growth || []);
    renderTaskChart(analytics.task_stats);
    renderFinanceChart(analytics.finance_stats);
  } catch (err) {
    console.error('Failed to load analytics:', err);
  }
}

function renderMemberGrowthChart(data) {
  const svg = document.getElementById('member-growth-chart');
  svg.innerHTML = '';
  if (data.length === 0) {
    svg.innerHTML = `<text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="var(--text-muted)">No member growth data available</text>`;
    return;
  }
  
  const width = svg.clientWidth || 350;
  const height = svg.clientHeight || 200;
  const padding = 35;
  
  const maxVal = Math.max(...data.map(d => d.count), 5);
  
  const points = data.map((d, index) => {
    const x = padding + (index / Math.max(data.length - 1, 1)) * (width - 2 * padding);
    const y = height - padding - (d.count / maxVal) * (height - 2 * padding);
    return { x, y, month: d.month, count: d.count };
  });
  
  let pathD = '';
  points.forEach((p, idx) => {
    if (idx === 0) pathD += `M ${p.x} ${p.y}`;
    else pathD += ` L ${p.x} ${p.y}`;
  });

  let areaD = pathD + ` L ${points[points.length-1].x} ${height - padding} L ${points[0].x} ${height - padding} Z`;
  
  svg.innerHTML = `
    <defs>
      <linearGradient id="chartGrad" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="var(--accent-blue)" stop-opacity="0.35"/>
        <stop offset="100%" stop-color="var(--accent-blue)" stop-opacity="0.0"/>
      </linearGradient>
      <linearGradient id="lineGrad" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0%" stop-color="var(--accent-blue)"/>
        <stop offset="100%" stop-color="var(--accent-purple)"/>
      </linearGradient>
    </defs>
    
    <path d="${areaD}" fill="url(#chartGrad)"></path>
    <path d="${pathD}" fill="none" stroke="url(#lineGrad)" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"></path>
    
    ${points.map((p) => `
      <circle cx="${p.x}" cy="${p.y}" r="4" fill="var(--bg-card)" stroke="var(--accent-blue)" stroke-width="2"></circle>
      <text x="${p.x}" y="${p.y - 10}" font-size="10" font-weight="600" text-anchor="middle" fill="white">${p.count}</text>
      <text x="${p.x}" y="${height - 10}" font-size="9" text-anchor="middle" fill="var(--text-secondary)">${p.month}</text>
    `).join('')}
  `;
}

function renderTaskChart(stats) {
  const svg = document.getElementById('task-chart');
  svg.innerHTML = '';
  const total = stats.total || 0;
  if (total === 0) {
    svg.innerHTML = `<text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="var(--text-muted)">No tasks created yet</text>`;
    return;
  }
  
  const width = svg.clientWidth || 350;
  const height = svg.clientHeight || 200;
  const cx = width / 2 - 40;
  const cy = height / 2;
  const r = 50;
  
  const data = [
    { label: 'Todo', value: stats.todo, color: 'var(--text-secondary)' },
    { label: 'In Progress', value: stats.in_progress, color: 'var(--accent-amber)' },
    { label: 'Completed', value: stats.done, color: 'var(--accent-green)' }
  ].filter(d => d.value > 0);
  
  let currentAngle = 0;
  let paths = [];
  
  data.forEach(d => {
    const percent = d.value / total;
    const angle = percent * 360;
    
    const radStart = (currentAngle - 90) * Math.PI / 180;
    const radEnd = (currentAngle + angle - 90) * Math.PI / 180;
    
    const x1 = cx + r * Math.cos(radStart);
    const y1 = cy + r * Math.sin(radStart);
    const x2 = cx + r * Math.cos(radEnd);
    const y2 = cy + r * Math.sin(radEnd);
    
    const largeArcFlag = angle > 180 ? 1 : 0;
    const pathD = `M ${x1} ${y1} A ${r} ${r} 0 ${largeArcFlag} 1 ${x2} ${y2}`;
    
    paths.push({
      d: pathD,
      color: d.color,
      label: d.label,
      value: d.value,
      percent: Math.round(percent * 100)
    });
    
    currentAngle += angle;
  });
  
  svg.innerHTML = `
    ${paths.map(p => `
      <path d="${p.d}" fill="none" stroke="${p.color}" stroke-width="14" stroke-linecap="round"></path>
    `).join('')}
    
    ${paths.map((p, idx) => `
      <g transform="translate(${width - 105}, ${40 + idx * 28})">
        <circle cx="0" cy="0" r="5" fill="${p.color}"></circle>
        <text x="12" y="4" font-size="11" fill="white" font-weight="600">${p.label}</text>
        <text x="12" y="16" font-size="10" fill="var(--text-secondary)">${p.value} (${p.percent}%)</text>
      </g>
    `).join('')}
    
    <text x="${cx}" y="${cy - 2}" text-anchor="middle" font-size="16" font-weight="700" fill="white">${total}</text>
    <text x="${cx}" y="${cy + 12}" text-anchor="middle" font-size="9" fill="var(--text-secondary)">TASKS</text>
  `;
}

function renderFinanceChart(stats) {
  const svg = document.getElementById('finance-chart');
  svg.innerHTML = '';
  
  const categories = stats.categories || [];
  if (categories.length === 0) {
    svg.innerHTML = `<text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="var(--text-muted)">No finance data available</text>`;
    return;
  }
  
  const width = svg.clientWidth || 350;
  const height = svg.clientHeight || 200;
  
  categories.sort((a, b) => b.Amount - a.Amount);
  
  const maxVal = Math.max(...categories.map(c => c.Amount), 1);
  const barHeight = 12;
  const gap = 15;
  const paddingLeft = 85;
  const paddingRight = 60;
  const paddingTop = 25;
  
  svg.innerHTML = categories.map((cat, idx) => {
    const y = paddingTop + idx * (barHeight + gap);
    const barWidth = ((width - paddingLeft - paddingRight) * (cat.Amount / maxVal));
    const isIncome = cat.Type === 'income';
    const color = isIncome ? 'var(--accent-green)' : 'var(--accent-rose)';
    const categoryName = cat.Category.length > 11 ? cat.Category.substring(0, 9) + '..' : cat.Category;

    return `
      <text x="${paddingLeft - 10}" y="${y + barHeight - 1}" font-size="11" text-anchor="end" font-weight="500" fill="var(--text-secondary)">${categoryName}</text>
      <rect x="${paddingLeft}" y="${y}" width="${width - paddingLeft - paddingRight}" height="${barHeight}" fill="var(--bg-muted)" rx="6"></rect>
      <rect x="${paddingLeft}" y="${y}" width="${Math.max(barWidth, 6)}" height="${barHeight}" fill="${color}" rx="6"></rect>
      <text x="${paddingLeft + barWidth + 8}" y="${y + barHeight - 1}" font-size="11" font-weight="600" fill="${color}">₹${cat.Amount.toLocaleString()}</text>
    `;
  }).join('');
}
