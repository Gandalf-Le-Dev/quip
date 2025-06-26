import type { IPasteService } from '../../core/ports';
import type { Paste, TTL } from '../../core/types';

export class PasteService implements IPasteService {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async create(content: string, language: string, title: string, ttl: TTL): Promise<Paste> {
    const response = await fetch(`${this.baseUrl}/api/paste`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        content,
        language,
        title,
        ttl,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to create paste');
    }

    return response.json();
  }

  async get(id: string): Promise<Paste> {
    const response = await fetch(`${this.baseUrl}/api/paste/${id}`);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to get paste');
    }

    return response.json();
  }

  async getRaw(id: string): Promise<string> {
    const response = await fetch(`${this.baseUrl}/api/paste/${id}/raw`);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to get raw paste content');
    }

    return response.text();
  }
}
