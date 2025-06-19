// ===== 主要JavaScript功能 =====

class CoursePlayer {
    constructor() {
        this.currentUser = null;
        this.currentCourse = null;
        this.isPlaying = false;
        this.init();
    }

    init() {
        this.bindEvents();
        this.checkAuthStatus();
        this.loadCourseData();
        this.initializeAnimations();
    }

    // ===== 事件绑定 =====
    bindEvents() {
        // 导航事件
        this.bindNavigationEvents();
        
        // 课程卡片事件
        this.bindCourseCardEvents();
        
        // 视频控制事件
        this.bindVideoControlEvents();
        
        // 用户配置文件事件
        this.bindUserProfileEvents();

        // 响应式事件
        this.bindResponsiveEvents();
    }

    bindNavigationEvents() {
        // 导航链接点击 - 使用标准的链接跳转行为
        document.querySelectorAll('.nav-link').forEach(link => {
            link.addEventListener('click', (e) => {
                // 获取实际的链接元素，可能是事件目标的父元素
                const actualLink = e.target.closest('.nav-link') || link;
                const href = actualLink.getAttribute('href');
                
                if (href && href !== '#') {
                    // 允许正常的页面跳转
                    console.log(`Navigating to: ${href}`);
                    // 不使用preventDefault()，让浏览器执行正常的导航
                } else {
                    // 只对无效链接阻止默认行为
                    e.preventDefault();
                    const linkText = actualLink.textContent.trim();
                    console.log(`Handling navigation for: ${linkText}`);
                    this.handleNavigation(linkText);
                }
            });
        });

        // Logo点击返回首页 - 也使用标准链接行为
        document.querySelector('.logo')?.addEventListener('click', (e) => {
            // 由于logo现在是<a>标签，这个事件处理器实际上可能不需要了
            // 但保留它以防万一有其他情况
            const logoLink = e.target.closest('.logo');
            const href = logoLink?.getAttribute('href');
            
            if (href && href !== '#') {
                // 允许正常的页面跳转
                console.log(`Navigating to home: ${href}`);
                // 不使用preventDefault()，让浏览器执行正常的导航
            } else {
                e.preventDefault();
                this.navigateToHome();
            }
        });
    }

