(function () {
  function getParam(name) {
    var m = new RegExp('([?&])' + name + '=([^&]*)').exec(location.search);
    return m ? decodeURIComponent(m[2]) : null;
  }

  function updateNav() {
    var span = document.getElementById('sessions-user-span');
    var navProfile = document.getElementById('sessions-nav-profile');
    var navAdmin = document.getElementById('sessions-nav-admin');
    if (!window.auth) return;
    if (navProfile) navProfile.style.display = auth.isCustomer() ? '' : 'none';
    if (navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
    if (!span) return;
    if (auth.isLoggedIn()) {
      var u = auth.getUser();
      span.innerHTML =
        (u.first_name || u.email) +
        ' | <a href="index.html">Home</a> | <a href="profile.html">Profile</a> | <button type="button" id="sessions-logout-btn" class="btn btn-ghost">Log out</button>';
      var btn = document.getElementById('sessions-logout-btn');
      if (btn) btn.onclick = function () { auth.logout(); updateNav(); };
    } else {
      span.innerHTML =
        '<button type="button" id="sessions-login-btn" class="btn btn-primary">Login / Register</button>';
      var loginBtn = document.getElementById('sessions-login-btn');
      if (loginBtn) loginBtn.onclick = function () { if (window.openAuthModal) openAuthModal({ onSuccess: updateNav }); };
    }
  }

  function buildDateRange(days) {
    var out = [];
    var now = new Date();
    for (var i = 0; i < days; i++) {
      var d = new Date(now);
      d.setDate(now.getDate() + i);
      out.push(d);
    }
    return out;
  }

  function sameDate(a, b) {
    return (
      a.getFullYear() === b.getFullYear() &&
      a.getMonth() === b.getMonth() &&
      a.getDate() === b.getDate()
    );
  }

  function renderSummary(movie, cinema) {
    var wrap = document.getElementById('sessions-summary');
    if (!wrap) return;
    wrap.innerHTML = '';
    if (!movie) return;
    var imgWrap = document.createElement('div');
    imgWrap.style.flex = '0 0 80px';
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
    meta.textContent = cinema && cinema.name ? cinema.name : '';
    info.appendChild(title);
    info.appendChild(meta);
    wrap.appendChild(imgWrap);
    wrap.appendChild(info);
  }

  function renderDateSelector(dates, activeDate, onSelect) {
    var bar = document.getElementById('date-selector');
    if (!bar) return;
    bar.innerHTML = '';
    dates.forEach(function (d) {
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'btn btn-ghost';
      var label =
        d.toLocaleDateString(undefined, { weekday: 'short' }) +
        ' ' +
        d.getDate() +
        '.' +
        (d.getMonth() + 1);
      btn.textContent = label;
      if (sameDate(d, activeDate)) {
        btn.classList.remove('btn-ghost');
        btn.classList.add('btn-primary');
      }
      btn.addEventListener('click', function () {
        onSelect(d);
      });
      bar.appendChild(btn);
    });
  }

  function filterSessionsByDate(sessions, date) {
    var now = new Date();
    return sessions.filter(function (s) {
      if (!s.start_time) return false;
      var start = new Date(s.start_time);
      if (start < now) {
        // Hide expired sessions (past now).
        return false;
      }
      return sameDate(start, date);
    });
  }

  function renderSessions(sessions) {
    var list = document.getElementById('sessions-list');
    if (!list) return;
    list.innerHTML = '';
    if (!sessions.length) {
      var p = document.createElement('p');
      p.textContent = 'No sessions available for this day.';
      list.appendChild(p);
      return;
    }
    sessions.forEach(function (s) {
      var card = document.createElement('button');
      card.type = 'button';
      card.className = 'card';
      card.style.textAlign = 'left';

      var body = document.createElement('div');
      body.className = 'card-body';

      var title = document.createElement('div');
      title.className = 'card-title';
      var start = s.start_time ? new Date(s.start_time) : null;
      title.textContent = start ? start.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '';

      var meta = document.createElement('div');
      meta.className = 'card-meta';
      meta.textContent = 'Price: ' + (s.price != null ? s.price : 0);

      var badges = document.createElement('div');
      badges.className = 'badge-row';

      if (s.format) {
        var formatBadge = document.createElement('span');
        formatBadge.className = 'badge badge-format';
        formatBadge.textContent = s.format;
        badges.appendChild(formatBadge);
      }
      if (s.language) {
        var langBadge = document.createElement('span');
        langBadge.className = 'badge badge-language';
        langBadge.textContent = s.language;
        badges.appendChild(langBadge);
      }
      if (s.age_restriction) {
        var ageBadge = document.createElement('span');
        ageBadge.className = 'badge badge-age';
        ageBadge.textContent = s.age_restriction;
        badges.appendChild(ageBadge);
      }

      body.appendChild(title);
      body.appendChild(meta);
      body.appendChild(badges);
      card.appendChild(body);

      card.addEventListener('click', function () {
        location.href = 'session.html?sessionId=' + encodeURIComponent(s.id || '');
      });

      list.appendChild(card);
    });
  }

  function init() {
    updateNav();
    var movieId = getParam('movieId');
    var cinemaId = getParam('cinemaId');
    if (!movieId) {
      var err = document.getElementById('sessions-error');
      if (err) {
        err.textContent = 'Specify movieId in URL.';
        err.style.display = 'block';
      }
      return;
    }

    var allSessions = [];
    var dates = buildDateRange(7);
    var activeDate = dates[0];

    function refresh() {
      renderDateSelector(dates, activeDate, function (d) {
        activeDate = d;
        refresh();
      });
      var filtered = allSessions;
      if (cinemaId) {
        filtered = filtered.filter(function (s) {
          return s.hall_id === cinemaId;
        });
      }
      filtered = filterSessionsByDate(filtered, activeDate);
      renderSessions(filtered);
    }

    Promise.all([
      window.api.fetchMovieDetails(movieId),
      window.api.fetchMovieSessions(movieId)
    ])
      .then(function (results) {
        var movie = results[0];
        allSessions = Array.isArray(results[1]) ? results[1] : [];
        renderSummary(movie, null);
        refresh();
      })
      .catch(function (err) {
        var el = document.getElementById('sessions-error');
        if (el) {
          el.textContent = (err && err.message) || 'Failed to load sessions.';
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

