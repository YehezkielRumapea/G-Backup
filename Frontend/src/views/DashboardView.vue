<script setup>
import { ref } from 'vue';

// --- KONTROL MODAL ---
//-- Update BackupView.vue - 20 Okt 2025 -//
const showBackupModal = ref(false);
const showRestoreModal = ref(false);

// --- DATA FORM UNTUK BACKUP ---
const backupForm = ref({
  jobName: '',
  target: 'file', // 'file' atau 'database'
  dbName: '',
  dbUser: '',
  dbPass: '',
  cloudPath: '/backup-folder/data',
  backupMode: 'auto', // 'auto' atau 'onetime'
  scheduleType: 'hourly', // 'hourly' atau 'daily'
  scheduleValue: 4,
  scheduleDays: [],
  backupOption: 'Overwrite',
  encrypt: false,
});

// --- DATA FORM UNTUK RESTORE ---
const restoreForm = ref({
  jobName: '',
  target: 'file', // 'file' atau 'database'
  sourcePath: '/backup-folder/data/file.zip',
  destDbName: '',
  destDbUser: '',
  destDbPass: '',
  restoreOption: 'Overwrite'
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
            <div class="modal-header">
                <h2>Backup Config</h2>
                <button @click="closeModal" class="close-button">&times;</button>
            </div>
            <form @submit.prevent="saveBackupJob">
                <div class="form-group">
                    <label for="jobName">Job Name*</label>
                    <input type="text" id="jobName" v-model="backupForm.jobName" required>
                </div>

                <div class="form-group">
                    <label for="backupTarget">Select Target</label>
                    <select id="backupTarget" v-model="backupForm.target">
                        <option value="file">File</option>
                        <option value="database">Database</option>
                    </select>
                </div>

                <!-- Opsi untuk Target File -->
                <div v-if="backupForm.target === 'file'" class="form-group">
                    <label for="fileInput">File</label>
                    <input type="file" id="fileInput" class="file-input">
                </div>

                <!-- Opsi untuk Target Database -->
                <div v-if="backupForm.target === 'database'" class="db-group">
                    <div class="form-group">
                        <label for="dbName">DB Name</label>
                        <input type="text" id="dbName" v-model="backupForm.dbName">
                    </div>
                    <div class="form-group">
                        <label for="dbUser">DB Username</label>
                        <input type="text" id="dbUser" v-model="backupForm.dbUser">
                    </div>
                     <div class="form-group">
                        <label for="dbPass">DB Password</label>
                        <input type="password" id="dbPass" v-model="backupForm.dbPass">
                    </div>
                </div>
                
                <div class="form-group">
                    <label for="cloudPath">Cloud Path</label>
                    <input type="text" id="cloudPath" v-model="backupForm.cloudPath" placeholder="e.g., /backup-folder/data">
                </div>

                <div class="form-group">
                    <label>Backup Mode</label>
                     <div class="radio-group">
                        <label><input type="radio" v-model="backupForm.backupMode" value="auto"> Auto Backup</label>
                        <label><input type="radio" v-model="backupForm.backupMode" value="onetime"> One Time Backup</label>
                    </div>
                </div>
                
                <!-- Opsi Skema Backup (hanya muncul jika Auto) -->
                <div v-if="backupForm.backupMode === 'auto'" class="backup-scheme">
                    <h3>Backup Scheme</h3>
                     <div class="radio-group">
                        <label><input type="radio" v-model="backupForm.scheduleType" value="hourly"> Hourly</label>
                        <label><input type="radio" v-model="backupForm.scheduleType" value="daily"> Daily</label>
                    </div>
                    <div v-if="backupForm.scheduleType === 'hourly'" class="schedule-options">
                        <span>Run every</span>
                        <input type="number" v-model="backupForm.scheduleValue" class="hour-input">
                        <span>hour(s)</span>
                    </div>
                </div>

                <div class="form-group-inline">
                    <div class="form-group">
                        <label>Backup Options</label>
                        <select v-model="backupForm.backupOption">
                            <option>Overwrite</option>
                            <option>Append</option>
                        </select>
                    </div>
                    <div class="checkbox-group">
                         <label><input type="checkbox" v-model="backupForm.encrypt"> Encrypt</label>
                    </div>
                </div>

                <div class="modal-actions">
                    <button type="submit" class="btn-primary full-width">Start Job</button>
                </div>
            </form>
        </div>
    </div>
    
    <!-- MODAL UNTUK RESTORE JOB -->
     <div v-if="showRestoreModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content">
             <div class="modal-header">
                <h2>Restore Config</h2>
                <button @click="closeModal" class="close-button">&times;</button>
            </div>
            <form @submit.prevent="startRestore">
                 <div class="form-group">
                    <label for="restoreJobName">Job Name*</label>
                    <select id="restoreJobName" v-model="restoreForm.jobName" required>
                        <option disabled value="">-- Select a backup job to restore --</option>
                         <option v-for="job in backupJobs" :key="job" :value="job">{{ job }}</option>
                    </select>
                </div>
                 <div class="form-group">
                    <label for="restoreTarget">Select Target</label>
                    <select id="restoreTarget" v-model="restoreForm.target">
                        <option value="file">File</option>
                        <option value="database">Database</option>
                    </select>
                </div>

                <!-- Opsi untuk Restore File -->
                 <div v-if="restoreForm.target === 'file'" class="form-group">
                    <label for="restoreSourcePath">Restore From Path</label>
                    <input type="text" id="restoreSourcePath" v-model="restoreForm.sourcePath" placeholder="e.g., /backup-folder/data/file.zip">
                </div>

                <!-- Opsi untuk Restore Database -->
                <div v-if="restoreForm.target === 'database'" class="db-group">
                     <div class="form-group">
                        <label for="destDbName">Destination DB Name</label>
                        <input type="text" id="destDbName" v-model="restoreForm.destDbName">
                    </div>
                    <div class="form-group">
                        <label for="destDbUser">DB Username</label>
                        <input type="text" id="destDbUser" v-model="restoreForm.destDbUser">
                    </div>
                     <div class="form-group">
                        <label for="destDbPass">DB Password</label>
                        <input type="password" id="destDbPass" v-model="restoreForm.destDbPass">
                    </div>
                </div>

                <div class="form-group">
                    <label>Restore Options</label>
                    <select v-model="restoreForm.restoreOption">
                        <option>Overwrite</option>
                        <option>Append</option>
                    </select>
                </div>

                <div class="modal-actions">
                    <button type="submit" class="btn-success full-width">Start Job</button>
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
.action-btn { flex-grow: 1; padding: 60px 20px; border: none; border-radius: 8px; cursor: pointer; font-size: 1.5rem; font-weight: 600; color: #fff; transition: opacity 0.3s; }
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

/* --- STYLE BARU UNTUK MODAL --- */
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
  padding: 20px;
}

.modal-content {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.3);
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid #eee;
}

.modal-header h2 {
    margin: 0;
    font-size: 1.25rem;
}

.close-button {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #888;
}

.modal-content form {
    padding: 20px;
    overflow-y: auto;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.form-group input[type="text"],
.form-group input[type="password"],
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}

.file-input {
    border: 1px solid #ccc;
    border-radius: 4px;
    padding: 5px;
}

.db-group {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
}

.db-group .form-group:last-child {
    grid-column: span 2;
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

.backup-scheme {
    border: 1px solid #eee;
    background-color: #f9f9f9;
    padding: 15px;
    border-radius: 6px;
    margin-top: 10px;
}

.backup-scheme h3 {
    margin: 0 0 10px 0;
    font-size: 1rem;
    font-weight: 600;
}

.schedule-options {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-top: 10px;
}
.hour-input {
    width: 60px;
    text-align: center;
    padding: 5px;
    border: 1px solid #ccc;
    border-radius: 4px;
}

.form-group-inline {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    gap: 20px;
}
.checkbox-group {
    padding-bottom: 10px;
}
.checkbox-group label {
    display: flex;
    align-items: center;
    gap: 5px;
}

.modal-actions {
  margin-top: 20px;
}

.btn-primary, .btn-secondary, .btn-success {
    padding: 12px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-weight: 500;
    font-size: 1rem;
}
.btn-primary { background-color: #007bff; color: white; }
.btn-success { background-color: #28a745; color: white; }
.btn-secondary { background-color: #6c757d; color: white; }

.full-width {
    width: 100%;
}

</style>

