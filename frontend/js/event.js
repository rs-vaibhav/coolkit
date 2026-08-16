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
  setupUploadDocForm(eventId);
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
      document.getElementById('btn-upload-doc').style.display = 'block';
      document.getElementById('btn-create-chat').style.display = 'block';
      document.getElementById('event-admin-actions').style.display = 'flex';
      
      // Show form creation/edit buttons based on survey status
      if (event.formbricks_survey_id) {
        document.getElementById('btn-edit-form').style.display = 'block';
        document.getElementById('btn-show-qr').style.display = 'block';
      } else {
        document.getElementById('btn-create-form').style.display = 'block';
      }
    }
    
    // Store formbricks info for later use
    if (event.formbricks_survey_id) {
      document.getElementById('qr-formbricks-link').value = `https://app.formbricks.com/s/${event.formbricks_survey_id}`;
      document.getElementById('btn-show-qr').style.display = 'block';
    }
    
    // Load QR code info if available
    if (event.qr_code_url) {
      document.getElementById('qr-formbricks-link').value = event.qr_code_url;
    }
    
    // Load Matrix room if available
    if (event.matrix_room_id) {
      renderMatrixChat(event.matrix_room_id);
    }
    
    loadRoles(id);
    loadTasks(id);
    loadFinance(id);
    loadRSVP(id);
    loadDocuments(id);
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
    
    const proofFile = document.getElementById('finance-proof').files[0];
    let proofUrl = null;
    
    // Upload proof image if provided
    if (proofFile) {
      const formData = new FormData();
      formData.append('file', proofFile);
      try {
        const uploadRes = await fetch('/api/v1/upload', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${window.CoolKitAPI.getToken()}`
          },
          body: formData
        });
        if (uploadRes.ok) {
          const uploadData = await uploadRes.json();
          proofUrl = uploadData.url;
        }
      } catch (err) {
        console.error('Failed to upload proof:', err);
      }
    }
    
    try {
      await window.CoolKitAPI.api(`/events/${eventId}/finance`, {
        method: 'POST',
        body: JSON.stringify({
          type: document.getElementById('finance-type').value,
          category: document.getElementById('finance-category').value,
          amount: parseFloat(document.getElementById('finance-amount').value),
          description: document.getElementById('finance-desc').value,
          date: new Date(document.getElementById('finance-date').value).toISOString(),
          proof_image: proofUrl
        })
      });
      document.getElementById('add-finance-modal').classList.remove('active');
      form.reset();
      loadFinance(eventId);
    } catch (err) { alert('Failed to add entry: ' + (err.message || '')); }
  });
}

// ─── RSVP Functions ──────────────────────────────
async function loadRSVP(eventId) {
  try {
    // Load RSVP counts
    const countsRes = await window.CoolKitAPI.api(`/events/${eventId}/rsvp-counts`);
    const counts = countsRes.data;
    document.getElementById('rsvp-going').textContent = counts.going || 0;
    document.getElementById('rsvp-interested').textContent = counts.interested || 0;
    document.getElementById('rsvp-not-going').textContent = counts.not_going || 0;
    
    // Load user's RSVP
    try {
      const userRes = await window.CoolKitAPI.api(`/events/${eventId}/rsvp`);
      const userRsvp = userRes.data;
      renderUserRSVP(userRsvp);
    } catch (err) {
      // No RSVP yet
      renderUserRSVP(null);
    }
    
    // Load attendee list (for admins)
    if (window._isAdmin) {
      const attendeesRes = await window.CoolKitAPI.api(`/events/${eventId}/rsvps`);
      renderAttendeeList(attendeesRes.data || []);
    }
  } catch (err) {
    console.error('Failed to load RSVP:', err);
  }
}

function renderUserRSVP(rsvp) {
  const statusDiv = document.getElementById('rsvp-user-status');
  const actionsDiv = document.getElementById('rsvp-actions');
  
  if (!rsvp) {
    statusDiv.innerHTML = '<p style="color: var(--text-muted);">You haven\'t responded yet.</p>';
    actionsDiv.style.display = 'flex';
    return;
  }
  
  const statusLabels = {
    'going': '✅ Going',
    'interested': '🤔 Interested',
    'not_going': '❌ Not Going'
  };
  
  statusDiv.innerHTML = `<span class="badge badge-primary" style="font-size: var(--text-base); padding: var(--spacing-2) var(--spacing-3);">${statusLabels[rsvp.status] || 'Unknown'}</span>`;
  actionsDiv.style.display = 'none';
}

async function submitRSVP(status) {
  try {
    await window.CoolKitAPI.api(`/events/${window._eventId}/rsvp`, {
      method: 'POST',
      body: JSON.stringify({ status })
    });
    loadRSVP(window._eventId);
  } catch (err) {
    alert('Failed to submit RSVP: ' + (err.message || ''));
  }
}

function renderAttendeeList(rsvps) {
  const container = document.getElementById('rsvp-attendees-list');
  container.innerHTML = '';
  
  if (rsvps.length === 0) {
    container.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: var(--spacing-4);">No responses yet.</p>';
    return;
  }
  
  const byStatus = { going: [], interested: [], not_going: [] };
  rsvps.forEach(r => {
    if (byStatus[r.status]) byStatus[r.status].push(r);
  });
  
  const renderSection = (status, title, icon) => {
    if (byStatus[status].length === 0) return;
    const section = document.createElement('div');
    section.style.marginBottom = 'var(--spacing-4)';
    section.innerHTML = `<h5 style="margin-bottom: var(--spacing-2); color: var(--text-secondary);">${icon} ${title} (${byStatus[status].length})</h5>`;
    const list = document.createElement('div');
    byStatus[status].forEach(r => {
      const item = document.createElement('div');
      item.className = 'card';
      item.style.marginBottom = 'var(--spacing-2)';
      item.style.padding = 'var(--spacing-2)';
      item.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div style="display: flex; align-items: center; gap: var(--spacing-2);">
            <div class="avatar" style="width: 32px; height: 32px; font-size: 14px;">${(r.user?.name || '?').charAt(0).toUpperCase()}</div>
            <div>
              <div style="font-weight: 500;">${r.user?.name || 'Unknown'}</div>
              <div style="font-size: var(--text-xs); color: var(--text-muted);">${r.user?.email || ''}</div>
            </div>
          </div>
          ${r.checked_in ? '<span class="badge badge-success">Checked In</span>' : ''}
        </div>
      `;
      list.appendChild(item);
    });
    section.appendChild(list);
    container.appendChild(section);
  };
  
  renderSection('going', 'Going', '✅');
  renderSection('interested', 'Interested', '🤔');
  renderSection('not_going', 'Not Going', '❌');
}

