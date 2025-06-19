// 登录页面逻辑
class LoginPage {
    constructor() {
        this.authForm = null;
        this.init();
    }
    
    init() {
        // 页面加载完成后初始化
        document.addEventListener('DOMContentLoaded', () => {
            this.hidePageLoader();
            this.showRegistrationSuccessAlert();
            this.initAuthForm();
            this.initMobileMenu();
            this.checkExistingAuth();
            this.handleURLParams();
        });
    }
    
    hidePageLoader() {
        const loader = document.getElementById('pageLoader');
        if (loader) {
            setTimeout(() => {
                loader.classList.add('hidden');
                setTimeout(() => {
                    loader.remove();
                }, 500);
            }, 300);
        }
    }
    
    showRegistrationSuccessAlert() {
        // 检查是否有注册成功信息
        const registrationSuccess = sessionStorage.getItem('registrationSuccess');
        if (registrationSuccess) {
            try {
                const successData = JSON.parse(registrationSuccess);
                
                // 检查时间戳，防止显示过期的消息
                const now = Date.now();
                const timestamp = successData.timestamp || 0;
                const timeElapsed = now - timestamp;
                
                // 只在5分钟内显示
                if (timeElapsed < 5 * 60 * 1000) {
                    const alertElement = document.getElementById('registrationSuccessAlert');
                    const messageElement = document.getElementById('successMessage');
                    
                    if (alertElement && messageElement) {
                        messageElement.textContent = successData.message || '注册成功！请使用您的账户信息登录';
                        alertElement.style.display = 'block';
                        
                        // 预填用户名（如果有）
                        if (successData.username) {
                            setTimeout(() => {
                                this.prefillUsername(successData.username);
                            }, 1000);
                        }
                        
                        // 绑定关闭按钮事件
                        const closeBtn = document.getElementById('closeSuccessAlert');
                        if (closeBtn) {
                            closeBtn.addEventListener('click', () => {
                                this.hideRegistrationSuccessAlert();
                            });
                        }
                        
                        // 10秒后自动隐藏
                        setTimeout(() => {
                            this.hideRegistrationSuccessAlert();
                        }, 10000);
                    }
                }
                
                // 清除sessionStorage中的数据
                sessionStorage.removeItem('registrationSuccess');
                
            } catch (error) {
                console.error('解析注册成功信息失败:', error);
                sessionStorage.removeItem('registrationSuccess');
            }
        }
    }
    
    hideRegistrationSuccessAlert() {
        const alertElement = document.getElementById('registrationSuccessAlert');
        if (alertElement) {
            alertElement.style.animation = 'slideOutToTop 0.3s ease-in forwards';
            setTimeout(() => {
                alertElement.style.display = 'none';
            }, 300);
        }
    }
    
    prefillUsername(username) {
        // 预填用户名到登录表单
        if (this.authForm) {
            const identifierInput = document.getElementById('identifier');
            if (identifierInput) {
                identifierInput.value = username;
                identifierInput.focus();
                
                // 光标移到末尾
                identifierInput.setSelectionRange(username.length, username.length);
            }
        }
    }
    
    checkExistingAuth() {
        // 检查用户是否已经登录
        const token = localStorage.getItem('authToken');
        if (token) {
            // 验证token是否有效
            this.validateToken(token);
        }
    }
    
