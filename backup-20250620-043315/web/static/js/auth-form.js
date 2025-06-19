// AuthForm 可复用组件
class AuthForm {
    constructor(container, options = {}) {
        this.container = container;
        this.isRegister = options.isRegister || false;
        this.onSubmit = options.onSubmit || (() => {});
        this.isLoading = false;
        this.formData = {};
        this.errors = {};
        
        this.init();
    }
    
    init() {
        this.render();
        this.bindEvents();
    }
    
    render() {
        const formHTML = `
            <div class="auth-form-container">
                <div class="auth-form-card">
                    <div class="auth-form-header">
                        <h2 class="auth-title">
                            ${this.isRegister ? '创建您的账户' : '欢迎回来'}
                        </h2>
                        <p class="auth-subtitle">
                            ${this.isRegister ? '加入我们的学习社区' : '继续您的学习之旅'}
                        </p>
                    </div>
                    
                    <div class="auth-error-area" id="authErrorArea" style="display: none;">
                        <div class="error-message">
                            <i class="fas fa-exclamation-circle"></i>
                            <span id="errorText"></span>
                        </div>
                    </div>
                    
                    <form class="auth-form" id="authForm">
                        ${this.renderFormFields()}
                        
                        <button type="submit" class="auth-submit-btn" id="submitBtn">
                            <span class="btn-text">
                                ${this.isRegister ? '创建账户' : '登录'}
                            </span>
                            <div class="btn-loading" style="display: none;">
                                <i class="fas fa-spinner fa-spin"></i>
                                处理中...
                            </div>
                        </button>
                    </form>
                    
                    <div class="auth-footer">
                        ${this.renderFooter()}
                    </div>
                </div>
            </div>
        `;
        
        this.container.innerHTML = formHTML;
    }
    
    renderFormFields() {
        if (this.isRegister) {
            return `
                <div class="form-group">
                    <label for="username" class="form-label">用户名</label>
                    <div class="input-wrapper">
                        <i class="fas fa-user input-icon"></i>
                        <input 
                            type="text" 
                            id="username" 
                            name="username" 
                            class="form-input" 
                            placeholder="请输入用户名"
                            required
                        >
                    </div>
                    <div class="field-error" id="usernameError"></div>
                </div>
                
                <div class="form-group">
                    <label for="email" class="form-label">邮箱地址</label>
                    <div class="input-wrapper">
                        <i class="fas fa-envelope input-icon"></i>
                        <input 
                            type="email" 
                            id="email" 
                            name="email" 
                            class="form-input" 
                            placeholder="请输入邮箱地址"
                            required
                        >
                    </div>
                    <div class="field-error" id="emailError"></div>
                </div>
                
                <div class="form-group">
                    <label for="password" class="form-label">密码</label>
                    <div class="input-wrapper">
                        <i class="fas fa-lock input-icon"></i>
                        <input 
                            type="password" 
                            id="password" 
                            name="password" 
                            class="form-input" 
                            placeholder="请输入密码（至少6位）"
                            required
                        >
                        <button type="button" class="password-toggle" id="passwordToggle">
                            <i class="fas fa-eye"></i>
                        </button>
                    </div>
                    <div class="field-error" id="passwordError"></div>
                </div>
            `;
        } else {
            return `
                <div class="form-group">
                    <label for="identifier" class="form-label">用户名或邮箱</label>
                    <div class="input-wrapper">
                        <i class="fas fa-user input-icon"></i>
                        <input 
                            type="text" 
                            id="identifier" 
                            name="identifier" 
                            class="form-input" 
                            placeholder="请输入用户名或邮箱"
                            required
                        >
                    </div>
                    <div class="field-error" id="identifierError"></div>
                </div>
                
                <div class="form-group">
                    <label for="password" class="form-label">密码</label>
                    <div class="input-wrapper">
                        <i class="fas fa-lock input-icon"></i>
                        <input 
                            type="password" 
                            id="password" 
                            name="password" 
                            class="form-input" 
                            placeholder="请输入密码"
                            required
                        >
                        <button type="button" class="password-toggle" id="passwordToggle">
                            <i class="fas fa-eye"></i>
                        </button>
                    </div>
                    <div class="field-error" id="passwordError"></div>
                </div>
                
                <div class="form-options">
                    <label class="checkbox-wrapper">
                        <input type="checkbox" id="rememberMe" name="rememberMe">
                        <span class="checkmark"></span>
                        <span class="checkbox-text">记住我</span>
                    </label>
                    <a href="#" class="forgot-password">忘记密码？</a>
                </div>
            `;
        }
    }
    
