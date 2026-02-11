(function () {
  function getStored() {
    try {
      var raw = localStorage.getItem('cinema_user');
      return raw ? JSON.parse(raw) : null;
    } catch (_) {
      return null;
    }
  }

  function setStored(user) {
    if (user) localStorage.setItem('cinema_user', JSON.stringify(user));
    else localStorage.removeItem('cinema_user');
  }

  window.auth = {
    isLoggedIn: function () {
      return !!getStored();
    },
    getUser: function () {
      return getStored();
    },
    logout: function () {
      setStored(null);
    },
    getAuthHeaders: function () {
      var u = getStored();
      if (!u) return {};
      var headers = { 'Content-Type': 'application/json' };
      if (u.token) {
        headers['Authorization'] = 'Bearer ' + u.token;
      } else if (u.id) {
        headers['X-User-ID'] = u.id;
        headers['X-User-Role'] = u.role || 'USER';
      }
      return headers;
    },

    hasRole: function (role) {
      var u = getStored();
      if (!u) return role === 'GUEST';
      return (u.role || '').toUpperCase() === (role || '').toUpperCase();
    },
    isAdmin: function () {
      return this.hasRole('ADMIN');
    },
    isCustomer: function () {
      var u = getStored();
      return u && ((u.role || '').toUpperCase() === 'USER' || (u.role || '').toUpperCase() === 'ADMIN');
    }
  };

  function buildModal() {
    var overlay = document.createElement('div');
    overlay.className = 'auth-modal-overlay';
    overlay.id = 'auth-modal-overlay';

    var wrap = document.createElement('div');
    wrap.className = 'auth-modal-wrap';

    var closeBtn = document.createElement('button');
    closeBtn.type = 'button';
    closeBtn.className = 'auth-modal-close';
    closeBtn.setAttribute('aria-label', 'Close');
    closeBtn.textContent = '\u00D7';
    closeBtn.onclick = function () { window.closeAuthModal(); };

    var title = document.createElement('h2');
    title.id = 'auth-modal-title';
    title.textContent = 'Login';

    var tabs = document.createElement('div');
    tabs.className = 'auth-modal-tabs';
    var btnLogin = document.createElement('button');
    btnLogin.type = 'button';
    btnLogin.textContent = 'Login';
    btnLogin.onclick = function () { switchMode('login'); };
    var btnReg = document.createElement('button');
    btnReg.type = 'button';
    btnReg.textContent = 'Register';
    btnReg.onclick = function () { switchMode('register'); };
    tabs.appendChild(btnLogin);
    tabs.appendChild(btnReg);

    var form = document.createElement('form');
    form.className = 'auth-form';
    form.id = 'auth-form';

    var errorEl = document.createElement('p');
    errorEl.className = 'auth-error';
    errorEl.id = 'auth-form-error';
    errorEl.style.display = 'none';

    var registerFields = document.createElement('div');
    registerFields.id = 'auth-register-fields';
    registerFields.style.display = 'none';
    var row1 = document.createElement('div');
    row1.className = 'row';
    var first = document.createElement('input');
    first.type = 'text';
    first.name = 'first_name';
    first.placeholder = 'First name';
    first.required = true;
    var last = document.createElement('input');
    last.type = 'text';
    last.name = 'last_name';
    last.placeholder = 'Last name';
    last.required = true;
    row1.appendChild(first);
    row1.appendChild(last);
    registerFields.appendChild(row1);
    var phoneWrapper = document.createElement('div');
    phoneWrapper.className = 'phone-input-wrapper';

    var phonePrefix = document.createElement('div');
    phonePrefix.className = 'phone-prefix';
    phonePrefix.id = 'phone-country-selector';
    phonePrefix.innerHTML = '<img src="https://flagcdn.com/w20/us.png" alt="Flag"> <span class="prefix-text">+1</span> <div class="prefix-arrow"></div>';

    var countryDropdown = document.createElement('div');
    countryDropdown.className = 'country-dropdown';
    var countries = [
      { code: 'us', prefix: '+1', name: 'USA' },
      { code: 'kz', prefix: '+7', name: 'Kazakhstan' },
      { code: 'ru', prefix: '+7', name: 'Russia' },
      { code: 'tr', prefix: '+90', name: 'Turkey' },
      { code: 'ae', prefix: '+971', name: 'UAE' }
    ];

    countries.forEach(function (c) {
      var item = document.createElement('div');
      item.className = 'country-item';
      item.innerHTML = `<img src="https://flagcdn.com/w20/${c.code}.png"> ${c.name} (${c.prefix})`;
      item.onclick = function (e) {
        e.stopPropagation();
        phonePrefix.querySelector('img').src = `https://flagcdn.com/w20/${c.code}.png`;
        phonePrefix.querySelector('.prefix-text').textContent = c.prefix;
        countryDropdown.classList.remove('active');
      };
      countryDropdown.appendChild(item);
    });

    phonePrefix.appendChild(countryDropdown);
    phonePrefix.onclick = function (e) {
      e.stopPropagation();
      countryDropdown.classList.toggle('active');
    };

    document.addEventListener('click', function () {
      countryDropdown.classList.remove('active');
    });

    var phone = document.createElement('input');
    phone.type = 'tel';
    phone.name = 'phone_number';
    phone.placeholder = 'Phone number';
    phone.required = true;
    phone.oninput = function (e) {
      this.value = this.value.replace(/[^0-9]/g, '');
    };

    phoneWrapper.appendChild(phonePrefix);
    phoneWrapper.appendChild(phone);
    registerFields.appendChild(phoneWrapper);

    var email = document.createElement('input');
    email.type = 'email';
    email.name = 'email';
    email.placeholder = 'Email';
    email.required = true;

    var pass = document.createElement('input');
    pass.type = 'password';
    pass.name = 'password';
    pass.placeholder = 'Password';
    pass.required = true;

    var submit = document.createElement('button');
    submit.type = 'submit';
    submit.id = 'auth-submit-btn';
    submit.textContent = 'Log in';

    var divider = document.createElement('p');
    divider.className = 'auth-divider';
    divider.textContent = 'or';

    var googleBtn = document.createElement('button');
    googleBtn.type = 'button';
    googleBtn.className = 'auth-google-btn';
    googleBtn.textContent = 'Sign in with Google (coming soon)';
    googleBtn.disabled = true;

    form.appendChild(errorEl);
    form.appendChild(registerFields);
    form.appendChild(email);
    form.appendChild(pass);
    form.appendChild(submit);
    form.appendChild(divider);
    form.appendChild(googleBtn);

    var modal = document.createElement('div');
    modal.className = 'auth-modal';
    modal.appendChild(wrap);
    wrap.appendChild(closeBtn);
    wrap.appendChild(title);
    wrap.appendChild(tabs);
    wrap.appendChild(form);

    overlay.appendChild(modal);
    overlay.addEventListener('click', function (e) {
      if (e.target === overlay) window.closeAuthModal();
    });

    var currentMode = 'login';
    var currentCb = null;

    function showError(msg) {
      if (msg) console.error('Auth Error:', msg);
      errorEl.textContent = typeof msg === 'string' ? msg : (msg.message || 'An unknown error occurred');
      errorEl.style.display = msg ? 'block' : 'none';
    }

    function switchMode(mode) {
      currentMode = mode;
      showError('');
      if (mode === 'login') {
        title.textContent = 'Login';
        registerFields.style.display = 'none';
        submit.textContent = 'Log in';
        tabs.querySelectorAll('button')[0].classList.add('active');
        tabs.querySelectorAll('button')[1].classList.remove('active');
        first.removeAttribute('required');
        last.removeAttribute('required');
        phone.removeAttribute('required');
      } else {
        title.textContent = 'Register';
        registerFields.style.display = 'flex';
        registerFields.style.flexDirection = 'column';
        submit.textContent = 'Register';
        tabs.querySelectorAll('button')[1].classList.add('active');
        tabs.querySelectorAll('button')[0].classList.remove('active');
        first.setAttribute('required', 'required');
        last.setAttribute('required', 'required');
        phone.setAttribute('required', 'required');
      }
    }

    function handleAuthSuccess(res) {
      if (!res) return;
      var user = res.user || res.data || res;

      var toStore = {
        id: user.id || user._id,
        role: user.role || 'USER',
        email: user.email,
        first_name: user.first_name,
        last_name: user.last_name,
        token: res.token || user.token
      };

      setStored(toStore);
    }

    form.onsubmit = function (e) {
      e.preventDefault();
      showError('');
      submit.disabled = true;

      var payload = {
        email: form.email.value.trim(),
        password: form.password.value
      };
      if (currentMode === 'register') {
        payload.first_name = form.first_name.value.trim();
        payload.last_name = form.last_name.value.trim();
        var prefix = document.querySelector('.prefix-text').textContent;
        payload.phone_number = prefix + ' ' + form.phone_number.value.trim();
      }

      var apiFn = null;
      if (window.api) {
        apiFn = currentMode === 'register' ? window.api.register : window.api.login;
      }

      if (!apiFn) {
        showError('Auth API not available.');
        submit.disabled = false;
        return;
      }

      apiFn(payload)
        .then(function (res) {

          handleAuthSuccess(res);
          window.closeAuthModal();

          if (typeof currentCb === 'function') {
            currentCb();
          } else {

            location.reload();
          }
          submit.disabled = false;
        })
        .catch(function (err) {

          showError((err && err.message) || 'Network error. Please try again.');
          submit.disabled = false;
        });
    };

    window.openAuthModal = function (options) {
      options = options || {};
      currentCb = options.onSuccess || null;
      switchMode(options.mode === 'register' ? 'register' : 'login');
      showError('');
      document.getElementById('auth-modal-overlay').classList.add('is-open');
    };

    window.closeAuthModal = function () {
      document.getElementById('auth-modal-overlay').classList.remove('is-open');
    };

    document.body.appendChild(overlay);
    switchMode('login');
  }

  function highlightActiveLink() {
    var path = window.location.pathname;
    var filename = path.substring(path.lastIndexOf('/') + 1) || 'index.html';
    var links = document.querySelectorAll('.site-nav a');
    links.forEach(function (link) {
      var href = link.getAttribute('href');
      if (href === filename || (filename === 'index.html' && href === 'index.html')) {
        link.classList.add('active');
      } else {
        link.classList.remove('active');
      }
    });
  }

  function init() {
    buildModal();
    highlightActiveLink();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
