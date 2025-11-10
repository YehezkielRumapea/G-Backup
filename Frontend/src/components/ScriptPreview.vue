<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isVisible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-container">
          <!-- HEADER -->
          <div class="modal-header">
            <div class="header-left">
              <span class="icon">📜</span>
              <h3>Script Preview</h3>
              <span class="job-id">Job ID: {{ jobId }}</span>
            </div>
            <button @click="handleClose" class="close-btn" aria-label="Close">
              <span>×</span>
            </button>
          </div>

          <!-- TOOLBAR -->
          <div class="modal-toolbar">
            <div class="toolbar-left">
              <span class="language-badge">{{ detectedLanguage }}</span>
              <span class="line-count">{{ lineCount }} lines</span>
            </div>
            <div class="toolbar-right">
              <button @click="copyToClipboard" class="toolbar-btn" :class="{ copied: isCopied }">
                <span v-if="!isCopied">📋 Copy</span>
                <span v-else>✅ Copied!</span>
              </button>
              <button @click="downloadScript" class="toolbar-btn">
                💾 Download
              </button>
            </div>
          </div>

          <!-- CODE CONTENT -->
          <div class="modal-body">
            <div class="code-container">
              <!-- Line Numbers -->
              <div class="line-numbers">
                <div v-for="n in lineCount" :key="n" class="line-number">
                  {{ n }}
                </div>
              </div>

              <!-- Code with Syntax Highlighting -->
              <pre class="code-content"><code 
                ref="codeElement" 
                :class="`language-${detectedLanguage}`"
                v-html="highlightedCode"
              ></code></pre>
            </div>
          </div>

          <!-- FOOTER -->
          <div class="modal-footer">
            <div class="footer-info">
              <span>💡 Tip: Scroll untuk melihat lebih banyak code</span>
            </div>
            <button @click="handleClose" class="btn-close">Close</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';
import hljs from 'highlight.js';

// Import theme CSS (pilih salah satu)
import 'highlight.js/styles/atom-one-dark.css'; // Dark theme
// import 'highlight.js/styles/github.css'; // Light theme
// import 'highlight.js/styles/monokai-sublime.css'; // Alternative dark
// import 'highlight.js/styles/vs2015.css'; // VS Code dark

const props = defineProps({
  isVisible: {
    type: Boolean,
    required: true
  },
  jobId: {
    type: [String, Number],
    default: null
  },
  scriptContent: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['close']);

const codeElement = ref(null);
const isCopied = ref(false);

// Auto-detect language from script content
const detectedLanguage = computed(() => {
  const content = props.scriptContent.toLowerCase();
  
  if (content.includes('#!/bin/bash') || content.includes('#!/bin/sh')) {
    return 'bash';
  } else if (content.includes('<?php')) {
    return 'php';
  } else if (content.includes('import ') || content.includes('def ')) {
    return 'python';
  } else if (content.includes('function') || content.includes('const ') || content.includes('let ')) {
    return 'javascript';
  } else if (content.includes('mysqldump') || content.includes('rclone')) {
    return 'bash';
  }
  
  return 'bash'; // default
});

// Count lines
const lineCount = computed(() => {
  if (!props.scriptContent) return 0;
  return props.scriptContent.split('\n').length;
});

// Highlighted code
const highlightedCode = computed(() => {
  if (!props.scriptContent) return '';
  
  try {
    const result = hljs.highlight(props.scriptContent, {
      language: detectedLanguage.value
    });
    return result.value;
  } catch (error) {
    console.error('Highlight error:', error);
    // Fallback: escape HTML manually
    return props.scriptContent
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;');
  }
});

// Watch for visibility changes
watch(() => props.isVisible, (newVal) => {
  if (newVal) {
    // Prevent body scroll when modal is open
    document.body.style.overflow = 'hidden';
  } else {
    // Restore body scroll
    document.body.style.overflow = '';
    isCopied.value = false;
  }
});

// Close modal
function handleClose() {
  emit('close');
}

// Copy to clipboard
async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(props.scriptContent);
    isCopied.value = true;
    
    setTimeout(() => {
      isCopied.value = false;
    }, 2000);
  } catch (error) {
    console.error('Copy failed:', error);
    alert('Failed to copy to clipboard');
  }
}

