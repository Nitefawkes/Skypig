# Propagation API Documentation

Real-time HF propagation conditions and band analysis for Ham-Radio Cloud.

## Overview

The Propagation API provides current space weather data and HF band condition forecasts to help operators choose the best bands for DX contacts. Data includes:

- **Solar indices** (Solar Flux, Sunspot Number)
- **Geomagnetic indices** (K-index, A-index)
- **Band-by-band conditions** (80m through 10m)
- **Rule-based forecasting** (day/night propagation)

**Update Frequency**: Data is cached for 15 minutes and refreshed automatically.

---

## Endpoints

### 1. Get Current Conditions

Retrieve current space weather and propagation data.

**Endpoint**: `GET /api/v1/propagation`

**Response**: `200 OK`
```json
{
  "data": {
    "timestamp": "2025-11-23T14:30:00Z",
    "solar_flux": 125.5,
    "sunspot_number": 45,
    "a_index": 8,
    "k_index": 2,
    "x_ray_flux": "C1.2",
    "solar_wind": 380.5,
    "bz_component": -2.3,
    "proton_flux": 0.5,
    "electron_flux": 1.2,
    "geomagnetic_storm": "None",
    "radio_blackout": "None",
    "solar_radiation": "None",
    "updated_at": "2025-11-23T14:30:00Z",
    "source": "Mock Data"
  },
  "meta": {
    "cache_ttl": 900
  }
}
```

**Example**:
```bash
curl http://localhost:8080/api/v1/propagation
```

---

### 2. Get Propagation Forecast

Get current conditions plus band-by-band analysis and summary.

**Endpoint**: `GET /api/v1/propagation/forecast`

**Response**: `200 OK`
```json
{
  "data": {
    "current": {
      "timestamp": "2025-11-23T14:30:00Z",
      "solar_flux": 125.5,
      "sunspot_number": 45,
      "a_index": 8,
      "k_index": 2,
      "updated_at": "2025-11-23T14:30:00Z",
      "source": "Mock Data"
    },
    "band_conditions": [
      {
        "band": "80m",
        "day": "fair",
        "night": "good",
        "reasoning": "Low bands favor nighttime propagation"
      },
      {
        "band": "40m",
        "day": "fair",
        "night": "good",
        "reasoning": "Low bands favor nighttime propagation"
      },
      {
        "band": "20m",
        "day": "excellent",
        "night": "fair",
        "reasoning": "Daytime DX band, depends on solar flux"
      },
      {
        "band": "15m",
        "day": "good",
        "night": "poor",
        "reasoning": "High bands need solar flux >120 (current: 126)"
      }
    ],
    "summary": "Solar Flux: 126 (Moderate), K-Index: 2 (Quiet). Best daytime bands: [20m, 15m]. Overall HF conditions: good.",
    "last_updated": "2025-11-23T14:30:00Z",
    "next_update": "2025-11-23T14:45:00Z"
  }
}
```

**Example**:
```bash
curl http://localhost:8080/api/v1/propagation/forecast
```

---

## Data Fields

### Space Weather Indices

| Field | Description | Range | Interpretation |
|-------|-------------|-------|----------------|
| `solar_flux` | Solar Flux Index (SFI) at 10.7 cm | 50-300 | Higher = better HF conditions |
| `sunspot_number` | Daily sunspot count | 0-200+ | More sunspots = higher solar activity |
| `k_index` | Planetary K-index (3-hour) | 0-9 | **0-2**: Quiet, **3-4**: Unsettled, **5+**: Storm |
| `a_index` | Planetary A-index (24-hour) | 0-400 | Lower = better propagation |
| `x_ray_flux` | Solar X-ray flux class | e.g., C1.2 | M/X class = solar flares |
| `solar_wind` | Solar wind speed (km/s) | 250-900 | Faster = potential disturbance |
| `bz_component` | IMF Bz component (nT) | -30 to +30 | Negative = geomag disturbance |

### Geomagnetic Storm Levels
- `None` - No geomagnetic activity
- `Minor` (G1) - K-index 5
- `Moderate` (G2) - K-index 6
- `Strong` (G3) - K-index 7
- `Severe` (G4-G5) - K-index 8-9

### Band Condition Ratings
- `excellent` - Optimal propagation expected
- `good` - Above-average conditions
- `fair` - Average/marginal conditions
- `poor` - Below-average or unlikely

---

## Band Analysis

### Rule-Based Forecasting

The system analyzes each HF band based on:

1. **Solar Activity** (SFI, sunspot number)
2. **Geomagnetic Conditions** (K-index, A-index)
3. **Time of Day** (UTC, day/night transition)
4. **Band Characteristics** (frequency-dependent behavior)

### Band Categories

#### Low Bands (80m, 40m)
- **Best**: Nighttime, low K-index
- **Characteristics**: Short-to-medium distance, reliable
- **Solar Dependency**: Low

#### Mid Bands (30m, 20m, 17m)
- **Best**: Daytime with moderate-to-high SFI
- **Characteristics**: Long-distance DX, sensitive to solar activity
- **Solar Dependency**: Moderate to High

#### High Bands (15m, 12m, 10m)
- **Best**: Daytime with high SFI (>120)
- **Characteristics**: Long-distance DX, peak solar cycle
- **Solar Dependency**: Very High

