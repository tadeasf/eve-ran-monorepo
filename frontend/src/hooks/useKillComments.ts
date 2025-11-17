import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiFetch, apiFetchJson } from '@/lib/api';

export interface Comment {
  id: number;
  kill_id: number;
  killmail_id: number;
  user_id: number;
  username: string;
  comment: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCommentData {
  username: string;
  comment: string;
}

// Fetch comments for a specific kill
export function useKillComments(killmailId: number) {
  return useQuery({
    queryKey: ['killComments', killmailId],
    queryFn: () => apiFetchJson<Comment[]>(`/kills/${killmailId}/comments`),
    enabled: !!killmailId,
  });
}

// Create a new comment
export function useCreateComment(killmailId: number) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCommentData) =>
      apiFetch(`/kills/${killmailId}/comments`, {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      // Invalidate and refetch comments
      queryClient.invalidateQueries({ queryKey: ['killComments', killmailId] });
    },
  });
}

// Update a comment
export function useUpdateComment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ commentId, comment }: { commentId: number; comment: string }) =>
      apiFetch(`/comments/${commentId}`, {
        method: 'PUT',
        body: JSON.stringify({ comment }),
      }),
    onSuccess: () => {
      // Invalidate all comment queries (we don't know which kill it belongs to)
      queryClient.invalidateQueries({ queryKey: ['killComments'] });
    },
  });
}

// Delete a comment
export function useDeleteComment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (commentId: number) =>
      apiFetch(`/comments/${commentId}`, {
        method: 'DELETE',
      }),
    onSuccess: () => {
      // Invalidate all comment queries
      queryClient.invalidateQueries({ queryKey: ['killComments'] });
    },
  });
}
