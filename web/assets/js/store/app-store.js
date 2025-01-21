// Make sure createAppStore is available globally
window.createAppStore = function() {
    return {
        searchQuery: '',
        showSearch: false,
        user: null,
        isAuthenticated: false,

        init() {
            console.log('Store init called');
            this.user = JSON.parse(localStorage.getItem('user') || 'null');
            this.isAuthenticated = !!localStorage.getItem('token');
        },

        toggleSearch(value) {
            console.log('toggleSearch called with:', value);
            this.showSearch = value;
            if (value) {
                console.log('Dispatching open-search event');
                window.dispatchEvent(new CustomEvent('open-search', { 
                    detail: { query: this.searchQuery }
                }));
            }
        },

        setSearchQuery(query) {
            console.log('setSearchQuery:', query);
            this.searchQuery = query;
        },

        logout() {
            localStorage.removeItem('token');
            this.user = null;
            window.location.reload();
        }
    };
}; 