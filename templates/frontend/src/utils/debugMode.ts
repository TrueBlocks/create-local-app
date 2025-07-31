import { GetDebugMode } from '@app';
import { appPreferencesStore } from '@stores';

export const isDebugMode = (): boolean => {
  return appPreferencesStore.getState().debugMode;
};

export const getDebugModeAsync = async (): Promise<boolean> => {
  try {
    return await GetDebugMode();
  } catch {
    return false;
  }
};
