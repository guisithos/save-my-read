function createAppStore() {
    return {
        searchQuery: '',
        showSearch: false,
        user: JSON.parse(localStorage.getItem('user') || 'null'),
        isAuthenticated: false,

        init() {
            // Initialize auth state
            this.isAuthenticated = !!localStorage.getItem('token');
            
            // Watch for user changes
            this.$watch('user', (value) => {
                if (value) {
                    localStorage.setItem('user', JSON.stringify(value));
                    this.isAuthenticated = true;
                } else {
                    localStorage.removeItem('user');
                    this.isAuthenticated = false;
                }
            });
        },

        toggleSearch(value) {
            this.showSearch = value;
            if (value) {
                this.$dispatch('open-search', { query: this.searchQuery });
            }
        },

        setSearchQuery(query) {
            this.searchQuery = query;
        },

        logout() {
            localStorage.removeItem('token');
            this.user = null;
            window.location.reload();
        }
    };
} 