    bindCourseCardEvents() {
        // 课程卡片点击播放
        document.querySelectorAll('.course-card').forEach(card => {
            card.addEventListener('click', (e) => {
                e.preventDefault();
                this.handleCourseCardClick(card);
            });

            // 悬停效果
            card.addEventListener('mouseenter', () => {
                this.handleCourseCardHover(card, true);
            });

            card.addEventListener('mouseleave', () => {
                this.handleCourseCardHover(card, false);
            });
        });

        // Watch Now 按钮
        document.querySelector('.btn-primary')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.handleWatchNowClick();
        });
    }

    bindVideoControlEvents() {
        // 视频控制按钮
        document.querySelectorAll('.control-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                this.handleVideoControl(btn);
            });
        });

        // 视频容器点击播放/暂停
        document.querySelector('.video-container')?.addEventListener('click', (e) => {
            if (e.target.classList.contains('video-container') || 
                e.target.classList.contains('video-placeholder')) {
                this.toggleVideoPlayback();
            }
        });
    }

    bindUserProfileEvents() {
        // 用户配置文件下拉菜单
        const userProfile = document.getElementById('userProfile');
        const logoutBtn = document.getElementById('logoutBtn');
        const profileSettings = document.getElementById('profileSettings');
        const orderHistory = document.getElementById('orderHistory');
        const helpCenter = document.getElementById('helpCenter');
        
        if (userProfile) {
            userProfile.addEventListener('click', (e) => {
                e.stopPropagation();
                userProfile.classList.toggle('active');
            });
            
            // 点击其他地方关闭下拉菜单
            document.addEventListener('click', () => {
                userProfile.classList.remove('active');
            });
        }
        
        if (logoutBtn) {
            logoutBtn.addEventListener('click', (e) => {
                e.preventDefault();
                userProfile.classList.remove('active');
                this.logout();
            });
        }
        
        if (profileSettings) {
            profileSettings.addEventListener('click', (e) => {
                e.preventDefault();
                userProfile.classList.remove('active');
                this.showProfileSettings();
            });
        }
        
        if (orderHistory) {
            orderHistory.addEventListener('click', (e) => {
                e.preventDefault();
                userProfile.classList.remove('active');
                this.showOrderHistory();
            });
        }
        
        if (helpCenter) {
            helpCenter.addEventListener('click', (e) => {
                e.preventDefault();
                userProfile.classList.remove('active');
                this.showHelpCenter();
            });
        }
    }

    bindResponsiveEvents() {
        // 窗口大小改变
        let resizeTimeout;
        window.addEventListener('resize', () => {
            clearTimeout(resizeTimeout);
            resizeTimeout = setTimeout(() => {
                this.handleResize();
            }, 250);
        });

        // 滚动事件
        let scrollTimeout;
        window.addEventListener('scroll', () => {
            clearTimeout(scrollTimeout);
            scrollTimeout = setTimeout(() => {
                this.handleScroll();
            }, 10);
        });
    }

    // ===== 导航处理 =====
    handleNavigation(section) {
        console.log(`Navigating to: ${section}`);
        
        // 更新活跃状态
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
        });
        
        // 根据不同section执行不同逻辑
        switch (section) {
            case 'Discover':
                this.showDiscoverContent();
                break;
            case 'My Progress':
                this.showProgressContent();
                break;
            case 'Library':
                this.showLibraryContent();
                break;
            default:
                console.log('Unknown navigation section');
        }
    }

    navigateToHome() {
        window.location.href = '/';
    }

    // ===== 课程处理 =====
    handleCourseCardClick(card) {
        const courseId = card.getAttribute('data-course-id');
        const courseTitle = card.querySelector('.course-title')?.textContent;
        
        console.log(`Navigating to course: ${courseTitle} (ID: ${courseId})`);
        
        if (courseId) {
            // 跳转到课程详情页面
            window.location.href = `/course/${courseId}`;
        } else {
            console.warn('Course ID not found');
        }
    }

    handleCourseCardHover(card, isHovering) {
        const playOverlay = card.querySelector('.play-overlay');
        if (playOverlay) {
            if (isHovering) {
                playOverlay.style.transform = 'translate(-50%, -50%) scale(1.1)';
                playOverlay.style.opacity = '1';
            } else {
                playOverlay.style.transform = 'translate(-50%, -50%) scale(1)';
                playOverlay.style.opacity = '0';
            }
        }
    }

    handleWatchNowClick() {
        console.log('Watch Now clicked');
        this.playCourse('French Pastry Fundamentals', 'Dominique Ansel');
    }

    playCourse(title, instructor) {
        // 更新主视频区域
        const videoTitle = document.querySelector('.video-title');
        if (videoTitle) {
            videoTitle.textContent = title;
        }

        // 显示播放状态
        this.isPlaying = true;
        this.updateVideoControls();
        
        // 这里可以集成真实的视频播放器
        console.log(`Now playing: ${title} by ${instructor}`);
        
        // 显示通知
        this.showNotification(`开始播放: ${title}`, 'success');
    }

    // ===== 视频控制 =====
    handleVideoControl(btn) {
        const icon = btn.querySelector('i');
        if (!icon) return;

        const action = this.getVideoAction(icon.className);
        console.log(`Video control: ${action}`);

        switch (action) {
            case 'play-pause':
                this.toggleVideoPlayback();
                break;
            case 'volume':
                this.toggleVolume();
                break;
            case 'captions':
                this.toggleCaptions();
                break;
            case 'previous':
                this.previousVideo();
                break;
            case 'next':
                this.nextVideo();
                break;
        }
    }

    getVideoAction(iconClass) {
        if (iconClass.includes('play') || iconClass.includes('pause')) return 'play-pause';
        if (iconClass.includes('volume')) return 'volume';
        if (iconClass.includes('caption')) return 'captions';
        if (iconClass.includes('chevron-left')) return 'previous';
        if (iconClass.includes('chevron-right')) return 'next';
        return 'unknown';
    }

    toggleVideoPlayback() {
        this.isPlaying = !this.isPlaying;
        this.updateVideoControls();
        
        const action = this.isPlaying ? '播放' : '暂停';
        console.log(`Video ${action}`);
        
        // 这里集成真实视频播放器的播放/暂停逻辑
    }

    updateVideoControls() {
        const playPauseBtn = document.querySelector('.control-btn i.fas.fa-pause, .control-btn i.fas.fa-play');
        if (playPauseBtn) {
            playPauseBtn.className = this.isPlaying ? 'fas fa-pause' : 'fas fa-play';
        }
    }

    toggleVolume() {
        const volumeBtn = document.querySelector('.control-btn i[class*="volume"]');
        if (volumeBtn) {
            const isMuted = volumeBtn.classList.contains('fa-volume-mute');
            volumeBtn.className = isMuted ? 'fas fa-volume-up' : 'fas fa-volume-mute';
            console.log(isMuted ? 'Volume unmuted' : 'Volume muted');
        }
    }

    toggleCaptions() {
        console.log('Captions toggled');
        this.showNotification('字幕切换', 'info');
    }

    previousVideo() {
        console.log('Previous video');
        this.showNotification('上一个视频', 'info');
    }

    nextVideo() {
        console.log('Next video');
        this.showNotification('下一个视频', 'info');
    }

    // ===== 用户菜单 =====
    toggleUserMenu() {
        console.log('User menu toggled');
        
        // 检查是否已登录
        if (!this.currentUser) {
            this.showLoginModal();
        } else {
            this.showUserDropdown();
        }
    }

    showLoginModal() {
        // 创建登录模态框
        const modal = this.createModal('login');
        document.body.appendChild(modal);
        
        // 延迟显示以触发动画
        setTimeout(() => {
            modal.classList.add('show');
        }, 10);
    }

    showUserDropdown() {
        console.log('Showing user dropdown');
        // 用户下拉菜单逻辑
    }

    // ===== 内容显示 =====
    showDiscoverContent() {
        console.log('Showing Discover content');
        // 显示发现页面内容
    }

    showProgressContent() {
        console.log('Showing Progress content');
        // 显示学习进度
    }

    showLibraryContent() {
        console.log('Showing Library content');
        // 显示课程库
    }

    // ===== 数据加载 =====
    async loadCourseData() {
        try {
            // 模拟API调用
            console.log('Loading course data...');
            
            // 这里调用实际的API
            // const response = await fetch('/api/v1/courses');
            // const courses = await response.json();
            
            // 模拟延迟
            setTimeout(() => {
                console.log('Course data loaded');
                this.animateContent();
            }, 1000);
            
        } catch (error) {
            console.error('Error loading course data:', error);
            this.showNotification('加载课程数据失败', 'error');
        }
    }

    checkAuthStatus() {
        // 检查用户认证状态
        const token = this.getAuthToken();
        if (token) {
            this.validateToken(token);
        } else {
            // 如果没有token，显示登录/注册按钮
            this.showAuthButtons();
        }
    }

    async validateToken(token) {
        try {
            console.log('🔍 验证用户认证状态...');
            
            // 首先尝试从本地存储获取用户信息
            const storedUserInfo = this.getStoredUserInfo();
            if (storedUserInfo) {
                this.currentUser = storedUserInfo;
                this.updateUserInterface();
                console.log('✅ 使用本地存储的用户信息');
            }
            
            // 验证token有效性
            const response = await fetch('/api/v1/validate-token', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ token })
            });
            
            if (response.ok) {
                // 如果没有本地用户信息，获取最新用户信息
                if (!this.currentUser) {
                    await this.fetchUserInfo(token);
                }
                console.log('✅ Token验证成功');
            } else {
                throw new Error('Token验证失败');
            }
            
        } catch (error) {
            console.error('❌ Token验证失败:', error);
            this.clearAuthData();
            this.showAuthButtons();
        }
    }

    async fetchUserInfo(token) {
        try {
            const response = await fetch('/api/v1/me', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });
            
            if (response.ok) {
                const data = await response.json();
                this.currentUser = data.user;
                this.updateUserInterface();
                console.log('✅ 获取用户信息成功:', this.currentUser);
            }
            
        } catch (error) {
            console.error('❌ 获取用户信息失败:', error);
        }
    }

    updateUserInterface() {
        if (!this.currentUser) return;
        
        // 显示用户头像和隐藏登录按钮
        const userProfile = document.getElementById('userProfile');
        const authButtons = document.getElementById('authButtons');
        const userName = document.getElementById('userName');
        
        if (userProfile) {
            userProfile.style.display = 'flex';
        }
        if (authButtons) {
            authButtons.style.display = 'none';
        }
        if (userName) {
            userName.textContent = this.currentUser.nickname || this.currentUser.username || '用户';
        }
        
        console.log('✅ 用户界面更新完成');
    }

    showAuthButtons() {
        // 显示登录/注册按钮，隐藏用户头像
        const userProfile = document.getElementById('userProfile');
        const authButtons = document.getElementById('authButtons');
        
        if (userProfile) {
            userProfile.style.display = 'none';
        }
        if (authButtons) {
            authButtons.style.display = 'flex';
        }
    }

    // 工具方法
    getAuthToken() {
        return localStorage.getItem('authToken') || sessionStorage.getItem('authToken');
    }

    getStoredUserInfo() {
        const userInfo = localStorage.getItem('userInfo') || sessionStorage.getItem('userInfo');
        return userInfo ? JSON.parse(userInfo) : null;
    }

    clearAuthData() {
        localStorage.removeItem('authToken');
        localStorage.removeItem('tokenExpiry');
        localStorage.removeItem('userInfo');
        sessionStorage.removeItem('authToken');
        sessionStorage.removeItem('tokenExpiry');
        sessionStorage.removeItem('userInfo');
    }

    logout() {
        console.log('👋 用户退出登录');
        
        // 清除认证数据
        this.clearAuthData();
        
        // 重置用户状态
        this.currentUser = null;
        
        // 更新界面
        this.showAuthButtons();
        
        // 显示通知
        this.showNotification('您已成功退出登录', 'success');
        
        // 可选：重定向到首页
        setTimeout(() => {
            window.location.reload();
        }, 1500);
    }

    showProfileSettings() {
        console.log('👤 显示个人资料设置');
        
        // 显示个人资料设置模态框或跳转到设置页面
        this.showNotification('个人资料设置功能开发中...', 'info');
        
        // 未来可以实现：
        // window.location.href = '/profile/settings';
    }

    showOrderHistory() {
        console.log('🧾 显示订单历史');
        
        // 显示订单历史页面或模态框
        this.showNotification('订单历史功能开发中...', 'info');
        
        // 未来可以实现：
        // window.location.href = '/orders';
    }

    showHelpCenter() {
        console.log('❓ 显示帮助中心');
        
        // 显示帮助中心页面或模态框
        this.showNotification('帮助中心功能开发中...', 'info');
        
        // 未来可以实现：
        // window.location.href = '/help';
        // 或者打开在线客服
    }

    // ===== 动画和UI效果 =====
    initializeAnimations() {
        // 初始化页面动画
        this.animateOnScroll();
        this.initParallaxEffects();
    }

    animateContent() {
        // 为课程卡片添加进入动画
        const cards = document.querySelectorAll('.course-card');
        cards.forEach((card, index) => {
            setTimeout(() => {
                card.style.opacity = '1';
                card.style.transform = 'translateY(0)';
            }, index * 100);
        });
    }

    animateOnScroll() {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('animate-in');
                }
            });
        }, { threshold: 0.1 });

        document.querySelectorAll('.course-card, .section-title').forEach(el => {
            observer.observe(el);
        });
    }

    initParallaxEffects() {
        // 轻微的视差效果
        window.addEventListener('scroll', () => {
            const scrolled = window.pageYOffset;
            const heroSection = document.querySelector('.featured-course');
            
            if (heroSection) {
                heroSection.style.transform = `translateY(${scrolled * 0.1}px)`;
            }
        });
    }

    handleResize() {
        // 响应式处理
        const width = window.innerWidth;
        
        if (width < 768) {
            this.enableMobileOptimizations();
        } else {
            this.disableMobileOptimizations();
        }
    }

    handleScroll() {
        const scrolled = window.pageYOffset;
        const navbar = document.querySelector('.navbar');
        
        // 导航栏透明度效果
        if (navbar) {
            const opacity = Math.min(scrolled / 100, 1);
            navbar.style.backgroundColor = `rgba(10, 10, 10, ${0.95 * opacity})`;
        }
    }

    enableMobileOptimizations() {
        // 移动端优化
        document.body.classList.add('mobile-optimized');
    }

    disableMobileOptimizations() {
        // 桌面端优化
        document.body.classList.remove('mobile-optimized');
    }

    // ===== 工具函数 =====
    createModal(type) {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        
        const modalContent = document.createElement('div');
        modalContent.className = 'modal-content';
        
        if (type === 'login') {
            modalContent.innerHTML = `
                <div class="modal-header">
                    <h2>登录</h2>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form class="login-form">
                        <input type="text" placeholder="用户名或邮箱" required>
                        <input type="password" placeholder="密码" required>
                        <button type="submit" class="btn-primary">登录</button>
                    </form>
                    <p>还没有账户？<a href="#" class="switch-to-register">注册</a></p>
                </div>
            `;
        }
        
        modal.appendChild(modalContent);
        
        // 绑定关闭事件
        modal.addEventListener('click', (e) => {
            if (e.target === modal || e.target.classList.contains('modal-close')) {
                this.closeModal(modal);
            }
        });
        
        return modal;
    }

    closeModal(modal) {
        modal.classList.remove('show');
        setTimeout(() => {
            document.body.removeChild(modal);
        }, 300);
    }

    showNotification(message, type = 'info') {
        // 创建通知
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.textContent = message;
        
        // 添加到页面
        document.body.appendChild(notification);
        
        // 显示动画
        setTimeout(() => notification.classList.add('show'), 10);
        
        // 自动消失
        setTimeout(() => {
            notification.classList.remove('show');
            setTimeout(() => {
                if (document.body.contains(notification)) {
                    document.body.removeChild(notification);
                }
            }, 300);
        }, 3000);
    }

    // ===== 键盘快捷键 =====
    initKeyboardShortcuts() {
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey || e.metaKey) {
                switch (e.key) {
                    case 'k':
                        e.preventDefault();
                        this.focusSearch();
                        break;
                    case 'p':
                        e.preventDefault();
                        this.toggleVideoPlayback();
                        break;
                }
            }
            
            switch (e.key) {
                case 'Escape':
                    this.closeAllModals();
                    break;
                case ' ':
                    if (e.target.tagName !== 'INPUT' && e.target.tagName !== 'TEXTAREA') {
                        e.preventDefault();
                        this.toggleVideoPlayback();
                    }
                    break;
            }
        });
    }

    focusSearch() {
        const searchInput = document.querySelector('input[type="search"]');
        if (searchInput) {
            searchInput.focus();
        }
    }

    closeAllModals() {
        const modals = document.querySelectorAll('.modal-overlay');
        modals.forEach(modal => this.closeModal(modal));
    }
}

