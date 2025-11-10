<template>
    <tr>
        <td><strong>{{ job.job_name }}</strong></td>
        
        <td>{{ job.type }}</td> 

        <td>{{ job.gdrive_target }}</td>

        <td>{{ job.mode }}</td>

        <td>{{ job.last_run || 'N/A' }}</td>

        <td>
            <span class="status" :class="job.status.toLowerCase()">
                {{ job.status }}
            </span>
        </td>
        
        <td>{{ job.next_run || 'N/A' }}</td>
        
        <td class="job-actions-col"> 
          <div class="actions">
            <button @click="$emit('trigger', job.id)" class="action-btn play" title="Run Now (Ad-hoc)">
              ▶
            </button>

            <button @click="$emit('view-script', job.id)" class="action-btn view" title="View Script">
              (i)
            </button>

            <!-- <button disabled class="action-btn delete disabled" title="Hapus dinonaktifkan untuk Job Terjadwal">
              🗑
            </button> -->
          </div>
        </td>
    </tr>
</template>

<script setup>
import jobService from '@/services/jobService'; // Import service untuk aksi
import { computed } from 'vue';

// 1. Terima data 'job' dari parent
const props = defineProps({
  job: {
    type: Object,
    required: true
  }
});

// 2. Definisikan event yang akan dikirim ke parent
const emit = defineEmits(['trigger', 'view-script']);

// Fungsi placeholder untuk menghapus (tetapi dinonaktifkan di template)
function handleDelete() {
    alert("Penghapusan Job Terjadwal harus melalui menu Edit.");
}
</script>
<style scoped>
/* --- STYLING BARIS JOB --- */
.status {
    padding: 4px 8px;
    border-radius: 12px;
    font-weight: bold;
    font-size: 0.8rem;
    color: white;
}
/* Tambahkan semua status queue yang Anda miliki */
.status.pending { background-color: #aaa; }
.status.running { background-color: #3498db; }
.status.completed, .status.success { background-color: #2ecc71; }
.status.failed, .status.fail, .status.fail_pre_script, .status.fail_rclone { background-color: #e74c3c; }

/* Styling kolom aksi */
.job-runs-col {
    width: 150px; /* Lebar lebih besar untuk tombol */
}
.actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}
.action-btn {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 5px;
  transition: opacity 0.2s;
}
.action-btn:hover {
    opacity: 0.8;
}

/* Warna tombol aksi */
.action-btn.play { color: #f39c12; }
.action-btn.view { color: #3499db; font-weight: bold; }
.action-btn.delete { color: #e74c3c; }
.action-btn.delete.disabled {
    color: #ccc;
    cursor: not-allowed;
}
</style>