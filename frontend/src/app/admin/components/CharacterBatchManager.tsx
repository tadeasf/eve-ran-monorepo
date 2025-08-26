"use client";

import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { CheckCircle, XCircle, Loader2, Plus } from "lucide-react";

interface BatchResult {
  success: string[];
  failed: string[];
  duplicates: string[];
}

export function CharacterBatchManager() {
  const [characterIds, setCharacterIds] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<BatchResult | null>(null);

  const handleBatchAdd = async () => {
    if (!characterIds.trim()) return;

    setLoading(true);
    setResult(null);

    try {
      // Parse character IDs from textarea
      const ids = characterIds
        .split('\n')
        .map(line => line.trim())
        .filter(line => line.length > 0)
        .map(line => {
          // Extract just the ID if it's a full line with names/etc
          const match = line.match(/\d+/);
          return match ? match[0] : line;
        })
        .filter(id => /^\d+$/.test(id)); // Only keep valid numeric IDs

      if (ids.length === 0) {
        setResult({ success: [], failed: ['No valid character IDs found'], duplicates: [] });
        return;
      }

      // Simulate API call for now
      // In reality, this would call your backend API
      await new Promise(resolve => setTimeout(resolve, 2000));

      // Simulate some results
      const successCount = Math.floor(ids.length * 0.7);
      const duplicateCount = Math.floor(ids.length * 0.2);

      setResult({
        success: ids.slice(0, successCount),
        duplicates: ids.slice(successCount, successCount + duplicateCount),
        failed: ids.slice(successCount + duplicateCount),
      });

    } catch {
      setResult({ 
        success: [], 
        failed: ['Failed to process character IDs'], 
        duplicates: [] 
      });
    } finally {
      setLoading(false);
    }
  };

  const clearForm = () => {
    setCharacterIds('');
    setResult(null);
  };

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="character-ids">Character IDs</Label>
        <Textarea
          id="character-ids"
          placeholder="Enter character IDs, one per line:
95465499
1234567890
2345678901"
          value={characterIds}
          onChange={(e) => setCharacterIds(e.target.value)}
          rows={6}
          disabled={loading}
        />
        <p className="text-sm text-muted-foreground">
          Enter one character ID per line. Names and extra text will be ignored.
        </p>
      </div>

      <div className="flex gap-2">
        <Button 
          onClick={handleBatchAdd} 
          disabled={!characterIds.trim() || loading}
          className="flex-1"
        >
          {loading ? (
            <>
              <Loader2 className="mr-2 size-4 animate-spin" />
              Adding Characters...
            </>
          ) : (
            <>
              <Plus className="mr-2 size-4" />
              Add Characters
            </>
          )}
        </Button>
        <Button variant="outline" onClick={clearForm} disabled={loading}>
          Clear
        </Button>
      </div>

      {result && (
        <div className="space-y-4">
          <Separator />
          
          {result.success.length > 0 && (
            <Alert>
              <CheckCircle className="size-4" />
              <AlertDescription>
                <strong>Successfully added {result.success.length} character(s)</strong>
                <div className="mt-2 flex flex-wrap gap-1">
                  {result.success.slice(0, 5).map(id => (
                    <Badge key={id} variant="secondary" className="text-xs">
                      {id}
                    </Badge>
                  ))}
                  {result.success.length > 5 && (
                    <Badge variant="secondary" className="text-xs">
                      +{result.success.length - 5} more
                    </Badge>
                  )}
                </div>
              </AlertDescription>
            </Alert>
          )}

          {result.duplicates.length > 0 && (
            <Alert>
              <AlertDescription>
                <strong>{result.duplicates.length} character(s) already exist</strong>
                <div className="mt-2 flex flex-wrap gap-1">
                  {result.duplicates.slice(0, 5).map(id => (
                    <Badge key={id} variant="outline" className="text-xs">
                      {id}
                    </Badge>
                  ))}
                  {result.duplicates.length > 5 && (
                    <Badge variant="outline" className="text-xs">
                      +{result.duplicates.length - 5} more
                    </Badge>
                  )}
                </div>
              </AlertDescription>
            </Alert>
          )}

          {result.failed.length > 0 && (
            <Alert variant="destructive">
              <XCircle className="size-4" />
              <AlertDescription>
                <strong>Failed to add {result.failed.length} character(s)</strong>
                <div className="mt-2 flex flex-wrap gap-1">
                  {result.failed.slice(0, 5).map((id, index) => (
                    <Badge key={index} variant="destructive" className="text-xs">
                      {id}
                    </Badge>
                  ))}
                  {result.failed.length > 5 && (
                    <Badge variant="destructive" className="text-xs">
                      +{result.failed.length - 5} more
                    </Badge>
                  )}
                </div>
              </AlertDescription>
            </Alert>
          )}
        </div>
      )}
    </div>
  );
}
