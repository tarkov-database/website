/* global mapboxgl */

mapboxgl.workerUrl = '/resources/js/mapbox/mapbox-gl-csp-worker.js';

let map = {};
let loadedLayers = {};

export let locationID = '';

class APIRequest {
  constructor(addr, opts) {
    this.address = addr;
    this.options = opts || {};
  }

  async _request() {
    const req = new Request(this.address);

    try {
      const res = await fetch(req, this.options);
      const json = await res.json();
      if (!res.ok) return Promise.reject(new Error(`${json.code}: ${json.message}`));
      return json;
    } catch (err) {
      return Promise.reject(err);
    }
  }

  json() {
    this.options.headers = {'Content-Type': 'application/json'};
    return this._request();
  }

  geojson() {
    this.options.headers = {'Content-Type': 'application/geo+json'};
    return this._request();
  }
}

// const addDragMarker = async () => {
//   const el = document.getElementById('dragMarker');
//
//   let lngLat = [0, 0];
//   const hashParams = new URLSearchParams(window.location.hash.replace('#', ''));
//   const markerKey = 'marker';
//
//   if (hashParams.has(markerKey)) {
//     let zoom = map.getZoom();
//     hashParams
//       .get(markerKey)
//       .split('/')
//       .forEach((v, i) => {
//         const val = parseFloat(v);
//         if (i === 2) {
//           zoom = val;
//         } else {
//           lngLat[i] = val;
//         }
//       });
//     map.flyTo({center: lngLat, zoom});
//   }
//
//   const marker = new mapboxgl.Marker({
//     element: el,
//     draggable: true
//   })
//     .setLngLat(lngLat)
//     .addTo(map);
//
//   const setMarkerURL = () => {
//     const pos = marker.getLngLat();
//     hashParams.set(markerKey, `${pos.lng}/${pos.lat}/${map.getZoom()}`);
//     window.location.hash = hashParams.toString();
//   };
//
//   marker.on('dragend', setMarkerURL);
//
//   map.on('dblclick', e => {
//     marker.setLngLat(e.lngLat);
//     setMarkerURL();
//   });
//
//   el.hidden = false;
// };

const addLayer = (name, layer) => {
  const id = layer['id'];

  map.addLayer(layer);

  const popup = new mapboxgl.Popup({
    closeButton: false,
    closeOnClick: false
  });

  map.on('mouseenter', id, e => {
    map.getCanvas().style.cursor = 'pointer';

    const feature = e.features[0];
    const coordinates = feature.geometry.coordinates.slice();
    const content = `<center>${name}<br><b>${feature.properties.title || feature.name}</b></center>`;

    while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
      coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
    }

    popup
      .setLngLat(coordinates)
      .setHTML(content)
      .addTo(map);
  });

  map.on('mouseleave', id, () => {
    map.getCanvas().style.cursor = '';
    popup.remove();
  });

  loadedLayers[id] = {
    id,
    get: map.getLayer(id),
    get visible() {
      return map.getLayoutProperty(id, 'visibility') === 'visible';
    },
    toggleVisibility: function() {
      if (!this.visible) {
        map.setLayoutProperty(id, 'visibility', 'visible');
        return true;
      } else {
        map.setLayoutProperty(id, 'visibility', 'none');
        return false;
      }
    }
  };

  return loadedLayers[id];
};

const flyToCenter = () => {
  const center = map.getCenter();
  const zoom = 16.8;
  map.flyTo({center, zoom});
};

const flyToFeaturePopup = new mapboxgl.Popup({
  closeButton: false,
  closeOnClick: true,
  closeOnMove: true
});

export const flyToFeature = async feature => {
  const layerID = feature.group;
  const layer = loadedLayers[layerID];

  if (layer && !layer.visible) {
    const sl = loadedLayers['search'];
    if (sl && sl.visible) {
      toggleSearchLayer();
    } else {
      layer.toggleVisibility();
      setLayerUI(layer);
    }
  }

  const lngLat = feature.geometry.coordinates.slice();
  const zoom = 18.5;

  map.flyTo({center: lngLat, zoom});

  const content = `<center><b>${feature.properties.title || feature.name}</b></center>`;
  map.once('moveend', layerID, () => {
    flyToFeaturePopup
      .setLngLat(lngLat)
      .setHTML(content)
      .addTo(map);
  });

  const removePopup = () => {
    flyToFeaturePopup.remove();
    map.off('mouseenter', layerID, removePopup);
  };
  map.on('mouseenter', layerID, removePopup);
};

const featureFromURL = () => {
  const hashParams = new URLSearchParams(window.location.hash.replace('#', ''));
  const key = 'feature';

  if (hashParams.has(key)) {
    return hashParams.get(key);
  }

  return null;
};

export const getFeature = async (fID, gID) => {
  if (gID) {
    const source = map.getSource(gID);
    if (source) {
      for (const feature of source._data.features) {
        if (feature._id === fID) return feature;
      }
    }
  }

  const lID = locationID;
  let data = {};
  try {
    const url = new URL(`/location/${lID}/feature/${fID}`, window.location.href);
    const req = new APIRequest(url);
    data = await req.geojson();
  } catch (err) {
    return Promise.reject(err);
  }

  return data;
};

