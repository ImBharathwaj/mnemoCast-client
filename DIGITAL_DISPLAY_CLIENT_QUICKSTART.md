# 🚀 Digital Display Client - Quick Start Guide

**Quick implementation guide for getting the digital display client up and running.**

---

## 🎯 Quick Start (30 minutes)

### Step 1: Initialize Project

```bash
# Create new React + TypeScript project
npm create vite@latest digital-display-client -- --template react-ts
cd digital-display-client

# Install dependencies
npm install axios
npm install video.js @videojs/themes
npm install zustand  # For state management
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

### Step 2: Project Structure

```
src/
├── services/
│   └── api/
│       ├── client.ts          # Axios instance
│       ├── screen.ts          # Screen registration
│       ├── playlist.ts        # Playlist fetching
│       └── events.ts          # Event tracking
├── components/
│   ├── Player.tsx             # Main player component
│   └── StatusBar.tsx          # Connection status
├── hooks/
│   ├── usePlaylist.ts         # Playlist management
│   └── usePlayer.ts           # Playback logic
├── types/
│   └── index.ts               # TypeScript types
├── config/
│   └── default.ts             # Default configuration
└── App.tsx
```

### Step 3: Core Implementation Files

#### `src/services/api/client.ts`
```typescript
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
apiClient.interceptors.request.use(
  (config) => {
    console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('[API Error]', error.response?.data || error.message);
    return Promise.reject(error);
  }
);
```

#### `src/services/api/screen.ts`
```typescript
import { apiClient } from './client';

export interface ScreenLocation {
  city: string;
  area: string;
  venueType: string;
}

export interface Screen {
  id: string;
  name: string;
  location: ScreenLocation;
  classification: number;
  isOnline: boolean;
  lastHeartbeat?: string;
}

export interface RegisterScreenRequest {
  id: string;
  name: string;
  location: ScreenLocation;
  classification: number;
}

export const screenApi = {
  register: async (data: RegisterScreenRequest): Promise<Screen> => {
    const response = await apiClient.post<Screen>('/screens/register', data);
    return response.data;
  },

  heartbeat: async (screenId: string): Promise<void> => {
    await apiClient.post(`/screens/${screenId}/heartbeat`);
  },
};
```

#### `src/services/api/playlist.ts`
```typescript
import { apiClient } from './client';

export interface PlaylistItem {
  creativeId: string;
  campaignId: string;
  url: string;
  durationSeconds: number;
  order: number;
  type: 'image' | 'video';
}

export interface Playlist {
  screenId: string;
  generatedAt: string;
  durationMinutes: number;
  items: PlaylistItem[];
}

export const playlistApi = {
  fetch: async (screenId: string, durationMinutes: number = 3): Promise<Playlist> => {
    const response = await apiClient.get<Playlist>(
      `/screens/${screenId}/playlist`,
      { params: { durationMinutes } }
    );
    return response.data;
  },
};
```

#### `src/services/api/events.ts`
```typescript
import { apiClient } from './client';

export interface PlayEvent {
  screenId: string;
  creativeId: string;
  campaignId: string;
  timestamp: string;
}

export const eventsApi = {
  recordImpression: async (event: PlayEvent): Promise<void> => {
    await apiClient.post('/events/impression', event);
  },

  recordPlay: async (event: PlayEvent): Promise<void> => {
    await apiClient.post('/events/play', event);
  },
};
```

#### `src/types/index.ts`
```typescript
export interface ScreenConfig {
  screenId: string;
  name: string;
  location: {
    city: string;
    area: string;
    venueType: string;
  };
  classification: number;
  backendUrl?: string;
  playlistRefreshInterval?: number; // minutes
  heartbeatInterval?: number; // seconds
}

export type { Playlist, PlaylistItem } from '../services/api/playlist';
export type { Screen } from '../services/api/screen';
```

#### `src/config/default.ts`
```typescript
import { ScreenConfig } from '../types';

export const defaultConfig: ScreenConfig = {
  screenId: 'screen-1',
  name: 'Default Screen',
  location: {
    city: 'Chennai',
    area: 'Airport',
    venueType: 'airport',
  },
  classification: 1,
  backendUrl: 'http://localhost:8080',
  playlistRefreshInterval: 3, // minutes
  heartbeatInterval: 30, // seconds
};