// ===== 初始化 =====
document.addEventListener('DOMContentLoaded', () => {
    const app = new CoursePlayer();
    
    // 全局错误处理
    window.addEventListener('error', (e) => {
        console.error('Global error:', e.error);
    });
    
    // 注册Service Worker（如果需要离线功能）
    if ('serviceWorker' in navigator) {
        navigator.serviceWorker.register('/sw.js')
            .then(registration => {
                console.log('SW registered:', registration);
            })
            .catch(error => {
                console.log('SW registration failed:', error);
            });
    }
});

// ===== CSS动态注入（用于模态框等动态元素）=====
const additionalStyles = `
    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 2000;
        opacity: 0;
        transition: opacity 0.3s ease;
    }
    
    .modal-overlay.show {
        opacity: 1;
    }
    
    .modal-content {
        background: var(--bg-secondary);
        border-radius: 12px;
        padding: 2rem;
        max-width: 400px;
        width: 90%;
        transform: translateY(-20px);
        transition: transform 0.3s ease;
    }
    
    .modal-overlay.show .modal-content {
        transform: translateY(0);
    }
    
    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
    }
    
    .modal-close {
        background: none;
        border: none;
        font-size: 1.5rem;
        color: var(--text-secondary);
        cursor: pointer;
    }
    
    .login-form input {
        width: 100%;
        padding: 0.75rem;
        margin-bottom: 1rem;
        border: 1px solid var(--border-color);
        border-radius: 8px;
        background: var(--bg-tertiary);
        color: var(--text-primary);
        font-family: inherit;
    }
    
    .notification {
        position: fixed;
        top: 100px;
        right: 20px;
        padding: 1rem 1.5rem;
        border-radius: 8px;
        color: white;
        font-weight: 500;
        transform: translateX(100%);
        transition: transform 0.3s ease;
        z-index: 3000;
    }
    
    .notification.show {
        transform: translateX(0);
    }
    
    .notification-success {
        background: #10b981;
    }
    
    .notification-error {
        background: #ef4444;
    }
    
    .notification-info {
        background: #3b82f6;
    }
    
    .animate-in {
        animation: fadeInUp 0.6s ease forwards;
    }
    
    @media (max-width: 768px) {
        .modal-content {
            margin: 1rem;
            padding: 1.5rem;
        }
        
        .notification {
            right: 10px;
            left: 10px;
            transform: translateY(-100%);
        }
        
        .notification.show {
            transform: translateY(0);
        }
    }
`;

