<script setup>
import { ref, onMounted, computed } from 'vue';
import apiClient from '@/api'; // Gunakan jembatan API yang sudah kita buat

// --- STATE UNTUK KONTROL UI ---
const showBackupModal = ref(false);
const showRestoreModal = ref(false);

// --- STATE UNTUK DATA DARI API ---
const remotes = ref([]);
const logs = ref([]);
const jobs = ref([]);
const isLoading = ref(true);
const errorMsg = ref('');

// --- STATE UNTUK FORM BACKUP ---
// Strukturnya disesuaikan agar cocok dengan struct `BackupRequest` di backend Go (json tag)
const backupForm = ref({
    job_name: '',
    job_type: 'FILE', // 'FILE' atau 'DB'
    source_path: '', // Untuk FILE: path di VM1, Untuk DB: nama database di VM1
    
    // Detail koneksi SSH ke server aplikasi (VM1)
    ssh_host: '192.168.56.142', // <-- Sesuaikan jika IP VM1 Anda berbeda
    ssh_user: 'idven',          // <-- Sesuaikan jika user di VM1 Anda berbeda
    
    // Detail Remote Rclone
    remote_name: '',
    destination_path: '/gbackup_data/', // Path tujuan di Google Drive
    
    // Penjadwalan (kosongkan untuk manual)
    schedule_cron: '',
    
    // Opsi Tambahan
    is_encrypted: false,

    // Kredensial DB (jika job_type = 'DB')
    db_name: '', // Hanya untuk UI, akan disalin ke source_path
    db_user: '',
    db_pass: '',
});

// --- STATE UNTUK FORM RESTORE ---
// Strukturnya disesuaikan dengan struct `RestoreRequest` di backend Go
const restoreForm = ref({
    job_name: 'Manual Restore',
    remote_name: '',
    source_path: '', // Path file backup di Google Drive
    restore_path: '', // Path tujuan restore di server
    job_type: 'FILE', // 'FILE' atau 'DB'
    destDbName: '', // Hanya untuk UI
    db_user: '',
    db_pass: ''
});

// --- FUNGSI UNTUK MENGAMBIL DATA AWAL ---
const fetchData = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    try {
        // Panggil beberapa endpoint sekaligus untuk efisiensi
        const [remotesRes, logsRes, jobsRes] = await Promise.all([
            apiClient.get('/remotes/status'),
            apiClient.get('/logs'), // Pastikan endpoint ini ada di Go
            apiClient.get('/jobs/scheduled')
        ]);
        
        remotes.value = remotesRes.data || [];
        logs.value = (logsRes.data || []).slice(0, 5); // Ambil 5 log terbaru
        jobs.value = jobsRes.data || [];

        // Set remote pertama sebagai default di form jika belum ada
        if (remotes.value.length > 0 && !backupForm.value.remote_name) {
            backupForm.value.remote_name = remotes.value[0].RemoteName;
        }
        if (remotes.value.length > 0 && !restoreForm.value.remote_name) {
            restoreForm.value.remote_name = remotes.value[0].RemoteName;
        }

    } catch (error) {
        console.error("Gagal memuat data dashboard:", error);
         if (error.response?.status !== 401) {
            errorMsg.value = "Gagal memuat data dari server. Pastikan backend berjalan dan terhubung.";
         }
    } finally {
        isLoading.value = false;
    }
};

// Panggil fetchData() saat komponen pertama kali dimuat
onMounted(fetchData);

// --- FUNGSI UNTUK MENGIRIM DATA FORM KE BACKEND ---
const saveBackupJob = async () => {
    // Siapkan payload yang akan dikirim, sesuai dengan `BackupRequest` di Go
    let payload = { ...backupForm.value };

    // Logika "satu nalar": Untuk tipe DB, `source_path` di backend adalah nama database
    if (payload.job_type === 'DB') {
        payload.source_path = payload.db_name;
    }

    // Hapus field yang hanya untuk UI dan tidak perlu dikirim ke backend
    delete payload.db_name; 
    
    console.log("Mengirim payload backup:", payload); // Untuk debugging

    try {
        const response = await apiClient.post('/jobs/new', payload);
        alert(response.data.message || "Job backup berhasil dibuat!");
        closeModal();
        fetchData(); // Muat ulang data di dashboard setelah job dibuat
    } catch (error) {
        console.error("Error saat menyimpan job backup:", error.response?.data || error.message);
        alert("Gagal membuat job backup: " + (error.response?.data?.error || "Error tidak diketahui"));
    }
};

