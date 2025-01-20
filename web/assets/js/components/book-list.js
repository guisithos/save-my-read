function bookList() {
    return {
        books: [],
        isLoading: true,
        currentStatus: 'ALL',
        statusModal: {
            show: false,
            book: null
        },
        statuses: [
            { value: 'ALL', label: 'All Books', icon: 'fas fa-books' },
            { value: 'TO_READ', label: 'To Read', icon: 'fas fa-bookmark' },
            { value: 'READING', label: 'Reading', icon: 'fas fa-book-reader' },
            { value: 'COMPLETED', label: 'Completed', icon: 'fas fa-check-circle' },
            { value: 'DNF', label: 'Did Not Finish', icon: 'fas fa-times-circle' }
        ],

        async init() {
            try {
                const response = await api.getBooks();
                this.books = response.data;
            } catch (error) {
                console.error('Failed to fetch books:', error);
            } finally {
                this.isLoading = false;
            }
        },

        get filteredBooks() {
            if (this.currentStatus === 'ALL') {
                return this.books;
            }
            return this.books.filter(book => book.status === this.currentStatus);
        },

        getBookCountByStatus(status) {
            if (status === 'ALL') {
                return this.books.length;
            }
            return this.books.filter(book => book.status === status).length;
        },

        getStatusLabel(status) {
            return this.statuses.find(s => s.value === status)?.label || status;
        },

        openStatusModal(book) {
            this.statusModal.book = book;
            this.statusModal.show = true;
        },

        closeStatusModal() {
            this.statusModal.show = false;
            this.statusModal.book = null;
        },

        async updateBookStatus(bookId, newStatus) {
            try {
                await api.updateBookStatus(bookId, newStatus);
                this.books = this.books.map(book => 
                    book.id === bookId 
                        ? { ...book, status: newStatus }
                        : book
                );
                this.closeStatusModal();
            } catch (error) {
                console.error('Failed to update book status:', error);
            }
        },

        async removeBook(bookId) {
            if (!confirm('Are you sure you want to remove this book?')) {
                return;
            }

            try {
                await api.deleteBook(bookId);
                this.books = this.books.filter(book => book.id !== bookId);
            } catch (error) {
                console.error('Failed to remove book:', error);
            }
        }
    }
} 