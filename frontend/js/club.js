// ─── Global State ─────────────────────────────
let _clubId = null;
let _currentUserRole = 'member';
let _orgData = { members: [], hierarchy_levels: [], domains: [] };
let _navState = { view: 'root', domainId: null, levelPosition: null };
let _assignTarget = null;

document.addEventListener('DOMContentLoaded', () => {
  const { getUser } = window.CoolKitAPI;
  
  const user = getUser();
  if (user) {
    document.getElementById('user-name').textContent = user.name;
    document.getElementById('user-avatar').textContent = user.name.charAt(0).toUpperCase();
  }

  const urlParams = new URLSearchParams(window.location.search);
  _clubId = urlParams.get('id');

  if (!_clubId) {
    window.location.href = '/dashboard';
    return;
  }

  loadClub(_clubId);
  setupTabs();
  setupEventForm(_clubId);
  setupAnnouncementForm(_clubId);
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

// ─── Organization Tree ──────────────────────────────
async function loadMembers(id) {
  try {
    // Load organization data (members + hierarchy + domains)
    const orgRes = await window.CoolKitAPI.api(`/clubs/${id}/organization`);
    _orgData = orgRes.data || { members: [], hierarchy_levels: [], domains: [] };
    // Ensure arrays
    _orgData.members = _orgData.members || [];
    _orgData.hierarchy_levels = _orgData.hierarchy_levels || [];
    _orgData.domains = _orgData.domains || [];
    
    const currentUser = window.CoolKitAPI.getUser();
    
    // Find current user's role
    for (const m of _orgData.members) {
      if (m.user && m.user.id === currentUser?.id) {
        _currentUserRole = m.role;
        break;
      }
    }

    const isAdmin = _currentUserRole === 'owner' || _currentUserRole === 'admin';
    
    // Show/hide admin-only buttons
    const createEventBtn = document.getElementById('btn-create-event');
    const createAnnouncementBtn = document.getElementById('btn-create-announcement');
    const orgToolbar = document.getElementById('org-admin-toolbar');
    if (createEventBtn) createEventBtn.style.display = isAdmin ? 'block' : 'none';
    if (createAnnouncementBtn) createAnnouncementBtn.style.display = isAdmin ? 'block' : 'none';
    if (orgToolbar) orgToolbar.style.display = isAdmin ? 'flex' : 'none';
    
    // Show Leave Club button (hide for owners)
    const leaveBtn = document.getElementById('btn-leave-club');
    if (leaveBtn) leaveBtn.style.display = _currentUserRole !== 'owner' ? 'block' : 'none';
    
    document.getElementById('member-count').textContent = `${_orgData.members.length} Members`;
    document.getElementById('stat-members').textContent = _orgData.members.length;
    
    // Load requests if admin
    if (isAdmin) {
      document.getElementById('tab-btn-requests').style.display = 'inline-block';
      loadJoinRequests(id);
    }
    
    // Render the organization tree
    renderOrganizationTree();
    
    // Also load domain list in domain modal
    renderDomainListEditor();
  } catch (err) {
    console.error('Failed to load organization:', err);
    // Fallback: try loading members the old way
    try {
      const res = await window.CoolKitAPI.api(`/clubs/${id}/members`);
      _orgData.members = res.data || [];
      renderOrganizationTree();
    } catch (e2) {
      console.error('Fallback also failed:', e2);
    }
  }
}

function renderOrganizationTree() {
  const container = document.getElementById('org-tree');
  container.innerHTML = '';
  
  const { view, domainId, levelPosition } = _navState;
  
  if (view === 'root') {
    renderRootView(container);
  } else if (view === 'domain') {
    renderDomainView(container, domainId);
  } else if (view === 'level') {
    renderLevelView(container, domainId, levelPosition);
  }
  
  updateBreadcrumb();
}

// ─── Root View: Leadership + Domains + Unassigned ─────
function renderRootView(container) {
  const isAdmin = _currentUserRole === 'owner' || _currentUserRole === 'admin';
  
  // 1. Leadership section (owners + admins)
  const leaders = _orgData.members.filter(m => m.role === 'owner' || m.role === 'admin');
  if (leaders.length > 0) {
    const section = createSection('👑', 'Leadership', leaders.length);
    const grid = document.createElement('div');
    grid.className = 'org-members-grid';
    leaders.forEach(m => grid.appendChild(createMemberCard(m, isAdmin)));
    section.appendChild(grid);
    container.appendChild(section);
  }
  
  // 2. Domains section
  if (_orgData.domains.length > 0) {
    const section = createSection('📁', 'Domains', _orgData.domains.length);
    const grid = document.createElement('div');
    grid.className = 'org-domains-grid';
    
    _orgData.domains.forEach(domain => {
      const domainMembers = _orgData.members.filter(m => 
        m.domain_id === domain.id && m.role !== 'owner' && m.role !== 'admin'
      );
      grid.appendChild(createDomainCard(domain, domainMembers.length));
    });
    
    section.appendChild(grid);
    container.appendChild(section);
  }
  
  // 3. Unassigned members (no domain, not leadership)
  const unassigned = _orgData.members.filter(m => 
    !m.domain_id && m.role !== 'owner' && m.role !== 'admin'
  );
  if (unassigned.length > 0) {
    const section = createSection('👤', `Unassigned`, unassigned.length);
    const grid = document.createElement('div');
    grid.className = 'org-members-grid';
    unassigned.forEach(m => grid.appendChild(createMemberCard(m, isAdmin)));
    section.appendChild(grid);
    container.appendChild(section);
  }
  
  // Empty state
  if (_orgData.members.length === 0) {
    container.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-8);">No members yet.</p>';
  }
}

// ─── Domain View: members at top level + level cards ─────
function renderDomainView(container, domainId) {
  const isAdmin = _currentUserRole === 'owner' || _currentUserRole === 'admin';
  const domain = _orgData.domains.find(d => d.id === domainId);
  if (!domain) { navigateTo('root'); return; }
  
  const domainMembers = _orgData.members.filter(m => m.domain_id === domainId && m.role !== 'owner' && m.role !== 'admin');
  const levels = _orgData.hierarchy_levels;
  
  if (levels.length === 0) {
    // No hierarchy defined — show all domain members flat
    const section = createSection('👥', domain.name, domainMembers.length);
    const grid = document.createElement('div');
    grid.className = 'org-members-grid';
    domainMembers.forEach(m => grid.appendChild(createMemberCard(m, isAdmin)));
    section.appendChild(grid);
    container.appendChild(section);
    return;
  }
  
  // Show members at the highest hierarchy level (position 1)
  const topLevel = levels[0]; // sorted by position ASC
  const topMembers = domainMembers.filter(m => m.hierarchy_level_id === topLevel.id);
  
  if (topMembers.length > 0) {
    const section = createSection('🏆', topLevel.name, topMembers.length);
    const grid = document.createElement('div');
    grid.className = 'org-members-grid';
    topMembers.forEach(m => grid.appendChild(createMemberCard(m, isAdmin)));
    section.appendChild(grid);
    container.appendChild(section);
  }
  
  // Show clickable cards for remaining levels
  const remainingLevels = levels.slice(1);
  if (remainingLevels.length > 0) {
    const levelsGrid = document.createElement('div');
    levelsGrid.className = 'org-levels-grid';
    
    remainingLevels.forEach(level => {
      const levelMembers = domainMembers.filter(m => m.hierarchy_level_id === level.id);
      levelsGrid.appendChild(createLevelCard(level, levelMembers.length, domainId));
    });
    
    container.appendChild(levelsGrid);
  }
  
  // Members with no hierarchy level in this domain
  const noLevel = domainMembers.filter(m => !m.hierarchy_level_id);
  if (noLevel.length > 0) {
    const section = createSection('👤', 'No Level Assigned', noLevel.length);
    section.style.marginTop = 'var(--spacing-6)';
    const grid = document.createElement('div');
    grid.className = 'org-members-grid';
    noLevel.forEach(m => grid.appendChild(createMemberCard(m, isAdmin)));
    section.appendChild(grid);
    container.appendChild(section);
  }
}

// ─── Level View: members at specific level + deeper levels ─────
function renderLevelView(container, domainId, levelPosition) {
  const isAdmin = _currentUserRole === 'owner' || _currentUserRole === 'admin';
  const levels = _orgData.hierarchy_levels;
  const currentLevel = levels.find(l => l.position === levelPosition);
  if (!currentLevel) { navigateTo('domain', domainId); return; }
  
  const domainMembers = _orgData.members.filter(m => m.domain_id === domainId && m.role !== 'owner' && m.role !== 'admin');
  const levelMembers = domainMembers.filter(m => m.hierarchy_level_id === currentLevel.id);
  
  // Show members at this level
  if (levelMembers.length > 0) {
    const section = createSection('🏅', currentLevel.name, levelMembers.length);
    const grid = document.createElement('div');
    grid.className = 'org-members-grid';
    levelMembers.forEach(m => grid.appendChild(createMemberCard(m, isAdmin)));
    section.appendChild(grid);
    container.appendChild(section);
  } else {
    const empty = document.createElement('p');
    empty.style.cssText = 'color: var(--text-muted); text-align: center; padding: var(--spacing-4);';
    empty.textContent = `No members at ${currentLevel.name} level yet.`;
    container.appendChild(empty);
  }
  
  // Show deeper levels as clickable cards
  const deeperLevels = levels.filter(l => l.position > levelPosition);
  if (deeperLevels.length > 0) {
    const levelsSection = document.createElement('div');
    levelsSection.style.marginTop = 'var(--spacing-6)';
    const levelsGrid = document.createElement('div');
    levelsGrid.className = 'org-levels-grid';
    
    deeperLevels.forEach(level => {
      const lm = domainMembers.filter(m => m.hierarchy_level_id === level.id);
      levelsGrid.appendChild(createLevelCard(level, lm.length, domainId));
    });
    
    levelsSection.appendChild(levelsGrid);
    container.appendChild(levelsSection);
  }
}

// ─── Navigation ─────────────────────────────
function navigateTo(view, domainId = null, levelPosition = null) {
  _navState = { view, domainId, levelPosition };
  renderOrganizationTree();
}

function updateBreadcrumb() {
  const bc = document.getElementById('org-breadcrumb');
  const { view, domainId, levelPosition } = _navState;
  
  if (view === 'root') {
    bc.style.display = 'none';
    return;
  }
  
  bc.style.display = 'flex';
  bc.innerHTML = '';
  
  // Root link
  const rootLink = document.createElement('span');
  rootLink.className = 'org-breadcrumb-item';
  rootLink.textContent = 'All Members';
  rootLink.onclick = () => navigateTo('root');
  bc.appendChild(rootLink);
  
  if (view === 'domain' || view === 'level') {
    const domain = _orgData.domains.find(d => d.id === domainId);
    bc.appendChild(createBreadcrumbSeparator());
    
    if (view === 'domain') {
      const current = document.createElement('span');
      current.className = 'org-breadcrumb-current';
      current.textContent = domain?.name || 'Domain';
      bc.appendChild(current);
    } else {
      const domainLink = document.createElement('span');
      domainLink.className = 'org-breadcrumb-item';
      domainLink.textContent = domain?.name || 'Domain';
      domainLink.onclick = () => navigateTo('domain', domainId);
      bc.appendChild(domainLink);
      
      bc.appendChild(createBreadcrumbSeparator());
      const level = _orgData.hierarchy_levels.find(l => l.position === levelPosition);
      const current = document.createElement('span');
      current.className = 'org-breadcrumb-current';
      current.textContent = level?.name || 'Level';
      bc.appendChild(current);
    }
  }
}

function createBreadcrumbSeparator() {
  const sep = document.createElement('span');
  sep.className = 'org-breadcrumb-separator';
  sep.textContent = '›';
  return sep;
}

// ─── UI Component Builders ─────────────────────────────
function createSection(icon, title, count) {
  const section = document.createElement('div');
  section.className = 'org-section';
  
  const header = document.createElement('div');
  header.className = 'org-section-header';
  header.innerHTML = `
    <span class="org-section-title">${title}</span>
    <span class="org-section-count">(${count})</span>
  `;
  section.appendChild(header);
  return section;
}

function createMemberCard(member, isAdmin) {
  const card = document.createElement('div');
  card.className = 'org-member-card';
  
  const user = member.user || {};
  const initial = user.name ? user.name.charAt(0).toUpperCase() : '?';
  const avatarClass = member.role || 'member';
  
  let levelBadge = '';
  if (member.hierarchy_level) {
    levelBadge = `<span class="badge badge-secondary">${member.hierarchy_level.name}</span>`;
  }
  
  let domainBadge = '';
  if (member.domain && _navState.view === 'root') {
    domainBadge = `<span class="badge badge-secondary">${member.domain.name}</span>`;
  }
  
  let roleBadge = '';
  if (member.role === 'owner') roleBadge = '<span class="badge badge-primary">Owner</span>';
  else if (member.role === 'admin') roleBadge = '<span class="badge badge-warning">Admin</span>';
  
  let adminBtn = '';
  if (isAdmin && member.role !== 'owner') {
    adminBtn = `<button class="btn btn-ghost btn-sm" style="font-size: 11px; padding: 2px 8px; margin-top: 4px;" onclick="event.stopPropagation(); openAssignModal('${member.user.id}', '${user.name}', '${member.domain_id || ''}', '${member.hierarchy_level_id || ''}')">Assign ✏️</button>`;
  }
  
  card.innerHTML = `
    <div class="org-member-avatar ${avatarClass}">${initial}</div>
    <div class="org-member-details">
      <div class="org-member-name">${user.name || 'Unknown'}</div>
      <div class="org-member-email">${user.email || ''}</div>
      <div class="org-member-meta">
        ${roleBadge}${levelBadge}${domainBadge}
        ${adminBtn}
      </div>
    </div>
  `;
  return card;
}

function createDomainCard(domain, memberCount) {
  const card = document.createElement('div');
  card.className = 'org-domain-card';
  card.onclick = () => navigateTo('domain', domain.id);
  
  card.innerHTML = `
    <div class="org-domain-name">${domain.name}</div>
    <div class="org-domain-desc">${domain.description || 'No description'}</div>
    <div class="org-domain-footer">
      <span class="org-domain-count">${memberCount} member${memberCount !== 1 ? 's' : ''}</span>
      <span class="org-domain-arrow">View →</span>
    </div>
  `;
  return card;
}

function createLevelCard(level, memberCount, domainId) {
  const card = document.createElement('div');
  card.className = 'org-level-card';
  card.onclick = () => navigateTo('level', domainId, level.position);
  
  card.innerHTML = `
    <span class="org-level-name">${level.name}</span>
    <span class="org-level-count">${memberCount} member${memberCount !== 1 ? 's' : ''} →</span>
  `;
  return card;
}

// ─── Hierarchy Editor ─────────────────────────────
async function loadHierarchyEditor() {
  try {
    const res = await window.CoolKitAPI.api(`/clubs/${_clubId}/hierarchy`);
    const levels = res.data || [];
    
    const editor = document.getElementById('hierarchy-editor');
    editor.innerHTML = '';
    
    if (levels.length === 0) {
      // Start with a default template
      addHierarchyLevel('Head');
      addHierarchyLevel('Lead');
      addHierarchyLevel('Member');
    } else {
      levels.forEach(l => addHierarchyLevel(l.name));
    }
  } catch (err) {
    console.error('Failed to load hierarchy:', err);
    // Start with empty
    const editor = document.getElementById('hierarchy-editor');
    editor.innerHTML = '';
    addHierarchyLevel('Head');
    addHierarchyLevel('Lead');
    addHierarchyLevel('Member');
  }
}

function addHierarchyLevel(name = '') {
  const editor = document.getElementById('hierarchy-editor');
  const position = editor.children.length + 1;
  
  const row = document.createElement('div');
  row.className = 'hierarchy-level-row';
  row.innerHTML = `
    <span class="hierarchy-level-position">${position}</span>
    <input type="text" class="hierarchy-level-input" value="${name}" placeholder="Level name (e.g., Lead)">
    <button class="hierarchy-level-remove" onclick="this.parentElement.remove(); renumberHierarchy()">&times;</button>
  `;
  editor.appendChild(row);
}

function renumberHierarchy() {
  const editor = document.getElementById('hierarchy-editor');
  const rows = editor.querySelectorAll('.hierarchy-level-row');
  rows.forEach((row, i) => {
    row.querySelector('.hierarchy-level-position').textContent = i + 1;
  });
}

async function saveHierarchy() {
  const editor = document.getElementById('hierarchy-editor');
  const rows = editor.querySelectorAll('.hierarchy-level-row');
  const levels = [];
  
  rows.forEach((row, i) => {
    const name = row.querySelector('.hierarchy-level-input').value.trim();
    if (name) {
      levels.push({ name, position: i + 1 });
    }
  });
  
  if (levels.length === 0) {
    alert('Please add at least one hierarchy level.');
    return;
  }
  
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/hierarchy`, {
      method: 'POST',
      body: JSON.stringify({ levels })
    });
    
    document.getElementById('hierarchy-modal').classList.remove('active');
    loadMembers(_clubId); // Reload to get updated hierarchy
  } catch (err) {
    alert('Failed to save hierarchy: ' + (err.message || 'Unknown error'));
  }
}

// ─── Domain Management ─────────────────────────────
function renderDomainListEditor() {
  const container = document.getElementById('domain-list-editor');
  if (!container) return;
  container.innerHTML = '';
  
  if (_orgData.domains.length === 0) {
    container.innerHTML = '<p style="color: var(--text-muted); font-size: var(--text-sm);">No domains created yet.</p>';
    return;
  }
  
  _orgData.domains.forEach(domain => {
    const row = document.createElement('div');
    row.style.cssText = 'display: flex; justify-content: space-between; align-items: center; padding: var(--spacing-3); background: var(--bg-tertiary); border-radius: var(--radius-md); margin-bottom: var(--spacing-2);';
    
    const memberCount = _orgData.members.filter(m => m.domain_id === domain.id).length;
    
    row.innerHTML = `
      <div>
        <strong>${domain.name}</strong>
        <span style="color: var(--text-muted); font-size: var(--text-sm); margin-left: var(--spacing-2);">(${memberCount} members)</span>
        ${domain.description ? `<div style="color: var(--text-secondary); font-size: var(--text-xs);">${domain.description}</div>` : ''}
      </div>
      <button class="btn btn-ghost btn-sm" style="color: var(--accent-rose); font-size: 12px;" onclick="deleteDomain('${domain.id}')">Delete</button>
    `;
    container.appendChild(row);
  });
}

async function createDomain() {
  const name = document.getElementById('new-domain-name').value.trim();
  const desc = document.getElementById('new-domain-desc').value.trim();
  
  if (!name) {
    alert('Please enter a domain name.');
    return;
  }
  
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/domains`, {
      method: 'POST',
      body: JSON.stringify({ name, description: desc })
    });
    
    document.getElementById('new-domain-name').value = '';
    document.getElementById('new-domain-desc').value = '';
    
    // Reload data
    await loadMembers(_clubId);
    renderDomainListEditor();
  } catch (err) {
    alert('Failed to create domain: ' + (err.message || 'Unknown error'));
  }
}

