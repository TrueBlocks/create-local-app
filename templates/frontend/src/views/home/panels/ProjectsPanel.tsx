import { useEffect, useState } from 'react';

import { GetOpenProjects } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useEvent, useIconSets } from '@hooks';
import { Badge, Button, Group, Stack, Text } from '@mantine/core';
import { msgs } from '@models';
import { Log } from '@utils';

interface Project {
  id: string;
  name: string;
  path: string;
  isActive: boolean;
  isDirty: boolean;
  lastOpened: string;
  createdAt: string;
}

interface ProjectsPanelProps {
  onViewAll?: () => void;
  onNewProject?: () => void;
}

export const ProjectsPanel = ({
  onViewAll,
  onNewProject,
}: ProjectsPanelProps) => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { File } = useIconSets();

  const fetchProjects = async (showLoading = true) => {
    try {
      if (showLoading) setLoading(true);
      const openProjects = await GetOpenProjects();
      const typedProjects = openProjects as unknown as Project[];
      const sorted = typedProjects
        .sort(
          (a, b) =>
            new Date(b.lastOpened).getTime() - new Date(a.lastOpened).getTime(),
        )
        .slice(0, 3);
      setProjects(sorted);
      setError(null);
    } catch (err) {
      Log(`Error fetching projects: ${err}`);
      setError('Failed to load projects');
    } finally {
      if (showLoading) setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  // Listen for project changes and refresh silently (no loading state)
  useEvent(msgs.EventType.MANAGER, () => {
    fetchProjects(false);
  });

  const activeProject = projects.find((p) => p.isActive);
  const recentProjects = projects.filter((p) => !p.isActive);

  return (
    <DashboardCard
      title="Projects"
      subtitle={`${projects.length} open`}
      icon={<File size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={
              activeProject
                ? 'healthy'
                : projects.length > 0
                  ? 'warning'
                  : 'inactive'
            }
            label={activeProject ? 'Active Project' : 'Project Status'}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          <Badge size="sm" variant="light" color="blue">
            Total: {projects.length}
          </Badge>
          <Badge size="sm" variant="light" color="green">
            Active: {activeProject ? 1 : 0}
          </Badge>
          <Badge size="sm" variant="light" color="orange">
            Recent: {recentProjects.length}
          </Badge>
          <Badge size="sm" variant="light" color="purple">
            Dirty: {projects.filter((p) => p.isDirty).length}
          </Badge>
        </div>

        <div>
          {activeProject ? (
            <Text size="xs" c="dimmed">
              Active: {activeProject.name}
              {activeProject.isDirty && ' *'}
            </Text>
          ) : recentProjects.length > 0 ? (
            <Text size="xs" c="dimmed">
              Most recent: {recentProjects[0]?.name}
              {recentProjects[0]?.isDirty && ' *'}
            </Text>
          ) : (
            <Text size="xs" c="dimmed">
              Manage {{ORG_NAME}} configuration projects
            </Text>
          )}
        </div>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onNewProject}>
            New Project
          </Button>
          {projects.length > 0 && (
            <Button size="xs" variant="outline" onClick={onViewAll}>
              View All
            </Button>
          )}
        </Group>
      </Stack>
    </DashboardCard>
  );
};
