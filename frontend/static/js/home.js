(function () {
  var allMovies = [];

  function updateNav() {
    var span = document.getElementById('user-span');
    var navProfile = document.getElementById('nav-profile');
    var navAdmin = document.getElementById('nav-admin');
    if (!window.auth) return;
    if (navProfile) navProfile.style.display = auth.isCustomer() ? '' : 'none';
    if (navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
    if (!span) return;

    if (auth.isLoggedIn()) {
      var u = auth.getUser();
      span.innerHTML =
        'Hi, ' +
        (u.first_name || u.email) +
        ' | <button type="button" id="logout-btn" class="btn btn-ghost">Log out</button>';
      var btn = document.getElementById('logout-btn');
      if (btn) btn.onclick = function () {
        auth.logout();
        updateNav();
      };
    } else {
      span.innerHTML =
        '<button type="button" id="login-btn" class="btn btn-primary">Login / Register</button>';
      var loginBtn = document.getElementById('login-btn');
      if (loginBtn) loginBtn.onclick = function () {
        if (window.openAuthModal) openAuthModal();
      };
    }
  }

  function ensureGenres(movies) {
    var select = document.getElementById('filter-genre');
    if (!select) return;
    var genres = {};
    movies.forEach(function (m) {
      if (!m.genres) return;
      var list = Array.isArray(m.genres) ? m.genres : String(m.genres).split(',');
      list.forEach(function (g) {
        var trimmed = String(g || '').trim();
        if (trimmed) genres[trimmed] = true;
      });
    });
    var current = select.value;
    select.innerHTML = '<option value="">All</option>';
    Object.keys(genres)
      .sort()
      .forEach(function (g) {
        var opt = document.createElement('option');
        opt.value = g;
        opt.textContent = g;
        select.appendChild(opt);
      });
    if (current && genres[current]) {
      select.value = current;
    }
  }

  function renderMovies(movies) {
    var grid = document.getElementById('movies-list');
    if (!grid) return;
    grid.innerHTML = '';
    var searchInput = document.getElementById('search-title');
    var genreSelect = document.getElementById('filter-genre');
    var q = (searchInput && searchInput.value ? searchInput.value : '').toLowerCase().trim();
    var genre = (genreSelect && genreSelect.value) || '';

    var filtered = (movies || []).filter(function (m) {
      var title = (m.name || m.title || '').toLowerCase();
      var matchTitle = !q || title.indexOf(q) !== -1;
      var genres = Array.isArray(m.genres) ? m.genres : (m.genres ? String(m.genres).split(',') : []);
      var matchGenre =
        !genre ||
        genres.some(function (g) {
          return String(g || '').trim() === genre;
        });
      return matchTitle && matchGenre;
    });

    if (filtered.length === 0) {
      var empty = document.createElement('p');
      empty.textContent = 'No movies found.';
      grid.appendChild(empty);
      return;
    }

    filtered.forEach(function (m) {
      var cardLink = document.createElement('a');
      cardLink.href = 'movie.html?id=' + encodeURIComponent(m.id || '');
      cardLink.className = 'card';

      var media = document.createElement('div');
      media.className = 'card-media';
      var img = document.createElement('img');
      img.src = m.poster_url || '';
      img.alt = m.name || 'Movie poster';
      media.appendChild(img);

      var body = document.createElement('div');
      body.className = 'card-body';

      var titleEl = document.createElement('div');
      titleEl.className = 'card-title';
      titleEl.textContent = m.name || 'Untitled';

      var meta = document.createElement('div');
      meta.className = 'card-meta';
      var duration = m.duration || m.duration_minutes;
      meta.textContent = duration ? duration + ' min' : '';

      var badgesRow = document.createElement('div');
      badgesRow.className = 'badge-row';
      var genres = Array.isArray(m.genres) ? m.genres : (m.genres ? String(m.genres).split(',') : []);
      genres.forEach(function (g) {
        var badge = document.createElement('span');
        badge.className = 'badge badge-genre';
        badge.textContent = String(g || '').trim();
        badgesRow.appendChild(badge);
      });
      if (m.age_restriction) {
        var ageBadge = document.createElement('span');
        ageBadge.className = 'badge badge-age';
        ageBadge.textContent = m.age_restriction;
        badgesRow.appendChild(ageBadge);
      }

      body.appendChild(titleEl);
      if (duration) body.appendChild(meta);
      if (badgesRow.children.length) body.appendChild(badgesRow);

      cardLink.appendChild(media);
      cardLink.appendChild(body);
      grid.appendChild(cardLink);
    });
  }

  function load() {
    var errorEl = document.getElementById('movies-error');
    if (errorEl) {
      errorEl.style.display = 'none';
      errorEl.textContent = '';
    }
    window.api
      .fetchMovies()
      .then(function (movies) {
        allMovies = movies;
        ensureGenres(allMovies);
        renderMovies(allMovies);
      })
      .catch(function (err) {
        if (errorEl) {
          errorEl.textContent = err.message || 'Failed to load movies.';
          errorEl.style.display = 'block';
        }
      });
  }

  function init() {
    updateNav();
    var search = document.getElementById('search-title');
    var genre = document.getElementById('filter-genre');
    if (search) search.addEventListener('input', function () { renderMovies(allMovies); });
    if (genre) genre.addEventListener('change', function () { renderMovies(allMovies); });
    load();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();

