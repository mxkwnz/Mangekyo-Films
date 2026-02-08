(function () {
  var API_BASE = window.API_BASE || '/api';

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
      if (!u || !u.id) return {};
      return {
        'X-User-ID': u.id,
        'X-User-Role': u.role || 'USER',
        'Content-Type': 'application/json'
      };
    },
    /** Role checks: GUEST (not logged in), USER (customer), ADMIN */
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
    closeBtn.setAttribute('aria-label', 'Закрыть');
    closeBtn.textContent = '\u00D7';
    closeBtn.onclick = function () { window.closeAuthModal(); };

    var title = document.createElement('h2');
    title.id = 'auth-modal-title';
    title.textContent = 'Вход';

    var tabs = document.createElement('div');
    tabs.className = 'auth-modal-tabs';
    var btnLogin = document.createElement('button');
    btnLogin.type = 'button';
    btnLogin.textContent = 'Вход';
    btnLogin.onclick = function () { switchMode('login'); };
    var btnReg = document.createElement('button');
    btnReg.type = 'button';
    btnReg.textContent = 'Регистрация';
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
    first.placeholder = 'Имя';
    first.required = true;
    var last = document.createElement('input');
    last.type = 'text';
    last.name = 'last_name';
    last.placeholder = 'Фамилия';
    last.required = true;
    row1.appendChild(first);
    row1.appendChild(last);
    registerFields.appendChild(row1);
    var phone = document.createElement('input');
    phone.type = 'tel';
    phone.name = 'phone_number';
    phone.placeholder = 'Телефон';
    phone.required = true;
    registerFields.appendChild(phone);

    var email = document.createElement('input');
    email.type = 'email';
    email.name = 'email';
    email.placeholder = 'Email';
  email.required = true;

    var pass = document.createElement('input');
    pass.type = 'password';
    pass.name = 'password';
    pass.placeholder = 'Пароль';
    pass.required = true;

    var submit = document.createElement('button');
    submit.type = 'submit';
    submit.id = 'auth-submit-btn';
    submit.textContent = 'Войти';

    var divider = document.createElement('p');
    divider.className = 'auth-divider';
    divider.textContent = 'или';

    var googleBtn = document.createElement('button');
    googleBtn.type = 'button';
    googleBtn.className = 'auth-google-btn';
    googleBtn.textContent = 'Войти через Google (скоро)';
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
      errorEl.textContent = msg || '';
      errorEl.style.display = msg ? 'block' : 'none';
    }

    function switchMode(mode) {
      currentMode = mode;
      showError('');
      if (mode === 'login') {
        document.getElementById('auth-modal-title').textContent = 'Вход';
        document.getElementById('auth-register-fields').style.display = 'none';
        submit.textContent = 'Войти';
        tabs.querySelectorAll('button')[0].classList.add('active');
        tabs.querySelectorAll('button')[1].classList.remove('active');
        first.removeAttribute('required');
        last.removeAttribute('required');
        phone.removeAttribute('required');
      } else {
        document.getElementById('auth-modal-title').textContent = 'Регистрация';
        document.getElementById('auth-register-fields').style.display = 'flex';
        document.getElementById('auth-register-fields').style.flexDirection = 'column';
        submit.textContent = 'Зарегистрироваться';
        tabs.querySelectorAll('button')[1].classList.add('active');
        tabs.querySelectorAll('button')[0].classList.remove('active');
        first.setAttribute('required', 'required');
        last.setAttribute('required', 'required');
        phone.setAttribute('required', 'required');
      }
    }

    form.onsubmit = function (e) {
      e.preventDefault();
      showError('');
      submit.disabled = true;

      var payload = {
        email: form.email.value.trim(),
        password: form.password.value
      };
      var url = API_BASE + '/auth/login';
      if (currentMode === 'register') {
        url = API_BASE + '/auth/register';
        payload.first_name = form.first_name.value.trim();
        payload.last_name = form.last_name.value.trim();
        payload.phone_number = form.phone_number.value.trim();
      }

      fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })
        .then(function (res) { return res.json().then(function (data) { return { ok: res.ok, data: data }; }); })
        .then(function (result) {
          if (!result.ok) {
            showError(result.data.error || 'Ошибка входа');
            submit.disabled = false;
            return;
          }
          var user = result.data.user;
          setStored({
            id: user.id,
            role: user.role || 'USER',
            email: user.email,
            first_name: user.first_name,
            last_name: user.last_name
          });
          window.closeAuthModal();
          if (typeof currentCb === 'function') currentCb();
          submit.disabled = false;
        })
        .catch(function () {
          showError('Ошибка сети. Попробуйте позже.');
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

    switchMode('login');
    document.body.appendChild(overlay);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', buildModal);
  } else {
    buildModal();
  }
})();
