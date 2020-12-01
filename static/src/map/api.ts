import { Feature, Geometry, FeatureCollection, GeoJsonProperties } from "geojson";

export interface CustomFeature<G extends Geometry | null = Geometry, P = GeoJsonProperties> extends Feature<G, P> {
    _id: string
    name: string
    group: string
    _location: string
}

export interface FeatureGroup {
    _id: string
    name: string
    description: string
    tags: string[]
    _location: string
}

export interface FeatureGroupResult {
    total: number
    items: FeatureGroup[]
}

// TODO: Add options like limit and pagination
export class LocationAPI {
    private _location: string;
    private _url: URL;
    private _options: RequestInit;

    constructor(locationId: string) {
        this._location = locationId;
        this._url = new URL(`/location/${this._location}`, window.location.href);
        this._options = {};
    }

    private async _request(url: URL, opts: RequestInit) {
        const req = new Request(url.toString());

        try {
            const res = await fetch(req, opts);
            const json = await res.json();
            if (!res.ok) return Promise.reject(new Error(`${json.code}: ${json.message}`));
            return json;
        } catch (err) {
            return Promise.reject(err);
        }
    }

    private _json(url: URL) {
        return this._request(url, {
            ...this._options,
            headers: {
                ...this._options.headers,
                'Content-Type': 'application/json'
            }
        });
    }

    private _geojson(url: URL) {
        return this._request(url, {
            ...this._options,
            headers: {
                ...this._options.headers,
                'Content-Type': 'application/geo+json'
            }
        });
    }

    async featureGroups(): Promise<FeatureGroupResult> {
        const url = this._url;
        url.pathname += '/featuregroup';
        return await this._json(url);
    }

    async featureGroupByID(id: string): Promise<FeatureGroup> {
        const url = this._url;
        url.pathname += `/featuregroup/${id}`;
        return await this._json(url);
    }

    async featureByID(id: string): Promise<CustomFeature> {
        const url = this._url;
        url.pathname += `/feature/${id}`;
        return await this._geojson(url);
    }

    async featuresByGroup(id: string): Promise<FeatureCollection> {
        const url = this._url;
        url.pathname += '/feature';
        url.search = new URLSearchParams({ group: id }).toString();
        return await this._geojson(url);
    }

    async featuresByText(keyword: string): Promise<FeatureCollection> {
        const url = this._url;
        url.pathname += '/feature';
        url.search = new URLSearchParams({ text: keyword }).toString();
        return await this._geojson(url);
    }
}
