import type { IFileService } from '../../core/ports';
import type { File, TTL } from '../../core/types';

export class FileService implements IFileService {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async upload(file: any, ttl: TTL): Promise<File> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('ttl', ttl);

    const response = await fetch(`${this.baseUrl}/api/file`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to upload file');
    }

    return response.json();
  }

  async download(id: string): Promise<{ reader: ReadableStream<Uint8Array>; file: File }> {
    const response = await fetch(`${this.baseUrl}/api/file/${id}`);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to download file');
    }

    const fileInfoResponse = await fetch(`${this.baseUrl}/api/file/${id}/info`);
    if (!fileInfoResponse.ok) {
      const errorData = await fileInfoResponse.json();
      throw new Error(errorData.message || 'Failed to get file info');
    }
    const fileInfo: File = await fileInfoResponse.json();

    if (!response.body) {
      throw new Error('No readable stream found for download');
    }

    return { reader: response.body, file: fileInfo };
  }

  async getInfo(id: string): Promise<File> {
    const response = await fetch(`${this.baseUrl}/api/file/${id}/info`);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to get file info');
    }

    return response.json();
  }
}
