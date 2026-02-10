(function () {
  function getParam(name) {
    var m = new RegExp('([?&])' + name + '=([^&]*)').exec(location.search);
    return m ? decodeURIComponent(m[2]) : null;
  }

  function updateNav() {
    var span = document.getElementById('cinema-user-span');
    var navProfile = document.getElementById('cinema-nav-profile');
    var navAdmin = document.getElementById('cinema-nav-admin');
    if (!window.auth) return;
    if (navProfile) navProfile.style.display = auth.isCustomer() ? '' : 'none';
    if (navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
    if (!span) return;
    if (auth.isLoggedIn()) {
      var u = auth.getUser();
      span.innerHTML =
        'Hi, ' +
        (u.first_name || u.email) +
        ' | <button type="button" id="cinema-logout-btn" class="btn btn-ghost">Log out</button>';
      var btn = document.getElementById('cinema-logout-btn');
      if (btn) btn.onclick = function () { auth.logout(); updateNav(); };
    } else {
      span.innerHTML =
        '<button type="button" id="cinema-login-btn" class="btn btn-primary">Login / Register</button>';
      var loginBtn = document.getElementById('cinema-login-btn');
      if (loginBtn) loginBtn.onclick = function () { if (window.openAuthModal) openAuthModal(); };
    }
  }

  function deriveCinemaType(hallName) {
    var name = (hallName || '').toLowerCase();
    if (name.indexOf('imax') !== -1) return 'IMAX';
    if (name.indexOf('vip') !== -1) return 'VIP';
    return 'Standard';
  }

  function buildCinemasForMovie(movieId) {
    // Derive pseudo-cinemas from halls that have upcoming sessions for this movie.
    return window.api.fetchMovieSessions(movieId).then(function (sessions) {
      var byHall = {};
      sessions.forEach(function (s) {
        if (!s.hall_id) return;
        if (!byHall[s.hall_id]) {
          byHall[s.hall_id] = { hallId: s.hall_id, sessions: [] };
        }
        byHall[s.hall_id].sessions.push(s);
      });
      var hallIds = Object.keys(byHall);
      if (hallIds.length === 0) return [];

      return Promise.all(
        hallIds.map(function (id) {
          return window.api.fetchHall(id).then(function (hall) {
            byHall[id].hall = hall;
            return byHall[id];
          });
        })
      ).then(function (list) {
        return list.map(function (entry) {
          var h = entry.hall || {};
          return {
            id: h.id,
            name: h.name || 'Mangekyo Cinema',
            photoUrl: '',
            type: deriveCinemaType(h.name),
            address: h.location || 'Main building',
            hallsSummary: 'Rows: ' + (h.total_rows || 0) + ', seats/row: ' + (h.seats_per_row || 0),
            upcomingSessionsCount: entry.sessions.length
          };
        });
      });
    });
  }

  function renderMovieSummary(movie) {
    var wrap = document.getElementById('cinema-movie-summary');
    if (!wrap || !movie) return;
    wrap.innerHTML = '';
    var imgWrap = document.createElement('div');
    imgWrap.style.flex = '0 0 80px';
    imgWrap.style.maxWidth = '80px';
    var img = document.createElement('img');
    img.src = movie.poster_url || '';
    img.alt = movie.name || 'Movie poster';
    img.style.borderRadius = '8px';
    imgWrap.appendChild(img);
    var info = document.createElement('div');
    var title = document.createElement('div');
    title.className = 'card-title';
    title.textContent = movie.name || 'Untitled';
    var meta = document.createElement('div');
    meta.className = 'card-meta';
    meta.textContent = (movie.duration || 0) + ' min';
    info.appendChild(title);
    info.appendChild(meta);
    wrap.appendChild(imgWrap);
    wrap.appendChild(info);
  }

  function renderCinemas(cinemas, movieId) {
    var list = document.getElementById('cinema-list');
    var err = document.getElementById('cinema-error');
    if (err) {
      err.style.display = 'none';
      err.textContent = '';
    }
    if (!list) return;
    list.innerHTML = '';
    if (!cinemas.length) {
      var p = document.createElement('p');
      p.textContent = 'No cinemas with upcoming sessions for this movie.';
      list.appendChild(p);
      return;
    }
    cinemas.forEach(function (c) {
      var card = document.createElement('button');
      card.type = 'button';
      card.className = 'card';
      card.style.textAlign = 'left';

      var body = document.createElement('div');
      body.className = 'card-body';

      var title = document.createElement('div');
      title.className = 'card-title';
      title.textContent = c.name;

      var meta = document.createElement('div');
      meta.className = 'card-meta';
      meta.textContent = c.address;

      var badges = document.createElement('div');
      badges.className = 'badge-row';
      var typeBadge = document.createElement('span');
      typeBadge.className = 'badge badge-format';
      typeBadge.textContent = c.type;
      badges.appendChild(typeBadge);
      var sessionsBadge = document.createElement('span');
      sessionsBadge.className = 'badge';
      sessionsBadge.textContent = c.upcomingSessionsCount + ' upcoming sessions';
      badges.appendChild(sessionsBadge);

      body.appendChild(title);
      body.appendChild(meta);
      body.appendChild(badges);
      card.appendChild(body);

      card.addEventListener('click', function () {
        location.href =
          'sessions.html?movieId=' +
          encodeURIComponent(movieId) +
          '&cinemaId=' +
          encodeURIComponent(c.id || '');
      });

      list.appendChild(card);
    });
  }

  function init() {
    updateNav();
    var movieId = getParam('movieId');
    if (!movieId) {
      var err = document.getElementById('cinema-error');
      if (err) {
        err.textContent = 'Specify movie in URL: cinemas.html?movieId=...';
        err.style.display = 'block';
      }
      return;
    }
    window.api
      .fetchMovieDetails(movieId)
      .then(function (movie) {
        renderMovieSummary(movie);
        return buildCinemasForMovie(movieId);
      })
      .then(function (cinemas) {
        renderCinemas(cinemas, movieId);
      })
      .catch(function (err) {
        var el = document.getElementById('cinema-error');
        if (el) {
          el.textContent = (err && err.message) || 'Failed to load cinemas.';
          el.style.display = 'block';
        }
      });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();