    async validateToken(token) {
        try {
            const response = await fetch('/api/v1/validate-token', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (response.ok) {
                // Token有效，重定向到首页或指定页面
                const redirectUrl = this.getRedirectUrl();
                this.showNotification('您已登录，正在跳转...', 'info');
                setTimeout(() => {
                    window.location.href = redirectUrl;
                }, 1500);
            }
        } catch (error) {
            console.log('Token验证失败:', error);
            // Token无效，清除并继续
            localStorage.removeItem('authToken');
        }
    }
    
    getRedirectUrl() {
        // 获取重定向URL
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get('redirect') || '/';
    }
    
    handleURLParams() {
        // 处理URL参数
        const urlParams = new URLSearchParams(window.location.search);
        const error = urlParams.get('error');
        const message = urlParams.get('message');
        
        if (error) {
            setTimeout(() => {
                this.showNotification(decodeURIComponent(error), 'error');
            }, 1000);
        }
        
        if (message) {
            setTimeout(() => {
                this.showNotification(decodeURIComponent(message), 'info');
            }, 1000);
        }
    }
    
    initAuthForm() {
        const container = document.getElementById('authFormContainer');
        if (!container) {
            console.error('找不到AuthForm容器');
            return;
        }
        
        // 创建AuthForm实例
        this.authForm = new AuthForm(container, {
            isRegister: false,
            onSubmit: (formData) => this.handleLogin(formData)
        });
        
        console.log('登录表单初始化完成');
    }
    
    initMobileMenu() {
        const mobileMenuBtn = document.getElementById('mobileMenuBtn');
        const navMenu = document.querySelector('.nav-menu');
        
        if (mobileMenuBtn && navMenu) {
            mobileMenuBtn.addEventListener('click', () => {
                navMenu.classList.toggle('active');
                mobileMenuBtn.classList.toggle('active');
            });
            
            // 点击菜单项时关闭移动菜单
            const navLinks = navMenu.querySelectorAll('.nav-link');
            navLinks.forEach(link => {
                link.addEventListener('click', () => {
                    navMenu.classList.remove('active');
                    mobileMenuBtn.classList.remove('active');
                });
            });
        }
    }
    
    async handleLogin(formData) {
        console.log('处理登录请求:', formData);
        
        try {
            // 验证表单数据
            this.validateLoginData(formData);
            
            // 发送登录请求
            const response = await fetch('/api/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    identifier: formData.identifier, // 用户名或邮箱
                    password: formData.password,
                    rememberMe: formData.rememberMe === 'on'
                })
            });
            
            const result = await response.json();
            