function showQRCode() {
  const link = document.getElementById('qr-formbricks-link').value;
  const container = document.getElementById('qr-code-container');
  
  if (!link) {
    container.innerHTML = '<p style="color: var(--accent-rose);">No QR code configured for this event. Admin should update the event with Formbricks survey link.</p>';
  } else {
    // Generate QR code using a simple API
    const qrUrl = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(link)}`;
    container.innerHTML = `<img src="${qrUrl}" alt="QR Code" style="border: 1px solid var(--border-color); border-radius: var(--radius-md); padding: var(--spacing-3); background: white;">`;
  }
  
  document.getElementById('qr-code-modal').classList.add('active');
}

// ─── Documents Functions ──────────────────────────────
async function loadDocuments(eventId) {
  try {
    const res = await window.CoolKitAPI.api(`/events/${eventId}/documents`);
    const docs = res.data || [];
    renderDocuments(docs);
  } catch (err) {
    console.error('Failed to load documents:', err);
  }
}

function renderDocuments(docs) {
  const container = document.getElementById('documents-list');
  container.innerHTML = '';
  
  if (docs.length === 0) {
    container.innerHTML = '<div class="empty-state"><div style="font-size: 48px; margin-bottom: var(--spacing-4);">📄</div><h3>No documents uploaded</h3><p>Admins can upload permission letters, invoices, and other event documents.</p></div>';
    return;
  }
  
  docs.forEach(doc => {
    const el = document.createElement('div');
    el.className = 'card';
    el.style.marginBottom = 'var(--spacing-3)';
    const date = new Date(doc.created_at).toLocaleDateString();
    const typeIcons = {
      'permission_letter': '📜',
      'invoice': '🧾',
      'receipt': '💵',
      'contract': '📋',
      'other': '📄'
    };
    
    el.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <div style="display: flex; align-items: center; gap: var(--spacing-3);">
          <div style="font-size: 32px;">${typeIcons[doc.document_type] || '📄'}</div>
          <div>
            <div style="font-weight: 600;">${doc.title}</div>
            <div style="font-size: var(--text-sm); color: var(--text-muted);">${doc.description || 'No description'} • Uploaded ${date}</div>
          </div>
        </div>
        <div style="display: flex; align-items: center; gap: var(--spacing-2);">
          <a href="${doc.file_url}" target="_blank" class="btn btn-primary btn-sm">Download</a>
          ${window._isAdmin ? `<button class="btn btn-ghost btn-sm" style="color: var(--accent-rose);" onclick="deleteDocument('${doc.id}')">🗑</button>` : ''}
        </div>
      </div>
    `;
    container.appendChild(el);
  });
}

