import { useEffect, useState } from 'react';

import { CloseProject, GetOpenProjects, SwitchToProject } from '@app';
import { Action } from '@components';
import { useEvent, useIconSets } from '@hooks';
import { Group, List, Paper, Text, ThemeIcon, Title } from '@mantine/core';
import { msgs } from '@models';
import { Log } from '@utils';

// TODO: Use Project from project model
interface Project {
  id: string;
  name: string;
  path: string;
  isActive: boolean;
  isDirty: boolean;
  lastOpened: string;
  createdAt: string;
}

export const ProjectsList = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const { File } = useIconSets();

  const refreshProjects = async () => {
    try {
      const openProjects = await GetOpenProjects();
      const typedProjects = openProjects as unknown as Project[];

      const sortedProjects = typedProjects.sort((a, b) => {
        const lastOpenedDiff =
          new Date(b.lastOpened).getTime() - new Date(a.lastOpened).getTime();
        if (lastOpenedDiff !== 0) {
          return lastOpenedDiff;
        }

        return (
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
        );
      });

      setProjects(sortedProjects);
    } catch (error) {
      Log(`Error fetching projects: ${error}`);
    }
  };

  useEffect(() => {
    refreshProjects();
  }, []);

  useEvent(msgs.EventType.MANAGER, (_message?: string) => refreshProjects());

  const handleSwitchProject = async (id: string) => {
    try {
      await SwitchToProject(id);
      refreshProjects();
    } catch (error) {
      Log(`Error switching projects: ${error}`);
    }
  };

  const handleCloseProject = async (id: string, isDirty: boolean) => {
    if (isDirty) {
      // You might want to show a confirmation dialog here
      const confirmed = window.confirm(
        'This project has unsaved changes. Close anyway?',
      );
      if (!confirmed) return;
    }

    try {
      await CloseProject(id);
      refreshProjects();
    } catch (error) {
      Log(`Error closing project: ${error}`);
    }
  };

  if (projects.length === 0) {
    return null; // Don't show anything if there are no open projects
  }

  return (
    <Paper p="md" withBorder radius="md" mb="md">
      <Title order={4} mb="sm">
        Open Projects
      </Title>
      <List spacing="xs">
        {projects.map((project) => (
          <List.Item
            key={project.id}
            icon={
              <ThemeIcon
                color={project.isActive ? 'blue' : 'gray'}
                size={24}
                radius="xl"
              >
                <File size={16} />
              </ThemeIcon>
            }
          >
            <Group justify="space-between" wrap="nowrap">
              <div>
                <Text fw={project.isActive ? 700 : 400}>
                  {project.name} {project.isDirty && '*'}
                </Text>
                <Text size="xs" c="dimmed">
                  {project.path}
                </Text>
              </div>
              <Group gap="xs">
                {!project.isActive && (
                  <Action
                    icon="Switch"
                    title="Switch to this project"
                    color="blue"
                    variant="light"
                    onClick={() => handleSwitchProject(project.id)}
                  />
                )}
                <Action
                  icon="Delete"
                  title="Close project"
                  color="red"
                  variant="light"
                  onClick={() =>
                    handleCloseProject(project.id, project.isDirty)
                  }
                />
              </Group>
            </Group>
          </List.Item>
        ))}
      </List>
    </Paper>
  );
};
