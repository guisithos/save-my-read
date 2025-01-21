function search() {
    return {
        isOpen: false,
        query: '',
        results: [],
        isLoading: false,
        addedBooks: new Set(),

        init() {
            console.log('Search component initialized');
            
            // Watch store changes
            this.$watch('$store.app.showSearch', (value) => {
                console.log('Store showSearch changed:', value);
                if (value) {
                    document.body.style.overflow = 'hidden';
                    this.isOpen = true;
                    const store = Alpine.store('app');
                    if (store?.searchQuery) {
                        console.log('Setting query from store:', store.searchQuery);
                        this.query = store.searchQuery;
                        this.handleSearch();
                    }
                } else {
                    document.body.style.overflow = '';
                    this.isOpen = false;
                    this.query = '';
                    this.results = [];
                }
            });

            // Listen for search events
            window.addEventListener('open-search', (event) => {
                console.log('Received open-search event:', event);
                if (event.detail?.query) {
                    this.query = event.detail.query;
                    this.handleSearch();
                }
            });
        },

        async handleSearch() {
            console.log('Handling search for query:', this.query);
            
            if (!this.query.trim()) {
                console.log('Empty query, clearing results');
                this.results = [];
                return;
            }

            this.isLoading = true;
            try {
                console.log('Fetching search results...');
                const response = await fetch(`/api/books/search?q=${encodeURIComponent(this.query)}`, {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });
                
                if (!response.ok) {
                    throw new Error(`Search failed with status: ${response.status}`);
                }
                
                const data = await response.json();
                console.log('Search API response:', data);

                if (!Array.isArray(data)) {
                    console.error('Unexpected response format:', data);
                    throw new Error('Invalid response format from server');
                }

                this.results = data.slice(0, 10).map(book => ({
                    id: book.id,
                    title: book.title,
                    authors: book.authors || [],
                    categories: book.categories || [],
                    description: book.description || '',
                    imageURL: book.imageURL || '/assets/images/no-cover.png'
                }));
                console.log('Processed search results:', this.results);
            } catch (error) {
                console.error('Search failed:', error);
                this.$dispatch('notification', {
                    type: 'error',
                    message: 'Failed to search books: ' + error.message
                });
                this.results = [];
            } finally {
                this.isLoading = false;
                console.log('Search complete. Results:', this.results.length, 'Loading:', this.isLoading);
            }
        },

        async addBook(book) {
            try {
                const response = await fetch('/api/books', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    },
                    body: JSON.stringify({
                        id: book.id,
                        title: book.title,
                        authors: book.authors,
                        description: book.description,
                        categories: book.categories,
                        imageURL: book.imageURL,
                        status: 'TO_READ'
                    })
                });

                if (!response.ok) {
                    throw new Error(`Failed to add book: ${response.statusText}`);
                }

                this.addedBooks.add(book.id);
                this.$dispatch('notification', {
                    type: 'success',
                    message: 'Book added successfully'
                });
            } catch (error) {
                console.error('Failed to add book:', error);
                this.$dispatch('notification', {
                    type: 'error',
                    message: 'Failed to add book: ' + error.message
                });
            }
        },

        isBookAdded(bookId) {
            return this.addedBooks.has(bookId);
        },

        close() {
            if (this.$store.app) {
                this.$store.app.toggleSearch(false);
            }
        }
    }
} 