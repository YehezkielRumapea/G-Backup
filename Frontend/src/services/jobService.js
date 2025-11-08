import apiClient from './api' // Import instance Axios (sudah ada interceptor JWT)

export default {
  /**
   * Mengirim konfigurasi Job Backup baru (Auto atau Manual) ke backend.
   * Endpoint: POST /api/v1/jobs/new
   * @param {object} backupConfig - DTO (JSON) dari form Backup Config.
   */
  async createBackupJob(backupConfig) {
    try {
      // Panggil endpoint /jobs/new
      const response = await apiClient.post('/jobs/new', backupConfig)
      // Mengembalikan pesan sukses (misalnya, "Job diterima")
      return response.data
    } catch (error) {
      console.error("Error creating backup job:", error)
      throw error // Lempar error ke komponen Vue untuk ditampilkan
    }
  },

  /**
   * Mengirim konfigurasi Job Restore baru ke backend.
   * Endpoint: POST /api/v1/jobs/restore
   * @param {object} restoreConfig - DTO (JSON) dari form Restore Config.
   */
  async createRestoreJob(restoreConfig) {
    try {
      // Panggil endpoint /jobs/restore
      const response = await apiClient.post('/jobs/restore', restoreConfig)
      return response.data
    } catch (error) {
      console.error("Error creating restore job:", error)
      throw error
    }
  },

  /**
   * Memicu Job (Auto atau Manual) yang sudah ada di DB untuk berjalan SEKARANG.
   * Endpoint: POST /api/v1/jobs/trigger/:id
   * @param {number} jobId - ID dari Job yang akan dijalankan.
   */
  async triggerManualJob(jobId) {
    try {
      // Panggil endpoint /jobs/trigger/[id]
      const response = await apiClient.post(`/jobs/trigger/${jobId}`)
      return response.data
    } catch (error) {
      console.error(`Error triggering job ${jobId}:`, error)
      throw error
    }
  },

  /**
   * Mengambil pratinjau script lengkap untuk Job ID tertentu.
   * Endpoint: GET /api/v1/jobs/script/:id
   * @param {number} jobId - ID dari Job.
   * @returns {Promise<object>} Object berisi "script_preview".
   */
  async getJobScript(jobId) {
    try {
      const response = await apiClient.get(`/jobs/script/${jobId}`)
      return response.data
    } catch (error) {
      console.error(`Error fetching job script ${jobId}:`, error)
      throw error
    }
  },

  /**
   * Menjelajahi file di GDrive (Untuk UI Restore).
   * Endpoint: GET /api/v1/browser/list
   * @param {string} remoteName - Nama remote GDrive.
   * @param {string} path - Path folder yang akan dijelajahi.
   * @returns {Promise<Array>} Daftar file/folder.
   */
  async browseFiles(remoteName, path) {
    try {
      const response = await apiClient.get('/browser/list', {
        params: {
          remote: remoteName,
          path: path,
        },
      })
      return response.data
    } catch (error) {
      console.error(`Error browsing files at ${remoteName}:${path}:`, error)
      throw error
    }
  },
}