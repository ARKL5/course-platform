// ===== 创作者工作台 JavaScript =====

class CreatorDashboard {
    constructor() {
        this.currentUser = null;
        this.currentCourseId = null;
        this.currentCourseTitle = '';
        this.uploadedFiles = [];
        this.isUploading = false;
        
        this.init();
    }

    async init() {
        console.log('Creator Dashboard initializing...');
        
        // 验证用户身份（可选）
        await this.validateToken();
        
        // 初始化事件监听器
        this.initEventListeners();
        this.initDragAndDrop();
        
        // 加载统计数据（演示模式或真实数据）
        this.loadUserStats();
        
        // 显示演示模式提示（如果未登录）
        if (!this.currentUser) {
            this.showDemoModeNotice();
        }
        
        console.log('Creator Dashboard initialized successfully');
    }

    // 验证身份
    async validateToken() {
        try {
            const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken');
            if (!token) {
                return false;
            }

            const response = await fetch('/api/v1/me', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                this.currentUser = await response.json();
                this.updateUserProfile();
                return true;
            } else {
                this.clearAuthData();
                return false;
            }
        } catch (error) {
            console.error('Auth validation error:', error);
            this.clearAuthData();
            return false;
        }
    }

    updateUserProfile() {
        const userProfile = document.getElementById('userProfile');
        const authButtons = document.getElementById('authButtons');
        const userName = document.getElementById('userName');
        
        if (this.currentUser && userProfile && authButtons && userName) {
            userProfile.style.display = 'flex';
            authButtons.style.display = 'none';
            userName.textContent = this.currentUser.nickname || this.currentUser.username;
        }
    }

    clearAuthData() {
        localStorage.removeItem('authToken');
        localStorage.removeItem('userInfo');
        sessionStorage.removeItem('authToken');
        sessionStorage.removeItem('userInfo');
        this.currentUser = null;
    }

    // 显示演示模式通知
    showDemoModeNotice() {
        this.showNotification('演示模式：您可以体验所有功能，登录后可保存数据', 'info');
        
        // 在页面头部添加演示模式标识
        const header = document.querySelector('.creator-header');
        if (header && !document.querySelector('.demo-notice')) {
            const demoNotice = document.createElement('div');
            demoNotice.className = 'demo-notice';
            demoNotice.innerHTML = `
                <div class="demo-notice-content">
                    <i class="fas fa-flask"></i>
                    <span>演示模式 - 所有功能均可体验</span>
                    <a href="/login" class="demo-login-btn">登录保存数据</a>
                </div>
            `;
            header.appendChild(demoNotice);
        }
    }

