(function () {
    var sessionsPage = 1;
    var moviesPage = 1;
    var bookingsPage = 1;
    var perPage = 5;

    var allSessions = [];
    var allMovies = [];
    var allBookings = [];
    var allGenres = [];

    function showError(msg) {
        if (msg) {
            console.error('Admin Error:', msg);
            window.showToast(typeof msg === 'string' ? msg : (msg.message || 'An unknown error occurred'), 'error');
        }
    }

    function updateNav() {
        var span = document.getElementById('user-span');
        if (!span) return;
        if (window.auth && auth.isLoggedIn()) {
            var u = auth.getUser();
            span.innerHTML = 'Hi, ' + (u.first_name || u.email) + ' | <button type="button" id="logout-btn" class="btn btn-ghost">Log out</button>';
            document.getElementById('logout-btn').onclick = function () {
                auth.logout();
                location.href = 'index.html';
            };
        } else {
            span.innerHTML = '<a href="index.html" class="btn btn-ghost">Back to Home</a>';
        }
    }

    function resetForm(form) {
        form.reset();
        if (form.id === 'form-movie') {
            form.querySelectorAll('input[type="checkbox"]').forEach(function (cb) { cb.checked = false; });
        }
        var idInput = form.querySelector('input[name="id"]');
        if (idInput) idInput.value = '';

        var submitBtn = form.querySelector('button[type="submit"]');
        var cancelBtn = form.querySelector('.cancel-btn');
        if (submitBtn) {
            var formId = form.getAttribute('id') || '';
            var type = formId.replace('form-', '');
            submitBtn.textContent = 'Add ' + type.charAt(0).toUpperCase() + type.slice(1);
        }
        if (cancelBtn) cancelBtn.style.display = 'none';
    }

    function startEdit(type, data) {
        var form = document.getElementById('form-' + type);
        if (!form) return;

        form.reset();
        var idInput = form.querySelector('input[name="id"]');
        if (idInput) idInput.value = data.id;

        if (type === 'hall') {
            form.name.value = data.name;
            form.location.value = data.location;
            form.type.value = data.type || 'STANDARD';
            form.total_rows.value = data.total_rows;
            form.seats_per_row.value = data.seats_per_row;
        } else if (type === 'movie') {
            form.name.value = data.name;
            form.duration.value = data.duration;
            form.rating.value = data.rating;
            form.description.value = data.description || '';
            form.poster_url.value = data.poster_url || '';
            form.trailer_url.value = data.trailer_url || '';
            form.age_limit.value = data.age_limit || 0;
            form.is_coming_soon.checked = !!data.is_coming_soon;

            var genreIds = data.genre_ids || [];
            form.querySelectorAll('input[name="genre"]').forEach(function (cb) {
                cb.checked = genreIds.indexOf(cb.value) !== -1;
            });
        } else if (type === 'session') {
            form.movie_id.value = data.movie_id;
            form.hall_id.value = data.hall_id;
            if (data.start_time) {
                var dt = new Date(data.start_time);
                var pad = function (n) { return n < 10 ? '0' + n : n; };
                var localISO = dt.getFullYear() + '-' + pad(dt.getMonth() + 1) + '-' + pad(dt.getDate()) + 'T' + pad(dt.getHours()) + ':' + pad(dt.getMinutes());
                form.start_time.value = localISO;
            }
            form.price.value = data.price;
        } else if (type === 'genre') {
            form.name.value = data.name;
        }

        var submitBtn = form.querySelector('button[type="submit"]');
        var cancelBtn = form.querySelector('.cancel-btn');
        if (submitBtn) submitBtn.textContent = 'Update ' + type.charAt(0).toUpperCase() + type.slice(1);
        if (cancelBtn) cancelBtn.style.display = 'inline-flex';

        window.scrollTo({ top: form.offsetTop - 100, behavior: 'smooth' });
    }

    function loadGenres() {
        return window.api.adminFetchGenres()
            .then(function (data) {
                allGenres = data || [];
                renderGenreCheckboxes();
                renderGenreList();
            })
            .catch(function (err) { showError(err.message); });
    }

    function renderGenreCheckboxes() {
        var container = document.getElementById('movie-genres-selection');
        if (!container) return;
        container.innerHTML = '';
        allGenres.forEach(function (g) {
            var label = document.createElement('label');
            label.className = 'admin-checkbox-label';
            label.innerHTML = `<input type="checkbox" name="genre" value="${g.id}"> ${g.name}`;
            container.appendChild(label);
        });
    }

    function renderGenreList() {
        var ul = document.getElementById('admin-genres');
        if (!ul) return;
        ul.innerHTML = '';
        allGenres.forEach(function (g) {
            var li = document.createElement('li');
            li.className = 'admin-list-item';
            li.innerHTML = `
                <div class="item-info"><strong>${g.name}</strong></div>
                <div class="item-actions">
                    <button type="button" class="btn btn-ghost btn-sm edit-genre">Edit</button>
                    <button type="button" class="btn btn-danger btn-sm del-genre" data-id="${g.id}">Delete</button>
                </div>
            `;
            li.querySelector('.edit-genre').onclick = function () { startEdit('genre', g); };
            ul.appendChild(li);
        });
        ul.querySelectorAll('.del-genre').forEach(function (btn) {
            btn.onclick = function () {
                var id = btn.dataset.id;
                if (confirm('Delete this genre?')) {
                    window.api.adminDeleteGenre(id)
                        .then(function () {
                            loadGenres();
                            showToast('Genre deleted successfully');
                        })
                        .catch(function (err) { showError(err.message); });
                }
            };
        });
    }

    function loadHalls() {
        window.api.adminFetchHalls()
            .then(function (data) {
                var ul = document.getElementById('admin-halls');
                var select = document.querySelector('#form-session select[name="hall_id"]');
                if (!ul) return;
                ul.innerHTML = '';
                if (select) select.innerHTML = '<option value="">Select Hall</option>';
                data.forEach(function (h) {
                    var li = document.createElement('li');
                    li.className = 'admin-list-item';
                    li.innerHTML = `
                        <div class="item-info">
                        <strong>${h.name}</strong> (${h.type || 'STANDARD'}) — ${h.location}<br>
                        <small>Rows: ${h.total_rows}, Seats/Row: ${h.seats_per_row}</small>
                        </div>
                        <div class="item-actions">
                            <button type="button" class="btn btn-ghost btn-sm edit-hall">Edit</button>
                            <button type="button" class="btn btn-danger btn-sm del-hall" data-id="${h.id}">Delete</button>
                        </div>
                    `;
                    li.querySelector('.edit-hall').onclick = function () { startEdit('hall', h); };
                    ul.appendChild(li);

                    if (select) {
                        var opt = document.createElement('option');
                        opt.value = h.id;
                        opt.textContent = h.name + ' (' + h.location + ')';
                        select.appendChild(opt);
                    }
                });
                ul.querySelectorAll('.del-hall').forEach(function (btn) {
                    btn.onclick = function () {
                        var id = btn.dataset.id;
                        if (confirm('Delete this hall?')) {
                            window.api.adminDeleteHall(id)
                                .then(function () {
                                    loadHalls();
                                    showToast('Hall deleted successfully');
                                })
                                .catch(function (err) { showError(err.message); });
                        }
                    };
                });
            })
            .catch(function (err) { showError(err.message); });
    }

    function loadMovies() {
        window.api.fetchMovies()
            .then(function (data) {
                allMovies = data || [];
                renderMovies();
            })
            .catch(function (err) { showError(err.message); });
    }

    function renderMovies() {
        var ul = document.getElementById('admin-movies');
        var pageInfo = document.getElementById('movies-page-info');
        var select = document.querySelector('#form-session select[name="movie_id"]');
        if (!ul || !pageInfo) return;

        ul.innerHTML = '';
        if (select) select.innerHTML = '<option value="">Select Movie</option>';

        allMovies.forEach(function (m) {
            if (select) {
                var opt = document.createElement('option');
                opt.value = m.id;
                opt.textContent = m.name;
                select.appendChild(opt);
            }
        });

        var startIdx = (moviesPage - 1) * perPage;
        var endIdx = startIdx + perPage;
        var pageData = allMovies.slice(startIdx, endIdx);

        pageData.forEach(function (m) {
            var li = document.createElement('li');
            li.className = 'admin-list-item';
            var genreNames = Array.isArray(m.genres) ? m.genres.join(', ') : '';
            li.innerHTML = `
                <div class="item-info">
                  <strong>${m.name}</strong> (${m.duration} min)<br>
                  <small>Genres: ${genreNames || 'None'}</small><br>
                  <small>ID: ${m.id}</small>
                </div>
                <div class="item-actions">
                    <button type="button" class="btn btn-ghost btn-sm edit-movie">Edit</button>
                    <button type="button" class="btn btn-danger btn-sm del-movie" data-id="${m.id}">Delete</button>
                </div>
            `;
            li.querySelector('.edit-movie').onclick = function () { startEdit('movie', m); };
            ul.appendChild(li);
        });

        var totalPages = Math.ceil(allMovies.length / perPage) || 1;
        pageInfo.textContent = 'Page ' + moviesPage + ' of ' + totalPages;

        document.getElementById('prev-movies').disabled = moviesPage === 1;
        document.getElementById('next-movies').disabled = moviesPage >= totalPages;

        ul.querySelectorAll('.del-movie').forEach(function (btn) {
            btn.onclick = function () {
                var id = btn.dataset.id;
                if (confirm('Delete this movie?')) {
                    window.api.adminDeleteMovie(id)
                        .then(function () {
                            loadMovies();
                            showToast('Movie deleted successfully');
                        })
                        .catch(function (err) { showError(err.message); });
                }
            };
        });
    }

    function loadSessions() {
        window.api.fetchUpcomingSessions()
            .then(function (data) {
                allSessions = data || [];
                renderSessions();
            })
            .catch(function (err) { showError(err.message); });
    }

    function renderSessions() {
        var ul = document.getElementById('admin-sessions');
        var pageInfo = document.getElementById('sessions-page-info');
        if (!ul || !pageInfo) return;

        ul.innerHTML = '';
        var startIdx = (sessionsPage - 1) * perPage;
        var endIdx = startIdx + perPage;
        var pageData = allSessions.slice(startIdx, endIdx);

        pageData.forEach(function (s) {
            var start = s.start_time ? new Date(s.start_time).toLocaleString() : '';
            var li = document.createElement('li');
            li.className = 'admin-list-item';
            li.innerHTML = `
                <div class="item-info">
                  <strong>${start}</strong> — Price: ${s.price} ₸<br>
                  <small>Movie: ${s.movie_id}, Hall: ${s.hall_id}</small>
                </div>
                <div class="item-actions">
                    <button type="button" class="btn btn-ghost btn-sm edit-session">Edit</button>
                    <button type="button" class="btn btn-danger btn-sm del-session" data-id="${s.id}">Delete</button>
                </div>
            `;
            li.querySelector('.edit-session').onclick = function () { startEdit('session', s); };
            ul.appendChild(li);
        });

        var totalPages = Math.ceil(allSessions.length / perPage) || 1;
        pageInfo.textContent = 'Page ' + sessionsPage + ' of ' + totalPages;

        document.getElementById('prev-sessions').disabled = sessionsPage === 1;
        document.getElementById('next-sessions').disabled = sessionsPage >= totalPages;

        ul.querySelectorAll('.del-session').forEach(function (btn) {
            btn.onclick = function () {
                var id = btn.dataset.id;
                if (confirm('Delete this session?')) {
                    window.api.adminDeleteSession(id)
                        .then(function () {
                            loadSessions();
                            showToast('Session deleted successfully');
                        })
                        .catch(function (err) { showError(err.message); });
                }
            };
        });
    }

    function loadBookings() {
        window.api.adminFetchBookings()
            .then(function (data) {
                allBookings = data || [];
                renderBookings();
            })
            .catch(function (err) { showError(err.message); });
    }

    function renderBookings() {
        var ul = document.getElementById('admin-bookings');
        var pageInfo = document.getElementById('bookings-page-info');
        if (!ul || !pageInfo) return;

        ul.innerHTML = '';
        var startIdx = (bookingsPage - 1) * perPage;
        var endIdx = startIdx + perPage;
        var pageData = allBookings.slice(startIdx, endIdx);

        pageData.forEach(function (b) {
            var li = document.createElement('li');
            li.className = 'admin-list-item';
            li.innerHTML = `
                <div class="item-info">
                  <strong>Session: ${b.session_id}</strong> — Row ${b.row_number}, Seat ${b.seat_number}<br>
                  <small>Status: ${b.status} | User: ${b.user_id}</small>
                </div>
            `;
            ul.appendChild(li);
        });

        var totalPages = Math.ceil(allBookings.length / perPage) || 1;
        pageInfo.textContent = 'Page ' + bookingsPage + ' of ' + totalPages;

        document.getElementById('prev-bookings').disabled = bookingsPage === 1;
        document.getElementById('next-bookings').disabled = bookingsPage >= totalPages;
    }

    function init() {
        updateNav();
        if (!window.auth || !auth.isLoggedIn() || !auth.isAdmin()) {
            location.replace('index.html');
            return;
        }
        document.getElementById('admin-forbidden').style.display = 'none';
        document.getElementById('admin-content').style.display = 'block';

        loadHalls();
        loadMovies();
        loadSessions();
        loadBookings();
        loadGenres();

        document.querySelectorAll('.cancel-btn').forEach(function (btn) {
            btn.onclick = function () {
                resetForm(btn.closest('form'));
            };
        });

        document.getElementById('prev-sessions').onclick = function () {
            if (sessionsPage > 1) {
                sessionsPage--;
                renderSessions();
            }
        };
        document.getElementById('next-sessions').onclick = function () {
            var totalPages = Math.ceil(allSessions.length / perPage);
            if (sessionsPage < totalPages) {
                sessionsPage++;
                renderSessions();
            }
        };
        document.getElementById('prev-movies').onclick = function () {
            if (moviesPage > 1) {
                moviesPage--;
                renderMovies();
            }
        };
        document.getElementById('next-movies').onclick = function () {
            var totalPages = Math.ceil(allMovies.length / perPage);
            if (moviesPage < totalPages) {
                moviesPage++;
                renderMovies();
            }
        };
        document.getElementById('prev-bookings').onclick = function () {
            if (bookingsPage > 1) {
                bookingsPage--;
                renderBookings();
            }
        };
        document.getElementById('next-bookings').onclick = function () {
            var totalPages = Math.ceil(allBookings.length / perPage);
            if (bookingsPage < totalPages) {
                bookingsPage++;
                renderBookings();
            }
        };

        document.getElementById('form-hall').onsubmit = function (e) {
            e.preventDefault();
            var idInput = this.querySelector('input[name="id"]');
            var id = idInput ? idInput.value : '';
            var payload = {
                name: this.name.value.trim(),
                location: this.location.value.trim(),
                type: this.type.value,
                total_rows: parseInt(this.total_rows.value, 10) || 0,
                seats_per_row: parseInt(this.seats_per_row.value, 10) || 0
            };
            var promise = id
                ? window.api.adminUpdateHall(id, payload)
                : window.api.adminCreateHall(payload);

            promise.then(function () {
                resetForm(e.target);
                loadHalls();
                showToast(id ? 'Hall updated successfully!' : 'Hall created successfully!');
            }).catch(function (err) { showError(err.message); });
        };

        document.getElementById('form-movie').onsubmit = function (e) {
            e.preventDefault();
            var idInput = this.querySelector('input[name="id"]');
            var id = idInput ? idInput.value : '';
            var genresList = [];
            this.querySelectorAll('input[name="genre"]:checked').forEach(function (cb) {
                genresList.push(cb.value);
            });

            var payload = {
                name: this.name.value.trim(),
                duration: parseInt(this.duration.value, 10) || 0,
                description: this.description.value.trim(),
                poster_url: this.poster_url.value.trim(),
                trailer_url: this.trailer_url ? this.trailer_url.value.trim() : '',
                age_limit: parseInt(this.age_limit ? this.age_limit.value : 0, 10) || 0,
                rating: parseFloat(this.rating.value) || 0,
                genre_ids: genresList,
                is_coming_soon: this.is_coming_soon ? this.is_coming_soon.checked : false
            };
            var promise = id
                ? window.api.adminUpdateMovie(id, payload)
                : window.api.adminCreateMovie(payload);

            promise.then(function () {
                resetForm(e.target);
                loadMovies();
                showToast(id ? 'Movie updated successfully!' : 'Movie added successfully!');
            }).catch(function (err) { showError(err.message); });
        };

        document.getElementById('form-genre').onsubmit = function (e) {
            e.preventDefault();
            var idInput = this.querySelector('input[name="id"]');
            var id = idInput ? idInput.value : '';
            var payload = { name: this.name.value.trim() };
            var promise = id
                ? window.api.adminUpdateGenre(id, payload)
                : window.api.adminCreateGenre(payload);

            promise.then(function () {
                resetForm(e.target);
                loadGenres();
                showToast(id ? 'Genre updated successfully!' : 'Genre added successfully!');
            }).catch(function (err) { showError(err.message); });
        };

        document.getElementById('form-session').onsubmit = function (e) {
            e.preventDefault();
            var idInput = this.querySelector('input[name="id"]');
            var id = idInput ? idInput.value : '';
            var startVal = this.start_time.value;
            if (!startVal) { showError('Please select a start time'); return; }
            var start = new Date(startVal);
            if (isNaN(start.getTime())) { showError('Invalid start time'); return; }
            var payload = {
                movie_id: this.movie_id.value,
                hall_id: this.hall_id.value,
                start_time: start.toISOString(),
                price: parseFloat(this.price.value) || 0
            };
            if (!payload.movie_id || !payload.hall_id) { showError('Please select a movie and a hall'); return; }
            var promise = id
                ? window.api.adminUpdateSession(id, payload)
                : window.api.adminCreateSession(payload);

            promise.then(function () {
                resetForm(e.target);
                loadSessions();
                showToast(id ? 'Session updated successfully!' : 'Session scheduled successfully!');
            }).catch(function (err) { showError(err.message); });
        };
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
