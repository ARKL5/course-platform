// ===== 课程详情页面交互功能 =====

// 全局变量
let currentLessonId = null;
let courseId = null;
let isPlaying = false;
let currentTime = 0;
let duration = 0;

// DOM 元素
const videoPlayer = document.getElementById('videoPlayer');
const playPauseBtn = document.getElementById('playPauseBtn');
const lessonsList = document.getElementById('lessonsList');

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    initializePage();
    bindEventListeners();
    updateVideoControls();
    updateMainCtaButton();
});

// 初始化页面
function initializePage() {
    // 获取当前课程和课时信息
    const activeLesson = document.querySelector('.lesson-item.active');
    if (activeLesson) {
        currentLessonId = activeLesson.getAttribute('data-lesson-id');
    }
    
    // 从URL获取课程ID
    const urlParams = new URLSearchParams(window.location.search);
    courseId = urlParams.get('course_id') || getQueryParam('id');
    
    console.log('课程详情页面初始化完成', {
        courseId: courseId,
        currentLessonId: currentLessonId
    });
}

// 绑定事件监听器
function bindEventListeners() {
    // 播放/暂停按钮
    if (playPauseBtn) {
        playPauseBtn.addEventListener('click', togglePlayPause);
    }
    
    // 进度条点击
    const progressBar = document.querySelector('.progress-bar');
    if (progressBar) {
        progressBar.addEventListener('click', handleProgressClick);
    }
    
    // 课程列表点击事件
    if (lessonsList) {
        lessonsList.addEventListener('click', handleLessonClick);
    }
    
    // 操作按钮事件
    bindActionButtons();
    
    // 键盘快捷键
    document.addEventListener('keydown', handleKeyboardShortcuts);
}

// 选择课程
function selectLesson(lessonId, lessonTitle, courseIdParam) {
    console.log('选择课程:', { lessonId, lessonTitle, courseIdParam });
    
    // 更新当前课程信息
    currentLessonId = lessonId;
    courseId = courseIdParam;
    
    // 更新UI
    updateActiveLesson(lessonId);
    updateVideoPlayer(lessonTitle);
    
    // 记录学习行为
    trackLessonSelection(lessonId, courseIdParam);
    
    // 更新URL参数
    updateURLParams(lessonId);
}

// 更新活跃课程状态
function updateActiveLesson(lessonId) {
    // 移除所有活跃状态
    document.querySelectorAll('.lesson-item').forEach(item => {
        item.classList.remove('active');
        const playIcon = item.querySelector('.lesson-status i');
        if (playIcon) {
            playIcon.className = 'fas fa-play-circle';
        }
    });
    
    // 添加新的活跃状态
    const selectedLesson = document.querySelector(`[data-lesson-id="${lessonId}"]`);
    if (selectedLesson) {
        selectedLesson.classList.add('active');
        const playIcon = selectedLesson.querySelector('.lesson-status i');
        if (playIcon) {
            playIcon.className = 'fas fa-play-circle playing-icon';
        }
    }
}

// 更新视频播放器
function updateVideoPlayer(lessonTitle) {
    const currentLessonTitleEl = document.querySelector('.current-lesson-title');
    const videoContent = document.querySelector('.video-content');
    
    if (currentLessonTitleEl) {
        currentLessonTitleEl.textContent = lessonTitle;
    }
    
    // 显示视频内容区域
    if (videoContent) {
        videoContent.style.display = 'flex';
    }
    
    // 重置播放状态
    isPlaying = false;
    updatePlayButton();
}

// 处理课程列表点击
function handleLessonClick(event) {
    const lessonItem = event.target.closest('.lesson-item');
    if (lessonItem) {
        const lessonId = lessonItem.getAttribute('data-lesson-id');
        const lessonTitle = lessonItem.querySelector('.lesson-title').textContent;
        selectLesson(parseInt(lessonId), lessonTitle, courseId);
    }
}

