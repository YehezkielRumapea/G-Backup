<script setup>
import { ref } from 'vue';

// --- KONTROL MODAL ---
const showBackupModal = ref(false);
const showRestoreModal = ref(false);

// --- DATA FORM UNTUK BACKUP ---
const backupForm = ref({
  jobName: '',
  remoteType: 'Cloud',
  selectedRemote: '',
  cloudPath: '/backups/website/',
  sourcePath: 'C:/projects/my-website/',
  backupMode: 'Append',
});

// --- DATA FORM UNTUK RESTORE ---
const restoreForm = ref({
  selectedJob: '',
  sourcePath: '/backups/website/2025-10-07',
  restorePath: 'C:/restore/my-website/',
});


// --- DATA SIMULASI UNTUK DROPDOWN ---
const remotes = ref(['Remote1', 'Azure-West', 'S3-Primary-Bucket']);
const backupJobs = ref(['Backup1', 'Sync-Images', 'DB_Archive']);

// Data simulasi untuk Log Activity
const logActivities = ref([
  { id: 1, time: '10:30 AM', message: 'Backup_Daily_DB completed successfully.'},
  { id: 2, time: '10:25 AM', message: 'Failed to connect to Azure-West.'},
  { id: 3, time: '10:15 AM', message: 'Job Sync_Media_Weekly started.'},
  { id: 4, time: '09:50 AM', message: 'User admin logged in.'},
]);

// --- FUNGSI ---
const openBackupModal = () => {
  showBackupModal.value = true;
};

const openRestoreModal = () => {
  showRestoreModal.value = true;
};

const closeModal = () => {
  showBackupModal.value = false;
  showRestoreModal.value = false;
};

const saveBackupJob = () => {
  // Logika untuk menyimpan backup job
  alert(`Menyimpan Backup Job:\n${JSON.stringify(backupForm.value, null, 2)}`);
  closeModal();
};

const startRestore = () => {
  // Logika untuk memulai restore
  alert(`Memulai Restore Job:\n${JSON.stringify(restoreForm.value, null, 2)}`);
  closeModal();
};
</script>

<template>
    <div class="dashboard-container">
        <!-- Header Halaman -->
        <header class="main-header">
            <h1>Dashboard</h1>
        </header>

        <!-- Konten Utama Dashboard -->
        <div class="dashboard-grid">
            
            <!-- Kolom Kiri -->
            <div class="left-column">
                <div class="top-stats">
                    <div class="stat-card">
                        <h3>Total Remote</h3>
                        <p>5</p>
                    </div>
                    <div class="stat-card">
                        <h3>Total Job</h3>
                        <p>12</p>
                    </div>
                </div>

                <div class="card quick-actions">
                    <h3>Quick Action</h3>
                    <div class="action-buttons">
                        <button @click="openBackupModal" class="action-btn backup-btn">Backup</button>
                        <button @click="openRestoreModal" class="action-btn restore-btn">Restore</button>
                    </div>
                </div>

                <div class="card next-job">
                    <h3>Next Job</h3>
                    <p><strong>Backup_Finance_Docs</strong></p>
                    <p>Scheduled for: <strong>Today, 11:00 PM</strong></p>
                </div>
            </div>

            <!-- Kolom Kanan (Log Activity) -->
            <div class="right-column">
                <div class="card log-activity">
                    <h3>Log Activity</h3>
                    <ul>
                        <li v-for="log in logActivities" :key="log.id">
                            <span class="log-time">{{ log.time }}</span>
                            <span class="log-message">{{ log.message }}</span>
                        </li>
                    </ul>
                </div>
            </div>

        </div>
    </div>

    <!-- MODAL UNTUK BACKUP JOB -->
    <div v-if="showBackupModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content">
            <h2>Backup Job / Config</h2>
            <form @submit.prevent="saveBackupJob">
                <div class="form-group">
                    <label for="jobName">Job Name</label>
                    <input type="text" id="jobName" v-model="backupForm.jobName" placeholder="e.g., Daily_Website_Backup" required>
                </div>

                <div class="form-group">
                    <label>Select Remote</label>
                    <div class="radio-group">
                        <label><input type="radio" v-model="backupForm.remoteType" value="Cloud"> Cloud</label>
                        <label><input type="radio" v-model="backupForm.remoteType" value="Local"> Local</label>
                    </div>
                </div>

                <div class="form-group">
                    <select v-model="backupForm.selectedRemote" required>
                        <option disabled value="">-- Please select a remote --</option>
                        <option v-for="remote in remotes" :key="remote" :value="remote">{{ remote }}</option>
                    </select>
                </div>
                
                <div class="form-group">
                    <label for="cloudPath">Cloud Path</label>
                    <input type="text" id="cloudPath" v-model="backupForm.cloudPath">
                </div>
                
                <div class="form-group">
                    <label for="sourcePath">Source Path</label>
                    <input type="text" id="sourcePath" v-model="backupForm.sourcePath">
                </div>

                <div class="form-group">
                    <label>Backup Mode</label>
                     <div class="radio-group">
                        <label><input type="radio" v-model="backupForm.backupMode" value="Append"> Append</label>
                        <label><input type="radio" v-model="backupForm.backupMode" value="Mirror"> Mirror</label>
                        <label><input type="radio" v-model="backupForm.backupMode" value="Update"> Update</label>
                    </div>
                </div>

                <div class="modal-actions">
                    <button type="button" @click="closeModal" class="btn-secondary">Cancel</button>
                    <button type="submit" class="btn-primary">Save Job</button>
                </div>
            </form>
        </div>
    </div>
    
    <!-- MODAL UNTUK RESTORE JOB -->
     <div v-if="showRestoreModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content">
            <h2>Restore Job / Config</h2>
            <form @submit.prevent="startRestore">
                 <div class="form-group">
                    <label for="restoreJobName">Job Name</label>
                    <select id="restoreJobName" v-model="restoreForm.selectedJob" required>
                        <option disabled value="">-- Select a backup job to restore --</option>
                         <option v-for="job in backupJobs" :key="job" :value="job">{{ job }}</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="restoreSourcePath">Source Path (from remote)</label>
                    <input type="text" id="restoreSourcePath" v-model="restoreForm.sourcePath">
                </div>
                 <div class="form-group">
                    <label for="restorePath">Restore Path (to local)</label>
                    <input type="text" id="restorePath" v-model="restoreForm.restorePath">
                </div>

                <div class="modal-actions">
                    <button type="button" @click="closeModal" class="btn-secondary">Cancel</button>
                    <button type="submit" class="btn-success">Start Restore</button>
                </div>
            </form>
        </div>
    </div>