async function deleteDocument(docId) {
  if (!confirm('Delete this document?')) return;
  try {
    await window.CoolKitAPI.api(`/documents/${docId}`, { method: 'DELETE' });
    loadDocuments(window._eventId);
  } catch (err) {
    alert('Failed to delete document.');
  }
}

function setupUploadDocForm(eventId) {
  const form = document.getElementById('upload-doc-form');
  if (!form) return;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const file = document.getElementById('doc-file').files[0];
    if (!file) {
      alert('Please select a file to upload.');
      return;
    }
    
    const formData = new FormData();
    formData.append('file', file);
    formData.append('title', document.getElementById('doc-title').value);
    formData.append('document_type', document.getElementById('doc-type').value);
    formData.append('description', document.getElementById('doc-desc').value);
    
    try {
      await window.CoolKitAPI.api(`/events/${eventId}/documents`, {
        method: 'POST',
        body: formData,
        headers: {} // Let browser set Content-Type for FormData
      });
      document.getElementById('upload-doc-modal').classList.remove('active');
      form.reset();
      loadDocuments(eventId);
    } catch (err) {
      alert('Failed to upload document: ' + (err.message || ''));
    }
  });
}

// ─── Matrix Chat Functions ──────────────────────────────
function createMatrixRoom() {
  const roomName = `Event-${currentEvent.title.replace(/\s+/g, '-')}`;
  alert('Matrix room creation requires server-side integration. Please configure matrix_room_id in the event settings manually or use Element web interface to create a room and share the room ID.');
  // In production, this would call an API endpoint to create a Matrix room
}

