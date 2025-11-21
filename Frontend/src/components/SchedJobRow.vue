<template>
  <tr>
    <div class="job-name">
    <td><strong>{{ job.job_name }}</strong></td>
        </div>
    <td>{{ job.source_path }}</td>
    <td>{{ job.gdrive_target }}</td>
    <td>{{ job.last_run || 'N/A' }}</td>

    <td>
      <span class="status" :class="job.status.toLowerCase()">
        {{ job.status }}
      </span>
    </td>
    
    <td>{{ job.next_run || 'N/A' }}</td>
    
    <td class="job-actions-col"> 
      <div class="actions">
        <!-- ✅ Edit Button -->
        <button 
          @click="$emit('edit', job.id)" 
          class="action-btn edit" 
          title="Edit Job"
        >
          ✏️
        </button>

        <!-- View Script Button -->
        <button 
          @click="$emit('view-script', job.id)" 
          class="action-btn view" 
          title="View Script"
        >
          (i)
        </button>

        <!-- Delete Button -->
        <button 
          @click="confirmDelete" 
          class="action-btn delete" 
          title="Delete Job"
        >
          🗑
        </button>
      </div>
    </td>
  </tr>
</template>

<script setup>
const props = defineProps({
  job: {
    type: Object,
    required: true
  }
});

// ✅ Tambah 'edit' ke emits
const emit = defineEmits(['trigger', 'view-script', 'edit', 'delete']); 

function confirmDelete() {
  if (confirm(`PERINGATAN! Menghapus Job Terjadwal '${props.job.job_name}' juga akan menghapus jadwal CRON.\nApakah Anda yakin ingin menghapus job ini?`)) {
    emit('delete', props.job.id, props.job.job_name);
  }
}
</script>

<style scoped>
/* --- STATUS BADGE (mengikuti style remote) --- */
.status {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  color: white;
}

.status.pending {
  background-color: #7f8c8d;
}

.status.running {
  background-color: #3498db;
}

.status.completed,
.status.success {
  background-color: #27ae60;
}

.status.failed,
.status.fail_source_check, 
.status.fail_pre_script, 
.status.fail_rclone,
.status.fail_post_script {
  background-color: #e74c3c;
}

/* --- TABEL CELL --- */
td {
  padding: 0.75rem;
  vertical-align: middle;
}

strong {
  color: #1a1a1a;
}

/* --- KOLUM A K S I --- */
.job-actions-col {
  width: 150px;
}

.job-name {
  font-size: 0.90rem; /* atau 0.65rem kalau mau lebih kecil lagi */
  margin-bottom: 0.25rem;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 10px;
}

/* Tombol aksi – konsisten & modern */
.action-btn {
  border: none;
  background: none;
  cursor: pointer;
  padding: 6px;
  font-size: 1.15rem;
  transition: transform 0.15s ease, opacity 0.2s ease;
}

.action-btn:hover {
  opacity: 0.85;
  transform: scale(1.15);
}

/* Warna ikon sesuai kategori */
.action-btn.edit {
  color: #f39c12;
}

.action-btn.view {
  color: #3498db;
}

.action-btn.delete {
  color: #e74c3c;
}

.action-btn.delete.disabled {
  color: #ccc;
  cursor: not-allowed;
  transform: none;
  opacity: 0.5;
}

</style>