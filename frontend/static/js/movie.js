(function () {
    function getParam(name) {
        var m = new RegExp('([?&])' + name + '=([^&]*)').exec(location.search);
        return m ? decodeURIComponent(m[2]) : null;
    }

    function updateNav() {
        var navProfile = document.getElementById('nav-profile');
        var navAdmin = document.getElementById('nav-admin');
        var span = document.getElementById('user-span');
        var reviewWrap = document.getElementById('add-review-wrap');

        var isLoggedIn = window.auth && auth.isLoggedIn && auth.isLoggedIn();

        if (window.auth && auth.isCustomer && navProfile) navProfile.style.display = auth.isCustomer() ? '' : 'none';
        if (window.auth && auth.isAdmin && navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
        if (reviewWrap) reviewWrap.style.display = isLoggedIn ? 'block' : 'none';

        if (!span) return;
        if (isLoggedIn) {
            var u = auth.getUser();
            span.innerHTML =
                'Hi, ' + (u.first_name || u.email) +
                ' | <button type="button" id="logout-btn" class="btn btn-ghost">Log out</button>';
            var btn = document.getElementById('logout-btn');
            if (btn) btn.onclick = function () { auth.logout(); updateNav(); };
        } else {
            span.innerHTML = '<button type="button" id="login-btn" class="btn btn-primary">Login / Register</button>';
            var loginBtn = document.getElementById('login-btn');
            if (loginBtn) loginBtn.onclick = function () { if (window.openAuthModal) openAuthModal({ onSuccess: function () { location.reload(); } }); };
        }
    }

    var movieId = getParam('id');
    if (!movieId) {
        var main = document.getElementById('movie-main');
        if (main) main.innerHTML = '<div class="page-shell"><p>Specify movie in URL: movie.html?id=...</p></div>';
        return;
    }

    function loadMovieDetails() {
        window.api
            .fetchMovieDetails(movieId)
            .then(function (movie) {
                document.getElementById('movie-poster').src = movie.poster_url || '';
                document.getElementById('movie-poster').alt = movie.name || 'Poster';
                document.getElementById('movie-title').textContent = movie.name || 'Untitled';
                document.getElementById('movie-duration').textContent = 'Duration: ' + (movie.duration || 0) + ' min';
                var genres = Array.isArray(movie.genres) ? movie.genres : (movie.genres ? String(movie.genres).split(',') : []);
                document.getElementById('movie-meta-genres').textContent =
                    genres.length ? 'Genres: ' + genres.join(', ') : '';

                var ageEl = document.getElementById('movie-age-limit');
                if (ageEl) {
                    ageEl.textContent = (movie.age_limit || 0) + '+';
                }

                document.getElementById('movie-description').textContent = movie.description || '';

                var ratEl = document.getElementById('movie-rating');
                if (ratEl) {
                    ratEl.textContent = 'Rating: ' + (movie.rating != null ? movie.rating.toFixed(1) : '—') + '/10';
                }
                var trailer = document.getElementById('movie-trailer');
                if (movie.trailer_url) {
                    trailer.src = movie.trailer_url.replace('watch?v=', 'embed/');
                    document.getElementById('movie-trailer-wrap').style.display = 'block';
                } else {
                    document.getElementById('movie-trailer-wrap').style.display = 'none';
                }
                return window.api.fetchMovieReviews(movieId);
            })
            .then(function (reviews) {
                var ul = document.getElementById('reviews-list');
                ul.innerHTML = '';
                if (!reviews || reviews.length === 0) {
                    ul.innerHTML = '<li>No reviews yet.</li>';
                    return;
                }
                reviews.forEach(function (r) {
                    var li = document.createElement('li');
                    li.className = 'review-item';
                    li.innerHTML = `
                        <div class="review-rating">${r.rating || 0}/10</div>
                        <div class="card-meta" style="margin-bottom: 0.5rem; color: var(--accent);">By: ${r.user_name || 'Anonymous'}</div>
                        <p class="review-comment">${r.comment || ''}</p>
                    `;
                    ul.appendChild(li);
                });
            })
            .catch(function (err) {
                var errEl = document.getElementById('movie-error');
                if (errEl) {
                    errEl.textContent = err.message || 'Failed to load.';
                    errEl.style.display = 'block';
                }
            });
    }

    function handleBuyTickets() {
        if (window.auth && auth.isLoggedIn && auth.isLoggedIn()) {
            location.href = 'cinemas.html?movieId=' + encodeURIComponent(movieId);
        } else if (window.openAuthModal) {
            openAuthModal({
                onSuccess: function () {
                    updateNav();
                    location.href = 'cinemas.html?movieId=' + encodeURIComponent(movieId);
                }
            });
        } else {
            location.href = 'cinemas.html?movieId=' + encodeURIComponent(movieId);
        }
    }

    function init() {
        updateNav();
        loadMovieDetails();
        var buyBtn = document.getElementById('buy-tickets-btn');
        if (buyBtn && movieId) {
            buyBtn.addEventListener('click', handleBuyTickets);
        }

        var formReview = document.getElementById('form-review');
        if (formReview) {
            formReview.onsubmit = function (e) {
                e.preventDefault();
                var payload = {
                    movie_id: movieId,
                    rating: parseInt(this.rating.value, 10),
                    comment: this.comment.value.trim()
                };
                window.api.createReview(payload)
                    .then(function () {
                        formReview.reset();
                        loadMovieDetails();
                    })
                    .catch(function (err) {
                        alert(err.message || 'Failed to submit review');
                    });
            };
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