    // 初始化事件监听器
    initEventListeners() {
        const courseForm = document.getElementById('courseCreationForm');
        if (courseForm) {
            courseForm.addEventListener('submit', this.handleCourseCreation.bind(this));
        }

        const resetBtn = document.getElementById('resetFormBtn');
        if (resetBtn) {
            resetBtn.addEventListener('click', this.resetForm.bind(this));
        }

        const backBtn = document.getElementById('backToFormBtn');
        if (backBtn) {
            backBtn.addEventListener('click', this.backToForm.bind(this));
        }

        const coverInput = document.getElementById('courseCover');
        if (coverInput) {
            coverInput.addEventListener('change', this.handleCoverUpload.bind(this));
        }

        const removeCoverBtn = document.getElementById('removeCoverBtn');
        if (removeCoverBtn) {
            removeCoverBtn.addEventListener('click', this.removeCover.bind(this));
        }

        const fileInput = document.getElementById('contentFileInput');
        if (fileInput) {
            fileInput.addEventListener('change', this.handleFileUpload.bind(this));
        }

        const refreshBtn = document.getElementById('refreshFilesBtn');
        if (refreshBtn) {
            refreshBtn.addEventListener('click', this.loadCourseFiles.bind(this));
        }

        const previewBtn = document.getElementById('previewCourseBtn');
        const saveDraftBtn = document.getElementById('saveDraftBtn');
        const publishBtn = document.getElementById('publishCourseBtn');

        if (previewBtn) previewBtn.addEventListener('click', this.previewCourse.bind(this));
        if (saveDraftBtn) saveDraftBtn.addEventListener('click', this.saveDraft.bind(this));
        if (publishBtn) publishBtn.addEventListener('click', this.publishCourse.bind(this));

        const logoutBtn = document.getElementById('logoutBtn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', this.handleLogout.bind(this));
        }
    }

    // 拖拽上传功能
    initDragAndDrop() {
        const fileUploadArea = document.getElementById('fileUploadArea');
        if (!fileUploadArea) return;

        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            fileUploadArea.addEventListener(eventName, this.preventDefaults, false);
        });

        ['dragenter', 'dragover'].forEach(eventName => {
            fileUploadArea.addEventListener(eventName, this.highlight.bind(this), false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            fileUploadArea.addEventListener(eventName, this.unhighlight.bind(this), false);
        });

        fileUploadArea.addEventListener('drop', this.handleDrop.bind(this), false);
    }

    preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    highlight(e) {
        const uploadZone = e.currentTarget.querySelector('.upload-zone');
        if (uploadZone) {
            uploadZone.style.borderColor = 'var(--accent-primary)';
            uploadZone.style.background = 'rgba(229, 9, 20, 0.1)';
        }
    }

    unhighlight(e) {
        const uploadZone = e.currentTarget.querySelector('.upload-zone');
        if (uploadZone) {
            uploadZone.style.borderColor = 'var(--border-color)';
            uploadZone.style.background = 'var(--bg-tertiary)';
        }
    }

    handleDrop(e) {
        const files = e.dataTransfer.files;
        if (files.length > 0 && this.currentCourseId) {
            this.uploadFiles(files);
        }
    }

    // 课程创建
    async handleCourseCreation(e) {
        e.preventDefault();
        
        if (this.isUploading) {
            this.showNotification('正在处理中，请稍候...', 'warning');
            return;
        }

        const formData = new FormData(e.target);
        
        // 手动构建课程数据对象，确保符合后端API要求
        const courseData = {
            title: formData.get('title') || '',
            description: formData.get('description') || '',
            instructor_id: 1, // 默认设置instructor_id为1，实际应用中可从用户信息获取
            category_id: parseInt(formData.get('category_id')) || 1,
            price: parseFloat(formData.get('price')) || 0,
            cover_image: '' // 暂时设为空字符串，文件上传在后续处理
        };

        if (!courseData.title || !courseData.description || !courseData.category_id || courseData.price === '') {
            this.showNotification('请填写所有必填字段', 'error');
            return;
        }

        // 保存选择的封面文件（如果有的话）
        const coverFile = formData.get('cover_image');
        if (coverFile && coverFile.size > 0) {
            // 存储文件信息，在课程创建成功后处理
            this.selectedCoverFile = coverFile;
        }

        try {
            this.showLoading(true);
            const createBtn = document.getElementById('createCourseBtn');
            if (createBtn) {
                createBtn.disabled = true;
                createBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> 创建中...';
            }

            const token = localStorage.getItem('authToken');
            
            // 如果没有token，使用演示模式
            if (!token) {
                // 演示模式：模拟课程创建
                await new Promise(resolve => setTimeout(resolve, 1500)); // 模拟网络延迟
                
                const mockCourseId = Date.now(); // 使用时间戳作为模拟ID
                this.currentCourseId = mockCourseId;
                this.currentCourseTitle = courseData.title;
                
                // 演示模式下处理封面文件
                if (this.selectedCoverFile) {
                    this.showNotification('演示模式：课程和封面创建成功！（数据仅用于演示）', 'success');
                } else {
                    this.showNotification('演示模式：课程创建成功！（数据仅用于演示）', 'success');
                }
                this.transitionToContentManagement();
            } else {
                // 正常模式：实际API调用
                console.log('发送课程数据:', courseData); // 调试信息
                
                const response = await fetch('/api/v1/courses', {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(courseData)
                });

                if (response.ok) {
                    const result = await response.json();
                    console.log('课程创建响应:', result); // 调试信息
                    
                    // 根据实际API响应结构获取课程数据
                    const newCourse = result.data || result;
                    this.currentCourseId = newCourse.id;
                    this.currentCourseTitle = newCourse.title;
                    
                    // 如果有选择封面文件，尝试上传
                    if (this.selectedCoverFile) {
                        try {
                            await this.uploadCoverImage(this.currentCourseId, this.selectedCoverFile);
                            this.showNotification('课程和封面创建成功！', 'success');
                        } catch (error) {
                            console.error('封面上传失败:', error);
                            this.showNotification('课程创建成功，但封面上传失败', 'warning');
                        }
                    } else {
                        this.showNotification('课程创建成功！', 'success');
                    }
                    
                    this.transitionToContentManagement();
                } else {
                    const errorData = await response.json();
                    console.error('API错误响应:', errorData); // 调试信息
                    throw new Error(errorData.message || '课程创建失败');
                }
            }
        } catch (error) {
            console.error('Course creation error:', error);
            this.showNotification('创建失败：' + error.message, 'error');
        } finally {
            this.showLoading(false);
            const createBtn = document.getElementById('createCourseBtn');
            if (createBtn) {
                createBtn.disabled = false;
                createBtn.innerHTML = '<i class="fas fa-rocket"></i> 创建课程并进入下一步';
            }
        }
    }

    // 切换到内容管理页面
    transitionToContentManagement() {
        const courseInfoSection = document.getElementById('courseInfoSection');
        const contentSection = document.getElementById('contentManagementSection');
        const currentTitleSpan = document.getElementById('currentCourseTitle');

        if (courseInfoSection) courseInfoSection.style.display = 'none';
        if (contentSection) contentSection.style.display = 'block';
        if (currentTitleSpan) currentTitleSpan.textContent = this.currentCourseTitle;

        this.loadCourseFiles();
        contentSection?.scrollIntoView({ behavior: 'smooth' });
    }

    // 返回表单
    backToForm() {
        const courseInfoSection = document.getElementById('courseInfoSection');
        const contentSection = document.getElementById('contentManagementSection');

        if (courseInfoSection) courseInfoSection.style.display = 'block';
        if (contentSection) contentSection.style.display = 'none';
        courseInfoSection?.scrollIntoView({ behavior: 'smooth' });
    }

    // 重置表单
    resetForm() {
        const form = document.getElementById('courseCreationForm');
        if (form) {
            form.reset();
            this.removeCover();
            this.showNotification('表单已重置', 'info');
        }
    }

    // 封面图片处理
    handleCoverUpload(e) {
        const file = e.target.files[0];
        if (!file) return;

        if (!file.type.startsWith('image/')) {
            this.showNotification('请选择图片文件', 'error');
            return;
        }

        if (file.size > 5 * 1024 * 1024) {
            this.showNotification('图片文件不能超过5MB', 'error');
            return;
        }

        const reader = new FileReader();
        reader.onload = (e) => {
            this.showCoverPreview(e.target.result);
        };
        reader.readAsDataURL(file);
    }

    showCoverPreview(imageSrc) {
        const placeholder = document.querySelector('.upload-placeholder');
        const preview = document.querySelector('.cover-preview');
        const previewImg = document.getElementById('coverPreview');

        if (placeholder) placeholder.style.display = 'none';
        if (preview) preview.style.display = 'block';
        if (previewImg) previewImg.src = imageSrc;
    }

    removeCover() {
        const placeholder = document.querySelector('.upload-placeholder');
        const preview = document.querySelector('.cover-preview');
        const coverInput = document.getElementById('courseCover');

        if (placeholder) placeholder.style.display = 'block';
        if (preview) preview.style.display = 'none';
        if (coverInput) coverInput.value = '';
        
        // 清除选择的封面文件
        this.selectedCoverFile = null;
    }

    // 上传封面图片
    async uploadCoverImage(courseId, file) {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('course_id', courseId);
        formData.append('file_type', 'image'); // 封面图片使用image类型

        const token = localStorage.getItem('authToken');
        const response = await fetch('/api/v1/content/upload', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + token
            },
            body: formData
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || '封面上传失败');
        }

        const result = await response.json();
        console.log('封面上传成功:', result);
        return result;
    }

    // 文件上传
    async handleFileUpload(e) {
        const files = e.target.files;
        if (files.length > 0 && this.currentCourseId) {
            await this.uploadFiles(files);
        }
    }

    async uploadFiles(files) {
        if (!this.currentCourseId) {
            this.showNotification('请先创建课程', 'warning');
            return;
        }

        if (this.isUploading) {
            this.showNotification('正在上传中，请稍候...', 'warning');
            return;
        }

        this.isUploading = true;
        const progressContainer = document.getElementById('uploadProgress');
        const progressFill = document.getElementById('progressFill');
        const progressPercentage = document.getElementById('uploadPercentage');
        const uploadStatus = document.getElementById('uploadStatus');

        try {
            if (progressContainer) progressContainer.style.display = 'block';

            const token = localStorage.getItem('authToken');

            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                
                if (uploadStatus) uploadStatus.textContent = '正在上传: ' + file.name;

                if (!token) {
                    // 演示模式：模拟上传进度
                    await new Promise(resolve => setTimeout(resolve, 800)); // 模拟上传时间
                    
                    // 添加到模拟文件列表
                    const mockFile = {
                        id: Date.now() + i,
                        filename: file.name,
                        file_type: this.getFileTypeFromName(file.name),
                        file_size: this.formatFileSize(file.size),
                        upload_date: new Date().toISOString(),
                        course_id: this.currentCourseId
                    };
                    
                    if (!this.uploadedFiles) this.uploadedFiles = [];
                    this.uploadedFiles.push(mockFile);
                    
                    const progress = Math.round(((i + 1) / files.length) * 100);
                    if (progressFill) progressFill.style.width = progress + '%';
                    if (progressPercentage) progressPercentage.textContent = progress + '%';
                } else {
                    // 正常模式：实际API调用
                    const formData = new FormData();
                    formData.append('file', file);
                    formData.append('course_id', this.currentCourseId);
                    formData.append('file_type', this.getFileTypeFromName(file.name)); // 添加文件类型

                    const response = await fetch('/api/v1/content/upload', {
                        method: 'POST',
                        headers: {
                            'Authorization': 'Bearer ' + token
                        },
                        body: formData
                    });

                    if (response.ok) {
                        const result = await response.json();
                        console.log('文件上传成功:', result);
                        
                        const progress = Math.round(((i + 1) / files.length) * 100);
                        if (progressFill) progressFill.style.width = progress + '%';
                        if (progressPercentage) progressPercentage.textContent = progress + '%';
                    } else {
                        const error = await response.json();
                        throw new Error(error.message || '文件 ' + file.name + ' 上传失败');
                    }
                }
            }

            const mode = token ? '' : '演示模式：';
            this.showNotification(mode + '成功上传 ' + files.length + ' 个文件', 'success');
            await this.loadCourseFiles();
            
        } catch (error) {
            console.error('文件上传错误:', error);
            this.showNotification('上传失败：' + error.message, 'error');
        } finally {
            this.isUploading = false;
            
            setTimeout(() => {
                if (progressContainer) progressContainer.style.display = 'none';
                if (progressFill) progressFill.style.width = '0%';
                if (progressPercentage) progressPercentage.textContent = '0%';
                if (uploadStatus) uploadStatus.textContent = '准备上传...';
            }, 2000);

            const fileInput = document.getElementById('contentFileInput');
            if (fileInput) fileInput.value = '';
        }
    }

    // 加载课程文件列表
    async loadCourseFiles() {
        if (!this.currentCourseId) return;

        try {
            const token = localStorage.getItem('authToken');
            
            if (!token) {
                // 演示模式：使用模拟数据或已上传的文件
                if (!this.uploadedFiles) {
                    this.uploadedFiles = [];
                }
                this.renderFilesList();
                return;
            }

            // 正常模式：从API获取
            const response = await fetch('/api/v1/content/files?course_id=' + this.currentCourseId, {
                method: 'GET',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                const result = await response.json();
                // API返回的数据格式：{code: "SUCCESS", data: {files: [...], total: N}}
                this.uploadedFiles = (result.data && result.data.files) ? result.data.files : [];
                this.renderFilesList();
            } else {
                console.error('加载文件列表失败');
                this.uploadedFiles = [];
                this.renderFilesList();
            }
        } catch (error) {
            console.error('加载文件列表错误:', error);
            this.uploadedFiles = [];
            this.renderFilesList();
        }
    }

    // 渲染文件列表
    renderFilesList() {
        const filesEmpty = document.getElementById('filesEmpty');
        const filesList = document.getElementById('filesList');
        const totalFiles = document.getElementById('totalFiles');

        if (totalFiles) {
            totalFiles.textContent = this.uploadedFiles.length;
        }

        if (this.uploadedFiles.length === 0) {
            if (filesEmpty) filesEmpty.style.display = 'flex';
            if (filesList) filesList.style.display = 'none';
            return;
        }

        if (filesEmpty) filesEmpty.style.display = 'none';
        if (filesList) {
            filesList.style.display = 'block';
            filesList.innerHTML = this.uploadedFiles.map(file => this.createFileItem(file)).join('');
        }
    }

    // 创建文件项
    createFileItem(file) {
        // 兼容不同的字段名格式（API返回的是snake_case）
        const fileId = file.file_id || file.id;
        const fileName = file.file_name || file.filename;
        const fileType = file.file_type || file.type;
        const fileSize = file.file_size || file.size;
        const createdAt = file.created_at || file.upload_date;
        
        const fileIcon = this.getFileIcon(fileType);
        const formattedSize = this.formatFileSize(fileSize);
        
        return `
            <div class="file-item" data-file-id="${fileId}">
                <div class="file-icon">
                    <i class="fas ${fileIcon}"></i>
                </div>
                <div class="file-info">
                    <h4 class="file-name" title="${fileName}">${fileName}</h4>
                    <div class="file-meta">
                        <span>类型: ${fileType}</span>
                        <span>大小: ${formattedSize}</span>
                        <span>上传时间: ${this.formatDate(createdAt)}</span>
                    </div>
                </div>
                <div class="file-actions">
                    <button class="file-action-btn" onclick="creatorDashboard.downloadFile('${fileId}')" title="下载">
                        <i class="fas fa-download"></i>
                    </button>
                    <button class="file-action-btn delete" onclick="creatorDashboard.deleteFile('${fileId}')" title="删除">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `;
    }

    // 获取文件图标
    getFileIcon(fileType) {
        const type = fileType.toLowerCase();
        if (type.includes('video')) return 'fa-video';
        if (type.includes('audio')) return 'fa-music';
        if (type.includes('pdf')) return 'fa-file-pdf';
        if (type.includes('doc')) return 'fa-file-word';
        if (type.includes('ppt')) return 'fa-file-powerpoint';
        if (type.includes('xls')) return 'fa-file-excel';
        if (type.includes('image')) return 'fa-file-image';
        if (type.includes('text')) return 'fa-file-alt';
        return 'fa-file';
    }

    // 从文件名获取文件类型（后端支持：image, video, document, audio, other）
    getFileTypeFromName(filename) {
        const ext = filename.toLowerCase().split('.').pop();
        const typeMap = {
            // 视频文件
            'mp4': 'video', 'avi': 'video', 'mov': 'video', 'wmv': 'video', 'flv': 'video', 'webm': 'video',
            // 图片文件
            'jpg': 'image', 'jpeg': 'image', 'png': 'image', 'gif': 'image', 'bmp': 'image', 'webp': 'image',
            // 文档文件
            'pdf': 'document', 'doc': 'document', 'docx': 'document', 'txt': 'document', 'rtf': 'document',
            'ppt': 'document', 'pptx': 'document', 'md': 'document',
            // 音频文件
            'mp3': 'audio', 'wav': 'audio', 'aac': 'audio', 'flac': 'audio', 'ogg': 'audio'
        };
        return typeMap[ext] || 'other';
    }

    // 格式化文件大小
    formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    // 格式化日期
    formatDate(dateString) {
        if (!dateString) return '未知';
        const date = new Date(dateString);
        return date.toLocaleDateString('zh-CN');
    }

    // 删除文件
    async deleteFile(fileId) {
        if (!confirm('确定要删除这个文件吗？此操作不可撤销。')) {
            return;
        }

        try {
            const token = localStorage.getItem('authToken');
            
            if (!token) {
                // 演示模式：从本地数组中删除
                this.uploadedFiles = this.uploadedFiles.filter(file => file.id !== fileId);
                this.showNotification('演示模式：文件删除成功', 'success');
                this.renderFilesList();
                return;
            }

            // 正常模式：API删除
            const response = await fetch('/api/v1/content/files/' + fileId, {
                method: 'DELETE',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                this.showNotification('文件删除成功', 'success');
                await this.loadCourseFiles();
            } else {
                const error = await response.json();
                throw new Error(error.message || '删除失败');
            }
        } catch (error) {
            console.error('删除文件错误:', error);
            this.showNotification('删除失败：' + error.message, 'error');
        }
    }

    // 下载文件
    downloadFile(fileId) {
        const token = localStorage.getItem('authToken');
        
        if (!token) {
            // 演示模式：显示提示
            this.showNotification('演示模式：下载功能需要登录后使用', 'info');
            return;
        }

        // 正常模式：实际下载
        const downloadUrl = '/api/v1/content/files/' + fileId + '/download?token=' + encodeURIComponent(token);
        
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.style.display = 'none';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    // 预览课程
    previewCourse() {
        if (!this.currentCourseId) {
            this.showNotification('请先创建课程', 'warning');
            return;
        }
        
        window.open('/course/' + this.currentCourseId, '_blank');
    }

    // 保存草稿
    async saveDraft() {
        if (!this.currentCourseId) {
            this.showNotification('请先创建课程', 'warning');
            return;
        }

        this.showNotification('草稿已自动保存', 'success');
    }

    // 发布课程
    async publishCourse() {
        if (!this.currentCourseId) {
            this.showNotification('请先创建课程', 'warning');
            return;
        }

        if (this.uploadedFiles.length === 0) {
            this.showNotification('请至少上传一个课程文件', 'warning');
            return;
        }

        if (!confirm('确定要发布这个课程吗？发布后学员将可以看到这个课程。')) {
            return;
        }

        try {
            const token = localStorage.getItem('authToken');
            
            if (!token) {
                // 演示模式：模拟发布
                await new Promise(resolve => setTimeout(resolve, 1000));
                this.showNotification('演示模式：课程发布成功！（仅用于演示）', 'success');
                setTimeout(() => {
                    this.showNotification('演示模式完成，您可以登录后使用完整功能', 'info');
                }, 2000);
                return;
            }

            // 正常模式：实际发布
            const response = await fetch('/api/v1/courses/' + this.currentCourseId + '/publish', {
                method: 'POST',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                this.showNotification('课程发布成功！', 'success');
                setTimeout(() => {
                    window.location.href = '/dashboard';
                }, 2000);
            } else {
                const error = await response.json();
                throw new Error(error.message || '发布失败');
            }
        } catch (error) {
            console.error('发布课程错误:', error);
            this.showNotification('发布失败：' + error.message, 'error');
        }
    }

    // 加载用户统计
    async loadUserStats() {
        try {
            const token = localStorage.getItem('authToken');
            
            if (!token) {
                // 演示模式：使用演示数据
                const demoStats = {
                    total_courses: 8,
                    total_students: 1250,
                    total_revenue: 8450.50,
                    active_courses: 6,
                    draft_courses: 2,
                    this_month_students: 180,
                    this_month_revenue: 1250.00
                };
                this.updateStatsDisplay(demoStats);
                return;
            }

            // 正常模式：从API获取
            const response = await fetch('/api/v1/creator/stats', {
                method: 'GET',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                const stats = await response.json();
                this.updateStatsDisplay(stats);
            } else {
                throw new Error('获取统计数据失败');
            }
        } catch (error) {
            console.error('加载统计信息错误:', error);
            this.updateStatsDisplay({
                total_courses: 0,
                total_students: 0,
                total_revenue: 0
            });
        }
    }

    // 更新统计显示
    updateStatsDisplay(stats) {
        const totalCourses = document.getElementById('totalCourses');
        const totalStudents = document.getElementById('totalStudents');
        const totalRevenue = document.getElementById('totalRevenue');

        if (totalCourses) totalCourses.textContent = stats.total_courses || 0;
        if (totalStudents) totalStudents.textContent = stats.total_students || 0;
        if (totalRevenue) totalRevenue.textContent = '¥' + (stats.total_revenue || 0).toFixed(2);
    }

    // 显示通知
    showNotification(message, type = 'info') {
        const container = document.getElementById('notificationContainer');
        if (!container) return;

        const notification = document.createElement('div');
        notification.className = 'notification notification-' + type;
        
        const iconMap = {
            success: 'fa-check-circle',
            error: 'fa-exclamation-circle',
            warning: 'fa-exclamation-triangle',
            info: 'fa-info-circle'
        };

        notification.innerHTML = `
            <div class="notification-content">
                <i class="fas ${iconMap[type]}"></i>
                <span>${message}</span>
            </div>
            <button class="notification-close">
                <i class="fas fa-times"></i>
            </button>
        `;

        notification.querySelector('.notification-close').addEventListener('click', () => {
            this.removeNotification(notification);
        });

        container.appendChild(notification);
        setTimeout(() => notification.classList.add('show'), 100);
        setTimeout(() => this.removeNotification(notification), 5000);
    }

    // 移除通知
    removeNotification(notification) {
        if (notification && notification.parentNode) {
            notification.classList.remove('show');
            setTimeout(() => {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
            }, 300);
        }
    }

    // 显示加载状态
    showLoading(show) {
        const overlay = document.getElementById('loadingOverlay');
        if (overlay) {
            overlay.style.display = show ? 'flex' : 'none';
        }
    }

    // 退出登录
    handleLogout() {
        if (confirm('确定要退出登录吗？')) {
            this.clearAuthData();
            window.location.href = '/';
        }
    }
}

// 全局实例
let creatorDashboard;

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    creatorDashboard = new CreatorDashboard();
});

// 导出给全局使用
window.creatorDashboard = creatorDashboard;