const startRestore = async () => {
    // Siapkan payload yang akan dikirim, sesuai dengan `RestoreRequest` di Go
    let payload = { ...restoreForm.value };

    // Logika "satu nalar": Untuk tipe DB, `restore_path` di backend adalah nama DB tujuan
    if (payload.job_type === 'DB') {
        payload.restore_path = payload.destDbName;
    }

    // Hapus field yang hanya untuk UI
    delete payload.destDbName;

    console.log("Mengirim payload restore:", payload); // Untuk debugging

    try {
        const response = await apiClient.post('/restore', payload);
        alert(response.data.message || "Job restore berhasil dimulai!");
        closeModal();
        fetchData(); // Muat ulang log
    } catch (error) {
         console.error("Error saat memulai restore:", error.response?.data || error.message);
        alert("Gagal memulai restore: " + (error.response?.data?.error || "Error tidak diketahui"));
    }
};

// --- FUNGSI UTILITAS ---
const closeModal = () => {
  showBackupModal.value = false;
  showRestoreModal.value = false;
};
const formatTimestamp = (timestamp) => {
    if (!timestamp) return '-';
    try { return new Date(timestamp).toLocaleString(); } catch(e) { return timestamp; }
};

// --- COMPUTED PROPERTIES UNTUK STATISTIK DINAMIS ---
const totalRemotes = computed(() => remotes.value.length);
const totalJobs = computed(() => jobs.value.length);
const nextJob = computed(() => {
    // Cari job dengan `nextRun` terdekat di masa depan
    const futureJobs = jobs.value
        .map(j => ({ ...j, nextRunDate: new Date(j.nextRun) }))
        .filter(j => j.nextRunDate > new Date())
        .sort((a, b) => a.nextRunDate - b.nextRunDate);
    return futureJobs.length > 0 ? futureJobs[0] : null;
});
</script>

