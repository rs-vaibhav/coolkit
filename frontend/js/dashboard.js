document.addEventListener('DOMContentLoaded', () => {
  const { getUser, api } = window.CoolKitAPI;
  
  // Load user info
  const user = getUser();
  if (user) {
    document.getElementById('user-name').textContent = user.name;
    document.getElementById('user-avatar').textContent = getInitial(user.name);
  }

  // Load clubs
  loadClubs();

  // Setup Create Club Form
  const createForm = document.getElementById('create-club-form');
  if (createForm) {
    createForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const name = document.getElementById('create-name').value;
      const desc = document.getElementById('create-desc').value;
      const errorDiv = document.getElementById('create-error');
      
      try {
        errorDiv.style.display = 'none';
        await api('/clubs', { method: 'POST', body: JSON.stringify({ name, description: desc }) });
        
        document.getElementById('create-modal').classList.remove('active');
        createForm.reset();
        loadClubs();
      } catch (err) {
        errorDiv.textContent = err.message;
        errorDiv.style.display = 'block';
      }
    });
  }

  // Setup Join Club Form
  const joinForm = document.getElementById('join-club-form');
  if (joinForm) {
    joinForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const clubId = document.getElementById('join-id').value;
      const errorDiv = document.getElementById('join-error');
      
      try {
        errorDiv.style.display = 'none';
        await api(`/clubs/${clubId}/join`, { method: 'POST' });
        
        document.getElementById('join-modal').classList.remove('active');
        joinForm.reset();
        loadClubs();
      } catch (err) {
        errorDiv.textContent = err.message;
        errorDiv.style.display = 'block';
      }
    });
  }
});

async function loadClubs() {
  try {
    const res = await window.CoolKitAPI.api('/clubs');
    let clubs = res.data || [];
    
    if (clubs.length === 0) {
      renderEmptyState();
    } else {
      renderClubs(clubs);
    }
  } catch (err) {
    console.error('Failed to load clubs:', err);
  }
}

function renderClubs(clubs) {
  const grid = document.getElementById('clubs-grid');
  const empty = document.getElementById('empty-state');
  
  grid.style.display = 'grid';
  empty.style.display = 'none';
  grid.innerHTML = '';

  clubs.forEach(club => {
    const card = document.createElement('div');
    card.className = 'card club-card';
    card.onclick = () => window.location.href = `/club?id=${club.id}`;
    
    const initial = getInitial(club.name);
    const role = club.owner_id === window.CoolKitAPI.getUser()?.id ? 'owner' : 'member';
    const badgeClass = getRoleBadgeClass(role);
    const date = formatDate(club.created_at);
    
    card.innerHTML = `
      <div class="club-card-header">
        <div class="club-card-logo">${initial}</div>
        <div class="club-card-info">
          <div class="club-card-name">${club.name}</div>
          <span class="badge ${badgeClass}">${role}</span>
        </div>
      </div>
      <div class="club-card-description">${club.description || 'No description provided.'}</div>
      <div class="club-card-footer">
        <span>Est. ${date}</span>
      </div>
    `;
    grid.appendChild(card);
  });

  // Add "Create Club" card
  const createCard = document.createElement('div');
  createCard.className = 'create-club-card';
  createCard.onclick = () => document.getElementById('create-modal').classList.add('active');
  createCard.innerHTML = `
    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
    <span style="font-weight: 600;">Create New Club</span>
  `;
  grid.appendChild(createCard);
}

function renderEmptyState() {
  document.getElementById('clubs-grid').style.display = 'none';
  document.getElementById('empty-state').style.display = 'flex';
}

function getInitial(name) {
  return name && name.length > 0 ? name.charAt(0).toUpperCase() : '?';
}

function formatDate(dateString) {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
}

function getRoleBadgeClass(role) {
  switch (role?.toLowerCase()) {
    case 'owner': return 'badge-primary';
    case 'admin': return 'badge-warning';
    default: return 'badge-secondary';
  }
}
