(function () {
    var movieId = null;

    function getParam(name) {
        var m = new RegExp('([?&])' + name + '=([^&]*)').exec(location.search);
        return m ? decodeURIComponent(m[2]) : null;
    }

    function updateNav() {
        var navProfile = document.getElementById('movie-nav-profile');
        var navAdmin = document.getElementById('movie-nav-admin');
        var span = document.getElementById('movie-user-span');
        
        if (!window.auth) return;
        
        if (navProfile) navProfile.style.display = auth.isCustomer() ? '' : 'none';
        if (navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
        
        if (!span) return;

        if (auth.isLoggedIn()) {
            var u = auth.getUser();
            span.innerHTML =
                'Hi, ' + (u.first_name || u.email) +
                ' | <button type="button" id="movie-logout-btn" class="btn btn-ghost btn-sm">Log out</button>';
            var btn = document.getElementById('movie-logout-btn');
            if (btn) {
                btn.onclick = function () {
                    auth.logout();
                    updateNav();
                };
            }
        } else {
            span.innerHTML =
                '<button type="button" id="movie-login-btn" class="btn btn-primary btn-sm">Login / Register</button>';
            var loginBtn = document.getElementById('movie-login-btn');
            if (loginBtn) {
                loginBtn.onclick = function () {
                    if (window.openAuthModal) {
                        openAuthModal({ onSuccess: updateNav });
                    }
                };
            }
        }
    }

    function handleBuyTickets() {
        if (!movieId) return;
        
        if (window.auth && auth.isLoggedIn()) {
            // User is logged in, go directly to cinema selection
            location.href = 'cinemas.html?movieId=' + encodeURIComponent(movieId);
        } else if (window.openAuthModal) {
            // User not logged in, show auth modal
            openAuthModal({
                onSuccess: function () {
                    updateNav();
                    // After successful login, redirect to cinema selection
                    location.href = 'cinemas.html?movieId=' + encodeURIComponent(movieId);
                }
            });
        } else {
            // Fallback: go to cinema selection anyway
            location.href = 'cinemas.html?movieId=' + encodeURIComponent(movieId);
        }
    }

    function checkIfUserCanRate() {
        // Check if user attended a past session of this movie
        // This requires backend support - for now, we'll show the form if logged in
        // In real implementation, backend should return whether user can rate
        if (!window.auth || !auth.isLoggedIn()) {
            return;
        }

        // For now, we'll hide the rating form
        // In production, you'd check if user has attended sessions
        // var ratingWrap = document.getElementById('rating-form-wrap');
        // if (ratingWrap) ratingWrap.style.display = 'none';
    }

    function handleRatingSubmit(e) {
        e.preventDefault();
        
        var ratingInput = document.getElementById('rating-input');
        var commentInput = document.getElementById('comment-input');
        var errorEl = document.getElementById('rating-error');
        
        if (!ratingInput || !ratingInput.value) {
            if (errorEl) {
                errorEl.textContent = 'Please select a rating.';
                errorEl.style.display = 'block';
            }
            return;
        }

        var payload = {
            movie_id: movieId,
            rating: parseInt(ratingInput.value, 10),
            comment: commentInput ? commentInput.value.trim() : ''
        };

        window.api.createReview(payload)
            .then(function () {
                // Clear form
                ratingInput.value = '';
                if (commentInput) commentInput.value = '';
                if (errorEl) errorEl.style.display = 'none';
                
                // Hide rating form
                var ratingWrap = document.getElementById('rating-form-wrap');
                if (ratingWrap) ratingWrap.style.display = 'none';
                
                // Reload reviews
                loadReviews();
                
                alert('Thank you for your rating!');
            })
            .catch(function (err) {
                if (errorEl) {
                    errorEl.textContent = err.message || 'Failed to submit rating.';
                    errorEl.style.display = 'block';
                }
            });
    }

    function loadReviews() {
        window.api.fetchMovieReviews(movieId)
            .then(function (reviews) {
                var ul = document.getElementById('reviews-list');
                var noReviews = document.getElementById('no-reviews');
                
                if (!ul) return;
                
                ul.innerHTML = '';
                
                if (!Array.isArray(reviews) || reviews.length === 0) {
                    if (noReviews) noReviews.style.display = 'block';
                    return;
                }
                
                if (noReviews) noReviews.style.display = 'none';
                
                reviews.forEach(function (r) {
                    var li = document.createElement('li');
                    li.className = 'review-item';
                    
                    var content = '';
                    if (r.rating != null) {
                        content += '<strong>' + r.rating + '/5</strong>';
                    }
                    if (r.comment) {
                        content += '<p style="margin-top: 8px;">' + r.comment + '</p>';
                    }
                    
                    li.innerHTML = content;
                    ul.appendChild(li);
                });
            })
            .catch(function () {
                // Reviews are optional, don't show error
                var noReviews = document.getElementById('no-reviews');
                if (noReviews) noReviews.style.display = 'block';
            });
    }

    function loadMovie() {
        if (!movieId) {
            document.getElementById('movie-main').innerHTML = 
                '<div class="page-shell"><p>Please specify a movie ID in the URL.</p></div>';
            return;
        }

        var errorEl = document.getElementById('movie-error');
        
        window.api.fetchMovieDetails(movieId)
            .then(function (movie) {
                // Poster
                var poster = document.getElementById('movie-poster');
                if (poster) {
                    poster.src = movie.poster_url || '';
                    poster.alt = (movie.name || 'Movie') + ' poster';
                }

                // Title
                var title = document.getElementById('movie-title');
                if (title) title.textContent = movie.name || 'Untitled';

                // Badges (Age restriction)
                var badgesEl = document.getElementById('movie-badges');
                if (badgesEl && movie.age_restriction) {
                    var ageBadge = document.createElement('span');
                    ageBadge.className = 'badge badge-age';
                    ageBadge.textContent = movie.age_restriction;
                    badgesEl.appendChild(ageBadge);
                }

                // Year
                var yearEl = document.getElementById('movie-meta-year');
                if (yearEl && movie.year) {
                    yearEl.textContent = 'Year: ' + movie.year;
                }

                // Duration
                var duration = document.getElementById('movie-duration');
                if (duration) {
                    duration.textContent = 'Duration: ' + (movie.duration || 0) + ' min';
                }

                // Genres
                var genres = Array.isArray(movie.genres) 
                    ? movie.genres 
                    : (movie.genres ? String(movie.genres).split(',').map(function(g) { return g.trim(); }) : []);
                var genresEl = document.getElementById('movie-meta-genres');
                if (genresEl && genres.length) {
                    genresEl.textContent = 'Genres: ' + genres.join(', ');
                }

                // Rating
                var ratingEl = document.getElementById('movie-rating');
                if (ratingEl) {
                    ratingEl.textContent = 'Rating: ' + (movie.rating != null ? movie.rating + '/10' : 'Not rated yet');
                }

                // Description
                var desc = document.getElementById('movie-description');
                if (desc) {
                    desc.textContent = movie.description || 'No description available.';
                }

                // Trailer
                var trailer = document.getElementById('movie-trailer');
                var trailerWrap = document.getElementById('movie-trailer-wrap');
                if (trailer && trailerWrap) {
                    if (movie.trailer_url) {
                        trailer.src = movie.trailer_url;
                        trailerWrap.style.display = 'block';
                    } else {
                        trailerWrap.style.display = 'none';
                    }
                }

                // Load reviews
                loadReviews();
                
                // Check if user can rate
                checkIfUserCanRate();
            })
            .catch(function (err) {
                if (errorEl) {
                    errorEl.textContent = err.message || 'Failed to load movie details.';
                    errorEl.style.display = 'block';
                }
            });
    }

    function init() {
        updateNav();
        movieId = getParam('id');
        
        // Buy tickets button
        var buyBtn = document.getElementById('buy-tickets-btn');
        if (buyBtn) {
            buyBtn.addEventListener('click', handleBuyTickets);
        }

        // Rating form
        var ratingForm = document.getElementById('rating-form');
        if (ratingForm) {
            ratingForm.addEventListener('submit', handleRatingSubmit);
        }

        loadMovie();
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
