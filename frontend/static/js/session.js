(function () {
    var API_BASE = window.API_BASE || '/api';
    var sessionId = null;
    var session = null;
    var movie = null;
    var hall = null;
    var bookedSeats = [];
    var selectedSeat = null;
    var seatMapEl = null;
    var confirmBtn = null;
    var successEl = null;
    var confirmErrorEl = null;

    function getParam(name) {
        var m = new RegExp('([?&])' + name + '=([^&]*)').exec(location.search);
        return m ? decodeURIComponent(m[2]) : null;
    }

    function authHeaders() {
        return window.auth && window.auth.getAuthHeaders ? window.auth.getAuthHeaders() : {};
    }

    function updateNav() {
        var navProfile = document.getElementById('session-nav-profile');
        var navAdmin = document.getElementById('session-nav-admin');
        var span = document.getElementById('session-user-span');
        
        if (!window.auth) return;
        
        if (navProfile) navProfile.style.display = auth.isCustomer() ? '' : 'none';
        if (navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
        
        if (!span) return;

        if (auth.isLoggedIn()) {
            var u = auth.getUser();
            span.innerHTML =
                'Hi, ' + (u.first_name || u.email) +
                ' | <button type="button" id="session-logout-btn" class="btn btn-ghost btn-sm">Log out</button>';
            var btn = document.getElementById('session-logout-btn');
            if (btn) {
                btn.onclick = function () {
                    auth.logout();
                    updateNav();
                };
            }
        } else {
            span.innerHTML =
                '<button type="button" id="session-login-btn" class="btn btn-primary btn-sm">Login / Register</button>';
            var loginBtn = document.getElementById('session-login-btn');
            if (loginBtn) {
                loginBtn.onclick = function () {
                    if (window.openAuthModal) {
                        openAuthModal({ onSuccess: updateNav });
                    }
                };
            }
        }
    }

    function ensureLoggedIn(callback) {
        if (window.auth && auth.isLoggedIn()) {
            if (callback) callback();
            return true;
        }
        
        if (window.openAuthModal) {
            openAuthModal({
                onSuccess: function () {
                    updateNav();
                    if (callback) callback();
                }
            });
        }
        return false;
    }

    function setLoading(loading) {
        if (confirmBtn) confirmBtn.disabled = loading;
    }

    function showSuccess() {
        if (successEl) successEl.style.display = 'block';
        if (confirmBtn) confirmBtn.style.display = 'none';
        if (seatMapEl) seatMapEl.style.pointerEvents = 'none';
    }

    function showConfirmError(msg) {
        if (!confirmErrorEl) return;
        confirmErrorEl.textContent = msg || '';
        confirmErrorEl.style.display = msg ? 'block' : 'none';
    }

    function fetchSession() {
        if (!sessionId) return Promise.reject(new Error('No sessionId'));
        return fetch(API_BASE + '/sessions/' + sessionId)
            .then(function (r) { return r.json(); });
    }

    function fetchMovie(id) {
        return fetch(API_BASE + '/movies/' + id)
            .then(function (r) { return r.json(); });
    }

    function fetchHall(id) {
        return fetch(API_BASE + '/halls/' + id)
            .then(function (r) { return r.json(); });
    }

    function fetchBookedSeats() {
        return fetch(API_BASE + '/sessions/' + sessionId + '/booked-seats')
            .then(function (r) { return r.json(); });
    }

    function isBooked(row, seat) {
        return bookedSeats.some(function (s) {
            return s.row_number === row && s.seat_number === seat;
        });
    }

    function renderSeatMap() {
        if (!hall || !seatMapEl) return;
        
        var rows = hall.total_rows;
        var seatsPerRow = hall.seats_per_row;
        seatMapEl.innerHTML = '';

        for (var r = 1; r <= rows; r++) {
            var rowDiv = document.createElement('div');
            rowDiv.className = 'seat-row';

            // Left row number
            var leftNum = document.createElement('span');
            leftNum.className = 'row-num';
            leftNum.textContent = r;
            rowDiv.appendChild(leftNum);

            // Seats
            var cells = document.createElement('div');
            cells.className = 'seat-row-cells';
            
            for (var s = 1; s <= seatsPerRow; s++) {
                var cell = document.createElement('button');
                cell.type = 'button';
                cell.className = 'seat';
                cell.dataset.row = r;
                cell.dataset.seat = s;
                cell.setAttribute('aria-label', 'Row ' + r + ', Seat ' + s);
                
                var occupied = isBooked(r, s);
                if (occupied) {
                    cell.classList.add('occupied');
                    cell.disabled = true;
                } else {
                    cell.classList.add('free');
                    cell.addEventListener('click', function (ev) {
                        var row = parseInt(ev.currentTarget.dataset.row, 10);
                        var seat = parseInt(ev.currentTarget.dataset.seat, 10);
                        selectSeat(row, seat);
                    });
                }
                
                cells.appendChild(cell);
            }
            
            rowDiv.appendChild(cells);

            // Right row number
            var rightNum = document.createElement('span');
            rightNum.className = 'row-num';
            rightNum.textContent = r;
            rowDiv.appendChild(rightNum);

            seatMapEl.appendChild(rowDiv);
        }
        
        updateSeatSelection();
    }

    function selectSeat(row, seat) {
        if (isBooked(row, seat)) return;
        
        selectedSeat = { row: row, seat: seat };
        updateSeatSelection();
        confirmBtn.disabled = false;
        
        // Update seat info display
        var seatInfoText = document.getElementById('seat-info-text');
        if (seatInfoText) {
            seatInfoText.textContent = 'Selected: Row ' + row + ', Seat ' + seat;
            seatInfoText.style.color = '#46d369';
        }
    }

    function updateSeatSelection() {
        var seats = seatMapEl ? seatMapEl.querySelectorAll('.seat:not(.occupied)') : [];
        
        for (var i = 0; i < seats.length; i++) {
            var el = seats[i];
            var r = parseInt(el.dataset.row, 10);
            var s = parseInt(el.dataset.seat, 10);
            
            if (selectedSeat && selectedSeat.row === r && selectedSeat.seat === s) {
                el.classList.add('selected');
                el.classList.remove('free');
            } else {
                el.classList.remove('selected');
                el.classList.add('free');
            }
        }
    }

    function fillHeader() {
        var poster = document.getElementById('session-poster');
        var title = document.getElementById('session-movie-title');
        var duration = document.getElementById('session-duration');
        var hallEl = document.getElementById('session-hall');
        var timeEl = document.getElementById('session-time');
        var priceEl = document.getElementById('session-price');

        if (poster && movie) {
            poster.src = movie.poster_url || '';
            poster.alt = (movie.name || 'Movie') + ' poster';
        }

        if (title && movie) {
            title.textContent = movie.name || 'Untitled';
        }

        if (duration && movie) {
            duration.textContent = 'Duration: ' + (movie.duration || 0) + ' min';
        }

        if (hallEl && hall) {
            hallEl.textContent = 'Hall: ' + (hall.name || '') + 
                (hall.location ? ' (' + hall.location + ')' : '');
        }

        if (timeEl && session && session.start_time) {
            var start = new Date(session.start_time);
            timeEl.textContent = 'Session: ' + start.toLocaleString();
        }

        if (priceEl && session) {
            priceEl.textContent = 'Price: $' + (session.price != null ? session.price : 0);
        }
    }

    function doBooking() {
        if (!selectedSeat || !sessionId) return;

        ensureLoggedIn(function () {
            setLoading(true);
            showConfirmError('');

            fetch(API_BASE + '/bookings', {
                method: 'POST',
                headers: authHeaders(),
                body: JSON.stringify({
                    session_id: sessionId,
                    row_number: selectedSeat.row,
                    seat_number: selectedSeat.seat
                })
            })
                .then(function (res) {
                    return res.json().then(function (data) {
                        return { ok: res.ok, data: data };
                    });
                })
                .then(function (result) {
                    setLoading(false);
                    if (result.ok) {
                        showSuccess();
                        bookedSeats.push({
                            row_number: selectedSeat.row,
                            seat_number: selectedSeat.seat
                        });
                        selectedSeat = null;
                        renderSeatMap();
                    } else {
                        showConfirmError(result.data.error || 'Booking failed. Please try again.');
                    }
                })
                .catch(function () {
                    setLoading(false);
                    showConfirmError('Network error. Please check your connection.');
                });
        });
    }

    function init() {
        updateNav();
        
        sessionId = getParam('sessionId') || getParam('session_id');
        seatMapEl = document.getElementById('seat-map');
        confirmBtn = document.getElementById('confirm-booking-btn');
        successEl = document.getElementById('success-message');
        confirmErrorEl = document.getElementById('confirm-error');

        if (!sessionId) {
            document.getElementById('session-page').innerHTML =
                '<div class="page-shell"><p>Please specify a session ID in the URL.</p></div>';
            return;
        }

        // Load session data
        Promise.all([fetchSession(), fetchBookedSeats()])
            .then(function (arr) {
                session = arr[0];
                if (session.error) throw new Error(session.error || 'Session not found');
                
                bookedSeats = Array.isArray(arr[1]) ? arr[1] : [];
                
                return fetchMovie(session.movie_id).then(function (m) {
                    movie = m;
                    return fetchHall(session.hall_id);
                });
            })
            .then(function (h) {
                hall = h;
                if (!hall || !hall.total_rows) {
                    throw new Error('Hall not found');
                }
                
                fillHeader();
                renderSeatMap();
            })
            .catch(function (err) {
                document.getElementById('session-page').innerHTML =
                    '<div class="page-shell"><p class="error-msg">Error: ' + 
                    (err.message || 'Failed to load session') + '</p></div>';
            });

        // Confirm booking button
        if (confirmBtn) {
            confirmBtn.addEventListener('click', function () {
                if (!selectedSeat) return;
                doBooking();
            });
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
