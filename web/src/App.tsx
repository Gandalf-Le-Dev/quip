import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Upload, Edit3 } from 'lucide-react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { FileUploader } from './components/FileUploader';
import { PasteEditor } from './components/PasteEditor';
// import { ContentViewer } from './components/ContentViewer';

function App() {
  return (
    <Router>
      <Routes>
        {/* <Route path="/view/:id" element={<ContentViewer />} />
        <Route path="/:id" element={<ContentViewer />} /> */}
        <Route path="/" element={
          <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
            <div className="container mx-auto px-4 py-12">
              {/* Header */}
              <div className="text-center mb-12">
                <h1 className="text-5xl font-bold bg-gradient-to-r from-slate-900 to-slate-600 bg-clip-text text-transparent mb-4">
                  File Share & Pastebin
                </h1>
                <p className="text-lg text-slate-600 max-w-2xl mx-auto">
                  Quickly share files or create text snippets with secure, temporary links
                </p>
              </div>
              
              {/* Main Content */}
              <div className="max-w-4xl mx-auto">
                <Card className="shadow-xl border-0 bg-white/80 backdrop-blur-sm">
                  <CardHeader className="pb-6">
                    <CardTitle className="text-2xl text-center text-slate-800">
                      Choose your sharing method
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <Tabs defaultValue="file" className="w-full">
                      <TabsList className="grid w-full grid-cols-2 mb-8 h-12">
                        <TabsTrigger 
                          value="file" 
                          className="flex items-center gap-2 text-base font-medium data-[state=active]:bg-slate-900 data-[state=active]:text-white"
                        >
                          <Upload className="w-4 h-4" />
                          Upload File
                        </TabsTrigger>
                        <TabsTrigger 
                          value="paste"
                          className="flex items-center gap-2 text-base font-medium data-[state=active]:bg-slate-900 data-[state=active]:text-white"
                        >
                          <Edit3 className="w-4 h-4" />
                          Create Paste
                        </TabsTrigger>
                      </TabsList>
                      
                      <TabsContent value="file" className="mt-0">
                        <Card className="border-dashed border-2 border-slate-200 bg-slate-50/50">
                          <CardContent className="p-8">
                            <FileUploader />
                          </CardContent>
                        </Card>
                      </TabsContent>
                      
                      <TabsContent value="paste" className="mt-0">
                        <Card className="border-dashed border-2 border-slate-200 bg-slate-50/50">
                          <CardContent className="p-8">
                            <PasteEditor />
                          </CardContent>
                        </Card>
                      </TabsContent>
                    </Tabs>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        } />
      </Routes>
    </Router>
  );
}

export default App;