// Load config from environment or localStorage
export const loadConfig = (): ScreenConfig => {
  const envConfig = {
    screenId: import.meta.env.VITE_SCREEN_ID,
    backendUrl: import.meta.env.VITE_API_URL,
  };

  const storedConfig = localStorage.getItem('screenConfig');
  const parsedConfig = storedConfig ? JSON.parse(storedConfig) : {};

  return {
    ...defaultConfig,
    ...parsedConfig,
    ...envConfig,
  };
};
```

#### `src/hooks/usePlaylist.ts`
```typescript
import { useState, useEffect, useCallback } from 'react';
import { playlistApi, Playlist } from '../services/api/playlist';
import { loadConfig } from '../config/default';

export const usePlaylist = () => {
  const [playlist, setPlaylist] = useState<Playlist | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const config = loadConfig();

  const fetchPlaylist = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await playlistApi.fetch(config.screenId, 3);
      setPlaylist(data);
      console.log('[Playlist] Fetched:', data.items.length, 'items');
    } catch (err) {
      setError(err as Error);
      console.error('[Playlist] Fetch failed:', err);
    } finally {
      setLoading(false);
    }
  }, [config.screenId]);

  useEffect(() => {
    // Initial fetch
    fetchPlaylist();

    // Refresh every N minutes
    const interval = setInterval(
      fetchPlaylist,
      (config.playlistRefreshInterval || 3) * 60 * 1000
    );

    return () => clearInterval(interval);
  }, [fetchPlaylist, config.playlistRefreshInterval]);

  return { playlist, loading, error, refresh: fetchPlaylist };
};
```

#### `src/hooks/usePlayer.ts`
```typescript
import { useState, useEffect, useRef } from 'react';
import { PlaylistItem } from '../services/api/playlist';
import { eventsApi } from '../services/api/events';
import { loadConfig } from '../config/default';

export const usePlayer = (items: PlaylistItem[]) => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const config = loadConfig();

  const playItem = (item: PlaylistItem) => {
    // Send impression event
    eventsApi.recordImpression({
      screenId: config.screenId,
      creativeId: item.creativeId,
      campaignId: item.campaignId,
      timestamp: new Date().toISOString(),
    }).catch(console.error);

    // Set timer for duration
    timerRef.current = setTimeout(() => {
      // Send play event
      eventsApi.recordPlay({
        screenId: config.screenId,
        creativeId: item.creativeId,
        campaignId: item.campaignId,
        timestamp: new Date().toISOString(),
      }).catch(console.error);

      // Move to next item
      setCurrentIndex((prev) => (prev + 1) % items.length);
    }, item.durationSeconds * 1000);
  };

  useEffect(() => {
    if (items.length === 0) return;

    setIsPlaying(true);
    playItem(items[currentIndex]);

    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, [items, currentIndex]);

  return {
    currentItem: items[currentIndex] || null,
    isPlaying,
    currentIndex,
  };
};
```

#### `src/components/Player.tsx`
```typescript
import { useEffect, useRef } from 'react';
import { PlaylistItem } from '../services/api/playlist';

interface PlayerProps {
  item: PlaylistItem | null;
}

export const Player = ({ item }: PlayerProps) => {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (!item) return;

    if (item.type === 'video' && videoRef.current) {
      videoRef.current.src = item.url;
      videoRef.current.play().catch(console.error);
    }
  }, [item]);

  if (!item) {
    return (
      <div className="flex items-center justify-center h-screen bg-black text-white">
        <div>Loading...</div>
      </div>
    );
  }

  if (item.type === 'image') {
    return (
      <div className="flex items-center justify-center h-screen bg-black">
        <img
          src={item.url}
          alt={`Creative ${item.creativeId}`}
          className="max-w-full max-h-full object-contain"
        />
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center h-screen bg-black">
      <video
        ref={videoRef}
        src={item.url}
        className="max-w-full max-h-full"
        autoPlay
        muted
        playsInline
        onError={(e) => console.error('Video playback error:', e)}
      />
    </div>
  );
};
```

#### `src/components/StatusBar.tsx`
```typescript
import { useEffect, useState } from 'react';
import { screenApi } from '../services/api/screen';
import { loadConfig } from '../config/default';

