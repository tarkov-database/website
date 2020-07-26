let map = {};

const initMap = async () => {
  const el = document.getElementById('map');
  if (el === null) return;

  const libPath = document.getElementById('mapLib').src;

  try {
    map = await import(libPath);
    await map.init(el);
  } catch (err) {
    console.error(err);
  }
};

initMap();


const registerTabs = () => {
  const openTab = e => {
    const cl = 'active';
    document.querySelector(`.tab.${cl}`).classList.remove(cl);
    document.querySelector(`.tab-btn.${cl}`).classList.remove(cl);

    const el = e.currentTarget;
    document.getElementById(el.dataset.tab).classList.add(cl);
    el.classList.add(cl);
  };

  const tablinks = document.getElementsByClassName('tab-btn');
  for (const el of tablinks) el.addEventListener('click', openTab);
};

registerTabs();


const loadImage = async () => {
  const image = document.getElementById('itemImage');
  if (image === null) return;

  const staticURL = (() => {
    const host = window.location.host;
    const parts = host.split('.');
    return `//static.${parts.length > 2 ? [parts[parts.length-2], parts[parts.length-1]].join('.'): host}`;
  })();
  const imageID = image.dataset.id;
  const request = new Request(`${staticURL}/image/icon/1-1/${imageID}.png`);

  try {
    const response = await fetch(request);
    if(!response.ok) throw new Error(response.status);

    const objectURL = URL.createObjectURL(await response.blob());
    const img = new Image();
    img.src = objectURL;
    img.onload = () => {
      const imgWidth = img.naturalWeight, imgHeight = img.naturalHeight;
      const boxWidth = image.boxWidth, boxHeight = image.offsetHeight;

      if (imgHeight > boxHeight || imgWidth > boxWidth) image.style.backgroundSize = 'contain';

      image.style.backgroundImage = `url("${objectURL}")`;
    };
  } catch (err) {
    if (err.message === '404') {
      console.warn('Image can not be found');
    } else {
      console.error('Image load failed:', err.message);
    }
  }
};

loadImage();


const sortTables = () => {
  const tables = document.querySelectorAll('.sort-table.client-sort thead th');
  if (tables.len === 0) return;

  const getCellValue = (tr, idx) =>
    tr.children[idx].dataset.value || tr.children[idx].innerText || tr.children[idx].textContent;

  const comparer = (idx, asc) => (a, b) => ((v1, v2) =>
    v1 !== '' && v2 !== '' && !isNaN(v1) && !isNaN(v2) ? v1 - v2 : v1.toString().localeCompare(v2)
  )(getCellValue(asc ? a : b, idx), getCellValue(asc ? b : a, idx));


  for (const th of tables) {
    th.addEventListener('click', (e => {
      const el = e.currentTarget;
      const table = el.closest('table');
      const body = table.getElementsByTagName('tbody')[0];

      const clSorted = 'sorted-by';
      const clAsc = 'up';
      const clDesc = 'down';

      if (el.asc === undefined) {
        if (el.classList.contains(clAsc)) {
          el.asc = true;
        } else if (el.classList.contains(clDesc)) {
          el.asc = false;
        }
      }

      Array.from(body.getElementsByTagName('tr'))
        .sort(comparer(Array.from(el.parentNode.children).indexOf(th), el.asc = !el.asc))
        .forEach(tr => body.appendChild(tr));

      const sorted = table.getElementsByClassName(clSorted)[0];
      if (th !== sorted) {
        sorted.classList.remove(clSorted);
        el.classList.add(clSorted);
      }

      if (el.asc) {
        if (el.classList.contains(clDesc)) {
          el.classList.replace(clDesc, clAsc);
        } else {
          el.classList.add(clAsc);
        }
      } else {
        if (el.classList.contains(clAsc)) {
          el.classList.replace(clAsc, clDesc);
        } else {
          el.classList.add(clDesc);
        }
      }
    }));

    th.style.cursor = 'pointer';
  }
};

sortTables();