// 注入样式
const styleSheet = document.createElement('style');
styleSheet.textContent = additionalStyles;
document.head.appendChild(styleSheet);

// ===== 播放统计功能 =====
function trackVideoPlay(section, videoId) {
    const playData = {
        section: section,
        videoId: videoId,
        timestamp: new Date().toISOString(),
        userAgent: navigator.userAgent,
        sessionId: generateSessionId()
    };
    
    console.log('📊 Video play tracked:', playData);
    
    // 发送统计数据到后端
    sendPlayAnalytics(playData);
    
    // 显示播放提示
    showPlayFeedback(section, videoId);
}

function generateSessionId() {
    return 'session_' + Math.random().toString(36).substr(2, 9) + '_' + Date.now();
}

async function sendPlayAnalytics(data) {
    try {
        const response = await fetch('/api/v1/analytics/play', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            console.log('✅ Play analytics sent successfully');
        } else {
            console.log('📊 Analytics endpoint not ready, data logged locally');
        }
    } catch (error) {
        console.log('📊 Analytics will be implemented later, data stored locally:', data);
        // 存储到localStorage作为备用
        storeAnalyticsLocally(data);
    }
}

function storeAnalyticsLocally(data) {
    try {
        const stored = JSON.parse(localStorage.getItem('playAnalytics') || '[]');
        stored.push(data);
        // 只保留最近100条记录
        if (stored.length > 100) {
            stored.splice(0, stored.length - 100);
        }
        localStorage.setItem('playAnalytics', JSON.stringify(stored));
    } catch (error) {
        console.log('Unable to store analytics locally');
    }
}

