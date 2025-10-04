-- =========================================================
-- 1. Tabel users (Otentikasi dan Manajemen Akses)
-- =========================================================
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- =========================================================
-- 2. Tabel scheduled_jobs (Konfigurasi Job Otomatis/Terjadwal)
-- =========================================================
CREATE TABLE scheduled_jobs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    job_name VARCHAR(100) NOT NULL,
    
    -- Konfigurasi Operasi
    job_type ENUM('FILE', 'DB') NOT NULL, 
    rclone_mode ENUM('COPY', 'SYNC', 'APPEND', 'OVERWRITE') NOT NULL,
    
    -- Path dan Remote
    source_path VARCHAR(255) NOT NULL,       
    remote_name VARCHAR(100) NOT NULL,       
    
    -- Penjadwalan dan Enkripsi
    schedule_cron VARCHAR(50) NOT NULL,      
    priority INT DEFAULT 5,                  
    is_encrypted BOOLEAN DEFAULT FALSE,
    encrypt_key VARCHAR(255) NULL,           -- Kunci Enkripsi Runtime Rclone
    
    -- Status dan Kontrol
    is_active BOOLEAN DEFAULT TRUE,          
    last_run_at DATETIME NULL,               
    
    -- Foreign Key ke tabel users
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- =========================================================
-- 3. Tabel logs (Riwayat Eksekusi: Otomatis, Manual, dan Restore)
-- =========================================================
CREATE TABLE logs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Key ke scheduled_jobs.id (Bisa NULL untuk Job Manual/Restore)
    job_id INT UNSIGNED NULL, 
    
    operation_type ENUM('BACKUP', 'RESTORE', 'MANUAL_BACKUP') NOT NULL,
    status ENUM('SUCCESS', 'FAIL', 'ERROR') NOT NULL,
    
    -- Config Snapshot (Menyimpan konfigurasi Job Manual dalam format JSON)
    config_snapshot JSON NULL, 
    
    message TEXT,                               
    checksum VARCHAR(255) NULL,                 
    duration_sec INT NULL,                      
    
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Relasi: Jika scheduled_job dihapus, log tetap ada dan job_id menjadi NULL
    FOREIGN KEY (job_id) REFERENCES scheduled_jobs(id) ON DELETE SET NULL
);

-- =========================================================
-- 4. Tabel monitoring (Status Cloud Storage)
-- =========================================================
CREATE TABLE monitoring (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    remote_name VARCHAR(100) UNIQUE NOT NULL,   
    status_connect ENUM('CONNECTED', 'DISCONNECTED') DEFAULT 'DISCONNECTED',
    
    -- Penyimpanan dalam satuan GB (FLOAT untuk presisi desimal)
    total_storage_gb FLOAT NULL,
    used_storage_gb FLOAT NULL,
    free_storage_gb FLOAT NULL,
    
    last_checked_at DATETIME DEFAULT CURRENT_TIMESTAMP
);