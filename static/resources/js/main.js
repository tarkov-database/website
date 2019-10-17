'use strict';

const openTab = (evt, tabName) => { // eslint-disable-line
  const tabcontent = document.getElementsByClassName('tab');
  for (const el of tabcontent) {
    el.style.display = 'none';
  }

  const clName = 'active';
  const tablinks = document.getElementsByClassName('tab-btn');
  for (const el of tablinks) {
    el.classList.remove(clName);
  }
  document.getElementById(tabName).style.display = 'block';
  evt.currentTarget.classList.add(clName);
};


const loadImage = async () => {
  const image = document.getElementById('image');
  if (image === null) return;

  const staticURL = `//static.${window.location.host}`;
  const imageID = image.dataset.id;
  const request = new Request(`${staticURL}/image/icon/1-1/${imageID}.png`);

  try {
    const response = await fetch(request);
    if(!response.ok) throw new Error(response.status);

    const objectURL = URL.createObjectURL(await response.blob());
    const img = new Image();
    img.src = objectURL;
    img.onload = () => {
      const imgWidth = img.naturalWeight;
      const imgHeight = img.naturalHeight;
      const boxWidth = image.boxWidth;
      const boxHeight = image.offsetHeight;

      if (imgHeight > boxHeight || imgWidth > boxWidth) image.style.backgroundSize = 'contain';
    };

    image.style.backgroundImage = `url("${objectURL}")`;
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
  const tables = document.querySelectorAll('.sort-table thead th');
  if (tables.len === 0) return;

  const getCellValue = (tr, idx) => tr.children[idx].innerText || tr.children[idx].textContent;

  const comparer = (idx, asc) => (a, b) => ((v1, v2) =>
    v1 !== '' && v2 !== '' && !isNaN(v1) && !isNaN(v2) ? v1 - v2 : v1.toString().localeCompare(v2)
  )(getCellValue(asc ? a : b, idx), getCellValue(asc ? b : a, idx));


  for (const th of tables) {
    if (th.getElementsByTagName('a').length !== 0) continue;

    th.addEventListener('click', (() => {
      const table = th.closest('table');
      const body = table.querySelector('tbody');

      Array.from(body.querySelectorAll('tr'))
        .sort(comparer(Array.from(th.parentNode.children).indexOf(th), this.asc = !this.asc))
        .forEach(tr => body.appendChild(tr));

      const clSorted = 'sorted-by';
      const clUp = 'up';
      const clDown = 'down';

      if (!th.classList.contains(clSorted)) {
        const last = table.querySelector(`.${clSorted}`);
        if (last) {
          last.classList.remove(clSorted);
          if (last.classList.contains(clUp)) {
            last.classList.remove(clUp);
            th.classList.add(clDown);
          } else {
            last.classList.remove(clDown);
            th.classList.add(clUp);
          }
        } else {
          th.classList.add(clUp);
        }
        th.classList.add(clSorted);
        return;
      }

      if (th.classList.contains(clUp)) {
        th.classList.replace(clUp, clDown);
      } else {
        th.classList.replace(clDown, clUp);
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

    let proto = '';
    if (window.location.protocol === 'https:') {
      proto = 'wss';
    } else {
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
        socket = await openSocket(onMessage);
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

    if (data.items.length === 0) {
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

      div.innerHTML = '&nbsp;';
      div.className = `icon ${item.category}`;
      a.appendChild(div);

      const highlightMatches = str => {
        const re = new RegExp(`(${keywords})`, 'gi');
        const matches = str.match(re);
        if (!matches) return str;
        let bold = str;
        matches.forEach(m => bold = bold.replace(m, `<b>${m}</b>`));
        return bold;
      };

      text.innerHTML = highlightMatches(item.name);
      a.append(text);

      a.title = item.name;
      a.href = `/item/${item.category}/${item.id}`;
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

    const data = {
      id: qCount,
      text: val
    };

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

  input.addEventListener('focusin', onFocusIn);
  input.addEventListener('input', onInput);
};

initSearchSocket();