// 播放/暂停切换
function togglePlayPause() {
    if (!currentLessonId) {
        showNotification('请先选择一个课程', 'warning');
        return;
    }
    
    isPlaying = !isPlaying;
    updatePlayButton();
    
    if (isPlaying) {
        trackVideoPlay('course-detail', `${courseId}-${currentLessonId}`);
        simulateVideoProgress();
    }
    
    console.log(isPlaying ? '开始播放' : '暂停播放');
}

// 更新播放按钮状态
function updatePlayButton() {
    const playBtnIcon = playPauseBtn?.querySelector('i');
    if (playBtnIcon) {
        playBtnIcon.className = isPlaying ? 'fas fa-pause' : 'fas fa-play';
    }
}

// 处理进度条点击
function handleProgressClick(event) {
    if (!currentLessonId) return;
    
    const progressBar = event.currentTarget;
    const rect = progressBar.getBoundingClientRect();
    const clickPosition = (event.clientX - rect.left) / rect.width;
    
    // 模拟设置播放进度
    const newTime = clickPosition * duration;
    setVideoProgress(clickPosition);
    
    console.log('跳转到进度:', Math.floor(clickPosition * 100) + '%');
}

// 设置视频进度
function setVideoProgress(progress) {
    const progressFilled = document.querySelector('.progress-filled');
    const progressHandle = document.querySelector('.progress-handle');
    
    if (progressFilled) {
        progressFilled.style.width = (progress * 100) + '%';
    }
    
    currentTime = progress * duration;
    updateTimeDisplay();
}

// 模拟视频播放进度
function simulateVideoProgress() {
    if (!isPlaying) return;
    
    // 模拟总时长
    if (duration === 0) {
        duration = 300; // 5分钟示例
    }
    
    const interval = setInterval(() => {
        if (!isPlaying) {
            clearInterval(interval);
            return;
        }
        
        currentTime += 1;
        if (currentTime >= duration) {
            currentTime = duration;
            isPlaying = false;
            updatePlayButton();
            clearInterval(interval);
            onVideoEnded();
            return;
        }
        
        const progress = currentTime / duration;
        setVideoProgress(progress);
    }, 1000);
}

// 视频播放结束
function onVideoEnded() {
    console.log('课程播放完成');
    showNotification('课程播放完成！', 'success');
    
    // 自动播放下一课程
    const currentLesson = document.querySelector('.lesson-item.active');
    const nextLesson = currentLesson?.nextElementSibling;
    
    if (nextLesson && nextLesson.classList.contains('lesson-item')) {
        setTimeout(() => {
            nextLesson.click();
            showNotification('自动播放下一课程', 'info');
        }, 2000);
    }
}

// 更新时间显示
function updateTimeDisplay() {
    const timeDisplay = document.querySelector('.time-display');
    if (timeDisplay) {
        const currentFormatted = formatTime(currentTime);
        const durationFormatted = formatTime(duration);
        timeDisplay.textContent = `${currentFormatted} / ${durationFormatted}`;
    }
}

// 格式化时间
function formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = Math.floor(seconds % 60);
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}

// 绑定操作按钮事件
function bindActionButtons() {
    // 主操作按钮
    const mainCtaBtn = document.getElementById('mainCtaBtn');
    if (mainCtaBtn) {
        mainCtaBtn.addEventListener('click', () => {
            handleMainCtaClick();
        });
    }
    
    // 收藏课程
    const bookmarkBtn = document.getElementById('bookmarkBtn');
    if (bookmarkBtn) {
        bookmarkBtn.addEventListener('click', () => {
            toggleBookmark();
        });
    }
    
    // 分享课程
    const shareBtn = document.querySelector('.action-btn:nth-child(2)');
    if (shareBtn) {
        shareBtn.addEventListener('click', () => {
            shareCourse();
        });
    }
    
    // 添加到列表
    const addToListBtn = document.querySelector('.action-btn:nth-child(3)');
    if (addToListBtn) {
        addToListBtn.addEventListener('click', () => {
            addToPlaylist();
        });
    }
    
    // 下载资料
    const downloadBtn = document.querySelector('.action-btn:nth-child(4)');
    if (downloadBtn) {
        downloadBtn.addEventListener('click', () => {
            downloadMaterials();
        });
    }
}

