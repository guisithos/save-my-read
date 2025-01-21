function auth() {
    return {
        isRegistering: false,
        showPassword: false,
        loginForm: {
            email: '',
            password: '',
            isLoading: false,
            errors: {},
        },
        registerForm: {
            email: '',
            password: '',
            name: '',
            genres: [],
            isLoading: false,
            errors: {},
        },

        validateForm(type) {
            const form = type === 'login' ? this.loginForm : this.registerForm;
            form.errors = {};
            
            if (!form.email) {
                form.errors.email = 'Email is required';
            } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
                form.errors.email = 'Invalid email format';
            }

            if (!form.password) {
                form.errors.password = 'Password is required';
            } else if (form.password.length < 8) {
                form.errors.password = 'Password must be at least 8 characters';
            }

            if (type === 'register') {
                if (!form.name) {
                    form.errors.name = 'Name is required';
                }
            }

            return Object.keys(form.errors).length === 0;
        },

        async handleLogin() {
            if (!this.validateForm('login')) return;
            
            this.loginForm.isLoading = true;
            this.loginForm.errors = {};

            try {
                const response = await api.login({
                    email: this.loginForm.email,
                    password: this.loginForm.password,
                });

                auth.setToken(response.data.token);
                window.location.reload();
            } catch (error) {
                this.loginForm.errors.general = error.message;
            } finally {
                this.loginForm.isLoading = false;
            }
        },

        // Similar handleRegister method...
    }
} 