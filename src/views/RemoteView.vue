<script setup>
import { ref, reactive } from 'vue';

// --- DATA SIMULASI ---
const remotes = ref([
  { id: 1, name: 'Google Drive Utama', type: 'Google Drive', status: 'Online', lastCheck: '2025-10-07 08:15' },
  { id: 2, name: 'S3-Archive-Backup', type: 'Amazon S3', status: 'Online', lastCheck: '2025-10-07 09:01' },
  { id: 3, name: 'Azure-West-Storage', type: 'Microsoft Azure', status: 'Offline', lastCheck: '2025-10-06 17:30' },
  { id: 4, name: 'Local-NAS-01', type: 'SFTP', status: 'Online', lastCheck: '2025-10-07 09:10' },
]);

// --- STATE MANAGEMENT UNTUK FORM ---
const showAddForm = ref(false); // Mengontrol visibilitas form
const isEditing = ref(false); // Menandakan apakah form dalam mode edit atau tambah baru

// Objek reaktif untuk menampung data dari form
const newRemoteForm = reactive({
  id: null,
  name: '',
  type: 'Google Drive', // Nilai default
});

// --- FUNGSI-FUNGSI ---

// Fungsi untuk menampilkan form tambah data
const handleAddNew = () => {
  isEditing.value = false;
  // Reset form
  newRemoteForm.id = null;
  newRemoteForm.name = '';
  newRemoteForm.type = 'Google Drive';
  showAddForm.value = true;
};

// Fungsi untuk menampilkan form edit data
const handleEdit = (remote) => {
  isEditing.value = true;
  // Isi form dengan data yang ada
  newRemoteForm.id = remote.id;
  newRemoteForm.name = remote.name;
  newRemoteForm.type = remote.type;
  showAddForm.value = true;
};

// Fungsi untuk menyimpan data (baik baru maupun editan)
const handleSubmit = () => {
  if (!newRemoteForm.name) {
    alert('Remote name is required!');
    return;
  }

  if (isEditing.value) {
    // Logika untuk update data yang ada
    const index = remotes.value.findIndex(r => r.id === newRemoteForm.id);
    if (index !== -1) {
      remotes.value[index] = { ...remotes.value[index], ...newRemoteForm };
    }
  } else {
    // Logika untuk menambah data baru
    const newId = remotes.value.length > 0 ? Math.max(...remotes.value.map(r => r.id)) + 1 : 1;
    remotes.value.push({
      id: newId,
      name: newRemoteForm.name,
      type: newRemoteForm.type,
      status: 'Pending', // Status default untuk remote baru
      lastCheck: new Date().toISOString().slice(0, 16).replace('T', ' '),
    });
  }

  // Sembunyikan form setelah submit
  showAddForm.value = false;
};

// Fungsi untuk menghapus data
const handleDelete = (id) => {
  if (confirm('Are you sure you want to delete this remote?')) {
    remotes.value = remotes.value.filter(remote => remote.id !== id);
  }
};

// Helper untuk styling status
const getStatusClass = (status) => {
  if (status === 'Online') return 'status-online';
  if (status === 'Offline') return 'status-offline';
  return 'status-pending';
};
</script>