const getFeaturesByText = async keyword => {
  const lID = locationID;

  let data = {};
  try {
    const url = new URL(`/location/${lID}/feature`, window.location.href);
    url.search = new URLSearchParams({text: keyword});
    const req = new APIRequest(url);
    data = await req.geojson();
  } catch (err) {
    return Promise.reject(err);
  }

  const id = 'search';

  let search = loadedLayers[id];
  if (search) {
    map.getSource(id).setData(data);
  } else {
    map.addSource(id, {type: 'geojson', data});
    search = addLayer('Search result', {
      'id': id,
      'type': 'circle',
      'source': id,
      'layout': {
        'visibility': 'none'
      },
      'paint': {
        'circle-radius': {
          'base': 3,
          'stops': [
            [0, 7],
            [22, 18]
          ]
        },
        'circle-color': 'rgb(118, 43, 170)'
      }
    });
  }

  if (!search.visible) toggleSearchLayer();
};

let layerWasVisible = {};

const toggleSearchLayer = () => {
  const search = loadedLayers['search'];
  if (!search) return;

  if (!search.visible) {
    for (const id in loadedLayers) {
      const l = loadedLayers[id];
      if (id !== search.id && l.visible) {
        l.toggleVisibility();
        layerWasVisible[id] = true;
      }
    }
  } else {
    for (const id in loadedLayers) {
      const l = loadedLayers[id];
      if (layerWasVisible[id]) {
        l.toggleVisibility();
      }
    }
  }

  search.toggleVisibility();
};

const setLayerUI = layer => {
  const el = document.querySelector(`#layers input[value="${layer.id}"]`);
  if (layer.visible) {
    el.parentElement.classList.add('active');
    el.checked = true;
  } else {
    el.parentElement.classList.remove('active');
    el.checked = false;
  }
};

const addGroupLayer = async group => {
  const lID = locationID;
  const id = group._id;
  const layerName = group.name;

  let data = {};
  try {
    const url = new URL(`/location/${lID}/feature`, window.location.href);
    url.search = new URLSearchParams({group: id});
    const req = new APIRequest(url);
    data = await req.geojson();
  } catch (err) {
    return Promise.reject(err);
  }

  map.addSource(id, {type: 'geojson', data});
  const layer = addLayer(layerName, {
    'id': id,
    'type': 'circle',
    'source': id,
    'layout': {
      'visibility': 'visible'
    },
    'paint': {
      'circle-radius': {
        'base': 2,
        'stops': [
          [0, 7],
          [22, 18]
        ]
      },
      'circle-color': 'rgb(43, 159, 26)'
    }
  });

  const layers = document.getElementById('layers');
  const ul = document.createElement('ul');
  ul.classList.add('list');

  const li = document.createElement('li');

  const label = document.createElement('label');
  label.classList.add('active');

  const icon = document.createElement('div');
  icon.classList.add('icon');

  label.appendChild(icon);

  const input = document.createElement('input');
  input.value = id;
  input.name = id;
  input.type = 'checkbox';
  input.checked = true;

  input.addEventListener('change', e => {
    const sl = loadedLayers['search'];
    if (sl && sl.visible) return;

    const visible = layer.toggleVisibility();
    const el = e.target;

    if (!visible) {
      el.checked = false;
      el.parentElement.classList.remove('active');
    } else {
      el.checked = true;
      el.parentElement.classList.add('active');
    }
  });

  const span = document.createElement('span');
  span.innerText = layerName;

  label.appendChild(span);
  label.appendChild(input);
  li.appendChild(label);
  ul.appendChild(li);
  layers.appendChild(ul);
};

const getGroups = async () => {
  const lID = locationID;

  let data = {};
  try {
    const url = new URL(`/location/${lID}/featuregroup`, window.location.href);
    const req = new APIRequest(url);
    data = await req.json();
  } catch (err) {
    return Promise.reject(err);
  }

  return data.items;
};

const registerSearch = () => {
  const form = document.getElementById('search');
  const input = form.querySelector('input[type=search]');

  form.addEventListener('submit', e => {
    e.preventDefault();
    if (e.isTrusted && input.value.length >= 3) {
      getFeaturesByText(input.value);
      flyToCenter();
    }
  });

  input.addEventListener('input', e => {
    const el = e.target;
    if (el.value.length === 0) toggleSearchLayer();
  });
};

export const init = async el => {
  const lID = el.dataset.id;
  locationID = lID;

  let style = {};
  try {
    const url = new URL(`/resources/style/${lID}.json`, window.location.href);
    const req = new APIRequest(url);
    style = await req.json();
  } catch (err) {
    return Promise.reject(err);
  }

  map = await new mapboxgl.Map({
    container: el.id,
    style,
    maxBounds: style.metadata['mapbox:maxBounds'],
    antialias: true,
    doubleClickZoom: false
  });

  map.on('load', async () => {
    // addDragMarker(map);
    for (const g of await getGroups()) {
      addGroupLayer(g);
    }
    const fID = featureFromURL();
    if (fID) flyToFeature(await getFeature(fID));
    registerSearch();
  });
};
