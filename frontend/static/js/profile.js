(function () {
    function updateNav() {
        var span = document.getElementById('user-span');
        var navAdmin = document.getElementById('nav-admin');
        
        if (!window.auth) return;
        
        if (navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
        
        if (!span) return;

        if (auth.isLoggedIn()) {
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
            span.innerHTML =
                '<button type="button" id="login-btn" class="btn btn-primary btn-sm">Login / Register</button>';
            var loginBtn = document.getElementById('login-btn');
            if (loginBtn) {
                loginBtn.onclick = function () {
                    if (window.openAuthModal) {
                        openAuthModal({ onSuccess: load });
                    }
                };
            }
        }
    }

    function showGuest() {
        document.getElementById('profile-guest').style.display = 'block';
        document.getElementById('profile-content').style.display = 'none';
        
        var loginBtn = document.getElementById('profile-login-btn');
        if (loginBtn) {
            loginBtn.onclick = function () {
                if (window.openAuthModal) {
                    openAuthModal({
                        onSuccess: function () {
                            load();
                            updateNav();
                        }
                    });
                }
            };
        }
    }

    function showError(msg) {
        var el = document.getElementById('profile-error');
        if (!el) return;
        el.textContent = msg || '';
        el.style.display = msg ? 'block' : 'none';
    }

    function load() {
        if (!window.auth || !auth.isLoggedIn()) {
            showGuest();
            updateNav();
            return;
        }

        document.getElementById('profile-guest').style.display = 'none';
        document.getElementById('profile-content').style.display = 'block';
        showError('');

        var u = auth.getUser();
        
        // Display user info
        var infoEl = document.getElementById('profile-info');
        if (infoEl) {
            infoEl.innerHTML =
                '<strong>Name:</strong> ' + (u.first_name || '') + ' ' + (u.last_name || '') + '<br>' +
                '<strong>Email:</strong> ' + (u.email || '') + '<br>' +
                '<strong>Role:</strong> ' + (u.role || 'USER');
        }

        // Show admin link if user is admin
        var adminLink = document.getElementById('profile-admin-link');
        if (adminLink) {
            adminLink.style.display = auth.isAdmin() ? 'inline-block' : 'none';
        }

        // Load bookings
        window.api.fetchMyBookings()
            .then(function (data) {
                var list = document.getElementById('bookings-list');
                var empty = document.getElementById('bookings-empty');
                
                if (!list) return;
                
                list.innerHTML = '';

                if (!Array.isArray(data) || data.length === 0) {
                    if (empty) empty.style.display = 'block';
                    return;
                }

                if (empty) empty.style.display = 'none';

                data.forEach(function (booking) {
                    var li = document.createElement('li');
                    li.dataset.id = booking.id;
                    
                    var info = document.createElement('div');
                    info.innerHTML =
                        '<strong>Session ID:</strong> ' + (booking.session_id || '') + '<br>' +
                        '<strong>Seat:</strong> Row ' + (booking.row_number || '') + ', Seat ' + (booking.seat_number || '') + '<br>' +
                        '<strong>Status:</strong> ' + (booking.status || 'confirmed');
                    
                    var actions = document.createElement('div');
                    var cancelBtn = document.createElement('button');
                    cancelBtn.type = 'button';
                    cancelBtn.className = 'btn btn-danger btn-sm cancel-ticket-btn';
                    cancelBtn.dataset.id = booking.id || '';
                    cancelBtn.textContent = 'Cancel';
                    actions.appendChild(cancelBtn);
                    
                    li.appendChild(info);
                    li.appendChild(actions);
                    list.appendChild(li);
                });

                // Add cancel handlers
                list.querySelectorAll('.cancel-ticket-btn').forEach(function (btn) {
                    btn.onclick = function () {
                        var id = btn.getAttribute('data-id');
                        if (!id) return;
                        
                        if (!confirm('Are you sure you want to cancel this booking?')) return;

                        window.api.cancelBooking(id)
                            .then(function () {
                                load(); // Reload bookings
                            })
                            .catch(function (err) {
                                showError((err && err.message) || 'Failed to cancel booking.');
                            });
                    };
                });
            })
            .catch(function (err) {
                showError((err && err.message) || 'Failed to load bookings.');
            });
    }

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