<template>
  <div class="remote-content">
    <header class="main-header">
      <h1>Remote Management</h1>
      <p>Manage your cloud and local storage connections.</p>
    </header>

    <!-- FORM UNTUK TAMBAH/EDIT REMOTE -->
    <div v-if="showAddForm" class="card form-card">
      <h2>{{ isEditing ? 'Edit Remote' : 'Add New Remote' }}</h2>
      <form @submit.prevent="handleSubmit" class="add-remote-form">
        <div class="form-group">
          <label for="remote-name">Remote Name</label>
          <input id="remote-name" v-model="newRemoteForm.name" type="text" placeholder="e.g., My Google Drive" required class="input-field" />
        </div>
        <div class="form-group">
          <label for="remote-type">Remote Type</label>
          <select id="remote-type" v-model="newRemoteForm.type" class="input-field">
            <option>Google Drive</option>
            <option>Amazon S3</option>
            <option>Microsoft Azure</option>
            <option>SFTP</option>
            <option>Local</option>
          </select>
        </div>
        <div class="form-actions">
          <button type="button" @click="showAddForm = false" class="btn btn-secondary">Cancel</button>
          <button type="submit" class="btn btn-primary">{{ isEditing ? 'Update' : 'Save' }}</button>
        </div>
      </form>
    </div>

    <!-- TABEL DAFTAR REMOTE -->
    <div class="card card-table">
      <div class="table-toolbar">
        <h2>Remote List</h2>
        <button @click="handleAddNew" class="btn btn-primary">
          <!-- SVG icon for plus sign -->
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-plus-lg" viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M8 2a.5.5 0 0 1 .5.5v5h5a.5.5 0 0 1 0 1h-5v5a.5.5 0 0 1-1 0v-5h-5a.5.5 0 0 1 0-1h5v-5A.5.5 0 0 1 8 2z"/>
          </svg>
          Add New
        </button>
      </div>
      <div class="table-responsive">
        <table class="data-table">
          <thead>
            <tr>
              <th>Remote Name</th>
              <th>Type</th>
              <th>Status</th>
              <th>Last Check</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="remotes.length === 0">
              <td colspan="5" class="text-center">No remote connections found.</td>
            </tr>
            <tr v-for="remote in remotes" :key="remote.id">
              <td>{{ remote.name }}</td>
              <td>{{ remote.type }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(remote.status)]">{{ remote.status }}</span>
              </td>
              <td>{{ remote.lastCheck }}</td>
              <td class="action-buttons">
                <button @click="handleEdit(remote)" class="btn-icon btn-edit" title="Edit">
                  <!-- SVG icon for edit -->
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708l-3-3zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207l6.5-6.5zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.499.499 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11l.178-.178z"/></svg>
                </button>
                <button @click="handleDelete(remote.id)" class="btn-icon btn-delete" title="Delete">
                  <!-- SVG icon for delete -->
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/><path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/></svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* --- MAIN LAYOUT & HEADER --- */
.remote-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.main-header {
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
}

.main-header h1 {
  font-size: 1.8rem;
  font-weight: 700;
  margin-bottom: 4px;
}

.main-header p {
  color: #6c757d;
}

/* --- FORM STYLING --- */
.form-card {
  padding: 25px;
  background-color: #f8f9fa;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 10px;
}

.form-card h2 {
  font-size: 1.4rem;
  margin-bottom: 20px;
}

.add-remote-form {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  align-items: flex-end;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  margin-bottom: 8px;
  font-weight: 500;
  color: #495057;
}

.input-field {
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 1rem;
}

.input-field:focus {
  border-color: var(--primary-color);
  outline: none;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  grid-column: 1 / -1; /* Make it span all columns */
  margin-top: 10px;
}

/* --- CARD & TABLE STYLING --- */
.card {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-color);
  overflow: hidden; /* Ensures padding and border radius work together */
}

.card-table {
  padding: 0;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.table-toolbar h2 {
  font-size: 1.3rem;
  margin: 0;
}

.table-responsive {
  overflow-x: auto;
  padding: 0 20px 20px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 15px;
}

.data-table th, .data-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.data-table th {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #495057;
}

.data-table tbody tr:hover {
  background-color: #f1f3f5;
}

.text-center {
  text-align: center;
  color: #6c757d;
  padding: 20px;
}

/* --- STATUS BADGES --- */
.status-badge {
  display: inline-block;
  padding: 5px 10px;
  border-radius: 15px;
  font-weight: 500;
  font-size: 0.8rem;
  color: white;
  min-width: 70px;
  text-align: center;
}

.status-online { background-color: #28a745; } /* Green */
.status-offline { background-color: #dc3545; } /* Red */
.status-pending { background-color: #6c757d; } /* Gray */

/* --- BUTTONS --- */
.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}
.btn-primary:hover {
  background-color: var(--primary-hover);
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}
.btn-secondary:hover {
  background-color: #5a6268;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  padding: 5px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.btn-icon svg {
  width: 16px;
  height: 16px;
}

.btn-edit { color: #007bff; }
.btn-edit:hover { background-color: rgba(0, 123, 255, 0.1); }

.btn-delete { color: #dc3545; }
.btn-delete:hover { background-color: rgba(220, 53, 69, 0.1); }
</style>
