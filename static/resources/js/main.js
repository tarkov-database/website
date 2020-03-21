let map = {};

const registerTabs = () => {
  const openTab = e => {
    const tabcontent = document.getElementsByClassName('tab');
    for (const el of tabcontent) {
      el.style.display = 'none';
    }

    const clName = 'active';
    const tablinks = document.getElementsByClassName('tab-btn');
    for (const el of tablinks) {
      el.classList.remove(clName);
    }

    const el = e.currentTarget;
    document.getElementById(el.dataset.tab).style.display = 'block';
    el.classList.add(clName);
  };

  const tablinks = document.getElementsByClassName('tab-btn');
  for (const el of tablinks) {
    el.addEventListener('click', openTab);
  }
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
  const sugg = form.querySelector('.suggestion');

  let socket;
  let idleTimeout;
  let noReconnect = false;
  let qCount = 0;

  let keyword = '';

  const wait = time => new Promise(resolve => setTimeout(() => resolve(), time));

  const regexMeta = new RegExp(/[.*+?^${}()|[\]\\]/, 'g');
  const regexQuoteMeta = str => str.replace(regexMeta, '\\$&');

  const regexValidQuery = new RegExp(/^[\w-".\-,#!?& ]{2,32}$/);

  const showSuggestions = () => sugg.classList.replace('hide', 'show');
  const hideSuggestions = () => sugg.classList.replace('show', 'hide');

  const errErrorClosure = new Error('Search socket was closed with an error');

  const openSocket = async msgListener => {
    const host = window.location.host;
    const path = 'search/ws';

    let proto = 'wss';
    if (window.location.protocol !== 'https:') {
      proto = 'ws';
      console.warn('Insecure WebSocket protocol is used');
    }

    let socket;
    try {
      socket = await new WebSocket(`${proto}://${host}/${path}`);
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

    if (data.id !== qCount) return;

    const ul = sugg.querySelector('ul');

    if (!data.items || data.items.length === 0) {
      hideSuggestions();
      return;
    }

    const keywords = regexQuoteMeta(keyword).replace(' ', '|');

    const newUl = document.createElement('ul');
    for (const item of data.items) {
      const a = document.createElement('a');
      const li = document.createElement('li');
      const text = document.createElement('span');
      const div = document.createElement('div');

      const highlightMatches = str => {
        const re = new RegExp(`(${keywords})`, 'gi');
        const matches = str.match(re);
        if (!matches) return str;
        let bold = str;
        matches.forEach(m => bold = bold.replace(m, `<b>${m}</b>`));
        return bold;
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
        a.href += window.location.pathname.endsWith('/map') ? '/map': '';
        div.className = 'icon location';
        break;
      case 2:
        a.addEventListener('click', async () => map.flyToFeature(await map.getFeature(item.id, item.parent)));
        a.href = `#feature=${item.id}`;
        div.className = 'icon feature';
        break;
      }

      div.innerHTML = '&nbsp;';

      a.addEventListener('keydown', e => {
        const current = e.target.parentNode;

        let next = {};
        switch (e.key) {
        case 'ArrowUp':
          e.preventDefault();
          next = current.previousElementSibling;
          if (!next) input.focus();
          break;
        case 'ArrowDown':
          e.preventDefault();
          next = current.nextElementSibling;
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

    sugg.replaceChild(newUl, ul);

    showSuggestions();
  };

  const onInput = async event => {
    const count = qCount;

    try {
      await connect();
    } catch (err) {
      if (err !== errErrorClosure) console.error(err);
      return;
    }

    const val = event.target.value.trim();
    if (val === keyword || !regexValidQuery.test(val)) return;
    if (val.length < 3) {
      hideSuggestions();
      return;
    }

    if (count < qCount) return;

    qCount++;

    let data = {
      id: qCount,
      keyword: val
    };

    if (map.locationID) {
      Object.assign(data, {
        location: map.locationID,
        locations: true,
        features: true
      });
    } else {
      Object.assign(data, {
        items: true,
        locations: true
      });
    }

    keyword = val;

    socket.send(JSON.stringify(data));
  };

  const onFocusIn = async () => {
    try {
      await connect();
    } catch (err) {
      if (err !== errErrorClosure) console.error(err);
      return;
    }

    sugg.hidden = false;
  };

  input.addEventListener('keydown', e => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      const next = sugg.querySelector('ul > li > a');
      if (next) next.focus();
    }
  });

  input.addEventListener('focusin', onFocusIn);
  input.addEventListener('input', onInput);
};

initSearchSocket();


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