    renderFooter() {
        if (this.isRegister) {
            return `
                <p class="auth-link">
                    已有账户？
                    <a href="/login" class="link-primary">立即登录</a>
                </p>
            `;
        } else {
            return `
                <p class="auth-link">
                    还没有账户？
                    <a href="/register" class="link-primary">立即注册</a>
                </p>
            `;
        }
    }
    
    bindEvents() {
        const form = this.container.querySelector('#authForm');
        const passwordToggle = this.container.querySelector('#passwordToggle');
        
        // 表单提交事件
        form.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleSubmit();
        });
        
        // 密码可见性切换
        if (passwordToggle) {
            passwordToggle.addEventListener('click', () => {
                this.togglePasswordVisibility();
            });
        }
        
        // 实时验证
        const inputs = this.container.querySelectorAll('.form-input');
        inputs.forEach(input => {
            input.addEventListener('blur', () => {
                this.validateField(input);
            });
            
            input.addEventListener('input', () => {
                this.clearFieldError(input.name);
            });
        });
        
        // 忘记密码点击事件
        const forgotPasswordLink = this.container.querySelector('.forgot-password');
        if (forgotPasswordLink) {
            forgotPasswordLink.addEventListener('click', (e) => {
                e.preventDefault();
                this.showNotification('密码重置功能暂未开放', 'info');
            });
        }
    }
    
    togglePasswordVisibility() {
        const passwordInput = this.container.querySelector('#password');
        const toggleIcon = this.container.querySelector('#passwordToggle i');
        
        if (passwordInput.type === 'password') {
            passwordInput.type = 'text';
            toggleIcon.classList.remove('fa-eye');
            toggleIcon.classList.add('fa-eye-slash');
        } else {
            passwordInput.type = 'password';
            toggleIcon.classList.remove('fa-eye-slash');
            toggleIcon.classList.add('fa-eye');
        }
    }
    
    validateField(input) {
        const value = input.value.trim();
        const name = input.name;
        let isValid = true;
        let errorMessage = '';
        
        // 清除之前的错误
        this.clearFieldError(name);
        
        // 基础验证
        if (!value) {
            isValid = false;
            errorMessage = '此字段为必填项';
        } else {
            // 特定字段验证
            switch (name) {
                case 'username':
                    if (value.length < 3) {
                        isValid = false;
                        errorMessage = '用户名至少需要3个字符';
                    } else if (!/^[a-zA-Z0-9_]+$/.test(value)) {
                        isValid = false;
                        errorMessage = '用户名只能包含字母、数字和下划线';
                    }
                    break;
                    
                case 'email':
                    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                    if (!emailRegex.test(value)) {
                        isValid = false;
                        errorMessage = '请输入有效的邮箱地址';
                    }
                    break;
                    
                case 'password':
                    if (value.length < 6) {
                        isValid = false;
                        errorMessage = '密码至少需要6个字符';
                    }
                    break;
                    
                case 'identifier':
                    if (value.length < 3) {
                        isValid = false;
                        errorMessage = '用户名或邮箱长度不能少于3个字符';
                    }
                    break;
            }
        }
        
        if (!isValid) {
            this.showFieldError(name, errorMessage);
        }
        
        return isValid;
    }
    
    showFieldError(fieldName, message) {
        const errorElement = this.container.querySelector(`#${fieldName}Error`);
        const inputWrapper = this.container.querySelector(`input[name="${fieldName}"]`).closest('.input-wrapper');
        
        if (errorElement) {
            errorElement.textContent = message;
            errorElement.style.display = 'block';
        }
        
        if (inputWrapper) {
            inputWrapper.classList.add('error');
        }
    }
    
    clearFieldError(fieldName) {
        const errorElement = this.container.querySelector(`#${fieldName}Error`);
        const inputWrapper = this.container.querySelector(`input[name="${fieldName}"]`).closest('.input-wrapper');
        
        if (errorElement) {
            errorElement.style.display = 'none';
        }
        
        if (inputWrapper) {
            inputWrapper.classList.remove('error');
        }
    }
    
    showError(message) {
        const errorArea = this.container.querySelector('#authErrorArea');
        const errorText = this.container.querySelector('#errorText');
        
        errorText.textContent = message;
        errorArea.style.display = 'block';
        
        // 滚动到错误区域
        errorArea.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
    
    hideError() {
        const errorArea = this.container.querySelector('#authErrorArea');
        errorArea.style.display = 'none';
    }
    
    setLoading(loading) {
        this.isLoading = loading;
        const submitBtn = this.container.querySelector('#submitBtn');
        const btnText = submitBtn.querySelector('.btn-text');
        const btnLoading = submitBtn.querySelector('.btn-loading');
        
        if (loading) {
            btnText.style.display = 'none';
            btnLoading.style.display = 'flex';
            submitBtn.disabled = true;
        } else {
            btnText.style.display = 'block';
            btnLoading.style.display = 'none';
            submitBtn.disabled = false;
        }
    }
    
    collectFormData() {
        const form = this.container.querySelector('#authForm');
        const formData = new FormData(form);
        const data = {};
        
        for (let [key, value] of formData.entries()) {
            data[key] = value.trim();
        }
        
        return data;
    }
    
    validateForm() {
        const inputs = this.container.querySelectorAll('.form-input');
        let isValid = true;
        
        inputs.forEach(input => {
            if (!this.validateField(input)) {
                isValid = false;
            }
        });
        
        return isValid;
    }
    
    async handleSubmit() {
        this.hideError();
        
        if (!this.validateForm()) {
            this.showError('请修正表单中的错误后重试');
            return;
        }
        
        const formData = this.collectFormData();
        
        this.setLoading(true);
        
        try {
            await this.onSubmit(formData);
        } catch (error) {
            console.error('表单提交错误:', error);
            this.showError(error.message || '提交失败，请重试');
        } finally {
            this.setLoading(false);
        }
    }
    
    // 通知系统
    showNotification(message, type = 'info') {
        // 创建通知元素
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.innerHTML = `
            <div class="notification-content">
                <i class="fas ${this.getNotificationIcon(type)}"></i>
                <span>${message}</span>
            </div>
            <button class="notification-close">
                <i class="fas fa-times"></i>
            </button>
        `;
        
        // 添加到页面
        document.body.appendChild(notification);
        
        // 绑定关闭事件
        const closeBtn = notification.querySelector('.notification-close');
        closeBtn.addEventListener('click', () => {
            this.removeNotification(notification);
        });
        
        // 自动关闭
        setTimeout(() => {
            this.removeNotification(notification);
        }, 5000);
        
        // 显示动画
        setTimeout(() => {
            notification.classList.add('show');
        }, 100);
    }
    
    getNotificationIcon(type) {
        const icons = {
            success: 'fa-check-circle',
            error: 'fa-exclamation-circle',
            warning: 'fa-exclamation-triangle',
            info: 'fa-info-circle'
        };
        return icons[type] || icons.info;
    }
    
    removeNotification(notification) {
        notification.classList.remove('show');
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }
}

// 导出给全局使用
window.AuthForm = AuthForm; 