<template>
    <div class="dashboard-container">
        <header class="main-header">
            <h1>Dashboard</h1>
        </header>

        <div v-if="isLoading" class="text-center p-8">Memuat data dashboard...</div>
        <div v-else-if="errorMsg" class="error-message p-4">{{ errorMsg }}</div>

        <div v-else class="dashboard-grid">
            <div class="left-column">
                <div class="top-stats">
                    <div class="stat-card">
                        <h3>Total Remote</h3>
                        <p>{{ totalRemotes }}</p>
                    </div>
                    <div class="stat-card">
                        <h3>Total Job</h3>
                        <p>{{ totalJobs }}</p>
                    </div>
                </div>

                <div class="card quick-actions">
                    <h3>Quick Action</h3>
                    <div class="action-buttons">
                        <button @click="showBackupModal = true" class="action-btn backup-btn">Create Backup</button>
                        <button @click="showRestoreModal = true" class="action-btn restore-btn">Start Restore</button>
                    </div>
                </div>

                <div class="card next-job">
                    <h3>Next Scheduled Job</h3>
                    <div v-if="nextJob">
                        <p><strong>{{ nextJob.name }}</strong></p>
                        <p>Scheduled for: <strong>{{ formatTimestamp(nextJob.nextRun) }}</strong></p>
                    </div>
                     <div v-else>
                        <p>Tidak ada job terjadwal berikutnya.</p>
                    </div>
                </div>
            </div>

            <div class="right-column">
                <div class="card log-activity">
                    <h3>Recent Log Activity</h3>
                    <ul v-if="logs.length > 0">
                        <li v-for="log in logs" :key="log.ID">
                            <span class="log-time">{{ formatTimestamp(log.Timestamp) }}</span>
                            <span class="log-message">[{{ log.OperationType }} - {{ log.Status }}] {{ log.Message }}</span>
                        </li>
                    </ul>
                     <p v-else>Belum ada aktivitas log.</p>
                </div>
            </div>
        </div>

        <div v-if="showBackupModal" class="modal-overlay" @click.self="closeModal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>Backup Configuration</h2>
                    <button @click="closeModal" class="close-button">&times;</button>
                </div>
                <form @submit.prevent="saveBackupJob">
                    <div class="form-group">
                        <label for="backupJobName">Job Name*</label>
                        <input type="text" id="backupJobName" v-model="backupForm.job_name" required>
                    </div>

                    <div class="form-group">
                        <label for="backupJobType">Target Type</label>
                        <select id="backupJobType" v-model="backupForm.job_type">
                            <option value="FILE">File / Folder</option>
                            <option value="DB">Database (MariaDB/MySQL)</option>
                        </select>
                    </div>
                    
                    <div class="form-group">
                        <label for="sshHost">Application Server IP (VM1)</label>
                        <input type="text" id="sshHost" v-model="backupForm.ssh_host" required>
                    </div>
                    <div class="form-group">
                        <label for="sshUser">SSH User (on VM1)</label>
                        <input type="text" id="sshUser" v-model="backupForm.ssh_user" required>
                    </div>

                    <div v-if="backupForm.job_type === 'FILE'">
                        <div class="form-group">
                            <label for="backupSourcePathFile">Source Path (on VM1)</label>
                            <input type="text" id="backupSourcePathFile" v-model="backupForm.source_path" placeholder="/var/www/html">
                        </div>
                    </div>

                    <div v-if="backupForm.job_type === 'DB'" class="db-group">
                        <div class="form-group">
                            <label for="backupDbName">DB Name (on VM1)</label>
                            <input type="text" id="backupDbName" v-model="backupForm.db_name">
                        </div>
                        <div class="form-group">
                            <label for="backupDbUser">DB Username</label>
                            <input type="text" id="backupDbUser" v-model="backupForm.db_user">
                        </div>
                         <div class="form-group">
                            <label for="backupDbPass">DB Password</label>
                            <input type="password" id="backupDbPass" v-model="backupForm.db_pass">
                        </div>
                    </div>
                    
                    <div class="form-group">
                        <label for="backupRemoteName">Cloud Remote Destination</label>
                        <select id="backupRemoteName" v-model="backupForm.remote_name" required>
                            <option disabled value="">-- Select Remote --</option>
                            <option v-for="remote in remotes" :key="remote.RemoteName" :value="remote.RemoteName">
                                {{ remote.RemoteName }}
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="backupDestPath">Cloud Destination Path</label>
                        <input type="text" id="backupDestPath" v-model="backupForm.destination_path" placeholder="/folder-backup/">
                    </div>

                    <div class="modal-actions">
                         <button type="button" class="btn-secondary" @click="closeModal">Cancel</button>
                        <button type="submit" class="btn-primary">Start Manual Backup</button>
                    </div>
                </form>
            </div>
        </div>
        
        <div v-if="showRestoreModal" class="modal-overlay" @click.self="closeModal">
            <div class="modal-content">
                 <div class="modal-header">
                    <h2>Restore Configuration</h2>
                    <button @click="closeModal" class="close-button">&times;</button>
                </div>
                <form @submit.prevent="startRestore">
                     <div class="form-group">
                        <label for="restoreRemoteName">Cloud Remote Source</label>
                         <select id="restoreRemoteName" v-model="restoreForm.remote_name" required>
                             <option disabled value="">-- Select Remote --</option>
                             <option v-for="remote in remotes" :key="remote.RemoteName" :value="remote.RemoteName">
                                {{ remote.RemoteName }}
                            </option>
                        </select>
                    </div>
                     <div class="form-group">
                        <label for="restoreSourcePath">Source Path (on Cloud)</label>
                        <input type="text" id="restoreSourcePath" v-model="restoreForm.source_path" placeholder="/folder-backup/namafile.zip atau .sql">
                    </div>

                     <div class="form-group">
                        <label for="restoreJobType">Restore Target Type</label>
                        <select id="restoreJobType" v-model="restoreForm.job_type">
                            <option value="FILE">File / Folder</option>
                            <option value="DB">Database (MariaDB/MySQL)</option>
                        </select>
                    </div>

                    <div v-if="restoreForm.job_type === 'FILE'">
                        <div class="form-group">
                            <label for="restoreDestPathFile">Destination Path (on Server)</label>
                            <input type="text" id="restoreDestPathFile" v-model="restoreForm.restore_path" placeholder="/var/www/restored_files/">
                        </div>
                    </div>

                    <div v-if="restoreForm.job_type === 'DB'" class="db-group">
                         <div class="form-group">
                            <label for="restoreDestDbName">Destination DB Name</label>
                            <input type="text" id="restoreDestDbName" v-model="restoreForm.destDbName">
                        </div>
                        <div class="form-group">
                            <label for="restoreDbUser">DB Username (Server)</label>
                            <input type="text" id="restoreDbUser" v-model="restoreForm.db_user">
                        </div>
                         <div class="form-group">
                            <label for="restoreDbPass">DB Password (Server)</label>
                            <input type="password" id="restoreDbPass" v-model="restoreForm.db_pass">
                        </div>
                    </div>

                    <div class="modal-actions">
                        <button type="button" @click="closeModal" class="btn-secondary">Cancel</button>
                        <button type="submit" class="btn-success">Start Restore</button>
                    </div>
                </form>
            </div>
        </div>

    </div>
