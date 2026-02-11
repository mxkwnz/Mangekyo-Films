(function () {
    function updateNav() {
        var span = document.getElementById('user-span');
        var navAdmin = document.getElementById('nav-admin');
        if (window.auth && auth.isAdmin && navAdmin) navAdmin.style.display = auth.isAdmin() ? '' : 'none';
        if (!span) return;
        if (window.auth && auth.isLoggedIn && auth.isLoggedIn()) {
            var u = auth.getUser();
            span.innerHTML = 'Hi, ' + (u.first_name || u.email) + ' | <button type="button" id="logout-btn" class="btn btn-ghost">Log out</button>';
            var btn = document.getElementById('logout-btn');
            if (btn) btn.onclick = function () { auth.logout(); location.href = 'index.html'; };
        } else {
            span.innerHTML = '<button type="button" id="login-btn" class="btn btn-primary">Login / Register</button>';
            var loginBtn = document.getElementById('login-btn');
            if (loginBtn) loginBtn.onclick = function () { if (window.openAuthModal) openAuthModal({ onSuccess: function () { location.reload(); } }); };
        }
    }

    function showError(msg, isSuccess) {
        if (msg) {
            window.showToast(msg, isSuccess ? 'success' : 'error');
        }
        var el = document.getElementById('topup-error');
        if (!el) return;
        el.textContent = msg || '';
        el.style.display = msg ? 'block' : 'none';
        el.style.color = isSuccess ? '#10b981' : '#ef4444';
    }

    function loadCards() {
        window.api.fetchMyCards()
            .then(function (cards) {
                var list = document.getElementById('cards-list');
                var select = document.getElementById('select-card');
                if (!list || !select) return;

                list.innerHTML = '';
                select.innerHTML = '<option value="">Select a card</option>';

                if (cards.length === 0) {
                    list.innerHTML = '<p style="color: var(--text-dim)">No cards saved.</p>';
                }

                cards.forEach(function (card) {
                    var cardEl = document.createElement('div');
                    cardEl.className = 'card';
                    cardEl.style.border = '1px solid var(--border)';
                    cardEl.innerHTML = `
                        <div class="card-body">
                            <div class="card-title">**** **** **** ${card.card_number.slice(-4)}</div>
                            <div class="card-meta">${card.card_holder_name} | Exp: ${card.expiry_date}</div>
                            <button type="button" class="btn btn-ghost btn-del-card" data-id="${card.id}" style="margin-top: 1rem; color: #ef4444;">Remove Card</button>
                        </div>
                    `;
                    list.appendChild(cardEl);

                    var opt = document.createElement('option');
                    opt.value = card.id;
                    opt.textContent = '****' + card.card_number.slice(-4) + ' (' + card.card_holder_name + ')';
                    select.appendChild(opt);
                });

                list.querySelectorAll('.btn-del-card').forEach(function (btn) {
                    btn.onclick = function () {
                        var id = btn.dataset.id;
                        if (confirm('Remove this card?')) {
                            window.api.deleteCard(id).then(loadCards).catch(function (err) { showError(err.message); });
                        }
                    };
                });
            })
            .catch(function (err) { console.error('Cards load failed', err); });
    }

    function init() {
        updateNav();
        if (!(window.auth && auth.isLoggedIn && auth.isLoggedIn())) {
            showError('Please log in to manage your balance.');
            return;
        }
        showError('');
        loadCards();

        document.getElementById('form-add-card').onsubmit = function (e) {
            e.preventDefault();
            var payload = {
                card_holder_name: this.card_holder_name.value.trim(),
                card_number: this.card_number.value.trim(),
                expiry_date: this.expiry_date.value.trim(),
                cvv: this.cvv.value.trim()
            };
            window.api.createCard(payload)
                .then(function () {
                    document.getElementById('form-add-card').reset();
                    loadCards();
                    showError('Card added successfully!', true);
                })
                .catch(function (err) { showError(err.message); });
        };

        document.getElementById('form-topup').onsubmit = function (e) {
            e.preventDefault();
            var payload = {
                payment_card_id: this.card_id.value,
                amount: parseFloat(this.amount.value)
            };
            window.api.topUpBalance(payload)
                .then(function () {
                    showError('Balance topped up successfully!', true);
                    this.reset();
                }.bind(this))
                .catch(function (err) { showError(err.message); });
        };

    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
