import { useEffect, useState } from 'react';

import { GetChunksSummary } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useEvent, useIconSets } from '@hooks';
import { Badge, Button, Group, Stack, Text } from '@mantine/core';
import { msgs, types } from '@models';
import { Log } from '@utils';

interface ChunksPanelProps {
  onViewAll?: () => void;
}

// Extended summary type that includes UI-specific health status
type ChunksSummaryState = types.Summary & {
  indexHealth: 'healthy' | 'warning' | 'error' | 'loading';
};

export const ChunksPanel = ({ onViewAll }: ChunksPanelProps) => {
  const [summary, setSummary] = useState<ChunksSummaryState>({
    totalCount: 0,
    facetCounts: {},
    customData: {
      statsCount: 0,
      indexCount: 0,
      bloomsCount: 0,
      manifestCount: 0,
    },
    lastUpdated: Date.now(),
    indexHealth: 'loading',
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { Chunks } = useIconSets();

  const fetchChunksSummary = async (showLoading = true) => {
    try {
      if (showLoading) setLoading(true);

      // Use the new GetChunksSummary API that returns pre-computed summaries
      const summaryData = await GetChunksSummary(types.Payload.createFrom({}));

      if (summaryData) {
        // Use the backend Summary structure directly
        setSummary({
          ...summaryData,
          indexHealth: summaryData.totalCount > 0 ? 'healthy' : 'warning',
        });
      }
      setError(null);
    } catch (err) {
      Log(`Error fetching chunks summary: ${err}`);
      setError('Failed to load chunks');
      setSummary((prev) => ({ ...prev, indexHealth: 'error' }));
    } finally {
      if (showLoading) setLoading(false);
    }
  };

  useEffect(() => {
    fetchChunksSummary();
  }, []);

  // Listen for the new consolidated collection state changes
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      // Update when chunks collection data changes
      if (payload?.collection === 'chunks') {
        // Extract summary directly from the event payload - no API call needed!
        const summary = payload.summary as types.Summary | undefined;
        if (summary) {
          setSummary({
            ...summary,
            indexHealth: summary.totalCount > 0 ? 'healthy' : 'warning',
          });
        }

        // Update loading state based on collection state
        const state = payload.state as types.LoadState | undefined;
        setLoading(state === types.LoadState.FETCHING);

        // Handle errors
        const error = payload.error as string | undefined;
        if (error) {
          setError(error);
          setSummary((prev) => ({ ...prev, indexHealth: 'error' }));
        } else {
          setError(null);
        }
      }
    },
  );

  const formatNumber = (num: number) => {
    if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
    return num.toString();
  };

  return (
    <DashboardCard
      title="Chunks"
      subtitle="Index statistics"
      icon={<Chunks size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={summary.indexHealth}
            label="Index Health"
            count={summary.totalCount}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          <Badge size="sm" variant="light" color="blue">
            Stats:{' '}
            {formatNumber((summary.customData?.statsCount as number) || 0)}
          </Badge>
          <Badge size="sm" variant="light" color="green">
            Index:{' '}
            {formatNumber((summary.customData?.indexCount as number) || 0)}
          </Badge>
          <Badge size="sm" variant="light" color="purple">
            Blooms:{' '}
            {formatNumber((summary.customData?.bloomsCount as number) || 0)}
          </Badge>
          <Badge size="sm" variant="light" color="orange">
            Manifest:{' '}
            {formatNumber((summary.customData?.manifestCount as number) || 0)}
          </Badge>
        </div>

        <Text size="xs" c="dimmed">
          View stats, index, blooms, and manifest data
        </Text>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View Details
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
