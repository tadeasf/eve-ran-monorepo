'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { KillListSkeleton } from '@/components/ui/skeleton';
import { formatDistanceToNow } from 'date-fns';
import { useInfiniteKills } from '@/hooks/useInfiniteKills';
import { useKillComments, useCreateComment } from '@/hooks/useKillComments';

interface Achievement {
  id: number;
  character_id: number;
  character_name: string;
  type: 'milestone' | 'record' | 'streak';
  title: string;
  description: string;
  created_at: string;
}

// Component to display comments for a single kill
function KillComments({ killmailId }: { killmailId: number }) {
  const [newCommentText, setNewCommentText] = useState('');
  
  const { data: comments = [], isLoading } = useKillComments(killmailId);
  const createComment = useCreateComment(killmailId);

  const handleAddComment = async () => {
    if (!newCommentText.trim()) return;

    try {
      await createComment.mutateAsync({
        comment: newCommentText,
        username: 'Current User', // TODO: Get from auth context
      });
      setNewCommentText('');
    } catch (error) {
      console.error('Failed to add comment:', error);
    }
  };

  if (isLoading) {
    return <div className="text-sm text-muted-foreground">Loading comments...</div>;
  }

  return (
    <div className="space-y-4">
      <h3 className="font-semibold text-sm">Comments</h3>
      
      {comments.map((comment) => (
        <div key={comment.id} className="flex gap-3">
          <Avatar className="h-8 w-8">
            <AvatarFallback>
              {comment.username.charAt(0).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 space-y-1">
            <div className="flex items-center gap-2">
              <span className="font-medium text-sm">
                {comment.username}
              </span>
              <span className="text-xs text-muted-foreground">
                {formatDistanceToNow(new Date(comment.created_at), {
                  addSuffix: true,
                })}
              </span>
            </div>
            <p className="text-sm">{comment.comment}</p>
          </div>
        </div>
      ))}

      {/* Add Comment */}
      <div className="flex gap-3">
        <Avatar className="h-8 w-8">
          <AvatarFallback>U</AvatarFallback>
        </Avatar>
        <div className="flex-1 space-y-2">
          <Textarea
            placeholder="Add a comment..."
            value={newCommentText}
            onChange={(e) => setNewCommentText(e.target.value)}
            className="min-h-[60px]"
          />
          <Button
            size="sm"
            onClick={handleAddComment}
            disabled={!newCommentText.trim() || createComment.isPending}
          >
            {createComment.isPending ? 'Posting...' : 'Post Comment'}
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function ActivityFeed() {
  const [achievements] = useState<Achievement[]>([]); // TODO: Implement achievements API
  const [selectedTab, setSelectedTab] = useState<'kills' | 'achievements'>('kills');
  
  // Fetch recent kills using infinite query (limited to first 20)
  const {
    data: killsData,
    isLoading,
  } = useInfiniteKills({});

  // Get the first page of kills (most recent 20)
  const recentKills = killsData?.pages[0]?.data || [];

  const formatISK = (value: number) => {
    if (value >= 1e9) return `${(value / 1e9).toFixed(2)}B`;
    if (value >= 1e6) return `${(value / 1e6).toFixed(2)}M`;
    return `${value.toLocaleString()}`;
  };

  if (isLoading) {
    return (
      <div className="container mx-auto py-8">
        <h1 className="text-3xl font-bold mb-6">Activity Feed</h1>
        <KillListSkeleton />
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 px-4 max-w-6xl">
      <h1 className="text-3xl font-bold mb-6">Activity Feed</h1>

      {/* Tab Navigation */}
      <div className="flex gap-2 mb-6">
        <Button
          variant={selectedTab === 'kills' ? 'default' : 'outline'}
          onClick={() => setSelectedTab('kills')}
        >
          Recent Kills
        </Button>
        <Button
          variant={selectedTab === 'achievements' ? 'default' : 'outline'}
          onClick={() => setSelectedTab('achievements')}
        >
          Achievements
        </Button>
      </div>

      {/* Recent Kills Tab */}
      {selectedTab === 'kills' && (
        <div className="space-y-4">
          {recentKills.length === 0 ? (
            <Card>
              <CardContent className="py-8 text-center text-muted-foreground">
                No recent kills to display
              </CardContent>
            </Card>
          ) : (
            recentKills.map((kill) => (
              <Card key={kill.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex justify-between items-start">
                    <div className="space-y-2">
                      <CardTitle className="text-lg">
                        Killmail #{kill.killmail_id}
                      </CardTitle>
                      <div className="flex gap-4 text-sm text-muted-foreground">
                        <span>System: {kill.solar_system_id}</span>
                        <span>Ship: {kill.victim_ship_type_id}</span>
                        <span>
                          {formatDistanceToNow(new Date(kill.killmail_time), {
                            addSuffix: true,
                          })}
                        </span>
                      </div>
                    </div>
                    <div className="text-right space-y-1">
                      {kill.zkill_data && (
                        <>
                          <Badge variant="secondary">
                            {formatISK(kill.zkill_data.total_value)} ISK
                          </Badge>
                          <div className="text-sm text-muted-foreground">
                            {kill.zkill_data.points} points
                          </div>
                        </>
                      )}
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <Separator className="mb-4" />
                  
                  {/* Comments Section */}
                  <KillComments killmailId={kill.killmail_id} />
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}

      {/* Achievements Tab */}
      {selectedTab === 'achievements' && (
        <div className="grid gap-4 md:grid-cols-2">
          {achievements.length === 0 ? (
            <Card className="col-span-2">
              <CardContent className="py-8 text-center text-muted-foreground">
                No achievements to display
              </CardContent>
            </Card>
          ) : (
            achievements.map((achievement) => (
              <Card
                key={achievement.id}
                className="hover:shadow-lg transition-shadow"
              >
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="space-y-1">
                      <Badge
                        variant={
                          achievement.type === 'milestone'
                            ? 'default'
                            : achievement.type === 'record'
                            ? 'destructive'
                            : 'secondary'
                        }
                      >
                        {achievement.type}
                      </Badge>
                      <CardTitle className="text-lg">
                        {achievement.title}
                      </CardTitle>
                    </div>
                    <span className="text-xs text-muted-foreground">
                      {formatDistanceToNow(new Date(achievement.created_at), {
                        addSuffix: true,
                      })}
                    </span>
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    {achievement.description}
                  </p>
                  <p className="text-sm font-medium mt-2">
                    {achievement.character_name}
                  </p>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}
    </div>
  );
}
