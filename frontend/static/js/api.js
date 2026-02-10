(function () {
  var API_BASE = window.API_BASE || '/api';

  function withJson(response) {
    return response.text().then(function (text) {
      var data = null;
      try {
        if (text) data = JSON.parse(text);
      } catch (e) {
        return { ok: false, status: response.status, data: { error: 'Invalid server response: ' + text.substring(0, 100) } };
      }
      return { ok: response.ok, status: response.status, data: data };
    });
  }

  function authHeaders() {
    var u = window.auth && window.auth.getUser ? window.auth.getUser() : null;
    if (u && u.token) {
      return { 'Authorization': 'Bearer ' + u.token };
    }
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
    }).then(function (response) {
      if (response.status === 401) {

      }
      return withJson(response);
    });
  }

  window.api = {

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
    },
    fetchMyReviews: function () {
      return request('GET', '/reviews/my', {
        headers: authHeaders()
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load reviews');
        return Array.isArray(res.data) ? res.data : [];
      });
    },
    updateReview: function (id, payload) {
      return request('PUT', '/reviews/' + encodeURIComponent(id), {
        headers: authHeaders(),
        body: payload
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to update review');
        return res.data;
      });
    },
    deleteReview: function (id) {
      return request('DELETE', '/reviews/' + encodeURIComponent(id), {
        headers: authHeaders()
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to delete review');
        return res.data;
      });
    },


    fetchMyCards: function () {
      return request('GET', '/payment-cards', { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load cards');
        var data = res.data;
        if (data && data.cards) return data.cards;
        return Array.isArray(data) ? data : [];
      });
    },
    createCard: function (payload) {
      return request('POST', '/payment-cards', { headers: authHeaders(), body: payload }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to create card');
        return res.data;
      });
    },
    deleteCard: function (id) {
      return request('DELETE', '/payment-cards/' + encodeURIComponent(id), { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to delete card');
        return res.data;
      });
    },


    topUpBalance: function (payload) {
      return request('POST', '/payments/topup', { headers: authHeaders(), body: payload }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Top-up failed');
        return res.data;
      });
    },
    fetchMyPayments: function () {
      return request('GET', '/payments', { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load payments');
        return Array.isArray(res.data) ? res.data : [];
      });
    },


    adminFetchHalls: function () {
      return request('GET', '/admin/halls', { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load halls');
        return Array.isArray(res.data) ? res.data : [];
      });
    },
    adminCreateHall: function (payload) {
      return request('POST', '/admin/halls', { headers: authHeaders(), body: payload }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to create hall');
        return res.data;
      });
    },
    adminDeleteHall: function (id) {
      return request('DELETE', '/admin/halls/' + encodeURIComponent(id), { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to delete hall');
        return res.data;
      });
    },
    adminCreateMovie: function (payload) {
      return request('POST', '/admin/movies', { headers: authHeaders(), body: payload }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to create movie');
        return res.data;
      });
    },
    adminDeleteMovie: function (id) {
      return request('DELETE', '/admin/movies/' + encodeURIComponent(id), { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to delete movie');
        return res.data;
      });
    },
    adminCreateSession: function (payload) {
      return request('POST', '/admin/sessions', { headers: authHeaders(), body: payload }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to create session');
        return res.data;
      });
    },
    adminDeleteSession: function (id) {
      return request('DELETE', '/admin/sessions/' + encodeURIComponent(id), { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to delete session');
        return res.data;
      });
    },
    adminFetchBookingsBySession: function (sessionId) {
      return request('GET', '/admin/bookings/session/' + sessionId, { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load session bookings');
        return Array.isArray(res.data) ? res.data : [];
      });
    },


    fetchMe: function () {
      return request('GET', '/auth/me', { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load profile');
        return res.data.user || res.data;
      });
    },
    updateProfile: function (payload) {
      return request('PUT', '/auth/me', {
        headers: authHeaders(),
        body: payload
      }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to update profile');
        return res.data;
      });
    },

    adminFetchBookings: function () {
      return request('GET', '/admin/bookings', { headers: authHeaders() }).then(function (res) {
        if (!res.ok) throw new Error(res.data.error || 'Failed to load bookings');
        return Array.isArray(res.data) ? res.data : [];
      });
    }
  };
})();
