/* global mapboxgl */

const loadJSON = async addr => {
  const req = new Request(addr);

  try {
    const res = await fetch(req);
    if (!res.ok) throw new Error(res.status);
    return res.json();
  } catch (err) {
    return new Error(err);
  }
};

const addDragMarker = async map => {
  const el = document.getElementById('dragMarker');

  let lngLat = [0, 0];
  const hashParams = new URLSearchParams(window.location.hash.replace('#', ''));
  const markerKey = 'marker';

  if (hashParams.has(markerKey)) {
    let zoom = map.getZoom();
    hashParams
      .get(markerKey)
      .split('/')
      .forEach((v, i) => {
        const val = parseFloat(v);
        if (i === 2) {
          zoom = val;
        } else {
          lngLat[i] = val;
        }
      });
    map.flyTo({center: lngLat, zoom});
  }

  const marker = new mapboxgl.Marker({
    element: el,
    draggable: true
  })
    .setLngLat(lngLat)
    .addTo(map);

  const setMarkerURL = () => {
    const pos = marker.getLngLat();
    hashParams.set(markerKey, `${pos.lng}/${pos.lat}/${map.getZoom()}`);
    window.location.hash = hashParams.toString();
  };

  marker.on('dragend', setMarkerURL);

  map.on('dblclick', e => {
    marker.setLngLat(e.lngLat);
    setMarkerURL();
  });

  el.hidden = false;
};

const addLayer = async (layer, id, map) => {
  let data = {};
  try {
    data = await loadJSON(`/resources/layer/${id}_${layer}.geojson`);
  } catch (err) {
    throw new Error(err);
  }

  const layerName = data.name;

  map.addSource(layerName, {type: 'geojson', data});

  map.addLayer({
    'id': layerName,
    'type': 'circle',
    'source': layerName,
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

  const popup = new mapboxgl.Popup({
    closeButton: false,
    closeOnClick: false
  });

  map.on('mouseenter', layerName, e => {
    map.getCanvas().style.cursor = 'pointer';

    const coordinates = e.features[0].geometry.coordinates.slice();
    const content = `<center>${layerName}<br><b>${e.features[0].properties.name}</b></center>`;

    while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
      coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
    }

    popup
      .setLngLat(coordinates)
      .setHTML(content)
      .addTo(map);
  });

  map.on('mouseleave', layerName, () => {
    map.getCanvas().style.cursor = '';
    popup.remove();
  });
};

const initMap = async () => {
  const elID = 'map';
  const el = document.getElementById(elID);
  const locID = el.dataset.id;

  let style = {};
  try {
    style = await loadJSON(`/resources/style/${locID}.json`);
  } catch (err) {
    throw new Error(err);
  }

  mapboxgl.workerUrl = '/resources/js/mapbox/mapbox-gl-csp-worker.js';
  const map = await new mapboxgl.Map({
    container: elID,
    style,
    maxBounds: style.metadata['mapbox:maxBounds'],
    antialias: true,
    doubleClickZoom: false
  });

  map.on('load', async () => {
    addDragMarker(map);
    addLayer('exfil', locID, map);
  });
};

initMap();