// 收藏课程切换
function toggleBookmark() {
    const bookmarkBtn = document.querySelector('.action-btn:nth-child(1)');
    const icon = bookmarkBtn?.querySelector('i');
    
    if (icon?.classList.contains('fas')) {
        icon.className = 'far fa-bookmark';
        showNotification('已取消收藏', 'info');
    } else {
        icon.className = 'fas fa-bookmark';
        showNotification('已收藏课程', 'success');
    }
    
    // 这里可以添加API调用来保存收藏状态
    console.log('切换收藏状态');
}

// 分享课程
function shareCourse() {
    const shareData = {
        title: document.querySelector('.course-title')?.textContent || '精彩课程',
        text: '我在Course Platform上发现了一门很棒的课程，推荐给你！',
        url: window.location.href
    };
    
    if (navigator.share) {
        navigator.share(shareData).catch(console.error);
    } else {
        // 复制链接到剪贴板
        navigator.clipboard.writeText(window.location.href).then(() => {
            showNotification('课程链接已复制到剪贴板', 'success');
        }).catch(() => {
            showNotification('分享失败，请手动复制链接', 'error');
        });
    }
}

// 添加到播放列表
function addToPlaylist() {
    showNotification('已添加到我的学习列表', 'success');
    console.log('添加到播放列表');
}

// 下载课程资料
function downloadMaterials() {
    showNotification('开始下载课程资料...', 'info');
    
    // 模拟下载过程
    setTimeout(() => {
        showNotification('课程资料下载完成', 'success');
    }, 2000);
    
    console.log('下载课程资料');
}

// 键盘快捷键处理
function handleKeyboardShortcuts(event) {
    // 只在非输入框元素上响应快捷键
    if (event.target.tagName === 'INPUT' || event.target.tagName === 'TEXTAREA') {
        return;
    }
    
    switch (event.code) {
        case 'Space':
            event.preventDefault();
            togglePlayPause();
            break;
        case 'ArrowLeft':
            event.preventDefault();
            seekVideo(-10);
            break;
        case 'ArrowRight':
            event.preventDefault();
            seekVideo(10);
            break;
        case 'ArrowUp':
            event.preventDefault();
            selectPreviousLesson();
            break;
        case 'ArrowDown':
            event.preventDefault();
            selectNextLesson();
            break;
    }
}

// 视频快进/快退
function seekVideo(seconds) {
    if (!currentLessonId) return;
    
    currentTime = Math.max(0, Math.min(duration, currentTime + seconds));
    const progress = currentTime / duration;
    setVideoProgress(progress);
    
    showNotification(`${seconds > 0 ? '快进' : '快退'} ${Math.abs(seconds)} 秒`, 'info');
}

// 选择上一课程
function selectPreviousLesson() {
    const currentLesson = document.querySelector('.lesson-item.active');
    const previousLesson = currentLesson?.previousElementSibling;
    
    if (previousLesson && previousLesson.classList.contains('lesson-item')) {
        previousLesson.click();
    }
}

// 选择下一课程
function selectNextLesson() {
    const currentLesson = document.querySelector('.lesson-item.active');
    const nextLesson = currentLesson?.nextElementSibling;
    
    if (nextLesson && nextLesson.classList.contains('lesson-item')) {
        nextLesson.click();
    }
}

// 更新URL参数
function updateURLParams(lessonId) {
    const url = new URL(window.location);
    url.searchParams.set('lesson_id', lessonId);
    window.history.replaceState({}, '', url);
}

// 获取URL参数
function getQueryParam(param) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(param);
}

