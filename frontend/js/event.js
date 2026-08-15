document.addEventListener('DOMContentLoaded', () => {
  const { getUser } = window.CoolKitAPI;
  const user = getUser();
  if (user) {
    document.getElementById('user-name').textContent = user.name;
    document.getElementById('user-avatar').textContent = user.name.charAt(0).toUpperCase();
  }

  const urlParams = new URLSearchParams(window.location.search);
  const eventId = urlParams.get('id');
  if (!eventId) { window.location.href = '/dashboard'; return; }

  window._eventId = eventId;
  window._isAdmin = false;
  window._members = [];

  loadEvent(eventId);
  setupTabs();
  setupAssignRoleForm(eventId);
  setupEditEventForm(eventId);
  setupCreateTaskForm(eventId);
  setupFinanceForm(eventId);
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
    
    // Pre-fill edit form
    document.getElementById('edit-event-title').value = event.title;
    document.getElementById('edit-event-desc').value = event.description || '';
    document.getElementById('edit-event-loc').value = event.location || '';
    const d = new Date(event.date);
    document.getElementById('edit-event-date').value = d.toISOString().slice(0, 16);

    // Fetch club members
    const clubRes = await window.CoolKitAPI.api(`/clubs/${event.club_id}/members`);
    const members = clubRes.data || [];
    window._members = members;
    const currentUser = window.CoolKitAPI.getUser();
    
    // Populate dropdowns
    const roleSelect = document.getElementById('role-user-id');
    const taskSelect = document.getElementById('task-assignee');
    roleSelect.innerHTML = '<option value="" disabled selected>Select a member...</option>';
    taskSelect.innerHTML = '<option value="" disabled selected>Select a member...</option>';
    
    members.forEach(m => {
      if (m.user.id === currentUser?.id && (m.role === 'owner' || m.role === 'admin')) {
        window._isAdmin = true;
      }
      const opt1 = document.createElement('option');
      opt1.value = m.user.id;
      opt1.textContent = `${m.user.name} (${m.user.email})`;
      roleSelect.appendChild(opt1);
      const opt2 = opt1.cloneNode(true);
      taskSelect.appendChild(opt2);
    });

    if (window._isAdmin) {
      document.getElementById('btn-assign-role').style.display = 'block';
      document.getElementById('btn-create-task').style.display = 'block';
      document.getElementById('btn-add-finance').style.display = 'block';
      document.getElementById('event-admin-actions').style.display = 'flex';
    }
    
    loadRoles(id);
    loadTasks(id);
    loadFinance(id);
  } catch (err) {
    console.error('Failed to load event:', err);
  }
}

// ─── Roles ──────────────────────────────
async function loadRoles(id) {
  try {
    const res = await window.CoolKitAPI.api(`/events/${id}/roles`);
    renderRoles(res.data || []);
  } catch (err) { console.error('Failed to load roles:', err); }
}