const initSearchSocket = async() => {
  const form = document.getElementById('search');
  if (form === null) return;

  const input = form.querySelector('input[type="search"]');
  const suggBox = form.querySelector('.suggestion');
  const suggInline = form.querySelector('.inline-suggestion');

  let socket;
  let idleTimeout;
  let noReconnect = false;
  let qCount = 0;

  let lastTerm = '';

  const wait = time => new Promise(resolve => setTimeout(() => resolve(), time));

  const regexMeta = new RegExp(/[.*+?^${}()|[\]\\]/, 'g');
  const quoteMeta = str => str.replace(regexMeta, '\\$&');

  const regexFilterPrefix = new RegExp(/(?:^\s*(?<key>\w+):)(?:\s*(?<value>\w+)\s)?(?:\s*(?<term>.*))/, 'i');
  const getFilterPrefix = str => regexFilterPrefix.exec(str);

  const showElement = el => el.classList.replace('hide', 'show');
  const hideElement = el => el.classList.replace('show', 'hide');

  const isInteractiveMap = window.location.pathname.endsWith('/map');

  const errErrorClosure = new Error('Search socket was closed with an error');

  const updateBoxSuggestions = (items, term) => {
    const newUl = document.createElement('ul');

    for (const item of items) {
      const a = document.createElement('a');
      const li = document.createElement('li');
      const text = document.createElement('span');
      const div = document.createElement('div');

      const highlightMatches = str => {
        const regex = new RegExp(`(${term})`, 'gi');
        const matches = str.match(regex);
        if (matches) {
          for (const m of matches) {
            str = str.replace(m, `<b>${m}</b>`);
          }
        }

        return str;
      };

      text.innerHTML = highlightMatches(item.name);
      a.title = item.name;

      switch (item.type) {
      case 0:
        a.href = `/item/${item.parent}/${item.id}`;
        div.className = `icon ${item.parent}`;
        break;
      case 1:
        a.href = `/location/${item.id}`;
        a.href += isInteractiveMap ? '/map': '';
        div.className = 'icon location';
        break;
      case 2:
        a.addEventListener('click', async () => {
          hideElement(suggBox);
          map.flyToFeature(await map.getFeature(item.id, item.parent));
        });
        a.href = `#feature=${item.id}`;
        div.className = 'icon feature';
        break;
      }

      div.innerHTML = '&nbsp;';

      a.addEventListener('keydown', e => {
        const current = e.target.parentNode;

        let next = null;
        switch (e.key) {
        case 'ArrowUp':
          e.preventDefault();
          next = current.previousElementSibling;
          if (!next) input.focus();
          break;
        case 'ArrowDown':
          e.preventDefault();
          next = current.nextElementSibling;
          if (!next) input.focus();
          break;
        case 'Escape':
          current.blur();
          hideElement(suggBox);
          break;
        default:
          return;
        }

        if (next) next.querySelector('a').focus();
      });

      a.appendChild(div);
      a.append(text);
      li.appendChild(a);
      newUl.appendChild(li);
    }

    const ul = suggBox.querySelector('ul');
    suggBox.replaceChild(newUl, ul);
  };

  const updateInlineSuggestion = val => {
    suggInline.innerHTML = val.replace(' ', '&nbsp;');
    suggInline.dataset.value = val;
  };

  const filters = {
    _data: {
      item: {
        available: !isInteractiveMap,
        values: [
          'ammunition',
          'armor',
          'backpack',
          'barrel',
          'barter',
          'bipod',
          'charge',
          'clothing',
          'common',
          'container',
          'device',
          'firearm',
          'food',
          'foregrip',
          'gasblock',
          'goggles',
          'grenade',
          'handguard',
          'headphone',
          'key',
          'launcher',
          'magazine',
          'map',
          'medical',
          'melee',
          'mod-other',
          'money',
          'mount',
          'muzzle',
          'pistolgrip',
          'receiver',
          'sight-special',
          'sight',
          'stock',
          'tacticalrig'
        ],
      }
    },

    contains(key) {
      return Object.prototype.hasOwnProperty.call(this._data, key) && 
        this._data[key].available;
    },

    includes(key, value) {
      const entry = this._data[key];
      return entry.available && entry.values.includes(value);
    },

    firstMatch(key)  {
      const match = Object.entries(this._data).find(v => v[1].available && v[0].startsWith(key));
      return match ? match[0]: '';
    },

    firstMatchIn(key, value)  {
      const entry = this._data[key];
      return entry.available ? entry.values.find(v => v.startsWith(value)): '';
    }
  };

  const openSocket = msgListener => {
    const host = window.location.host;
    const path = 'search/ws';

    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    if (proto === 'ws') console.warn('Insecure WebSocket protocol is used');

    let socket;
    try {
      socket = new WebSocket(`${proto}://${host}/${path}`);
    } catch (err) {
      return Promise.reject(err);
    }

    socket.addEventListener('open', () => {
      console.info('Search socket opened');
    });

    socket.addEventListener('close', e => {
      if (e.wasClean) {
        const msg = 'Search socket closed';
        console.info(e.reason ? `${msg}: ${e.reason}`: msg);
        if (e.code !== 1000) noReconnect = true;
        if (e.code === 1012) wait(1000).then(noReconnect = false);
      } else {
        const msg = 'Search socket closed unexpectedly';
        console.error(e.reason ? `${msg}: ${e.reason}`: msg);
        noReconnect = true;
      }
    });

    socket.addEventListener('error', e => {
      console.error('Search socket error: %s', e);
    });

    socket.addEventListener('message', msgListener);

    return socket;
  };

  const connect = async() => {
    if (!socket || socket.readyState > 1) {
      if (noReconnect) return Promise.reject(errErrorClosure);
      try {
        socket = await openSocket(onMessage); // eslint-disable-line require-atomic-updates
      } catch (err) {
        return Promise.reject(err);
      }
    }

    while (socket.readyState === 0) await wait(50);

    clearTimeout(idleTimeout);
    idleTimeout = setTimeout(() => socket.close(1000, 'Idle timeout'), 30*1000);

    return;
  };

  const onMessage = async event => {
    const data = JSON.parse(event.data);

    if (data.error !== null) {
      console.error(data.error);
      return;
    }

    if (data.id < qCount) return;

    if (!data.items || data.items.length === 0) {
      hideElement(suggBox);
      return;
    }

    const term = quoteMeta(lastTerm).replace(' ', '|');
    
    updateBoxSuggestions(data.items, term);
    showElement(suggBox);
  };

  const onInput = async event => {
    const count = qCount;

    try {
      await connect();
    } catch (err) {
      if (err !== errErrorClosure) console.error(err);
      return;
    }

    const input = event.target;
    const val = input.value;

    if (!input.validity.valid) {
      hideElement(suggBox);
      hideElement(suggInline);
      return;
    }

    const keyCompletion = val => {
      const match = filters.firstMatch(val);
      if (match) {
        const v = input.value.replace(val, `${match}:`);
        updateInlineSuggestion(v);
        showElement(suggInline);
      } else {
        hideElement(suggInline);
      }
    };

    const valueCompletion = (key, val) => {
      const match = filters.firstMatchIn(key, val);
      if (match) {
        const v = input.value.replace(val, `${match} `);
        updateInlineSuggestion(v);
        showElement(suggInline);
      } else {
        hideElement(suggInline);
      }
    };

    const prefix = getFilterPrefix(val);

    if (prefix !== null && filters.contains(prefix.groups.key)) {
      const {key, value, term} = prefix.groups;
      
      if (!value && term.length >= 2) {
        valueCompletion(key, term);
      } else {
        hideElement(suggInline);
      }
    } else {
      keyCompletion(val.trim());
    }

    let currentTerm = '';
    let filter = {};

    if (prefix !== null) {
      const {key, value, term} = prefix.groups;

      if (!filters.contains(key) || !filters.includes(key, value)) return;

      currentTerm = term;
      filter[key] = value;
    } else {
      currentTerm = val.trim();
    }

    if (currentTerm.length < 3 || currentTerm.length > 32) return;

    if (currentTerm === lastTerm) {
      showElement(suggBox);
      return;
    }

    if (map.locationID) {
      filter.location = map.locationID;
    }

    if (count < qCount) return;

    qCount++;

    const data = {
      id: qCount,
      term: currentTerm,
      items: !filter.location,
      locations: !filter.item,
      features: !!filter.location
    };

    if (Object.keys(filter).length) data.filter = filter;

    socket.send(JSON.stringify(data));

    lastTerm = currentTerm;
  };

  const onInputFocusIn = async e => {
    try {
      await connect();
    } catch (err) {
      if (err !== errErrorClosure) console.error(err);
      return;
    }

    if (e.target.dataset.valid === 'true' && suggBox.querySelector('ul > li')) showElement(suggBox);
  };

  const onInputKeydown = e => {
    const selFirstBoxSugg = () => {
      const next = suggBox.querySelector('ul > li:first-child > a');
      if (next) next.focus();
    };
    const selLastBoxSugg = () => {
      const next = suggBox.querySelector('ul > li:last-child > a');
      if (next) next.focus();
    };
    const applyInlineSugg = () => {
      const {dataset, classList} = suggInline;
      if (classList.contains('show') && dataset.value) input.value = dataset.value;
    };

    switch (e.key) {
    case 'ArrowDown':
      e.preventDefault();
      selFirstBoxSugg();
      break;
    case 'ArrowUp':
      e.preventDefault();
      selLastBoxSugg();
      break;
    case 'Tab':
      e.preventDefault();
      applyInlineSugg();
      hideElement(suggInline);
      break;
    case 'Escape':
      hideElement(suggBox);
      break;
    }
  };

  input.addEventListener('keydown', onInputKeydown);
  input.addEventListener('focusin', onInputFocusIn);
  input.addEventListener('input', onInput);

  const onDocKeydown = e => {
    if (document.activeElement.nodeName === 'INPUT') return;
    if (e.ctrlKey && e.key === 'f' || e.key === 'F3' || e.key === '/') {
      e.preventDefault();
      input.focus();
    }
  };

  document.addEventListener('keydown', onDocKeydown);
};

initSearchSocket();


const initListFilter = () => {
  const el = document.getElementById('listFilter');
  if (el === null) return;

  const changeFilter = e => {
    const params = new URLSearchParams(window.location.search);
    const name = e.target.name, value = e.target.value;

    if (!value || value === 'all') {
      params.delete(name);
    } else {
      params.set(name, value);
    }

    const pageKey = 'p';
    if (params.has(pageKey)) params.delete(pageKey);

    window.location.search = params.toString();
  };

  const ul = el.querySelectorAll('ul > li');
  for (const li of ul) {
    const sel = li.getElementsByTagName('select')[0];
    sel.addEventListener('change', changeFilter);
  }
};

initListFilter();
