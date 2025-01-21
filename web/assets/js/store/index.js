// Initialize Alpine.js store
document.addEventListener('alpine:init', () => {
    console.log('Initializing Alpine store');
    
    Alpine.store('app', {
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
    });

    // Initialize store
    Alpine.store('app').init();
});

// Verify store initialization
document.addEventListener('alpine:initialized', () => {
    console.log('Alpine initialized, store status:', Alpine.store('app'));
});

// Add more debugging
window.addEventListener('load', () => {
    console.log('Window loaded, store status:', Alpine.store('app'));
}); 