function renderRoles(roles) {
  const list = document.getElementById('roles-list');
  list.innerHTML = '';
  if (roles.length === 0) {
    list.innerHTML = `<div class="empty-state"><div style="font-size: 48px; margin-bottom: var(--spacing-4);">👥</div><h3>No roles assigned yet</h3><p>Assign responsibilities to club members for this event.</p></div>`;
    return;
  }
  roles.forEach(role => {
    const el = document.createElement('div');
    el.className = 'card member-card';
    el.innerHTML = `
      <div class="member-info" style="display: flex; align-items: center; gap: var(--spacing-3);">
        <div class="avatar">${role.user.name.charAt(0).toUpperCase()}</div>
        <div>
          <div class="member-name">${role.user.name}</div>
          <div class="member-email">${role.user.email}</div>
        </div>
      </div>
      <div style="display: flex; align-items: center; gap: var(--spacing-2);">
        <span class="badge badge-primary">${role.role_name}</span>
        ${window._isAdmin ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose);" onclick="removeRole('${role.id}')">✕</button>` : ''}
      </div>
    `;
    list.appendChild(el);
  });
}

async function removeRole(roleId) {
  if (!confirm('Remove this role assignment?')) return;
  try {
    await window.CoolKitAPI.api(`/events/${window._eventId}/roles/${roleId}`, { method: 'DELETE' });
    loadRoles(window._eventId);
  } catch (err) { alert('Failed to remove role.'); }
}

function setupAssignRoleForm(eventId) {
  const form = document.getElementById('assign-role-form');
  if (!form) return;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
      await window.CoolKitAPI.api(`/events/${eventId}/roles`, {
        method: 'POST',
        body: JSON.stringify({
          user_id: document.getElementById('role-user-id').value,
          role_name: document.getElementById('role-name').value
        })
      });
      document.getElementById('assign-role-modal').classList.remove('active');
      form.reset();
      loadRoles(eventId);
    } catch (err) {
      const errorDiv = document.getElementById('role-error');
      errorDiv.textContent = err.message || 'Failed to assign role.';
      errorDiv.style.display = 'block';
    }
  });
}

// ─── Edit / Delete Event ──────────────────────────────
function setupEditEventForm(eventId) {
  const form = document.getElementById('edit-event-form');
  if (!form) return;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
      await window.CoolKitAPI.api(`/events/${eventId}`, {
        method: 'PUT',
        body: JSON.stringify({
          title: document.getElementById('edit-event-title').value,
          description: document.getElementById('edit-event-desc').value,
          date: new Date(document.getElementById('edit-event-date').value).toISOString(),
          location: document.getElementById('edit-event-loc').value
        })
      });
      document.getElementById('edit-event-modal').classList.remove('active');
      loadEvent(eventId);
    } catch (err) { alert('Failed to update event: ' + (err.message || '')); }
  });
}

async function deleteEvent() {
  if (!confirm('Are you sure you want to delete this event? This action cannot be undone.')) return;
  try {
    await window.CoolKitAPI.api(`/events/${window._eventId}`, { method: 'DELETE' });
    window.location.href = `/club?id=${currentEvent.club_id}`;
  } catch (err) { alert('Failed to delete event.'); }
}

// ─── Tasks ──────────────────────────────
async function loadTasks(id) {
  try {
    const res = await window.CoolKitAPI.api(`/events/${id}/tasks`);
    renderTasks(res.data || []);
  } catch (err) { console.error('Failed to load tasks:', err); }
}

function renderTasks(tasks) {
  const todoEl = document.getElementById('tasks-todo');
  const ipEl = document.getElementById('tasks-in-progress');
  const doneEl = document.getElementById('tasks-done');
  todoEl.innerHTML = ''; ipEl.innerHTML = ''; doneEl.innerHTML = '';

  const byStatus = { todo: [], in_progress: [], done: [] };
  tasks.forEach(t => {
    if (!byStatus[t.status]) byStatus[t.status] = [];
    byStatus[t.status].push(t);
  });

  const renderCard = (task) => {
    const card = document.createElement('div');
    card.className = 'card';
    card.style.marginBottom = 'var(--spacing-3)';
    const assignee = task.assigned_to?.name || 'Unassigned';
    const due = task.due_date ? new Date(task.due_date).toLocaleDateString() : 'No deadline';
    
    let statusBtns = '';
    if (task.status === 'todo') statusBtns = `<button class="btn btn-ghost btn-sm" onclick="updateTaskStatus('${task.id}', 'in_progress')">▶ Start</button>`;
    else if (task.status === 'in_progress') statusBtns = `<button class="btn btn-ghost btn-sm" onclick="updateTaskStatus('${task.id}', 'done')">✅ Done</button>`;
    else statusBtns = `<span class="badge badge-success">Completed</span>`;

    card.innerHTML = `
      <h4 style="margin: 0 0 var(--spacing-2) 0; font-size: var(--text-base);">${task.title}</h4>
      <p style="margin: 0 0 var(--spacing-2) 0; color: var(--text-secondary); font-size: var(--text-sm);">${task.description || ''}</p>
      <div style="display: flex; justify-content: space-between; align-items: center; font-size: var(--text-xs); color: var(--text-muted);">
        <span>👤 ${assignee} • 📅 ${due}</span>
        <div style="display: flex; gap: var(--spacing-1);">${statusBtns}${window._isAdmin ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose);" onclick="deleteTask('${task.id}')">🗑</button>` : ''}</div>
      </div>
    `;
    return card;
  };

  byStatus.todo.forEach(t => todoEl.appendChild(renderCard(t)));
  byStatus.in_progress.forEach(t => ipEl.appendChild(renderCard(t)));
  byStatus.done.forEach(t => doneEl.appendChild(renderCard(t)));

  if (byStatus.todo.length === 0) todoEl.innerHTML = '<p style="color: var(--text-muted); font-size: var(--text-sm);">No tasks</p>';
  if (byStatus.in_progress.length === 0) ipEl.innerHTML = '<p style="color: var(--text-muted); font-size: var(--text-sm);">No tasks</p>';
  if (byStatus.done.length === 0) doneEl.innerHTML = '<p style="color: var(--text-muted); font-size: var(--text-sm);">No tasks</p>';
}

async function updateTaskStatus(taskId, status) {
  try {
    await window.CoolKitAPI.api(`/tasks/${taskId}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status })
    });
    loadTasks(window._eventId);
  } catch (err) { alert('Failed to update task: ' + (err.message || '')); }
}

async function deleteTask(taskId) {
  if (!confirm('Delete this task?')) return;
  try {
    await window.CoolKitAPI.api(`/tasks/${taskId}`, { method: 'DELETE' });
    loadTasks(window._eventId);
  } catch (err) { alert('Failed to delete task.'); }
}