</template>

<style scoped>
/* --- STYLE DASHBOARD YANG SUDAH ADA --- */
.main-header { margin-bottom: 20px; }
.main-header h1 { font-size: 2rem; font-weight: 700; color: #333; }
.dashboard-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 20px; }
.left-column, .right-column { display: flex; flex-direction: column; gap: 20px; }
.card { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); }
.card h3 { margin-bottom: 15px; font-size: 1.1rem; font-weight: 600; border-bottom: 1px solid #eee; padding-bottom: 10px; }
.top-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.stat-card { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); }
.stat-card h3 { font-size: 1rem; color: #6c757d; margin-bottom: 10px; }
.stat-card p { font-size: 2.2rem; font-weight: 700; color: var(--primary-color); }
.quick-actions .action-buttons { display: flex; gap: 15px; }
.action-btn { flex-grow: 1; padding: 25px; border: none; border-radius: 6px; cursor: pointer; font-size: 1.2rem; font-weight: 600; color: #fff; transition: opacity 0.3s; }
.action-btn:hover { opacity: 0.9; }
.backup-btn { background-color: #007bff; }
.restore-btn { background-color: #28a745; }
.next-job p { font-size: 1rem; color: #333; line-height: 1.6; }
.log-activity { flex-grow: 1; }
.log-activity ul { list-style: none; padding: 0; }
.log-activity li { padding: 8px 0; border-bottom: 1px solid #f0f0f0; display: flex; gap: 15px; font-size: 0.9rem; }
.log-activity li:last-child { border-bottom: none; }
.log-time { color: #6c757d; flex-shrink: 0; }
.log-message { color: #343a40; }

@media (max-width: 992px) { .dashboard-grid { grid-template-columns: 1fr; } }

/* --- STYLE UNTUK MODAL --- */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.3);
  width: 90%;
  max-width: 500px;
}

.modal-content h2 {
  margin-top: 0;
  margin-bottom: 25px;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.form-group input[type="text"],
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}

.radio-group {
    display: flex;
    gap: 20px;
    align-items: center;
    margin-top: 10px;
}
.radio-group label {
    font-weight: 400;
    display: flex;
    align-items: center;
    gap: 5px;
}


.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 30px;
}

.btn-primary, .btn-secondary, .btn-success {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-weight: 500;
}
.btn-primary { background-color: #007bff; color: white; }
.btn-success { background-color: #28a745; color: white; }
.btn-secondary { background-color: #6c757d; color: white; }

</style>

