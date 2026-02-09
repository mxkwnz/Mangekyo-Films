(function () {
  var API_BASE = window.API_BASE || '/api';

  function withJson(response) {
    return response.json().then(function (data) {
      return { ok: response.ok, status: response.status, data: data };
    });
  }

  function authHeaders() {
    return window.auth && window.auth.getAuthHeaders ? window.auth.getAuthHeaders() : {};
  }

  function request(method, path, options) {
    options = options || {};
    var headers = options.headers || {};
    var body = options.body;
    var isJson = body && typeof body === 'object' && !(body instanceof FormData);
    if (isJson) {
      headers['Content-Type'] = 'application/json';
      body = JSON.stringify(body);
    }

    return fetch(API_BASE + path, {
      method: method,
      headers: headers,
      body: body
    }).then(withJson);
  }

  window.api = {
    /* Movies */
    fetchMovies: function () {
      return request('GET', '/movies').then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load movies');
        return Array.isArray(res.data) ? res.data : [];
      });
    },
    fetchMovieDetails: function (id) {
      return request('GET', '/movies/' + encodeURIComponent(id)).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load movie');
        return res.data;
      });
    },

    /* Sessions & halls */
    fetchMovieSessions: function (movieId) {
      return request('GET', '/sessions/movie/' + encodeURIComponent(movieId)).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load sessions');
        return Array.isArray(res.data) ? res.data : [];
      });
    },
    fetchUpcomingSessions: function () {
      return request('GET', '/sessions/upcoming').then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load sessions');
        return Array.isArray(res.data) ? res.data : [];
      });
    },
    fetchSession: function (id) {
      return request('GET', '/sessions/' + encodeURIComponent(id)).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load session');
        return res.data;
      });
    },
    fetchHall: function (id) {
      return request('GET', '/halls/' + encodeURIComponent(id)).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load hall');
        return res.data;
      });
    },
    fetchSessionBookedSeats: function (sessionId) {
      return request('GET', '/sessions/' + encodeURIComponent(sessionId) + '/booked-seats').then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load seats');
        return Array.isArray(res.data) ? res.data : [];
      });
    },

    /**
     * Pseudo-cinema list for a movie, derived from halls that
     * have upcoming sessions for the given movie.
     */
    fetchCinemasForMovie: function (movieId) {
      var self = this;
      return self.fetchMovieSessions(movieId).then(function (sessions) {
        var byHall = {};
        sessions.forEach(function (s) {
          if (!s.hall_id) return;
          if (!byHall[s.hall_id]) {
            byHall[s.hall_id] = { hallId: s.hall_id, sessions: [] };
          }
          byHall[s.hall_id].sessions.push(s);
        });
        var hallIds = Object.keys(byHall);
        if (!hallIds.length) return [];
        return Promise.all(
          hallIds.map(function (id) {
            return self.fetchHall(id).then(function (hall) {
              byHall[id].hall = hall;
              return byHall[id];
            });
          })
        ).then(function (entries) {
          return entries.map(function (entry) {
            var h = entry.hall || {};
            var name = (h.name || 'Mangekyo Cinema');
            var lower = name.toLowerCase();
            var type = 'Standard';
            if (lower.indexOf('imax') !== -1) type = 'IMAX';
            else if (lower.indexOf('vip') !== -1) type = 'VIP';
            return {
              id: h.id,
              name: name,
              photoUrl: '',
              type: type,
              address: h.location || 'Main building',
              hallsSummary: 'Rows: ' + (h.total_rows || 0) + ', seats/row: ' + (h.seats_per_row || 0),
              upcomingSessionsCount: entry.sessions.length
            };
          });
        });
      });
    },

    /* Auth */
    login: function (payload) {
      return request('POST', '/auth/login', { body: payload }).then(function (res) {
        if (!res.ok) throw new Error(res.data && res.data.error || 'Login failed');
        return res.data;
      });
    },
    register: function (payload) {
      return request('POST', '/auth/register', { body: payload }).then(function (res) {
        if (!res.ok) throw new Error(res.data && res.data.error || 'Registration failed');
        return res.data;
      });
    },

    /* Bookings */
    createBooking: function (payload) {
      return request('POST', '/bookings', {
        headers: authHeaders(),
        body: payload
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Booking failed');
        return res.data;
      });
    },
    cancelBooking: function (id) {
      return request('DELETE', '/bookings/' + encodeURIComponent(id), {
        headers: authHeaders()
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Cancel failed');
        return res.data;
      });
    },
    fetchMyBookings: function () {
      return request('GET', '/bookings/my', {
        headers: authHeaders()
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load bookings');
        return Array.isArray(res.data) ? res.data : [];
      });
    },

    /* Reviews (ratings) */
    fetchMovieReviews: function (movieId) {
      return request('GET', '/reviews/movie/' + encodeURIComponent(movieId), {
        headers: authHeaders()
      }).then(function (res) {
        if (!res.ok) return [];
        return Array.isArray(res.data) ? res.data : [];
      });
    },
    createReview: function (payload) {
      return request('POST', '/reviews', {
        headers: authHeaders(),
        body: payload
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to submit review');
        return res.data;
      });
    }
  };
})();

