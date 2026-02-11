(function () {
    function updateNav() {
        var span = document.getElementById('user-span');
        var navAdmin = document.getElementById('nav-admin');
        if (window.auth && auth.isAdmin && navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
        if (!span) return;
        if (window.auth && auth.isLoggedIn && auth.isLoggedIn()) {
            var u = auth.getUser();
            span.innerHTML = 'Hi, ' + (u.first_name || u.email) + ' | <button type="button" id="logout-btn" class="btn btn-ghost">Log out</button>';
            var btn = document.getElementById('logout-btn');
            if (btn) btn.onclick = function () { auth.logout(); location.href = 'index.html'; };
        } else {
            span.innerHTML = '<button type="button" id="login-btn" class="btn btn-primary">Login / Register</button>';
            var loginBtn = document.getElementById('login-btn');
            if (loginBtn) loginBtn.onclick = function () { if (window.openAuthModal) openAuthModal({ onSuccess: load }); };
        }
    }

    function showGuest() {
        document.getElementById('profile-guest').style.display = 'block';
        document.getElementById('profile-content').style.display = 'none';
        document.getElementById('profile-login-btn').onclick = function () {
            openAuthModal({ onSuccess: function () { location.reload(); } });
        };
    }

    function showError(msg, isSuccess) {
        if (msg) {
            window.showToast(msg, isSuccess ? 'success' : 'error');
        }
        var el = document.getElementById('profile-error');
        if (el) {
            el.textContent = msg || '';
            el.style.display = msg ? 'block' : 'none';
            el.style.color = isSuccess ? '#2e7d32' : '#c62828';
        }
    }

    function load() {
        if (!(window.auth && auth.isLoggedIn && auth.isLoggedIn())) {
            showGuest();
            updateNav();
            return;
        }
        document.getElementById('profile-guest').style.display = 'none';
        document.getElementById('profile-content').style.display = 'block';
        showError('');


        window.api.fetchMe()
            .then(function (user) {

                var stored = auth.getUser();
                if (stored && user) {
                    user.token = stored.token;
                    localStorage.setItem('cinema_user', JSON.stringify(user));
                }
                renderProfile(user);
            })
            .catch(function () {
                renderProfile(auth.getUser());
            });

        loadBookings();
        loadReviews();
    }

    function renderProfile(u) {
        if (!u) return;
        document.getElementById('profile-info').innerHTML =
            `<strong>Name:</strong> ${u.first_name || ''} ${u.last_name || ''}<br>` +
            `<strong>Email:</strong> ${u.email || ''}<br>` +
            `<strong>Phone:</strong> ${u.phone_number || 'Not provided'}`;

        var adminLink = document.getElementById('profile-admin-link');
        if (adminLink) adminLink.style.display = (u.role === 'ADMIN') ? 'inline-block' : 'none';

        document.getElementById('user-balance').textContent = u.balance != null ? u.balance.toFixed(2) : '0.00';
    }

    function loadBookings() {
        window.api.fetchMyBookings()
            .then(function (data) {
                var list = document.getElementById('bookings-list');
                var empty = document.getElementById('bookings-empty');
                if (!list) return;
                list.innerHTML = '';
                if (!Array.isArray(data) || data.length === 0) {
                    empty.textContent = 'You have no bookings yet.';
                    empty.style.display = 'block';
                    return;
                }
                empty.style.display = 'none';
                data.forEach(function (ticket) {
                    var card = document.createElement('div');
                    card.className = 'card';
                    card.innerHTML = `
                        <div class="card-body">
                            <div class="card-title">${ticket.movie_title || 'Movie Booking'}</div>
                            <div class="card-meta">Row ${ticket.row_number}, Seat ${ticket.seat_number}</div>
                            <div class="card-meta" style="font-size: 0.8rem; color: var(--text-dim);">Ticket: ${ticket.type}</div>
                            <div class="badge-row">
                                <span class="badge ${ticket.status === 'PAID' ? 'badge-format' : ''}">${ticket.status}</span>
                            </div>
                            ${ticket.status !== 'CANCELLED' ? `
                                <button type="button" class="btn btn-ghost cancel-ticket-btn" data-id="${ticket.id}" style="margin-top: 1rem; color: #ef4444; padding-left: 0;">Cancel Booking</button>
                            ` : ''}
                        </div>
                    `;
                    list.appendChild(card);
                });
                list.querySelectorAll('.cancel-ticket-btn').forEach(function (btn) {
                    btn.onclick = function () {
                        var id = btn.dataset.id;
                        if (!id) return;
                        if (!confirm('Cancel this booking?')) return;
                        window.api.cancelBooking(id)
                            .then(function () { load(); })
                            .catch(function (err) { showError((err && err.message) || 'Cancel failed'); });
                    };
                });
            })
            .catch(function (err) {

                var empty = document.getElementById('bookings-empty');
                if (empty) {
                    empty.textContent = 'Error loading bookings. Please try again later.';
                    empty.style.display = 'block';
                    empty.style.color = '#ef4444';
                }
            });
    }

    function loadReviews() {
        window.api.fetchMyReviews()
            .then(function (reviews) {
                var list = document.getElementById('reviews-list');
                var empty = document.getElementById('reviews-empty');
                if (!list) return;
                list.innerHTML = '';
                if (!reviews || reviews.length === 0) {
                    empty.textContent = "You haven't written any reviews yet.";
                    empty.style.display = 'block';
                    return;
                }
                empty.style.display = 'none';
                reviews.forEach(function (r) {
                    var card = document.createElement('div');
                    card.className = 'card';
                    card.innerHTML = `
                        <div class="card-body">
                            <div class="card-title" style="color: var(--accent)">Rating: ${r.rating}/10</div>
                            <div class="card-meta" style="margin-bottom: 0.5rem; color: var(--text-dim); font-weight: 600;">Movie: ${r.movie_title || 'Unknown Movie'}</div>
                            <p class="card-meta" style="margin-top: 0.5rem; color: white;">"${r.comment}"</p>
                            <div style="margin-top: 1rem; border-top: 1px solid #333; padding-top: 0.75rem; display: flex; gap: 0.75rem;">
                                <button type="button" class="btn btn-ghost edit-review-btn" 
                                    data-id="${r.id}" 
                                    data-rating="${r.rating}" 
                                    data-comment="${r.comment.replace(/"/g, '&quot;')}"
                                    style="color: var(--accent);">Edit</button>
                                <button type="button" class="btn btn-ghost delete-review-btn" 
                                    data-id="${r.id}" 
                                    style="color: #ef4444;">Delete</button>
                            </div>
                        </div>
                    `;
                    list.appendChild(card);
                });

                list.querySelectorAll('.edit-review-btn').forEach(function (btn) {
                    btn.onclick = function () {
                        openEditModal({
                            id: btn.dataset.id,
                            rating: btn.dataset.rating,
                            comment: btn.dataset.comment
                        });
                    };
                });

                list.querySelectorAll('.delete-review-btn').forEach(function (btn) {
                    btn.onclick = function () {
                        var id = btn.dataset.id;
                        if (!confirm('Are you sure you want to delete this review?')) return;
                        window.api.deleteReview(id)
                            .then(function () { load(); })
                            .catch(function (err) { showError((err && err.message) || 'Delete failed'); });
                    };
                });
            })
            .catch(function (err) {

                var empty = document.getElementById('reviews-empty');
                if (empty) {
                    empty.textContent = 'Error loading review history. Please try again later.';
                    empty.style.display = 'block';
                    empty.style.color = '#ef4444';
                }
            });
    }

    function openEditModal(review) {
        document.getElementById('edit-review-id').value = review.id;
        document.getElementById('edit-rating').value = review.rating;
        document.getElementById('edit-comment').value = review.comment;
        var modal = document.getElementById('edit-review-modal');
        if (modal) modal.classList.add('is-open');
    }

    window.closeEditModal = function () {
        var modal = document.getElementById('edit-review-modal');
        if (modal) modal.classList.remove('is-open');
    };

    document.getElementById('edit-review-form').onsubmit = function (e) {
        e.preventDefault();
        var id = document.getElementById('edit-review-id').value;
        var payload = {
            rating: parseInt(document.getElementById('edit-rating').value),
            comment: document.getElementById('edit-comment').value
        };

        window.api.updateReview(id, payload)
            .then(function () {
                closeEditModal();
                load();
            })
            .catch(function (err) {
                alert((err && err.message) || 'Update failed');
            });
    };

    function openProfileModal(u) {
        document.getElementById('edit-first-name').value = u.first_name || '';
        document.getElementById('edit-last-name').value = u.last_name || '';
        document.getElementById('edit-email').value = u.email || '';
        document.getElementById('edit-phone').value = u.phone_number || '';
        var modal = document.getElementById('edit-profile-modal');
        if (modal) modal.classList.add('is-open');
    }

    window.closeProfileModal = function () {
        var modal = document.getElementById('edit-profile-modal');
        if (modal) modal.classList.remove('is-open');
    };

    document.getElementById('edit-profile-form').onsubmit = function (e) {
        e.preventDefault();
        var payload = {
            first_name: document.getElementById('edit-first-name').value,
            last_name: document.getElementById('edit-last-name').value,
            email: document.getElementById('edit-email').value,
            phone_number: document.getElementById('edit-phone').value
        };

        window.api.updateProfile(payload)
            .then(function () {
                closeProfileModal();
                load();
            })
            .catch(function (err) {
                alert((err && err.message) || 'Update failed');
            });
    };

    document.getElementById('edit-profile-btn').onclick = function () {
        var u = window.auth.getUser();
        if (u) openProfileModal(u);
    };

    function init() {
        updateNav();
        load();
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
