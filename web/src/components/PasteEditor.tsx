import { useState } from 'react';
import { 
  FileText, 
  Code, 
  Clock, 
  Eye, 
  Edit3, 
  Copy, 
  CheckCircle, 
  Loader2,
  Send,
  Hash
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { PasteService } from '../adapters/api/pasteService';
import type { Paste, TTL } from '../core/types';

const languages = [
  { value: 'text', label: 'Plain Text', icon: FileText },
  { value: 'javascript', label: 'JavaScript', icon: Code },
  { value: 'python', label: 'Python', icon: Code },
  { value: 'go', label: 'Go', icon: Code },
  { value: 'java', label: 'Java', icon: Code },
  { value: 'cpp', label: 'C++', icon: Code },
  { value: 'rust', label: 'Rust', icon: Code },
  { value: 'sql', label: 'SQL', icon: Code },
  { value: 'json', label: 'JSON', icon: Code },
  { value: 'yaml', label: 'YAML', icon: Code },
  { value: 'markdown', label: 'Markdown', icon: Hash },
];

const pasteService = new PasteService(import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080');

export function PasteEditor() {
  const [content, setContent] = useState('');
  const [language, setLanguage] = useState('text');
  const [title, setTitle] = useState('');
  const [ttl, setTtl] = useState<TTL>('24h');
  const [activeTab, setActiveTab] = useState('edit');
  const [result, setResult] = useState<Paste | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [copiedUrl, setCopiedUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async () => {
    if (!content.trim()) return;
    
    setSubmitting(true);
    setError(null);
    setResult(null);

    try {
      const createdPaste = await pasteService.create(content, language, title, ttl);
      setResult(createdPaste);
    } catch (err: any) {
      setError(err.message || 'An unknown error occurred during paste creation.');
    } finally {
      setSubmitting(false);
    }
  };

  const copyToClipboard = async (url: string, type: string) => {
    try {
      await navigator.clipboard.writeText(`${window.location.origin}${url}`);
      setCopiedUrl(type);
      setTimeout(() => setCopiedUrl(null), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const renderPreview = () => {
    if (!content.trim()) {
      return (
        <div className="flex items-center justify-center h-64 text-slate-400">
          <div className="text-center">
            <Code className="w-12 h-12 mx-auto mb-2 opacity-50" />
            <p>Enter some content to see the preview</p>
          </div>
        </div>
      );
    }

    return (
      <div className="bg-slate-900 rounded-lg p-4 font-mono text-sm overflow-auto max-h-96">
        <pre className="text-slate-100 whitespace-pre-wrap">{content}</pre>
      </div>
    );
  };

  const selectedLanguage = languages.find(lang => lang.value === language);

  return (
    <div className="space-y-6">
      {/* Header Form */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Title Input */}
        <div className="space-y-2">
          <Label htmlFor="title" className="text-sm font-medium flex items-center gap-2">
            <FileText className="w-4 h-4" />
            Title (optional)
          </Label>
          <Input
            id="title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Untitled paste"
            className="w-full"
          />
        </div>
        
        {/* Language Selection */}
        <div className="space-y-2">
          <Label className="text-sm font-medium flex items-center gap-2">
            <Code className="w-4 h-4" />
            Language
          </Label>
          <Select value={language} onValueChange={setLanguage}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Select language" />
            </SelectTrigger>
            <SelectContent>
              {languages.map(lang => {
                const IconComponent = lang.icon;
                return (
                  <SelectItem key={lang.value} value={lang.value}>
                    <div className="flex items-center gap-2">
                      <IconComponent className="w-4 h-4" />
                      {lang.label}
                    </div>
                  </SelectItem>
                );
              })}
            </SelectContent>
          </Select>
        </div>
        
        {/* TTL Selection */}
        <div className="space-y-2">
          <Label className="text-sm font-medium flex items-center gap-2">
            <Clock className="w-4 h-4" />
            Expiration
          </Label>
          <Select value={ttl} onValueChange={(value) => setTtl(value as TTL)}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Select expiration" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1h">1 hour</SelectItem>
              <SelectItem value="24h">24 hours</SelectItem>
              <SelectItem value="72h">3 days</SelectItem>
              <SelectItem value="168h">1 week</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Content Editor */}
      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg flex items-center gap-2">
              {selectedLanguage && <selectedLanguage.icon className="w-5 h-5" />}
              Content Editor
            </CardTitle>
            <div className="text-sm text-slate-500">
              {content.length} characters
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
            <TabsList className="grid w-full grid-cols-2 mb-4">
              <TabsTrigger value="edit" className="flex items-center gap-2">
                <Edit3 className="w-4 h-4" />
                Edit
              </TabsTrigger>
              <TabsTrigger value="preview" className="flex items-center gap-2">
                <Eye className="w-4 h-4" />
                Preview
              </TabsTrigger>
            </TabsList>
            
            <TabsContent value="edit" className="mt-0">
              <Textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className="min-h-[400px] font-mono text-sm resize-none"
                placeholder="Paste your content here..."
              />
            </TabsContent>
            
            <TabsContent value="preview" className="mt-0">
              <div className="border rounded-lg p-4 min-h-[400px]">
                {renderPreview()}
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* Submit Button */}
      <Button
        onClick={handleSubmit}
        disabled={!content.trim() || submitting}
        className="w-full h-12 text-base"
        size="lg"
      >
        {submitting ? (
          <>
            <Loader2 className="w-5 h-5 mr-2 animate-spin" />
            Creating Paste...
          </>
        ) : (
          <>
            <Send className="w-5 h-5 mr-2" />
            Create Paste
          </>
        )}
      </Button>

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
                <h3 className="font-semibold text-green-800">Paste Created Successfully!</h3>
              </div>
              
              {/* Paste Info */}
              <div className="bg-white rounded-lg p-3 border border-green-200">
                <div className="text-sm space-y-1">
                  <div className="flex justify-between">
                    <span className="font-medium">Title:</span>
                    <span className="text-slate-600">{result.Title || 'Untitled'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="font-medium">Language:</span>
                    <span className="text-slate-600">
                      {languages.find(l => l.value === result.Language)?.label}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="font-medium">ID:</span>
                    <code className="text-slate-600 font-mono text-xs">{result.ID}</code>
                  </div>
                </div>
              </div>

              {/* Action Links */}
              <div className="space-y-3">
                {/* Raw URL */}
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-green-800">Raw Content URL:</Label>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 px-3 py-2 bg-slate-100 rounded text-sm font-mono border">
                      {window.location.origin}{result.raw}
                    </code>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => copyToClipboard(result.raw, 'raw')}
                    >
                      {copiedUrl === 'raw' ? (
                        <CheckCircle className="w-4 h-4 text-green-600" />
                      ) : (
                        <Copy className="w-4 h-4" />
                      )}
                    </Button>
                  </div>
                </div>

                {/* View URL */}
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

                {/* cURL Command */}
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-green-800">cURL Command:</Label>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 px-3 py-2 bg-slate-900 text-green-400 rounded text-sm font-mono border">
                      curl {window.location.origin}{result.raw}
                    </code>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => copyToClipboard(`curl ${window.location.origin}${result.raw}`, 'curl')}
                    >
                      {copiedUrl === 'curl' ? (
                        <CheckCircle className="w-4 h-4 text-green-600" />
                      ) : (
                        <Copy className="w-4 h-4" />
                      )}
                    </Button>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex gap-2 pt-2">
                  <a href={result.view} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="sm" className="w-full">
                      <Eye className="w-4 h-4 mr-2" />
                      View Paste
                    </Button>
                  </a>
                  <a href={result.raw} target="_blank" rel="noopener noreferrer" className="flex-1">
                    <Button size="sm" variant="outline" className="w-full">
                      <FileText className="w-4 h-4 mr-2" />
                      Raw Content
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
