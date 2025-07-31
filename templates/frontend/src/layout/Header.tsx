import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { Action, WalletConnectButton } from '@components';
import { useActiveProject } from '@hooks';
import { AppShell, Group, Text, useMantineColorScheme } from '@mantine/core';

export const Header = () => {
  const [appName, setAppName] = useState('AppName');
  const { colorScheme, setColorScheme } = useMantineColorScheme();
  const { toggleDarkMode, isDarkMode, toggleDebugMode, debugMode } =
    useActiveProject();

  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
    });
  }, []);

  useEffect(() => {
    setColorScheme(isDarkMode ? 'dark' : 'light');
  }, [isDarkMode, setColorScheme]);

  const handleToggleTheme = async () => {
    await toggleDarkMode();
  };

  const handleToggleDebug = async () => {
    await toggleDebugMode();
  };

  return (
    <AppShell.Header>
      <Group justify="space-between" p="md" h="100%">
        <Text size="xl" fw={700}>
          {appName}
        </Text>
        <Group justify="flex-end" align="center" gap="xs">
          <Action
            icon="DebugOn"
            iconOff="DebugOff"
            isOn={debugMode}
            onClick={handleToggleDebug}
            title={
              debugMode
                ? 'Debug mode ON - Click to disable'
                : 'Debug mode OFF - Click to enable'
            }
            variant="default"
            color={debugMode ? 'red' : 'gray'}
          />
          <Action
            icon="Light"
            iconOff="Dark"
            isOn={colorScheme === 'light'}
            onClick={handleToggleTheme}
            title="Toggle color scheme"
            variant="default"
            color={colorScheme === 'dark' ? 'yellow' : 'blue'}
          />
          <WalletConnectButton />
        </Group>
      </Group>
    </AppShell.Header>
  );
};