export const StatusBar = () => {
  const [connected, setConnected] = useState(false);
  const config = loadConfig();

  useEffect(() => {
    // Send heartbeat every N seconds
    const interval = setInterval(async () => {
      try {
        await screenApi.heartbeat(config.screenId);
        setConnected(true);
      } catch (error) {
        setConnected(false);
        console.error('[Heartbeat] Failed:', error);
      }
    }, (config.heartbeatInterval || 30) * 1000);

    // Initial heartbeat
    screenApi.heartbeat(config.screenId).catch(console.error);

    return () => clearInterval(interval);
  }, [config.screenId, config.heartbeatInterval]);

  return (
    <div className="fixed top-0 right-0 p-2 bg-black bg-opacity-50 text-white text-xs">
      <div className={`w-2 h-2 rounded-full inline-block mr-2 ${connected ? 'bg-green-500' : 'bg-red-500'}`} />
      {connected ? 'Connected' : 'Disconnected'}
    </div>
  );
};
```

#### `src/App.tsx`
```typescript
import { useEffect } from 'react';
import { screenApi } from './services/api/screen';
import { usePlaylist } from './hooks/usePlaylist';
import { usePlayer } from './hooks/usePlayer';
import { Player } from './components/Player';
import { StatusBar } from './components/StatusBar';
import { loadConfig } from './config/default';

function App() {
  const config = loadConfig();
  const { playlist, loading, error } = usePlaylist();
  const { currentItem } = usePlayer(playlist?.items || []);

  useEffect(() => {
    // Register screen on mount
    screenApi.register({
      id: config.screenId,
      name: config.name,
      location: config.location,
      classification: config.classification,
    }).then(() => {
      console.log('[Screen] Registered successfully');
    }).catch((err) => {
      console.error('[Screen] Registration failed:', err);
    });
  }, [config]);

  if (loading && !playlist) {
    return (
      <div className="flex items-center justify-center h-screen bg-black text-white">
        <div>Loading playlist...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen bg-black text-white">
        <div>Error: {error.message}</div>
      </div>
    );
  }

  return (
    <div className="h-screen w-screen overflow-hidden">
      <StatusBar />
      <Player item={currentItem} />
    </div>
  );
}

export default App;
```

### Step 4: Environment Configuration

Create `.env` file:
```env
VITE_API_URL=http://localhost:8080
VITE_SCREEN_ID=screen-1
```

### Step 5: Run the Application

```bash
# Development
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

---

## 🧪 Testing with Backend

1. **Start Backend:**
   ```bash
   cd backend
   sbt run
   ```

2. **Verify Backend is Running:**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

3. **Start Client:**
   ```bash
   cd digital-display-client
   npm run dev
   ```

4. **Check Browser Console:**
   - Should see registration success
   - Should see playlist fetch
   - Should see heartbeat messages

---

## 📝 Configuration

### Option 1: Environment Variables
```env
VITE_API_URL=http://localhost:8080
VITE_SCREEN_ID=screen-1
```

### Option 2: localStorage
```javascript
localStorage.setItem('screenConfig', JSON.stringify({
  screenId: 'screen-1',
  name: 'My Screen',
  location: {
    city: 'Chennai',
    area: 'Airport',
    venueType: 'airport'
  },
  classification: 1
}));
```

### Option 3: Config File
Create `public/config.json`:
```json
{
  "screenId": "screen-1",
  "name": "Chennai Airport Screen 1",
  "location": {
    "city": "Chennai",
    "area": "Airport",
    "venueType": "airport"
  },
  "classification": 1,
  "backendUrl": "http://localhost:8080"
}
```

---

## 🐛 Troubleshooting

### Issue: Cannot connect to backend
- **Check:** Backend is running on port 8080
- **Check:** CORS is enabled in backend
- **Check:** API URL in config is correct

### Issue: Playlist is empty
- **Check:** Screen is registered in backend
- **Check:** There are active campaigns
- **Check:** Campaign targeting matches screen location

### Issue: Events not sending
- **Check:** Network connection
- **Check:** Backend is accepting events
- **Check:** Browser console for errors

---

## 🚀 Next Steps

1. **Add Video.js** for better video playback
2. **Add Error Recovery** for network failures
3. **Add Offline Mode** with cached playlists
4. **Add Configuration UI** for easy setup
5. **Add Health Monitoring** dashboard

---

## 📚 Full Documentation

See `docs/DIGITAL_DISPLAY_CLIENT_PLAN.md` for complete implementation plan.

