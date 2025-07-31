import { CancelFetch, Reload } from '@app';
import { useActiveProject } from '@hooks';
import { msgs, types } from '@models';
import { Log, emitEvent, registerHotkeys, useEmitters } from '@utils';
import type { HotkeyConfig, RegisterHotkeyOptions } from '@utils';
import { MenuItems } from 'src/Menu';
import { useLocation } from 'wouter';

interface BaseHotkey {
  type: 'navigation' | 'dev' | 'toggle';
  hotkey: string;
  label: string;
}

interface NavigationHotkey extends BaseHotkey {
  type: 'navigation';
  path: string;
}

interface DevHotkey extends BaseHotkey {
  type: 'dev';
  path: string;
  action?: () => Promise<void>;
}

interface ToggleHotkey extends BaseHotkey {
  type: 'toggle';
  action: () => void;
}

type Hotkey = NavigationHotkey | DevHotkey | ToggleHotkey;

export const useAppHotkeys = (): void => {
  const [currentLocation] = useLocation();
  const {
    lastTab,
    menuCollapsed,
    setMenuCollapsed,
    helpCollapsed,
    setHelpCollapsed,
  } = useActiveProject();
  const currentTab = lastTab[currentLocation] || '';

  const [, navigate] = useLocation();
  const { emitStatus, emitError } = useEmitters();

  const handleHotkey = async (
    hkType: Hotkey,
    _e: KeyboardEvent,
  ): Promise<void> => {
    try {
      switch (hkType.type) {
        case 'navigation':
          if (
            currentLocation === hkType.path ||
            currentLocation.startsWith(hkType.path + '/')
          ) {
            emitEvent(msgs.EventType.TAB_CYCLE, 'Tab cycle triggered', {
              route: hkType.path,
              key: hkType.hotkey,
            });
          } else {
            navigate(hkType.path);
          }
          break;

        case 'dev':
          if (!import.meta.env.DEV) return;
          if (hkType.action) await hkType.action();
          navigate(hkType.path);
          break;

        case 'toggle':
          hkType.action();
          break;
      }
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      Log(errorMessage);

      if (
        (hkType.type === 'navigation' || hkType.type === 'dev') &&
        hkType.path
      ) {
        window.location.href = hkType.path;
      }
    }
  };

  const toggleHotkeys: HotkeyConfig[] = [
    {
      key: 'mod+h',
      handler: (e: KeyboardEvent) =>
        handleHotkey(
          {
            type: 'toggle',
            hotkey: 'mod+h',
            label: 'Toggle help panel',
            action: () => {
              const next = !helpCollapsed;
              setHelpCollapsed(next);
            },
          },
          e,
        ),
      options: { preventDefault: true, enableOnFormTags: true },
    },
    {
      key: 'mod+m',
      handler: (e: KeyboardEvent) =>
        handleHotkey(
          {
            type: 'toggle',
            hotkey: 'mod+m',
            label: 'Toggle menu panel',
            action: () => {
              const next = !menuCollapsed;
              setMenuCollapsed(next);
            },
          },
          e,
        ),
      options: { preventDefault: true, enableOnFormTags: true },
    },
    {
      key: 'mod+r',
      handler: (e: KeyboardEvent) =>
        handleHotkey(
          {
            type: 'toggle',
            hotkey: 'mod+r',
            label: 'Reload',
            action: () => {
              Reload(types.Payload.createFrom({ dataFacet: currentTab })).then(
                () => {
                  // do nothing
                },
              );
            },
          },
          e,
        ),
      options: { preventDefault: true, enableOnFormTags: true },
    },
  ];

  const menuItemHotkeys = MenuItems.flatMap((item) => {
    const hotkeyConfigs: HotkeyConfig[] = [];

    if (item.hotkey) {
      hotkeyConfigs.push({
        key: item.hotkey,
        handler: (e: KeyboardEvent) => {
          let hotkeyObj: Hotkey;
          const hotkey = item.hotkey || '';

          switch (item.type) {
            case 'dev':
              hotkeyObj = {
                type: 'dev',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label}`,
                action: item.action as () => Promise<void>,
              };
              break;
            case 'toggle':
              hotkeyObj = {
                type: 'toggle',
                hotkey,
                label: `Toggle ${item.label}`,
                action: item.action as () => void,
              };
              break;
            case 'navigation':
            default:
              hotkeyObj = {
                type: 'navigation',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label}`,
              };
              break;
          }
          handleHotkey(hotkeyObj, e);
        },
        options: { preventDefault: true },
      });
    }

    if (item.altHotkey) {
      hotkeyConfigs.push({
        key: item.altHotkey,
        handler: (e: KeyboardEvent) => {
          let hotkeyObj: Hotkey;
          const hotkey = item.altHotkey || '';

          switch (item.type) {
            case 'dev':
              hotkeyObj = {
                type: 'dev',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label} (reverse)`,
                action: item.action as () => Promise<void>,
              };
              break;
            case 'toggle':
              hotkeyObj = {
                type: 'toggle',
                hotkey,
                label: `Toggle ${item.label} (reverse)`,
                action: item.action as () => void,
              };
              break;
            case 'navigation':
            default:
              hotkeyObj = {
                type: 'navigation',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label} (reverse)`,
              };
              break;
          }
          handleHotkey(hotkeyObj, e);
        },
        options: { preventDefault: true },
      });
    }

    return hotkeyConfigs;
  });

  const editHotkeys: HotkeyConfig[] = [
    {
      key: 'mod+c',
      handler: (_e: KeyboardEvent) => {
        // Allow default browser action
      },
      options: {
        preventDefault: false,
        enableOnFormTags: true,
      } as RegisterHotkeyOptions,
    },
    {
      key: 'mod+v',
      handler: (_e: KeyboardEvent) => {
        // Allow default browser action
      },
      options: {
        preventDefault: false,
        enableOnFormTags: true,
      } as RegisterHotkeyOptions,
    },
    {
      key: 'mod+x',
      handler: (_e: KeyboardEvent) => {
        // Allow default browser action
      },
      options: {
        preventDefault: false,
        enableOnFormTags: true,
      } as RegisterHotkeyOptions,
    },
    {
      key: 'mod+z',
      handler: (_e: KeyboardEvent) => {
        // Allow default browser action
      },
      options: {
        preventDefault: false,
        enableOnFormTags: true,
      } as RegisterHotkeyOptions,
    },
    {
      key: 'mod+shift+z',
      handler: (_e: KeyboardEvent) => {
        // Allow default browser action
      },
      options: {
        preventDefault: false,
        enableOnFormTags: true,
      } as RegisterHotkeyOptions,
    },
    {
      key: 'mod+y', // Redo (alternative to mod+shift+z)
      handler: (_e: KeyboardEvent) => {
        // Allow default browser action
      },
      options: {
        preventDefault: false,
        enableOnFormTags: true,
      } as RegisterHotkeyOptions,
    },
  ];

  const globalHotkeys: HotkeyConfig[] = [
    {
      key: 'escape',
      handler: (_e: KeyboardEvent) => {
        CancelFetch(currentTab as types.DataFacet)
          .then(() => {
            emitStatus('Cancellation request processed.');
          })
          .catch((err: Error) => {
            emitError(
              `Failed to send cancellation request via Escape key: ${
                err.message || 'Unknown error'
              }`,
            );
          });
      },
      options: {
        enableOnFormTags: true,
        preventDefault: true,
      } as RegisterHotkeyOptions,
    },
  ];

  const hotkeysConfig: HotkeyConfig[] = [
    ...menuItemHotkeys,
    ...toggleHotkeys,
    ...editHotkeys,
    ...globalHotkeys,
  ];
  registerHotkeys(hotkeysConfig);
};
