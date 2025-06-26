import type { File, Paste, TTL } from './types';

export interface IFileService {
  upload(file: File, ttl: TTL): Promise<File>;
  download(id: string): Promise<{ reader: ReadableStream<Uint8Array>; file: File }>;
  getInfo(id: string): Promise<File>;
}

export interface IPasteService {
  create(content: string, language: string, title: string, ttl: TTL): Promise<Paste>;
  get(id: string): Promise<Paste>;
  getRaw(id: string): Promise<string>;
}
