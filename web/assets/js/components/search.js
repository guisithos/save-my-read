function search() {
    return {
        isOpen: false,
        query: '',
        results: [],
        isLoading: false,
        addedBooks: new Set(),

        init() {
            this.$watch('isOpen', value => {
                if (value) {
                    document.body.style.overflow = 'hidden';
                } else {
                    document.body.style.overflow = '';
                    this.query = '';
                    this.results = [];
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
                const response = await api.searchBooks(this.query);
                this.results = response.data;
            } catch (error) {
                console.error('Search failed:', error);
            } finally {
                this.isLoading = false;
            }
        },

        async addBook(book) {
            try {
                await api.addBook({
                    googleBookId: book.id,
                    title: book.title,
                    authors: book.authors,
                    description: book.description,
                    categories: book.categories,
                    imageURL: book.imageURL,
                    status: 'TO_READ'
                });

                this.addedBooks.add(book.id);
                this.$dispatch('book-added', book);
            } catch (error) {
                console.error('Failed to add book:', error);
            }
        },

        isBookAdded(bookId) {
            return this.addedBooks.has(bookId);
        },

        close() {
            this.isOpen = false;
        }
    }
} 