function setupCreateTaskForm(eventId) {
  const form = document.getElementById('create-task-form');
  if (!form) return;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const body = {
      title: document.getElementById('task-title').value,
      description: document.getElementById('task-desc').value,
      assigned_to_id: document.getElementById('task-assignee').value
    };
    const due = document.getElementById('task-due').value;
    if (due) body.due_date = new Date(due).toISOString();
    
    try {
      await window.CoolKitAPI.api(`/events/${eventId}/tasks`, {
        method: 'POST',
        body: JSON.stringify(body)
      });
      document.getElementById('create-task-modal').classList.remove('active');
      form.reset();
      loadTasks(eventId);
    } catch (err) { alert('Failed to create task: ' + (err.message || '')); }
  });
}

// ─── Finance ──────────────────────────────
async function loadFinance(id) {
  try {
    const res = await window.CoolKitAPI.api(`/events/${id}/finance`);
    const data = res.data;
    
    const income = data.total_income || 0;
    const expense = data.total_expense || 0;
    const balance = data.balance || 0;
    
    document.getElementById('finance-income').textContent = `₹${income.toLocaleString()}`;
    document.getElementById('finance-expense').textContent = `₹${expense.toLocaleString()}`;
    document.getElementById('finance-balance').textContent = `₹${balance.toLocaleString()}`;
    document.getElementById('finance-balance').style.color = balance >= 0 ? 'var(--accent-green)' : 'var(--accent-rose)';
    
    renderFinanceEntries(data.entries || []);
  } catch (err) { console.error('Failed to load finance:', err); }
}

function renderFinanceEntries(entries) {
  const container = document.getElementById('finance-entries');
  container.innerHTML = '';
  
  if (entries.length === 0) {
    container.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-4);">No finance entries yet.</p>';
    return;
  }
  
  entries.forEach(entry => {
    const el = document.createElement('div');
    el.className = 'card';
    el.style.marginBottom = 'var(--spacing-3)';
    const isIncome = entry.type === 'income';
    const date = new Date(entry.date).toLocaleDateString();
    const creator = entry.created_by?.name || 'Unknown';
    
    el.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <div>
          <div style="display: flex; align-items: center; gap: var(--spacing-2);">
            <span style="font-size: 20px;">${isIncome ? '💰' : '💸'}</span>
            <strong>${entry.description || entry.category}</strong>
            <span class="badge ${isIncome ? 'badge-success' : 'badge-danger'}">${entry.category}</span>
          </div>
          <div style="font-size: var(--text-xs); color: var(--text-muted); margin-top: var(--spacing-1);">By ${creator} • ${date}</div>
        </div>
        <div style="display: flex; align-items: center; gap: var(--spacing-2);">
          <span style="font-weight: 700; font-size: var(--text-lg); color: ${isIncome ? 'var(--accent-green)' : 'var(--accent-rose)'};">${isIncome ? '+' : '-'}₹${entry.amount.toLocaleString()}</span>
          ${window._isAdmin ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose);" onclick="deleteFinanceEntry('${entry.id}')">✕</button>` : ''}
        </div>
      </div>
    `;
    container.appendChild(el);
  });
}

async function deleteFinanceEntry(id) {
  if (!confirm('Delete this finance entry?')) return;
  try {
    await window.CoolKitAPI.api(`/finance/${id}`, { method: 'DELETE' });
    loadFinance(window._eventId);
  } catch (err) { alert('Failed to delete entry.'); }
}

function setupFinanceForm(eventId) {
  const form = document.getElementById('add-finance-form');
  if (!form) return;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
      await window.CoolKitAPI.api(`/events/${eventId}/finance`, {
        method: 'POST',
        body: JSON.stringify({
          type: document.getElementById('finance-type').value,
          category: document.getElementById('finance-category').value,
          amount: parseFloat(document.getElementById('finance-amount').value),
          description: document.getElementById('finance-desc').value,
          date: new Date(document.getElementById('finance-date').value).toISOString()
        })
      });
      document.getElementById('add-finance-modal').classList.remove('active');
      form.reset();
      loadFinance(eventId);
    } catch (err) { alert('Failed to add entry: ' + (err.message || '')); }
  });
}

// ─── Tabs ──────────────────────────────
function setupTabs() {
  const buttons = document.querySelectorAll('.tab-btn');
  const contents = document.querySelectorAll('.tab-content');
  buttons.forEach(btn => {
    btn.addEventListener('click', () => {
      buttons.forEach(b => b.classList.remove('active'));
      contents.forEach(c => { c.classList.remove('active'); c.style.display = 'none'; });
      btn.classList.add('active');
      const target = document.getElementById(btn.dataset.tab);
      target.classList.add('active');
      target.style.display = 'block';
    });
  });
}