// 更新视频控制UI
function updateVideoControls() {
    // 初始化时间显示
    updateTimeDisplay();
    
    // 初始化其他控制元素的事件
    const volumeBtn = document.querySelector('.control-btn:nth-child(2)');
    if (volumeBtn) {
        volumeBtn.addEventListener('click', toggleMute);
    }
    
    const fullscreenBtn = document.querySelector('.controls-right .control-btn:last-child');
    if (fullscreenBtn) {
        fullscreenBtn.addEventListener('click', toggleFullscreen);
    }
}

// 静音切换
function toggleMute() {
    const volumeBtn = document.querySelector('.control-btn:nth-child(2) i');
    if (volumeBtn) {
        if (volumeBtn.classList.contains('fa-volume-up')) {
            volumeBtn.className = 'fas fa-volume-mute';
            showNotification('已静音', 'info');
        } else {
            volumeBtn.className = 'fas fa-volume-up';
            showNotification('已取消静音', 'info');
        }
    }
}

// 全屏切换
function toggleFullscreen() {
    const videoPlayerSection = document.querySelector('.video-player-section');
    
    if (document.fullscreenElement) {
        document.exitFullscreen();
    } else if (videoPlayerSection) {
        videoPlayerSection.requestFullscreen().catch(console.error);
    }
}

