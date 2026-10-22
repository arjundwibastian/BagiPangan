import { useCallback, useEffect, useRef, useState } from "react";
import {
  Circle,
  MapContainer,
  Marker,
  TileLayer,
  useMap,
  useMapEvents,
} from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";

type LeafletMap = ReturnType<typeof L.map>;

const pinIcon = L.divIcon({
  className: "map-pin",
  html: '<span class="map-pin__dot"></span>',
  iconSize: [22, 22],
  iconAnchor: [11, 11],
});

interface MapPickerProps {
  latitude: number;
  longitude: number;
  onChange: (latitude: number, longitude: number) => void;
  radiusKm?: number;
}

function ClickToMove({
  onChange,
}: {
  onChange: (latitude: number, longitude: number) => void;
}) {
  useMapEvents({
    click(event) {
      onChange(event.latlng.lat, event.latlng.lng);
    },
  });
  return null;
}

function CaptureMap({ onReady }: { onReady: (map: LeafletMap) => void }) {
  const map = useMap();

  useEffect(() => {
    onReady(map);
  }, [map, onReady]);

  return null;
}

export function MapPicker({
  latitude,
  longitude,
  onChange,
  radiusKm,
}: MapPickerProps) {
  const mapRef = useRef<LeafletMap | null>(null);
  const [locating, setLocating] = useState(false);

  const handleReady = useCallback((map: LeafletMap) => {
    mapRef.current = map;
  }, []);

  function handleLocate() {
    if (!navigator.geolocation) return;
    setLocating(true);

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const lat = position.coords.latitude;
        const lng = position.coords.longitude;
        onChange(lat, lng);
        mapRef.current?.flyTo([lat, lng], 15);
        setLocating(false);
      },
      () => setLocating(false),
      { enableHighAccuracy: true, timeout: 8000 },
    );
  }

  const position: [number, number] = [latitude, longitude];

  return (
    <div className="map-picker">
      <div className="map-picker__toolbar">
        <span className="muted">Click the map or drag the pin.</span>
        <button
          type="button"
          className="btn-secondary"
          onClick={handleLocate}
          disabled={locating}
        >
          {locating ? "Locating..." : "Use my location"}
        </button>
      </div>

      <div className="map-picker__canvas">
        <MapContainer
          center={position}
          zoom={13}
          scrollWheelZoom
          style={{ height: "100%", width: "100%" }}
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <Marker
            position={position}
            icon={pinIcon}
            draggable
            eventHandlers={{
              dragend(event) {
                const point = event.target.getLatLng();
                onChange(point.lat, point.lng);
              },
            }}
          />
          {typeof radiusKm === "number" && radiusKm > 0 && (
            <Circle
              center={position}
              radius={radiusKm * 1000}
              pathOptions={{ color: "#1f4d2e", weight: 2, fillOpacity: 0.12 }}
            />
          )}
          <ClickToMove onChange={onChange} />
          <CaptureMap onReady={handleReady} />
        </MapContainer>
      </div>

      <p className="muted">
        Coordinates: {latitude.toFixed(6)}, {longitude.toFixed(6)}
      </p>
    </div>
  );
}
