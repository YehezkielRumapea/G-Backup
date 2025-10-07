<script setup>
import { ref, reactive } from 'vue';

// --- STATE UNTUK FORM ---
const backupJobForm = reactive({
  jobName: '',
  remoteType: 'Cloud',
  selectedRemote: '',
  cloudPath: '',
  backupMode: 'Append',
  sourcePath: '',
});

const restoreJobForm = reactive({
  jobName: '',
  remoteType: 'Cloud',
  restorePath: '',
  sourcePath: '',
});

// --- FUNGSI SUBMIT (PLACEHOLDER) ---
const handleSaveBackup = () => {
  alert('Backup configuration saved!\n' + JSON.stringify(backupJobForm, null, 2));
};

const handleStartRestore = () => {
  alert('Restore process started!\n' + JSON.stringify(restoreJobForm, null, 2));
};
</script>

<template>
  <div class="config-content">
    <header class="main-header">
      <h1>Configuration</h1>
      <p>Setup and manage your backup and restore jobs.</p>
    </header>

    <div class="config-grid">
      <!-- KARTU UNTUK BACKUP JOB -->
      <div class="card">
        <h2>Backup Job / Config</h2>
        <form @submit.prevent="handleSaveBackup" class="config-form">
          <div class="form-group">
            <label for="job-name">Job Name</label>
            <input id="job-name" v-model="backupJobForm.jobName" type="text" placeholder="e.g., Daily_Website_Backup" required>
          </div>

          <div class="form-group">
            <label>Select Remote</label>
            <div class="radio-group">
              <label><input type="radio" v-model="backupJobForm.remoteType" value="Cloud"> Cloud</label>
              <label><input type="radio" v-model="backupJobForm.remoteType" value="Local"> Local</label>
            </div>
            <select v-model="backupJobForm.selectedRemote" required>
              <option disabled value="">-- Please select a remote --</option>
              <option>Google Drive Utama</option>
              <option>S3-Archive</option>
              <option>Local-NAS-01</option>
            </select>
          </div>

          <div class="form-group">
            <label for="cloud-path">Cloud Path</label>
            <input id="cloud-path" v-model="backupJobForm.cloudPath" type="text" placeholder="/backups/website/">
          </div>
          
          <div class="form-group">
            <label for="source-path-backup">Source Path</label>
            <input id="source-path-backup" v-model="backupJobForm.sourcePath" type="text" placeholder="C:/projects/my-website/">
          </div>

          <div class="form-group">
            <label>Backup Mode</label>
             <div class="radio-group">
                <label><input type="radio" v-model="backupJobForm.backupMode" value="Append"> Append</label>
                <label><input type="radio" v-model="backupJobForm.backupMode" value="Mirror"> Mirror</label>
                <label><input type="radio" v-model="backupJobForm.backupMode" value="Update"> Update</label>
            </div>
          </div>
          
          <div class="form-actions">
            <button type="submit" class="btn btn-primary">Save Job</button>
          </div>
        </form>
      </div>

      <!-- KARTU UNTUK RESTORE JOB -->
      <div class="card">
        <h2>Restore Job / Config</h2>
        <form @submit.prevent="handleStartRestore" class="config-form">
           <div class="form-group">
            <label for="restore-job-name">Job Name</label>
            <select id="restore-job-name" v-model="restoreJobForm.jobName" required>
                <option disabled value="">-- Select a backup job to restore --</option>
                <option>Daily_Website_Backup</option>
                <option>Archive_Docs_Monthly</option>
            </select>
          </div>

           <div class="form-group">
            <label for="source-path-restore">Source Path (from remote)</label>
            <input id="source-path-restore" v-model="restoreJobForm.sourcePath" type="text" placeholder="/backups/website/2025-10-07">
          </div>

          <div class="form-group">
            <label for="restore-path">Restore Path (to local)</label>
            <input id="restore-path" v-model="restoreJobForm.restorePath" type="text" placeholder="C:/restore/my-website/">
          </div>
          
          <div class="form-actions">
            <button type="submit" class="btn btn-success">Start Restore</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.main-header {
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}
.main-header h1 { font-size: 1.8rem; }
.main-header p { color: #6c757d; }

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.card {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.05);
  border: 1px solid var(--border-color);
  padding: 25px;
}

.card h2 {
  font-size: 1.4rem;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 10px;
}

.config-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
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

.form-group input[type="text"],
.form-group select {
  padding: 10px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  font-size: 0.9rem;
  width: 100%;
}

.radio-group {
  display: flex;
  gap: 20px;
  margin-bottom: 10px;
}

.radio-group label {
  font-weight: normal;
  display: flex;
  align-items: center;
  gap: 5px;
}

.form-actions {
  margin-top: 15px;
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
}
.btn-primary { background-color: var(--primary-color); color: white; }
.btn-primary:hover { background-color: var(--primary-hover); }
.btn-success { background-color: #28a745; color: white; }
.btn-success:hover { background-color: #218838; }
</style>
