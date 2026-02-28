/* shared.js — Village Square common utilities */
var VS = (function () {
  var MAX_TOASTS = 3;

  return {
    /* ---- Toast notification ---- */
    toast: function (message, type) {
      type = type || 'info';
      var container = document.getElementById('toastContainer');
      if (!container) return;

      // Limit visible toasts
      while (container.children.length >= MAX_TOASTS) {
        container.removeChild(container.lastChild);
      }

      var el = document.createElement('div');
      el.className = 'toast ' + type;
      el.innerHTML =
        '<span>' + VS.escapeHTML(message) + '</span>' +
        '<button class="toast-close" aria-label="Close">&times;</button>' +
        '<div class="toast-progress"></div>';

      container.insertBefore(el, container.firstChild);

      var dismiss = function () {
        if (!el.parentNode) return;
        el.style.animation = 'toastFadeOut 0.3s ease forwards';
        setTimeout(function () { if (el.parentNode) el.parentNode.removeChild(el); }, 300);
      };

      el.querySelector('.toast-close').addEventListener('click', dismiss);

      var timer = setTimeout(dismiss, 4000);
      el.addEventListener('mouseenter', function () { clearTimeout(timer); });
      el.addEventListener('mouseleave', function () { timer = setTimeout(dismiss, 2000); });
    },

    /* backward compat — redirect to toast */
    showBanner: function (message, type) {
      VS.toast(message, type);
    },

    /* ---- Auth guard ---- */
    authGuard: function (onSuccess) {
      fetch('/api/me', { credentials: 'same-origin' })
        .then(function (r) {
          if (!r.ok) { window.location.href = '/index.html'; return null; }
          return r.json();
        })
        .then(function (user) { if (user) onSuccess(user); })
        .catch(function () { window.location.href = '/index.html'; });
    },

    /* ---- Logout ---- */
    setupLogout: function () {
      var btn = document.getElementById('logoutBtn');
      if (!btn) return;
      btn.addEventListener('click', function () {
        btn.disabled = true;
        fetch('/api/logout', { method: 'POST', credentials: 'same-origin' })
          .then(function () { window.location.href = '/index.html'; })
          .catch(function () {
            VS.toast('Logout failed. Please try again.', 'error');
            btn.disabled = false;
          });
      });
    },

    /* ---- Time helpers ---- */
    timeAgo: function (dateStr) {
      var secs = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000);
      if (secs < 60) return 'just now';
      var mins = Math.floor(secs / 60);
      if (mins < 60) return mins + (mins === 1 ? ' minute ago' : ' minutes ago');
      var hrs = Math.floor(mins / 60);
      if (hrs < 24) return hrs + (hrs === 1 ? ' hour ago' : ' hours ago');
      var days = Math.floor(hrs / 24);
      if (days < 30) return days + (days === 1 ? ' day ago' : ' days ago');
      var months = Math.floor(days / 30);
      return months + (months === 1 ? ' month ago' : ' months ago');
    },

    escapeHTML: function (str) {
      var d = document.createElement('div');
      d.appendChild(document.createTextNode(str));
      return d.innerHTML;
    },

    /* ---- Inline confirm (replaces window.confirm) ---- */
    inlineConfirm: function (btn, onConfirm) {
      if (btn.dataset.confirming) return;
      btn.dataset.confirming = '1';

      var origHTML = btn.innerHTML;
      var origClass = btn.className;

      var group = document.createElement('span');
      group.className = 'inline-confirm-group';

      var cancelBtn = document.createElement('button');
      cancelBtn.type = 'button';
      cancelBtn.className = 'inline-confirm-cancel';
      cancelBtn.textContent = 'Cancel';

      var confirmBtn = document.createElement('button');
      confirmBtn.type = 'button';
      confirmBtn.className = 'inline-confirm-yes';
      confirmBtn.textContent = 'Confirm';

      group.appendChild(cancelBtn);
      group.appendChild(confirmBtn);

      btn.style.display = 'none';
      btn.parentNode.insertBefore(group, btn.nextSibling);

      var revert = function () {
        clearTimeout(timeout);
        if (group.parentNode) group.parentNode.removeChild(group);
        btn.style.display = '';
        delete btn.dataset.confirming;
      };

      var timeout = setTimeout(revert, 3000);

      cancelBtn.addEventListener('click', function (e) { e.stopPropagation(); revert(); });
      confirmBtn.addEventListener('click', function (e) {
        e.stopPropagation();
        revert();
        onConfirm();
      });
    },

    /* ---- Skeleton loader helper ---- */
    skeletonCards: function (count, cssClass) {
      cssClass = cssClass || 'skeleton-card';
      var html = '';
      for (var i = 0; i < count; i++) {
        html += '<div class="skeleton ' + cssClass + '"></div>';
      }
      return html;
    },

    /* ---- Wire up input listeners to clear field errors ---- */
    clearErrorOnInput: function (inputEl, errorEl) {
      var handler = function () {
        errorEl.textContent = '';
        errorEl.classList.remove('visible');
        inputEl.classList.remove('input-error');
      };
      inputEl.addEventListener('input', handler);
      inputEl.addEventListener('change', handler);
    },

    /* ---- Hamburger menu setup ---- */
    setupHamburger: function () {
      var btn = document.getElementById('hamburgerBtn');
      var menu = document.getElementById('headerMenu');
      if (!btn || !menu) return;

      btn.addEventListener('click', function (e) {
        e.stopPropagation();
        var isOpen = menu.classList.toggle('open');
        btn.classList.toggle('open', isOpen);
        btn.setAttribute('aria-expanded', isOpen ? 'true' : 'false');
      });

      // Close on outside click
      document.addEventListener('click', function (e) {
        if (!menu.classList.contains('open')) return;
        if (!menu.contains(e.target) && e.target !== btn) {
          menu.classList.remove('open');
          btn.classList.remove('open');
          btn.setAttribute('aria-expanded', 'false');
        }
      });

      // Close when nav link clicked
      var links = menu.querySelectorAll('.nav-link');
      for (var i = 0; i < links.length; i++) {
        links[i].addEventListener('click', function () {
          menu.classList.remove('open');
          btn.classList.remove('open');
          btn.setAttribute('aria-expanded', 'false');
        });
      }
    }
  };
})();
