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
                    // Get query from store if exists
                    if (this.$store.app?.searchQuery) {
                        this.query = this.$store.app.searchQuery;
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
            if (!this.query.trim()) {
                this.results = [];
                return;
            }

            this.isLoading = true;
            try {
                const response = await fetch(`/api/books/search?q=${encodeURIComponent(this.query)}`);
                if (!response.ok) throw new Error('Search failed');
                
                const data = await response.json();
                this.results = data.items.slice(0, 10).map(book => ({
                    id: book.id,
                    title: book.volumeInfo.title,
                    authors: book.volumeInfo.authors || [],
                    categories: book.volumeInfo.categories || [],
                    description: book.volumeInfo.description || '',
                    imageURL: book.volumeInfo.imageLinks?.thumbnail || '/assets/images/no-cover.png'
                }));
            } catch (error) {
                console.error('Search failed:', error);
                this.$dispatch('notification', {
                    type: 'error',
                    message: 'Failed to search books'
                });
                this.results = [];
            } finally {
                this.isLoading = false;
            }
        },

        async addBook(book) {
            try {
                await fetch('/api/books', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        googleBookId: book.id,
                        title: book.title,
                        authors: book.authors,
                        description: book.description,
                        categories: book.categories,
                        imageURL: book.imageURL,
                        status: 'TO_READ'
                    })
                });

                this.addedBooks.add(book.id);
                this.$dispatch('notification', {
                    type: 'success',
                    message: 'Book added successfully'
                });
            } catch (error) {
                console.error('Failed to add book:', error);
                this.$dispatch('notification', {
                    type: 'error',
                    message: 'Failed to add book'
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