async function deleteDomain(domainId) {
  if (!confirm('Are you sure? Members in this domain will become unassigned.')) return;
  
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/domains/${domainId}`, {
      method: 'DELETE'
    });
    
    // If we're viewing this domain, go back to root
    if (_navState.domainId === domainId) {
      navigateTo('root');
    }
    
    await loadMembers(_clubId);
    renderDomainListEditor();
  } catch (err) {
    alert('Failed to delete domain: ' + (err.message || 'Unknown error'));
  }
}

// ─── Assign Member Organization ─────────────────────────────
function openAssignModal(userId, userName, currentDomainId, currentLevelId) {
  _assignTarget = userId;
  document.getElementById('assign-member-name').textContent = userName;
  
  // Populate domain dropdown
  const domainSelect = document.getElementById('assign-domain');
  domainSelect.innerHTML = '<option value="">No domain (unassigned)</option>';
  _orgData.domains.forEach(d => {
    const opt = document.createElement('option');
    opt.value = d.id;
    opt.textContent = d.name;
    if (d.id === currentDomainId) opt.selected = true;
    domainSelect.appendChild(opt);
  });
  
  // Populate level dropdown
  const levelSelect = document.getElementById('assign-level');
  levelSelect.innerHTML = '<option value="">No level</option>';
  _orgData.hierarchy_levels.forEach(l => {
    const opt = document.createElement('option');
    opt.value = l.id;
    opt.textContent = `${l.position}. ${l.name}`;
    if (l.id === currentLevelId) opt.selected = true;
    levelSelect.appendChild(opt);
  });
  
  document.getElementById('assign-modal').classList.add('active');
}

async function saveAssignment() {
  const domainId = document.getElementById('assign-domain').value || null;
  const levelId = document.getElementById('assign-level').value || null;
  
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/members/${_assignTarget}/organization`, {
      method: 'PUT',
      body: JSON.stringify({
        domain_id: domainId,
        hierarchy_level_id: levelId
      })
    });
    
    document.getElementById('assign-modal').classList.remove('active');
    await loadMembers(_clubId);
  } catch (err) {
    alert('Failed to assign: ' + (err.message || 'Unknown error'));
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
    card.style.cssText = 'display: flex; align-items: center; gap: var(--spacing-4); padding: var(--spacing-4);';
    
    const initial = req.user.name.charAt(0).toUpperCase();
    const requestedAt = new Date(req.created_at).toLocaleDateString();
    const domainLabel = req.domain ? ` — Domain: <strong>${req.domain.name}</strong>` : '';
    
    card.innerHTML = `
      <div class="avatar">${initial}</div>
      <div class="member-info" style="flex-grow: 1;">
        <div class="member-name">${req.user.name}</div>
        <div class="member-email">${req.user.email}${domainLabel}</div>
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
    await window.CoolKitAPI.api(`/clubs/${_clubId}/requests/${requestId}/approve`, { method: 'POST' });
    loadJoinRequests(_clubId);
    loadMembers(_clubId);
  } catch (err) {
    alert('Failed to approve request: ' + (err.message || 'Unknown error'));
  }
}

async function rejectRequest(requestId) {
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/requests/${requestId}/reject`, { method: 'POST' });
    loadJoinRequests(_clubId);
  } catch (err) {
    alert('Failed to reject request: ' + (err.message || 'Unknown error'));
  }
}

