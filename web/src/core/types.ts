export interface File {
  ID: string;
  OriginalName: string;
  Size: number;
  ContentType: string;
  StorageKey: string;
  Downloads: number;
  MaxDownloads: number;
  CreatedAt: string;
  ExpiresAt: string;
  download: string; // API generated download URL
  view: string;     // API generated view URL
}

export interface Paste {
  ID: string;
  Content: string;
  Language: string;
  Title: string;
  Views: number;
  MaxViews: number;
  CreatedAt: string;
  ExpiresAt: string;
  raw: string;      // API generated raw URL
  view: string;     // API generated view URL
}

export type TTL = '1h' | '24h' | '72h' | '168h'; // Added for recompilation
export const DUMMY_EXPORT = true;
