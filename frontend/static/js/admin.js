(function () {
    var sessionsPage = 1;
    var moviesPage = 1;
    var bookingsPage = 1;
    var perPage = 5;

    var allSessions = [];
    var allMovies = [];
    var allBookings = [];

    function showError(msg) {
        var el = document.getElementById('admin-error');
        if (el) {
            el.textContent = msg || '';
            el.style.display = msg ? 'block' : 'none';
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

    function loadHalls() {
        window.api.adminFetchHalls()
            .then(function (data) {
                var ul = document.getElementById('admin-halls');
                var select = document.querySelector('#form-session select[name="hall_id"]');
                ul.innerHTML = '';
                if (select) select.innerHTML = '<option value="">Select Hall</option>';
                data.forEach(function (h) {
                    var li = document.createElement('li');
                    li.className = 'admin-list-item';
                    li.innerHTML = `
            <div class="item-info">
              <strong>${h.name}</strong> — ${h.location}<br>
              <small>Rows: ${h.total_rows}, Seats/Row: ${h.seats_per_row}</small>
            </div>
            <button type="button" class="btn btn-danger btn-sm del-hall" data-id="${h.id}">Delete</button>
          `;
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
                            window.api.adminDeleteHall(id).then(loadHalls).catch(function (err) { showError(err.message); });
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

        // Update all dropdown regardless of pagination for movies list
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
            li.innerHTML = `
                <div class="item-info">
                  <strong>${m.name}</strong> (${m.duration} min)<br>
                  <small>ID: ${m.id}</small>
                </div>
                <button type="button" class="btn btn-danger btn-sm del-movie" data-id="${m.id}">Delete</button>
            `;
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
                    window.api.adminDeleteMovie(id).then(loadMovies).catch(function (err) { showError(err.message); });
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
                <button type="button" class="btn btn-danger btn-sm del-session" data-id="${s.id}">Delete</button>
            `;
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
                    window.api.adminDeleteSession(id).then(loadSessions).catch(function (err) { showError(err.message); });
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
            var payload = {
                name: this.name.value.trim(),
                location: this.location.value.trim(),
                total_rows: parseInt(this.total_rows.value, 10) || 0,
                seats_per_row: parseInt(this.seats_per_row.value, 10) || 0
            };
            window.api.adminCreateHall(payload)
                .then(function () { e.target.reset(); loadHalls(); showError(''); })
                .catch(function (err) { showError(err.message); });
        };

        document.getElementById('form-movie').onsubmit = function (e) {
            e.preventDefault();
            var payload = {
                name: this.name.value.trim(),
                duration: parseInt(this.duration.value, 10) || 0,
                description: this.description.value.trim(),
                poster_url: this.poster_url.value.trim(),
                rating: parseFloat(this.rating.value) || 0
            };
            window.api.adminCreateMovie(payload)
                .then(function () { e.target.reset(); loadMovies(); showError(''); })
                .catch(function (err) { showError(err.message); });
        };

        document.getElementById('form-session').onsubmit = function (e) {
            e.preventDefault();
            var start = new Date(this.start_time.value);
            var payload = {
                movie_id: this.movie_id.value,
                hall_id: this.hall_id.value,
                start_time: start.toISOString(),
                price: parseFloat(this.price.value) || 0
            };
            if (!payload.movie_id || !payload.hall_id) {
                showError('Please select a movie and a hall');
                return;
            }
            window.api.adminCreateSession(payload)
                .then(function () { e.target.reset(); loadSessions(); showError(''); })
                .catch(function (err) { showError(err.message); });
        };
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
