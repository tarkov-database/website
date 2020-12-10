import {
    AnyLayer,
    GeoJSONSource,
    LngLat,
    Map as MapboxMap,
    MapMouseEvent,
    Popup,
    Style,
} from "mapbox-gl";
import {
    FeatureCollection,
    Feature,
    Point,
    LineString,
    Geometry,
} from "geojson";
import measureLength from "@turf/length";

import { LocationAPI, CustomFeature, FeatureGroup } from "./api";

const getCSSVariable = (v: string) =>
    getComputedStyle(document.documentElement).getPropertyValue(v);

const layerColors: { [key: string]: string } = {
    search: getCSSVariable("--layer-search-color"),
    exfil: getCSSVariable("--layer-exfil-color"),
    cache: getCSSVariable("--layer-cache-color"),
};

let map: MapboxMap;
const loadedLayers: Map<string, ActiveLayer> = new Map();

export let locationID = "";

class ActiveLayer {
    readonly id: string;

    constructor(id: string) {
        this.id = id;
    }

    get get() {
        return map.getLayer(this.id);
    }

    get visible() {
        return map.getLayoutProperty(this.id, "visibility") === "visible";
    }

    toggleVisibility() {
        if (!this.visible) {
            map.setLayoutProperty(this.id, "visibility", "visible");
            return true;
        } else {
            map.setLayoutProperty(this.id, "visibility", "none");
            return false;
        }
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

const enableMeasureLines = () => {
    const valueContainer = document.getElementById("distance");
    if (valueContainer === null) return;

    if (valueContainer.style.visibility !== "visible") {
        valueContainer.style.visibility = "visible";
    } else {
        return;
    }

    const measurement: FeatureCollection<Geometry> = {
        type: "FeatureCollection",
        features: [],
    };

    const sourceID = "measurement";
    map.addSource(sourceID, {
        type: "geojson",
        data: measurement,
    });

    const pointLayer = "measure-points";
    map.addLayer({
        id: pointLayer,
        type: "circle",
        source: sourceID,
        paint: {
            "circle-radius": 5,
            "circle-color": "rgb(157, 59, 255)",
        },
        filter: ["in", "$type", "Point"],
    });

    const lineLayer = "measure-lines";
    map.addLayer({
        id: lineLayer,
        type: "line",
        source: sourceID,
        layout: {
            "line-cap": "round",
            "line-join": "round",
        },
        paint: {
            "line-color": "rgb(157, 59, 255)",
            "line-width": 2.5,
        },
        filter: ["in", "$type", "LineString"],
    });

    const addPoint = ({ point, lngLat }: MapMouseEvent) => {
        const features = map.queryRenderedFeatures(point, {
            layers: [pointLayer],
        });

        if (measurement.features.length > 1) measurement.features.pop();

        if (features.length) {
            const id = features[0].properties?.id;
            measurement.features = measurement.features.filter(
                ({ properties }) => properties?.id !== id
            );
        } else {
            const point: Feature<Point> = {
                type: "Feature",
                geometry: {
                    type: "Point",
                    coordinates: [lngLat.lng, lngLat.lat],
                },
                properties: {
                    id: String(new Date().getTime()),
                },
            };

            measurement.features.push(point);
        }

        if (measurement.features.length > 1) {
            const linestring: Feature<LineString> = {
                type: "Feature",
                geometry: {
                    type: "LineString",
                    coordinates: (measurement.features.filter(
                        ({ geometry }) => geometry.type === "Point"
                    ) as Feature<Point>[]).map(
                        ({ geometry }) => geometry.coordinates
                    ),
                },
                properties: {},
            };
            measurement.features.push(linestring);

            const val = valueContainer.getElementsByTagName("span")[0];
            val.innerText = `${(
                measureLength(linestring) * 1000
            ).toLocaleString()}m`;
        }

        (map.getSource(sourceID) as GeoJSONSource).setData(measurement);
    };

    map.on("click", addPoint);

    const changeCursor = ({ point }: MapMouseEvent) => {
        const features = map.queryRenderedFeatures(point, {
            layers: [pointLayer],
        });
        map.getCanvas().style.cursor = features.length
            ? "pointer"
            : "crosshair";
    };

    map.on("mousemove", changeCursor);

    const removeMeasurement = () => {
        valueContainer.style.visibility = "hidden";
        valueContainer.getElementsByTagName("span")[0].innerText = "0.0m";
        map.off("click", addPoint);
        map.off("mousemove", changeCursor);
        map.getCanvas().style.cursor = "";
        map.removeLayer(pointLayer);
        map.removeLayer(lineLayer);
        map.removeSource(sourceID);
    };

    const mapEl = map.getContainer();
    const escKey = (e: KeyboardEvent) => {
        closeEl?.removeEventListener("click", closeClick);
        if (e.key === "Escape") removeMeasurement();
    };
    mapEl.addEventListener("keydown", escKey, { once: true });

    const closeEl = document.getElementById("closeMeasurement");
    if (closeEl === null) return;

    closeEl.style.cursor = "pointer";

    const closeClick = () => {
        mapEl.removeEventListener("keydown", escKey);
        removeMeasurement();
    };
    closeEl?.addEventListener("click", closeClick, { once: true });
};

const registerMenu = () => {
    const measure = document.getElementById("measure");
    if (measure === null) return;
    measure.addEventListener("click", () => enableMeasureLines());
    measure.style.cursor = "pointer";
};

const getRandomLayerColor = () => {
    const random = (min: number, max: number) => {
        (min = Math.ceil(min)), (max = Math.floor(max));
        return Math.floor(Math.random() * (max - min + 1)) + min;
    };

    const h = random(80, 320);
    const s = random(70, 90);
    const l = random(55, 65);

    return `hsl(${h}, ${s}%, ${l}%)`;
};

const addLayer = (name: string, layer: AnyLayer) => {
    const id = layer["id"];

    map.addLayer(layer);

    const popup = new Popup({
        closeButton: false,
        closeOnClick: false,
    });

    map.on("mouseenter", id, ({ features, lngLat }) => {
        map.getCanvas().style.cursor = "pointer";

        if (!features) return;

        const feature: CustomFeature = features[0] as any;

        if (feature.geometry.type === "Point") {
            const coords = feature.geometry.coordinates;
            while (Math.abs(lngLat.lng - coords[0]) > 180) {
                coords[0] += lngLat.lng > coords[0] ? 360 : -360;
            }
            lngLat = new LngLat(coords[0], coords[1]);
        }

        const content = `<center>${name}<br><b>${
            feature.properties?.title || feature.name
        }</b></center>`;

        popup.setLngLat(lngLat).setHTML(content).addTo(map);
    });

    map.on("mouseleave", id, () => {
        map.getCanvas().style.cursor = "";
        popup.remove();
    });

    const lr = new ActiveLayer(id);
    loadedLayers.set(id, lr);

    return lr;
};

const flyToCenter = () => {
    const center = map.getCenter();
    const zoom = 16.8;
    map.flyTo({ center, zoom });
};

const flyToFeaturePopup = new Popup({
    closeButton: false,
    closeOnClick: true,
    closeOnMove: true,
});

export const flyToFeature = async (feature: CustomFeature): Promise<void> => {
    const layerID = feature.group;
    const layer = loadedLayers.get(layerID);

    if (layer && !layer.visible) {
        const sl = loadedLayers.get("search");
        if (sl && sl.visible) {
            toggleSearchLayer();
        } else {
            layer.toggleVisibility();
            setLayerUI(layer);
        }
    }

    // Temporary workaround
    if (feature.geometry.type !== "Point") return;

    const coords = feature.geometry.coordinates.slice();
    const lngLat = new LngLat(coords[0], coords[1]);
    const zoom = 18.5;

    map.flyTo({ center: lngLat, zoom });

    const content = `<center><b>${
        feature.properties?.title || feature.name
    }</b></center>`;
    map.once("moveend", () => {
        flyToFeaturePopup.setLngLat(lngLat).setHTML(content).addTo(map);
    });

    const removePopup = () => {
        flyToFeaturePopup.remove();
        map.off("mouseenter", layerID, removePopup);
    };
    map.on("mouseenter", layerID, removePopup);
};

const featureFromURL = () => {
    const hashParams = new URLSearchParams(
        window.location.hash.replace("#", "")
    );
    const key = "feature";

    if (hashParams.has(key)) {
        return hashParams.get(key);
    }

    return null;
};

export const getFeature = async (fID: string): Promise<CustomFeature> => {
    // if (gID) {
    //     const source = map.getSource(gID) as GeoJSONSource;
    //     if (source) {
    //         for (const feature of source._data.features) {
    //             if (feature._id === fID) return feature;
    //         }
    //     }
    // }

    const location = new LocationAPI(locationID);
    const data = await location.featureByID(fID);

    return data;
};

const getFeaturesByText = async (keyword: string) => {
    const location = new LocationAPI(locationID);
    const data = await location.featuresByText(keyword);

    const id = "search";

    let search = loadedLayers.get(id);
    if (search) {
        (map.getSource(id) as GeoJSONSource).setData(data);
    } else {
        map.addSource(id, { type: "geojson", data });
        search = addLayer("search-result", {
            id: id,
            type: "circle",
            source: id,
            layout: {
                visibility: "none",
            },
            paint: {
                "circle-radius": {
                    base: 3,
                    stops: [
                        [0, 7],
                        [22, 18],
                    ],
                },
                "circle-color": layerColors[id],
            },
        });
    }

    if (!search?.visible) toggleSearchLayer();
};

const layerWasVisible: { [key: string]: boolean } = {};

const toggleSearchLayer = () => {
    const search = loadedLayers.get("search");
    if (!search) return;

    if (!search.visible) {
        for (const id in loadedLayers) {
            const l = loadedLayers.get(id);
            if (id !== search.id && l?.visible) {
                l.toggleVisibility();
                layerWasVisible[id] = true;
            }
        }
    } else {
        for (const id in loadedLayers) {
            const l = loadedLayers.get(id);
            if (layerWasVisible[id]) {
                l?.toggleVisibility();
            }
        }
    }

    search.toggleVisibility();
};

const setLayerUI = (layer: ActiveLayer) => {
    const el = document.querySelector<HTMLInputElement>(
        `#layers input[value="${layer.id}"]`
    );
    if (el === null) return;

    if (layer.visible) {
        el.parentElement?.classList.add("active");
        el.checked = true;
    } else {
        el.parentElement?.classList.remove("active");
        el.checked = false;
    }
};

const addGroupLayer = async (group: FeatureGroup) => {
    const id = group._id;
    const layerName = group.name;

    const location = new LocationAPI(locationID);
    const features = await location.featuresByGroup(id);

    const color = layerColors[group.tags[0]] || getRandomLayerColor();

    map.addSource(id, { type: "geojson", data: features });
    const layer = addLayer(layerName, {
        id: id,
        type: "circle",
        source: id,
        layout: {
            visibility: "visible",
        },
        paint: {
            "circle-radius": {
                base: 2,
                stops: [
                    [0, 7],
                    [22, 18],
                ],
            },
            "circle-color": color,
        },
    });

    const layers = document.getElementById("layers");
    const ul = layers?.getElementsByTagName("ul")[0];

    const li = document.createElement("li");

    const label = document.createElement("label");
    label.classList.add("active");

    const dot = document.createElement("span");
    dot.classList.add("dot");
    dot.style.backgroundColor = color;

    label.appendChild(dot);

    const input = document.createElement("input");
    input.value = id;
    input.name = id;
    input.type = "checkbox";
    input.checked = true;

    input.addEventListener("change", function (this: HTMLInputElement) {
        const sl = loadedLayers.get("search");
        if (sl?.visible) return;

        const visible = layer.toggleVisibility();

        if (!visible) {
            this.checked = false;
            this.parentElement?.classList.remove("active");
        } else {
            this.checked = true;
            this.parentElement?.classList.add("active");
        }
    });

    const span = document.createElement("span");
    span.innerText = layerName;

    label.appendChild(span);
    label.appendChild(input);
    li.appendChild(label);
    ul?.appendChild(li);
};

const getGroups = async () => {
    const location = new LocationAPI(locationID);
    const groups = await location.featureGroups();

    return groups.items;
};

const registerSearch = () => {
    const form = document.getElementById("search");
    const input = form?.querySelector<HTMLInputElement>("input[type=search]");
    if (!input) return;

    form?.addEventListener("submit", (e) => {
        e.preventDefault();
        if (e.isTrusted && input.value.length >= 3) {
            getFeaturesByText(input.value);
            flyToCenter();
        }
    });

    input.addEventListener("input", function (this: HTMLInputElement) {
        if (this.value.length === 0) toggleSearchLayer();
    });
};

export const init = async (el: HTMLElement): Promise<void> => {
    const lID = el.dataset.id || "";
    locationID = lID;

    let style: Style;
    try {
        const url = new URL(
            `/resources/style/${lID}.json`,
            window.location.href
        );
        const req = await fetch(url.toString());
        style = await req.json();
    } catch (err) {
        return Promise.reject(err);
    }

    map = new MapboxMap({
        container: el.id,
        style,
        maxBounds: style.metadata["mapbox:maxBounds"],
        antialias: true,
        doubleClickZoom: false,
    });

    map.on("load", async () => {
        for (const g of await getGroups()) addGroupLayer(g);
        const fID = featureFromURL();
        if (fID) flyToFeature(await getFeature(fID));
        registerSearch();
        registerMenu();
    });
};
