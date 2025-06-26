import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { FileService } from '../adapters/api/fileService';
import { PasteService } from '../adapters/api/pasteService';
import type { File, Paste } from '../core/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Loader2, FileText, Code, Download, Eye } from 'lucide-react';
import { Button } from '@/components/ui/button';

const fileService = new FileService(import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080');
const pasteService = new PasteService(import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080');

export function ContentViewer() {
  const { id } = useParams<{ id: string }>();
  const [content, setContent] = useState<File | Paste | null>(null);
  const [rawContent, setRawContent] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const isFile = content && 'OriginalName' in content && typeof (content as any).OriginalName === 'string';
  const isPaste = content && 'Content' in content && typeof (content as any).Content === 'string';

  useEffect(() => {
    const fetchContent = async () => {
      setLoading(true);
      setError(null);
      setContent(null);
      setRawContent(null);

      if (!id) {
        setError('No content ID provided.');
        setLoading(false);
        return;
      }

      try {
        // Try fetching as paste first
        const paste = await pasteService.get(id);
        setContent(paste);
        const raw = await pasteService.getRaw(id);
        setRawContent(raw);
      } catch (pasteError: any) {
        // If not a paste, try fetching as file
        try {
          const file = await fileService.getInfo(id);
          setContent(file);
          // For files, we don't fetch raw content directly here, just info
        } catch (fileError: any) {
          setError(`Content not found or expired: ${fileError.message}`);
        }
      } finally {
        setLoading(false);
        console.log('Content fetched:', content);
        console.log('Is File:', isFile);
        console.log('Is Paste:', isPaste);
      }
    };

    fetchContent();
  }, [id, isFile, isPaste]);

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100">
        <Loader2 className="w-16 h-16 text-blue-500 animate-spin" />
        <p className="text-lg text-slate-700 ml-4">Loading content...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 p-4">
        <Alert className="border-red-200 bg-red-50 max-w-md">
          <AlertDescription>
            <h3 className="font-semibold text-red-800">Error:</h3>
            <p className="text-red-700">{error}</p>
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  if (!content) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 p-4">
        <Alert className="border-yellow-200 bg-yellow-50 max-w-md">
          <AlertDescription>
            <h3 className="font-semibold text-yellow-800">Content Not Found</h3>
            <p className="text-yellow-700">The requested content could not be found or has expired.</p>
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 py-12">
      <div className="container mx-auto px-4 max-w-4xl">
        <Card className="shadow-xl border-0 bg-white/80 backdrop-blur-sm">
          <CardHeader className="pb-6 text-center">
            <CardTitle className="text-3xl font-bold text-slate-800 flex items-center justify-center gap-3">
              {isFile ? <FileText className="w-8 h-8" /> : <Code className="w-8 h-8" />}
              {isFile ? (content as File).OriginalName : (content as Paste).Title || 'Untitled Paste'}
            </CardTitle>
            <p className="text-slate-600 text-lg">
              {isFile ? 'File Details' : 'Paste Content'}
            </p>
          </CardHeader>
          <CardContent>
            {isFile ? (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4 text-lg">
                  <div className="font-medium">ID:</div>
                  <div>{(content as File).ID}</div>
                  <div className="font-medium">Size:</div>
                  <div>{formatFileSize((content as File).Size)}</div>
                  <div className="font-medium">Content Type:</div>
                  <div>{(content as File).ContentType}</div>
                  <div className="font-medium">Uploaded At:</div>
                  <div>{new Date((content as File).CreatedAt).toLocaleString()}</div>
                  <div className="font-medium">Expires At:</div>
                  <div>{new Date((content as File).ExpiresAt).toLocaleString()}</div>
                </div>
                <div className="flex gap-4 mt-6">
                  <a href={(content as File).download} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="lg" className="w-full">
                      <Download className="w-5 h-5 mr-2" />
                      Download File
                    </Button>
                  </a>
                  <a href={(content as File).view} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="lg" variant="outline" className="w-full">
                      <Eye className="w-5 h-5 mr-2" />
                      View Raw (if applicable)
                    </Button>
                  </a>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4 text-lg">
                  <div className="font-medium">ID:</div>
                  <div>{(content as Paste).ID}</div>
                  <div className="font-medium">Language:</div>
                  <div>{(content as Paste).Language}</div>
                  <div className="font-medium">Views:</div>
                  <div>{(content as Paste).Views}</div>
                  <div className="font-medium">Uploaded At:</div>
                  <div>{new Date((content as Paste).CreatedAt).toLocaleString()}</div>
                  <div className="font-medium">Expires At:</div>
                  <div>{new Date((content as Paste).ExpiresAt).toLocaleString()}</div>
                </div>
                <div className="bg-slate-900 rounded-lg p-6 mt-6 overflow-auto max-h-[500px]">
                  <pre className="text-slate-100 whitespace-pre-wrap text-sm">{rawContent}</pre>
                </div>
                <div className="flex gap-4 mt-6">
                  <a href={(content as Paste).raw} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="lg" className="w-full">
                      <FileText className="w-5 h-5 mr-2" />
                      View Raw
                    </Button>
                  </a>
                  <a href={(content as Paste).view} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="lg" variant="outline" className="w-full">
                      <Eye className="w-5 h-5 mr-2" />
                      View Formatted
                    </Button>
                  </a>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
