/* global mapboxgl */

mapboxgl.workerUrl = '/resources/js/mapbox/mapbox-gl-csp-worker.js';

let map = {};

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

    const coordinates = e.features[0].geometry.coordinates.slice();
    const content = `<center>${name}<br><b>${e.features[0].properties.title}</b></center>`;

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
};

const flyToPopup = new mapboxgl.Popup({
  closeButton: false,
  closeOnClick: true,
  closeOnMove: true
});

export const flyToFeature = async feature => {
  const layerID = feature.group;
  const layer = map.getLayer(layerID);
  if (layer && map.getLayoutProperty(layerID, 'visibility') === 'none') {
    toggleLayerVisibility(layerID);
  }

  const lngLat = feature.geometry.coordinates.slice();
  const zoom = 18.5;

  map.flyTo({center: lngLat, zoom});

  const content = `<center><b>${feature.properties.title}</b></center>`;
  map.once('moveend', layerID, () => {
    flyToPopup
      .setLngLat(lngLat)
      .setHTML(content)
      .addTo(map);
  });

  const removePopup = () => {
    flyToPopup.remove();
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

const toggleLayerVisibility = id => {
  const el = document.querySelector(`#layers > .list input[name="${id}"]`);
  const cl = 'active';
  if (map.getLayoutProperty(id, 'visibility') !== 'visible') {
    map.setLayoutProperty(id, 'visibility', 'visible');
    el.checked = true;
    el.parentElement.classList.add(cl);
  } else {
    map.setLayoutProperty(id, 'visibility', 'none');
    el.checked = false;
    el.parentElement.classList.remove(cl);
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
  addLayer(layerName, {
    'id': id,
    'type': 'circle',
    'source': id,
    'layout': {
      'visibility': 'none'
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

  toggleLayerVisibility(id);
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

  const layers = document.getElementById('layers');
  const ul = document.createElement('ul');
  ul.classList.add('list');

  for (const g of data.items) {
    const li = document.createElement('li');

    const label = document.createElement('label');

    const icon = document.createElement('div');
    icon.classList.add('icon');

    label.appendChild(icon);

    const input = document.createElement('input');
    input.value = g._id;
    input.name = g._id;
    input.type = 'checkbox';

    input.addEventListener('change', e => toggleLayerVisibility(e.target.value));

    const span = document.createElement('span');
    span.innerText = g.name;

    label.appendChild(span);
    label.appendChild(input);
    li.appendChild(label);
    ul.appendChild(li);
  }

  layers.appendChild(ul);

  return data.items;
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
  });
};