// ─── Role & Member Management (kept from original) ─────
async function changeRole(userId, newRole) {
  if (!newRole) return;
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/members/${userId}/role`, {
      method: 'PUT',
      body: JSON.stringify({ role: newRole })
    });
    loadMembers(_clubId);
  } catch (err) {
    alert('Failed to change role: ' + (err.message || 'Unknown error'));
  }
}

async function removeMember(userId) {
  if (!confirm('Are you sure you want to remove this member?')) return;
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/members/${userId}`, { method: 'DELETE' });
    loadMembers(_clubId);
  } catch (err) {
    alert('Failed to remove member: ' + (err.message || 'Unknown error'));
  }
}

async function leaveClub() {
  if (!confirm('Are you sure you want to leave this club?')) return;
  try {
    await window.CoolKitAPI.api(`/clubs/${_clubId}/members/me`, { method: 'DELETE' });
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
  const feed = document.getElementById('announcements-feed');
  feed.innerHTML = '';
  
  if (announcements.length === 0) {
    feed.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-4);">No announcements yet.</p>';
    return;
  }
  
  announcements.forEach(a => {
    const card = document.createElement('div');
    card.className = 'card';
    card.style.marginBottom = 'var(--spacing-3)';
    
    let priorityBadge = '';
    if (a.priority === 'urgent') priorityBadge = '<span class="badge badge-danger">Urgent</span>';
    else if (a.priority === 'important') priorityBadge = '<span class="badge badge-warning">Important</span>';
    
    const date = new Date(a.created_at).toLocaleString();
    const authorName = a.author?.name || 'Unknown';
    
    const isAdmin = _currentUserRole === 'owner' || _currentUserRole === 'admin';
    const deleteBtn = isAdmin ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose);" onclick="deleteAnnouncement('${a.id}')">Delete</button>` : '';
    
    card.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: flex-start;">
        <div>
          <h4 style="margin: 0 0 var(--spacing-2) 0;">${a.title} ${priorityBadge}</h4>
          <p style="margin: 0 0 var(--spacing-3) 0; color: var(--text-secondary); white-space: pre-wrap;">${a.content}</p>
        </div>
        ${deleteBtn}
      </div>
      <div style="color: var(--text-muted); font-size: var(--text-sm);">
        Posted by ${authorName} · ${date}
      </div>
    `;
    feed.appendChild(card);
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
    loadAnnouncements(_clubId);
  } catch (err) {
    alert('Failed to delete: ' + (err.message || 'Unknown error'));
  }
}

// ─── Tabs ──────────────────────────────
function setupTabs() {
  const tabs = document.querySelectorAll('.tab-btn');
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      tabs.forEach(t => t.classList.remove('active'));
      tab.classList.add('active');
      
      document.querySelectorAll('.tab-content').forEach(c => c.style.display = 'none');
      const target = tab.dataset.tab;
      const el = document.getElementById(`tab-${target}`);
      if (el) el.style.display = 'block';
    });
  });
}