</template>

<style scoped>
/* --- STYLE DASHBOARD --- */
.main-header { margin-bottom: 20px; }
.main-header h1 { font-size: 2rem; font-weight: 700; color: #333; }
.dashboard-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 20px; }
.left-column, .right-column { display: flex; flex-direction: column; gap: 20px; }
.card { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); }
.card h3 { margin-bottom: 15px; font-size: 1.1rem; font-weight: 600; border-bottom: 1px solid #eee; padding-bottom: 10px; }
.top-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.stat-card { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); text-align: center;}
.stat-card h3 { font-size: 1rem; color: #6c757d; margin-bottom: 10px; }
.stat-card p { font-size: 2.2rem; font-weight: 700; color: var(--primary-color); }
.quick-actions .action-buttons { display: flex; gap: 15px; }
.action-btn { flex-grow: 1; padding: 30px 15px; border: none; border-radius: 8px; cursor: pointer; font-size: 1.2rem; font-weight: 600; color: #fff; text-align: center; }
.backup-btn { background-color: #0d6efd; }
.restore-btn { background-color: #198754; }
.next-job p { font-size: 1rem; color: #333; line-height: 1.6; }
.log-activity { flex-grow: 1; }
.log-activity ul { list-style: none; padding: 0; margin: 0; max-height: 300px; overflow-y: auto;}
.log-activity li { padding: 8px 0; border-bottom: 1px solid #f0f0f0; display: flex; gap: 10px; font-size: 0.9rem; flex-wrap: wrap;}
.log-activity li:last-child { border-bottom: none; }
.log-time { color: #6c757d; flex-shrink: 0; font-size: 0.8rem; }
.log-message { color: #343a40; word-break: break-word; }
@media (max-width: 992px) { .dashboard-grid { grid-template-columns: 1fr; } }
.error-message { color: #dc3545; }
.text-center { text-align: center; }

/* --- STYLE MODAL --- */
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); display: flex; justify-content: center; align-items: center; z-index: 1000; padding: 20px; }
.modal-content { background-color: #fff; border-radius: 8px; box-shadow: 0 5px 15px rgba(0,0,0,0.3); width: 90%; max-width: 550px; max-height: 90vh; display: flex; flex-direction: column; }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 15px 20px; border-bottom: 1px solid #eee; }
.modal-header h2 { margin: 0; font-size: 1.25rem; }
.close-button { background: none; border: none; font-size: 1.7rem; cursor: pointer; color: #888; line-height: 1; }
.modal-content form { padding: 20px; overflow-y: auto; display: flex; flex-direction: column; gap: 15px; }
.form-group label { display: block; margin-bottom: 5px; font-weight: 500; font-size: 0.9rem; }
.form-group input, .form-group select { width: 100%; padding: 10px; border: 1px solid #ccc; border-radius: 4px; font-size: 1rem; box-sizing: border-box; }
.db-group { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 15px; }
.modal-actions { margin-top: 20px; display: flex; justify-content: flex-end; gap: 10px; padding: 15px 20px; border-top: 1px solid #eee; }
.btn-primary, .btn-secondary, .btn-success { padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; font-weight: 500; font-size: 0.9rem; }
.btn-primary { background-color: #0d6efd; color: white; }
.btn-success { background-color: #198754; color: white; }
.btn-secondary { background-color: #6c757d; color: white; }
</style>