const filterStatus = {
  'category': '',
  'type': '',
  'class': '',
  'caliber': '',
  'available': ''
};

const setFilter = (key, value) => {
  filterStatus[key] = value;
  objectToHash(filterStatus);
  updateFilter();
};

const updateFilter = () => {
  const elements = document.getElementById('list').querySelectorAll('li');
  for (const el of elements) {
    let i = 0;
    const dataset = Object.entries(el.dataset);
    for (const [key, value] of dataset) {
      if (filterStatus[key] !== '' && filterStatus[key] !== value) {
        if (!el.hidden) el.hidden = true;
        break;
      }
      i++;
      if (i === dataset.length && el.hidden) el.hidden = false;
    }
  }
};

const hashToObject = () => {
  const hash = window.location.hash.slice(1);
  const filter = {};
  for (const pair of hash.split('&')) {
    const arr = pair.split('=');
    if (arr[1]) filter[arr[0]] = arr[1];
  }
  return filter;
};

const objectToHash = obj => {
  const pairs = [];
  Object.entries(obj).forEach(([k, v]) => v !== null && v !== '' ? pairs.push(`${k}=${v}`) : false);
  window.location.hash = pairs.join('&');
};

// const updateFilter = name => {
//   const filters = document.getElementById('filters').querySelectorAll('li > select');
//   const list = document.getElementById('list');
//   for (const filter of filters) {
//     // if (filter.name === name) continue;
//     for (const opt of filter.childNodes) {
//       if (opt.value === '') continue;
//       let i = 0;
//       for (const el of list.querySelectorAll(`li[data-${filter.name}="${opt.value}"]`)) {
//         if (el.style.display !== 'none') i++;
//       }
//       opt.innerHTML = opt.innerHTML.replace(/\(\d+\)/, `(${i})`);
//       // if (i > 0) {
//       //   opt.disabled = false;
//       // } else {
//       //   opt.disabled = true;
//       // }
//     }
//   }
// };

const camelToTitle = camelCase => camelCase.length <= 3 ? camelCase.toUpperCase() : camelCase
  .replace(/([A-Z])/g, match => ` ${match}`)
  .replace(/^./, match => match.toUpperCase());

const initFilters = () => {
  const list = document.getElementById('list');
  if (list === null) return;

  const filters = {
    'category': [],
    'type': [],
    'class': [],
    'caliber': [],
    'available': []
  };

  const elements = list.querySelectorAll('li');
  for (const el of elements) {
    for (let [key, value] of Object.entries(el.dataset)) {
      if (filters[key]) filters[key].push(value);
    }
  }

  for (const prop in filters) {
    if (filters.hasOwnProperty(prop)) { // eslint-disable-line no-prototype-builtins
      const arr = filters[prop].sort();
      filters[prop] = {};

      for (const key of arr) {
        if (filters[prop]) {
          filters[prop][key] = filters[prop][key] + 1 || 1;
        }
      }
    }
  }

  const filterNode = document.getElementById('filters');

  for (const prop in filters) {
    if (!filters.hasOwnProperty(prop) || Object.keys(filters[prop]).length < 2) continue; // eslint-disable-line no-prototype-builtins
    const el = document.createElement('li');
    const sel = document.createElement('select');
    sel.name = prop;

    const defOpt = document.createElement('option');
    defOpt.value = '';
    defOpt.innerHTML = 'All';
    sel.appendChild(defOpt);

    sel.addEventListener('change', e => { setFilter(prop, e.target.value); });

    for (const [name] of Object.entries(filters[prop])) {
      const opt = document.createElement('option');

      let title = name;
      if (prop === 'category' || prop === 'class' || prop === 'type') title = camelToTitle(title);
      if (title === 'true') title = 'Yes';
      if (title === 'false') title = 'No';
      opt.innerHTML = title;
      opt.value = name;

      sel.appendChild(opt);
    }

    const label = document.createElement('label');
    label.innerHTML = camelToTitle(prop);
    el.appendChild(label);

    el.appendChild(sel);
    filterNode.appendChild(el);
  }

  if (filterNode.childNodes.length > 1) {
    document.querySelector('.filters').hidden = false;
  }

  const hash = Object.entries(hashToObject());
  for (const [k, v] of hash) {
    if (!filters.hasOwnProperty(k)) continue; // eslint-disable-line no-prototype-builtins
    filterStatus[k] = v;
    document.querySelector(`#filters > li select[name="${k}"]`).value = v;
  }
  updateFilter();
};

initFilters();