// Download script
function downloadScript() {
  const blob = new Blob([props.scriptContent], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  
  link.href = url;
  link.download = `job_${props.jobId}_script.sh`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

// Handle ESC key
function handleKeydown(e) {
  if (e.key === 'Escape' && props.isVisible) {
    handleClose();
  }
}

// Add/remove event listener
watch(() => props.isVisible, (newVal) => {
  if (newVal) {
    window.addEventListener('keydown', handleKeydown);
  } else {
    window.removeEventListener('keydown', handleKeydown);
  }
});
</script>

<style scoped>
/* ==================== MODAL OVERLAY ==================== */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

/* ==================== MODAL CONTAINER ==================== */
.modal-container {
  background: #1e1e1e;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  width: 90%;
  max-width: 1000px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* ==================== HEADER ==================== */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-bottom: 2px solid rgba(255, 255, 255, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.icon {
  font-size: 1.5rem;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.job-id {
  background: rgba(255, 255, 255, 0.2);
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.close-btn {
  background: rgba(255, 255, 255, 0.15);
  border: none;
  color: white;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  font-size: 1.8rem;
  line-height: 1;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: scale(1.1);
}

/* ==================== TOOLBAR ==================== */
.modal-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: #252526;
  border-bottom: 1px solid #3e3e42;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.language-badge {
  background: #007acc;
  color: white;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.line-count {
  color: #858585;
  font-size: 0.85rem;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.toolbar-btn {
  background: #3c3c3c;
  border: 1px solid #555;
  color: #d4d4d4;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.toolbar-btn:hover {
  background: #505050;
  border-color: #666;
}

.toolbar-btn.copied {
  background: #0e7c3e;
  border-color: #0e7c3e;
  color: white;
}

/* ==================== BODY (CODE CONTAINER) ==================== */
.modal-body {
  flex: 1;
  overflow: auto;
  background: #1e1e1e;
}

.code-container {
  display: flex;
  min-height: 100%;
}

/* Line Numbers */
.line-numbers {
  background: #1e1e1e;
  color: #858585;
  padding: 16px 12px;
  text-align: right;
  user-select: none;
  border-right: 1px solid #3e3e42;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  min-width: 50px;
}

.line-number {
  height: 20.8px; /* Match code line height */
}

/* Code Content */
.code-content {
  flex: 1;
  margin: 0;
  padding: 16px 20px;
  background: transparent;
  color: #d4d4d4;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  overflow-x: auto;
  white-space: pre;
  tab-size: 4;
}

.code-content code {
  display: block;
  color: inherit;
}

/* Custom scrollbar untuk code area */
.modal-body::-webkit-scrollbar {
  width: 12px;
  height: 12px;
}

.modal-body::-webkit-scrollbar-track {
  background: #1e1e1e;
}

.modal-body::-webkit-scrollbar-thumb {
  background: #424242;
  border-radius: 6px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: #4e4e4e;
}

/* ==================== FOOTER ==================== */
.modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #252526;
  border-top: 1px solid #3e3e42;
}

.footer-info {
  color: #858585;
  font-size: 0.85rem;
}

.btn-close {
  background: #0e639c;
  border: none;
  color: white;
  padding: 8px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #1177bb;
  transform: translateY(-1px);
}

/* ==================== TRANSITIONS ==================== */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9) translateY(-20px);
}

/* ==================== RESPONSIVE ==================== */
@media (max-width: 768px) {
  .modal-container {
    width: 95%;
    max-height: 85vh;
  }
  
  .modal-header {
    padding: 16px;
  }
  
  .header-left {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .modal-header h3 {
    font-size: 1.1rem;
  }
  
  .modal-toolbar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .toolbar-right {
    justify-content: stretch;
  }
  
  .toolbar-btn {
    flex: 1;
    justify-content: center;
  }
  
  .line-numbers {
    min-width: 40px;
    padding: 16px 8px;
  }
  
  .code-content {
    font-size: 12px;
    padding: 16px 12px;
  }
}
</style>