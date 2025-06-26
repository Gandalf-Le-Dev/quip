import React, { useState, useCallback } from 'react';
import { Upload, Clock, Download, Eye, Copy, CheckCircle, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Card, CardContent } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { FileService } from '../adapters/api/fileService';
import type { File, TTL } from '../core/types';

const fileService = new FileService(import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080');

export function FileUploader() {
  const [uploading, setUploading] = useState(false);
  const [result, setResult] = useState<File | null>(null);
  const [ttl, setTtl] = useState<TTL>('24h');
  const [isDragActive, setIsDragActive] = useState(false);
  const [copiedUrl, setCopiedUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleFileUpload = useCallback(async (files: FileList | null) => {
    if (!files || files.length === 0) return;
    
    const fileToUpload = files[0];
    setUploading(true);
    setError(null);
    setResult(null);

    try {
      const uploadedFile = await fileService.upload(fileToUpload, ttl);
      setResult(uploadedFile);
    } catch (err: any) {
      setError(err.message || 'An unknown error occurred during upload.');
    } finally {
      setUploading(false);
    }
  }, [ttl]);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragActive(false);
    handleFileUpload(e.dataTransfer.files);
  }, [handleFileUpload]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragActive(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragActive(false);
  }, []);

  const handleFileInput = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileUpload(e.target.files);
  }, [handleFileUpload]);

  const copyToClipboard = async (url: string, type: string) => {
    try {
      await navigator.clipboard.writeText(`${window.location.origin}${url}`);
      setCopiedUrl(type);
      setTimeout(() => setCopiedUrl(null), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="space-y-6">
      {/* TTL Selection */}
      <div className="space-y-2">
        <Label htmlFor="ttl-select" className="text-sm font-medium flex items-center gap-2">
          <Clock className="w-4 h-4" />
          Expiration Time
        </Label>
        <Select value={ttl} onValueChange={(value) => setTtl(value as TTL)}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select expiration time" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1h">1 hour</SelectItem>
            <SelectItem value="24h">24 hours</SelectItem>
            <SelectItem value="72h">3 days</SelectItem>
            <SelectItem value="168h">1 week</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Upload Area */}
      <Card className={`transition-all duration-200 ${
        isDragActive 
          ? 'border-2 border-blue-500 bg-blue-50 shadow-lg scale-105' 
          : 'border-2 border-dashed border-slate-300 hover:border-slate-400'
      }`}>
        <CardContent className="p-8">
          <div
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            className="text-center cursor-pointer"
            onClick={() => document.getElementById('file-input')?.click()}
          >
            <input
              id="file-input"
              type="file"
              className="hidden"
              onChange={handleFileInput}
              disabled={uploading}
            />
            
            <div className="space-y-4">
              {uploading ? (
                <>
                  <Loader2 className="w-12 h-12 mx-auto text-blue-500 animate-spin" />
                  <div className="space-y-2">
                    <p className="text-lg font-medium text-slate-700">Uploading...</p>
                    <div className="w-64 mx-auto bg-slate-200 rounded-full h-2">
                      <div className="bg-blue-500 h-2 rounded-full animate-pulse w-3/4"></div>
                    </div>
                  </div>
                </>
              ) : isDragActive ? (
                <>
                  <Upload className="w-12 h-12 mx-auto text-blue-500" />
                  <p className="text-lg font-medium text-blue-600">Drop the file here!</p>
                </>
              ) : (
                <>
                  <Upload className="w-12 h-12 mx-auto text-slate-400" />
                  <div className="space-y-2">
                    <p className="text-lg font-medium text-slate-700">
                      Drag and drop your file here
                    </p>
                    <p className="text-sm text-slate-500">or click to browse files</p>
                  </div>
                  <Button variant="outline" className="mt-4">
                    <Upload className="w-4 h-4 mr-2" />
                    Choose File
                  </Button>
                </>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {error && (
        <Alert className="border-red-200 bg-red-50">
          <AlertDescription>
            <div className="flex items-center gap-2">
              <h3 className="font-semibold text-red-800">Error:</h3>
              <p className="text-red-700">{error}</p>
            </div>
          </AlertDescription>
        </Alert>
      )}

      {/* Success Result */}
      {result && (
        <Alert className="border-green-200 bg-green-50">
          <CheckCircle className="h-4 w-4 text-green-600" />
          <AlertDescription>
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                <h3 className="font-semibold text-green-800">Upload Successful!</h3>
              </div>
              
              {/* File Info */}
              <div className="bg-white rounded-lg p-3 border border-green-200">
                <div className="text-sm space-y-1">
                  <div className="flex justify-between">
                    <span className="font-medium">Filename:</span>
                    <span className="text-slate-600">{result.OriginalName}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="font-medium">Size:</span>
                    <span className="text-slate-600">{formatFileSize(result.Size)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="font-medium">Type:</span>
                    <span className="text-slate-600">{result.ContentType}</span>
                  </div>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="space-y-3">
                {/* Download Link */}
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-green-800">Download URL:</Label>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 px-3 py-2 bg-slate-100 rounded text-sm font-mono border">
                      {window.location.origin}{result.download}
                    </code>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => copyToClipboard(result.download, 'download')}
                    >
                      {copiedUrl === 'download' ? (
                        <CheckCircle className="w-4 h-4 text-green-600" />
                      ) : (
                        <Copy className="w-4 h-4" />
                      )}
                    </Button>
                  </div>
                </div>

                {/* View Link */}
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-green-800">View URL:</Label>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 px-3 py-2 bg-slate-100 rounded text-sm font-mono border">
                      {window.location.origin}{result.view}
                    </code>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => copyToClipboard(result.view, 'view')}
                    >
                      {copiedUrl === 'view' ? (
                        <CheckCircle className="w-4 h-4 text-green-600" />
                      ) : (
                        <Copy className="w-4 h-4" />
                      )}
                    </Button>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex gap-2 pt-2">
                  <a href={result.download} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="sm" className="w-full">
                      <Download className="w-4 h-4 mr-2" />
                      Download
                    </Button>
                  </a>
                  <a href={result.view} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="sm" variant="outline" className="w-full">
                      <Eye className="w-4 h-4 mr-2" />
                      View File
                    </Button>
                  </a>
                </div>
              </div>
            </div>
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}
