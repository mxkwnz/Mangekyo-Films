(function () {
    var API_BASE = window.API_BASE || '/api';
    var ADMIN_BASE = API_BASE + '/admin';

    function authHeaders() {
        return window.auth && window.auth.getAuthHeaders ? window.auth.getAuthHeaders() : {};
    }

    function showError(msg) {
        var el = document.getElementById('admin-error');
        if (!el) return;
        el.textContent = msg || '';
        el.style.display = msg ? 'block' : 'none';
    }

    function updateNav() {
        var span = document.getElementById('user-span');
        if (!span) return;
        
        if (window.auth && auth.isLoggedIn()) {
            var u = auth.getUser();
            span.innerHTML =
                'Hi, ' + (u.first_name || u.email) +
                ' | <button type="button" id="logout-btn" class="btn btn-ghost btn-sm">Log out</button>';
            var btn = document.getElementById('logout-btn');
            if (btn) {
                btn.onclick = function () {
                    auth.logout();
                    location.href = 'index.html';
                };
            }
        } else {
            span.innerHTML = '<a href="index.html">Back to Home</a>';
        }
    }

    // ========== HALLS ==========
    
    function loadHalls() {
        fetch(ADMIN_BASE + '/halls', { headers: authHeaders() })
            .then(function (r) { return r.json(); })
            .then(function (data) {
                var ul = document.getElementById('admin-halls');
                if (!ul) return;
                
                ul.innerHTML = '';
                var halls = Array.isArray(data) ? data : [];
                
                if (halls.length === 0) {
                    var li = document.createElement('li');
                    li.textContent = 'No halls yet.';
                    li.style.color = '#666';
                    ul.appendChild(li);
                    return;
                }
                
                halls.forEach(function (h) {
                    var li = document.createElement('li');
                    
                    var info = document.createElement('span');
                    info.textContent =
                        (h.name || '') + ' – ' + (h.location || '') +
                        ' (Rows: ' + (h.total_rows || 0) + ', Seats/row: ' + (h.seats_per_row || 0) + ')';
                    
                    var delBtn = document.createElement('button');
                    delBtn.type = 'button';
                    delBtn.className = 'btn btn-danger btn-sm';
                    delBtn.textContent = 'Delete';
                    delBtn.onclick = function () {
                        if (!confirm('Delete hall "' + (h.name || '') + '"?')) return;
                        deleteHall(h.id);
                    };
                    
                    li.appendChild(info);
                    li.appendChild(delBtn);
                    ul.appendChild(li);
                });
            })
            .catch(function () {
                showError('Failed to load halls');
            });
    }

    function deleteHall(id) {
        if (!id) return;
        fetch(ADMIN_BASE + '/halls/' + id, {
            method: 'DELETE',
            headers: authHeaders()
        })
            .then(function (r) { return r.json(); })
            .then(function () {
                loadHalls();
                showError('');
            })
            .catch(function () {
                showError('Failed to delete hall');
            });
    }

    // ========== MOVIES ==========
    
    function loadMovies() {
        fetch(API_BASE + '/movies')
            .then(function (r) { return r.json(); })
            .then(function (data) {
                var ul = document.getElementById('admin-movies');
                if (!ul) return;
                
                ul.innerHTML = '';
                var movies = Array.isArray(data) ? data : [];
                
                if (movies.length === 0) {
                    var li = document.createElement('li');
                    li.textContent = 'No movies yet.';
                    li.style.color = '#666';
                    ul.appendChild(li);
                    return;
                }
                
                movies.forEach(function (m) {
                    var li = document.createElement('li');
                    
                    var info = document.createElement('span');
                    info.textContent =
                        (m.name || '') + ' (' + (m.year || '') + ') – ' +
                        (m.duration || 0) + ' min – ID: ' + (m.id || '');
                    
                    var delBtn = document.createElement('button');
                    delBtn.type = 'button';
                    delBtn.className = 'btn btn-danger btn-sm';
                    delBtn.textContent = 'Delete';
                    delBtn.onclick = function () {
                        if (!confirm('Delete movie "' + (m.name || '') + '"?')) return;
                        deleteMovie(m.id);
                    };
                    
                    li.appendChild(info);
                    li.appendChild(delBtn);
                    ul.appendChild(li);
                });
            })
            .catch(function () {
                showError('Failed to load movies');
            });
    }

    function deleteMovie(id) {
        if (!id) return;
        fetch(ADMIN_BASE + '/movies/' + id, {
            method: 'DELETE',
            headers: authHeaders()
        })
            .then(function (r) { return r.json(); })
            .then(function () {
                loadMovies();
                showError('');
            })
            .catch(function () {
                showError('Failed to delete movie');
            });
    }

    function loadSessions() {
        fetch(API_BASE + '/sessions/upcoming')
            .then(function (r) { return r.json(); })
            .then(function (data) {
                var ul = document.getElementById('admin-sessions');
                if (!ul) return;
                
                ul.innerHTML = '';
                var sessions = Array.isArray(data) ? data : [];
                
                if (sessions.length === 0) {
                    var li = document.createElement('li');
                    li.textContent = 'No upcoming sessions.';
                    li.style.color = '#666';
                    ul.appendChild(li);
                    return;
                }
                
                sessions.forEach(function (s) {
                    var li = document.createElement('li');
                    
                    var start = s.start_time ? new Date(s.start_time).toLocaleString() : '';
                    
                    var info = document.createElement('span');
                    info.textContent =
                        'ID: ' + (s.id || '') + ' – Movie: ' + (s.movie_id || '') +
                        ', Hall: ' + (s.hall_id || '') + ' – ' + start +
                        ' ($' + (s.price != null ? s.price : 0) + ')';
                    
                    var delBtn = document.createElement('button');
                    delBtn.type = 'button';
                    delBtn.className = 'btn btn-danger btn-sm';
                    delBtn.textContent = 'Delete';
                    delBtn.onclick = function () {
                        if (!confirm('Delete this session?')) return;
                        deleteSession(s.id);
                    };
                    
                    li.appendChild(info);
                    li.appendChild(delBtn);
                    ul.appendChild(li);
                });
            })
            .catch(function () {
                showError('Failed to load sessions');
            });
    }

    function deleteSession(id) {
        if (!id) return;
        fetch(ADMIN_BASE + '/sessions/' + id, {
            method: 'DELETE',
            headers: authHeaders()
        })
            .then(function (r) { return r.json(); })
            .then(function () {
                loadSessions();
                showError('');
            })
            .catch(function () {
                showError('Failed to delete session');
            });
    }

    // ========== BOOKINGS ==========
    
    function loadBookings() {
        fetch(ADMIN_BASE + '/bookings', { headers: authHeaders() })
            .then(function (r) { return r.json(); })
            .then(function (data) {
                var ul = document.getElementById('admin-bookings');
                var noBookings = document.getElementById('no-bookings');
                if (!ul) return;
                
                ul.innerHTML = '';
                var bookings = Array.isArray(data) ? data : [];
                
                if (bookings.length === 0) {
                    if (noBookings) noBookings.style.display = 'block';
                    return;
                }
                
                if (noBookings) noBookings.style.display = 'none';
                
                bookings.forEach(function (b) {
                    var li = document.createElement('li');
                    li.textContent =
                        'Session: ' + (b.session_id || '') +
                        ' – Row ' + (b.row_number || '') + ', Seat ' + (b.seat_number || '') +
                        ' – ' + (b.status || 'confirmed') +
                        ' – User ID: ' + (b.user_id || '');
                    ul.appendChild(li);
                });
            })
            .catch(function () {
                showError('Failed to load bookings');
            });
    }

    // ========== FORM HANDLERS ==========
    
    function setupHallForm() {
        var form = document.getElementById('form-hall');
        if (!form) return;
        
        form.onsubmit = function (e) {
            e.preventDefault();
            
            var payload = {
                name: form.name.value.trim(),
                location: form.location.value.trim(),
                total_rows: parseInt(form.total_rows.value, 10) || 0,
                seats_per_row: parseInt(form.seats_per_row.value, 10) || 0
            };

            fetch(ADMIN_BASE + '/halls', {
                method: 'POST',
                headers: authHeaders(),
                body: JSON.stringify(payload)
            })
                .then(function (r) {
                    return r.json().then(function (d) {
                        return { ok: r.ok, data: d };
                    });
                })
                .then(function (res) {
                    if (res.ok) {
                        form.reset();
                        loadHalls();
                        showError('');
                    } else {
                        showError(res.data.error || 'Error creating hall');
                    }
                })
                .catch(function () {
                    showError('Network error');
                });
        };
    }

    function setupMovieForm() {
        var form = document.getElementById('form-movie');
        if (!form) return;
        
        form.onsubmit = function (e) {
            e.preventDefault();
            
            var payload = {
                name: form.name.value.trim(),
                year: parseInt(form.year.value, 10) || 2024,
                duration: parseInt(form.duration.value, 10) || 0,
                description: form.description.value.trim(),
                genres: form.genres.value.trim(),
                age_restriction: form.age_restriction.value.trim(),
                poster_url: form.poster_url.value.trim(),
                trailer_url: form.trailer_url.value.trim(),
                rating: parseFloat(form.rating.value) || 0
            };

            fetch(ADMIN_BASE + '/movies', {
                method: 'POST',
                headers: authHeaders(),
                body: JSON.stringify(payload)
            })
                .then(function (r) {
                    return r.json().then(function (d) {
                        return { ok: r.ok, data: d };
                    });
                })
                .then(function (res) {
                    if (res.ok) {
                        form.reset();
                        loadMovies();
                        showError('');
                    } else {
                        showError(res.data.error || 'Error creating movie');
                    }
                })
                .catch(function () {
                    showError('Network error');
                });
        };
    }

    function setupSessionForm() {
        var form = document.getElementById('form-session');
        if (!form) return;
        
        form.onsubmit = function (e) {
            e.preventDefault();
            
            var startStr = form.start_time.value;
            if (!startStr) {
                showError('Enter date and time');
                return;
            }
            
            var start = new Date(startStr);
            
            var payload = {
                movie_id: form.movie_id.value.trim(),
                hall_id: form.hall_id.value.trim(),
                start_time: start.toISOString(),
                price: parseFloat(form.price.value) || 0,
                format: form.format.value,
                language: form.language.value,
                age_restriction: form.age_restriction.value.trim()
            };

            fetch(ADMIN_BASE + '/sessions', {
                method: 'POST',
                headers: authHeaders(),
                body: JSON.stringify(payload)
            })
                .then(function (r) {
                    return r.json().then(function (d) {
                        return { ok: r.ok, data: d };
                    });
                })
                .then(function (res) {
                    if (res.ok) {
                        form.reset();
                        loadSessions();
                        showError('');
                    } else {
                        showError(res.data.error || 'Error creating session');
                    }
                })
                .catch(function () {
                    showError('Network error');
                });
        };
    }

    // ========== INIT ==========
    
    function init() {
        updateNav();

        // Check if user is admin
        if (!window.auth || !auth.isLoggedIn() || !auth.isAdmin()) {
            document.getElementById('admin-forbidden').style.display = 'block';
            document.getElementById('admin-content').style.display = 'none';
            return;
        }

        document.getElementById('admin-forbidden').style.display = 'none';
        document.getElementById('admin-content').style.display = 'block';

        // Load all data
        loadHalls();
        loadMovies();
        loadSessions();
        loadBookings();

        // Setup forms
        setupHallForm();
        setupMovieForm();
        setupSessionForm();
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
