import { useEffect, useState } from 'react';

import { GetAbisSummary } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useEvent, useIconSets } from '@hooks';
import { Badge, Button, Group, Stack, Text } from '@mantine/core';
import { msgs, types } from '@models';
import { Log } from '@utils';

interface AbisPanelProps {
  onViewAll?: () => void;
  onAddAbi?: () => void;
}

export const AbisPanel = ({ onViewAll, onAddAbi }: AbisPanelProps) => {
  const [summary, setSummary] = useState<types.Summary>({
    totalCount: 0,
    facetCounts: {},
    customData: {
      functionsCount: 0,
      eventsCount: 0,
      knownCount: 0,
      downloadedCount: 0,
    },
    lastUpdated: Date.now(),
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { Abis } = useIconSets();

  const fetchAbisSummary = async (showLoading = true) => {
    try {
      if (showLoading) setLoading(true);

      // Use the new GetAbisSummary API that returns pre-computed summaries
      const summaryData = await GetAbisSummary(types.Payload.createFrom({}));

      if (summaryData) {
        // Use the backend Summary structure directly
        setSummary(summaryData);
      }
      setError(null);
    } catch (err) {
      Log(`Error fetching Abis summary: ${err}`);
      setError('Failed to load Abis');
    } finally {
      if (showLoading) setLoading(false);
    }
  };

  useEffect(() => {
    fetchAbisSummary();
  }, []);

  // Listen for the new consolidated collection state changes
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      // Update when abis collection data changes
      if (payload?.collection === 'abis') {
        // Extract summary directly from the event payload - no API call needed!
        const summary = payload.summary as types.Summary | undefined;
        if (summary) {
          setSummary(summary);
        }

        // Update loading state based on collection state
        const state = payload.state as types.LoadState | undefined;
        setLoading(state === types.LoadState.FETCHING);

        // Handle errors
        const error = payload.error as string | undefined;
        if (error) {
          setError(error);
        } else {
          setError(null);
        }
      }
    },
  );

  return (
    <DashboardCard
      title="Abis"
      subtitle={`${summary.totalCount} contracts`}
      icon={<Abis size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={summary.totalCount > 0 ? 'healthy' : 'inactive'}
            label="ABI Database"
            count={summary.totalCount}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          <Badge size="sm" variant="light" color="blue">
            Functions: {(summary.customData?.functionsCount as number) || 0}
          </Badge>
          <Badge size="sm" variant="light" color="purple">
            Events: {(summary.customData?.eventsCount as number) || 0}
          </Badge>
          <Badge size="sm" variant="light" color="green">
            Known: {(summary.customData?.knownCount as number) || 0}
          </Badge>
          <Badge size="sm" variant="light" color="orange">
            Downloaded: {(summary.customData?.downloadedCount as number) || 0}
          </Badge>
        </div>

        <Text size="xs" c="dimmed">
          Smart contract interfaces and function signatures
        </Text>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onAddAbi}>
            Add ABI
          </Button>
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View All
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
