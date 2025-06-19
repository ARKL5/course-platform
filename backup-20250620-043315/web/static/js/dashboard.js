// ===== 用户中心主要功能 =====
class UserDashboard {
    constructor() {
        this.currentUser = null;
        this.currentSection = 'learning-progress';
        this.init();
    }

    init() {
        console.log('🎯 初始化用户中心...');
        this.checkAuthStatus();
        this.bindEvents();
        this.initializeAnimations();
        this.loadDefaultSection();
    }

    // ===== 用户认证检查 =====
    checkAuthStatus() {
        const token = this.getAuthToken();
        if (!token) {
            console.log('❌ 未登录，重定向到登录页面');
            window.location.href = '/login';
            return;
        }

        this.validateToken(token);
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
            window.location.href = '/login';
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
                
                // 更新本地存储
                localStorage.setItem('userInfo', JSON.stringify(this.currentUser));
                sessionStorage.setItem('userInfo', JSON.stringify(this.currentUser));
                
                this.updateUserInterface();
                console.log('✅ 获取用户信息成功:', this.currentUser);
            }
            
        } catch (error) {
            console.error('❌ 获取用户信息失败:', error);
        }
    }

    updateUserInterface() {
        if (!this.currentUser) return;
        
        // 更新导航栏用户信息
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

        // 更新个人资料表单
        this.loadProfileData();
        
        console.log('✅ 用户界面更新完成');
    }

    // ===== 事件绑定 =====
    bindEvents() {
        // 左侧导航菜单事件
        this.bindSidebarNavEvents();
        
        // 用户下拉菜单事件
        this.bindUserProfileEvents();
        
        // 表单事件
        this.bindFormEvents();
        
        // 移动端菜单事件
        this.bindMobileEvents();
        
        // 课程过滤事件
        this.bindCourseFilterEvents();

        // 安全设置事件
        this.bindSecurityEvents();
    }

    bindSidebarNavEvents() {
        const navLinks = document.querySelectorAll('.sidebar .nav-link');
        
        navLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const section = link.getAttribute('data-section');
                this.switchSection(section);
            });
        });
    }

    bindUserProfileEvents() {
        const userProfile = document.getElementById('userProfile');
        const logoutBtn = document.getElementById('logoutBtn');
        
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
    }

    bindFormEvents() {
        // 个人资料表单
        const profileForm = document.getElementById('profileForm');
        const saveProfileBtn = document.getElementById('saveProfile');
        const cancelProfileBtn = document.getElementById('cancelProfile');
        
        if (profileForm) {
            profileForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.saveProfile();
            });
        }
        
        if (cancelProfileBtn) {
            cancelProfileBtn.addEventListener('click', () => {
                this.loadProfileData();
                this.showNotification('已取消修改', 'info');
            });
        }

        // 头像上传
        const avatarChangeBtn = document.getElementById('avatarChangeBtn');
        const avatarFileInput = document.getElementById('avatarFileInput');
        
        if (avatarChangeBtn && avatarFileInput) {
            avatarChangeBtn.addEventListener('click', () => {
                avatarFileInput.click();
            });
            
            avatarFileInput.addEventListener('change', (e) => {
                this.handleAvatarUpload(e);
            });
        }

        // 密码修改表单
        const passwordForm = document.getElementById('passwordForm');
        const cancelPasswordBtn = document.getElementById('cancelPassword');
        
        if (passwordForm) {
            passwordForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.changePassword();
            });
        }
        
        if (cancelPasswordBtn) {
            cancelPasswordBtn.addEventListener('click', () => {
                this.resetPasswordForm();
                this.showNotification('已取消密码修改', 'info');
            });
        }

        // 密码显示/隐藏切换
        this.bindPasswordToggle();
        
        // 密码强度检查
        this.bindPasswordValidation();
    }

    bindPasswordToggle() {
        const toggleButtons = document.querySelectorAll('.password-toggle');
        
        toggleButtons.forEach(button => {
            button.addEventListener('click', () => {
                const targetId = button.getAttribute('data-target');
                const input = document.getElementById(targetId);
                const icon = button.querySelector('i');
                
                if (input && icon) {
                    if (input.type === 'password') {
                        input.type = 'text';
                        icon.className = 'far fa-eye-slash';
                    } else {
                        input.type = 'password';
                        icon.className = 'far fa-eye';
                    }
                }
            });
        });
    }

    bindPasswordValidation() {
        const newPasswordInput = document.getElementById('newPassword');
        const confirmPasswordInput = document.getElementById('confirmPassword');
        
        if (newPasswordInput) {
            newPasswordInput.addEventListener('input', () => {
                this.checkPasswordStrength(newPasswordInput.value);
                this.checkPasswordMatch();
            });
        }
        
        if (confirmPasswordInput) {
            confirmPasswordInput.addEventListener('input', () => {
                this.checkPasswordMatch();
            });
        }
    }

    checkPasswordStrength(password) {
        const strengthIndicator = document.getElementById('passwordStrength');
        const requirements = {
            length: password.length >= 8,
            uppercase: /[A-Z]/.test(password),
            lowercase: /[a-z]/.test(password),
            number: /\d/.test(password)
        };
        
        // 更新需求列表
        Object.keys(requirements).forEach(req => {
            const element = document.getElementById(`req-${req}`);
            if (element) {
                element.classList.toggle('valid', requirements[req]);
            }
        });
        
        // 计算强度
        const validCount = Object.values(requirements).filter(Boolean).length;
        let strength = 'weak';
        
        if (validCount >= 4) {
            strength = 'strong';
        } else if (validCount >= 2) {
            strength = 'medium';
        }
        
        // 更新强度指示器
        if (strengthIndicator) {
            strengthIndicator.className = `password-strength ${strength}`;
            strengthIndicator.innerHTML = `<div class="strength-bar"></div>`;
        }
    }

    checkPasswordMatch() {
        const newPassword = document.getElementById('newPassword')?.value;
        const confirmPassword = document.getElementById('confirmPassword')?.value;
        const matchIndicator = document.getElementById('passwordMatch');
        
        if (matchIndicator && confirmPassword) {
            if (newPassword === confirmPassword) {
                matchIndicator.textContent = '密码匹配';
                matchIndicator.className = 'password-match match';
            } else {
                matchIndicator.textContent = '密码不匹配';
                matchIndicator.className = 'password-match no-match';
            }
        }
    }

    async handleAvatarUpload(event) {
        const file = event.target.files[0];
        if (!file) return;
        
        // 验证文件类型
        if (!file.type.startsWith('image/')) {
            this.showNotification('请选择图片文件', 'error');
            return;
        }
        
        // 验证文件大小 (2MB)
        if (file.size > 2 * 1024 * 1024) {
            this.showNotification('图片大小不能超过2MB', 'error');
            return;
        }
        
        // 显示预览
        const reader = new FileReader();
        reader.onload = (e) => {
            const profileAvatar = document.getElementById('profileAvatar');
            if (profileAvatar) {
                profileAvatar.src = e.target.result;
            }
        };
        reader.readAsDataURL(file);
        
        // 存储文件用于后续上传
        this.pendingAvatarFile = file;
        this.showNotification('头像预览已更新，请点击保存更改', 'info');
    }

    async changePassword() {
        console.log('🔐 修改密码...');
        
        const currentPassword = document.getElementById('currentPassword')?.value;
        const newPassword = document.getElementById('newPassword')?.value;
        const confirmPassword = document.getElementById('confirmPassword')?.value;
        
        // 验证表单
        if (!currentPassword || !newPassword || !confirmPassword) {
            this.showNotification('请填写所有密码字段', 'error');
            return;
        }
        
        if (newPassword !== confirmPassword) {
            this.showNotification('新密码和确认密码不匹配', 'error');
            return;
        }
        
        if (newPassword.length < 8) {
            this.showNotification('新密码至少需要8个字符', 'error');
            return;
        }
        
        try {
            // 这里应该调用密码修改API
            const token = this.getAuthToken();
            const response = await fetch('/api/v1/user/password', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    currentPassword,
                    newPassword
                })
            });
            
            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || '密码修改失败');
            }
            
            this.showNotification('密码修改成功！', 'success');
            this.resetPasswordForm();
            
        } catch (error) {
            console.error('密码修改失败:', error);
            this.showNotification(error.message || '密码修改失败，请重试', 'error');
        }
    }

    resetPasswordForm() {
        const passwordForm = document.getElementById('passwordForm');
        if (passwordForm) {
            passwordForm.reset();
        }
        
        // 重置密码强度指示器
        const strengthIndicator = document.getElementById('passwordStrength');
        if (strengthIndicator) {
            strengthIndicator.className = 'password-strength';
            strengthIndicator.innerHTML = '';
        }
        
        // 重置密码匹配指示器
        const matchIndicator = document.getElementById('passwordMatch');
        if (matchIndicator) {
            matchIndicator.textContent = '';
            matchIndicator.className = 'password-match';
        }
        
        // 重置需求列表
        ['length', 'uppercase', 'lowercase', 'number'].forEach(req => {
            const element = document.getElementById(`req-${req}`);
            if (element) {
                element.classList.remove('valid');
            }
        });
    }

    bindMobileEvents() {
        const mobileMenuBtn = document.getElementById('mobileMenuBtn');
        const sidebar = document.querySelector('.sidebar');
        
        if (mobileMenuBtn && sidebar) {
            mobileMenuBtn.addEventListener('click', () => {
                sidebar.classList.toggle('open');
            });
            
            // 点击内容区域关闭侧边栏
            document.querySelector('.content-area')?.addEventListener('click', () => {
                sidebar.classList.remove('open');
            });
        }
    }

    bindCourseFilterEvents() {
        const filterBtns = document.querySelectorAll('.filter-btn');
        
        filterBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                // 更新活跃状态
                filterBtns.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                
                // 过滤课程
                const filter = btn.getAttribute('data-filter');
                this.filterCourses(filter);
            });
        });
    }

    bindSecurityEvents() {
        const changePasswordBtn = document.getElementById('changePasswordBtn');
        const bindPhoneBtn = document.getElementById('bindPhoneBtn');
        const enable2FABtn = document.getElementById('enable2FABtn');
        
        if (changePasswordBtn) {
            changePasswordBtn.addEventListener('click', () => {
                this.showNotification('修改密码功能开发中...', 'info');
            });
        }
        
        if (bindPhoneBtn) {
            bindPhoneBtn.addEventListener('click', () => {
                this.showNotification('手机绑定功能开发中...', 'info');
            });
        }
        
        if (enable2FABtn) {
            enable2FABtn.addEventListener('click', () => {
                this.showNotification('两步验证功能开发中...', 'info');
            });
        }
    }

    // ===== 页面切换功能 =====
    switchSection(sectionName) {
        console.log(`🔄 切换到: ${sectionName}`);
        
        // 更新当前区域
        this.currentSection = sectionName;
        
        // 更新左侧导航状态
        this.updateSidebarNav(sectionName);
        
        // 更新面包屑和标题
        this.updatePageHeader(sectionName);
        
        // 显示对应内容区域
        this.showContentSection(sectionName);
        
        // 加载对应数据
        this.loadSectionData(sectionName);
    }

    updateSidebarNav(sectionName) {
        const navLinks = document.querySelectorAll('.sidebar .nav-link');
        navLinks.forEach(link => {
            link.classList.remove('active');
            if (link.getAttribute('data-section') === sectionName) {
                link.classList.add('active');
            }
        });
    }

    updatePageHeader(sectionName) {
        const currentPageTitle = document.getElementById('currentPageTitle');
        const pageTitle = document.getElementById('pageTitle');
        const pageSubtitle = document.getElementById('pageSubtitle');
        
        const sectionConfig = {
            'learning-progress': {
                breadcrumb: '学习进度',
                title: '学习仪表盘',
                subtitle: '掌握你的学习进度，继续你的知识之旅'
            },
            'my-courses': {
                breadcrumb: '我的课程',
                title: '我的课程',
                subtitle: '管理和继续你的学习课程'
            },
            'profile': {
                breadcrumb: '个人资料',
                title: '个人资料设置',
                subtitle: '更新你的个人信息和偏好设置'
            },
            'security': {
                breadcrumb: '账户安全',
                title: '账户安全设置',
                subtitle: '保护你的账户安全，管理登录和验证设置'
            },
            'orders': {
                breadcrumb: '订单管理',
                title: '订单管理',
                subtitle: '查看你的购买历史和订单状态'
            },
            'settings': {
                breadcrumb: '系统设置',
                title: '系统设置',
                subtitle: '个性化你的使用体验和通知偏好'
            }
        };

        const config = sectionConfig[sectionName];
        if (config) {
            if (currentPageTitle) currentPageTitle.textContent = config.breadcrumb;
            if (pageTitle) pageTitle.textContent = config.title;
            if (pageSubtitle) pageSubtitle.textContent = config.subtitle;
        }
    }

    showContentSection(sectionName) {
        // 隐藏所有内容区域
        const sections = document.querySelectorAll('.content-section');
        sections.forEach(section => {
            section.classList.remove('active');
        });
        
        // 显示目标区域
        const targetSection = document.getElementById(`${sectionName}-section`);
        if (targetSection) {
            targetSection.classList.add('active');
        }
    }

    loadSectionData(sectionName) {
        switch (sectionName) {
            case 'learning-progress':
                this.loadLearningProgress();
                break;
            case 'my-courses':
                this.loadMyCourses();
                break;
            case 'profile':
                this.loadProfileData();
                break;
            default:
                console.log(`📄 加载 ${sectionName} 数据...`);
        }
    }

    // ===== 学习进度数据加载 =====
    loadLearningProgress() {
        console.log('📊 加载学习进度数据...');
        
        // 动画显示统计数字
        this.animateStats();
        
        // 加载最近学习的课程
        this.loadRecentCourses();
    }

    animateStats() {
        const statNumbers = document.querySelectorAll('.stat-number[data-count]');
        
        statNumbers.forEach(element => {
            const target = parseInt(element.getAttribute('data-count'));
            const duration = 2000; // 2秒
            const increment = target / (duration / 16); // 60fps
            let current = 0;
            
            const updateCounter = () => {
                current += increment;
                if (current < target) {
                    element.textContent = Math.floor(current);
                    requestAnimationFrame(updateCounter);
                } else {
                    element.textContent = target;
                }
            };
            
            updateCounter();
        });
    }

    loadRecentCourses() {
        const recentCoursesGrid = document.getElementById('recentCoursesGrid');
        if (!recentCoursesGrid) return;

        // 模拟最近学习的课程数据
        const recentCourses = [
            {
                id: 1,
                title: 'Go语言微服务架构实战',
                instructor: '张三',
                progress: 75,
                coverImage: '/static/images/pastry-cover.svg',
                lastStudied: '2小时前'
            },
            {
                id: 2,
                title: 'Docker容器化部署',
                instructor: '李四',
                progress: 45,
                coverImage: '/static/images/pastry-cover.svg',
                lastStudied: '1天前'
            },
            {
                id: 3,
                title: 'React前端开发进阶',
                instructor: '王五',
                progress: 90,
                coverImage: '/static/images/pastry-cover.svg',
                lastStudied: '3天前'
            }
        ];

        recentCoursesGrid.innerHTML = recentCourses.map(course => `
            <div class="enrolled-course-card" data-course-id="${course.id}">
                <div class="course-image">
                    <img src="${course.coverImage}" alt="${course.title}" loading="lazy">
                    <div class="course-status ${course.progress === 100 ? 'completed' : course.progress > 0 ? 'in-progress' : 'not-started'}">
                        ${course.progress === 100 ? '已完成' : course.progress > 0 ? '进行中' : '未开始'}
                    </div>
                </div>
                <div class="course-content">
                    <h3 class="course-title">${course.title}</h3>
                    <div class="course-meta">
                        <span><i class="fas fa-user"></i> ${course.instructor}</span>
                        <span><i class="fas fa-clock"></i> ${course.lastStudied}</span>
                    </div>
                    <div class="course-progress">
                        <div class="progress-header">
                            <span class="progress-label">学习进度</span>
                            <span class="progress-percentage">${course.progress}%</span>
                        </div>
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: ${course.progress}%"></div>
                        </div>
                    </div>
                    <div class="course-actions">
                        <button class="action-btn primary">
                            <i class="fas fa-play"></i>
                            继续学习
                        </button>
                        <button class="action-btn secondary">
                            <i class="fas fa-info-circle"></i>
                            详情
                        </button>
                    </div>
                </div>
            </div>
        `).join('');

        // 绑定课程卡片事件
        this.bindCourseCardEvents(recentCoursesGrid);
    }

    // ===== 我的课程数据加载 =====
    loadMyCourses() {
        console.log('📚 加载我的课程数据...');
        
        const coursesGrid = document.getElementById('coursesGrid');
        const emptyState = document.getElementById('emptyState');
        
        if (!coursesGrid) return;

        // 模拟课程数据
        const allCourses = [
            {
                id: 1,
                title: 'Go语言微服务架构实战',
                instructor: '张三',
                progress: 75,
                status: 'in-progress',
                coverImage: '/static/images/pastry-cover.svg',
                rating: 4.8,
                duration: '12小时'
            },
            {
                id: 2,
                title: 'Docker容器化部署',
                instructor: '李四',
                progress: 45,
                status: 'in-progress',
                coverImage: '/static/images/pastry-cover.svg',
                rating: 4.7,
                duration: '8小时'
            },
            {
                id: 3,
                title: 'React前端开发进阶',
                instructor: '王五',
                progress: 100,
                status: 'completed',
                coverImage: '/static/images/pastry-cover.svg',
                rating: 4.9,
                duration: '15小时'
            },
            {
                id: 4,
                title: 'Kubernetes集群管理',
                instructor: '赵六',
                progress: 0,
                status: 'not-started',
                coverImage: '/static/images/pastry-cover.svg',
                rating: 4.6,
                duration: '20小时'
            }
        ];

        this.allCourses = allCourses;
        this.filterCourses('all');
    }

    filterCourses(filter) {
        const coursesGrid = document.getElementById('coursesGrid');
        const emptyState = document.getElementById('emptyState');
        
        if (!coursesGrid || !this.allCourses) return;

        let filteredCourses = this.allCourses;
        
        if (filter !== 'all') {
            filteredCourses = this.allCourses.filter(course => course.status === filter);
        }

        if (filteredCourses.length === 0) {
            coursesGrid.style.display = 'none';
            if (emptyState) emptyState.style.display = 'block';
            return;
        }

        coursesGrid.style.display = 'grid';
        if (emptyState) emptyState.style.display = 'none';

        coursesGrid.innerHTML = filteredCourses.map(course => `
            <div class="enrolled-course-card" data-course-id="${course.id}">
                <div class="course-image">
                    <img src="${course.coverImage}" alt="${course.title}" loading="lazy">
                    <div class="course-status ${course.status}">
                        ${course.status === 'completed' ? '已完成' : course.status === 'in-progress' ? '进行中' : '未开始'}
                    </div>
                </div>
                <div class="course-content">
                    <h3 class="course-title">${course.title}</h3>
                    <div class="course-meta">
                        <span><i class="fas fa-user"></i> ${course.instructor}</span>
                        <span><i class="fas fa-clock"></i> ${course.duration}</span>
                        <span><i class="fas fa-star"></i> ${course.rating}</span>
                    </div>
                    <div class="course-progress">
                        <div class="progress-header">
                            <span class="progress-label">学习进度</span>
                            <span class="progress-percentage">${course.progress}%</span>
                        </div>
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: ${course.progress}%"></div>
                        </div>
                    </div>
                    <div class="course-actions">
                        <button class="action-btn primary">
                            <i class="fas fa-${course.progress === 100 ? 'redo' : 'play'}"></i>
                            ${course.progress === 100 ? '重新学习' : course.progress > 0 ? '继续学习' : '开始学习'}
                        </button>
                        <button class="action-btn secondary">
                            <i class="fas fa-info-circle"></i>
                            详情
                        </button>
                    </div>
                </div>
            </div>
        `).join('');

        // 绑定课程卡片事件
        this.bindCourseCardEvents(coursesGrid);
    }

    bindCourseCardEvents(container) {
        const courseCards = container.querySelectorAll('.enrolled-course-card');
        
        courseCards.forEach(card => {
            const continueBtn = card.querySelector('.action-btn.primary');
            const detailBtn = card.querySelector('.action-btn.secondary');
            const courseId = card.getAttribute('data-course-id');
            
            if (continueBtn) {
                continueBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    this.continueCourse(courseId);
                });
            }
            
            if (detailBtn) {
                detailBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    this.showCourseDetail(courseId);
                });
            }
            
            // 卡片点击事件
            card.addEventListener('click', () => {
                this.showCourseDetail(courseId);
            });
        });
    }

    continueCourse(courseId) {
        console.log(`▶️ 继续学习课程: ${courseId}`);
        this.showNotification('正在跳转到课程...', 'info');
        // 这里可以跳转到课程学习页面
        // window.location.href = `/course/${courseId}/learn`;
    }

    showCourseDetail(courseId) {
        console.log(`ℹ️ 查看课程详情: ${courseId}`);
        this.showNotification('正在跳转到课程详情...', 'info');
        // 这里可以跳转到课程详情页面
        // window.location.href = `/course/${courseId}`;
    }

    // ===== 个人资料管理 =====
    loadProfileData() {
        if (!this.currentUser) return;
        
        console.log('👤 加载个人资料数据...');
        
        // 填充表单数据
        const nickname = document.getElementById('nickname');
        const username = document.getElementById('username');
        const email = document.getElementById('email');
        const phone = document.getElementById('phone');
        const bio = document.getElementById('bio');
        const profileAvatar = document.getElementById('profileAvatar');
        
        if (nickname) nickname.value = this.currentUser.nickname || '';
        if (username) username.value = this.currentUser.username || '';
        if (email) email.value = this.currentUser.email || '';
        if (phone) phone.value = this.currentUser.phone || '';
        if (bio) bio.value = this.currentUser.bio || '';
        if (profileAvatar) profileAvatar.src = this.currentUser.avatar || '/static/images/default-avatar.svg';
    }

    async saveProfile() {
        console.log('💾 保存个人资料...');
        
        const nickname = document.getElementById('nickname')?.value;
        const phone = document.getElementById('phone')?.value;
        const bio = document.getElementById('bio')?.value;

        try {
            const token = this.getAuthToken();
            console.log('🔑 获取到的Token:', token ? token.substring(0, 30) + '...' : 'null');
            
            if (!token) {
                throw new Error('未找到认证Token，请重新登录');
            }
            
            let avatarUrl = this.currentUser.avatar;
            
            // 如果有新头像，尝试上传头像（可选）
            if (this.pendingAvatarFile) {
                console.log('📤 尝试上传头像...');
                
                try {
                    const avatarFormData = new FormData();
                    avatarFormData.append('file', this.pendingAvatarFile);
                    
                    const uploadResponse = await fetch('/api/v1/content/upload', {
                        method: 'POST',
                        headers: {
                            'Authorization': `Bearer ${token}`
                        },
                        body: avatarFormData
                    });
                    
                    if (uploadResponse.ok) {
                        const uploadResult = await uploadResponse.json();
                        avatarUrl = uploadResult.data.url;
                        console.log('✅ 头像上传成功:', avatarUrl);
                    } else {
                        console.warn('⚠️ 头像上传失败，使用默认头像继续保存资料');
                        // 使用默认头像URL或保持原有头像
                        avatarUrl = this.currentUser.avatar || '/static/images/default-avatar.svg';
                    }
                } catch (uploadError) {
                    console.warn('⚠️ 头像上传异常，使用默认头像继续保存资料:', uploadError.message);
                    // 头像上传失败不影响资料保存
                    avatarUrl = this.currentUser.avatar || '/static/images/default-avatar.svg';
                }
            }
            
            // 保存个人资料（即使头像上传失败也继续）
            const profileData = {
                nickname: nickname,
                phone: phone,
                bio: bio,
                avatar: avatarUrl
            };
            
            console.log('📤 发送个人资料数据:', profileData);
            console.log('🔗 使用Authorization:', `Bearer ${token.substring(0, 30)}...`);
            
            const response = await fetch('/api/v1/user/profile', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(profileData)
            });
            
            if (!response.ok) {
                console.log('❌ 响应状态:', response.status, response.statusText);
                let errorMessage = '保存失败';
                try {
                    const error = await response.json();
                    errorMessage = error.message || error.error || '保存失败';
                    console.log('❌ 错误详情:', error);
                } catch (e) {
                    console.log('❌ 无法解析错误响应');
                }
                throw new Error(errorMessage);
            }
            
            const result = await response.json();
            
            // 更新本地用户信息
            Object.assign(this.currentUser, result.user || profileData);
            
            // 同步更新本地存储的用户信息
            localStorage.setItem('userInfo', JSON.stringify(this.currentUser));
            sessionStorage.setItem('userInfo', JSON.stringify(this.currentUser));
            
            this.updateUserInterface();
            
            // 清除待上传的头像文件
            this.pendingAvatarFile = null;
            
            this.showNotification('个人资料保存成功！', 'success');
            console.log('✅ 个人资料保存成功');
            
        } catch (error) {
            console.error('保存个人资料失败:', error);
            this.showNotification(error.message || '保存失败，请重试', 'error');
        }
    }

    // ===== 用户操作 =====
    logout() {
        console.log('👋 用户退出登录');
        
        // 清除认证数据
        this.clearAuthData();
        
        // 显示通知
        this.showNotification('您已成功退出登录', 'success');
        
        // 重定向到首页
        setTimeout(() => {
            window.location.href = '/';
        }, 1500);
    }

    // ===== 工具方法 =====
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

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    loadDefaultSection() {
        this.switchSection('learning-progress');
    }

    initializeAnimations() {
        // 页面加载动画
        const contentSections = document.querySelectorAll('.content-section');
        contentSections.forEach((section, index) => {
            section.style.animationDelay = `${index * 0.1}s`;
        });
    }

    showNotification(message, type = 'info') {
        const container = document.getElementById('notificationContainer') || document.body;
        
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.innerHTML = `
            <div class="notification-content">
                <i class="fas fa-${this.getNotificationIcon(type)}"></i>
                <span>${message}</span>
            </div>
            <button class="notification-close">
                <i class="fas fa-times"></i>
            </button>
        `;
        
        container.appendChild(notification);
        
        // 显示动画
        setTimeout(() => notification.classList.add('show'), 10);
        
        // 关闭按钮事件
        const closeBtn = notification.querySelector('.notification-close');
        closeBtn.addEventListener('click', () => {
            this.hideNotification(notification);
        });
        
        // 自动消失
        setTimeout(() => {
            this.hideNotification(notification);
        }, 5000);
    }

    hideNotification(notification) {
        notification.classList.remove('show');
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }

    getNotificationIcon(type) {
        const icons = {
            success: 'check-circle',
            error: 'exclamation-circle',
            warning: 'exclamation-triangle',
            info: 'info-circle'
        };
        return icons[type] || 'info-circle';
    }
}