function renderMatrixChat(roomId) {
  const container = document.getElementById('chat-container');
  container.innerHTML = `
    <div style="width: 100%; height: 500px; display: flex; flex-direction: column; align-items: center; justify-content: center; background: var(--bg-secondary); border-radius: var(--radius-lg); padding: var(--spacing-6);">
      <div style="font-size: 48px; margin-bottom: var(--spacing-4);">💬</div>
      <h4 style="margin-bottom: var(--spacing-2);">Matrix Chat Room</h4>
      <p style="color: var(--text-muted); margin-bottom: var(--spacing-4);">Room ID: ${roomId}</p>
      <a href="https://app.element.io/#/room/${roomId}" target="_blank" class="btn btn-primary" style="display: inline-flex; align-items: center; gap: var(--spacing-2);">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/></svg>
        Open in Element
      </a>
      <p style="font-size: var(--text-sm); color: var(--text-muted); margin-top: var(--spacing-3);">Or join using any Matrix client with room ID: ${roomId}</p>
    </div>
  `;
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

// ─── Formbricks Functions ──────────────────────────────
function showFormbricksSetup() {
  document.getElementById('formbricks-setup-modal').classList.add('active');
}

document.getElementById('setup-formbricks-form')?.addEventListener('submit', async (e) => {
  e.preventDefault();
  const envId = document.getElementById('formbricks-env-id').value.trim();
  
  if (!envId) {
    alert('Please enter a valid Environment ID');
    return;
  }
  
  try {
    await window.CoolKitAPI.api(`/events/${window._eventId}`, {
      method: 'PUT',
      body: JSON.stringify({ formbricks_env_id: envId })
    });
    
    alert('Formbricks Environment ID saved successfully!');
    document.getElementById('formbricks-setup-modal').classList.remove('active');
    document.getElementById('btn-setup-formbricks').style.display = 'none';
    document.getElementById('btn-create-form').style.display = 'block';
    
    // Reload event to update state
    loadEvent(window._eventId);
  } catch (err) {
    alert('Failed to save Environment ID: ' + (err.message || ''));
  }
});

async function createForm() {
  const modal = document.getElementById('create-form-modal');
  const resultDiv = document.getElementById('create-form-result');
  const openBtn = document.getElementById('btn-open-formbricks');
  
  document.getElementById('create-form-title').textContent = 'Create Registration Form';
  resultDiv.innerHTML = '<p style="color: var(--text-secondary);">Creating form in Formbricks...</p>';
  openBtn.style.display = 'none';
  modal.classList.add('active');
  
  try {
    // Create form using backend API with stored API key
    const res = await window.CoolKitAPI.api(`/events/${window._eventId}/formbricks/create`, {
      method: 'POST'
    });
    
    const survey = res.data;
    const surveyLink = `https://app.formbricks.com/s/${survey.id}`;
    
    resultDiv.innerHTML = `
      <div style="background: var(--bg-success); padding: var(--spacing-4); border-radius: var(--radius-md); margin-bottom: var(--spacing-3);">
        <p style="margin-bottom: var(--spacing-2); color: var(--accent-green);"><strong>✅ Form Created Successfully!</strong></p>
        <p style="font-size: var(--text-sm); color: var(--text-secondary); margin-bottom: var(--spacing-2);">
          A default registration form has been created for your event.
        </p>
        <p style="font-size: var(--text-sm); color: var(--text-secondary);">
          Survey ID: ${survey.id}
        </p>
      </div>
      <p style="font-size: var(--text-sm); color: var(--text-muted);">
        Click the button below to customize form fields (add/remove questions) in Formbricks editor.
      </p>
    `;
    
    openBtn.dataset.surveyId = survey.id;
    openBtn.style.display = 'inline-block';
    
    // Update button visibility
    document.getElementById('btn-create-form').style.display = 'none';
    document.getElementById('btn-edit-form').style.display = 'block';
    document.getElementById('btn-show-qr').style.display = 'block';
    
    // Reload event to update formbricks_survey_id
    loadEvent(window._eventId);
  } catch (err) {
    resultDiv.innerHTML = `<p style="color: var(--accent-rose);">❌ Failed to create form: ${err.message || 'Unknown error'}</p>`;
  }
}

async function editForm() {
  const surveyId = currentEvent.formbricks_survey_id;
  if (!surveyId) {
    alert('No form found for this event. Please create one first.');
    return;
  }
  
  openFormbricksEditor(surveyId);
}

function openFormbricksEditor(surveyId) {
  const id = surveyId || document.getElementById('btn-open-formbricks')?.dataset.surveyId;
  if (!id) {
    alert('No survey ID available');
    return;
  }
  
  const editorUrl = `https://app.formbricks.com/surveys/${id}/edit`;
  window.open(editorUrl, '_blank');
}
