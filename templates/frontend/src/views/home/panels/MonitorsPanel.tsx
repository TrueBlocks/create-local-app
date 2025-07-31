import { useEffect, useState } from 'react';

import { GetMonitorsSummary } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useEvent, useIconSets } from '@hooks';
import { Badge, Button, Group, Stack, Text } from '@mantine/core';
import { msgs, types } from '@models';
import { Log } from '@utils';

interface MonitorsPanelProps {
  onViewAll?: () => void;
  onAddMonitor?: () => void;
}

export const MonitorsPanel = ({
  onViewAll,
  onAddMonitor,
}: MonitorsPanelProps) => {
  const [summary, setSummary] = useState<types.Summary>({
    totalCount: 0,
    facetCounts: {},
    customData: {
      activeCount: 0,
      stagedCount: 0,
      deletedCount: 0,
      emptyCount: 0,
    },
    lastUpdated: Date.now(),
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { Monitors } = useIconSets();

  const fetchMonitorsSummary = async (showLoading = true) => {
    try {
      if (showLoading) setLoading(true);

      // Use the new GetMonitorsSummary API that returns pre-computed summaries
      const summaryData = await GetMonitorsSummary(
        types.Payload.createFrom({}),
      );

      if (summaryData) {
        // Use the backend Summary structure directly
        setSummary(summaryData);
      }
      setError(null);
    } catch (err) {
      Log(`Error fetching monitors summary: ${err}`);
      setError('Failed to load monitors');
    } finally {
      if (showLoading) setLoading(false);
    }
  };

  useEffect(() => {
    fetchMonitorsSummary();
  }, []);

  // Listen for the new consolidated collection state changes
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      // Update when monitors collection data changes
      if (payload?.collection === 'monitors') {
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

  const getHealthStatus = () => {
    if (summary.totalCount === 0) return 'inactive';
    const deletedCount = (summary.customData?.deletedCount as number) || 0;
    const activeCount = (summary.customData?.activeCount as number) || 0;
    if (deletedCount > activeCount) return 'warning';
    return 'healthy';
  };

  return (
    <DashboardCard
      title="Monitors"
      subtitle={`${summary.totalCount} addresses`}
      icon={<Monitors size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={getHealthStatus()}
            label="Monitor Health"
            count={(summary.customData?.activeCount as number) || 0}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          <Badge size="sm" variant="light" color="green">
            Active: {(summary.customData?.activeCount as number) || 0}
          </Badge>
          <Badge size="sm" variant="light" color="blue">
            Staged: {(summary.customData?.stagedCount as number) || 0}
          </Badge>
          <Badge size="sm" variant="light" color="red">
            Deleted: {(summary.customData?.deletedCount as number) || 0}
          </Badge>
          <Badge size="sm" variant="light" color="gray">
            Empty: {(summary.customData?.emptyCount as number) || 0}
          </Badge>
        </div>

        <Text size="xs" c="dimmed">
          Track and manage Ethereum address monitoring
        </Text>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onAddMonitor}>
            Add Monitor
          </Button>
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View All
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