// ===== 初始化用户中心 =====
document.addEventListener('DOMContentLoaded', () => {
    new UserDashboard();
});

// ===== 添加通知样式 =====
const notificationStyles = `
.notification-container {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 10000;
    pointer-events: none;
}

.notification {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    padding: 1rem 1.5rem;
    margin-bottom: 0.5rem;
    min-width: 300px;
    max-width: 400px;
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(10px);
    transform: translateX(400px);
    transition: all 0.3s ease;
    pointer-events: auto;
}

.notification.show {
    transform: translateX(0);
}

.notification-content {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
}

.notification-success {
    border-left: 4px solid #22c55e;
}

.notification-success .notification-content i {
    color: #22c55e;
}

.notification-error {
    border-left: 4px solid #ef4444;
}

.notification-error .notification-content i {
    color: #ef4444;
}

.notification-warning {
    border-left: 4px solid #f59e0b;
}

.notification-warning .notification-content i {
    color: #f59e0b;
}

.notification-info {
    border-left: 4px solid #3b82f6;
}

.notification-info .notification-content i {
    color: #3b82f6;
}

.notification-close {
    background: none;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    padding: 0.25rem;
    border-radius: 4px;
    transition: all 0.3s ease;
}

.notification-close:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
}

@media (max-width: 768px) {
    .notification-container {
        left: 20px;
        right: 20px;
    }
    
    .notification {
        min-width: auto;
        max-width: none;
    }
}
`;

// 添加样式到页面
const styleSheet = document.createElement('style');
styleSheet.textContent = notificationStyles;
document.head.appendChild(styleSheet); 