function showPlayFeedback(section, videoId) {
    // 创建播放反馈提示
    const feedback = document.createElement('div');
    feedback.className = 'play-feedback';
    feedback.innerHTML = `
        <div class="feedback-content">
            <i class="fas fa-play-circle"></i>
            <span>播放记录已统计</span>
        </div>
    `;
    
    // 添加样式
    feedback.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: rgba(229, 9, 20, 0.9);
        color: white;
        padding: 12px 20px;
        border-radius: 8px;
        font-size: 14px;
        font-weight: 500;
        z-index: 10000;
        animation: slideInRight 0.3s ease-out;
        backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.2);
    `;
    
    // 添加到页面
    document.body.appendChild(feedback);
    
    // 3秒后自动移除
    setTimeout(() => {
        feedback.style.animation = 'slideOutRight 0.3s ease-in forwards';
        setTimeout(() => {
            if (feedback.parentNode) {
                feedback.parentNode.removeChild(feedback);
            }
        }, 300);
    }, 3000);
}

// 为播放反馈添加动画样式
const feedbackStyles = `
    @keyframes slideInRight {
        from { transform: translateX(100%); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
    }
    @keyframes slideOutRight {
        from { transform: translateX(0); opacity: 1; }
        to { transform: translateX(100%); opacity: 0; }
    }
    .feedback-content {
        display: flex;
        align-items: center;
        gap: 8px;
    }
`;

const feedbackStyleSheet = document.createElement('style');
feedbackStyleSheet.textContent = feedbackStyles;
document.head.appendChild(feedbackStyleSheet); 