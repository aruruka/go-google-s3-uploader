// Upload Page Specific JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log('📷 Upload page JavaScript loaded');
    
    const fileInput = document.getElementById('file');
    const filePreview = document.getElementById('filePreview');
    const previewImage = document.getElementById('previewImage');
    const fileName = document.getElementById('fileName');
    const uploadForm = document.getElementById('uploadForm');
    const uploadBtn = document.getElementById('uploadBtn');
    const progressBar = document.getElementById('progressBar');
    const progressFill = document.getElementById('progressFill');

    // File preview functionality
    if (fileInput) {
        fileInput.addEventListener('change', function(e) {
            const file = e.target.files[0];
            if (file) {
                // File size validation
                const maxSize = 10 * 1024 * 1024; // 10MB
                if (file.size > maxSize) {
                    alert('File size exceeds 10MB limit. Please choose a smaller file.');
                    fileInput.value = '';
                    return;
                }
                
                // File type validation
                const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
                if (!allowedTypes.includes(file.type)) {
                    alert('Please select a valid image file (JPG, PNG, GIF, or WebP).');
                    fileInput.value = '';
                    return;
                }

                // Show preview
                const reader = new FileReader();
                reader.onload = function(e) {
                    if (previewImage) {
                        previewImage.src = e.target.result;
                    }
                    if (fileName) {
                        const fileSizeFormatted = window.AppUtils ? 
                            window.AppUtils.formatFileSize(file.size) : 
                            (file.size / 1024 / 1024).toFixed(2) + ' MB';
                        fileName.textContent = `Selected: ${file.name} (${fileSizeFormatted})`;
                    }
                    if (filePreview) {
                        filePreview.style.display = 'block';
                    }
                };
                reader.readAsDataURL(file);
            } else {
                if (filePreview) {
                    filePreview.style.display = 'none';
                }
            }
        });
    }

    // Enhanced form submission with progress tracking
    if (uploadForm) {
        uploadForm.addEventListener('submit', function(e) {
            const file = fileInput ? fileInput.files[0] : null;
            
            if (!file) {
                e.preventDefault();
                alert('Please select a file first.');
                return;
            }

            // Show loading state
            if (uploadBtn) {
                if (window.AppUtils) {
                    window.AppUtils.showLoading(uploadBtn, '⏳ Uploading...');
                } else {
                    uploadBtn.disabled = true;
                    uploadBtn.textContent = '⏳ Uploading...';
                }
            }

            // Show progress bar
            if (progressBar) {
                progressBar.style.display = 'block';
            }

            // For actual uploads, we'll use XMLHttpRequest to track progress
            if (window.location.pathname !== '/upload' || uploadForm.action.includes('api/upload')) {
                // This is a real form submission, let it proceed
                return;
            }

            // Demo mode - simulate upload progress
            e.preventDefault();
            simulateUploadProgress();
        });
    }

    // Upload progress simulation (for demo/testing)
    function simulateUploadProgress() {
        let progress = 0;
        const progressInterval = setInterval(() => {
            progress += Math.random() * 15;
            if (progress > 100) {
                progress = 100;
                clearInterval(progressInterval);
                
                setTimeout(() => {
                    if (window.AppUtils) {
                        window.AppUtils.showNotification('Upload completed successfully!', 'success');
                    }
                    
                    // Reset form
                    if (uploadForm) {
                        uploadForm.reset();
                    }
                    if (filePreview) {
                        filePreview.style.display = 'none';
                    }
                    if (progressBar) {
                        progressBar.style.display = 'none';
                    }
                    if (uploadBtn) {
                        if (window.AppUtils) {
                            window.AppUtils.hideLoading(uploadBtn);
                        } else {
                            uploadBtn.disabled = false;
                            uploadBtn.textContent = '🚀 Upload to S3';
                        }
                    }
                }, 500);
            }
            
            if (progressFill) {
                progressFill.style.width = progress + '%';
            }
        }, 200);
    }

    // Real upload progress tracking using XMLHttpRequest
    function uploadWithProgress(formData, url) {
        return new Promise((resolve, reject) => {
            const xhr = new XMLHttpRequest();
            
            // Track upload progress
            xhr.upload.addEventListener('progress', function(e) {
                if (e.lengthComputable) {
                    const percentComplete = (e.loaded / e.total) * 100;
                    if (progressFill) {
                        progressFill.style.width = percentComplete + '%';
                    }
                }
            });
            
            xhr.addEventListener('load', function() {
                if (xhr.status >= 200 && xhr.status < 300) {
                    resolve(xhr.response);
                } else {
                    reject(new Error('Upload failed'));
                }
            });
            
            xhr.addEventListener('error', function() {
                reject(new Error('Upload failed'));
            });
            
            xhr.open('POST', url);
            xhr.send(formData);
        });
    }

    // Drag and drop support
    const fileInputArea = document.querySelector('.file-input');
    if (fileInputArea) {
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            fileInputArea.addEventListener(eventName, preventDefaults, false);
        });

        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }

        ['dragenter', 'dragover'].forEach(eventName => {
            fileInputArea.addEventListener(eventName, highlight, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            fileInputArea.addEventListener(eventName, unhighlight, false);
        });

        function highlight(e) {
            fileInputArea.style.borderColor = '#007bff';
            fileInputArea.style.backgroundColor = '#e7f3ff';
        }

        function unhighlight(e) {
            fileInputArea.style.borderColor = '#ddd';
            fileInputArea.style.backgroundColor = '#f8f9fa';
        }

        fileInputArea.addEventListener('drop', handleDrop, false);

        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;
            
            if (files.length > 0 && fileInput) {
                fileInput.files = files;
                // Trigger change event
                const event = new Event('change', { bubbles: true });
                fileInput.dispatchEvent(event);
            }
        }
    }
});

// Export utility functions for upload functionality
window.UploadUtils = {
    // Reset upload form
    resetForm: function() {
        const uploadForm = document.getElementById('uploadForm');
        const filePreview = document.getElementById('filePreview');
        const progressBar = document.getElementById('progressBar');
        
        if (uploadForm) {
            uploadForm.reset();
        }
        if (filePreview) {
            filePreview.style.display = 'none';
        }
        if (progressBar) {
            progressBar.style.display = 'none';
        }
    },
    
    // Show upload success
    showSuccess: function(fileName) {
        if (window.AppUtils) {
            window.AppUtils.showNotification(`Successfully uploaded: ${fileName}`, 'success');
        }
        this.resetForm();
    },
    
    // Show upload error
    showError: function(message) {
        if (window.AppUtils) {
            window.AppUtils.showNotification(`Upload failed: ${message}`, 'error');
        }
        
        const uploadBtn = document.getElementById('uploadBtn');
        if (uploadBtn && window.AppUtils) {
            window.AppUtils.hideLoading(uploadBtn);
        }
    }
};
