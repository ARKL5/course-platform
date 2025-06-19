// 注册页面逻辑
class RegisterPage {
    constructor() {
        this.authForm = null;
        this.init();
    }
    
    init() {
        // 页面加载完成后初始化
        document.addEventListener('DOMContentLoaded', () => {
            this.hidePageLoader();
            this.initAuthForm();
            this.initMobileMenu();
            this.checkExistingAuth();
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
                // Token有效，重定向到首页
                this.showNotification('您已登录，正在跳转到首页...', 'info');
                setTimeout(() => {
                    window.location.href = '/';
                }, 1500);
            }
        } catch (error) {
            console.log('Token验证失败:', error);
            // Token无效，清除并继续
            localStorage.removeItem('authToken');
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
            isRegister: true,
            onSubmit: (formData) => this.handleRegister(formData)
        });
        
        console.log('注册表单初始化完成');
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
    
    async handleRegister(formData) {
        console.log('处理注册请求:', formData);
        
        try {
            // 验证表单数据
            this.validateRegistrationData(formData);
            
            // 发送注册请求
            const response = await fetch('/api/v1/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: formData.username,
                    email: formData.email,
                    password: formData.password
                })
            });
            
            const result = await response.json();
            
            if (response.ok) {
                // 注册成功
                console.log('注册成功:', result);
                this.handleRegisterSuccess(result);
            } else {
                // 注册失败
                console.error('注册失败:', result);
                this.handleRegisterError(result);
            }
            
        } catch (error) {
            console.error('注册请求错误:', error);
            this.handleRegisterError({
                message: error.message || '网络错误，请检查网络连接后重试'
            });
        }
    }
    
    validateRegistrationData(formData) {
        // 用户名验证
        if (!formData.username || formData.username.length < 3) {
            throw new Error('用户名至少需要3个字符');
        }
        
        if (!/^[a-zA-Z0-9_]+$/.test(formData.username)) {
            throw new Error('用户名只能包含字母、数字和下划线');
        }
        
        // 邮箱验证
        if (!formData.email) {
            throw new Error('邮箱地址不能为空');
        }
        
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(formData.email)) {
            throw new Error('请输入有效的邮箱地址');
        }
        
        // 密码验证
        if (!formData.password || formData.password.length < 6) {
            throw new Error('密码至少需要6个字符');
        }
        
        if (formData.password.length > 100) {
            throw new Error('密码长度不能超过100个字符');
        }
    }
    
    handleRegisterSuccess(result) {
        // 显示成功消息
        this.showNotification('注册成功！正在跳转到登录页面...', 'success');
        
        // 保存成功信息到sessionStorage，供登录页面显示
        sessionStorage.setItem('registrationSuccess', JSON.stringify({
            message: '注册成功！请使用您的账户信息登录',
            username: result.username || '',
            timestamp: Date.now()
        }));
        
        // 延迟跳转到登录页面
        setTimeout(() => {
            window.location.href = '/login';
        }, 2000);
        
        // 记录注册成功事件
        this.trackEvent('user_register_success', {
            username: result.username,
            email: result.email,
            timestamp: new Date().toISOString()
        });
    }
    
    handleRegisterError(error) {
        let errorMessage = '注册失败，请重试';
        
        // 解析具体错误信息
        if (error.message) {
            errorMessage = error.message;
        } else if (error.error) {
            errorMessage = error.error;
        } else if (error.details) {
            errorMessage = error.details;
        }
        
        // 处理常见错误
        if (errorMessage.includes('username') && errorMessage.includes('exists')) {
            errorMessage = '用户名已存在，请选择其他用户名';
        } else if (errorMessage.includes('email') && errorMessage.includes('exists')) {
            errorMessage = '邮箱地址已被注册，请使用其他邮箱或直接登录';
        } else if (errorMessage.includes('validation')) {
            errorMessage = '输入信息格式不正确，请检查后重试';
        } else if (errorMessage.includes('server')) {
            errorMessage = '服务器暂时不可用，请稍后重试';
        }
        
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
    
    // 工具方法：更新页面meta描述
    updateMetaDescription(description) {
        const metaDescription = document.querySelector('meta[name="description"]');
        if (metaDescription) {
            metaDescription.setAttribute('content', description);
        }
    }
}

// 实例化注册页面
const registerPage = new RegisterPage();

// 导出到全局作用域，方便调试
window.registerPage = registerPage; 