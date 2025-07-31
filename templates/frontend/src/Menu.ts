// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
import { SetInitialized } from '@app';
import { DalleDress, Home, Khedra, Settings } from '@views';
import { Abis, Chunks, Contracts, Exports, Monitors, Names, Status } from '@views';
import { Wizard } from '@wizards';

export interface MenuItem {
  label: string;
  path: string;
  position: 'top' | 'bottom' | 'hidden';
  component?: React.ComponentType;
  hotkey?: string;
  altHotkey?: string;
  type?: 'navigation' | 'dev' | 'toggle';
  action?: () => void | Promise<void>;
}

export const MenuItems: MenuItem[] = [
  {
    label: 'Home',
    path: '/',
    position: 'top',
    component: Home,
    hotkey: 'mod+1',
    altHotkey: 'alt+1',
    type: 'navigation',
  },
  {
    label: 'Exports',
    path: '/exports',
    position: 'top',
    component: Exports,
    hotkey: 'mod+2',
    altHotkey: 'alt+2',
    type: 'navigation',
  },
  {
    label: 'Monitors',
    path: '/monitors',
    position: 'top',
    component: Monitors,
    hotkey: 'mod+3',
    altHotkey: 'alt+3',
    type: 'navigation',
  },
  {
    label: 'Abis',
    path: '/abis',
    position: 'top',
    component: Abis,
    hotkey: 'mod+4',
    altHotkey: 'alt+4',
    type: 'navigation',
  },
  {
    label: 'Names',
    path: '/names',
    position: 'top',
    component: Names,
    hotkey: 'mod+5',
    altHotkey: 'alt+5',
    type: 'navigation',
  },
  {
    label: 'Chunks',
    path: '/chunks',
    position: 'top',
    component: Chunks,
    hotkey: 'mod+6',
    altHotkey: 'alt+6',
    type: 'navigation',
  },
  {
    label: 'Contracts',
    path: '/contracts',
    position: 'top',
    component: Contracts,
    hotkey: 'mod+7',
    altHotkey: 'alt+7',
    type: 'navigation',
  },
  {
    label: 'Status',
    path: '/status',
    position: 'top',
    component: Status,
    hotkey: 'mod+8',
    altHotkey: 'alt+8',
    type: 'navigation',
  },
  {
    label: 'DalleDress',
    path: '/{{PROJECT_NAME}}',
    position: 'top',
    component: DalleDress,
    hotkey: 'mod+9',
    altHotkey: 'alt+9',
    type: 'navigation',
  },
  {
    label: 'Khedra',
    path: '/khedra',
    position: 'bottom',
    component: Khedra,
    hotkey: 'mod+0',
    altHotkey: 'alt+0',
    type: 'navigation',
  },
  {
    label: 'Settings',
    path: '/settings',
    position: 'bottom',
    component: Settings,
    hotkey: 'mod+shift+1',
    altHotkey: 'alt+shift+1',
    type: 'navigation',
  },
  {
    path: '/wizard',
    label: 'Wizard',
    position: 'hidden',
    component: Wizard,
    hotkey: 'mod+shift+2',
    altHotkey: 'alt+shift+2',
    type: 'dev',
    action: async () => {
      await SetInitialized(false);
    },
  },
];