            if (response.ok) {
                // 登录成功
                console.log('登录成功:', result);
                this.handleLoginSuccess(result, formData.rememberMe === 'on');
            } else {
                // 登录失败
                console.error('登录失败:', result);
                this.handleLoginError(result);
            }
            
        } catch (error) {
            console.error('登录请求错误:', error);
            this.handleLoginError({
                message: error.message || '网络错误，请检查网络连接后重试'
            });
        }
    }
    
    validateLoginData(formData) {
        // 用户名/邮箱验证
        if (!formData.identifier || formData.identifier.trim().length < 3) {
            throw new Error('请输入有效的用户名或邮箱地址');
        }
        
        // 密码验证
        if (!formData.password || formData.password.length < 6) {
            throw new Error('密码至少需要6个字符');
        }
    }
    
    handleLoginSuccess(result, rememberMe) {
        // 保存认证token
        if (result.token) {
            if (rememberMe) {
                // 记住我：保存到localStorage（长期有效）
                localStorage.setItem('authToken', result.token);
                localStorage.setItem('tokenExpiry', result.expiry || '');
            } else {
                // 不记住：保存到sessionStorage（会话有效）
                sessionStorage.setItem('authToken', result.token);
                sessionStorage.setItem('tokenExpiry', result.expiry || '');
            }
        }
        
        // 保存用户信息
        if (result.user) {
            const userInfo = {
                id: result.user.id,
                username: result.user.username,
                email: result.user.email,
                loginTime: new Date().toISOString()
            };
            
            if (rememberMe) {
                localStorage.setItem('userInfo', JSON.stringify(userInfo));
            } else {
                sessionStorage.setItem('userInfo', JSON.stringify(userInfo));
            }
        }
        
        // 显示成功消息
        this.showNotification('登录成功！正在跳转...', 'success');
        
        // 记录登录成功事件
        this.trackEvent('user_login_success', {
            username: result.user?.username,
            rememberMe: rememberMe,
            timestamp: new Date().toISOString()
        });
        
        // 延迟跳转
        setTimeout(() => {
            const redirectUrl = this.getRedirectUrl();
            window.location.href = redirectUrl;
        }, 1500);
    }
    
    handleLoginError(error) {
        let errorMessage = '登录失败，请重试';
        
        // 解析具体错误信息
        if (error.message) {
            errorMessage = error.message;
        } else if (error.error) {
            errorMessage = error.error;
        } else if (error.details) {
            errorMessage = error.details;
        }
        
        // 处理常见错误
        if (errorMessage.includes('password') && errorMessage.includes('incorrect')) {
            errorMessage = '密码错误，请检查后重试';
        } else if (errorMessage.includes('user') && errorMessage.includes('not found')) {
            errorMessage = '用户不存在，请检查用户名或邮箱';
        } else if (errorMessage.includes('account') && errorMessage.includes('locked')) {
            errorMessage = '账户已被锁定，请联系管理员';
        } else if (errorMessage.includes('credentials')) {
            errorMessage = '用户名或密码错误';
        } else if (errorMessage.includes('server')) {
            errorMessage = '服务器暂时不可用，请稍后重试';
        }
        
        // 记录登录失败事件
        this.trackEvent('user_login_failed', {
            error: errorMessage,
            timestamp: new Date().toISOString()
        });
        
        // 抛出错误，让AuthForm处理显示
        throw new Error(errorMessage);
    }
    
    showNotification(message, type = 'info') {
        // 使用AuthForm的通知系统
        if (this.authForm) {
            this.authForm.showNotification(message, type);
        } else {
            // 备用通知方法
            console.log(`[${type.toUpperCase()}] ${message}`);
            alert(message);
        }
    }
    
    trackEvent(eventName, data) {
        // 用户行为追踪
        try {
            console.log('追踪事件:', eventName, data);
            
            // 这里可以集成第三方分析工具
            // 例如: gtag('event', eventName, data);
            // 或者: analytics.track(eventName, data);
            
            // 发送到后端分析接口
            if (navigator.sendBeacon) {
                const analyticsData = {
                    event: eventName,
                    data: data,
                    url: window.location.href,
                    userAgent: navigator.userAgent,
                    timestamp: Date.now()
                };
                
                navigator.sendBeacon('/api/v1/analytics', JSON.stringify(analyticsData));
            }
        } catch (error) {
            console.error('事件追踪失败:', error);
        }
    }
    
    // 工具方法：获取URL参数
    getUrlParameter(name) {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get(name);
    }
    
    // 工具方法：设置页面标题
    setPageTitle(title) {
        document.title = title ? `${title} - Course Platform` : 'Course Platform';
    }
    
    // 工具方法：检查当前用户是否已登录
    isLoggedIn() {
        const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken');
        return !!token;
    }
    
    // 工具方法：获取当前用户信息
    getCurrentUser() {
        const userInfo = localStorage.getItem('userInfo') || sessionStorage.getItem('userInfo');
        if (userInfo) {
            try {
                return JSON.parse(userInfo);
            } catch (error) {
                console.error('解析用户信息失败:', error);
                return null;
            }
        }
        return null;
    }
    
    // 工具方法：登出
    logout() {
        // 清除所有认证相关数据
        localStorage.removeItem('authToken');
        localStorage.removeItem('tokenExpiry');
        localStorage.removeItem('userInfo');
        sessionStorage.removeItem('authToken');
        sessionStorage.removeItem('tokenExpiry');
        sessionStorage.removeItem('userInfo');
        
        // 重定向到登录页面
        window.location.href = '/login';
    }
}

// 添加动画样式
const style = document.createElement('style');
style.textContent = `
    @keyframes slideOutToTop {
        from {
            opacity: 1;
            transform: translateY(0);
        }
        to {
            opacity: 0;
            transform: translateY(-20px);
        }
    }
`;
document.head.appendChild(style);

// 实例化登录页面
const loginPage = new LoginPage();

// 导出到全局作用域，方便调试
window.loginPage = loginPage; 