---

## Frontend Integration

### PropagationWidget Component

A pre-built Svelte component displays propagation data:

**Location**: `frontend/src/lib/components/PropagationWidget.svelte`

**Features**:
- Current solar indices (SFI, K-index, A-index, sunspots)
- Visual band condition chart (80m-10m)
- Color-coded ratings (green/blue/yellow/red)
- Auto-refresh every 15 minutes
- Text summary

**Usage**:
```svelte
<script>
  import PropagationWidget from '$components/PropagationWidget.svelte';
</script>

<PropagationWidget />
```

---

## Understanding the Data

### Solar Flux Index (SFI)

The SFI is measured at 10.7 cm wavelength and indicates solar activity:

- **< 70**: Very Low (poor high-band conditions)
- **70-100**: Low (fair mid-band, good low-band)
- **100-150**: Moderate (good all-around)
- **150-200**: High (excellent high-band)
- **> 200**: Very High (exceptional DX conditions)

### K-Index

The K-index measures geomagnetic activity over 3 hours:

- **0-1**: Very Quiet
- **2**: Quiet ✅ (best for HF)
- **3-4**: Unsettled ⚠️
- **5-6**: Active/Storm 🔴
- **7-9**: Major Storm 🔴🔴

**Lower is better for HF propagation!**

### How to Use This Data

**For DXing**:
- Look for: High SFI (>120), Low K-index (≤3)
- Best bands: Check "excellent" or "good" ratings
- Timing: Daytime for high/mid bands, night for low bands

**For Contesting**:
- Monitor K-index closely (≤4 ideal)
- Plan band changes based on forecast
- Watch for sudden geomagnetic storms

**For Local QSOs**:
- Low bands (80m/40m) work in most conditions
- Less affected by solar activity

---

## Current Limitations (MVP)

### Mock Data
- Currently using realistic mock data
- Varies by time of day for testing
- External API integration pending

### Planned Enhancements
- [ ] Real-time data from NOAA/HamQSL APIs
- [ ] Historical data and trends
- [ ] Aurora alerts
- [ ] Solar flare warnings
- [ ] Band opening predictions (ML-based)
- [ ] Per-location analysis (grid square)
- [ ] Email/SMS alerts for favorable conditions

---

## External Data Sources (Future)

### Primary Sources
- **NOAA SWPC**: Official space weather data
- **HamQSL**: Ham radio-specific propagation
- **WWV/WWVH**: Audio propagation forecasts

### API Integration Roadmap
1. **Phase 1 (MVP)**: Mock data for testing ✅
2. **Phase 2**: HamQSL API integration
3. **Phase 3**: NOAA SWPC integration
4. **Phase 4**: Historical database & trends

---

## Error Handling

### API Errors

**Failed to Fetch**:
```json
{
  "error": {
    "code": "FETCH_FAILED",
    "message": "Failed to fetch propagation data",
    "details": "Network timeout"
  }
}
```

**Frontend Fallback**:
- Shows cached data (if available)
- Displays error message
- Retries on next interval (15 min)

---

## Cache Behavior

- **TTL**: 15 minutes
- **Strategy**: Stale-while-revalidate
- **Benefits**: Reduces API calls, faster responses
- **Refresh**: Automatic background updates

---

## Testing

### Manual Testing

```bash
# Get current conditions
curl http://localhost:8080/api/v1/propagation | jq

# Get full forecast
curl http://localhost:8080/api/v1/propagation/forecast | jq

# Check cache behavior (should be instant on second call)
time curl http://localhost:8080/api/v1/propagation
time curl http://localhost:8080/api/v1/propagation
```

### Frontend Testing

1. Open http://localhost:5173
2. Scroll to "Propagation Conditions" widget
3. Verify:
   - Solar data displays correctly
   - Band conditions show color-coded ratings
   - Summary text is readable
   - Auto-refresh works (check timestamp)

---

## Implementation Details

### Models
- **Location**: `backend/internal/models/propagation.go`
- **Structures**: `PropagationData`, `BandConditions`, `PropagationForecast`

### Service
- **Location**: `backend/internal/services/propagation_service.go`
- **Features**: Caching, mock data, band analysis

### Handler
- **Location**: `backend/internal/handlers/propagation_handler.go`
- **Endpoints**: Current conditions, forecast

### Frontend
- **Component**: `frontend/src/lib/components/PropagationWidget.svelte`
- **Integration**: Homepage (`frontend/src/routes/+page.svelte`)

---

## Future: Machine Learning Forecasts

**Post-MVP Enhancement**:

Train LSTM model to predict:
- Band openings 1-24 hours ahead
- Peak propagation times per band
- DX window predictions by grid square

**Data Sources**:
- Historical solar indices
- QSO success rates (from user logs)
- Ionospheric sounding data

---

## References

- [NOAA Space Weather Prediction Center](https://www.swpc.noaa.gov/)
- [HamQSL Solar Terrestrial Data](https://www.hamqsl.com/solar.html)
- [ARRL Propagation Resources](http://www.arrl.org/propagation)
- [Understanding K-index](https://www.spaceweatherlive.com/en/help/the-k-index.html)

---

**Next Steps**: See [LOTW.md](./LOTW.md) (future) for integrating propagation data with LoTW sync strategy.
