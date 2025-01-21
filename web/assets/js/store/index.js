// Initialize Alpine.js store
document.addEventListener('alpine:init', () => {
    console.log('Initializing Alpine store');
    
    if (typeof createAppStore !== 'function') {
        console.error('createAppStore is not defined! Check if app-store.js is loaded.');
        return;
    }

    // Create and initialize store
    const store = createAppStore();
    Alpine.store('app', store);
    
    // Initialize store after creation
    if (typeof store.init === 'function') {
        store.init();
    }

    console.log('Store created and initialized:', Alpine.store('app'));
});

// Verify store initialization
document.addEventListener('alpine:initialized', () => {
    const store = Alpine.store('app');
    console.log('Alpine initialized, verifying store:', store);
    if (!store) {
        console.error('Store not properly initialized!');
    }
});

// Add more debugging
window.addEventListener('load', () => {
    console.log('Window loaded, store status:', Alpine.store('app'));
}); 