document.addEventListener('DOMContentLoaded', () => {
  const { getUser } = window.CoolKitAPI;
  
  // Load user info
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

  loadClub(clubId);
  setupTabs();
  setupEventForm(clubId);
});

async function loadClub(id) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${id}`);
    const club = res.data;
    
    document.title = `${club.name} — CoolKit`;
    document.getElementById('club-name').textContent = club.name;
    document.getElementById('club-desc').textContent = club.description;
    
    const roleBadge = document.getElementById('user-role');
    const role = club.owner_id === window.CoolKitAPI.getUser()?.id ? 'owner' : 'member';
    roleBadge.textContent = role.toUpperCase();
    
    // Set badge color based on role
    roleBadge.className = 'badge'; // reset
    if (role === 'owner') roleBadge.classList.add('badge-primary');
    else if (role === 'admin') roleBadge.classList.add('badge-warning');
    else roleBadge.classList.add('badge-secondary');
    
    // We don't have member count in club response, will fetch members and count
    // document.getElementById('member-count').textContent = `${club.memberCount} Members`;
    // document.getElementById('stat-members').textContent = club.memberCount;
    
    if (role === 'owner' || role === 'admin') {
      document.getElementById('btn-create-event').style.display = 'block';
    } else {
      document.getElementById('btn-create-event').style.display = 'none';
    }
    
    loadMembers(id);
    loadEvents(id);
  } catch (err) {
    console.error('Failed to load club:', err);
    alert('Failed to load club details.');
  }
}

async function loadMembers(id) {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${id}/members`);
    const members = res.data;
    
    document.getElementById('member-count').textContent = `${members.length} Members`;
    document.getElementById('stat-members').textContent = members.length;
    renderMembers(members);
  } catch (err) {
    console.error('Failed to load members:', err);
  }
}

function renderMembers(members) {
  const list = document.getElementById('members-list');
  list.innerHTML = '';
  
  members.forEach(member => {
    const card = document.createElement('div');
    card.className = 'card member-card';
    
    const initial = member.user.name.charAt(0).toUpperCase();
    const joined = new Date(member.joined_at).toLocaleDateString();
    
    let badgeClass = 'badge-secondary';
    if (member.role === 'owner') badgeClass = 'badge-primary';
    else if (member.role === 'admin') badgeClass = 'badge-warning';

    card.innerHTML = `
      <div class="avatar">${initial}</div>
      <div class="member-info">
        <div class="member-name">${member.user.name}</div>
        <div class="member-email">${member.user.email}</div>
        <div class="member-joined">Joined ${joined}</div>
      </div>
      <div>
        <span class="badge ${badgeClass}">${member.role}</span>
      </div>
    `;
    list.appendChild(card);
  });
}

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
      const payload = {
        title: title,
        description: desc,
        date: new Date(date).toISOString(),
        location: loc
      };
      
      await window.CoolKitAPI.api(`/clubs/${clubId}/events`, {
        method: 'POST',
        body: JSON.stringify(payload)
      });
      
      document.getElementById('create-event-modal').classList.remove('active');
      form.reset();
      loadEvents(clubId);
      alert('Event created successfully!');
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to create event.';
      errorDiv.style.display = 'block';
      alert('Error creating event: ' + (err.message || 'Unknown error'));
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
      contents.forEach(c => c.style.display = 'none');
      
      // Add active to clicked
      btn.classList.add('active');
      const tabId = `tab-${btn.dataset.tab}`;
      document.getElementById(tabId).style.display = 'block';
    });
  });
}
