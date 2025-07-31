import { useCallback } from 'react';

import { useEmitters } from '@utils';
import { useHotkeys } from 'react-hotkeys-hook';

import { CancelAllFetches } from '../../wailsjs/go/app/App';

interface UseGlobalEscapeProps {
  enabled?: boolean;
  onEscape?: () => void;
}

/**
 * Global escape key handler that cancels all active fetches
 * and optionally runs a custom callback
 */
export const useGlobalEscape = ({
  enabled = true,
  onEscape,
}: UseGlobalEscapeProps = {}) => {
  const { emitStatus } = useEmitters();

  const handleEscape = useCallback(async () => {
    try {
      const cancelledCount = await CancelAllFetches();
      if (cancelledCount > 0) {
        emitStatus(
          `Cancelled ${cancelledCount} active fetch${cancelledCount === 1 ? '' : 'es'}`,
        );
      }
      onEscape?.();
    } catch (error) {
      console.error('Failed to cancel fetches:', error);
      emitStatus('Failed to cancel active fetches');
    }
  }, [emitStatus, onEscape]);

  useHotkeys(
    'esc',
    (e) => {
      e.preventDefault();
      handleEscape();
    },
    {
      enableOnFormTags: true,
      enabled,
    },
  );

  return { cancelAllFetches: handleEscape };
};