// 显示通知消息
function showNotification(message, type = 'info') {
    // 创建通知元素
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.innerHTML = `
        <div class="notification-content">
            <i class="fas ${getNotificationIcon(type)}"></i>
            <span>${message}</span>
        </div>
    `;
    
    // 添加样式
    notification.style.cssText = `
        position: fixed;
        top: 90px;
        right: 20px;
        background-color: var(--bg-secondary);
        color: var(--text-primary);
        padding: 12px 16px;
        border-radius: 8px;
        border: 1px solid var(--border-color);
        box-shadow: var(--shadow-lg);
        z-index: 10000;
        transform: translateX(100%);
        transition: transform 0.3s ease;
        max-width: 300px;
    `;
    
    // 添加类型特定样式
    if (type === 'success') {
        notification.style.borderLeftColor = '#4CAF50';
    } else if (type === 'error') {
        notification.style.borderLeftColor = '#f44336';
    } else if (type === 'warning') {
        notification.style.borderLeftColor = '#ff9800';
    }
    
    document.body.appendChild(notification);
    
    // 显示动画
    setTimeout(() => {
        notification.style.transform = 'translateX(0)';
    }, 100);
    
    // 自动隐藏
    setTimeout(() => {
        notification.style.transform = 'translateX(100%)';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// 获取通知图标
function getNotificationIcon(type) {
    switch (type) {
        case 'success': return 'fa-check-circle';
        case 'error': return 'fa-exclamation-circle';
        case 'warning': return 'fa-exclamation-triangle';
        default: return 'fa-info-circle';
    }
}

// 学习行为追踪
function trackLessonSelection(lessonId, courseId) {
    console.log('学习行为追踪:', {
        action: 'lesson_selected',
        courseId: courseId,
        lessonId: lessonId,
        timestamp: new Date().toISOString()
    });
    
    // 这里可以添加实际的分析追踪代码
}

// 视频播放追踪（从main.js继承）
function trackVideoPlay(context, videoId) {
    console.log('视频播放追踪:', {
        action: 'video_play',
        context: context,
        videoId: videoId,
        timestamp: new Date().toISOString()
    });
}

// 文档处理功能
function downloadDocument(lessonId) {
    showNotification('开始下载文档...', 'info');
    
    // 模拟下载过程
    setTimeout(() => {
        showNotification('文档下载完成！', 'success');
        console.log('下载文档:', lessonId);
    }, 1500);
}

function previewDocument(lessonId) {
    showNotification('正在加载文档预览...', 'info');
    
    // 模拟预览加载
    setTimeout(() => {
        showNotification('文档预览已打开', 'success');
        console.log('预览文档:', lessonId);
        // 这里可以打开模态框或新窗口显示文档预览
    }, 1000);
}

// 动态更新主操作按钮
function updateMainCtaButton() {
    const mainBtn = document.getElementById('mainCtaBtn');
    const ctaText = mainBtn?.querySelector('.cta-text');
    const ctaIcon = mainBtn?.querySelector('i');
    
    if (!mainBtn) return;
    
    // 模拟用户状态检查
    const hasStartedCourse = localStorage.getItem('course_progress_' + courseId);
    const isEnrolled = localStorage.getItem('enrolled_course_' + courseId);
    
    if (hasStartedCourse) {
        ctaIcon.className = 'fas fa-play';
        ctaText.textContent = '继续学习';
        mainBtn.style.background = 'linear-gradient(135deg, var(--accent-primary) 0%, var(--accent-secondary) 100%)';
    } else if (isEnrolled) {
        ctaIcon.className = 'fas fa-play-circle';
        ctaText.textContent = '开始学习';
        mainBtn.style.background = 'linear-gradient(135deg, #10b981 0%, #059669 100%)';
    } else {
        ctaIcon.className = 'fas fa-shopping-cart';
        ctaText.textContent = '立即加入';
        mainBtn.style.background = 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)';
    }
}

// 处理主操作按钮点击
function handleMainCtaClick() {
    const mainBtn = document.getElementById('mainCtaBtn');
    const ctaText = mainBtn?.querySelector('.cta-text')?.textContent;
    
    switch (ctaText) {
        case '继续学习':
        case '开始学习':
            // 开始学习当前课程
            if (currentLessonId) {
                togglePlayPause();
                localStorage.setItem('course_progress_' + courseId, currentLessonId);
            }
            break;
        case '立即加入':
            // 显示购买/注册流程
            showEnrollmentModal();
            break;
    }
}

// 显示注册模态框
function showEnrollmentModal() {
    showNotification('正在跳转到课程注册页面...', 'info');
    // 模拟注册流程
    setTimeout(() => {
        localStorage.setItem('enrolled_course_' + courseId, 'true');
        updateMainCtaButton();
        showNotification('恭喜！您已成功加入课程', 'success');
    }, 2000);
}

// 课程类型检测和UI更新
function updatePlayerInterface(lesson) {
    const videoPlayer = document.getElementById('videoPlayer');
    const isVideo = lesson.FileName && (
        lesson.FileName.endsWith('.mp4') || 
        lesson.FileName.endsWith('.avi') || 
        lesson.FileName.endsWith('.mov')
    );
    
    // 更新播放器界面类型
    if (isVideo) {
        videoPlayer.classList.remove('document-mode');
        videoPlayer.classList.add('video-mode');
    } else {
        videoPlayer.classList.remove('video-mode');
        videoPlayer.classList.add('document-mode');
    }
}

// 动态更新播放器内容
function updatePlayerContent(lesson) {
    const videoPlayer = document.getElementById('videoPlayer');
    const isVideo = lesson.FileName && (
        lesson.FileName.endsWith('.mp4') || 
        lesson.FileName.endsWith('.avi') || 
        lesson.FileName.endsWith('.mov')
    );
    
    // 清空当前内容
    videoPlayer.innerHTML = '';
    
    if (isVideo) {
        // 创建视频播放器界面
        videoPlayer.innerHTML = `
            <div class="video-placeholder">
                <div class="video-content">
                    <div class="play-button-container">
                        <button class="video-play-btn" onclick="trackVideoPlay('course-detail', '${courseId}-${lesson.Id}')">
                            <i class="fas fa-play"></i>
                        </button>
                    </div>
                    <div class="video-info">
                        <h3 class="current-lesson-title">${lesson.FileName}</h3>
                        <p class="current-lesson-meta">课程: Go语言进阶开发 • ${getEstimatedDuration(lesson.FileName)}</p>
                    </div>
                </div>
                <div class="video-overlay"></div>
            </div>
        `;
        
        // 重新绑定播放按钮事件
        const playBtn = videoPlayer.querySelector('.video-play-btn');
        if (playBtn) {
            playBtn.addEventListener('click', togglePlayPause);
        }
    } else {
        // 创建文档预览界面
        const fileType = getFileTypeFromName(lesson.FileName);
        const iconClass = getDocumentIconClass(fileType);
        
        videoPlayer.innerHTML = `
            <div class="document-preview">
                <div class="document-content">
                    <div class="document-icon-container">
                        <i class="${iconClass} document-icon"></i>
                    </div>
                    <div class="document-info">
                        <h3 class="current-lesson-title">${lesson.FileName}</h3>
                        <p class="current-lesson-meta">${fileType} • ${getEstimatedFileSize(lesson.FileName)}</p>
                        <div class="document-actions">
                            <button class="btn-primary document-btn" onclick="downloadDocument('${lesson.Id}')">
                                <i class="fas fa-download"></i>
                                下载资料
                            </button>
                            <button class="btn-secondary document-btn" onclick="previewDocument('${lesson.Id}')">
                                <i class="fas fa-eye"></i>
                                在线预览
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }
}

// 获取文件类型
function getFileTypeFromName(filename) {
    if (!filename) return "文件";
    
    if (filename.endsWith('.mp4') || filename.endsWith('.avi') || filename.endsWith('.mov')) {
        return "视频";
    } else if (filename.endsWith('.pdf')) {
        return "PDF";
    } else if (filename.endsWith('.ppt') || filename.endsWith('.pptx')) {
        return "演示文稿";
    }
    return "文档";
}

// 获取文档图标类名
function getDocumentIconClass(fileType) {
    switch (fileType) {
        case 'PDF':
            return 'fas fa-file-pdf';
        case '演示文稿':
            return 'fas fa-file-powerpoint';
        default:
            return 'fas fa-file-alt';
    }
}

// 获取预估播放时长
function getEstimatedDuration(filename) {
    const durations = ["8:45", "12:30", "15:20", "9:15", "18:40", "11:25", "7:30", "14:55"];
    let hash = 0;
    for (let i = 0; i < filename.length; i++) {
        hash += filename.charCodeAt(i);
    }
    return durations[hash % durations.length];
}

// 获取预估文件大小
function getEstimatedFileSize(filename) {
    const sizes = ["12.5 MB", "3.2 MB", "8.7 MB", "5.1 MB", "15.3 MB", "2.8 MB"];
    let hash = 0;
    for (let i = 0; i < filename.length; i++) {
        hash += filename.charCodeAt(i);
    }
    return sizes[hash % sizes.length];
}

// 增强的课程选择功能
function selectLesson(lessonId, lessonTitle, courseIdParam) {
    console.log('选择课程:', { lessonId, lessonTitle, courseIdParam });
    
    // 更新当前课程信息
    currentLessonId = lessonId;
    courseId = courseIdParam;
    
    // 获取课程对象
    const lessonItem = document.querySelector(`[data-lesson-id="${lessonId}"]`);
    const lesson = {
        Id: lessonId,
        FileName: lessonTitle
    };
    
    // 更新UI
    updateActiveLesson(lessonId);
    updateVideoPlayer(lessonTitle);
    updatePlayerInterface(lesson);
    
    // 动态更新视频播放器内容
    updatePlayerContent(lesson);
    
    // 记录学习行为
    trackLessonSelection(lessonId, courseIdParam);
    
    // 更新URL参数（不跳转页面）
    updateURLParams(lessonId);
}

// 相关推荐课程信息显示
function showRelatedCourseInfo(courseName, instructor) {
    showNotification(`${courseName} - 讲师：${instructor}`, 'info');
    console.log('查看相关课程:', courseName, instructor);
}

// 导出函数供HTML模板使用
window.selectLesson = selectLesson;
window.trackVideoPlay = trackVideoPlay;
window.downloadDocument = downloadDocument;
window.previewDocument = previewDocument;
window.showRelatedCourseInfo = showRelatedCourseInfo; 